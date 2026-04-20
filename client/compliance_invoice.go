package client

import (
	"context"
	"net/http"

	"github.com/Muhammed-Magdy-dev/zatca-go/internal/httpx"
)

const complianceInvoicePath = "/compliance/invoices"

type ComplianceInvoiceRequest struct {
	InvoiceHash string `json:"invoiceHash"`
	UUID        string `json:"uuid"`
	Invoice     string `json:"invoice"`
}

func (c *Client) SendComplianceInvoice(
	ctx context.Context,
	binarySecurityToken string,
	secret string,
	payload *ComplianceInvoiceRequest,
) (map[string]any, error) {
	headers := map[string]string{
		"Accept-Version":  "V2",
		"Accept-Language": "en",
		"Authorization":   buildBasicAuth(binarySecurityToken, secret),
		"Content-Type":    "application/json",
	}

	var out map[string]any

	err := httpx.DoRequest(
		ctx,
		http.MethodPost,
		c.BaseURL+complianceInvoicePath,
		payload,
		headers,
		&out,
	)
	if err != nil {
		return nil, err
	}

	return out, nil
}
