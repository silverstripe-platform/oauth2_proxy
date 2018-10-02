package providers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/silverstripeltd/oauth2_proxy/api"
)

type SilverStripeProvider struct {
	*ProviderData
}

func NewSilverStripeProvider(p *ProviderData) *SilverStripeProvider {
	p.ProviderName = "SilverStripe"
	if p.Scope == "" {
		p.Scope = "read_profile"
	}
	return &SilverStripeProvider{ProviderData: p}
}

func (p *SilverStripeProvider) getProfileField(s *SessionState, field string) (string, error) {
	req, err := http.NewRequest("GET", p.ValidateURL.String(), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))
	if err != nil {
		log.Printf("failed building request %s", err)
		return "", err
	}
	json, err := api.Request(req)
	if err != nil {
		log.Printf("failed making request %s", err)
		return "", err
	}
	return json.Get(field).String()
}

func (p *SilverStripeProvider) GetEmailAddress(s *SessionState) (string, error) {
	return p.getProfileField(s, "email")
}

func (p *SilverStripeProvider) GetUserName(s *SessionState) (string, error) {
	return p.getProfileField(s, "username")
}
