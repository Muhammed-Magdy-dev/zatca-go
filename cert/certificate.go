package cert

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
)

func ParseCertificate(certBytes []byte) (*x509.Certificate, error) {
	certBytes = bytes.TrimSpace(certBytes)
	return x509.ParseCertificate(certBytes)
}

func TokenToCertificateDER(binarySecurityToken string) ([]byte, error) {
	token := strings.TrimSpace(strings.Trim(binarySecurityToken, "\""))
	if token == "" {
		return nil, fmt.Errorf("binary security token is empty")
	}

	if strings.Contains(token, "BEGIN CERTIFICATE") {
		block, _ := pem.Decode([]byte(token))
		if block == nil {
			return nil, fmt.Errorf("failed to decode PEM certificate from token")
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PEM certificate: %w", err)
		}
		return cert.Raw, nil
	}

	firstDecode, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("failed to base64-decode token: %w", err)
	}

	if cert, err := x509.ParseCertificate(firstDecode); err == nil {
		return cert.Raw, nil
	}

	firstText := strings.TrimSpace(string(firstDecode))
	if firstText == "" {
		return nil, fmt.Errorf("decoded token is empty")
	}

	secondDecode, err := base64.StdEncoding.DecodeString(firstText)
	if err != nil {
		return nil, fmt.Errorf("double base64 decode failed: %w", err)
	}

	cert, err := x509.ParseCertificate(secondDecode)
	if err != nil {
		return nil, fmt.Errorf("decoded bytes are not a certificate: %w", err)
	}

	return cert.Raw, nil
}
