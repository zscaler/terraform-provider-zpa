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
	"github.com/zscaler/terraform-provider-zpa/v2/zpa/common/resourcetype"
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

var testResourcePrefix = "tf-acc-test-"

func TestRunForcedSweeper(t *testing.T) {
	if os.Getenv("ZPA_VCR_TF_ACC") != "" {
		t.Skip("forced sweeper is live and will never be run within VCR")
		return
	}
	if os.Getenv("ZPA_ACC_TEST_FORCE_SWEEPERS") == "" || os.Getenv("TF_ACC") == "" {
		t.Skipf("ENV vars %q and %q must not be blank to force running of the sweepers", "ZPA_ACC_TEST_FORCE_SWEEPERS", "TF_ACC")
		return
	}

	provider := Provider()
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
	sweepTestSegmentGroup(testClient)
	sweepTestServerGroup(testClient)
	sweepTestServiceEdgeGroup(testClient)

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

func sweepTestAppConnectorGroup(client *testClient) error {
	var errorList []error
	group, _, err := client.sdkClient.appconnectorgroup.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) {
			if _, err := client.sdkClient.appconnectorgroup.Delete(b.ID); err != nil {
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
	server, _, err := client.sdkClient.appservercontroller.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(server)))
	for _, b := range server {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) {
			if _, err := client.sdkClient.appservercontroller.Delete(b.ID); err != nil {
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
	appSegment, _, err := client.sdkClient.applicationsegment.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(appSegment)))
	for _, b := range appSegment {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) {
			if _, err := client.sdkClient.applicationsegment.Delete(b.ID); err != nil {
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
	appSegmentBA, _, err := client.sdkClient.browseraccess.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(appSegmentBA)))
	for _, b := range appSegmentBA {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) {
			if _, err := client.sdkClient.browseraccess.Delete(b.ID); err != nil {
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
	appInspection, _, err := client.sdkClient.applicationsegmentinspection.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(appInspection)))
	for _, b := range appInspection {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) {
			if _, err := client.sdkClient.applicationsegmentinspection.Delete(b.ID); err != nil {
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
	pra, _, err := client.sdkClient.applicationsegmentpra.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(pra)))
	for _, b := range pra {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) {
			if _, err := client.sdkClient.applicationsegmentpra.Delete(b.ID); err != nil {
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
	customControl, _, err := client.sdkClient.inspection_custom_controls.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(customControl)))
	for _, b := range customControl {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) {
			if _, err := client.sdkClient.inspection_custom_controls.Delete(b.ID); err != nil {
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
	profile, _, err := client.sdkClient.inspection_profile.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(profile)))
	for _, b := range profile {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) {
			if _, err := client.sdkClient.inspection_profile.Delete(b.ID); err != nil {
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
	lssConfig, _, err := client.sdkClient.lssconfigcontroller.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(lssConfig)))
	for _, b := range lssConfig {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.LSSConfig.Name, testResourcePrefix) {
			if _, err := client.sdkClient.lssconfigcontroller.Delete(b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZPAInspectionProfile, fmt.Sprintf(b.ID), b.LSSConfig.Name)
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

func sweepTestAccessPolicyRuleByType(client *testClient) error {
	var errorList []error

	policyTypes := []string{
		"ACCESS_POLICY",
		"TIMEOUT_POLICY",
		"CLIENT_FORWARDING_POLICY",
		"INSPECTION_POLICY",
		"ISOLATION_POLICY",
	}

	for _, policyType := range policyTypes {
		// Fetch the PolicySet details for the current policy type to get the PolicySetID
		policySet, _, err := client.sdkClient.policysetcontroller.GetByPolicyType(policyType)
		if err != nil {
			// If we fail to get a PolicySetID for a specific policy type, append the error and continue to the next type
			errorList = append(errorList, fmt.Errorf("Failed to get PolicySetID for policy type %s: %v", policyType, err))
			continue
		}
		policySetID := policySet.ID

		// Fetch all rules for the current policy type
		rules, _, err := client.sdkClient.policysetcontroller.GetAllByType(policyType)
		if err != nil {
			// If we fail to fetch rules for a specific policy type, append the error and continue to the next type
			errorList = append(errorList, fmt.Errorf("Failed to get rules for policy type %s: %v", policyType, err))
			continue
		}
		// Logging the number of identified resources before the deletion loop for the current policy type
		sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep for policy type %s", len(rules), policyType))

		for _, rule := range rules {
			// Check if the resource name has the required prefix before deleting it
			if strings.HasPrefix(rule.Name, testResourcePrefix) {
				// Use the fetched PolicySetID for deletion
				if _, err := client.sdkClient.policysetcontroller.Delete(policySetID, rule.ID); err != nil {
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
	provisioningKey, err := client.sdkClient.provisioningkey.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(provisioningKey)))
	for _, b := range provisioningKey {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) {
			// Assuming 'AssociationType' is a field in the provisioningKey object
			if _, err := client.sdkClient.provisioningkey.Delete(b.AssociationType, b.ID); err != nil {
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
	group, _, err := client.sdkClient.segmentgroup.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) {
			if _, err := client.sdkClient.segmentgroup.Delete(b.ID); err != nil {
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
	group, _, err := client.sdkClient.servergroup.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) {
			if _, err := client.sdkClient.servergroup.Delete(b.ID); err != nil {
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
	group, _, err := client.sdkClient.serviceedgegroup.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(group)))
	for _, b := range group {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) {
			if _, err := client.sdkClient.serviceedgegroup.Delete(b.ID); err != nil {
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
