package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Muhammed-Magdy-dev/zatca-go/internal/httpx"
)

const compliancePath = "/compliance"

type ComplianceRequest struct {
	CSR string `json:"csr"`
}

type ComplianceResponse struct {
	BinarySecurityToken string `json:"binarySecurityToken"`
	Secret              string `json:"secret"`
	RequestID           int    `json:"requestId"`
	DispositionMessage  string `json:"dispositionMessage"`

	CSR                   string `json:"csr,omitempty"`
	CSRPEM                string `json:"csr_pem,omitempty"`
	PrivateKeyPEM         string `json:"private_key_pem,omitempty"`
	SupportedInvoiceTypes string `json:"supported_invoice_types,omitempty"`
}

func (c *Client) Compliance(
	ctx context.Context,
	input *ComplianceRequest,
	otp string,
) (*ComplianceResponse, error) {
	if c == nil {
		return nil, fmt.Errorf("zatca client is nil")
	}

	if input == nil || input.CSR == "" {
		return nil, fmt.Errorf("zatca: csr is required")
	}

	headers := map[string]string{
		"OTP":            otp,
		"Accept-Version": "V2",
		"Content-Type":   "application/json",
	}

	var res ComplianceResponse

	if err := httpx.DoRequest(
		ctx,
		http.MethodPost,
		c.BaseURL+compliancePath,
		input,
		headers,
		&res,
	); err != nil {
		return nil, err
	}

	if res.BinarySecurityToken == "" || res.Secret == "" {
		return nil, fmt.Errorf("zatca: invalid response from ZATCA")
	}

	return &res, nil
}
