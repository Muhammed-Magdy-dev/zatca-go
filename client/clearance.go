package client

import (
	"context"
	"net/http"

	"github.com/Muhammed-Magdy-dev/zatca-go/internal/httpx"
)

const clearanceInvoicePath = "/invoices/clearance/single"

func (c *Client) SendClearanceInvoice(
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
		c.BaseURL+clearanceInvoicePath,
		payload,
		headers,
		&out,
	)

	if err != nil {
		return nil, err
	}

	return out, nil
}
