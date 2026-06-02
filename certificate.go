package main

import (
	"context"
	"crypto/sha1"
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	oidAppleDev         = asn1.ObjectIdentifier{1, 2, 840, 113635, 100, 6, 1, 2}
	oidAppleDist        = asn1.ObjectIdentifier{1, 2, 840, 113635, 100, 6, 1, 4}
	oidMacDev           = asn1.ObjectIdentifier{1, 2, 840, 113635, 100, 6, 1, 12}
	oidMacAppDev        = asn1.ObjectIdentifier{1, 2, 840, 113635, 100, 6, 1, 7}
	oidDevIDApplication = asn1.ObjectIdentifier{1, 2, 840, 113635, 100, 6, 1, 13}
)

var markerOIDs = []struct {
	OID      asn1.ObjectIdentifier
	CertType string
}{
	{oidAppleDev, "iOS Development"},
	{oidAppleDist, "iOS Distribution"},
	{oidMacDev, "Mac Development"},
	{oidMacAppDev, "Mac App Development"},
	{oidDevIDApplication, "Developer ID Application"},
}

var cnPrefixes = []struct {
	Prefix   string
	CertType string
}{
	{"Apple Development:", "Development"},
	{"Apple Distribution:", "Distribution"},
	{"iPhone Developer:", "iOS Development"},
	{"iPhone Distribution:", "iOS Distribution"},
	{"3rd Party Mac Developer", "Mac Distribution"},
}

var reasonMapping = map[string]string{
	"keyCompromise":        "Key Compromise",
	"caCompromise":         "CA Compromise",
	"affiliationChanged":   "Affiliation Changed",
	"superseded":           "Superseded",
	"cessationOfOperation": "Cessation of Operation",
	"certificateHold":      "Certificate Hold",
	"removeFromCRL":        "Remove from CRL",
	"privilegeWithdrawn":   "Privilege Withdrawn",
	"aACompromise":         "AA Compromise",
	"unspecified":          "Unspecified",
}

const opensslTimeout = 15 * time.Second

func GetCertificateInfo(plistData map[string]any) CertInfo {
	rawCerts := GetRawCertificates(plistData)
	if len(rawCerts) == 0 {
		return CertInfo{Name: "Unknown"}
	}

	cert, err := x509.ParseCertificate(rawCerts[0])
	if err != nil {
		return CertInfo{Name: "Unknown"}
	}

	info := CertInfo{
		Name:         cert.Subject.CommonName,
		SerialNumber: cert.SerialNumber.Text(16),
		CertType:     detectCertType(cert),
		SHA1:         formatSHA1(cert.Raw),
		NotBefore:    cert.NotBefore.Format("2006-01-02 15:04:05"),
		Expiration:   cert.NotAfter.Format("2006-01-02 15:04:05"),
		Expired:      cert.NotAfter.Before(time.Now()),
	}

	if len(cert.Subject.Organization) > 0 {
		info.Organization = cert.Subject.Organization[0]
	}

	if len(cert.Subject.OrganizationalUnit) > 0 {
		info.OrgUnit = cert.Subject.OrganizationalUnit[0]
	}

	info.Revoked, info.RevokeTime, info.RevokeReason = checkOCSPRevocation(rawCerts[0])

	return info
}

func detectCertType(cert *x509.Certificate) string {
	for _, m := range markerOIDs {
		for _, ext := range cert.Extensions {
			if ext.Id.Equal(m.OID) {
				return m.CertType
			}
		}
	}

	cn := cert.Subject.CommonName
	for _, p := range cnPrefixes {
		if strings.HasPrefix(cn, p.Prefix) {
			return p.CertType
		}
	}

	return ""
}

func formatSHA1(raw []byte) string {
	sum := sha1.Sum(raw)
	parts := make([]string, len(sum))
	for i, b := range sum {
		parts[i] = fmt.Sprintf("%02X", b)
	}
	return strings.Join(parts, ":")
}

func writeTempFile(pattern string, data []byte) (string, error) {
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	name := f.Name()
	_, writeErr := f.Write(data)
	closeErr := f.Close()
	if writeErr != nil {
		os.Remove(name)
		return "", writeErr
	}
	if closeErr != nil {
		os.Remove(name)
		return "", closeErr
	}
	return name, nil
}

func checkOCSPRevocation(certDER []byte) (bool, string, string) {
	certPath, err := writeTempFile("cert-*.der", certDER)
	if err != nil {
		return false, "", ""
	}
	defer os.Remove(certPath)

	ctx, cancel := context.WithTimeout(context.Background(), opensslTimeout)
	defer cancel()

	uriCmd := exec.CommandContext(ctx, "openssl", "x509", "-inform", "DER", "-in", certPath, "-ocsp_uri", "-noout")
	uriOut, err := uriCmd.Output()
	ocspURI := strings.TrimSpace(string(uriOut))
	if err != nil || ocspURI == "" {
		ocspURI = "http://ocsp.apple.com/ocsp03-wwdr01"
	}

	issuerDER, err := fetchIssuerCert()
	if err != nil {
		return false, "", ""
	}

	issuerPath, err := writeTempFile("issuer-*.cer", issuerDER)
	if err != nil {
		return false, "", ""
	}
	defer os.Remove(issuerPath)

	ocspCtx, ocspCancel := context.WithTimeout(context.Background(), opensslTimeout)
	defer ocspCancel()

	ocspCmd := exec.CommandContext(ocspCtx, "openssl", "ocsp",
		"-issuer", issuerPath,
		"-cert", certPath,
		"-url", ocspURI,
		"-resp_text",
	)

	ocspOut, ocspErr := ocspCmd.CombinedOutput()
	output := string(ocspOut)

	if ocspErr != nil && !strings.Contains(strings.ToLower(output), "revoked") {
		return false, "", ""
	}

	if !strings.Contains(strings.ToLower(output), "revoked") {
		return false, "", ""
	}

	var revokeTime, revokeReason string

	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)

		if strings.Contains(trimmed, "Revocation Time:") {
			parts := strings.SplitN(trimmed, "Revocation Time:", 2)
			if len(parts) == 2 {
				revokeTime = strings.TrimSpace(parts[1])
			}
		}

		if strings.Contains(trimmed, "Reason:") {
			parts := strings.SplitN(trimmed, "Reason:", 2)
			if len(parts) == 2 {
				raw := strings.TrimSpace(parts[1])
				if mapped, ok := reasonMapping[raw]; ok {
					revokeReason = mapped
				} else {
					revokeReason = raw
				}
			}
		}
	}

	return true, revokeTime, revokeReason
}

func fetchIssuerCert() ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://www.apple.com/certificateauthority/AppleWWDRCAG3.cer")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch issuer certificate: HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}

	return data, nil
}
