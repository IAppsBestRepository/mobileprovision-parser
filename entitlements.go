package main

import "fmt"

var EntitlementNames = map[string]string{
	"application-identifier":                          "Application Identifier",
	"com.apple.developer.team-identifier":             "Team Identifier",
	"keychain-access-groups":                          "Keychain Access Groups",
	"get-task-allow":                                  "Get Task Allow (Debuggable)",
	"com.apple.developer.user-fonts":                  "User Fonts",
	"com.apple.developer.accessibility.merchant-api-control":            "Accessibility Merchant API Control",
	"com.apple.developer.accessory-setup-discovery-extension":           "Accessory Setup Discovery Extension",
	"com.apple.developer.accessory-transport-extension":                 "Accessory Transport Extension",
	"com.apple.developer.marketplace.app-installation":                  "Alternative Marketplace App Installation",
	"com.apple.developer.parent-application-identifiers":                "Parent Application Identifiers",
	"com.apple.developer.associated-appclip-app-identifiers":            "Associated App Clip Identifiers",
	"com.apple.developer.on-demand-install-capable":                     "App Clip On-Demand Install",
	"com.apple.developer.side-button-access.allow":                      "Side Button Access",
	"com.apple.developer.app-migration.data-container-access":           "App Migration Data Container Access",
	"com.apple.developer.authentication-services.account-creation-requires-phone-number": "Account Creation Requires Phone Number",
	"com.apple.developer.authentication-services.autofill-credential-provider":            "AutoFill Credential Provider",
	"com.apple.developer.applesignin":                                   "Sign in with Apple",
	"com.apple.developer.background-tasks.continued-processing.gpu":     "Background GPU Access",
	"com.apple.developer.calling-app":                                   "Default Calling App",
	"com.apple.developer.carplay-audio":                                 "CarPlay Audio",
	"com.apple.developer.carplay-charging":                              "CarPlay Charging",
	"com.apple.developer.carplay-communication":                         "CarPlay Communication",
	"com.apple.developer.carplay-maps":                                  "CarPlay Maps",
	"com.apple.developer.carplay-parking":                               "CarPlay Parking",
	"com.apple.developer.carplay-quick-ordering":                        "CarPlay Quick Ordering",
	"com.apple.developer.carplay-messaging":                             "CarPlay Messaging (Deprecated)",
	"com.apple.developer.playable-content":                              "Playable Content (Deprecated)",
	"com.apple.developer.icloud-extended-share-access":                  "iCloud Extended Share Access",
	"com.apple.developer.contacts.notes":                                "Contacts Notes Access",
	"com.apple.CommCenter.fine-grained":                                 "CoreTelephony Services",
	"com.apple.developer.declared-age-range":                            "Declared Age Range",
	"com.apple.developer.automated-device-enrollment.add-devices":       "Automated Device Enrollment",
	"com.apple.developer.enrollment-sso-capable":                        "Enrollment Single Sign-On",
	"com.apple.developer.ClassKit-environment":                          "ClassKit Environment",
	"com.apple.developer.automatic-assessment-configuration":            "Automatic Assessment Configuration",
	"com.apple.developer.mail-client":                                   "Default Mail Client",
	"com.apple.developer.energykit":                                     "EnergyKit",
	"com.apple.developer.app-compute-category":                          "Increased Performance Headroom",
	"com.apple.developer.screen-capture.include-passthrough":            "Passthrough in Screen Capture",
	"com.apple.developer.arkit.main-camera-access.allow":                "Main Camera Access (visionOS)",
	"com.apple.developer.arkit.object-tracking-parameter-adjustment.allow": "Object Tracking Parameter Adjustment",
	"com.apple.developer.arkit.barcode-detection.allow":                 "Spatial Barcode/QR Code Scanning",
	"com.apple.developer.arkit.camera-region.allow":                     "Camera Region Access",
	"com.apple.developer.arkit.shared-coordinate-space.allow":           "Shared Coordinate Space Access",
	"com.apple.developer.protected-content":                             "App-Protected Content",
	"com.apple.developer.window-body-follow":                            "Window Follow Mode",
	"com.apple.developer.coreml.neural-engine-access":                   "Apple Neural Engine Access (Deprecated)",
	"com.apple.developer.avfoundation.uvc-device-access":                "UVC Device Access (visionOS, Deprecated)",
	"com.apple.developer.exposure-notification":                         "Exposure Notification",
	"com.apple.developer.family-controls":                               "Family Controls",
	"com.apple.developer.fileprovider.testing-mode":                     "FileProvider Testing Mode",
	"com.apple.developer.financekit":                                    "FinanceKit",
	"com.apple.developer.foundation-model-adapter":                      "Foundation Model Adapter",
	"com.apple.developer.fskit.fsmodule":                                "FSKit Module",
	"com.apple.developer.game-center":                                   "Game Center",
	"com.apple.developer.group-session":                                 "SharePlay / Group Activities",
	"com.apple.developer.healthkit":                                     "HealthKit",
	"com.apple.developer.healthkit.access":                              "HealthKit Capabilities",
	"com.apple.developer.healthkit.background-delivery":                 "HealthKit Background Delivery",
	"com.apple.developer.health.fall-detection":                         "Fall Detection Notifications",
	"com.apple.developer.healthkit.recalibrate-estimates":               "HealthKit Recalibrate Estimates",
	"com.apple.developer.homekit":                                       "HomeKit",
	"com.apple.developer.matter.allow-setup-payload":                    "Matter Allow Setup Payload",
	"com.apple.security.hypervisor":                                     "Hypervisor",
	"com.apple.vm.hypervisor":                                           "VM Hypervisor (Deprecated)",
	"com.apple.vm.device-access":                                        "VM Device Access",
	"com.apple.vm.networking":                                           "VM Networking",
	"com.apple.security.virtualization":                                  "Virtualization",
	"com.apple.developer.icloud-container-development-container-identifiers": "iCloud Development Container Identifiers",
	"com.apple.developer.icloud-container-environment":                  "iCloud Container Environment",
	"com.apple.developer.icloud-container-identifiers":                  "iCloud Container Identifiers",
	"com.apple.developer.icloud-services":                               "iCloud Services",
	"com.apple.developer.ubiquity-kvstore-identifier":                   "iCloud Key-Value Store",
	"com.apple.developer.ubiquity-container-identifiers":                "iCloud Ubiquity Container Identifiers",
	"com.apple.developer.identity-document-services.document-provider.mobile-document-types": "Digital Credentials Mobile Document Provider",
	"com.apple.developer.journal.allow":                                 "Journaling Suggestions",
	"com.apple.developer.dialing-app":                                   "Default Dialer App",
	"com.apple.developer.location.push":                                 "Location Push Service Extension",
	"com.apple.developer.managed-app-distribution.install-ui":           "Managed App Installation UI",
	"com.apple.developer.media-device-discovery-extension":              "Media Device Discovery Extension",
	"com.apple.developer.coremotion.head-pose":                          "Head Pose (Spatial Audio)",
	"com.apple.developer.spatial-audio.profile-access":                  "Spatial Audio Profile Access",
	"com.apple.developer.avfoundation.multitasking-camera-access":       "Multitasking Camera Access (Deprecated)",
	"com.apple.developer.coremedia.hls.low-latency":                     "Low Latency HLS",
	"com.apple.developer.devicecheck.appattest-environment":             "App Attest Environment",
	"com.apple.developer.kernel.increased-memory-limit":                 "Increased Memory Limit",
	"com.apple.developer.kernel.extended-virtual-addressing":            "Extended Virtual Addressing",
	"com.apple.developer.kernel.increased-debugging-memory-limit":       "Increased Debugging Memory Limit",
	"com.apple.developer.sustained-execution":                           "Sustained Execution",
	"com.apple.developer.messages.critical-messaging":                   "Critical Messaging",
	"com.apple.developer.messaging-app":                                 "Default Messaging App",
	"com.apple.developer.shared-with-you":                               "Shared with You",
	"com.apple.developer.shared-with-you.collaboration":                 "Messages Collaboration",
	"com.apple.developer.upi-device-validation":                         "UPI Device Validation",
	"com.apple.developer.coretelephony.sim-inserted":                    "SIM Inserted Detection",
	"com.apple.developer.proximity-reader.identity.display":             "ID Verifier - Display Only",
	"com.apple.developer.navigation-app":                                "Default Navigation App",
	"com.apple.developer.nearbyinteraction.dltdoa":                      "Nearby Interaction DL-TDOA",
	"com.apple.developer.networking.networkextension":                   "Network Extensions",
	"com.apple.developer.networking.vpn.api":                            "Personal VPN",
	"com.apple.developer.associated-domains":                            "Associated Domains",
	"com.apple.developer.networking.multicast":                          "Multicast Networking",
	"com.apple.developer.associated-domains.applinks.read-write":        "Universal Links Read-Write",
	"com.apple.developer.networking.manage-thread-network-credentials":  "Thread Network Credentials",
	"com.apple.developer.networking.slicing.appcategory":                "5G Network Slicing App Category",
	"com.apple.developer.networking.slicing.trafficcategory":            "5G Network Slicing Traffic Category",
	"com.apple.developer.networking.vmnet":                              "VMNet Networking",
	"com.apple.developer.networking.carrier-constrained.appcategory":    "Carrier Constrained App Category",
	"com.apple.developer.networking.carrier-constrained.app-optimized":  "Carrier Constrained App Optimized",
	"aps-environment":                                                   "Push Notifications (iOS)",
	"com.apple.developer.aps-environment":                               "Push Notifications (macOS)",
	"com.apple.developer.usernotifications.critical-alerts":             "Critical Alerts",
	"com.apple.developer.usernotifications.filtering":                   "Notification Filtering",
	"com.apple.developer.usernotifications.communication":               "Communication Notifications",
	"com.apple.developer.usernotifications.time-sensitive":              "Time Sensitive Notifications",
	"com.apple.developer.passkit.pass-presentation-suppression":         "Pass Presentation Suppression",
	"com.apple.developer.device-information.user-assigned-device-name":  "User Assigned Device Name",
	"com.apple.developer.push-to-talk":                                  "Push to Talk",
	"com.apple.developer.severe-vehicular-crash-event":                  "Crash Detection Events",
	"com.apple.developer.secure-element-credential":                     "Secure Element Credential",
	"com.apple.developer.secure-element-credential.default-contactless-app": "Default Contactless App",
	"com.apple.developer.sensitivecontentanalysis.client":               "Sensitive Content Analysis",
	"com.apple.developer.sensorkit.reader.allow":                        "SensorKit Reader",
	"com.apple.developer.siri":                                          "Siri",
	"com.apple.developer.storekit.external-link.account":                "StoreKit External Link Account",
	"com.apple.developer.storekit.external-purchase":                    "StoreKit External Purchase",
	"com.apple.developer.storekit.external-purchase-link":               "StoreKit External Purchase Link",
	"com.apple.developer.storekit.external-purchase-link-streaming":     "StoreKit External Purchase Link Streaming",
	"com.apple.developer.carrier-messaging-app":                         "Default Carrier Messaging App",
	"com.apple.developer.translation-app":                               "Default Translation App",
	"com.apple.developer.user-management":                               "User Management (tvOS)",
	"com.apple.developer.video-subscriber-single-sign-on":               "TV Provider Authentication",
	"com.apple.smoot.subscriptionservice":                               "Video Partner Program",
	"com.apple.developer.low-latency-streaming":                         "Low-Latency Streaming",
	"com.apple.developer.pass-type-identifiers":                         "Pass Type Identifiers",
	"com.apple.developer.in-app-payments":                               "Apple Pay (Merchant IDs)",
	"com.apple.developer.in-app-identity-presentment":                   "In-App Identity Presentment",
	"com.apple.developer.in-app-identity-presentment.merchant-identifiers": "In-App Identity Presentment Merchant IDs",
	"com.apple.developer.weatherkit":                                    "WeatherKit",
	"com.apple.developer.web-browser":                                   "Default Web Browser",
	"com.apple.developer.web-browser.public-key-credential":             "Web Browser Public Key Credential",
	"com.apple.developer.browser.app-installation":                      "Browser App Installation",
	"com.apple.developer.networking.wifi-info":                           "Access Wi-Fi Information",
	"com.apple.external-accessory.wireless-configuration":               "Wireless Accessory Configuration",
	"com.apple.developer.networking.multipath":                          "Multipath",
	"com.apple.developer.networking.HotspotConfiguration":               "Hotspot Configuration",
	"com.apple.developer.networking.HotspotHelper":                      "Hotspot Helper",
	"com.apple.developer.nfc.readersession.felica.systemcodes":          "NFC FeliCa System Codes",
	"com.apple.developer.nfc.readersession.formats":                     "NFC Tag Reader Session Formats",
	"com.apple.developer.nfc.readersession.iso7816.select-identifiers":  "NFC ISO7816 Select Identifiers",
	"com.apple.developer.nfc.hce":                                       "NFC Card Session (HCE)",
	"com.apple.developer.nfc.hce.iso7816.select-identifier-prefixes":    "NFC HCE ISO7816 Identifier Prefixes",
	"com.apple.developer.nfc.hce.default-contactless-app":               "NFC HCE Default Contactless App",
	"com.apple.developer.wireless-insights.service-predictions":         "Wireless Insights Service Predictions",
	"com.apple.developer.maps":                                          "Maps (Deprecated)",
	"inter-app-audio":                                                   "Inter-App Audio (Deprecated)",
	"com.apple.security.files.all":                                      "All Files Access (Deprecated)",
	"com.apple.developer.embedded-web-browser-engine":                   "Embedded Web Browser Engine",
	"com.apple.developer.memory.transfer_accept":                        "Memory Transfer Accept",
	"com.apple.developer.memory.transfer_send":                          "Memory Transfer Send",
	"com.apple.developer.web-browser-engine.host":                       "Web Browser Engine Host",
	"com.apple.developer.web-browser-engine.networking":                 "Web Browser Engine Networking",
	"com.apple.developer.web-browser-engine.rendering":                  "Web Browser Engine Rendering",
	"com.apple.developer.web-browser-engine.webcontent":                 "Web Browser Engine Web Content",
	"com.apple.developer.persistent-content-capture":                    "Persistent Content Capture",
	"com.apple.developer.wifi-aware":                                    "Wi-Fi Aware",
	"com.apple.developer.wifi-infrastructure":                           "Wi-Fi Infrastructure",
	"com.apple.security.application-groups":                             "App Groups",
	"com.apple.developer.default-data-protection":                       "Data Protection",
	"com.apple.developer.driverkit":                                     "DriverKit",
	"com.apple.developer.driverkit.allow-third-party-userclients":       "DriverKit Allow Third Party UserClients",
	"com.apple.developer.driverkit.transport.pci":                       "DriverKit PCI Transport",
	"com.apple.developer.driverkit.transport.usb":                       "DriverKit USB Transport",
	"com.apple.developer.driverkit.transport.hid":                       "DriverKit HID Transport",
	"com.apple.developer.system-extension.install":                      "System Extension Install",
}

var ConfigurableEntitlements = map[string]bool{
	"com.apple.security.application-groups":                true,
	"com.apple.developer.icloud-services":                  true,
	"com.apple.developer.icloud-container-identifiers":     true,
	"com.apple.developer.ubiquity-container-identifiers":   true,
	"com.apple.developer.associated-domains":               true,
	"com.apple.developer.in-app-payments":                  true,
	"com.apple.developer.pass-type-identifiers":            true,
	"com.apple.developer.networking.networkextension":      true,
	"com.apple.developer.nfc.readersession.formats":        true,
	"com.apple.developer.healthkit.access":                 true,
	"keychain-access-groups":                               true,
}

func GetEntitlementName(key string) string {
	if name, ok := EntitlementNames[key]; ok {
		return name
	}
	return key
}

func FormatEntitlementValue(value any) string {
	switch v := value.(type) {
	case bool:
		return ""
	case string:
		if v == "" {
			return ""
		}
		return "(" + v + ")"
	case []any:
		if len(v) == 0 {
			return ""
		}
		if len(v) <= 3 {
			s := "("
			for i, item := range v {
				if i > 0 {
					s += ", "
				}
				s += fmt.Sprintf("%v", item)
			}
			s += ")"
			return s
		}
		s := "("
		for i := 0; i < 3; i++ {
			if i > 0 {
				s += ", "
			}
			s += fmt.Sprintf("%v", v[i])
		}
		s += ", ...)"
		return s
	default:
		return fmt.Sprintf("%v", value)
	}
}

func isWildcardOrEmpty(value any) bool {
	if value == nil {
		return true
	}
	switch v := value.(type) {
	case string:
		if v == "" || v == "*" {
			return true
		}
		if len(v) >= 2 && v[len(v)-2:] == ".*" {
			return true
		}
		return false
	case []any:
		if len(v) == 0 {
			return true
		}
		for _, item := range v {
			s, ok := item.(string)
			if !ok {
				return false
			}
			if s != "*" && (len(s) < 2 || s[len(s)-2:] != ".*") {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func GetEntitlementStatus(key string, value any) int {
	if b, ok := value.(bool); ok {
		if b {
			return StatusEnabled
		}
		return StatusDisabled
	}
	if ConfigurableEntitlements[key] {
		if isWildcardOrEmpty(value) {
			return StatusConfigured
		}
		return StatusEnabled
	}
	if value == nil {
		return StatusDisabled
	}
	if s, ok := value.(string); ok && s == "" {
		return StatusDisabled
	}
	if a, ok := value.([]any); ok && len(a) == 0 {
		return StatusDisabled
	}
	return StatusEnabled
}
