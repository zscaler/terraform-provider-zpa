package gozscaler

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

const (
	defaultBaseURL           = "https://config.private.zscaler.com"
	defaultPrivateAPIBaseURL = "https://api.private.zscaler.com"
	defaultTimeout           = 240 * time.Second
	loggerPrefix             = "zpa-logger: "
	ZPA_CLIENT_ID            = "ZPA_CLIENT_ID"
	ZPA_CLIENT_SECRET        = "ZPA_CLIENT_SECRET"
	ZPA_CUSTOMER_ID          = "ZPA_CUSTOMER_ID"
)

// BackoffConfig contains all the configuration for the backoff and retry mechanism
type BackoffConfig struct {
	Enabled             bool // Set to true to enable backoff and retry mechanism
	RetryWaitMinSeconds int  // Minimum time to wait
	RetryWaitMaxSeconds int  // Maximum time to wait
	MaxNumOfRetries     int  // Maximum number of retries
}

// Need to implement exponential back off to comply with the API rate limit. https://help.zscaler.com/zpa/about-rate-limiting
// 20 times in a 10 second interval for a GET call.
// 10 times in a 10 second interval for any POST/PUT/DELETE call.
// See example: https://github.com/okta/terraform-provider-okta/blob/master/okta/config.go
type AuthToken struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	//ExpiresIn   string `json:"expires_in"`
}

// Config contains all the configuration data for the API client
type Config struct {
	BaseURL           *url.URL
	PrivateAPIBaseURL *url.URL
	httpClient        *http.Client
	// The logger writer interface to write logging messages to. Defaults to standard out.
	Logger *log.Logger
	// Credentials for basic authentication.
	ClientID, ClientSecret, CustomerID string
	// Backoff config
	BackoffConf *BackoffConfig
	AuthToken   *AuthToken
	sync.Mutex
}

/*
NewConfig returns a default configuration for the client.
By default it will try to read the access and te secret from the environment variable.
*/
// Need to implement exponential back off to comply with the API rate limit. https://help.zscaler.com/zpa/about-rate-limiting
// 20 times in a 10 second interval for a GET call.
// 10 times in a 10 second interval for any POST/PUT/DELETE call.
// See example: https://github.com/okta/terraform-provider-okta/blob/master/okta/config.go
// TODO Add healthCheck method to NewConfig
func NewConfig(clientID, clientSecret, customerID, rawUrl string) (*Config, error) {
	backoffConf := &BackoffConfig{
		Enabled:             true,
		MaxNumOfRetries:     100,
		RetryWaitMaxSeconds: 20,
		RetryWaitMinSeconds: 5,
	}
	if clientID == "" || clientSecret == "" || customerID == "" {
		clientID = os.Getenv(ZPA_CLIENT_ID)
		clientSecret = os.Getenv(ZPA_CLIENT_SECRET)
		customerID = os.Getenv(ZPA_CUSTOMER_ID)
	}
	if rawUrl == "" {
		rawUrl = defaultBaseURL
	}

	var logger *log.Logger
	if loggerEnv := os.Getenv("ZSCALER_SDK_LOG"); loggerEnv == "true" {
		logger = getDefaultLogger()
	}

	baseURL, err := url.Parse(rawUrl)
	if err != nil {
		log.Printf("[ERROR] error occured while configuring the client: %v", err)
	}
	privateAPIBaseURL, err := url.Parse(defaultPrivateAPIBaseURL)
	if err != nil {
		log.Printf("[ERROR] error occured while configuring the client: %v", err)
	}
	return &Config{
		BaseURL:           baseURL,
		PrivateAPIBaseURL: privateAPIBaseURL,
		Logger:            logger,
		httpClient:        nil,
		ClientID:          clientID,
		ClientSecret:      clientSecret,
		CustomerID:        customerID,
		BackoffConf:       backoffConf,
	}, err
}

func (c *Config) GetHTTPClient() *http.Client {
	if c.httpClient == nil {
		if c.BackoffConf != nil && c.BackoffConf.Enabled {
			retryableClient := retryablehttp.NewClient()
			retryableClient.RetryWaitMin = time.Second * time.Duration(c.BackoffConf.RetryWaitMinSeconds)
			retryableClient.RetryWaitMax = time.Second * time.Duration(c.BackoffConf.RetryWaitMaxSeconds)
			retryableClient.RetryMax = c.BackoffConf.MaxNumOfRetries
			retryableClient.HTTPClient.Transport = logging.NewTransport("gozscaler", retryableClient.HTTPClient.Transport)
			retryableClient.CheckRetry = checkRetry
			retryableClient.HTTPClient.Timeout = defaultTimeout
			c.httpClient = retryableClient.StandardClient()
		} else {
			c.httpClient = &http.Client{
				Timeout: defaultTimeout,
			}
		}
	}
	return c.httpClient
}

func getDefaultLogger() *log.Logger {
	return log.New(os.Stdout, loggerPrefix, log.LstdFlags|log.Lshortfile)
}

func containsInt(codes []int, code int) bool {
	for _, a := range codes {
		if a == code {
			return true
		}
	}
	return false
}

// getRetryOnStatusCodes return a list of http status codes we want to apply retry on.
// return empty slice to enable retry on all connection & server errors.
// or return []int{429}  to retry on only TooManyRequests error
func getRetryOnStatusCodes() []int {
	return []int{http.StatusTooManyRequests}
}

// Used to make http client retry on provided list of response status codes
func checkRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if resp != nil && containsInt(getRetryOnStatusCodes(), resp.StatusCode) {
		return true, nil
	}
	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}
