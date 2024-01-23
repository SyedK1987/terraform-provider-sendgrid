package sendgrid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type User struct {
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	// Address   string   `json:"address,omitempty"`
	// Address2  string   `json:"address2,omitempty"`
	// City      string   `json:"city,omitempty"`
	// State     string   `json:"state,omitempty"`
	// Zip       string   `json:"zip,omitempty"`
	// Country   string   `json:"country,omitempty"`
	// Company   string   `json:"company,omitempty"`
	// Website   string   `json:"website,omitempty"`
	// Phone     string   `json:"phone,omitempty"`
	IsAdmin bool `json:"is_admin,omitempty"`
	//IsSSO    bool     `json:"is_sso,omitempty"`
	UserType       string   `json:"user_type,omitempty"`
	Scopes         []string `json:"scopes,omitempty"`
	IsReadOnly     bool     `json:"is_read_only,omitempty"`
	ExpirationDate int64    `json:"expiration_date,omitempty"`
	Token          string   `json:"token,omitempty"`
}

type Users struct {
	Result []User `json:"result"`
}

// type PendingUser struct {
// 	Result []struct {
// 		Token          string   `json:"token,omitempty"`
// 		Email          string   `json:"email,omitempty"`
// 		IsAdmin        bool     `json:"is_admin,omitempty"`
// 		IsReadOnly     bool     `json:"is_read_only,omitempty"`
// 		ExpirationDate int      `json:"expiration_date,omitempty"`
// 		Scopes         []string `json:"scopes,omitempty"`
// 	} `json:"result"`
// }

func parseUser(respBody string) (*User, error) {
	var body User

	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, fmt.Errorf("failed parsing teammate: %w", err)
	}

	return &body, nil
}

func (c *Client) GetUsernameByEmail(ctx context.Context, email string) (string, error) {
	respBody, _, err := c.Get(ctx, "GET", "/teammates?limit=10000")
	if err != nil {
		return "", err
	}

	users := &Users{}

	decoder := json.NewDecoder(bytes.NewReader([]byte(respBody)))
	err = decoder.Decode(users)
	if err != nil {
		return "", err
	}

	for _, user := range users.Result {
		if user.Email == email && user.Username != "" {
			return user.Username, nil
		}
	}

	//return "", fmt.Errorf("username with email %s not found", email)
	return "", nil
}

func (c *Client) ReadUser(ctx context.Context, email string) (*User, error) {
	username, err := c.GetUsernameByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	respBody, _, err := c.Get(ctx, "GET", "/teammates/"+username)
	if err != nil {
		return nil, err
	}

	var u User
	err = json.Unmarshal([]byte(respBody), &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *Client) CreateTeammate(ctx context.Context, user User) (*User, error) {

	//_, err := c.RefreshTeammate(ctx, user.Email)
	//if err != nil {
	respBody, _, err := c.Post(ctx, "POST", "/teammates", user)
	if err != nil {
		return nil, fmt.Errorf("CreateUser: Bad Request:" + err.Error())
	}

	var body User

	err = json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, fmt.Errorf("CreateTeammate: failed parsing teammate: %w", err)
	}

	getcompletedetails, _, err := c.Get(ctx, "GET", "/teammates/pending")
	if err != nil {
		return nil, fmt.Errorf("CreateTeammate: unable to retrive invited user details: %w", err)
	}

	pendingUsers := &Users{}
	decoder1 := json.NewDecoder(bytes.NewReader([]byte(getcompletedetails)))
	err = decoder1.Decode(pendingUsers)
	if err != nil {
		return nil, err
	}

	for _, pu := range pendingUsers.Result {
		if pu.Email == user.Email {
			return &pu, nil
		}
	}

	//	return parseUser(respBody)
	//} else {
	return nil, fmt.Errorf("unable to retrive user data from pendinglist")
	//}
}

func (c *Client) RefreshTeammate(ctx context.Context, email string) (*User, error) {

	respBody, _, err := c.Get(ctx, "GET", "/teammates?limit=10000")
	if err != nil {
		return nil, err
	}

	users := &Users{}

	decoder := json.NewDecoder(bytes.NewReader([]byte(respBody)))
	err = decoder.Decode(users)
	if err != nil {
		return nil, err
	}

	for _, user := range users.Result {
		if user.Email == email && user.Username != "" {
			respBody, _, err := c.Get(ctx, "GET", "/teammates/"+user.Username)
			if err != nil {
				return nil, err
			}

			var u User
			err = json.Unmarshal([]byte(respBody), &u)
			if err != nil {
				return nil, err
			}
			return &u, nil
		}
	}

	respBody1, _, err := c.Get(ctx, "GET", "/teammates/pending")
	if err != nil {
		return nil, err
	}

	pendingUsers := &Users{}
	decoder1 := json.NewDecoder(bytes.NewReader([]byte(respBody1)))
	err = decoder1.Decode(pendingUsers)
	if err != nil {
		return nil, err
	}

	for _, pu := range pendingUsers.Result {
		if pu.Email == email {
			//golan build struct
			//var p User
			//return nil, fmt.Errorf("pending user with email not found:%+v", &pu)
			return &pu, nil
		}
	}

	return nil, fmt.Errorf("username with email %s not found", email)
}

// func (c *Client) GetPendingUser(ctx context.Context, email string) (string, error) {
// 	respBody, _, err := c.Get(ctx, "GET", "/teammates/pending")
// 	if err != nil {
// 		return "", err
// 	}

// 	pendingUsers := &Users{}
// 	decoder := json.NewDecoder(bytes.NewReader([]byte(respBody)))
// 	err = decoder.Decode(pendingUsers)
// 	if err != nil {
// 		return "", err
// 	}

// 	for _, user := range pendingUsers.Result {
// 		if user.Email == email {
// 			return user.Token, nil
// 		}
// 	}
// 	return "", fmt.Errorf("pending user with email %s not found", email)
// }

func (c *Client) UpdateTeammate(ctx context.Context, updateitems User) (*User, error) {

	username, err := c.RefreshTeammate(ctx, updateitems.Email)
	if err != nil || username.Username == "" {
		return nil, fmt.Errorf("User with email %s ", updateitems.Email+" not found, Please accept the invite first")
	}

	respBody1, _, err := c.Post(ctx, "PATCH", "/teammates/"+updateitems.Username, User{
		IsAdmin: updateitems.IsAdmin,
		Scopes:  updateitems.Scopes,
	})
	if err != nil {
		return nil, err
	}

	var body1 User

	err = json.Unmarshal([]byte(respBody1), &body1)
	if err != nil {
		return nil, fmt.Errorf("failed parsing teammate: %w", err)
	}

	return parseUser(respBody1)
}

func (c *Client) DeleteTeammate(ctx context.Context, email string) (bool, error) {

	getusername, err := c.RefreshTeammate(ctx, email)
	if err != nil {
		return false, fmt.Errorf("DeleteTeammate: User with email %s not found", email)
	}

	if getusername.Username != "" && getusername.Token == "" {
		_, _, err := c.Get(ctx, "DELETE", "/teammates/"+getusername.Username)
		if err != nil {
			return false, fmt.Errorf("DeleteTeammate: Failed to Delete:" + err.Error())
		}
		return true, nil
	} else {
		_, _, err := c.Get(ctx, "DELETE", "/teammates/pending/"+getusername.Token)
		if err != nil {
			return false, fmt.Errorf("DeleteTeammate: Failed to Delete:" + err.Error())
		}
		return true, nil
	}
}

// func (c *Client) ResendTmate(ctx context.Context, gettoekn string) (*User, error) {

// 	// token, err := c.GetPendingUser(ctx, email)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	respBody, _, err := c.Post(ctx, "POST", "/teammates/pending/"+gettoekn+"/resend", nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return c.RefreshTeammate(ctx, respBody)
// }
