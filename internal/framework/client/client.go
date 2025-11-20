package client

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa"

	"github.com/zscaler/terraform-provider-zpa/v4/version"
)

// Config contains our provider configuration values and Zscaler clients.
type Config struct {
	ClientID         string
	ClientSecret     string
	CustomerID       string
	MicrotenantID    string
	VanityDomain     string
	Cloud            string
	PrivateKey       string
	HTTPProxy        string
	RetryCount       int
	Parallelism      int
	Backoff          bool
	MinWait          int
	MaxWait          int
	RequestTimeout   int
	UseLegacyClient  bool
	TerraformVersion string
	ProviderVersion  string

	// Legacy SDK specific fields
	ZPAClientID     string
	ZPAClientSecret string
	ZPACustomerID   string
	ZPACloud        string
}

// Client wraps the Zscaler SDK client
type Client struct {
	Service *zscaler.Service
}

// NewClient creates a new ZPA client based on the configuration
func NewClient(config *Config) (*Client, error) {
	if config.UseLegacyClient {
		return newLegacyClient(config)
	}
	return newV3Client(config)
}

// newLegacyClient creates a legacy V2 client
func newLegacyClient(config *Config) (*Client, error) {
	applyDefaults(config)

	customUserAgent := generateUserAgent(config.TerraformVersion, config.CustomerID)

	setters := []zpa.ConfigSetter{
		zpa.WithRateLimitMaxRetries(int32(config.RetryCount)),
		zpa.WithRateLimitMinWait(time.Duration(config.MinWait) * time.Second),
		zpa.WithRateLimitMaxWait(time.Duration(config.MaxWait) * time.Second),
		zpa.WithRequestTimeout(time.Duration(config.RequestTimeout) * time.Second),
	}

	// Disable cache when running TF acceptance tests
	if os.Getenv("TF_ACC") != "1" {
		setters = append(setters, zpa.WithCache(true))
	}

	setters = append(setters,
		zpa.WithZPAClientID(config.ZPAClientID),
		zpa.WithZPAClientSecret(config.ZPAClientSecret),
		zpa.WithZPACustomerID(config.ZPACustomerID),
		zpa.WithZPACloud(config.ZPACloud),
	)

	if config.MicrotenantID != "" {
		setters = append(setters, zpa.WithZPAMicrotenantID(config.MicrotenantID))
	}

	if config.HTTPProxy != "" {
		_url, err := url.Parse(config.HTTPProxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %v", err)
		}
		setters = append(setters, zpa.WithProxyHost(_url.Hostname()))

		sPort := _url.Port()
		if sPort == "" {
			sPort = "80"
		}
		port64, err := strconv.ParseInt(sPort, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy port: %v", err)
		}
		if port64 < 1 || port64 > 65535 {
			return nil, fmt.Errorf("invalid port number: must be between 1 and 65535, got: %d", port64)
		}
		port32 := int32(port64)
		setters = append(setters, zpa.WithProxyPort(port32))
	}

	zpaCfg, err := zpa.NewConfiguration(setters...)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZPA configuration: %v", err)
	}
	zpaCfg.UserAgent = customUserAgent

	legacyClient, err := zscaler.NewLegacyZpaClient(zpaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZPA client: %v", err)
	}

	return &Client{
		Service: zscaler.NewService(legacyClient.Client, nil),
	}, nil
}

// newV3Client creates a V3 client
func newV3Client(config *Config) (*Client, error) {
	applyDefaults(config)

	customUserAgent := generateUserAgent(config.TerraformVersion, config.CustomerID)

	setters := []zscaler.ConfigSetter{
		zscaler.WithRateLimitMaxRetries(int32(config.RetryCount)),
		zscaler.WithRateLimitMinWait(time.Duration(config.MinWait) * time.Second),
		zscaler.WithRateLimitMaxWait(time.Duration(config.MaxWait) * time.Second),
		zscaler.WithRequestTimeout(time.Duration(config.RequestTimeout) * time.Second),
		zscaler.WithUserAgentExtra(""),
	}

	// Disable cache when running TF acceptance tests
	if os.Getenv("TF_ACC") != "1" {
		setters = append(setters,
			zscaler.WithCache(true),
			zscaler.WithCacheTtl(10*time.Minute),
			zscaler.WithCacheTti(8*time.Minute),
		)
	}

	tfLog := os.Getenv("TF_LOG")
	if tfLog == "DEBUG" || tfLog == "TRACE" {
		setters = append(setters, zscaler.WithDebug(false))
		log.Println("[DEBUG] SDK debug logging enabled")
	}

	if config.HTTPProxy != "" {
		_url, err := url.Parse(config.HTTPProxy)
		if err != nil {
			return nil, err
		}
		setters = append(setters, zscaler.WithProxyHost(_url.Hostname()))

		sPort := _url.Port()
		if sPort == "" {
			sPort = "80"
		}
		port64, err := strconv.ParseInt(sPort, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy port: %v", err)
		}
		if port64 < 1 || port64 > 65535 {
			return nil, fmt.Errorf("invalid port number: must be between 1 and 65535, got: %d", port64)
		}
		port32 := int32(port64)
		setters = append(setters, zscaler.WithProxyPort(port32))
	}

	switch {
	case config.ClientID != "" && config.ClientSecret != "" && config.VanityDomain != "" && config.CustomerID != "":
		setters = append(setters,
			zscaler.WithClientID(config.ClientID),
			zscaler.WithClientSecret(config.ClientSecret),
			zscaler.WithVanityDomain(config.VanityDomain),
			zscaler.WithZPACustomerID(config.CustomerID),
		)

		if config.MicrotenantID != "" {
			setters = append(setters, zscaler.WithZPAMicrotenantID(config.MicrotenantID))
		}
		if config.Cloud != "" {
			setters = append(setters, zscaler.WithZscalerCloud(config.Cloud))
		}

	case config.ClientID != "" && config.PrivateKey != "" && config.VanityDomain != "" && config.CustomerID != "":
		setters = append(setters,
			zscaler.WithClientID(config.ClientID),
			zscaler.WithPrivateKey(config.PrivateKey),
			zscaler.WithVanityDomain(config.VanityDomain),
			zscaler.WithZPACustomerID(config.CustomerID),
		)

		if config.MicrotenantID != "" {
			setters = append(setters, zscaler.WithZPAMicrotenantID(config.MicrotenantID))
		}
		if config.Cloud != "" {
			setters = append(setters, zscaler.WithZscalerCloud(config.Cloud))
		}

	default:
		return nil, fmt.Errorf("invalid authentication configuration: missing required parameters")
	}

	configSet, err := zscaler.NewConfiguration(setters...)
	if err != nil {
		return nil, fmt.Errorf("failed to create SDK V3 configuration: %v", err)
	}

	configSet.UserAgent = customUserAgent

	v3Client, err := zscaler.NewOneAPIClient(configSet)
	if err != nil {
		return nil, fmt.Errorf("failed to create Zscaler API client: %v", err)
	}

	return &Client{
		Service: zscaler.NewService(v3Client.Client, nil),
	}, nil
}

func applyDefaults(config *Config) {
	if config.RetryCount == 0 {
		config.RetryCount = 100
	}
	if config.MinWait == 0 {
		config.MinWait = 2
	}
	if config.MaxWait == 0 {
		config.MaxWait = 10
	}
	if config.RequestTimeout == 0 {
		config.RequestTimeout = 240
	}
}

func generateUserAgent(terraformVersion, customerID string) string {
	providerVersion := version.ProviderVersion
	if providerVersion == "" {
		providerVersion = "dev"
	}
	if terraformVersion == "" {
		terraformVersion = "unknown"
	}
	if customerID == "" {
		customerID = "unknown"
	}
	return fmt.Sprintf("(%s %s) Terraform/%s Provider/%s Customer/%s",
		runtime.GOOS,
		runtime.GOARCH,
		terraformVersion,
		providerVersion,
		customerID,
	)
}
