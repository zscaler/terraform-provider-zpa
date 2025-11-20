// Copyright (c) SecurityGeekIO, Inc.
// SPDX-License-Identifier: MPL-2.0

package acctest

import (
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
)

// testAccProtoV5ProviderFactories is the provider factory for ProtoV5 testing
var testAccProtoV5ProviderFactories map[string]func() (tfprotov5.ProviderServer, error)

func init() {
	testAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
		"zpa": func() (tfprotov5.ProviderServer, error) {
			provider := framework.New("test")
			return providerserver.NewProtocol5(provider)(), nil
		},
	}
}

// testAccPreCheck verifies required environment variables are set
func testAccPreCheck(t *testing.T) {
	// Check for mandatory environment variables for client_id + client_secret authentication
	if v := os.Getenv("ZSCALER_CLIENT_ID"); v == "" {
		t.Fatal("ZSCALER_CLIENT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("ZSCALER_CLIENT_SECRET"); v == "" {
		t.Fatal("ZSCALER_CLIENT_SECRET must be set for acceptance tests")
	}
	if v := os.Getenv("ZSCALER_VANITY_DOMAIN"); v == "" {
		t.Fatal("ZSCALER_VANITY_DOMAIN must be set for acceptance tests")
	}
	if v := os.Getenv("ZPA_CUSTOMER_ID"); v == "" {
		t.Fatal("ZPA_CUSTOMER_ID must be set for acceptance tests")
	}

	// Optional cloud configuration
	if v := os.Getenv("ZSCALER_CLOUD"); v == "" {
		t.Log("[INFO] ZSCALER_CLOUD is not set. Defaulting to production cloud.")
	}
}

// PreCheck is a public wrapper for testAccPreCheck
func PreCheck(t *testing.T) {
	testAccPreCheck(t)
}

// ProtoV5ProviderFactories returns the provider factories for ProtoV5 testing
func ProtoV5ProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return testAccProtoV5ProviderFactories
}

var (
	testClientOnce sync.Once
	testClient     *client.Client
	testClientErr  error
)

// TestClient returns a shared ZPA client configured from the acceptance test environment.
func TestClient(t *testing.T) *client.Client {
	t.Helper()

	testClientOnce.Do(func() {
		cfg := &client.Config{
			ClientID:         os.Getenv("ZSCALER_CLIENT_ID"),
			ClientSecret:     os.Getenv("ZSCALER_CLIENT_SECRET"),
			PrivateKey:       os.Getenv("ZSCALER_PRIVATE_KEY"),
			VanityDomain:     os.Getenv("ZSCALER_VANITY_DOMAIN"),
			Cloud:            os.Getenv("ZSCALER_CLOUD"),
			CustomerID:       os.Getenv("ZPA_CUSTOMER_ID"),
			MicrotenantID:    os.Getenv("ZPA_MICROTENANT_ID"),
			ZPAClientID:      os.Getenv("ZPA_CLIENT_ID"),
			ZPAClientSecret:  os.Getenv("ZPA_CLIENT_SECRET"),
			ZPACustomerID:    os.Getenv("ZPA_CUSTOMER_ID"),
			ZPACloud:         os.Getenv("ZPA_CLOUD"),
			UseLegacyClient:  strings.EqualFold(os.Getenv("ZSCALER_USE_LEGACY_CLIENT"), "true"),
			HTTPProxy:        os.Getenv("ZSCALER_HTTP_PROXY"),
			RetryCount:       100,
			Parallelism:      1,
			RequestTimeout:   240,
			MinWait:          2,
			MaxWait:          10,
			TerraformVersion: "",
			ProviderVersion:  "test",
		}

		testClient, testClientErr = client.NewClient(cfg)
	})

	if testClientErr != nil {
		t.Fatalf("failed to create acceptance test client: %v", testClientErr)
	}

	return testClient
}
