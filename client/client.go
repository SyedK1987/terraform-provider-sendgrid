package sendgrid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
)

const HostURL string = "https://api.sendgrid.com/v3"

type Client struct {
	ApiKey     string
	HTTPClient *http.Client
}

func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("clientgo: apikey is required")
	}
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		ApiKey:     apiKey,
	}
	return &c, nil
}

func bodyToJSON(body interface{}) ([]byte, error) {
	if body == nil {
		return nil, fmt.Errorf("clientgo: body could not be jsonified")
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("clietgo: body could not be jsonified: %w", err)
	}

	return jsonBody, nil
}

func (c *Client) Get(ctx context.Context, method rest.Method, endpoint string) (string, int, error) {

	var req rest.Request
	req = sendgrid.GetRequest(c.ApiKey, endpoint, HostURL)
	req.Method = method

	resp, err := sendgrid.API(req)
	if err != nil || resp.StatusCode >= 400 {
		return "", resp.StatusCode, fmt.Errorf("clientgetfunc:api response: http %d: %s, err: %v", resp.StatusCode, resp.Body, err)
	}

	return resp.Body, resp.StatusCode, nil

}

func (c *Client) Post(ctx context.Context, method rest.Method, endpoint string, body interface{}) (string, int, error) {
	var err error

	var req rest.Request
	req = sendgrid.GetRequest(c.ApiKey, endpoint, HostURL)
	req.Method = method

	if body != nil {
		req.Body, err = bodyToJSON(body)
	}

	if err != nil {
		return "", 0, fmt.Errorf("ClientGo: Failed preparing request body: %w", err)
	}

	resp, err := sendgrid.API(req)

	if err != nil || resp.StatusCode >= 400 {
		return "", resp.StatusCode, fmt.Errorf("api response: http %d: %s, err: %v", resp.StatusCode, resp.Body)
	}

	if err != nil {
		return "", 0, fmt.Errorf("clientgo: api post func error: %v", err)
	}

	return resp.Body, resp.StatusCode, nil
}
