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
	"time"

	"golang.org/x/oauth2"
)

type google struct {
	config  *oauth2.Config
	state   string
	name    string
	profile string
}

// NewGoogle ...
func NewGoogle(key, secret, callback string) Provider {
	return &google{
		name:    "google",
		state:   "",
		profile: "https://www.googleapis.com/oauth2/v3/userinfo?access_token=",
		config: &oauth2.Config{
			ClientID:     key,
			ClientSecret: secret,
			RedirectURL:  callback,
			Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://accounts.google.com/o/oauth2/auth",
				TokenURL:  "https://accounts.google.com/o/oauth2/token",
				AuthStyle: oauth2.AuthStyleInParams,
			},
		},
	}
}

// GetAuthURL ...
func (g *google) GetAuthURL() string {
	authURL, err := url.Parse(g.config.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", g.config.ClientID)
	parameters.Add("scope", strings.Join(g.config.Scopes, " "))
	parameters.Add("redirect_uri", g.config.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", g.state)
	authURL.RawQuery = parameters.Encode()
	return authURL.String()
}

// CallBack ...
func (g *google) CallBack(r *http.Request) (*ProviderUser, error) {
	state := r.FormValue("state")
	if state != g.state {
		return nil, fmt.Errorf("invalid oauth state, expected '%s', got '%s'", g.state, state)
	}

	code := r.FormValue("code")

	token, err := g.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("oauthConf.Exchange() failed with '%s'", err)
	}

	return g.GetProfile(token.AccessToken)
}

// GetProfile ...
func (g *google) GetProfile(token string) (*ProviderUser, error) {
	user := &ProviderUser{}
	client := http.Client{
		Timeout: time.Duration(3 * time.Second),
	}
	resp, err := client.Get(g.profile + url.QueryEscape(token))
	if err != nil {
		return nil, fmt.Errorf("Get: %s", err)
	}
	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return user, fmt.Errorf("ReadAll: %s", err)
	}

	//tglog.Logger.Infof("Google =====> responseBytes : %s\n", string(responseBytes))

	err = g.userFromReader(bytes.NewReader(responseBytes), user)
	return user, err
}

func (g *google) userFromReader(reader io.Reader, user *ProviderUser) error {
	u := struct {
		ID        string `json:"sub"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		FirstName string `json:"given_name"`
		LastName  string `json:"family_name"`
		Link      string `json:"link"`
		Picture   string `json:"picture"`
	}{}

	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		return err
	}

	user.Provider = g.name
	user.Name = u.Name
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.NickName = u.Name
	user.Email = u.Email
	user.AvatarURL = u.Picture
	user.AuthCode = u.ID

	return nil
}

// GetName ...
func (g *google) GetName() string {
	return g.name
}
