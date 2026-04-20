# zatca-go

Go SDK for ZATCA e-invoicing integration (Saudi Arabia)

---

## 🚀 Features

* CSR Generation (supports Standard / Simplified / Both)
* Compliance API (OTP onboarding)
* Reporting API (Simplified invoices)
* Clearance API (Standard invoices)
* Production CSID
* Invoice Signing (XAdES)
* QR Code Generation (TLV)

---

## 📦 Installation

```bash
go get github.com/Muhammed-Magdy-dev/zatca-go@latest
```

---

## ⚡ Quick Start

### Initialize Client

```go
client := client.New("https://gw-fatoora.zatca.gov.sa/e-invoicing/simulation")
```

---

# 🔄 Full Flow

---

## 1️⃣ Compliance (Onboarding)

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

## 🧠 CSR InvoiceType (VERY IMPORTANT)

| Value | Description                         |
| ----- | ----------------------------------- |
| 1100  | Standard + Simplified (Recommended) |
| 1000  | Standard only                       |
| 0100  | Simplified only                     |

Example:

```go
csrInput := &csr.CSRInput{
	VATNumber:        "...",
	OrganizationName: "...",
	BranchName:       "...",
	SectorName:       "...",
	Address:          "...",
	CountryCode:      "SA",
	InvoiceType:      "1100",
}
```

---

## 2️⃣ Build Invoice

```go
input := &invoice.InvoiceInput{
	ID:        "reporting:1.0",
	UUID:      uuid.NewString(),
	IssueDate: time.Now(),
	ICV:       1,

	InvoiceTypeCode: "388",

	// 👇 IMPORTANT
	IsSimplified: true, // true = simplified, false = standard

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

## 🧾 Standard vs Simplified

### ✔️ Simplified (B2C)

```go
input.IsSimplified = true
```

* Uses Reporting API
* Customer NOT required

---

### ✔️ Standard (B2B)

```go
input.IsSimplified = false
```

* Uses Clearance API
* Customer REQUIRED

```go
input.Customer = &invoice.CustomerInfo{
	VATNumber:        "123456789012345",
	RegistrationName: "Customer Name",
	Street:           "Riyadh",
	City:             "Riyadh",
	PostalCode:       "12345",
}
```

---

## ⚠️ Certificate Handling (VERY IMPORTANT)

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

## 3️⃣ Sign Invoice

```go
finalXML, err := invoice.PrepareSignedInvoice(
	input,
	[]byte(privateKeyPEM),
	certDER,
)
```

SDK handles:

* XML generation
* Signature
* SignedProperties
* QR Code

---

## 4️⃣ Send Invoice

### 🔹 Simplified (Reporting)

```go
res, err := client.SendReportingInvoice(
	ctx,
	binarySecurityToken,
	secret,
	payload,
)
```

---

### 🔹 Standard (Clearance)

```go
res, err := client.SendClearanceInvoice(
	ctx,
	binarySecurityToken,
	secret,
	payload,
)
```

---

## 5️⃣ Production CSID

```go
res, err := client.GetProductionCSID(
	ctx,
	binarySecurityToken,
	secret,
	requestID,
)
```

---

# 🔥 Flow Summary

1. Generate CSR
2. Compliance (OTP)
3. Save credentials
4. Build invoice
5. Generate Certificate Hash
6. Sign invoice
7. Send (Reporting / Clearance)
8. Move to Production

---

# ⚠️ Critical Rules

* `CertificateHash` must be EXACT
* Do NOT modify signed XML
* Standard requires Customer
* Simplified must NOT include Customer
* Always send previous invoice hash
* Use correct endpoint

---

# ❌ Common Mistakes

* Missing CertificateHash
* Using wrong Base64 encoding
* Sending unsigned XML
* Wrong InvoiceTypeCode
* Sending Standard via Reporting
* Sending Simplified via Clearance
* Customer missing in Standard

---

# 🧠 Best Practice

Always use:

```go
InvoiceType: "1100"
```

during onboarding to support all invoice types.

---

# 📄 License

MIT
