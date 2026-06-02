package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(headerStyle.Render("✦ MobileProvision Parser"))
		fmt.Println()
		fmt.Println(valueStyle.Render("Usage: mobileprovisionparser <file.mobileprovision>"))
		fmt.Println(dimStyle.Render("  Parse and display iOS/macOS provisioning profile details"))
		os.Exit(0)
	}

	info, err := ParseProfile(os.Args[1])
	if err != nil {
		fmt.Println(errorStyle.Render("Error: " + err.Error()))
		os.Exit(1)
	}

	PrintProfile(info)
}
