package zpa

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"context"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegmentinspection"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegmentpra"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appservercontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/bacertificate"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/browseraccess"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/cloudbrowserisolation/cbibannercontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/cloudbrowserisolation/cbicertificatecontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/cloudbrowserisolation/cbiprofilecontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/inspectioncontrol/inspection_custom_controls"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/inspectioncontrol/inspection_profile"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/lssconfigcontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/privilegedremoteaccess/praapproval"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/privilegedremoteaccess/praconsole"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/privilegedremoteaccess/pracredential"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/privilegedremoteaccess/praportal"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/provisioningkey"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/segmentgroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/servergroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/serviceedgegroup"
)

var (
	sweeperLogger   hclog.Logger
	sweeperLogLevel hclog.Level
)

func init() {
	sweeperLogLevel = hclog.Warn
	if os.Getenv("TF_LOG") != "" {
		sweeperLogLevel = hclog.LevelFromString(os.Getenv("TF_LOG"))
	}
	sweeperLogger = hclog.New(&hclog.LoggerOptions{
		Level:      sweeperLogLevel,
		TimeFormat: "2006/01/02 03:04:05",
	})
}

func logSweptResource(kind, id, nameOrLabel string) {
	sweeperLogger.Warn(fmt.Sprintf("sweeper found dangling %q %q %q", kind, id, nameOrLabel))
}

type testClient struct {
	sdkClient *Client
}

var (
	testResourcePrefix   = "tf-acc-test-"
	updateResourcePrefix = "tf-updated-"
)

func TestRunForcedSweeper(t *testing.T) {
	if os.Getenv("ZPA_VCR_TF_ACC") != "" {
		t.Skip("forced sweeper is live and will never be run within VCR")
		return
	}
	if os.Getenv("ZPA_ACC_TEST_FORCE_SWEEPERS") == "" || os.Getenv("TF_ACC") == "" {
		t.Skipf("ENV vars %q and %q must not be blank to force running of the sweepers", "ZPA_ACC_TEST_FORCE_SWEEPERS", "TF_ACC")
		return
	}

	provider := ZPAProvider()
	c := terraform.NewResourceConfigRaw(nil)
	diag := provider.Configure(context.TODO(), c)
	if diag.HasError() {
		t.Skipf("sweeper's provider configuration failed: %v", diag)
		return
	}

	sdkClient, err := sdkClientForTest()
	if err != nil {
		t.Fatalf("Failed to get SDK client: %s", err)
	}

	testClient := &testClient{
		sdkClient: sdkClient,
	}
	if *sweepFlag == "global" {
		sweepTestAppConnectorGroup(testClient)
		sweepTestApplicationServer(testClient)
		sweepTestApplicationSegment(testClient)
		sweepTestApplicationSegmentBA(testClient)
		sweepTestApplicationInspection(testClient)
		sweepTestApplicationPRA(testClient)
		sweepTestInspectionCustomControl(testClient)
		sweepTestInspectionProfile(testClient)
		sweepTestLSSConfigController(testClient) // TODO: Tests is failing on QA2 tenant. Needs further investigation.
		sweepTestAccessPolicyRuleByType(testClient)
		sweepTestProvisioningKey(testClient)
		sweepTestSegmentGroup(testClient)
		sweepTestServerGroup(testClient)
		sweepTestServiceEdgeGroup(testClient)
		sweepTestCBIBanner(testClient)
		sweepTestCBIExternalProfile(testClient)
		sweepTestPRACredentialController(testClient)
		sweepTestPRAConsoleController(testClient)
		sweepTestPRAPortalController(testClient)
		sweepTestPRAPrivilegedApprovalController(testClient)
	}
}

// Sets up sweeper to clean up dangling resources
func setupSweeper(resourceType string, del func(*testClient) error) {
	sdkClient, err := sdkClientForTest()
	if err != nil {
		// You might decide how to handle the error here. Using a panic for simplicity.
		panic(fmt.Sprintf("Failed to get SDK client: %s", err))
	}
	resource.AddTestSweepers(resourceType, &resource.Sweeper{
		Name: resourceType,
		F: func(_ string) error {
			return del(&testClient{sdkClient: sdkClient})
		},
	})
}

// TODO: Tests is failing on QA2 tenant. Needs further investigation.
func sweepTestAppConnectorGroup(client *testClient) error {
	var errorList []error
	group, _, err := appconnectorgroup.GetAll(client.sdkClient.AppConnectorGroup)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := appconnectorgroup.Delete(client.sdkClient.AppConnectorGroup, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAAppConnectorGroup, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestApplicationServer(client *testClient) error {
	var errorList []error
	server, _, err := appservercontroller.GetAll(client.sdkClient.AppServerController)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(server)))
	for _, b := range server {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := appservercontroller.Delete(client.sdkClient.AppServerController, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAApplicationServer, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestApplicationSegment(client *testClient) error {
	var errorList []error
	appSegment, _, err := applicationsegment.GetAll(client.sdkClient.ApplicationSegment)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(appSegment)))
	for _, b := range appSegment {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := applicationsegment.Delete(client.sdkClient.ApplicationSegment, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAApplicationSegment, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestApplicationSegmentBA(client *testClient) error {
	var errorList []error
	appSegmentBA, _, err := browseraccess.GetAll(client.sdkClient.BrowserAccess)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(appSegmentBA)))
	for _, b := range appSegmentBA {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := browseraccess.Delete(client.sdkClient.BrowserAccess, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAApplicationSegmentBrowserAccess, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestApplicationInspection(client *testClient) error {
	var errorList []error
	appInspection, _, err := applicationsegmentinspection.GetAll(client.sdkClient.ApplicationSegmentInspection)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(appInspection)))
	for _, b := range appInspection {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := applicationsegmentinspection.Delete(client.sdkClient.ApplicationSegmentInspection, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAApplicationSegmentInspection, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestApplicationPRA(client *testClient) error {
	var errorList []error
	pra, _, err := applicationsegmentpra.GetAll(client.sdkClient.ApplicationSegmentPRA)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(pra)))
	for _, b := range pra {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := applicationsegmentpra.Delete(client.sdkClient.ApplicationSegmentPRA, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAApplicationSegmentPRA, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestInspectionCustomControl(client *testClient) error {
	var errorList []error
	customControl, _, err := inspection_custom_controls.GetAll(client.sdkClient.InspectionCustomControls)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(customControl)))
	for _, b := range customControl {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := inspection_custom_controls.Delete(client.sdkClient.InspectionCustomControls, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAInspectionCustomControl, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestInspectionProfile(client *testClient) error {
	var errorList []error
	profile, _, err := inspection_profile.GetAll(client.sdkClient.InspectionProfile)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(profile)))
	for _, b := range profile {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := inspection_profile.Delete(client.sdkClient.InspectionProfile, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAInspectionProfile, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestLSSConfigController(client *testClient) error {
	var errorList []error

	lssConfig, _, err := lssconfigcontroller.GetAll(client.sdkClient.LSSConfigController)
	if err != nil {
		if strings.Contains(err.Error(), "resource.not.found") {
			// Log that the resource was not found and continue
			sweeperLogger.Info("No resources found to sweep.")
			return nil
		}
		// If any other error, return it
		return err
	}

	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(lssConfig)))

	for _, b := range lssConfig {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.LSSConfig.Name, testResourcePrefix) || strings.HasPrefix(b.LSSConfig.Name, updateResourcePrefix) {
			// Attempt to delete the resource
			_, err := lssconfigcontroller.Delete(client.sdkClient.LSSConfigController, b.ID)
			if err != nil {
				// Check if the error is because the resource doesn't exist
				if strings.Contains(err.Error(), "resource.not.found") {
					sweeperLogger.Info(fmt.Sprintf("Resource %s with ID %s was already deleted.", resourcetype.ZPALSSController, fmt.Sprintf(b.ID)))
					continue
				}
				// For any other error, append to the error list and continue
				errorList = append(errorList, err)
				continue
			}

			sweeperLogger.Info(fmt.Sprintf("Swept resource %s with ID %s named %s.", resourcetype.ZPALSSController, fmt.Sprintf(b.ID), b.LSSConfig.Name))
		}
	}

	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}

	return condenseError(errorList)
}

var defaultPolicyNames = map[string]string{
	"ACCESS_POLICY":                        "Global_Policy",
	"TIMEOUT_POLICY":                       "ReAuth_Policy",
	"CLIENT_FORWARDING_POLICY":             "Bypass_Policy",
	"INSPECTION_POLICY":                    "Inspection_Policy",
	"ISOLATION_POLICY":                     "Isolation_Policy",
	"SIEM_POLICY":                          "Siem_Policy",
	"CREDENTIAL_POLICY":                    "Credential_Policy",
	"CAPABILITIES_POLICY":                  "Capabilities_Policy",
	"CLIENTLESS_SESSION_PROTECTION_POLICY": "Clientless_Session_Protection_Policy",
	"REDIRECTION_POLICY":                   "ReDirection_Policy",
}

func sweepTestAccessPolicyRuleByType(client *testClient) error {
	var errorList []error

	policyTypes := []string{
		"ACCESS_POLICY",
		"TIMEOUT_POLICY",
		"CLIENT_FORWARDING_POLICY",
		"INSPECTION_POLICY",
		"ISOLATION_POLICY",
		"SIEM_POLICY",
		"CREDENTIAL_POLICY",
		"CAPABILITIES_POLICY",
		"CLIENTLESS_SESSION_PROTECTION_POLICY",
		"REDIRECTION_POLICY",
	}

	for _, policyType := range policyTypes {
		// Fetch the PolicySet details for the current policy type to get the PolicySetID
		policySet, _, err := policysetcontroller.GetByPolicyType(client.sdkClient.PolicySetController, policyType)
		if err != nil {
			// If we fail to get a PolicySetID for a specific policy type, append the error and continue to the next type
			errorList = append(errorList, fmt.Errorf("Failed to get PolicySetID for policy type %s: %v", policyType, err))
			continue
		}
		policySetID := policySet.ID

		// Fetch all rules for the current policy type
		rules, _, err := policysetcontroller.GetAllByType(client.sdkClient.PolicySetController, policyType)
		if err != nil {
			// If we fail to fetch rules for a specific policy type, append the error and continue to the next type
			errorList = append(errorList, fmt.Errorf("Failed to get rules for policy type %s: %v", policyType, err))
			continue
		}
		// Logging the number of identified resources before the deletion loop for the current policy type
		sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep for policy type %s", len(rules), policyType))

		for _, rule := range rules {
			// Check if the rule's name is a default name and skip it
			if rule.Name == defaultPolicyNames[policyType] {
				continue
			}

			// Check if the resource name has the required prefix before deleting it
			if strings.HasPrefix(rule.Name, testResourcePrefix) || strings.HasPrefix(rule.Name, updateResourcePrefix) {
				// Use the fetched PolicySetID for deletion
				if _, err := policysetcontroller.Delete(client.sdkClient.PolicySetController, policySetID, rule.ID); err != nil {
					errorList = append(errorList, err)
					continue
				}
				logSweptResource(resourcetype.ZPAPolicyAccessRule, rule.ID, rule.Name)
			}
		}
	}

	// Log errors encountered during the sweeping process
	for _, err := range errorList {
		sweeperLogger.Error(err.Error())
	}
	return condenseError(errorList)
}

func sweepTestProvisioningKey(client *testClient) error {
	var errorList []error
	provisioningKey, err := provisioningkey.GetAll(client.sdkClient.ProvisioningKey)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(provisioningKey)))
	for _, b := range provisioningKey {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			// Assuming 'AssociationType' is a field in the provisioningKey object
			if _, err := provisioningkey.Delete(client.sdkClient.ProvisioningKey, b.AssociationType, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAProvisioningKey, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestSegmentGroup(client *testClient) error {
	var errorList []error
	group, _, err := segmentgroup.GetAll(client.sdkClient.SegmentGroup)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := segmentgroup.Delete(client.sdkClient.SegmentGroup, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPASegmentGroup, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestServerGroup(client *testClient) error {
	var errorList []error
	group, _, err := servergroup.GetAll(client.sdkClient.ServerGroup)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := servergroup.Delete(client.sdkClient.ServerGroup, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAServerGroup, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestServiceEdgeGroup(client *testClient) error {
	var errorList []error
	group, _, err := serviceedgegroup.GetAll(client.sdkClient.ServiceEdgeGroup)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := serviceedgegroup.Delete(client.sdkClient.ServiceEdgeGroup, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAServiceEdgeGroup, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestCBIBanner(client *testClient) error {
	var errorList []error
	group, _, err := cbibannercontroller.GetAll(client.sdkClient.CBIBannerController)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := cbibannercontroller.Delete(client.sdkClient.CBIBannerController, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPACBIBannerController, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestCBIExternalProfile(client *testClient) error {
	var errorList []error
	group, _, err := cbiprofilecontroller.GetAll(client.sdkClient.CBIProfileController)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := cbiprofilecontroller.Delete(client.sdkClient.CBIProfileController, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPACBIExternalIsolationProfile, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestCBICertificate(client *testClient) error {
	var errorList []error
	group, _, err := cbicertificatecontroller.GetAll(client.sdkClient.CBICertificateController)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := cbicertificatecontroller.Delete(client.sdkClient.CBICertificateController, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPACBICertificate, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestBaCertificate(client *testClient) error {
	var errorList []error
	group, _, err := bacertificate.GetAll(client.sdkClient.BACertificate)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := bacertificate.Delete(client.sdkClient.BACertificate, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPABACertificate, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestPRACredentialController(client *testClient) error {
	var errorList []error
	credential, _, err := pracredential.GetAll(client.sdkClient.PRACredential)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(credential)))
	for _, b := range credential {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := pracredential.Delete(client.sdkClient.PRACredential, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAPRACredentialController, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestPRAConsoleController(client *testClient) error {
	var errorList []error
	console, _, err := praconsole.GetAll(client.sdkClient.PRAConsole)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(console)))
	for _, b := range console {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := praconsole.Delete(client.sdkClient.PRAConsole, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAPRAConsoleController, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestPRAPortalController(client *testClient) error {
	var errorList []error
	portal, _, err := praportal.GetAll(client.sdkClient.PRAPortal)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(portal)))
	for _, b := range portal {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := praportal.Delete(client.sdkClient.PRAPortal, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAPRAPortalController, fmt.Sprintf(b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestPRAPrivilegedApprovalController(client *testClient) error {
	var errorList []error
	// First, get all pra approval resources
	approvals, _, err := praapproval.GetAll(client.sdkClient.PRAApproval)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(approvals)))

	for _, approval := range approvals {
		// Iterate over email_ids of each resource to check for the "pra_user_" substring
		for _, emailID := range approval.EmailIDs {
			if strings.Contains(emailID, "pra_user_") {
				// If the emailID contains "pra_user_", delete the resource
				if _, err := praapproval.Delete(client.sdkClient.PRAApproval, approval.ID); err != nil {
					errorList = append(errorList, err)
					continue
				}
				logSweptResource(resourcetype.ZPAPRAApprovalController, fmt.Sprintf(approval.ID), strings.Join(approval.EmailIDs, ", "))
				break // Exit the loop after deletion to avoid multiple attempts
			}
		}
	}

	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}
