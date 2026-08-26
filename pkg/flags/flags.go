package flags

import (
	"flag"

	"github.com/prometheus/common/config"
)

// MilaFlags holds configuration for the Mila API client
type MilaFlags struct {
	Email    string
	Password config.Secret
	UnitID   string
}

// AddFlags adds Mila-specific flags to the flag set
func AddFlags(appName string) *MilaFlags {
	email := flag.String(
		appName+".email",
		"",
		"Mila account email address",
	)
	password := flag.String(
		appName+".password",
		"",
		"Mila account password",
	)
	unitID := flag.String(
		appName+".unit_id",
		"",
		"Mila Air unit ID to monitor",
	)

	return &MilaFlags{
		Email:    *email,
		Password: config.Secret(*password),
		UnitID:   *unitID,
	}
}
