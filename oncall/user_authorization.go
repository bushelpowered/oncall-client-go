package oncall

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// This type implements the http.RoundTripper interface
type UserAuthorizationRoundTripper struct {
	Proxied        http.RoundTripper
	LoginEndpoint  string
	UsernameGetter func() string
	PasswordGetter func() string
	csrfToken      *string
	cookieJar      http.CookieJar
}

func NewUserAuthorizationRoundTripper(src UserAuthorizationRoundTripper) UserAuthorizationRoundTripper {
	if src.cookieJar == nil {
		src.cookieJar, _ = cookiejar.New(nil)
	}
	src.csrfToken = new(string)
	return src
}

func (uart UserAuthorizationRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	log.Tracef("Going to roundtrip user auth for %s %s", req.Method, req.URL)
	if uart.PasswordGetter() != "" {
		csrfToken, err := uart.GetCSRFToken()
		if err != nil {
			e = errors.Wrap(err, "Getting CSRF Token")
			return
		}
		req.Header.Set("X-CSRF-TOKEN", csrfToken)
		for _, c := range uart.cookieJar.Cookies(req.URL) {
			req.AddCookie(c)
		}
	} else {
		log.Debug("Password not set, not going to set auth as user")
	}

	// Send the request, get the response (or the error)
	return uart.Proxied.RoundTrip(req)
}

// We don't actually login, just set the csrfToken to empty and it'll login again
func (uart UserAuthorizationRoundTripper) Login() error {
	*uart.csrfToken = ""
	return nil
}

func (uart UserAuthorizationRoundTripper) GetCSRFToken() (string, error) {
	var err error
	if uart.csrfToken == nil {
		return "", errors.New("csrfToken Pointer is nil. Please use the NewUserAuthorizationRoundTripper function")
	}

	if *uart.csrfToken != "" {
		log.Trace("Using existing csrf token")
		return *uart.csrfToken, nil
	}
	log.Debug("Getting new CSRF token")

	if uart.cookieJar == nil {
		return "", errors.New("Cookie jar is not set for User auth roundtripper. Please use the NewUserAuthorizationRoundTripper function")
	}

	client := &http.Client{
		Transport: uart.Proxied,
		Jar:       uart.cookieJar,
	}

	resp, err := client.PostForm(uart.LoginEndpoint,
		url.Values{
			"username": {uart.UsernameGetter()},
			"password": {uart.PasswordGetter()},
		},
	)
	if err != nil {
		return "", errors.Wrapf(err, "Logging into %s as %s", uart.LoginEndpoint, uart.UsernameGetter())
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read body while logging in")
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Failed to login (%d)", resp.StatusCode)
	}

	loginResponse := struct {
		CsrfToken string `json:"csrf_token"`
		God       int    `json:"god"`
	}{}

	err = json.Unmarshal(bodyBytes, &loginResponse)
	if err != nil {
		log.Tracef("Failed to unmarshal body: %s", string(bodyBytes))
		return "", errors.Wrap(err, "Failed to parse login JSON response")
	}

	*uart.csrfToken = loginResponse.CsrfToken
	return *uart.csrfToken, nil
}
