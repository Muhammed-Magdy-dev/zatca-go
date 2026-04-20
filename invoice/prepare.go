package invoice

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/Muhammed-Magdy-dev/zatca-go/cert"
	"github.com/Muhammed-Magdy-dev/zatca-go/sign"
)

func PrepareSignedInvoice(input *InvoiceInput, privateKeyPEM []byte, certPEM []byte) ([]byte, error) {
	if input == nil {
		return nil, fmt.Errorf("input is nil")
	}

	certificate, err := cert.ParseCertificate(certPEM)
	if err != nil {
		return nil, err
	}

	input.IssuerName = certificate.Issuer.ToRDNSequence().String()
	input.SerialNumber = certificate.SerialNumber.String()
	input.SigningTime = input.IssueDate.UTC().Format("2006-01-02T15:04:05")

	pubDER, err := x509.MarshalPKIXPublicKey(certificate.PublicKey)
	if err != nil {
		return nil, err
	}

	input.PublicKeyBytes = pubDER
	input.CertificateSignatureBytes = certificate.Signature
	input.X509Certificate = base64.StdEncoding.EncodeToString(certificate.Raw)

	input.InvoiceDigest = ""
	input.SignedPropsDigest = ""
	input.SignatureValue = ""
	input.QRCode = ""

	unsignedXML, err := BuildInvoiceXML(input)
	if err != nil {
		return nil, err
	}

	invoiceDigest, err := sign.HashInvoiceXML(unsignedXML)
	if err != nil {
		return nil, fmt.Errorf("invoice digest: %w", err)
	}
	input.InvoiceDigest = invoiceDigest

	signedPropsDigest, err := HashSignedProperties(input)
	if err != nil {
		return nil, fmt.Errorf("signed properties digest: %w", err)
	}
	input.SignedPropsDigest = signedPropsDigest

	xmlToSign, err := BuildInvoiceXML(input)
	if err != nil {
		return nil, err
	}

	signatureValue, err := sign.BuildSignatureValue(xmlToSign, privateKeyPEM)
	if err != nil {
		return nil, err
	}
	input.SignatureValue = signatureValue

	qr, err := BuildQRCodeTLVBase64(input)
	if err != nil {
		return nil, err
	}
	input.QRCode = qr

	finalXML, err := BuildInvoiceXML(input)
	if err != nil {
		return nil, err
	}

	return finalXML, nil
}

func HashSignedProperties(input *InvoiceInput) (string, error) {
	signedPropsXML, err := BuildSignedPropertiesXML(input)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(signedPropsXML)
	hexHash := hex.EncodeToString(sum[:])

	return base64.StdEncoding.EncodeToString([]byte(hexHash)), nil
}
