package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/Muhammed-Magdy-dev/zatca-go/internal/httpx"
)

const productionCSIDPath = "/production/csids"

type ProductionCSIDResponse struct {
	BinarySecurityToken string `json:"binarySecurityToken"`
	Secret              string `json:"secret"`
	RequestID           int    `json:"requestId"`
}

func buildBasicAuth(binarySecurityToken, secret string) string {
	raw := binarySecurityToken + ":" + secret
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(raw))
}

func (c *Client) GetProductionCSID(
	ctx context.Context,
	complianceToken string,
	complianceSecret string,
	complianceRequestID int64,
) (*ProductionCSIDResponse, error) {
	if c == nil {
		return nil, fmt.Errorf("zatca client is nil")
	}

	headers := map[string]string{
		"Authorization":  buildBasicAuth(complianceToken, complianceSecret),
		"Accept-Version": "V2",
		"Content-Type":   "application/json",
	}

	var res ProductionCSIDResponse

	err := httpx.DoRequest(
		ctx,
		http.MethodPost,
		c.BaseURL+productionCSIDPath,
		map[string]any{
			"compliance_request_id": complianceRequestID,
		},
		headers,
		&res,
	)
	if err != nil {
		return nil, err
	}

	if res.BinarySecurityToken == "" || res.Secret == "" {
		return nil, fmt.Errorf("invalid production csid response")
	}

	return &res, nil
}
