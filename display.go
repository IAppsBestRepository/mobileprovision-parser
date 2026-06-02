package main

import (
	"fmt"
	"strings"
)

func PrintProfile(info *ProfileInfo) {
	fmt.Println()
	fmt.Println(borderStyle.Render("╭" + strings.Repeat("─", 50) + "╮"))
	fmt.Println(borderStyle.Render("│") + "  " + headerStyle.Render("✦ MobileProvision Parser") + strings.Repeat(" ", 24) + borderStyle.Render("│"))
	fmt.Println(borderStyle.Render("╰" + strings.Repeat("─", 50) + "╯"))
	fmt.Println()

	row := func(label, value string) {
		fmt.Println("  " + labelStyle.Render(label) + valueStyle.Render(value))
	}

	rowStyled := func(label string, value string, style func(string) string) {
		fmt.Println("  " + labelStyle.Render(label) + style(value))
	}

	separator := borderStyle.Render(strings.Repeat("─", 50))

	row("Profile", info.FileName)
	row("Name", info.ProfileName)
	row("App ID Name", info.AppIDName)
	row("App ID", info.AppID)

	team := info.TeamID
	if info.TeamName != "" {
		team = info.TeamID + " (" + info.TeamName + ")"
	}
	row("Team", team)

	row("UUID", info.UUID)
	row("Type", info.ProfileType)
	row("Platform", strings.Join(info.Platforms, ", "))

	if info.IsXcodeManaged {
		row("Xcode Managed", "Yes")
	}

	row("Created", info.CreationDate)

	expires := info.ExpirationDate
	if info.DaysRemaining < 0 {
		expires += "  " + errorStyle.Render("(EXPIRED)")
	} else if info.DaysRemaining > 0 {
		expires += fmt.Sprintf("  (%d days remaining)", info.DaysRemaining)
	}
	row("Expires", expires)

	if info.TimeToLive > 0 {
		row("TTL", fmt.Sprintf("%d days", info.TimeToLive))
	}

	if info.Version > 0 {
		row("Version", fmt.Sprintf("%d", info.Version))
	}

	fmt.Println()
	fmt.Println("  " + sectionStyle.Render("Certificate"))
	fmt.Println("  " + separator)

	certRow := func(label, value string) {
		fmt.Println("    " + labelStyle.Render(label) + valueStyle.Render(value))
	}

	certRow("Name", info.Certificate.Name)

	if info.Certificate.CertType != "" {
		certRow("Type", info.Certificate.CertType)
	}

	if info.Certificate.SerialNumber != "" {
		certRow("Serial", info.Certificate.SerialNumber)
	}

	if info.Certificate.Organization != "" {
		certRow("Organization", info.Certificate.Organization)
	}

	if info.Certificate.OrgUnit != "" {
		certRow("Org Unit", info.Certificate.OrgUnit)
	}

	if info.Certificate.SHA1 != "" {
		certRow("SHA-1", info.Certificate.SHA1)
	}

	if info.Certificate.NotBefore != "" && info.Certificate.Expiration != "" {
		certRow("Valid", info.Certificate.NotBefore+" — "+info.Certificate.Expiration)
	} else if info.Certificate.Expiration != "" {
		certRow("Valid", "— "+info.Certificate.Expiration)
	}

	if info.Certificate.Revoked {
		revoked := errorStyle.Render("✗ Revoked")
		if info.Certificate.RevokeTime != "" {
			revoked += "  " + dimStyle.Render(info.Certificate.RevokeTime)
		}
		if info.Certificate.RevokeReason != "" {
			revoked += "  " + dimStyle.Render("("+info.Certificate.RevokeReason+")")
		}
		fmt.Println("    " + labelStyle.Render("Status") + revoked)
	} else if info.Certificate.Expired {
		fmt.Println("    " + labelStyle.Render("Status") + errorStyle.Render("✗ Expired"))
	} else {
		fmt.Println("    " + labelStyle.Render("Status") + successStyle.Render("✓ Signed"))
	}

	fmt.Println()

	if info.Validity.OK {
		status := successStyle.Render("✓ " + info.Validity.Status)
		if info.Validity.Status == "Expiring Soon" {
			status = warningStyle.Render("✓ " + info.Validity.Status)
		}
		if info.Validity.Detail != "" {
			status += "  " + dimStyle.Render(info.Validity.Detail)
		}
		rowStyled("Profile Status", "", func(_ string) string { return status })
	} else {
		status := errorStyle.Render("✗ " + info.Validity.Status)
		if info.Validity.Detail != "" {
			status += "  " + dimStyle.Render(info.Validity.Detail)
		}
		rowStyled("Profile Status", "", func(_ string) string { return status })
	}

	if info.PPQRequired {
		rowStyled("PPQ Check", "Enabled", func(s string) string { return errorStyle.Render(s) })
	} else {
		rowStyled("PPQ Check", "Disabled", func(s string) string { return successStyle.Render(s) })
	}

	if len(info.Devices) > 0 {
		fmt.Println()
		fmt.Println("  " + sectionStyle.Render(fmt.Sprintf("Devices (%d)", len(info.Devices))))

		limit := 5
		if len(info.Devices) < limit {
			limit = len(info.Devices)
		}
		for i := 0; i < limit; i++ {
			fmt.Println("    " + deviceStyle.Render(fmt.Sprintf("%d. %s", i+1, info.Devices[i])))
		}
		if len(info.Devices) > 5 {
			fmt.Println("    " + dimStyle.Render(fmt.Sprintf("... and %d more", len(info.Devices)-5)))
		}
	}

	if len(info.Entitlements) > 0 {
		fmt.Println()
		fmt.Println("  " + sectionStyle.Render(fmt.Sprintf("Entitlements (%d)", len(info.Entitlements))))
		fmt.Println("  " + separator)

		for _, ent := range info.Entitlements {
			switch ent.Status {
			case StatusEnabled:
				line := "    " + entEnabled.Render("✓") + " " + valueStyle.Render(ent.Name)
				if ent.Display != "" {
					line += "  " + dimStyle.Render(ent.Display)
				}
				fmt.Println(line)
			case StatusConfigured:
				line := "    " + entConfigured.Render("●") + " " + valueStyle.Render(ent.Name)
				if ent.Display != "" {
					line += "  " + dimStyle.Render(ent.Display)
				}
				fmt.Println(line)
			case StatusDisabled:
				fmt.Println("    " + entDisabled.Render("✗") + " " + dimStyle.Render(ent.Name))
			}
		}
	}

	fmt.Println()
}
