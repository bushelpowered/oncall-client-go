package oncall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Config struct {
	// Username is the username you authenticate to oncall with
	Username string
	// Password is the password for the associated username
	Password string
	// Endpoint is everything before "/api/v0/" in your url:
	// e.g.: https://example.com/oncall/api/v0/teams/
	// would be: https://example.com/oncall/
	// (Trailing slash is optional)
	Endpoint string
	// You can auth using either the API or user auth.
	// The API is limited when it can do using API auth
	AuthMethod AuthMethod
}

type AuthMethod string

const (
	AuthMethodAPI  AuthMethod = "api"
	AuthMethodUser            = "user"
)

type Client struct {
	Client           *http.Client
	Config           Config
	authRoundTripper AuthRoundtripper
}

type AuthRoundtripper interface {
	Login() error
	RoundTrip(req *http.Request) (res *http.Response, e error)
}

// New creates a new oncall client.
// client arg can be nil, which will default to http.DefaultClient
// config should be populated with username, passsword, and endpoint
func New(client *http.Client, config Config, logger LeveledLogger) (*Client, error) {
	if logger != nil {
		log = logger
	}

	if config.Endpoint == "" {
		return nil, errors.New("You must define at least an endpoint")
	}

	oncallClient := &Client{
		Config: config,
	}
	// Strip off the trailing slash if it's there
	oncallClient.Config.Endpoint = strings.TrimRight(oncallClient.Config.Endpoint, "/")

	if client == nil {
		client = http.DefaultClient
	}

	proxiedTransport := client.Transport
	if proxiedTransport == nil {
		proxiedTransport = http.DefaultTransport
	}

	if config.AuthMethod == AuthMethodUser {
		log.Debug("Using User AuthMethod")
		oncallClient.authRoundTripper = NewUserAuthorizationRoundTripper(UserAuthorizationRoundTripper{
			Proxied:        proxiedTransport,
			UsernameGetter: func() string { return oncallClient.Config.Username },
			PasswordGetter: func() string { return oncallClient.Config.Password },
			LoginEndpoint:  strings.TrimRight(config.Endpoint, "/") + "/login",
		})
	} else {
		log.Debug("Using API AuthMethod")
		oncallClient.authRoundTripper = APIAuthorizationRoundTripper{
			Proxied:        proxiedTransport,
			UsernameGetter: func() string { return oncallClient.Config.Username },
			PasswordGetter: func() string { return oncallClient.Config.Password },
		}
	}
	client.Transport = oncallClient.authRoundTripper

	oncallClient.Client = client

	return oncallClient, nil
}

// Request receives a result which, if not nil, will then json unmarshal the respone into
// It will also return the body bytes of the response
func (c *Client) Request(method string, path string, body string, result interface{}) ([]byte, error) {
	bodyReader := bytes.NewReader([]byte(body))
	req, err := http.NewRequest(method, c.Config.Endpoint+"/"+strings.TrimLeft(path, "/"), bodyReader)
	if err != nil {
		return []byte{}, errors.Wrap(err, "Failed to create new request")
	}

	var resp *http.Response

	doRequest := func() (*http.Response, []byte, error) {
		log.Tracef("Going to do request: %s %s", req.Method, req.URL)
		resp, err = c.Client.Do(req)
		if err != nil {
			return resp, []byte{}, errors.Wrap(err, "Failed to do http request")
		}
		defer resp.Body.Close()

		// read the body before checking status
		// This way we can use the bodyBytes in the error message
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		return resp, bodyBytes, errors.Wrap(err, "Failed to read response body")
	}

	resp, bodyBytes, err := doRequest()
	if err != nil {
		return bodyBytes, errors.Wrap(err, "Failed to do request")
	}

	if resp.StatusCode == 401 {
		log.Debug("Going to re-login due to 401")
		var err error
		for i := 0; i < 3; i++ {
			err = c.authRoundTripper.Login()
			if err == nil {
				break
			}
			time.Sleep(1 * time.Second)
		}
		if err != nil {
			return []byte{}, errors.Wrap(err, "Failed to login the auth roundtripper")
		}

		resp, bodyBytes, err = doRequest()
		if err != nil {
			return bodyBytes, errors.Wrap(err, "Failed to do request after re-login")
		}
	}

	if resp.StatusCode >= 400 {
		log.Debugf("Dump of body on error (%d) (%s %s): %s", resp.StatusCode, req.Method, req.URL, string(bodyBytes))
		return bodyBytes, fmt.Errorf("HTTP Request failed (%d) (%s)", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		err = json.Unmarshal(bodyBytes, result)
		if err != nil {
			log.Debugf("Dump of body on json error: %s", string(bodyBytes))
		}
	}
	return bodyBytes, errors.Wrap(err, "JSON Unmarshal Error")
}

func (c *Client) PayloadRequest(method string, path string, body interface{}, result interface{}) ([]byte, error) {
	var reqBody string

	if bodyAsString, ok := body.(string); ok {
		reqBody = bodyAsString
	} else {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return []byte{}, errors.Wrap(err, "Failed to json marshal body")
		}
		reqBody = string(jsonBody)
	}

	return c.Request(method, path, reqBody, result)
}

func (c *Client) Get(path string, result interface{}) ([]byte, error) {
	return c.Request("GET", path, "", result)
}

func (c *Client) Post(path string, body interface{}, result interface{}) ([]byte, error) {
	return c.PayloadRequest("POST", path, body, result)
}

func (c *Client) Put(path string, body interface{}, result interface{}) ([]byte, error) {
	return c.PayloadRequest("PUT", path, body, result)
}

func (c *Client) Delete(path string, body interface{}, result interface{}) ([]byte, error) {
	return c.PayloadRequest("DELETE", path, body, result)
}
