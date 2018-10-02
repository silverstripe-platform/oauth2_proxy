package providers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testSilverStripeProvider(hostname string) *SilverStripeProvider {
	p := NewSilverStripeProvider(
		&ProviderData{
			LoginURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/sso/authorize",
			},
			RedeemURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/sso/accessToken",
			},
			ValidateURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/sso/profile",
			},
		},
	)
	if hostname != "" {
		updateURL(p.Data().LoginURL, hostname)
		updateURL(p.Data().RedeemURL, hostname)
		updateURL(p.Data().ValidateURL, hostname)
	}
	return p
}

func testSilverStripeBackend(payload string) *httptest.Server {
	path := "/sso/profile"

	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			url := r.URL
			if url.Path != path {
				w.WriteHeader(404)
			} else {
				w.WriteHeader(200)
				w.Write([]byte(payload))
			}
		},
	))
}

func TestSilverStripeProviderGetEmailAddress(t *testing.T) {
	b := testSilverStripeBackend("{\"email\": \"someone@example.com\"}")
	defer b.Close()

	url, _ := url.Parse(b.URL)
	p := testSilverStripeProvider(url.Host)

	session := &SessionState{AccessToken: "imaginary_access_token"}
	email, err := p.GetEmailAddress(session)
	assert.Equal(t, nil, err)
	assert.Equal(t, "someone@example.com", email)
}

func TestSilverStripeProviderGetEmailAddressFailedRequest(t *testing.T) {
	b := testSilverStripeBackend("unused payload")
	defer b.Close()

	url, _ := url.Parse(b.URL)
	p := testSilverStripeProvider(url.Host)

	session := &SessionState{AccessToken: "unexpected_access_token"}
	email, err := p.GetEmailAddress(session)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", email)
}

func TestSilverStripeProviderGetEmailAddressEmailNotPresentInPayload(t *testing.T) {
	b := testSilverStripeBackend("{\"foo\": \"bar\"}")
	defer b.Close()

	url, _ := url.Parse(b.URL)
	p := testSilverStripeProvider(url.Host)

	session := &SessionState{AccessToken: "imaginary_access_token"}
	email, err := p.GetEmailAddress(session)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", email)
}

func TestSilverStripeProviderGetUserName(t *testing.T) {
	b := testSilverStripeBackend("{\"username\": \"someone\"}")
	defer b.Close()

	url, _ := url.Parse(b.URL)
	p := testSilverStripeProvider(url.Host)

	session := &SessionState{AccessToken: "imaginary_access_token"}
	email, err := p.GetUserName(session)
	assert.Equal(t, nil, err)
	assert.Equal(t, "someone", email)
}
