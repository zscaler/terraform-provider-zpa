package zpa

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
	// Define environment variables related to ZPA authentication
	envKeys := []string{
		"ZSCALER_CLIENT_ID",
		"ZSCALER_CLIENT_SECRET",
		"ZPA_CUSTOMER_ID",
		"ZSCALER_VANITY_DOMAIN",
		"ZSCALER_CLOUD",
	}
	envVals := make(map[string]string)
	// Save and clear ZPA env vars to test configuration cleanly
	for _, key := range envKeys {
		val := os.Getenv(key)
		if val != "" {
			envVals[key] = val
			os.Unsetenv(key)
		}
	}

	// Define test cases for the various configurations
	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		vanityDomain string
		customerID   string
		cloud        string // Optional field
		expectError  bool
	}{
		{"valid client_id + client_secret", "clientID", "clientSecret", "vanityDomain", "customerID", "cloud", false},
		{"missing client_id", "", "clientSecret", "vanityDomain", "customerID", "cloud", true},
		{"missing clientSecret", "clientID", "", "vanityDomain", "customerID", "cloud", true},
		{"missing vanity domain", "clientID", "clientSecret", "", "customerID", "cloud", true},
		{"valid client_id + client_secret without cloud", "clientID", "clientSecret", "vanityDomain", "customerID", "", false}, // Ensures cloud is optional
	}

	// Execute each test case
	for _, test := range tests {
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
			resourceConfig["cloud"] = test.cloud
		}

		config := terraform.NewResourceConfigRaw(resourceConfig)
		provider := ZPAProvider()
		err := provider.Validate(config)

		// Check expectations based on each test case setup
		if test.expectError && err == nil {
			t.Errorf("test %q: expected error but received none", test.name)
		}
		if !test.expectError && err != nil {
			t.Errorf("test %q: did not expect error but received error: %+v", test.name, err)
		}
	}

	// Restore environment variables after the tests
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
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SDK V3 client: %w", err)
	}

	testSdkV3Client = client
	return testSdkV3Client, nil
}
