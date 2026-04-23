package csr

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

const (
	defaultCountryCode = "SA"
	defaultInvoiceType = "0100"
	defaultAddress     = "Riyadh, Saudi Arabia"
	defaultOrgName     = "Unknown Organization"
	defaultSector      = "General Trading"
	defaultDeviceName  = "POS Device"
)

var vatDigitsOnly = regexp.MustCompile(`\D+`)

type CSRInput struct {
	VATNumber            string
	CommonName           string
	OrganizationUnitName string
	OrganizationName     string
	BranchName           string
	SectorName           string
	Address              string
	CountryCode          string
	PosID                *uuid.UUID
	InvoiceType          string
}

type CsrConfig struct {
	CountryName              string
	OrganizationName         string
	OrganizationUnitName     string
	CommonName               string
	SerialNumber             string
	OrganizationIdentifier   string
	LocationAddress          string
	IndustryBusinessCategory string
	InvoiceType              string
	Environment              string
}

func BuildCSRConfig(data *CSRInput, baseURL string) (*CsrConfig, error) {
	if data == nil {
		return nil, fmt.Errorf("zatca: data is nil")
	}

	if strings.TrimSpace(data.VATNumber) == "" {
		return nil, fmt.Errorf("zatca: VAT number is required")
	}

	vat := normalizeVAT(data.VATNumber)
	if len(vat) != 15 {
		return nil, fmt.Errorf("zatca: VAT number must be 15 digits, got %q", vat)
	}

	orgName := pickString(&data.OrganizationName, defaultOrgName)
	commonName := pickString(&data.CommonName, defaultDeviceName)
	orgUnit := pickString(&data.OrganizationUnitName, defaultDeviceName)
	sector := pickString(&data.SectorName, defaultSector)
	address := pickString(&data.Address, defaultAddress)

	country := defaultCountryCode
	if strings.TrimSpace(data.CountryCode) != "" {
		country = strings.ToUpper(strings.TrimSpace(data.CountryCode))
	}

	env := determineEnvironment(baseURL)
	serial := buildSerialNumber("Wasfa", commonName, data.PosID)
	invoiceType, err := normalizeInvoiceType(data.InvoiceType)
	if err != nil {
		return nil, err
	}
	cfg := &CsrConfig{
		CountryName:              country,
		OrganizationName:         orgName,
		OrganizationUnitName:     orgUnit,
		CommonName:               commonName,
		SerialNumber:             serial,
		OrganizationIdentifier:   vat,
		LocationAddress:          address,
		IndustryBusinessCategory: sector,
		InvoiceType:              invoiceType,
		Environment:              env,
	}

	return cfg, nil
}

func pickString(value *string, fallback string) string {
	if value == nil {
		return fallback
	}
	if trimmed := strings.TrimSpace(*value); trimmed != "" {
		return trimmed
	}
	return fallback
}

func normalizeVAT(v string) string {
	trimmed := strings.TrimSpace(v)
	digitsOnly := vatDigitsOnly.ReplaceAllString(trimmed, "")
	return digitsOnly
}

func determineEnvironment(baseURL string) string {
	lower := strings.ToLower(strings.TrimSpace(baseURL))
	switch {
	case strings.Contains(lower, "nonproduction"):
		return "nonProduction"
	case strings.Contains(lower, "production") && !strings.Contains(lower, "simulation"):
		return "production"
	default:
		return "simulation"
	}
}

func buildSerialNumber(orgName, deviceName string, posID *uuid.UUID) string {
	part1 := sanitizeSerialComponent(orgName, "WASFA")
	part2 := sanitizeSerialComponent(deviceName, "POS")

	part3 := "POS-DEVICE"
	if posID != nil {
		part3 = strings.ToUpper(strings.ReplaceAll(posID.String(), "-", ""))
	}

	return fmt.Sprintf("1-%s|2-%s|3-%s", part1, part2, part3)
}

func sanitizeSerialComponent(value string, fallback string) string {
	if value == "" {
		return fallback
	}

	builder := strings.Builder{}
	for _, r := range value {
		switch {
		case unicode.IsLetter(r):
			builder.WriteRune(unicode.ToUpper(r))
		case unicode.IsDigit(r):
			builder.WriteRune(r)
		case r == '-':
			builder.WriteRune('-')
		}
	}

	result := strings.Trim(builder.String(), "-")
	if result == "" {
		return fallback
	}

	return result
}
func normalizeInvoiceType(v string) (string, error) {
	t := strings.TrimSpace(strings.ToLower(v))

	switch t {
	case "", "both", "1100":
		return "1100", nil

	case "standard", "1000":
		return "1000", nil

	case "simplified", "0100":
		return "0100", nil

	default:
		return "", fmt.Errorf("zatca: invalid invoice type %q", v)
	}
}
