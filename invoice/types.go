package invoice

import "time"

type InvoiceInput struct {
	ID                        string
	UUID                      string
	SigningTime               string
	IssueDate                 time.Time
	ICV                       int
	PublicKeyBase64           string
	CertificateSignature      string
	PublicKeyBytes            []byte
	CertificateSignatureBytes []byte
	PreviousInvoiceHash       string
	Supplier                  SupplierInfo
	Customer                  *Party
	PaymentMeansCode          string
	Lines                     []InvoiceLine
	InvoiceLevelACs           []AllowanceCharge
	IsSimplified              bool
	QRCode                    string
	CertificateHash           string
	IssuerName                string
	SerialNumber              string
	X509Certificate           string
	InvoiceDigest             string
	SignedPropsDigest         string
	SignatureValue            string
	BillingReferenceID        string
	InvoiceTypeCode           string
	InstructionNote           string
}

type AllowanceCharge struct {
	Indicator       bool
	Reason          string
	ReasonCode      string
	Amount          float64
	VATRate         float64
	TaxCategoryCode string
}

type SupplierInfo struct {
	CRN              string
	VATNumber        string
	RegistrationName string
	Street           string
	BuildingNumber   string
	PlotID           string
	District         string
	City             string
	PostalCode       string
}

type InvoiceLine struct {
	ID            string
	Name          string
	Quantity      float64
	UnitCode      string
	Price         float64
	VATRate       float64
	ItemDiscounts []AllowanceCharge
}
