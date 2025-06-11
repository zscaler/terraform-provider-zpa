package zpa

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

var (
	testSdkV3Client          *zscaler.Client
	testAccProvider          *schema.Provider
	testAccProviders         map[string]*schema.Provider
	testAccProviderFactories map[string]func() (*schema.Provider, error)
)

func init() {
	testAccProvider = ZPAProvider()
	testAccProviders = map[string]*schema.Provider{
		"zpa": testAccProvider,
	}

	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"zpa": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

// TestMain overridden main testing function. Package level BeforeAll and AfterAll.
// It also delineates between acceptance tests and unit tests
func TestMain(m *testing.M) {
	os.Setenv("TF_VAR_hostname", fmt.Sprintf("%s.%s.%s.%s", os.Getenv("ZSCALER_CLIENT_ID"), os.Getenv("ZSCALER_CLIENT_SECRET"), os.Getenv("ZPA_CUSTOMER_ID"), os.Getenv("ZSCALER_CLOUD")))

	if os.Getenv("ZPA_VCR_TF_ACC") != "play" {
		setupSweeper(resourcetype.ZPAAppConnectorGroup, sweepTestAppConnectorGroup)
		setupSweeper(resourcetype.ZPAApplicationServer, sweepTestApplicationServer)
		setupSweeper(resourcetype.ZPAApplicationSegment, sweepTestApplicationSegment)
		setupSweeper(resourcetype.ZPAApplicationSegmentBrowserAccess, sweepTestApplicationSegmentBA)
		setupSweeper(resourcetype.ZPAApplicationSegmentInspection, sweepTestApplicationInspection)
		setupSweeper(resourcetype.ZPAApplicationSegmentPRA, sweepTestApplicationPRA)
		setupSweeper(resourcetype.ZPABACertificate, sweepTestBaCertificate)
		setupSweeper(resourcetype.ZPAInspectionCustomControl, sweepTestInspectionCustomControl)
		setupSweeper(resourcetype.ZPAInspectionProfile, sweepTestInspectionProfile)
		setupSweeper(resourcetype.ZPALSSController, sweepTestLSSConfigController)
		setupSweeper(resourcetype.ZPAServerGroup, sweepTestServerGroup)
		setupSweeper(resourcetype.ZPASegmentGroup, sweepTestSegmentGroup)
		setupSweeper(resourcetype.ZPAServiceEdgeGroup, sweepTestServiceEdgeGroup)
		setupSweeper(resourcetype.ZPAPolicyAccessRule, sweepTestAccessPolicyRuleByType)
		setupSweeper(resourcetype.ZPACBIBannerController, sweepTestCBIBanner)
		setupSweeper(resourcetype.ZPACBIExternalIsolationProfile, sweepTestCBIExternalProfile)
		setupSweeper(resourcetype.ZPACBICertificate, sweepTestCBICertificate)
		setupSweeper(resourcetype.ZPAPRAConsoleController, sweepTestPRAConsoleController)
		setupSweeper(resourcetype.ZPAPRACredentialController, sweepTestPRACredentialController)
		setupSweeper(resourcetype.ZPAPRAPortalController, sweepTestPRAPortalController)
		setupSweeper(resourcetype.ZPAPRAApprovalController, sweepTestPRAPrivilegedApprovalController)
	}

	resource.TestMain(m)
}

func TestProvider(t *testing.T) {
	if err := ZPAProvider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	_ = ZPAProvider()
}

func testAccPreCheck(t *testing.T) func() {
	return func() {
		err := accPreCheck()
		if err != nil {
			t.Fatalf("%v", err)
		}
	}
}

func accPreCheck() error {
	// Check for mandatory environment variables for client_id + client_secret authentication
	if v := os.Getenv("ZSCALER_CLIENT_ID"); v == "" {
		return errors.New("ZSCALER_CLIENT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("ZSCALER_CLIENT_SECRET"); v == "" {
		return errors.New("ZSCALER_CLIENT_SECRET must be set for acceptance tests")
	}
	if v := os.Getenv("ZSCALER_VANITY_DOMAIN"); v == "" {
		return errors.New("ZSCALER_VANITY_DOMAIN must be set for acceptance tests")
	}
	if v := os.Getenv("ZPA_CUSTOMER_ID"); v == "" {
		return errors.New("ZPA_CUSTOMER_ID must be set for acceptance tests")
	}

	// Optional cloud configuration
	if v := os.Getenv("ZSCALER_CLOUD"); v == "" {
		log.Println("[INFO] ZSCALER_CLOUD is not set. Defaulting to production cloud.")
	}

	return nil
}

func TestProviderValidate(t *testing.T) {
	// Save current env vars and clear them to simulate controlled test inputs
	envKeys := []string{
		"ZSCALER_CLIENT_ID",
		"ZSCALER_CLIENT_SECRET",
		"ZPA_CUSTOMER_ID",
		"ZSCALER_VANITY_DOMAIN",
		"ZSCALER_CLOUD",
	}
	envVals := make(map[string]string)
	for _, key := range envKeys {
		envVals[key] = os.Getenv(key)
		os.Unsetenv(key)
	}

	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		vanityDomain string
		customerID   string
		cloud        string
		expectError  bool
		skipIfCloud  bool
	}{
		{
			name:         "valid client_id + client_secret",
			clientID:     envVals["ZSCALER_CLIENT_ID"],
			clientSecret: envVals["ZSCALER_CLIENT_SECRET"],
			vanityDomain: envVals["ZSCALER_VANITY_DOMAIN"],
			customerID:   envVals["ZPA_CUSTOMER_ID"],
			cloud:        envVals["ZSCALER_CLOUD"],
			expectError:  false,
		},
		{
			name:         "missing client_id",
			clientID:     "",
			clientSecret: "clientSecret",
			vanityDomain: "vanityDomain",
			customerID:   "customerID",
			cloud:        "zscaler_cloud",
			expectError:  true,
		},
		{
			name:         "missing clientSecret",
			clientID:     "clientID",
			clientSecret: "",
			vanityDomain: "vanityDomain",
			customerID:   "customerID",
			cloud:        "zscaler_cloud",
			expectError:  true,
		},
		{
			name:         "missing vanity domain",
			clientID:     "clientID",
			clientSecret: "clientSecret",
			vanityDomain: "",
			customerID:   "customerID",
			cloud:        "zscaler_cloud",
			expectError:  true,
		},
		{
			name:         "valid client_id + client_secret without zscaler_cloud",
			clientID:     envVals["ZSCALER_CLIENT_ID"],
			clientSecret: envVals["ZSCALER_CLIENT_SECRET"],
			vanityDomain: envVals["ZSCALER_VANITY_DOMAIN"],
			customerID:   envVals["ZPA_CUSTOMER_ID"],
			cloud:        "", // simulate unset cloud
			expectError:  false,
			skipIfCloud:  true, // only run if cloud is NOT needed
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Skip if this test depends on cloud being unset and cloud is required
			if test.skipIfCloud && envVals["ZSCALER_CLOUD"] != "" {
				t.Skipf("Skipping '%s': cloud-specific credentials cannot authenticate without ZSCALER_CLOUD", test.name)
			}

			resourceConfig := map[string]interface{}{
				"vanity_domain": test.vanityDomain,
				"customer_id":   test.customerID,
			}
			if test.clientID != "" {
				resourceConfig["client_id"] = test.clientID
			}
			if test.clientSecret != "" {
				resourceConfig["client_secret"] = test.clientSecret
			}
			if test.cloud != "" {
				resourceConfig["zscaler_cloud"] = test.cloud
			}

			provider := ZPAProvider()
			rawData := schema.TestResourceDataRaw(t, provider.Schema, resourceConfig)
			_, diags := provider.ConfigureContextFunc(context.Background(), rawData)

			if test.expectError && !diags.HasError() {
				t.Errorf("expected error but received none")
			}
			if !test.expectError && diags.HasError() {
				t.Errorf("did not expect error but received: %+v", diags)
			}
		})
	}

	// Restore environment variables after test run
	for key, val := range envVals {
		os.Setenv(key, val)
	}
}

func sdkV3ClientForTest() (*zscaler.Client, error) {
	if testSdkV3Client != nil {
		return testSdkV3Client, nil
	}

	// Initialize the SDK V3 Client
	client, err := zscalerSDKV3Client(&Config{
		clientID:     os.Getenv("ZSCALER_CLIENT_ID"),
		clientSecret: os.Getenv("ZSCALER_CLIENT_SECRET"),
		customerID:   os.Getenv("ZPA_CUSTOMER_ID"),
		vanityDomain: os.Getenv("ZSCALER_VANITY_DOMAIN"),
		cloud:        os.Getenv("ZSCALER_CLOUD"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SDK V3 client: %w", err)
	}

	testSdkV3Client = client
	return testSdkV3Client, nil
}
