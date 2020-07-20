package oauth

import (
	"fmt"
	"net/http"
	"time"
)

// Provider ...
type Provider interface {
	GetAuthURL() string
	CallBack(r *http.Request) (*ProviderUser, error)
	GetProfile(token string) (*ProviderUser, error)
	GetName() string
}

// Providers ...
type Providers map[string]Provider

var providers = Providers{}

// ProviderUser ...
type ProviderUser struct {
	RawData           map[string]interface{}
	Provider          string `json:"provider"`
	Email             string `json:"email"`
	Name              string `json:"name"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	NickName          string `json:"nick_name"`
	Description       string `json:"description"`
	AuthCode          string `json:"auth_code"`
	AvatarURL         string `json:"avatar_url"`
	Location          string `json:"location"`
	AccessToken       string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
	RefreshToken      string `json:"refresh_token"`
	ExpiresAt         time.Time
}

// GetName 는 제공된 Name이 없으면 NickName 을 넣어준다.
// NickName 없으면 공백
func (p *ProviderUser) GetName() string {
	if len(p.Name) == 0 {
		return p.NickName
	}
	return p.Name
}

// GetNickName ...
func (p *ProviderUser) GetNickName() string {
	if len(p.NickName) == 0 {
		return p.Name
	}
	return p.NickName
}

// UseProviders ...
func UseProviders(p ...Provider) {
	for _, v := range p {
		providers[v.GetName()] = v
	}
}

// GetProvider ...
func GetProvider(name string) (Provider, error) {
	provider := providers[name]
	if provider == nil {
		return nil, fmt.Errorf("no provider for %s exists", name)
	}
	return provider, nil
}
