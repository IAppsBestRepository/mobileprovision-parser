<div align="center">

# MobileProvision Parser

**A fast, beautiful CLI tool for inspecting Apple provisioning profiles**

[![Go 1.22+](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![macOS](https://img.shields.io/badge/platform-macOS-lightgrey?logo=apple)](https://developer.apple.com)

</div>

---

## Overview

MobileProvision Parser is a command-line tool that decodes `.mobileprovision` and `.provisionprofile` files into a richly formatted terminal display. It extracts every important detail -- certificates, entitlements, device lists, profile validity, and OCSP revocation status -- and presents them with color-coded lipgloss styling so you can assess a profile's health at a glance. No more piping through `security cms` and squinting at raw XML.

<!-- screenshot -->

## Features

- **Full profile parsing** -- reads `.mobileprovision` and `.provisionprofile` files via macOS `security cms`
- **Certificate deep-dive** -- extracts Common Name, certificate type (via Apple OID markers), serial number, organization, SHA-1 fingerprint, and validity dates
- **Live OCSP revocation checking** -- queries Apple's OCSP responder in real time to detect revoked certificates, with revocation time and reason
- **PPQCheck detection** -- reads the `PPQCheck` key directly from the profile plist to flag App Store review risk
- **Smart profile type detection** -- uses the explicit `ProfileType` key when available, with heuristic fallback based on `ProvisionsAllDevices`, `ProvisionedDevices`, and `get-task-allow`
- **175 entitlement mappings** -- translates raw entitlement keys (e.g. `com.apple.developer.healthkit`) into human-readable names with 3-color status indicators
- **Expiring Soon warning** -- highlights profiles expiring within 30 days
- **Platform normalization** -- displays `xrOS` as `visionOS` and `OSX` as `macOS`
- **Metadata extraction** -- TeamName, TeamID, TimeToLive, Version, IsXcodeManaged, UUID, and more
- **Device list summary** -- shows provisioned device UDIDs with overflow count
- **Beautiful terminal output** -- styled with [lipgloss](https://github.com/charmbracelet/lipgloss) using a purple/indigo color scheme

## Quick Start

```bash
# Clone the repository
git clone https://github.com/IAppsBestRepository/mobileprovision-parser.git
cd mobileprovision-parser

# Build
go build -o mobileprovision-parser .

# Parse a provisioning profile
./mobileprovision-parser ~/Library/MobileDevice/Provisioning\ Profiles/MyApp.mobileprovision
```

## Output Example

```
╭──────────────────────────────────────────────────╮
│  ✦ MobileProvision Parser                        │
╰──────────────────────────────────────────────────╯

  Profile             MyApp_Dev.mobileprovision
  Name                MyApp Development
  App ID Name         MyApp
  App ID              A1B2C3D4E5.com.example.myapp
  Team                A1B2C3D4E5 (Example Corp)
  UUID                8f3a9b12-4c67-4e8a-b5d1-9e2f0a3b7c84
  Type                Development
  Platform            iOS
  Created             2026-03-15 09:30:00
  Expires             2027-03-15 09:30:00  (286 days remaining)
  TTL                 365 days
  Version             1

  Certificate
  ──────────────────────────────────────────────────
    Name              Apple Development: John Doe (X9Y8Z7W6V5)
    Type              iOS Development
    Serial            0a1b2c3d4e5f
    Organization      Example Corp
    SHA-1             AB:CD:EF:01:23:45:67:89:AB:CD:EF:01:23:45:67:89:AB:CD:EF:01
    Valid             2026-03-10 00:00:00 — 2027-03-10 23:59:59
    Status            ✓ Signed

  Profile Status      ✓ Signed
  PPQ Check           Disabled

  Entitlements (8)
  ──────────────────────────────────────────────────
    ✓ Application Identifier  (A1B2C3D4E5.com.example.myapp)
    ✓ Push Notifications (iOS)  (development)
    ✓ Sign in with Apple
    ✓ HealthKit
    ✓ Get Task Allow (Debuggable)
    ● App Groups  (group.com.example.myapp.*)
    ● Keychain Access Groups  (A1B2C3D4E5.*)
    ● Associated Domains  (*)
```

## What It Shows

| Field | Description |
|---|---|
| **Profile** | Source filename |
| **Name** | Profile display name (`Name` key) |
| **App ID Name** | Registered App ID name in the Developer Portal |
| **App ID** | Full application identifier (team prefix + bundle ID) |
| **Team** | Team Identifier and Team Name |
| **UUID** | Profile UUID |
| **Type** | Development, Ad Hoc, App Store, Distribution (In-House), Developer ID, or Development (Enterprise) |
| **Platform** | Target platforms, normalized (`xrOS` shown as `visionOS`, `OSX` as `macOS`) |
| **Xcode Managed** | Shown only for automatically managed profiles |
| **Created / Expires** | Profile creation and expiration timestamps with days-remaining counter |
| **TTL** | `TimeToLive` value in days |
| **Version** | Profile version number |
| **Certificate** | CN, type, serial, organization, org unit, SHA-1, validity range, and live OCSP status |
| **Profile Status** | Overall validity: Signed, Expiring Soon, Expired, Certificate Expired, or Revoked |
| **PPQ Check** | Whether Apple's PPQCheck flag is set in the profile |
| **Devices** | Provisioned device UDIDs (first 5 shown, remainder counted) |
| **Entitlements** | All entitlements with human-readable names and 3-color status |

## PPQ Check

**PPQ (Provisioning Profile Query)** is a flag (`PPQCheck`) embedded in provisioning profiles by Apple. When enabled, it signals that the app will undergo additional automated screening during App Store review. PPQ checks are typically associated with stricter compliance requirements.

This tool reads the `PPQCheck` boolean directly from the profile plist and displays it as:

- **Disabled** (green) -- standard review process
- **Enabled** (red) -- additional PPQ screening will apply

## Entitlement Status

Entitlements are displayed with three distinct status indicators:

| Symbol | Color | Meaning |
|:---:|---|---|
| **✓** | Green | **Enabled** -- the entitlement is active with a concrete, non-wildcard value |
| **●** | Amber | **Configured** -- the entitlement is present but uses a wildcard or placeholder value (e.g. `*`, `com.example.*`, or an empty array) and requires further configuration in Xcode |
| **✗** | Red | **Disabled** -- the entitlement is set to `false`, `nil`, or empty (these are hidden from output by default) |

The distinction between Enabled and Configured is applied to specific entitlements known to accept wildcard values, such as App Groups, Associated Domains, Keychain Access Groups, iCloud containers, and others.

## Architecture

```
.
├── main.go              Entry point, CLI argument handling
├── types.go             Data structures (ProfileInfo, CertInfo, EntitlementEntry, ValidityInfo)
├── extractor.go         Plist extraction via `security cms`, file validation, DER certificate extraction
├── certificate.go       X.509 parsing, Apple OID-based cert type detection, SHA-1 fingerprinting, OCSP revocation checking
├── profile.go           Profile type detection, validity checking, entitlement parsing, platform normalization, device extraction
├── entitlements.go      175 entitlement key-to-name mappings, status classification logic, value formatting
├── display.go           Terminal output rendering with structured layout
├── styles.go            Lipgloss style definitions (colors, formatting)
├── profile_test.go      Tests for profile type detection, validity checks, platform normalization, device parsing
└── entitlements_test.go Tests for entitlement name resolution, value formatting, status classification
```

## Certificate Type Detection

The tool identifies certificate types using two strategies:

1. **Apple OID markers** -- checks X.509 extension OIDs specific to Apple certificate types:
   - `1.2.840.113635.100.6.1.2` -- iOS Development
   - `1.2.840.113635.100.6.1.4` -- iOS Distribution
   - `1.2.840.113635.100.6.1.12` -- Mac Development
   - `1.2.840.113635.100.6.1.7` -- Mac App Development
   - `1.2.840.113635.100.6.1.13` -- Developer ID Application

2. **CN prefix fallback** -- when no OID marker is found, matches the certificate Common Name against known prefixes (`Apple Development:`, `Apple Distribution:`, `iPhone Developer:`, etc.)

## OCSP Revocation Checking

On every run, the tool performs a live OCSP check against Apple's certificate authority:

1. Extracts the OCSP responder URI from the certificate (falls back to `ocsp.apple.com`)
2. Fetches the Apple WWDR G3 issuer certificate
3. Queries the OCSP responder via `openssl ocsp`
4. Parses the response for revocation status, time, and reason

If the certificate has been revoked, the output shows the revocation timestamp and reason (Key Compromise, Superseded, etc.).

## Requirements

| Requirement | Notes |
|---|---|
| **Go 1.22+** | For building from source |
| **macOS** | Required -- uses `security cms` to decode CMS/PKCS#7 signed profiles |
| **openssl** | Used for OCSP revocation checks (pre-installed on macOS) |

> **Note:** This tool relies on macOS-specific commands (`security cms`) and cannot run on Linux or Windows.

## License

MIT
