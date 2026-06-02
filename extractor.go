package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"howett.net/plist"
)

const (
	securityCmdTimeout = 30 * time.Second
	maxProvisionSize   = 50 * 1024 * 1024
)

var validExtensions = map[string]struct{}{
	".mobileprovision": {},
	".provisionprofile": {},
}

func ExtractPlist(filePath string) (map[string]any, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("provisioning profile not found: %s", filePath)
		}
		if errors.Is(err, os.ErrPermission) {
			return nil, fmt.Errorf("permission denied reading provisioning profile: %s", filePath)
		}
		return nil, fmt.Errorf("cannot access provisioning profile %s: %w", filePath, err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("expected a file but got a directory: %s", filePath)
	}

	if info.Size() == 0 {
		return nil, fmt.Errorf("provisioning profile is empty: %s", filePath)
	}

	if info.Size() > maxProvisionSize {
		return nil, fmt.Errorf("provisioning profile too large (%d bytes, max %d): %s", info.Size(), maxProvisionSize, filePath)
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	if _, ok := validExtensions[ext]; !ok {
		return nil, fmt.Errorf("unsupported file type %q (expected .mobileprovision or .provisionprofile): %s", ext, filePath)
	}

	ctx, cancel := context.WithTimeout(context.Background(), securityCmdTimeout)
	defer cancel()

	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "security", "cms", "-D", "-i", filePath)
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("timed out after %s extracting plist from %s", securityCmdTimeout, filePath)
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			errMsg := strings.TrimSpace(stderr.String())
			if errMsg != "" {
				return nil, fmt.Errorf("security cms failed (exit %d) for %s: %s", exitErr.ExitCode(), filePath, errMsg)
			}
			return nil, fmt.Errorf("security cms failed (exit %d) for %s", exitErr.ExitCode(), filePath)
		}
		return nil, fmt.Errorf("failed to run security cms for %s: %w", filePath, err)
	}

	if len(output) == 0 {
		return nil, fmt.Errorf("security cms returned empty output for %s", filePath)
	}

	var result map[string]any
	if _, err := plist.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse plist from %s: %w", filePath, err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("provisioning profile contains no data: %s", filePath)
	}

	return result, nil
}

func GetRawCertificates(plistData map[string]any) [][]byte {
	if plistData == nil {
		return nil
	}

	certs, ok := plistData["DeveloperCertificates"]
	if !ok {
		return nil
	}

	certArray, ok := certs.([]any)
	if !ok {
		return nil
	}

	rawCerts := make([][]byte, 0, len(certArray))
	for _, cert := range certArray {
		if data, ok := cert.([]byte); ok && len(data) > 0 {
			rawCerts = append(rawCerts, data)
		}
	}

	if len(rawCerts) == 0 {
		return nil
	}

	return rawCerts
}
