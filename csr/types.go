package csr

import (
	"fmt"
	"strings"
)

func validateCsrConfig(cfg *CsrConfig) error {
	if cfg == nil {
		return fmt.Errorf("csr config is nil")
	}

	if strings.TrimSpace(cfg.CommonName) == "" {
		return fmt.Errorf("CommonName is required")
	}
	if strings.TrimSpace(cfg.SerialNumber) == "" {
		return fmt.Errorf("SerialNumber is required")
	}
	if strings.TrimSpace(cfg.OrganizationIdentifier) == "" {
		return fmt.Errorf("OrganizationIdentifier is required")
	}
	if strings.TrimSpace(cfg.OrganizationUnitName) == "" {
		return fmt.Errorf("OrganizationUnitName is required")
	}
	if strings.TrimSpace(cfg.OrganizationName) == "" {
		return fmt.Errorf("OrganizationName is required")
	}
	if strings.TrimSpace(cfg.CountryName) == "" {
		return fmt.Errorf("CountryName is required")
	}
	if strings.TrimSpace(cfg.InvoiceType) == "" {
		return fmt.Errorf("InvoiceType is required")
	}
	if strings.TrimSpace(cfg.LocationAddress) == "" {
		return fmt.Errorf("LocationAddress is required")
	}
	if strings.TrimSpace(cfg.IndustryBusinessCategory) == "" {
		return fmt.Errorf("IndustryBusinessCategory is required")
	}

	if strings.ToUpper(cfg.CountryName) != "SA" {
		return fmt.Errorf("CountryName must be SA")
	}

	if len(cfg.OrganizationIdentifier) != 15 {
		return fmt.Errorf("OrganizationIdentifier must be 15 digits")
	}

	for _, r := range cfg.OrganizationIdentifier {
		if r < '0' || r > '9' {
			return fmt.Errorf("OrganizationIdentifier must be digits only")
		}
	}

	if cfg.OrganizationIdentifier[0] != '3' || cfg.OrganizationIdentifier[len(cfg.OrganizationIdentifier)-1] != '3' {
		return fmt.Errorf("OrganizationIdentifier must start and end with 3")
	}

	if len(cfg.InvoiceType) != 4 {
		return fmt.Errorf("InvoiceType must be 4 digits")
	}

	for _, r := range cfg.InvoiceType {
		if r != '0' && r != '1' {
			return fmt.Errorf("InvoiceType must be binary")
		}
	}

	if cfg.InvoiceType == "0000" {
		return fmt.Errorf("InvoiceType cannot be 0000")
	}

	return nil
}
