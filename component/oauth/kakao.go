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
	"strconv"
	"strings"

	"golang.org/x/oauth2"
)

// Kakao ...
type kakao struct {
	config  *oauth2.Config
	state   string
	name    string
	profile string
}

// NewKakao ...
func NewKakao(key, secret, callback string) Provider {
	return &kakao{
		name:    "kakao",
		state:   "",
		profile: "https://kapi.kakao.com/v2/user/me",
		config: &oauth2.Config{
			ClientID:     key,
			ClientSecret: secret,
			RedirectURL:  callback,
			Scopes:       []string{"profile", "account_email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://kauth.kakao.com/oauth/authorize",
				TokenURL: "https://kauth.kakao.com/oauth/token",
			},
		},
	}
}

func (f *kakao) GetProfile(token string) (*ProviderUser, error) {
	return nil, nil
}

// GetAuthURL ...
func (k *kakao) GetAuthURL() string {
	authURL, err := url.Parse(k.config.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", k.config.ClientID)
	parameters.Add("scope", strings.Join(k.config.Scopes, " "))
	parameters.Add("redirect_uri", k.config.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", k.state)
	authURL.RawQuery = parameters.Encode()
	return authURL.String()
}

// CallBack ...
func (k *kakao) CallBack(r *http.Request) (*ProviderUser, error) {
	user := &ProviderUser{}
	state := r.FormValue("state")
	if state != k.state {
		return user, fmt.Errorf("invalid oauth state, expected '%s', got '%s'", k.state, state)
	}

	code := r.FormValue("code")

	token, err := k.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return user, fmt.Errorf("oauthConf.Exchange() failed with '%s'", err)
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", k.profile, nil)
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

	//tglog.Logger.Infof("KaKao =====> response : %s\n", string(responseBytes))

	err = k.userFromReader(bytes.NewReader(responseBytes), user)

	return user, err
}

// GetName ...
func (k *kakao) GetName() string {
	return k.name
}

func (k *kakao) userFromReader(reader io.Reader, user *ProviderUser) error {
	u := struct {
		ID         int `json:"id"`
		Properties struct {
			NickName       string `json:"nickname"`
			ProfileImage   string `json:"profile_image"`
			ThumbnailImage string `json:"thumbnail_image"`
		} `json:"properties"`
		KakaoAccount struct {
			HasEmail        bool   `json:"has_email"`
			IsEmailValid    bool   `json:"is_email_valid"`
			IsEmailVerified bool   `json:"is_email_verified"`
			HasBirthday     bool   `json:"has_birthday"`
			HasGender       bool   `json:"has_gender"`
			Email           string `json:"email"`
		} `json:"kakao_account"`
	}{}

	if err := json.NewDecoder(reader).Decode(&u); err != nil {
		return err
	}

	//tglog.Logger.Infof("KaKao =====> KakaoAccount : %+v\n", u)

	user.Provider = k.name
	user.Email = u.KakaoAccount.Email
	user.NickName = u.Properties.NickName
	user.AvatarURL = u.Properties.ThumbnailImage
	user.AuthCode = strconv.Itoa(u.ID)
	return nil
}
