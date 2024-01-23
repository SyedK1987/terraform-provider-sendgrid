package sendgrid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type DomainAuthDns struct {
	DKIM1        Domainauthdnsrecord `json:"dkim1,omitempty"`
	DKIM2        Domainauthdnsrecord `json:"dkim2,omitempty"`
	MailCNAME    Domainauthdnsrecord `json:"mail_cname,omitempty"`
	DKIM         Domainauthdnsrecord `json:"dkim,omitempty"`
	MailServer   Domainauthdnsrecord `json:"mail_server,omitempty"`
	SubDomainSPF Domainauthdnsrecord `json:"subdomain_spf,omitempty"`
}

type Domainauthdnsrecord struct {
	Host  string `json:"host,omitempty"`
	Type  string `json:"type,omitempty"`
	Data  string `json:"data,omitempty"`
	Valid bool   `json:"valid,omitempty"`
}

type DomainAuthSubuser struct {
	Username string `json:"username,omitempty"`
	UserID   int64  `json:"user_id,omitempty"`
}

type DomainAuth struct {
	ID            int64               `json:"id,omitempty"`
	UserId        int64               `json:"user_id,omitempty"`
	Domain        string              `json:"domain,omitempty"`
	Subdomain     string              `json:"subdomain,omitempty"`
	CustomDKIM    string              `json:"custom_dkim_selector,omitempty"`
	Username      string              `json:"username,omitempty"`
	Ips           []string            `json:"ips,omitempty"`
	CustomSPF     bool                `json:"custom_spf"`
	Defaultdomain bool                `json:"default"`
	Legacy        bool                `json:"legacy,omitempty"`
	Valid         bool                `json:"valid"`
	DNSDetails    DomainAuthDns       `json:"dns,omitempty"`
	Subusers      []DomainAuthSubuser `json:"subusers,omitempty"`
}

func (c *Client) GetDomainAuth(ctx context.Context, domainid DomainAuth) (*DomainAuth, error) {

	//return nil, fmt.Errorf("domainauth not found:%+v", domainid)

	respBody, _, err := c.Get(ctx, "GET", "/whitelabel/domains")
	if err != nil {
		return nil, err
	}

	var domainauth []DomainAuth
	err = json.Unmarshal([]byte(respBody), &domainauth)
	if err != nil {
		return nil, err
	}

	//return nil, fmt.Errorf("domainauth not found:%+v", domainauth)
	for _, domain := range domainauth {
		if domain.ID == domainid.ID {
			if len(domain.Subusers) == 0 {
				domain.Subusers = []DomainAuthSubuser{
					{
						Username: "",
						UserID:   0,
					},
				}
			}
			return &domain, nil
		}
	}

	return nil, fmt.Errorf("domainauth not found:%+v", domainid.ID)
}

func (c *Client) CreateDomainAuth(ctx context.Context, domainauth DomainAuth) (*DomainAuth, error) {

	createrespBody, _, err := c.Post(ctx, "POST", "/whitelabel/domains", DomainAuth{
		Domain:        domainauth.Domain,
		CustomDKIM:    domainauth.CustomDKIM,
		CustomSPF:     domainauth.CustomSPF,
		Legacy:        domainauth.Legacy,
		Ips:           domainauth.Ips,
		Subdomain:     domainauth.Subdomain,
		Defaultdomain: domainauth.Defaultdomain,
		//Username:      domainauth.Username,
	})
	if err != nil {
		return nil, err
	}

	var domainauthresp DomainAuth
	err = json.Unmarshal([]byte(createrespBody), &domainauthresp)
	if err != nil {
		return nil, fmt.Errorf("createdomainauth: domain creation failed:%s", err.Error())
	}

	return c.GetDomainAuth(ctx, domainauthresp)
	//return &domainauthresp, nil
}

func (c *Client) ValidateDomainAuth(ctx context.Context, domainauth DomainAuth) (*DomainAuth, error) {

	validdomain, statuscode, err := c.Post(ctx, "POST", "/whitelabel/domains/"+fmt.Sprintf("%d", domainauth.ID)+"/validate", nil)
	if err != nil && statuscode != http.StatusOK {
		return nil, fmt.Errorf("domain validation failed:%s", err.Error())
	}

	var validatedomainresp DomainAuth
	err = json.Unmarshal([]byte(validdomain), &validatedomainresp)
	if err != nil {
		return nil, fmt.Errorf("domain validation failed:%s", err.Error())
	}
	return &validatedomainresp, nil
}

func (c *Client) UpdateDomainAuth(ctx context.Context, updatedetails DomainAuth) (*DomainAuth, error) {

	updatedomaindets, statuscode, err := c.Post(ctx, "PATCH", "/whitelabel/domains/"+fmt.Sprintf("%d", updatedetails.ID), DomainAuth{
		Defaultdomain: updatedetails.Defaultdomain,
		CustomSPF:     updatedetails.CustomSPF,
	})
	if err != nil && statuscode != http.StatusOK {
		return nil, fmt.Errorf("updatedomainauth: domain update failed:%s", err.Error())
	}

	var updatedomainresp DomainAuth
	err = json.Unmarshal([]byte(updatedomaindets), &updatedomainresp)
	if err != nil {
		return nil, fmt.Errorf("updatedomainauth: domain update failed:%s", err.Error())
	}
	return &updatedomainresp, nil
}

func (c *Client) DeleteDomainAuth(ctx context.Context, domainid string) (bool, error) {
	if domainid == "" {
		return false, fmt.Errorf("domainid is empty")
	}

	respBody, statusCode, err := c.Get(ctx, "DELETE", "/whitelabel/domains/"+domainid)
	if err != nil {
		return false, err
	}

	if statusCode >= http.StatusMultipleChoices && statusCode != http.StatusNotFound {
		return false, fmt.Errorf("Error Failed Deleting Domain." + strconv.Itoa(statusCode) + "," + respBody)
	}

	return true, nil
}

func (c *Client) CreateDomainAuthSubuser(ctx context.Context, domainauth DomainAuth) (*DomainAuth, error) {

	createrespBody, _, err := c.Post(ctx, "POST", "/whitelabel/domains/"+fmt.Sprintf("%d", domainauth.ID)+"/subuser", DomainAuth{
		Username: domainauth.Username,
	})
	if err != nil {
		return nil, err
	}

	var domainauthresp DomainAuth
	err = json.Unmarshal([]byte(createrespBody), &domainauthresp)
	if err != nil {
		return nil, fmt.Errorf("associatesubuser: subuser association failed:%s", err.Error())
	}

	getuserdetails, _ := c.GetSubuser(ctx, Subuser{Username: domainauth.Username})
	if err != nil {
		return nil, fmt.Errorf("associatesubuser: Failed to retrieve subuser details:%s", err.Error())
	}

	domainauthresp.Subusers = []DomainAuthSubuser{
		{
			Username: getuserdetails.Username,
			UserID:   getuserdetails.ID,
		},
	}
	domainauthresp.Username = getuserdetails.Username

	//return c.GetDomainAuth(ctx, domainauthresp)
	return &domainauthresp, nil
	//return nil, fmt.Errorf("associatesubuser: subuser association failed:%v", domainauthresp)
}

func (c *Client) GetDomainSubuser(ctx context.Context, domainid DomainAuth) (*DomainAuth, error) {

	//return nil, fmt.Errorf("domainauth not found:%+v", domainid)

	respBody, _, err := c.Get(ctx, "GET", "/whitelabel/domains")
	if err != nil {
		return nil, err
	}

	var domainauth []DomainAuth
	err = json.Unmarshal([]byte(respBody), &domainauth)
	if err != nil {
		return nil, fmt.Errorf("getdomainsubuser: domain subuser parsing failed:%s", err.Error())
	}

	// return nil, fmt.Errorf("domainauth not found:%+v", domainauth)
	for _, domain := range domainauth {
		if domain.ID == domainid.ID {
			if len(domain.Subusers) > 0 {
				for _, userinlist := range domain.Subusers {
					if userinlist.Username == domainid.Username {
						domain.Subusers = []DomainAuthSubuser{
							{
								Username: userinlist.Username,
								UserID:   userinlist.UserID,
							},
						}
						//return nil, fmt.Errorf("domainauthsubuser:domainauth already exists:%+v", domainid)
						return &domain, nil
					}
				}
			} else {
				domain.Subusers = []DomainAuthSubuser{
					{
						Username: domainid.Username,
						UserID:   domain.UserId,
					},
				}
				return &domain, nil
			}
		}
	}

	return nil, fmt.Errorf("domainauthsubuser:domainauth not found:%+v", domainid.ID)
}

func (c *Client) DeleteDomainAuthSubuser(ctx context.Context, domainauth DomainAuth) (bool, error) {

	_, _, err := c.Get(ctx, "DELETE", "/whitelabel/domains/subuser?username="+domainauth.Username)
	if err != nil {
		return false, fmt.Errorf("removesubuserfromdomain: subuser disassociation failed:%s", err.Error())
	}

	//return c.GetDomainAuth(ctx, domainauthresp)
	return true, nil
}
