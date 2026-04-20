package csr

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func GenerateZATCACSR(cfg *CsrConfig, outputDir string) (string, string, error) {
	if cfg == nil {
		return "", "", fmt.Errorf("csr config is nil")
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", "", fmt.Errorf("failed to create output directory: %w", err)
	}

	keyPath := outputDir + "/private_key.pem"
	csrPath := outputDir + "/taxpayer.csr"
	configPath := outputDir + "/zatca.cnf"

	configContent := buildOpenSSLConfig(cfg)

	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		return "", "", fmt.Errorf("failed to write OpenSSL config: %w", err)
	}

	err := runCommand(
		"openssl",
		"ecparam",
		"-name", "prime256v1",
		"-genkey",
		"-noout",
		"-out", keyPath,
	)
	if err != nil {
		return "", "", fmt.Errorf("EC key generation failed: %w", err)
	}

	err = runCommand(
		"openssl",
		"req",
		"-new",
		"-key", keyPath,
		"-out", csrPath,
		"-config", configPath,
	)
	if err != nil {
		return "", "", fmt.Errorf("CSR generation failed: %w", err)
	}

	_ = runCommandAndPrint("openssl", "req", "-in", csrPath, "-text", "-noout")

	return keyPath, csrPath, nil
}

func buildOpenSSLConfig(cfg *CsrConfig) string {
	template := templateNameForEnv(cfg.Environment)
	industry := sanitizeIndustry(cfg.IndustryBusinessCategory)

	return `
oid_section = OIDs

[ OIDs ]
certificateTemplateName = 1.3.6.1.4.1.311.20.2

[ req ]
prompt             = no
default_md         = sha256
distinguished_name = dn
req_extensions     = req_ext
string_mask        = utf8only

[ dn ]
C  = ` + cfg.CountryName + `
OU = ` + cfg.OrganizationUnitName + `
O  = ` + cfg.OrganizationName + `
CN = ` + cfg.CommonName + `

[ req_ext ]
certificateTemplateName = ASN1:PRINTABLESTRING:` + template + `
subjectAltName = dirName:zatca_alt_names

[ zatca_alt_names ]
SN                = ` + cfg.SerialNumber + `
UID               = ` + cfg.OrganizationIdentifier + `
title             = ` + cfg.InvoiceType + `
registeredAddress = ` + cfg.LocationAddress + `
businessCategory  = ` + industry + `
`
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"command [%s %s] failed: %w\noutput: %s",
			name,
			strings.Join(args, " "),
			err,
			string(out),
		)
	}
	return nil
}

func runCommandAndPrint(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	_, err := cmd.CombinedOutput()
	return err
}

func sanitizeIndustry(s string) string {
	s = strings.TrimSpace(s)

	reg := regexp.MustCompile(`[^a-zA-Z0-9 ]`)
	s = reg.ReplaceAllString(s, "")

	s = strings.Join(strings.Fields(s), " ")

	if s == "" {
		s = "General"
	}

	return s
}
