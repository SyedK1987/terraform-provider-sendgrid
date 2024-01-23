package sendgrid

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Ipmgmt struct {
	CreatedAt int    `json:"created_at,omitempty"`
	ID        int64  `json:"id,omitempty"`
	IP        string `json:"ip,omitempty"`
	UpdatedAt int    `json:"updated_at,omitempty"`
}

type Ips struct {
	IPS []Ipmgmt `json:"ips,omitempty"`
}

type IPsresult struct {
	Ips []Ipmgmt `json:"result,omitempty"`
}

type Result struct {
	Result Ipmgmt `json:"result"`
}

func (c *Client) GetIPMgmt(ctx context.Context, ipmgmtid string) (*Result, error) {

	var getipResult Result

	if ipmgmtid == "" {
		return nil, fmt.Errorf("GetIPMgmt: ipmgmtid is empty")
	}

	if !validIP4(ipmgmtid) {

		respBody, _, err := c.Get(ctx, "GET", "/access_settings/whitelist/"+ipmgmtid)
		if err != nil {
			return nil, errors.New("Read Data[GetByIP]: Bad Request:" + err.Error())
		}

		err = json.Unmarshal([]byte(respBody), &getipResult)
		if err != nil {
			return nil, fmt.Errorf("GetIPMgmt: failed parsing ipmgmt: %w", err)
		}
	} else {
		respBody, _, err := c.Get(ctx, "GET", "/access_settings/whitelist")
		if err != nil {
			return nil, errors.New("getallips: Bad Request:" + err.Error())
		}

		var getipList IPsresult
		err = json.Unmarshal([]byte(respBody), &getipList)
		if err != nil {
			return nil, fmt.Errorf("getallips: failed parsing ipmgmt: %w", err)
		}

		for _, ipitem := range getipList.Ips {
			tempItem := ipitem.IP
			//tflog.Debug(ctx, "COMPARE", map[string]interface{}{"ipmgmtid": ipmgmtid, "tempItem": tempItem})
			if tempItem == ipmgmtid {
				getipResult.Result = ipitem
			}
		}
	}

	return &getipResult, nil

}

func (c *Client) CreateIPMgmt(ctx context.Context, ips Ipmgmt) (*Ipmgmt, error) {

	//return nil, errors.New(string(receivedips))

	collectedips := []Ipmgmt{}
	collectedips = append(collectedips, ips)
	respBody, statusCode, err := c.Post(ctx, "POST", "/access_settings/whitelist", Ips{IPS: collectedips})

	if err != nil && statusCode == 400 {
		return nil, errors.New("CreateIPMgMt: Bad Request:" + strconv.Itoa(statusCode) + "," + err.Error())
	}

	getResult := IPsresult{}

	err = json.Unmarshal([]byte(respBody), &getResult)
	if err != nil {
		return nil, fmt.Errorf("CreateIPMgmt: failed parsing ipmgmt: %w", err)
	}

	//return nil, fmt.Errorf("Error from Parsing: %+v", getResult)
	//return nil, errors.New(respBody)

	return &getResult.Ips[0], nil
}

func (c *Client) DeleteIPMgmt(ctx context.Context, ipmgmtid string) (bool, error) {
	if ipmgmtid == "" {
		return false, fmt.Errorf("ipmgmtid is empty")
	}

	//return false, errors.New("DeleteIPMgmt: Bad Request:" + ipmgmtid)

	respBody, statusCode, err := c.Get(ctx, "DELETE", "/access_settings/whitelist/"+ipmgmtid)

	if err != nil {
		return false, fmt.Errorf("DeleteIPMgmt: failed to delete resource: %w", err)
	}

	if respBody == "" && statusCode == 204 {
		return true, nil
	}

	return false, nil
}

func validIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")

	re, _ := regexp.Compile(`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-5][0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-5][0-5])`)
	// if re.MatchString(ipAddress) {
	// 	return true
	// }
	// return false
	return re.MatchString(ipAddress)
}
