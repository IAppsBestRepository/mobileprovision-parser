package main

import (
	"testing"
	"time"
)

func TestGetProfileType_DevelopmentEnterprise(t *testing.T) {
	plistData := map[string]any{
		"ProvisionsAllDevices": true,
		"Entitlements": map[string]any{
			"get-task-allow": true,
		},
	}
	result := GetProfileType(plistData)
	if result != "Development (Enterprise)" {
		t.Fatalf("expected Development (Enterprise), got %s", result)
	}
}

func TestGetProfileType_DistributionInHouse(t *testing.T) {
	plistData := map[string]any{
		"ProvisionsAllDevices": true,
		"Entitlements": map[string]any{
			"get-task-allow": false,
		},
	}
	result := GetProfileType(plistData)
	if result != "Distribution (In-House)" {
		t.Fatalf("expected Distribution (In-House), got %s", result)
	}
}

func TestGetProfileType_Development(t *testing.T) {
	plistData := map[string]any{
		"ProvisionedDevices": []any{"AAAA-BBBB-CCCC"},
		"Entitlements": map[string]any{
			"get-task-allow": true,
		},
	}
	result := GetProfileType(plistData)
	if result != "Development" {
		t.Fatalf("expected Development, got %s", result)
	}
}

func TestGetProfileType_AdHoc(t *testing.T) {
	plistData := map[string]any{
		"ProvisionedDevices": []any{"AAAA-BBBB-CCCC"},
		"Entitlements": map[string]any{
			"get-task-allow": false,
		},
	}
	result := GetProfileType(plistData)
	if result != "Ad Hoc" {
		t.Fatalf("expected Ad Hoc, got %s", result)
	}
}

func TestGetProfileType_AppStore(t *testing.T) {
	plistData := map[string]any{
		"Entitlements": map[string]any{
			"get-task-allow": false,
		},
	}
	result := GetProfileType(plistData)
	if result != "App Store" {
		t.Fatalf("expected App Store, got %s", result)
	}
}

func TestCheckValidity_Expired(t *testing.T) {
	plistData := map[string]any{
		"ExpirationDate": time.Now().Add(-time.Hour),
	}
	cert := CertInfo{}
	result := CheckValidity(plistData, cert)
	if result.Status != "Expired" || result.OK != false {
		t.Fatalf("expected Expired/false, got %s/%v", result.Status, result.OK)
	}
}

func TestCheckValidity_ExpiringSoon(t *testing.T) {
	plistData := map[string]any{
		"ExpirationDate": time.Now().Add(time.Hour * 24 * 15),
	}
	cert := CertInfo{}
	result := CheckValidity(plistData, cert)
	if result.Status != "Expiring Soon" {
		t.Fatalf("expected Expiring Soon, got %s", result.Status)
	}
	if result.OK != true {
		t.Fatalf("expected OK=true for Expiring Soon, got %v", result.OK)
	}
	if result.Detail == "" {
		t.Fatalf("expected non-empty Detail for Expiring Soon")
	}
}

func TestCheckValidity_CertificateExpired(t *testing.T) {
	plistData := map[string]any{
		"ExpirationDate": time.Now().Add(time.Hour * 24 * 365),
	}
	cert := CertInfo{Expired: true}
	result := CheckValidity(plistData, cert)
	if result.Status != "Certificate Expired" || result.OK != false {
		t.Fatalf("expected Certificate Expired/false, got %s/%v", result.Status, result.OK)
	}
}

func TestCheckValidity_Revoked(t *testing.T) {
	plistData := map[string]any{
		"ExpirationDate": time.Now().Add(time.Hour * 24 * 365),
	}
	cert := CertInfo{
		Revoked:      true,
		RevokeTime:   "2025-01-01",
		RevokeReason: "Key Compromise",
	}
	result := CheckValidity(plistData, cert)
	if result.Status != "Revoked" || result.OK != false {
		t.Fatalf("expected Revoked/false, got %s/%v", result.Status, result.OK)
	}
	expectedDetail := "2025-01-01 Key Compromise"
	if result.Detail != expectedDetail {
		t.Fatalf("expected detail %q, got %q", expectedDetail, result.Detail)
	}
}

func TestCheckValidity_Signed(t *testing.T) {
	plistData := map[string]any{
		"ExpirationDate": time.Now().Add(time.Hour * 24 * 365),
	}
	cert := CertInfo{}
	result := CheckValidity(plistData, cert)
	if result.Status != "Signed" || result.OK != true {
		t.Fatalf("expected Signed/true, got %s/%v", result.Status, result.OK)
	}
}

func TestGetDevices_WithDevices(t *testing.T) {
	plistData := map[string]any{
		"ProvisionedDevices": []any{"device1", "device2", "device3"},
	}
	devices := GetDevices(plistData)
	if len(devices) != 3 {
		t.Fatalf("expected 3 devices, got %d", len(devices))
	}
	if devices[0] != "device1" || devices[1] != "device2" || devices[2] != "device3" {
		t.Fatalf("unexpected devices: %v", devices)
	}
}

func TestGetDevices_WithoutDevices(t *testing.T) {
	plistData := map[string]any{}
	devices := GetDevices(plistData)
	if devices != nil {
		t.Fatalf("expected nil, got %v", devices)
	}
}

func TestGetDevices_Empty(t *testing.T) {
	plistData := map[string]any{
		"ProvisionedDevices": []any{},
	}
	devices := GetDevices(plistData)
	if len(devices) != 0 {
		t.Fatalf("expected 0 devices, got %d", len(devices))
	}
}

func TestNormalizePlatforms_WithPlatforms(t *testing.T) {
	plistData := map[string]any{
		"Platform": []any{"iOS", "macOS"},
	}
	platforms := NormalizePlatforms(plistData)
	if len(platforms) != 2 {
		t.Fatalf("expected 2 platforms, got %d", len(platforms))
	}
	if platforms[0] != "iOS" || platforms[1] != "macOS" {
		t.Fatalf("unexpected platforms: %v", platforms)
	}
}

func TestNormalizePlatforms_WithoutPlatforms(t *testing.T) {
	plistData := map[string]any{}
	platforms := NormalizePlatforms(plistData)
	if platforms != nil {
		t.Fatalf("expected nil, got %v", platforms)
	}
}

func TestNormalizePlatforms_XrOSToVisionOS(t *testing.T) {
	plistData := map[string]any{
		"Platform": []any{"iOS", "xrOS"},
	}
	platforms := NormalizePlatforms(plistData)
	if len(platforms) != 2 {
		t.Fatalf("expected 2 platforms, got %d", len(platforms))
	}
	if platforms[0] != "iOS" || platforms[1] != "visionOS" {
		t.Fatalf("expected [iOS visionOS], got %v", platforms)
	}
}

func TestNormalizePlatforms_OSXToMacOS(t *testing.T) {
	plistData := map[string]any{
		"Platform": []any{"OSX"},
	}
	platforms := NormalizePlatforms(plistData)
	if len(platforms) != 1 {
		t.Fatalf("expected 1 platform, got %d", len(platforms))
	}
	if platforms[0] != "macOS" {
		t.Fatalf("expected [macOS], got %v", platforms)
	}
}

func TestFormatDate(t *testing.T) {
	dt := time.Date(2025, 6, 15, 10, 30, 45, 0, time.UTC)
	result := formatDate(dt)
	expected := "2025-06-15 10:30:45"
	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}
