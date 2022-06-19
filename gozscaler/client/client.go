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
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/zscaler/terraform-provider-zpa/gozscaler"
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

func (client *Client) NewRequestDo(method, url string, options, body, v interface{}) (*http.Response, error) {
	return client.newRequestDoCustom(method, url, options, body, v)
}

func (client *Client) newRequestDoCustom(method, urlStr string, options, body, v interface{}) (*http.Response, error) {
	client.Config.Lock()
	defer client.Config.Unlock()
	if client.Config.AuthToken == nil || client.Config.AuthToken.AccessToken == "" {
		if client.Config.ClientID == "" || client.Config.ClientSecret == "" {
			log.Printf("[ERROR] No client credentials were provided. Please set %s, %s and %s enviroment variables.\n", gozscaler.ZPA_CLIENT_ID, gozscaler.ZPA_CLIENT_SECRET, gozscaler.ZPA_CUSTOMER_ID)
			return nil, errors.New("no client credentials were provided")
		}
		log.Printf("[TRACE] Getting access token for %s=%s\n", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID)
		data := url.Values{}
		data.Set("client_id", client.Config.ClientID)
		data.Set("client_secret", client.Config.ClientSecret)
		req, err := http.NewRequest("POST", client.Config.BaseURL.String()+"/signin", strings.NewReader(data.Encode()))
		if err != nil {
			log.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID, err)
			return nil, fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID, err)

		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Config.GetHTTPClient().Do(req)

		if err != nil {
			log.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID, err)
			return nil, fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID, err)
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID, err)
			return nil, fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID, err)
		}
		if resp.StatusCode >= 300 {
			log.Printf("[ERROR] Failed to signin the user %s=%s, got http status:%dn response body:%s\n", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID, resp.StatusCode, respBody)
			return nil, fmt.Errorf("[ERROR] Failed to signin the user %s=%s, got http status:%d, response body:%s", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID, resp.StatusCode, respBody)
		}
		var a gozscaler.AuthToken
		err = json.Unmarshal(respBody, &a)
		if err != nil {
			log.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID, err)
			return nil, fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", gozscaler.ZPA_CLIENT_ID, client.Config.ClientID, err)

		}
		// we need keep auth token for future http request
		client.Config.AuthToken = &a
	}
	req, err := client.newRequest(method, urlStr, options, body)
	if err != nil {
		return nil, err
	}
	client.logRequest(req)
	return client.do(req, v)
}

// Generating the Http request
func (client *Client) newRequest(method, urlPath string, options, body interface{}) (*http.Request, error) {
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
