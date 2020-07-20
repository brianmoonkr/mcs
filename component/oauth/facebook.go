package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

// facebook ...
type facebook struct {
	config  *oauth2.Config
	state   string
	name    string
	profile string
}

// NewFacebook ...
func NewFacebook(key, secret, callback string) Provider {
	return &facebook{
		name:    "facebook",
		state:   "",
		profile: "https://graph.facebook.com/me?fields=email,first_name,last_name,link,about,id,name,picture,location&access_token=",
		config: &oauth2.Config{
			ClientID:     key,
			ClientSecret: secret,
			RedirectURL:  callback,
			Scopes:       []string{"public_profile", "email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.facebook.com/dialog/oauth",
				TokenURL: "https://graph.facebook.com/oauth/access_token",
			},
		},
	}
}

// GetAuthURL ...
func (f *facebook) GetAuthURL() string {
	authURL, err := url.Parse(f.config.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", f.config.ClientID)
	parameters.Add("scope", strings.Join(f.config.Scopes, " "))
	parameters.Add("redirect_uri", f.config.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", f.state)
	authURL.RawQuery = parameters.Encode()
	return authURL.String()
}

func (f *facebook) GetProfile(token string) (*ProviderUser, error) {
	return nil, nil
}

// CallBack ...
func (f *facebook) CallBack(r *http.Request) (*ProviderUser, error) {
	user := &ProviderUser{}
	state := r.FormValue("state")
	if state != f.state {
		return user, fmt.Errorf("invalid oauth state, expected '%s', got '%s'", f.state, state)
	}

	code := r.FormValue("code")

	token, err := f.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return user, fmt.Errorf("oauthConf.Exchange() failed with '%s'", err)
	}

	resp, err := http.Get(f.profile + url.QueryEscape(token.AccessToken))
	if err != nil {
		return user, fmt.Errorf("Get: %s", err)
	}
	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return user, fmt.Errorf("ReadAll: %s", err)
	}

	err = f.userFromReader(bytes.NewReader(responseBytes), user)
	return user, err
}

// GetName ...
func (f *facebook) GetName() string {
	return f.name
}

func (f *facebook) userFromReader(reader io.Reader, user *ProviderUser) error {
	u := struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		About     string `json:"about"`
		Name      string `json:"name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Link      string `json:"link"`
		Picture   struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
		Location struct {
			Name string `json:"name"`
		} `json:"location"`
	}{}

	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		return err
	}

	user.Provider = f.name
	user.Name = u.Name
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.NickName = u.Name
	user.Email = u.Email
	user.Description = u.About
	user.AvatarURL = u.Picture.Data.URL
	user.AuthCode = u.ID
	user.Location = u.Location.Name

	return err
}
