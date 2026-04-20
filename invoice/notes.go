package invoice

import (
	"time"

	"github.com/google/uuid"
)

func CreateCreditNote(original *InvoiceInput) *InvoiceInput {
	return &InvoiceInput{
		ID:                  original.ID,
		UUID:                uuid.NewString(),
		IssueDate:           time.Now(),
		ICV:                 original.ICV + 1,
		InvoiceTypeCode:     "381",
		BillingReferenceID:  original.UUID,
		PreviousInvoiceHash: original.InvoiceDigest,
		Supplier:            original.Supplier,
		PaymentMeansCode:    original.PaymentMeansCode,
		Lines:               original.Lines,
	}
}

func CreateDebitNote(original *InvoiceInput) *InvoiceInput {
	return &InvoiceInput{
		ID:                  original.ID,
		UUID:                uuid.NewString(),
		IssueDate:           time.Now(),
		ICV:                 original.ICV + 1,
		InvoiceTypeCode:     "383",
		BillingReferenceID:  original.UUID,
		PreviousInvoiceHash: original.InvoiceDigest,
		Supplier:            original.Supplier,
		PaymentMeansCode:    original.PaymentMeansCode,
		Lines:               original.Lines,
	}
}
