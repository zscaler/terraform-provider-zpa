package gozscaler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

const (
	defaultBaseURL           = "https://config.private.zscaler.com"
	betaBaseURL              = "https://config.zpabeta.net"
	govBaseURL               = "https://config.zpagov.net"
	defaultTimeout           = 240 * time.Second
	loggerPrefix             = "zpa-logger: "
	ZPA_CLIENT_ID            = "ZPA_CLIENT_ID"
	ZPA_CLIENT_SECRET        = "ZPA_CLIENT_SECRET"
	ZPA_CUSTOMER_ID          = "ZPA_CUSTOMER_ID"
	ZPA_CLOUD                = "ZPA_CLOUD"
	configPath        string = ".zpa/credentials.json"
)

type BackoffConfig struct {
	Enabled             bool // Set to true to enable backoff and retry mechanism
	RetryWaitMinSeconds int  // Minimum time to wait
	RetryWaitMaxSeconds int  // Maximum time to wait
	MaxNumOfRetries     int  // Maximum number of retries
}
type AuthToken struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}
type CredentialsConfig struct {
	ClientID     string `json:"zpa_client_id"`
	ClientSecret string `json:"zpa_client_secret"`
	CustomerID   string `json:"zpa_customer_id"`
	ZpaCloud     string `json:"zpa_cloud"`
}

// Config contains all the configuration data for the API client
type Config struct {
	BaseURL    *url.URL
	httpClient *http.Client
	// The logger writer interface to write logging messages to. Defaults to standard out.
	Logger *log.Logger
	// Credentials for basic authentication.
	ClientID, ClientSecret, CustomerID string
	// Backoff config
	BackoffConf *BackoffConfig
	AuthToken   *AuthToken
	sync.Mutex
	UserAgent string
}

/*
NewConfig returns a default configuration for the client.
By default it will try to read the access and te secret from the environment variable.
*/
// Need to implement exponential back off to comply with the API rate limit. https://help.zscaler.com/zpa/about-rate-limiting
// 20 times in a 10 second interval for a GET call.
// 10 times in a 10 second interval for any POST/PUT/DELETE call.
// TODO Add healthCheck method to NewConfig
func NewConfig(clientID, clientSecret, customerID, cloud, userAgent string) (*Config, error) {
	backoffConf := &BackoffConfig{
		Enabled:             true,
		MaxNumOfRetries:     100,
		RetryWaitMaxSeconds: 20,
		RetryWaitMinSeconds: 5,
	}
	// if creds not provided in TF config, try loading from env vars
	if clientID == "" || clientSecret == "" || customerID == "" || cloud == "" || userAgent == "" {
		clientID = os.Getenv(ZPA_CLIENT_ID)
		clientSecret = os.Getenv(ZPA_CLIENT_SECRET)
		customerID = os.Getenv(ZPA_CUSTOMER_ID)
		cloud = os.Getenv(ZPA_CLOUD)
	}
	// last resort to configuration file:
	if clientID == "" || clientSecret == "" || customerID == "" {
		creds, err := loadCredentialsFromConfig()
		if err != nil || creds == nil {
			return nil, err
		}
		clientID = creds.ClientID
		clientSecret = creds.ClientSecret
		customerID = creds.CustomerID
		cloud = creds.ZpaCloud
	}
	rawUrl := defaultBaseURL
	if cloud == "" {
		cloud = os.Getenv(ZPA_CLOUD)
	}
	if strings.EqualFold(cloud, "BETA") {
		rawUrl = betaBaseURL
	}
	if strings.EqualFold(cloud, "GOV") {
		rawUrl = govBaseURL
	}

	var logger *log.Logger
	if loggerEnv := os.Getenv("ZSCALER_SDK_LOG"); loggerEnv == "true" {
		logger = getDefaultLogger()
	}

	baseURL, err := url.Parse(rawUrl)
	if err != nil {
		log.Printf("[ERROR] error occured while configuring the client: %v", err)
	}
	return &Config{
		BaseURL:      baseURL,
		Logger:       logger,
		httpClient:   nil,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		CustomerID:   customerID,
		BackoffConf:  backoffConf,
	}, err
}

// loadCredentialsFromConfig Returns the credentials found in a config file
func loadCredentialsFromConfig() (*CredentialsConfig, error) {
	usr, _ := user.Current()
	dir := usr.HomeDir
	path := filepath.Join(dir, configPath)
	log.Printf("[INFO]Loading configuration file at:%s", path)
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("Could not open credentials file, needs to contain one json object with keys: zpa_client_id, zpa_client_secret, zpa_customer_id, and zpa_cloud. " + err.Error())
	}
	configBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var config CredentialsConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil || config.ClientID == "" || config.ClientSecret == "" || config.CustomerID == "" || config.ZpaCloud == "" {
		return nil, fmt.Errorf("could not parse credentials file, needs to contain one json object with keys: zpa_client_id, zpa_client_secret, zpa_customer_id, and zpa_cloud. error: %v", err)
	}
	return &config, nil
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
