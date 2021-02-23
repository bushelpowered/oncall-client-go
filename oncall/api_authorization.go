package oncall

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// This type implements the http.RoundTripper interface
type APIAuthorizationRoundTripper struct {
	Proxied        http.RoundTripper
	UsernameGetter func() string
	PasswordGetter func() string
}

func (art APIAuthorizationRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	var err error

	if art.PasswordGetter() != "" {
		hmacTime := time.Now().Unix() / 5
		hmacMethod := req.Method
		hmacPath := req.URL.Path
		hmacBody := []byte{}

		if req.Body != nil {
			hmacBody, err = ioutil.ReadAll(req.Body)
			if err != nil {
				e = errors.Wrap(err, "Failed to read request body for hmac generation")
				return
			}
		}

		hmacData := fmt.Sprintf("%d %s %s %s", hmacTime, hmacMethod, hmacPath, string(hmacBody))
		log.Debugf("Setting auth header using this data: %s", hmacData)

		hmacSum := hmac512(art.PasswordGetter(), hmacData)

		req.Header.Set("Authorization", fmt.Sprintf("hmac %s:%s", art.UsernameGetter(), hmacSum))
		log.Tracef("Set auth header to: %s", req.Header.Get("Authorization"))
	} else {
		log.Debug("Password not set, not going to set Auth header")
	}

	// Send the request, get the response (or the error)
	return art.Proxied.RoundTrip(req)
}

func hmac512(secret, data string) string {
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha512.New, []byte(secret))

	// Write Data to it
	h.Write([]byte(data))

	// Get result and encode as hexadecimal string
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return sha
}
