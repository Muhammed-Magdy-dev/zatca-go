package httpx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/goccy/go-json"
)

func DoRequest(
	ctx context.Context,
	method string,
	url string,
	body any,
	headers map[string]string,
	out any,
) error {

	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := http.Client{
		Timeout: 60 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bd, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK &&
		res.StatusCode != http.StatusAccepted &&
		res.StatusCode != http.StatusPartialContent {

		return fmt.Errorf("zatca status %d: %s", res.StatusCode, string(bd))
	}

	if out != nil {
		return json.Unmarshal(bd, out)
	}

	return nil
}
