package main

const (
	StatusEnabled    = 0
	StatusConfigured = 1
	StatusDisabled   = 2
)

type CertInfo struct {
	Name         string
	SerialNumber string
	Organization string
	OrgUnit      string
	CertType     string
	SHA1         string
	NotBefore    string
	Expiration   string
	Expired      bool
	Revoked      bool
	RevokeTime   string
	RevokeReason string
}

type ValidityInfo struct {
	Status string
	OK     bool
	Detail string
}

type EntitlementEntry struct {
	Key     string
	Name    string
	Display string
	Status  int
}

type ProfileInfo struct {
	FileName       string
	ProfileName    string
	AppIDName      string
	AppID          string
	TeamID         string
	TeamName       string
	UUID           string
	ProfileType    string
	Platforms      []string
	CreationDate   string
	ExpirationDate string
	DaysRemaining  int
	TimeToLive     int
	Version        int
	IsXcodeManaged bool
	PPQRequired    bool
	Certificate    CertInfo
	Validity       ValidityInfo
	Entitlements   []EntitlementEntry
	Devices        []string
}
