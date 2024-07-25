package zpa

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
)

var (
	testSdkClient            *Client
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
	// TF_VAR_hostname allows the real hostname to be scripted into the config tests
	os.Setenv("TF_VAR_hostname", fmt.Sprintf("%s.%s.%s.%s", os.Getenv("ZPA_CLIENT_ID"), os.Getenv("ZPA_CLIENT_SECRET"), os.Getenv("ZPA_CUSTOMER_ID"), os.Getenv("ZPA_CLOUD")))

	// NOTE: Acceptance test sweepers are necessary to prevent dangling
	// resources.
	// NOTE: Don't run sweepers if we are playing back VCR as nothing should be
	// going over the wire
	if os.Getenv("ZPA_VCR_TF_ACC") != "play" {
		// TODO: Tests is failing on QA2 tenant. Needs further investigation.
		// setupSweeper(resourcetype.ZPAAppConnectorGroup, sweepTestAppConnectorGroup)
		setupSweeper(resourcetype.ZPAApplicationServer, sweepTestApplicationServer)
		setupSweeper(resourcetype.ZPAApplicationSegment, sweepTestApplicationSegment)
		setupSweeper(resourcetype.ZPAApplicationSegmentBrowserAccess, sweepTestApplicationSegmentBA)
		setupSweeper(resourcetype.ZPAApplicationSegmentInspection, sweepTestApplicationInspection)
		setupSweeper(resourcetype.ZPAApplicationSegmentPRA, sweepTestApplicationPRA)
		setupSweeper(resourcetype.ZPABACertificate, sweepTestBaCertificate)
		setupSweeper(resourcetype.ZPAInspectionCustomControl, sweepTestInspectionCustomControl)
		setupSweeper(resourcetype.ZPAInspectionProfile, sweepTestInspectionProfile)
		setupSweeper(resourcetype.ZPALSSController, sweepTestLSSConfigController)
		// setupSweeper(resourcetype.ZPASegmentGroup, sweepTestSegmentGroup)
		setupSweeper(resourcetype.ZPAServerGroup, sweepTestServerGroup)
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

func testAccPreCheck(t *testing.T) func() {
	return func() {
		err := accPreCheck()
		if err != nil {
			t.Fatalf("%v", err)
		}
	}
}

func accPreCheck() error {
	ClientID := os.Getenv("ZPA_CLIENT_ID")
	ClientSecret := os.Getenv("ZPA_CLIENT_SECRET")
	CustomerID := os.Getenv("ZPA_CUSTOMER_ID")
	BaseURL := os.Getenv("ZPA_CLOUD")

	// Check for the presence of necessary environment variables.
	if ClientID == "" {
		return errors.New("ZPA_CLIENT_ID must be set for acceptance tests")
	}

	if ClientSecret == "" {
		return errors.New("ZPA_CLIENT_SECRET must be set for acceptance tests")
	}

	if CustomerID == "" {
		return errors.New("ZPA_CUSTOMER_ID must be set for acceptance tests")
	}

	if BaseURL == "" {
		return errors.New("ZPA_CLOUD must be set for acceptance tests")
	}

	return nil
}

func sdkClientForTest() (*Client, error) {
	if testSdkClient == nil {
		sweeperLogger.Warn("testSdkClient is not initialized. Initializing now...")

		config := &Config{
			ClientID:     os.Getenv("ZPA_CLIENT_ID"),
			ClientSecret: os.Getenv("ZPA_CLIENT_SECRET"),
			CustomerID:   os.Getenv("ZPA_CUSTOMER_ID"),
			BaseURL:      os.Getenv("ZPA_CLOUD"),
			UserAgent:    "terraform-provider-zpa",
		}

		var err error
		testSdkClient, err = config.Client()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize testSdkClient: %v", err)
		}
	}
	return testSdkClient, nil
}
