package sendgrid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type LinkAuthDns struct {
	DCNAME Linkdnsrecord `json:"domain_cname,omitempty"`
	OCNAME Linkdnsrecord `json:"owner_cname,omitempty"`
}

type Linkdnsrecord struct {
	Host  string `json:"host,omitempty"`
	Type  string `json:"type,omitempty"`
	Data  string `json:"data,omitempty"`
	Valid bool   `json:"valid"`
}

type LinkAuth struct {
	ID            int64       `json:"id,omitempty"`
	UserId        int64       `json:"user_id,omitempty"`
	Domain        string      `json:"domain,omitempty"`
	Subdomain     string      `json:"subdomain,omitempty"`
	Username      string      `json:"username,omitempty"`
	Defaultdomain bool        `json:"default"`
	Legacy        bool        `json:"legacy"`
	Valid         bool        `json:"valid"`
	DNSDetails    LinkAuthDns `json:"dns,omitempty"`
}

func (c *Client) CreateLinkBrand(ctx context.Context, domainauth LinkAuth) (*LinkAuth, error) {

	createrespBody, _, err := c.Post(ctx, "POST", "/whitelabel/links", LinkAuth{
		Domain:        domainauth.Domain,
		Subdomain:     domainauth.Subdomain,
		Defaultdomain: domainauth.Defaultdomain,
	})
	if err != nil {
		return nil, err
	}

	var linkauthresp LinkAuth
	err = json.Unmarshal([]byte(createrespBody), &linkauthresp)
	if err != nil {
		return nil, fmt.Errorf("createlinkbrand: link branding creation failed:%s", err.Error())
	}

	//return c.GetDomainAuth(ctx, domainauthresp)
	return &linkauthresp, nil
}

func (c *Client) Getlinkbrand(ctx context.Context, domainid LinkAuth) (*LinkAuth, error) {

	//return nil, fmt.Errorf("domainauth not found:%+v", domainid)

	respBody, _, err := c.Get(ctx, "GET", "/whitelabel/links/"+fmt.Sprintf("%d", domainid.ID))
	if err != nil {
		return nil, err
	}

	var domainauth LinkAuth
	err = json.Unmarshal([]byte(respBody), &domainauth)
	if err != nil {
		return nil, err
	}

	return &domainauth, nil
}

func (c *Client) Updatelinkbrand(ctx context.Context, updatedetails LinkAuth) (*LinkAuth, error) {
	updatebranddets, statuscode, err := c.Post(ctx, "PATCH", "/whitelabel/links/"+fmt.Sprintf("%d", updatedetails.ID), DomainAuth{
		Defaultdomain: updatedetails.Defaultdomain,
	})
	if err != nil && statuscode != http.StatusOK {
		return nil, fmt.Errorf("updatelinkbrand: default domain update failed:%s", err.Error())
	}

	var updatebrandresp LinkAuth
	err = json.Unmarshal([]byte(updatebranddets), &updatebrandresp)
	if err != nil {
		return nil, fmt.Errorf("updatedomainauth: domain update failed:%s", err.Error())
	}
	return &updatebrandresp, nil
}

func (c *Client) Validatelinkbrand(ctx context.Context, validatedetails LinkAuth) (*LinkAuth, error) {

	//	return nil, fmt.Errorf("sk you are here")
	validbrand, statuscode, err := c.Post(ctx, "POST", "/whitelabel/links/"+fmt.Sprintf("%d", validatedetails.ID)+"/validate", nil)
	if err != nil && statuscode != http.StatusOK {
		return nil, fmt.Errorf("updatelinkbrand: link branding validation failed:%s", err.Error())
	}
	var validatebrandresp LinkAuth
	err = json.Unmarshal([]byte(validbrand), &validatebrandresp)
	if err != nil {
		return nil, fmt.Errorf("updatelinkbrand: Unable to unmarshal data:%s", err.Error())
	}

	return c.Getlinkbrand(ctx, validatebrandresp)
}

func (c *Client) Deletelinkbrand(ctx context.Context, domainid string) (bool, error) {
	if domainid == "" {
		return false, fmt.Errorf("linkbrand id is empty")
	}

	respBody, statusCode, err := c.Get(ctx, "DELETE", "/whitelabel/links/"+domainid)
	if err != nil {
		return false, err
	}

	if statusCode >= http.StatusMultipleChoices && statusCode != http.StatusNotFound {
		return false, fmt.Errorf("Error Failed Deleting Link Brand." + strconv.Itoa(statusCode) + "," + respBody)
	}

	return true, nil
}
