package sendgrid

import (
	"context"
	"encoding/json"
	"fmt"
)

type ChildApiKey struct {
	Name   string   `json:"name,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
	ID     string   `json:"api_key_id,omitempty"`
	Apikey string   `json:"api_key,omitempty"`
}

type CreateApikey struct {
	Name   string   `json:"name,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
}

func (c *Client) CreateApiKey(ctx context.Context, receiveditems ChildApiKey) (*ChildApiKey, error) {

	respBody, _, err := c.Post(ctx, "POST", "/api_keys", CreateApikey{
		Name:   receiveditems.Name,
		Scopes: receiveditems.Scopes,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateApiKey: Bad Request:" + err.Error())
	}

	var response ChildApiKey
	err = json.Unmarshal([]byte(respBody), &response)
	if err != nil {
		return nil, fmt.Errorf("CreateApiKey: failed parsing apikey: %w", err)
	}
	//return nil, fmt.Errorf("CreateApiKey: Bad Request:%+v", response)

	return &response, nil
}

func (c *Client) ReadApiKey(ctx context.Context, apikeyid string) (*ChildApiKey, error) {
	respBody, _, err := c.Get(ctx, "GET", "/api_keys/"+apikeyid)
	if err != nil {
		return nil, fmt.Errorf("ReadApiKey: Bad Request:" + err.Error())
	}

	var response ChildApiKey
	err = json.Unmarshal([]byte(respBody), &response)
	if err != nil {
		return nil, fmt.Errorf("ReadApiKey: failed parsing apikey: %w", err)
	}

	return &response, nil
}

func (c *Client) DeleteApiKey(ctx context.Context, apikeyid string) (bool, error) {
	delrespbody, statuscode, err := c.Get(ctx, "DELETE", "/api_keys/"+apikeyid)
	if err != nil {
		return false, fmt.Errorf("DeleteApiKey: Bad Request:" + err.Error())
	}

	if delrespbody == "" && statuscode == 204 {
		return true, nil
	}

	return false, nil
}

func (c *Client) UpdateApiKey(ctx context.Context, receivedupdateitems ChildApiKey) (*ChildApiKey, error) {

	updaterespBody, _, err := c.Post(ctx, "PUT", "/api_keys/"+receivedupdateitems.ID, ChildApiKey{
		Name:   receivedupdateitems.Name,
		Scopes: receivedupdateitems.Scopes,
	})
	if err != nil {
		return nil, fmt.Errorf("UpdateApiKey: Bad Request:" + err.Error())
	}

	var updateresponse ChildApiKey
	err = json.Unmarshal([]byte(updaterespBody), &updateresponse)
	if err != nil {
		return nil, fmt.Errorf("UpdateApiKey: failed parsing apikey: %w", err)
	}

	return &updateresponse, nil
}

func (c *Client) UpdateApiKeyName(ctx context.Context, nametoupdate ChildApiKey) (*ChildApiKey, error) {

	updaterespBody, _, err := c.Post(ctx, "PATCH", "/api_keys/"+nametoupdate.ID, ChildApiKey{
		Name: nametoupdate.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("UpdateApiKeyName: Bad Request:" + err.Error())
	}

	var updateresponse ChildApiKey
	err = json.Unmarshal([]byte(updaterespBody), &updateresponse)
	if err != nil {
		return nil, fmt.Errorf("UpdateApiKeyName: failed parsing apikey: %w", err)
	}
	//return nil, fmt.Errorf("UpdateApiKeyName: Bad Request:%+v", updateresponse)
	return &updateresponse, nil
}
