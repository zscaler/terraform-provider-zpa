package zpa

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentbrowseraccess"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentinspection"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentpra"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/bacertificate"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbibannercontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbicertificatecontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbiprofilecontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_custom_controls"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_profile"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/lssconfigcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praapproval"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praconsole"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredential"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praportal"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/provisioningkey"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgegroup"
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

// Adjusted testClient to contain only V3 client services.
type testClient struct {
	sdkV3Client *zscaler.Client
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

	// Handle the sdkV3ClientForTest() return values
	sdkClient, err := sdkV3ClientForTest()
	if err != nil {
		t.Fatalf("Failed to initialize SDK V3 client: %v", err)
	}

	testClient := &testClient{
		sdkV3Client: sdkClient,
	}

	// Call individual sweeper functions
	sweepTestAppConnectorGroup(testClient)
	sweepTestApplicationServer(testClient)
	sweepTestApplicationSegment(testClient)
	sweepTestApplicationSegmentBA(testClient)
	sweepTestApplicationInspection(testClient)
	sweepTestApplicationPRA(testClient)
	sweepTestInspectionCustomControl(testClient)
	sweepTestInspectionProfile(testClient)
	sweepTestLSSConfigController(testClient)
	sweepTestAccessPolicyRuleByType(testClient)
	sweepTestProvisioningKey(testClient)
	sweepTestServerGroup(testClient)
	sweepTestSegmentGroup(testClient)
	sweepTestServiceEdgeGroup(testClient)
	sweepTestCBIBanner(testClient)
	sweepTestCBIExternalProfile(testClient)
	sweepTestPRACredentialController(testClient)
	sweepTestPRAConsoleController(testClient)
	sweepTestPRAPortalController(testClient)
	sweepTestPRAPrivilegedApprovalController(testClient)
}

// Sets up sweeper to clean up dangling resources
func setupSweeper(resourceType string, del func(*testClient) error) {
	resource.AddTestSweepers(resourceType, &resource.Sweeper{
		Name: resourceType,
		F: func(_ string) error {
			// Retrieve the client and handle the error
			sdkClient, err := sdkV3ClientForTest()
			if err != nil {
				return fmt.Errorf("failed to initialize SDK V3 client for sweeper: %w", err)
			}

			// Pass the client to the deleter function
			return del(&testClient{sdkV3Client: sdkClient})
		},
	})
}

func sweepTestAppConnectorGroup(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	// Fetch all app connector groups
	group, _, err := appconnectorgroup.GetAll(context.Background(), service)
	if err != nil {
		return err
	}

	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))

	// Iterate over the groups and delete the ones with the required prefix
	for _, b := range group {
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := appconnectorgroup.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAAppConnectorGroup, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	server, _, err := appservercontroller.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(server)))
	for _, b := range server {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := appservercontroller.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAApplicationServer, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	appSegment, _, err := applicationsegment.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(appSegment)))
	for _, b := range appSegment {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := applicationsegment.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAApplicationSegment, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	appSegmentBA, _, err := applicationsegmentbrowseraccess.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(appSegmentBA)))
	for _, b := range appSegmentBA {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := applicationsegmentbrowseraccess.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAApplicationSegmentBrowserAccess, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	appInspection, _, err := applicationsegmentinspection.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(appInspection)))
	for _, b := range appInspection {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := applicationsegmentinspection.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAApplicationSegmentInspection, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	pra, _, err := applicationsegmentpra.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(pra)))
	for _, b := range pra {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := applicationsegmentpra.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAApplicationSegmentPRA, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	customControl, _, err := inspection_custom_controls.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(customControl)))
	for _, b := range customControl {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := inspection_custom_controls.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAInspectionCustomControl, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	profile, _, err := inspection_profile.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(profile)))
	for _, b := range profile {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := inspection_profile.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAInspectionProfile, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	lssConfig, _, err := lssconfigcontroller.GetAll(context.Background(), service)
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
			_, err := lssconfigcontroller.Delete(context.Background(), service, b.ID)
			if err != nil {
				// Check if the error is because the resource doesn't exist
				if strings.Contains(err.Error(), "resource.not.found") {
					sweeperLogger.Info(fmt.Sprintf("Resource %s with ID %s was already deleted.", resourcetype.ZPALSSController, (b.ID)))
					continue
				}
				// For any other error, append to the error list and continue
				errorList = append(errorList, err)
				continue
			}

			sweeperLogger.Info(fmt.Sprintf("Swept resource %s with ID %s named %s.", resourcetype.ZPALSSController, (b.ID), b.LSSConfig.Name))
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

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
		policySet, _, err := policysetcontroller.GetByPolicyType(context.Background(), service, policyType)
		if err != nil {
			// If we fail to get a PolicySetID for a specific policy type, append the error and continue to the next type
			errorList = append(errorList, fmt.Errorf("Failed to get PolicySetID for policy type %s: %v", policyType, err))
			continue
		}
		policySetID := policySet.ID

		// Fetch all rules for the current policy type
		rules, _, err := policysetcontroller.GetAllByType(context.Background(), service, policyType)
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
				if _, err := policysetcontroller.Delete(context.Background(), service, policySetID, rule.ID); err != nil {
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	provisioningKey, err := provisioningkey.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(provisioningKey)))
	for _, b := range provisioningKey {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			// Assuming 'AssociationType' is a field in the provisioningKey object
			if _, err := provisioningkey.Delete(context.Background(), service, b.AssociationType, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAProvisioningKey, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	group, _, err := segmentgroup.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := segmentgroup.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPASegmentGroup, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	group, _, err := servergroup.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := servergroup.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAServerGroup, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	group, _, err := serviceedgegroup.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := serviceedgegroup.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAServiceEdgeGroup, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	group, _, err := cbibannercontroller.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := cbibannercontroller.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPACBIBannerController, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	group, _, err := cbiprofilecontroller.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := cbiprofilecontroller.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPACBIExternalIsolationProfile, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	group, _, err := cbicertificatecontroller.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := cbicertificatecontroller.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPACBICertificate, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	group, _, err := bacertificate.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := bacertificate.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPABACertificate, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	credential, _, err := pracredential.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(credential)))
	for _, b := range credential {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := pracredential.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAPRACredentialController, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	console, _, err := praconsole.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(console)))
	for _, b := range console {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := praconsole.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAPRAConsoleController, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	portal, _, err := praportal.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(portal)))
	for _, b := range portal {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := praportal.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAPRAPortalController, (b.ID), b.Name)
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

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	// First, get all pra approval resources
	approvals, _, err := praapproval.GetAll(context.Background(), service)
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
				if _, err := praapproval.Delete(context.Background(), service, approval.ID); err != nil {
					errorList = append(errorList, err)
					continue
				}
				logSweptResource(resourcetype.ZPAPRAApprovalController, (approval.ID), strings.Join(approval.EmailIDs, ", "))
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
