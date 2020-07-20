package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/SermoDigital/jose"
	"golang.org/x/oauth2"
)

type line struct {
	config  *oauth2.Config
	state   string
	name    string
	profile string
}

// NewLine ...
func NewLine(key, secret, callback string) Provider {
	return &line{
		name:    "line",
		state:   "1234abcde",
		profile: "https://api.line.me/v2/profile",
		config: &oauth2.Config{
			ClientID:     key,
			ClientSecret: secret,
			RedirectURL:  callback,
			Scopes:       []string{"profile", "openid", "email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
				TokenURL: "https://api.line.me/oauth2/v2.1/token",
			},
		},
	}
}

// GetAuthURL ...
func (l *line) GetAuthURL() string {
	authURL, err := url.Parse(l.config.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", l.config.ClientID)
	parameters.Add("scope", strings.Join(l.config.Scopes, " "))
	parameters.Add("redirect_uri", l.config.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", l.state)
	authURL.RawQuery = parameters.Encode()
	return authURL.String()
}

func (f *line) GetProfile(token string) (*ProviderUser, error) {
	return nil, nil
}

// CallBack ...
func (l *line) CallBack(r *http.Request) (*ProviderUser, error) {
	user := &ProviderUser{}
	state := r.FormValue("state")
	if state != l.state {
		return user, fmt.Errorf("invalid oauth state, expected '%s', got '%s'", l.state, state)
	}

	code := r.FormValue("code")

	token, err := l.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return user, fmt.Errorf("oauthConf.Exchange() failed with '%s'", err)
	}

	tokenBytes, _ := jose.Base64Decode([]byte(strings.Split(token.Extra("id_token").(string), ".")[1]))

	err = l.userFromReader(bytes.NewReader(tokenBytes), user)

	return user, err
}

// GetName ...
func (l *line) GetName() string {
	return l.name
}

// userFromReader ...
func (l *line) userFromReader(reader io.Reader, user *ProviderUser) error {
	u := struct {
		ID      string `json:"sub"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}{}
	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		return err
	}

	user.Provider = l.name
	user.NickName = u.Name
	user.Email = u.Email
	user.AvatarURL = u.Picture
	user.AuthCode = u.ID

	return nil
}
