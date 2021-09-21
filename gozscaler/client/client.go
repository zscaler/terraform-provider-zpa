package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/willguibr/terraform-provider-zpa/gozscaler"
)

type Client struct {
	Config *gozscaler.Config
}

// NewClient returns a new client for the specified apiKey.
func NewClient(config *gozscaler.Config) (c *Client) {
	if config == nil {
		config, _ = gozscaler.NewConfig("", "", "", "")
	}
	c = &Client{Config: config}
	return
}

func (client *Client) NewPrivateRequestDo(method, url string, options, body, v interface{}) (*http.Response, error) {
	return client.newRequestDoCustom(method, url, true, options, body, v)
}

func (client *Client) NewRequestDo(method, url string, options, body, v interface{}) (*http.Response, error) {
	return client.newRequestDoCustom(method, url, false, options, body, v)
}

func (client *Client) newRequestDoCustom(method, url string, usePrivateAPI bool, options, body, v interface{}) (*http.Response, error) {
	client.Config.Lock()
	defer client.Config.Unlock()
	if client.Config.AuthToken == nil || client.Config.AuthToken.AccessToken == "" {
		if client.Config.ClientID == "" || client.Config.ClientSecret == "" {
			log.Printf("[ERROR] No client credentials were provided. Please set %s, %s and %s enviroment variables.\n", gozscaler.ZPA_CLIENT_ID, gozscaler.ZPA_CLIENT_SECRET, gozscaler.ZPA_CUSTOMER_ID)
			return nil, errors.New("no client credentials were provided")
		}
		log.Printf("[TRACE] Getting access token for %s=%s\n", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID)
		formData := []byte(fmt.Sprintf("client_id=%s&client_secret=%s",
			client.Config.ClientID,
			client.Config.ClientSecret,
		))

		req, err := http.NewRequest("POST", client.Config.BaseURL.String()+"/signin", bytes.NewBuffer(formData))
		if err != nil {
			log.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID, err)
			return nil, err
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Config.GetHTTPClient().Do(req)

		if err != nil {
			log.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID, err)
			return nil, err
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID, err)
			return nil, err
		}

		var a gozscaler.AuthToken
		err = json.Unmarshal(respBody, &a)
		if err != nil {
			log.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID, err)
			return nil, err
		}

		// we need keep auth token for future http request
		client.Config.AuthToken = &a
	}
	req, err := client.newRequest(method, url, usePrivateAPI, options, body)
	if err != nil {
		return nil, err
	}
	client.logRequest(req)
	return client.do(req, v)
}

// Generating the Http request
func (client *Client) newRequest(method, urlPath string, usePrivateAPI bool, options, body interface{}) (*http.Request, error) {
	if client.Config.AuthToken == nil || client.Config.AuthToken.AccessToken == "" {
		log.Printf("[ERROR] Failed to signin the user %s=%s\n", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID)
		return nil, fmt.Errorf("failed to signin the user %s=%s", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID)
	}
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	// Join the path to the base-url
	u := *client.Config.BaseURL
	if usePrivateAPI {
		u = *client.Config.PrivateAPIBaseURL
	}
	unescaped, err := url.PathUnescape(urlPath)
	if err != nil {
		return nil, err
	}

	// Set the encoded path data
	u.RawPath = u.Path + urlPath
	u.Path = u.Path + unescaped

	// Set the query parameters
	if options != nil {
		q, err := query.Values(options)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// We are using JWT not Basic auth
	// req.SetBasicAuth(client.Config.ClientID, client.Config.ClientSecret)
	// req.Header.Add("Accept", "application/x-www-form-urlencoded")
	// if body != nil {
	// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// }

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.Config.AuthToken.AccessToken))
	req.Header.Add("Content-Type", "application/json")
	//req.Header.Add("Accept", "application/json")
	return req, nil
}

func (client *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := client.Config.GetHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}

	// response is close before parsing response body due to  this function defer the body close
	// defer func() {
	// 	if rerr := resp.Body.Close(); err == nil {
	// 		err = rerr
	// 	}
	// }()

	if err := checkErrorInResponse(resp); err != nil {
		return resp, err
	}

	if v != nil {
		if err := decodeJSON(resp, v); err != nil {
			return resp, err
		}
	}
	client.logResponse(resp)

	return resp, nil
}

func decodeJSON(res *http.Response, v interface{}) error {
	return json.NewDecoder(res.Body).Decode(&v)
}
