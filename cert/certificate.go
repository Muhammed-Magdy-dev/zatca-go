package cert

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
)

var issuerAttributeNames = map[string]string{
	"2.5.4.6":                    "C",
	"2.5.4.10":                   "O",
	"2.5.4.11":                   "OU",
	"2.5.4.3":                    "CN",
	"2.5.4.5":                    "SERIALNUMBER",
	"2.5.4.7":                    "L",
	"2.5.4.8":                    "ST",
	"2.5.4.9":                    "STREET",
	"2.5.4.17":                   "POSTALCODE",
	"0.9.2342.19200300.100.1.25": "DC",
}

func ParseCertificate(certBytes []byte) (*x509.Certificate, error) {
	certBytes = bytes.TrimSpace(certBytes)
	return x509.ParseCertificate(certBytes)
}

func FormatIssuerName(certificate *x509.Certificate) string {
	if certificate == nil {
		return ""
	}

	if len(certificate.RawIssuer) == 0 {
		return certificate.Issuer.ToRDNSequence().String()
	}

	var rdns pkix.RDNSequence
	if _, err := asn1.Unmarshal(certificate.RawIssuer, &rdns); err != nil {
		return certificate.Issuer.ToRDNSequence().String()
	}

	var parts []string
	for i := len(rdns) - 1; i >= 0; i-- {
		for _, atv := range rdns[i] {
			name := atv.Type.String()
			if short, ok := issuerAttributeNames[name]; ok {
				name = short
			}

			parts = append(parts, name+"="+escapeIssuerValue(fmt.Sprint(atv.Value)))
		}
	}

	return strings.Join(parts, ", ")
}

func escapeIssuerValue(value string) string {
	if value == "" {
		return value
	}

	runes := []rune(value)
	escaped := make([]rune, 0, len(runes))
	for i, r := range runes {
		escape := false
		switch r {
		case ',', '+', '"', '\\', '<', '>', ';':
			escape = true
		case ' ':
			escape = i == 0 || i == len(runes)-1
		case '#':
			escape = i == 0
		}

		if escape {
			escaped = append(escaped, '\\', r)
			continue
		}

		escaped = append(escaped, r)
	}

	return string(escaped)
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
