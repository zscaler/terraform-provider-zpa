package zpa

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common"
)

// Resource names, defined in place, used throughout the provider and tests
const (
	zpaBrowserAccess = "zpa_application_segment_browser_access"
)

func ZPAProvider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"zpa_client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZPA_CLIENT_ID", nil),
				Description: "zpa client id",
			},
			"zpa_client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ZPA_CLIENT_SECRET", nil),
				Description: "zpa client secret",
			},
			"zpa_customer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ZPA_CUSTOMER_ID", nil),
				Description: "zpa customer id",
			},
			"zpa_cloud": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Cloud to use PRODUCTION, BETA, GOV, GOVUS, PREVIEW, DEV, QA, QA2",
				ValidateFunc: validation.StringInSlice([]string{"PRODUCTION", "BETA", "GOV", "GOVUS", "PREVIEW", "DEV", "QA", "QA2"}, true),
				DefaultFunc:  schema.EnvDefaultFunc("ZPA_CLOUD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			/*
			   terraform resource name: resource schema
			   resource formation: provider-resourcename-subresource
			*/
			"zpa_app_connector_assistant_schedule":         resourceAppConnectorAssistantSchedule(),
			"zpa_app_connector_group":                      resourceAppConnectorGroup(),
			"zpa_application_server":                       resourceApplicationServer(),
			"zpa_application_segment":                      resourceApplicationSegment(),
			"zpa_application_segment_pra":                  resourceApplicationSegmentPRA(),
			"zpa_application_segment_inspection":           resourceApplicationSegmentInspection(),
			"zpa_application_segment_browser_access":       resourceApplicationSegmentBrowserAccess(),
			"zpa_ba_certificate":                           resourceBaCertificate(),
			"zpa_cloud_browser_isolation_certificate":      resourceCBICertificates(),
			"zpa_cloud_browser_isolation_external_profile": resourceCBIExternalProfile(),
			"zpa_cloud_browser_isolation_banner":           resourceCBIBanners(),
			"zpa_segment_group":                            resourceSegmentGroup(),
			"zpa_server_group":                             resourceServerGroup(),
			"zpa_policy_access_rule_reorder":               resourcePolicyAccessRuleReorder(),
			"zpa_policy_access_rule":                       resourcePolicyAccessRule(),
			"zpa_policy_inspection_rule":                   resourcePolicyInspectionRule(),
			"zpa_policy_timeout_rule":                      resourcePolicyTimeoutRule(),
			"zpa_policy_forwarding_rule":                   resourcePolicyForwardingRule(),
			"zpa_policy_isolation_rule":                    resourcePolicyIsolationRule(),
			"zpa_policy_redirection_rule":                  resourcePolicyRedictionRule(),
			"zpa_provisioning_key":                         resourceProvisioningKey(),
			"zpa_service_edge_group":                       resourceServiceEdgeGroup(),
			"zpa_lss_config_controller":                    resourceLSSConfigController(),
			"zpa_inspection_custom_controls":               resourceInspectionCustomControls(),
			"zpa_inspection_profile":                       resourceInspectionProfile(),
			"zpa_microtenant_controller":                   resourceMicrotenantController(),
			"zpa_pra_approval_controller":                  resourcePRAPrivilegedApprovalController(),
			"zpa_pra_portal_controller":                    resourcePRAPortalController(),
			"zpa_pra_credential_controller":                resourcePRACredentialController(),
			"zpa_pra_console_controller":                   resourcePRAConsoleController(),

			// The day I realized I was naming stuff wrong :'-(
			"zpa_browser_access": deprecateIncorrectNaming(resourceApplicationSegmentBrowserAccess(), zpaBrowserAccess),
		},
		DataSourcesMap: map[string]*schema.Resource{
			// terraform data source name: data source schema
			"zpa_application_server":                       dataSourceApplicationServer(),
			"zpa_application_segment":                      dataSourceApplicationSegment(),
			"zpa_application_segment_pra":                  dataSourceApplicationSegmentPRA(),
			"zpa_application_segment_inspection":           dataSourceApplicationSegmentInspection(),
			"zpa_application_segment_browser_access":       dataSourceApplicationSegmentBrowserAccess(),
			"zpa_segment_group":                            dataSourceSegmentGroup(),
			"zpa_app_connector_group":                      dataSourceAppConnectorGroup(),
			"zpa_app_connector_controller":                 dataSourceAppConnectorController(),
			"zpa_app_connector_assistant_schedule":         dataSourceAppConnectorAssistantSchedule(),
			"zpa_ba_certificate":                           dataSourceBaCertificate(),
			"zpa_customer_version_profile":                 dataSourceCustomerVersionProfile(),
			"zpa_cloud_connector_group":                    dataSourceCloudConnectorGroup(),
			"zpa_idp_controller":                           dataSourceIdpController(),
			"zpa_machine_group":                            dataSourceMachineGroup(),
			"zpa_provisioning_key":                         dataSourceProvisioningKey(),
			"zpa_cloud_browser_isolation_region":           dataSourceCBIRegions(),
			"zpa_cloud_browser_isolation_certificate":      dataSourceCBICertificates(),
			"zpa_cloud_browser_isolation_zpa_profile":      dataSourceCBIZPAProfiles(),
			"zpa_cloud_browser_isolation_banner":           dataSourceCBIBanners(),
			"zpa_cloud_browser_isolation_external_profile": dataSourceCBIExternalProfile(),
			"zpa_policy_type":                              dataSourcePolicyType(),
			"zpa_isolation_profile":                        dataSourceIsolationProfile(),
			"zpa_posture_profile":                          dataSourcePostureProfile(),
			"zpa_service_edge_group":                       dataSourceServiceEdgeGroup(),
			"zpa_service_edge_controller":                  dataSourceServiceEdgeController(),
			"zpa_saml_attribute":                           dataSourceSamlAttribute(),
			"zpa_scim_groups":                              dataSourceScimGroup(),
			"zpa_scim_attribute_header":                    dataSourceScimAttributeHeader(),
			"zpa_server_group":                             dataSourceServerGroup(),
			"zpa_enrollment_cert":                          dataSourceEnrollmentCert(),
			"zpa_trusted_network":                          dataSourceTrustedNetwork(),
			"zpa_access_policy_platforms":                  dataSourceAccessPolicyPlatforms(),
			"zpa_access_policy_client_types":               dataSourceAccessPolicyClientTypes(),
			"zpa_lss_config_controller":                    dataSourceLSSConfigController(),
			"zpa_lss_config_client_types":                  dataSourceLSSClientTypes(),
			"zpa_lss_config_status_codes":                  dataSourceLSSStatusCodes(),
			"zpa_lss_config_log_type_formats":              dataSourceLSSLogTypeFormats(),
			"zpa_inspection_predefined_controls":           dataSourceInspectionPredefinedControls(),
			"zpa_inspection_all_predefined_controls":       dataSourceInspectionAllPredefinedControls(),
			"zpa_inspection_custom_controls":               dataSourceInspectionCustomControls(),
			"zpa_inspection_profile":                       dataSourceInspectionProfile(),
			"zpa_microtenant_controller":                   dataSourceMicrotenantController(),
			"zpa_pra_approval_controller":                  dataSourcePRAPrivilegedApprovalController(),
			"zpa_pra_portal_controller":                    dataSourcePRAPortalController(),
			"zpa_pra_credential_controller":                dataSourcePRACredentialController(),
			"zpa_pra_console_controller":                   dataSourcePRAConsoleController(),
		},
	}
	p.ConfigureContextFunc = func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		r, err := zscalerConfigure(d, terraformVersion)
		if err != nil {
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "failed configuring the provider",
					Detail:        fmt.Sprintf("error:%v", err),
					AttributePath: cty.Path{},
				},
			}
		}
		return r, nil
	}

	return p
}

func deprecateIncorrectNaming(d *schema.Resource, newResource string) *schema.Resource {
	d.DeprecationMessage = fmt.Sprintf("Resource is deprecated due to a correction in naming conventions, please use '%s' instead.", newResource)
	return d
}

func zscalerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	log.Printf("[INFO] Initializing ZPA client")
	clientID := d.Get("zpa_client_id").(string)
	clientSecret := d.Get("zpa_client_secret").(string)
	customerID := d.Get("zpa_customer_id").(string)

	// Retrieve zpa_cloud; if not set, default to the production environment.
	baseURL := d.Get("zpa_cloud").(string)
	if baseURL == "" {
		baseURL = "PRODUCTION" // Replace "PRODUCTION" with the actual URL or identifier of the production environment.
	}

	config := Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		CustomerID:   customerID,
		BaseURL:      baseURL,
		UserAgent:    fmt.Sprintf("(%s %s) Terraform/%s Provider/%s Customer/%s", runtime.GOOS, runtime.GOARCH, terraformVersion, common.Version(), customerID),
	}

	client, err := config.Client()
	if err != nil {
		log.Printf("[ERROR] Error initializing ZPA client: %s", err)
		return nil, err
	}

	return client, nil
}
