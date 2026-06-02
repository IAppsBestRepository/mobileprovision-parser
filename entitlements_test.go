package main

import "testing"

func TestGetEntitlementName_KnownKey(t *testing.T) {
	got := GetEntitlementName("get-task-allow")
	want := "Get Task Allow (Debuggable)"
	if got != want {
		t.Errorf("GetEntitlementName(%q) = %q, want %q", "get-task-allow", got, want)
	}
}

func TestGetEntitlementName_AnotherKnownKey(t *testing.T) {
	got := GetEntitlementName("application-identifier")
	want := "Application Identifier"
	if got != want {
		t.Errorf("GetEntitlementName(%q) = %q, want %q", "application-identifier", got, want)
	}
}

func TestGetEntitlementName_UnknownKey(t *testing.T) {
	key := "com.example.unknown-entitlement"
	got := GetEntitlementName(key)
	if got != key {
		t.Errorf("GetEntitlementName(%q) = %q, want raw key back", key, got)
	}
}

func TestGetEntitlementName_EmptyKey(t *testing.T) {
	got := GetEntitlementName("")
	if got != "" {
		t.Errorf("GetEntitlementName(%q) = %q, want %q", "", got, "")
	}
}

func TestFormatEntitlementValue_BoolTrue(t *testing.T) {
	got := FormatEntitlementValue(true)
	if got != "" {
		t.Errorf("FormatEntitlementValue(true) = %q, want %q", got, "")
	}
}

func TestFormatEntitlementValue_BoolFalse(t *testing.T) {
	got := FormatEntitlementValue(false)
	if got != "" {
		t.Errorf("FormatEntitlementValue(false) = %q, want %q", got, "")
	}
}

func TestFormatEntitlementValue_String(t *testing.T) {
	got := FormatEntitlementValue("production")
	want := "(production)"
	if got != want {
		t.Errorf("FormatEntitlementValue(%q) = %q, want %q", "production", got, want)
	}
}

func TestFormatEntitlementValue_EmptyString(t *testing.T) {
	got := FormatEntitlementValue("")
	if got != "" {
		t.Errorf("FormatEntitlementValue(%q) = %q, want %q", "", got, "")
	}
}

func TestFormatEntitlementValue_SliceTwoItems(t *testing.T) {
	input := []any{"alpha", "beta"}
	got := FormatEntitlementValue(input)
	want := "(alpha, beta)"
	if got != want {
		t.Errorf("FormatEntitlementValue(%v) = %q, want %q", input, got, want)
	}
}

func TestFormatEntitlementValue_SliceThreeItems(t *testing.T) {
	input := []any{"a", "b", "c"}
	got := FormatEntitlementValue(input)
	want := "(a, b, c)"
	if got != want {
		t.Errorf("FormatEntitlementValue(%v) = %q, want %q", input, got, want)
	}
}

func TestFormatEntitlementValue_SliceMoreThanThree(t *testing.T) {
	input := []any{"a", "b", "c", "d"}
	got := FormatEntitlementValue(input)
	want := "(a, b, c, ...)"
	if got != want {
		t.Errorf("FormatEntitlementValue(%v) = %q, want %q", input, got, want)
	}
}

func TestFormatEntitlementValue_SliceSingleItem(t *testing.T) {
	input := []any{"only"}
	got := FormatEntitlementValue(input)
	want := "(only)"
	if got != want {
		t.Errorf("FormatEntitlementValue(%v) = %q, want %q", input, got, want)
	}
}

func TestFormatEntitlementValue_EmptySlice(t *testing.T) {
	input := []any{}
	got := FormatEntitlementValue(input)
	if got != "" {
		t.Errorf("FormatEntitlementValue(%v) = %q, want %q", input, got, "")
	}
}

func TestFormatEntitlementValue_OtherTypeInt(t *testing.T) {
	got := FormatEntitlementValue(42)
	want := "42"
	if got != want {
		t.Errorf("FormatEntitlementValue(42) = %q, want %q", got, want)
	}
}

func TestFormatEntitlementValue_OtherTypeFloat(t *testing.T) {
	got := FormatEntitlementValue(3.14)
	want := "3.14"
	if got != want {
		t.Errorf("FormatEntitlementValue(3.14) = %q, want %q", got, want)
	}
}

func TestGetEntitlementStatus_BoolTrue(t *testing.T) {
	got := GetEntitlementStatus("get-task-allow", true)
	if got != StatusEnabled {
		t.Errorf("GetEntitlementStatus(_, true) = %d, want StatusEnabled (%d)", got, StatusEnabled)
	}
}

func TestGetEntitlementStatus_BoolFalse(t *testing.T) {
	got := GetEntitlementStatus("get-task-allow", false)
	if got != StatusDisabled {
		t.Errorf("GetEntitlementStatus(_, false) = %d, want StatusDisabled (%d)", got, StatusDisabled)
	}
}

func TestGetEntitlementStatus_ConfigurableWildcardString(t *testing.T) {
	got := GetEntitlementStatus("com.apple.developer.associated-domains", "*")
	if got != StatusConfigured {
		t.Errorf("GetEntitlementStatus(configurable, \"*\") = %d, want StatusConfigured (%d)", got, StatusConfigured)
	}
}

func TestGetEntitlementStatus_ConfigurableDotStarSuffix(t *testing.T) {
	got := GetEntitlementStatus("com.apple.developer.associated-domains", "com.example.*")
	if got != StatusConfigured {
		t.Errorf("GetEntitlementStatus(configurable, \"com.example.*\") = %d, want StatusConfigured (%d)", got, StatusConfigured)
	}
}

func TestGetEntitlementStatus_ConfigurableWithConcreteValue(t *testing.T) {
	got := GetEntitlementStatus("com.apple.developer.associated-domains", "applinks:example.com")
	if got != StatusEnabled {
		t.Errorf("GetEntitlementStatus(configurable, concrete) = %d, want StatusEnabled (%d)", got, StatusEnabled)
	}
}

func TestGetEntitlementStatus_ConfigurableNilValue(t *testing.T) {
	got := GetEntitlementStatus("com.apple.security.application-groups", nil)
	if got != StatusConfigured {
		t.Errorf("GetEntitlementStatus(configurable, nil) = %d, want StatusConfigured (%d)", got, StatusConfigured)
	}
}

func TestGetEntitlementStatus_ConfigurableEmptyString(t *testing.T) {
	got := GetEntitlementStatus("com.apple.security.application-groups", "")
	if got != StatusConfigured {
		t.Errorf("GetEntitlementStatus(configurable, \"\") = %d, want StatusConfigured (%d)", got, StatusConfigured)
	}
}

func TestGetEntitlementStatus_ConfigurableEmptySlice(t *testing.T) {
	got := GetEntitlementStatus("keychain-access-groups", []any{})
	if got != StatusConfigured {
		t.Errorf("GetEntitlementStatus(configurable, []) = %d, want StatusConfigured (%d)", got, StatusConfigured)
	}
}

func TestGetEntitlementStatus_ConfigurableSliceAllWildcards(t *testing.T) {
	got := GetEntitlementStatus("keychain-access-groups", []any{"com.example.*", "*"})
	if got != StatusConfigured {
		t.Errorf("GetEntitlementStatus(configurable, [wildcards]) = %d, want StatusConfigured (%d)", got, StatusConfigured)
	}
}

func TestGetEntitlementStatus_ConfigurableSliceMixed(t *testing.T) {
	got := GetEntitlementStatus("keychain-access-groups", []any{"com.example.*", "com.example.specific"})
	if got != StatusEnabled {
		t.Errorf("GetEntitlementStatus(configurable, [mixed]) = %d, want StatusEnabled (%d)", got, StatusEnabled)
	}
}

func TestGetEntitlementStatus_NonConfigurableNil(t *testing.T) {
	got := GetEntitlementStatus("com.apple.developer.siri", nil)
	if got != StatusDisabled {
		t.Errorf("GetEntitlementStatus(non-configurable, nil) = %d, want StatusDisabled (%d)", got, StatusDisabled)
	}
}

func TestGetEntitlementStatus_NonConfigurableEmptyString(t *testing.T) {
	got := GetEntitlementStatus("com.apple.developer.siri", "")
	if got != StatusDisabled {
		t.Errorf("GetEntitlementStatus(non-configurable, \"\") = %d, want StatusDisabled (%d)", got, StatusDisabled)
	}
}

func TestGetEntitlementStatus_NonConfigurableEmptySlice(t *testing.T) {
	got := GetEntitlementStatus("com.apple.developer.siri", []any{})
	if got != StatusDisabled {
		t.Errorf("GetEntitlementStatus(non-configurable, []) = %d, want StatusDisabled (%d)", got, StatusDisabled)
	}
}

func TestGetEntitlementStatus_NonConfigurableWithValue(t *testing.T) {
	got := GetEntitlementStatus("aps-environment", "production")
	if got != StatusEnabled {
		t.Errorf("GetEntitlementStatus(non-configurable, \"production\") = %d, want StatusEnabled (%d)", got, StatusEnabled)
	}
}

func TestGetEntitlementStatus_NonConfigurableSliceWithValues(t *testing.T) {
	got := GetEntitlementStatus("com.apple.developer.siri", []any{"feature-a", "feature-b"})
	if got != StatusEnabled {
		t.Errorf("GetEntitlementStatus(non-configurable, [values]) = %d, want StatusEnabled (%d)", got, StatusEnabled)
	}
}

func TestGetEntitlementStatus_NonConfigurableIntValue(t *testing.T) {
	got := GetEntitlementStatus("com.apple.developer.siri", 1)
	if got != StatusEnabled {
		t.Errorf("GetEntitlementStatus(non-configurable, 1) = %d, want StatusEnabled (%d)", got, StatusEnabled)
	}
}

func TestGetEntitlementStatus_BoolTrueOnConfigurableKey(t *testing.T) {
	got := GetEntitlementStatus("com.apple.security.application-groups", true)
	if got != StatusEnabled {
		t.Errorf("GetEntitlementStatus(configurable, true) = %d, want StatusEnabled (%d)", got, StatusEnabled)
	}
}

func TestGetEntitlementStatus_BoolFalseOnConfigurableKey(t *testing.T) {
	got := GetEntitlementStatus("com.apple.security.application-groups", false)
	if got != StatusDisabled {
		t.Errorf("GetEntitlementStatus(configurable, false) = %d, want StatusDisabled (%d)", got, StatusDisabled)
	}
}

func TestFormatEntitlementValue_SliceFiveItems(t *testing.T) {
	input := []any{"one", "two", "three", "four", "five"}
	got := FormatEntitlementValue(input)
	want := "(one, two, three, ...)"
	if got != want {
		t.Errorf("FormatEntitlementValue(%v) = %q, want %q", input, got, want)
	}
}

func TestFormatEntitlementValue_SliceWithNonStringItems(t *testing.T) {
	input := []any{1, true, "mixed"}
	got := FormatEntitlementValue(input)
	want := "(1, true, mixed)"
	if got != want {
		t.Errorf("FormatEntitlementValue(%v) = %q, want %q", input, got, want)
	}
}
