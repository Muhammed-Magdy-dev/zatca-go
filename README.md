# zatca-go

Go SDK for ZATCA e-invoicing integration (Saudi Arabia)

---

## Features

* CSR Generation
* Compliance API (OTP onboarding)
* Reporting API (Reporting Single)
* Clearance API
* Production CSID
* Invoice Signing (XAdES)
* QR Code Generation (TLV)

---

## Installation

```bash
go get github.com/Muhammed-Magdy-dev/zatca-go@latest
```

---

## Quick Start

### 1) Initialize Client

```go
client := client.New("https://gw-fatoora.zatca.gov.sa/e-invoicing/simulation")
```

---

## Full Flow

### Step 1: Compliance (Onboarding)

```go
res, err := client.Compliance(ctx, &client.ComplianceRequest{
	CSR: csrBase64,
}, otp)
```

Save:

* BinarySecurityToken
* Secret
* RequestID

---

## Step 2: Build Invoice

```go
input := &invoice.InvoiceInput{
	ID:        "reporting:1.0",
	UUID:      uuid.NewString(),
	IssueDate: time.Now(),
	ICV:       1,

	InvoiceTypeCode: "388",

	Supplier: invoice.SupplierInfo{
		VATNumber:        "123456789012345",
		RegistrationName: "My Company",
		CRN:              "1234567890",
		Street:           "Riyadh",
		City:             "Riyadh",
		PostalCode:       "12345",
	},

	PaymentMeansCode: "10",

	Lines: []invoice.InvoiceLine{
		{
			ID:       "1",
			Name:     "Item",
			Quantity: 1,
			UnitCode: "PCE",
			Price:    100,
			VATRate:  15,
		},
	},
}
```

---

## ⚠️ Step 3: Certificate Handling (VERY IMPORTANT)

You MUST prepare certificate and hash manually before signing.

```go
certDER, err := cert.TokenToCertificateDER(binarySecurityToken)
if err != nil {
	return err
}

certBytes, err := base64.StdEncoding.DecodeString(binarySecurityToken)
if err != nil {
	return err
}

hash := sha256.Sum256(certBytes)
hexStr := hex.EncodeToString(hash[:])

input.CertificateHash = base64.StdEncoding.EncodeToString([]byte(hexStr))
```

---

## Step 4: Sign Invoice

```go
finalXML, err := invoice.PrepareSignedInvoice(
	input,
	[]byte(privateKeyPEM),
	certDER,
)
```

SDK will handle:

* XML building
* SignedProperties
* SignatureValue
* QR Code

---

## Step 5: Send Reporting Invoice

```go
payload := &client.ComplianceInvoiceRequest{
	InvoiceHash: input.InvoiceDigest,
	UUID:        input.UUID,
	Invoice:     base64.StdEncoding.EncodeToString(finalXML),
}

res, err := client.SendReportingInvoice(
	ctx,
	binarySecurityToken,
	secret,
	payload,
)
```

---

## Step 6: Production CSID

```go
res, err := client.GetProductionCSID(
	ctx,
	binarySecurityToken,
	secret,
	requestID,
)
```

---

## Flow Summary

1. Generate CSR
2. Call Compliance (OTP)
3. Save credentials
4. Build invoice
5. Generate Certificate Hash (manual step)
6. Sign invoice
7. Send Reporting / Clearance
8. Move to Production

---

## Important Notes

* CertificateHash must be calculated exactly as shown
* Do NOT modify signed XML manually
* Always use previous invoice hash
* Use Simulation before Production
* OTP is required only once during onboarding

---

## Common Mistakes

❌ Forgetting CertificateHash
❌ Using wrong Base64 encoding
❌ Sending unsigned XML
❌ Using wrong InvoiceTypeCode

---

## License

MIT
