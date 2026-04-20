package csr

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

type CSRResult struct {
	CSRBase64  string
	PrivateKey string
	CSRPEM     string
}

func GenerateZatcaCSR(cfg *CsrConfig) (*CSRResult, error) {
	if err := validateCsrConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "zatca-csr-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "zatca.cnf")
	privateKeyPath := filepath.Join(tmpDir, "privatekey.pem")
	csrPath := filepath.Join(tmpDir, "taxpayer.csr")

	conf := buildOpenSSLConfig(cfg)

	if err := os.WriteFile(configPath, []byte(conf), 0o600); err != nil {
		return nil, fmt.Errorf("write config: %w", err)
	}

	if err := runCommand(
		"openssl",
		"ecparam",
		"-name", "prime256v1",
		"-genkey",
		"-noout",
		"-out", privateKeyPath,
	); err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}

	if err := runCommand(
		"openssl",
		"req",
		"-new",
		"-key", privateKeyPath,
		"-config", configPath,
		"-out", csrPath,
	); err != nil {
		return nil, fmt.Errorf("generate CSR: %w", err)
	}

	privateKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}

	csrPEM, err := os.ReadFile(csrPath)
	if err != nil {
		return nil, fmt.Errorf("read CSR: %w", err)
	}

	return &CSRResult{
		CSRBase64:  base64.StdEncoding.EncodeToString(csrPEM),
		PrivateKey: string(privateKeyPEM),
		CSRPEM:     string(csrPEM),
	}, nil
}
