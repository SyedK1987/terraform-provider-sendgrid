package sendgrid

import (
	"context"
	"encoding/json"
	"fmt"
)

// Singlesender is the model for the singlesender resource.
type Singlesender struct {
	Nickname    string `json:"nickname"`
	FromEmail   string `json:"from_email"`
	FromName    string `json:"from_name"`
	ReplyTo     string `json:"reply_to"`
	ReplyToName string `json:"reply_to_name"`
	Address     string `json:"address"`
	Address2    string `json:"address2"`
	State       string `json:"state"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Zip         string `json:"zip"`
	ID          int64  `json:"id"`
}

type ReturnSinglesender struct {
	Nickname    string `json:"nickname"`
	FromEmail   string `json:"from_email"`
	FromName    string `json:"from_name"`
	ReplyTo     string `json:"reply_to"`
	ReplyToName string `json:"reply_to_name"`
	Address     string `json:"address"`
	Address2    string `json:"address2"`
	State       string `json:"state"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Zip         string `json:"zip"`
	Verified    bool   `json:"verified"`
	Locked      bool   `json:"locked"`
	ID          int64  `json:"id"`
}

type SinglesenderResult struct {
	Result []ReturnSinglesender `json:"results"`
}

// SinglesenderCreate creates a new singlesender.
func (c *Client) CreateSingleSender(ctx context.Context, receiveditems Singlesender) (*ReturnSinglesender, error) {

	respBody, _, err := c.Post(ctx, "POST", "/verified_senders", receiveditems)
	if err != nil {
		return nil, fmt.Errorf("CreateSingleSender: Bad Request:" + err.Error())
	}

	var response ReturnSinglesender
	err = json.Unmarshal([]byte(respBody), &response)
	if err != nil {
		return nil, fmt.Errorf("CreateSingleSender: failed parsing singlesender: %w", err)
	}

	return &response, nil
}

// SinglesenderRead reads a singlesender.
func (c *Client) ReadSingleSender(ctx context.Context, id string) (*ReturnSinglesender, error) {

	var response SinglesenderResult
	var convertedsinglesender ReturnSinglesender

	respBody, _, err := c.Get(ctx, "GET", "/verified_senders")
	if err != nil {
		return nil, fmt.Errorf("ReadSingleSender: Bad Request:" + err.Error())
	}

	err = json.Unmarshal([]byte(respBody), &response)
	if err != nil {
		return nil, fmt.Errorf("ReadSingleSender: failed parsing singlesender: %w", err)
	}

	for _, item := range response.Result {
		if id == fmt.Sprintf("%d", item.ID) {
			convertedsinglesender = item
		}
	}

	return &convertedsinglesender, nil
}

// Singlesenderdelete updates a singlesender.
func (c *Client) DeleteSingleSender(ctx context.Context, id string) (bool, error) {

	respBody, statusCode, err := c.Get(ctx, "DELETE", "/verified_senders/"+id)
	if err != nil {
		return false, fmt.Errorf("DeleteSingleSender: Bad Request:" + err.Error())
	}

	if respBody == "" && statusCode == 204 {
		return true, nil
	}

	return false, nil
}

func (c *Client) UpdateSingleSender(ctx context.Context, updateitem Singlesender) (*ReturnSinglesender, error) {

	tmpitem := fmt.Sprintf("%d", updateitem.ID)
	updaterespBody, _, err := c.Post(ctx, "PATCH", "/verified_senders/"+tmpitem, updateitem)
	if err != nil {
		return nil, fmt.Errorf("UpdateSingleSender: Bad Request:" + err.Error())
	}

	var updateresponse ReturnSinglesender
	err = json.Unmarshal([]byte(updaterespBody), &updateresponse)
	if err != nil {
		return nil, fmt.Errorf("UpdateSingleSender: failed parsing singlesender: %w", err)
	}

	return &updateresponse, nil
}
