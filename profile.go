package main

import (
	"math"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

var profileTypeMap = map[string]string{
	"IOS_APP_DEVELOPMENT": "Development",
	"IOS_APP_STORE":       "App Store",
	"IOS_APP_ADHOC":       "Ad Hoc",
	"IOS_APP_INHOUSE":     "Distribution (In-House)",
	"TVOS_APP_DEVELOPMENT": "Development",
	"TVOS_APP_STORE":       "App Store",
	"TVOS_APP_ADHOC":       "Ad Hoc",
	"TVOS_APP_INHOUSE":     "Distribution (In-House)",
	"MAC_APP_DEVELOPMENT":  "Development",
	"MAC_APP_STORE":        "App Store",
	"MAC_APP_DIRECT":       "Developer ID",
	"MAC_APP_INHOUSE":      "Distribution (In-House)",
}

var platformNormalize = map[string]string{
	"xrOS": "visionOS",
	"OSX":  "macOS",
}

func GetProfileType(plistData map[string]any) string {
	if pt, ok := plistData["ProfileType"].(string); ok {
		if mapped, found := profileTypeMap[pt]; found {
			return mapped
		}
	}

	provisionsAll, _ := plistData["ProvisionsAllDevices"].(bool)

	var hasDevices bool
	if devs, ok := plistData["ProvisionedDevices"].([]any); ok && len(devs) > 0 {
		hasDevices = true
	}

	var getTaskAllow bool
	if ents, ok := plistData["Entitlements"].(map[string]any); ok {
		getTaskAllow, _ = ents["get-task-allow"].(bool)
	}

	switch {
	case provisionsAll && getTaskAllow:
		return "Development (Enterprise)"
	case provisionsAll && !getTaskAllow:
		return "Distribution (In-House)"
	case hasDevices && getTaskAllow:
		return "Development"
	case hasDevices && !getTaskAllow:
		return "Ad Hoc"
	case !getTaskAllow && !hasDevices:
		return "App Store"
	default:
		return "Unknown"
	}
}

func CheckValidity(plistData map[string]any, cert CertInfo) ValidityInfo {
	now := time.Now()

	if exp, ok := plistData["ExpirationDate"].(time.Time); ok {
		if exp.Before(now) {
			return ValidityInfo{Status: "Expired", OK: false}
		}
		days := int(math.Ceil(exp.Sub(now).Hours() / 24))
		if days <= 30 {
			return ValidityInfo{
				Status: "Expiring Soon",
				OK:     true,
				Detail: strconv.Itoa(days) + " days remaining",
			}
		}
	}

	if cert.Expired {
		return ValidityInfo{Status: "Certificate Expired", OK: false}
	}

	if cert.Revoked {
		return ValidityInfo{Status: "Revoked", OK: false, Detail: cert.RevokeTime + " " + cert.RevokeReason}
	}

	return ValidityInfo{Status: "Signed", OK: true}
}

func ParseEntitlements(plistData map[string]any) []EntitlementEntry {
	ents, ok := plistData["Entitlements"].(map[string]any)
	if !ok {
		return nil
	}

	entries := make([]EntitlementEntry, 0, len(ents))
	for key, value := range ents {
		status := GetEntitlementStatus(key, value)
		if status == StatusDisabled {
			continue
		}
		entries = append(entries, EntitlementEntry{
			Key:     key,
			Name:    GetEntitlementName(key),
			Display: FormatEntitlementValue(value),
			Status:  status,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Status != entries[j].Status {
			return entries[i].Status < entries[j].Status
		}
		return entries[i].Name < entries[j].Name
	})

	return entries
}

func NormalizePlatforms(plistData map[string]any) []string {
	plats, ok := plistData["Platform"].([]any)
	if !ok {
		return nil
	}

	seen := make(map[string]bool, len(plats))
	platforms := make([]string, 0, len(plats))
	for _, p := range plats {
		s, ok := p.(string)
		if !ok {
			continue
		}
		if normalized, found := platformNormalize[s]; found {
			s = normalized
		}
		if !seen[s] {
			seen[s] = true
			platforms = append(platforms, s)
		}
	}
	sort.Strings(platforms)
	return platforms
}

func GetDevices(plistData map[string]any) []string {
	devs, ok := plistData["ProvisionedDevices"].([]any)
	if !ok {
		return nil
	}

	devices := make([]string, 0, len(devs))
	for _, d := range devs {
		if s, ok := d.(string); ok {
			devices = append(devices, s)
		}
	}
	return devices
}

func plistInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case uint64:
		return int(n)
	case float64:
		return int(n)
	default:
		return 0
	}
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func ParseProfile(filePath string) (*ProfileInfo, error) {
	plistData, err := ExtractPlist(filePath)
	if err != nil {
		return nil, err
	}

	cert := GetCertificateInfo(plistData)
	validity := CheckValidity(plistData, cert)
	entitlements := ParseEntitlements(plistData)
	devices := GetDevices(plistData)
	platforms := NormalizePlatforms(plistData)
	profileType := GetProfileType(plistData)

	profileName, _ := plistData["Name"].(string)
	appIDName, _ := plistData["AppIDName"].(string)
	uuid, _ := plistData["UUID"].(string)
	teamName, _ := plistData["TeamName"].(string)
	isXcodeManaged, _ := plistData["IsXcodeManaged"].(bool)

	teamID := ""
	if teams, ok := plistData["TeamIdentifier"].([]any); ok && len(teams) > 0 {
		teamID, _ = teams[0].(string)
	}

	appID := ""
	if ents, ok := plistData["Entitlements"].(map[string]any); ok {
		if id, ok := ents["application-identifier"].(string); ok {
			appID = id
		} else if id, ok := ents["com.apple.application-identifier"].(string); ok {
			appID = id
		}
	}

	ttl := plistInt(plistData["TimeToLive"])
	version := plistInt(plistData["Version"])

	creationDate := ""
	if t, ok := plistData["CreationDate"].(time.Time); ok {
		creationDate = formatDate(t)
	}

	expirationDate := ""
	daysRemaining := 0
	if t, ok := plistData["ExpirationDate"].(time.Time); ok {
		expirationDate = formatDate(t)
		remaining := int(math.Ceil(t.Sub(time.Now()).Hours() / 24))
		if remaining > 0 {
			daysRemaining = remaining
		}
	}

	ppqRequired, _ := plistData["PPQCheck"].(bool)

	return &ProfileInfo{
		FileName:       filepath.Base(filePath),
		ProfileName:    profileName,
		AppIDName:      appIDName,
		AppID:          appID,
		TeamID:         teamID,
		TeamName:       teamName,
		UUID:           uuid,
		ProfileType:    profileType,
		Platforms:      platforms,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
		DaysRemaining:  daysRemaining,
		TimeToLive:     ttl,
		Version:        version,
		IsXcodeManaged: isXcodeManaged,
		PPQRequired:    ppqRequired,
		Certificate:    cert,
		Validity:       validity,
		Entitlements:   entitlements,
		Devices:        devices,
	}, nil
}
