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

// naver ...
type naver struct {
	config  *oauth2.Config
	state   string
	name    string
	profile string
}

// NewNaver ...
func NewNaver(key, secret, callback string) Provider {
	return &naver{
		name:    "naver",
		state:   "",
		profile: "https://openapi.naver.com/v1/nid/me",
		config: &oauth2.Config{
			ClientID:     key,
			ClientSecret: secret,
			RedirectURL:  callback,
			Scopes:       []string{""},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://nid.naver.com/oauth2.0/authorize",
				TokenURL: "https://nid.naver.com/oauth2.0/token",
			},
		},
	}
}

func (f *naver) GetProfile(token string) (*ProviderUser, error) {
	return nil, nil
}

// GetAuthURL ...
func (n *naver) GetAuthURL() string {
	authURL, err := url.Parse(n.config.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", n.config.ClientID)
	parameters.Add("scope", strings.Join(n.config.Scopes, " "))
	parameters.Add("redirect_uri", n.config.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", n.state)
	authURL.RawQuery = parameters.Encode()
	return authURL.String()
}

// CallBack ...
func (n *naver) CallBack(r *http.Request) (*ProviderUser, error) {
	user := &ProviderUser{}
	state := r.FormValue("state")
	if state != n.state {
		return user, fmt.Errorf("invalid oauth state, expected '%s', got '%s'", n.state, state)
	}

	code := r.FormValue("code")

	token, err := n.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return user, fmt.Errorf("oauthConf.Exchange() failed with '%s'", err)
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", n.profile, nil)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return user, fmt.Errorf("Get: %s", err)
	}
	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return user, fmt.Errorf("ReadAll: %s", err)
	}

	err = n.userFromReader(bytes.NewReader(responseBytes), user)

	return user, err
}

// GetName ...
func (n *naver) GetName() string {
	return n.name
}

func (n *naver) userFromReader(reader io.Reader, user *ProviderUser) error {
	u := struct {
		Response struct {
			ID           string
			Nickname     string
			Name         string
			Email        string
			Gender       string
			Age          string
			Birthday     string
			ProfileImage string `json:"profile_image"`
		}
	}{}

	if err := json.NewDecoder(reader).Decode(&u); err != nil {
		return err
	}
	r := u.Response
	user.Provider = n.name
	user.Email = r.Email
	user.Name = r.Name
	user.NickName = r.Nickname
	user.AvatarURL = r.ProfileImage
	user.AuthCode = r.ID
	user.Description = fmt.Sprintf(`{"gender":"%s","age":"%s","birthday":"%s"}`, r.Gender, r.Age, r.Birthday)

	return nil
}
