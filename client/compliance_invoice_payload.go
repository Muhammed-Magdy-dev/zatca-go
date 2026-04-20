package client

type BuiltInvoice struct {
	InvoiceHash string
	UUID        string
	Base64XML   string
}

func BuildCompliancePayload(b *BuiltInvoice) *ComplianceInvoiceRequest {
	if b == nil {
		return nil
	}

	return &ComplianceInvoiceRequest{
		InvoiceHash: b.InvoiceHash,
		UUID:        b.UUID,
		Invoice:     b.Base64XML,
	}
}
