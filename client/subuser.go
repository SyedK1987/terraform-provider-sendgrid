package sendgrid

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

type Subuser struct {
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
	Email    string   `json:"email,omitempty"`
	Ips      []string `json:"ips,omitempty"`
	Disabled bool     `json:"disabled,omitempty"`
	ID       int64    `json:"id,omitempty"`
}

func (c *Client) CreateSubuser(ctx context.Context, subuser Subuser) (*Subuser, error) {

	createRespBody, _, err := c.Post(ctx, "POST", "/subusers", subuser)
	if err != nil {
		return nil, fmt.Errorf("CreateSubuser: Failed to Create:" + err.Error())
	}

	var body Subuser

	err = json.Unmarshal([]byte(createRespBody), &body)
	if err != nil {
		return nil, fmt.Errorf("CreateSubuser: Failed to Unmarshal:" + err.Error())
	}

	return c.GetSubuser(ctx, body)
}

func (c *Client) GetSubuser(ctx context.Context, userdata Subuser) (*Subuser, error) {

	getRespBody, _, err := c.Get(ctx, "GET", "/subusers/"+userdata.Username)
	if err != nil {
		return nil, fmt.Errorf("GetSubuser: Failed to Get userdata:" + err.Error())
	}

	var body Subuser

	err = json.Unmarshal([]byte(getRespBody), &body)
	if err != nil {
		return nil, fmt.Errorf("GetSubuser: Failed to Unmarshal:" + err.Error())
	}

	return &body, nil
}

func (c *Client) UpdateSubuser(ctx context.Context, userdata Subuser) (*Subuser, error) {

	_, statusCode, err := c.Post(ctx, "PATCH", "/subusers/"+userdata.Username, Subuser{
		Disabled: userdata.Disabled,
	})

	if err != nil {
		return nil, fmt.Errorf("failed updating subUser: " + err.Error() + ". StatusCode: " + strconv.Itoa(statusCode))
	}

	return c.GetSubuser(ctx, userdata)
}

func (c *Client) DeleteSubuser(ctx context.Context, userdata string) (bool, error) {

	_, statusCode, err := c.Get(ctx, "DELETE", "/subusers/"+userdata)

	if err != nil {
		return false, fmt.Errorf("failed deleting subUser: " + err.Error() + ". StatusCode: " + strconv.Itoa(statusCode))
	}

	return true, nil
}

func (c *Client) UpdateIp(ctx context.Context, uip Subuser) (*Subuser, error) {

	_, statusCode, err := c.Post(ctx, "PATCH", "/subusers/"+uip.Username, Subuser{
		Ips: uip.Ips,
	})

	if err != nil {
		return nil, fmt.Errorf("failed updating subUser IP: " + err.Error() + ". StatusCode: " + strconv.Itoa(statusCode))
	}

	return c.GetSubuser(ctx, uip)
}

func (c *Client) ReadSubuser(ctx context.Context, userdata string) (*Subuser, error) {

	getRespBody, _, err := c.Get(ctx, "GET", "/subusers/"+userdata)
	if err != nil {
		return nil, fmt.Errorf("GetSubuser: Failed to Get userdata:" + err.Error())
	}

	var body Subuser

	err = json.Unmarshal([]byte(getRespBody), &body)
	if err != nil {
		return nil, fmt.Errorf("GetSubuser: Failed to Unmarshal:" + err.Error())
	}

	return &body, nil
}
