package invoice

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math"
	"strconv"

	"github.com/Muhammed-Magdy-dev/zatca-go/cert"
)

func formatFloat2(v float64) string {
	return strconv.FormatFloat(math.Round(v*100)/100, 'f', 2, 64)
}

func BuildQRCodeTLVBase64(input *InvoiceInput) (string, error) {
	if input == nil {
		return "", fmt.Errorf("input is nil")
	}

	totals := CalculateTotals(input)

	fields := make([][]byte, 0, 9)

	addStringField := func(tag byte, value string) error {
		field, err := encodeTLV(tag, []byte(value))
		if err != nil {
			return err
		}
		fields = append(fields, field)
		return nil
	}

	addBytesField := func(tag byte, value []byte) error {
		field, err := encodeTLV(tag, value)
		if err != nil {
			return err
		}
		fields = append(fields, field)
		return nil
	}

	if err := addStringField(1, input.Supplier.RegistrationName); err != nil {
		return "", fmt.Errorf("failed to add seller name: %w", err)
	}

	if err := addStringField(2, input.Supplier.VATNumber); err != nil {
		return "", fmt.Errorf("failed to add VAT number: %w", err)
	}

	if err := addStringField(3, input.IssueDate.Format("2006-01-02T15:04:05")); err != nil {
		return "", fmt.Errorf("failed to add timestamp: %w", err)
	}

	if err := addStringField(4, formatFloat2(totals.TaxInclusiveAmount)); err != nil {
		return "", fmt.Errorf("failed to add total with VAT: %w", err)
	}

	if err := addStringField(5, formatFloat2(totals.TaxAmount)); err != nil {
		return "", fmt.Errorf("failed to add VAT amount: %w", err)
	}

	if err := addStringField(6, input.InvoiceDigest); err != nil {
		return "", fmt.Errorf("failed to add invoice hash: %w", err)
	}

	if err := addStringField(7, input.SignatureValue); err != nil {
		return "", fmt.Errorf("failed to add signature: %w", err)
	}

	if err := addBytesField(8, input.PublicKeyBytes); err != nil {
		return "", fmt.Errorf("failed to add public key: %w", err)
	}

	if err := addBytesField(9, input.CertificateSignatureBytes); err != nil {
		return "", fmt.Errorf("failed to add certificate signature: %w", err)
	}

	var buffer bytes.Buffer
	for _, field := range fields {
		buffer.Write(field)
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func encodeTLV(tag byte, value []byte) ([]byte, error) {
	if len(value) > 255 {
		return nil, fmt.Errorf("TLV value exceeds maximum length of 255 bytes for tag %d", tag)
	}

	buffer := make([]byte, 0, 2+len(value))
	buffer = append(buffer, tag)
	buffer = append(buffer, byte(len(value)))
	buffer = append(buffer, value...)

	return buffer, nil
}

func PublicKeyDERFromPrivateKey(privateKeyPEM []byte) ([]byte, error) {
	key, err := cert.ParseECPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	return pubDER, nil
}

func PublicKeyBase64FromPrivateKey(privateKeyPEM []byte) (string, error) {
	pubDER, err := PublicKeyDERFromPrivateKey(privateKeyPEM)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(pubDER), nil
}

func CertificateSignatureBytesFromCertificate(certBytes []byte) ([]byte, error) {
	certBytes = bytes.TrimSpace(certBytes)
	if block, _ := pem.Decode(certBytes); block != nil {
		certBytes = block.Bytes
	}

	certificate, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return certificate.Signature, nil
}
