package zpa

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Resource names, defined in place, used throughout the provider and tests
const (
	zpaBrowserAccess = "zpa_application_segment_browser_access"
)

func ZPAProvider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "zpa client id",
			},
			"client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Description:   "zpa client secret",
				ConflictsWith: []string{"private_key"},
			},
			"private_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Description:   "zpa private key",
				ConflictsWith: []string{"client_secret"},
			},
			"vanity_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Zscaler Vanity Domain",
			},
			"customer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "zpa customer id",
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "zpa microtenant ID",
			},
			"zscaler_cloud": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Zscaler Cloud Name",
			},
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cloud to use PRODUCTION, ZPATWO, BETA, GOV, GOVUS, PREVIEW, DEV, QA, QA2",
				DefaultFunc: schema.EnvDefaultFunc("ZPA_CLOUD", nil),
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(string)
					if strings.HasPrefix(v, "https://") {
						return
					}
					return validation.StringInSlice([]string{"PRODUCTION", "ZPATWO", "BETA", "GOV", "GOVUS", "PREVIEW", "DEV", "QA", "QA2"}, true)(val, key)
				},
			},
			"use_legacy_client": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enables interaction with the ZPA legacy API framework",
			},
			"http_proxy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alternate HTTP proxy of scheme://hostname or scheme://hostname:port format",
			},
			"backoff": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Use exponential back off strategy for rate limits.",
			},
			"min_wait_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "minimum seconds to wait when rate limit is hit. We use exponential backoffs when backoff is enabled.",
			},
			"max_wait_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "maximum seconds to wait when rate limit is hit. We use exponential backoffs when backoff is enabled.",
			},
			"max_retries": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intAtMost(100),
				Description:      "maximum number of retries to attempt before erroring out.",
			},
			"parallelism": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of concurrent requests to make within a resource where bulk operations are not possible. Take note of https://help.zscaler.com/zpa/understanding-rate-limiting.",
			},
			"request_timeout": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(0, 300),
				Description:      "Timeout for single request (in seconds) which is made to Zscaler, the default is `0` (means no limit is set). The maximum value can be `300`.",
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
			"zpa_application_segment_multimatch_bulk":      resourceApplicationSegmentMultimatchBulk(),
			"zpa_application_segment_pra":                  resourceApplicationSegmentPRA(),
			"zpa_application_segment_inspection":           resourceApplicationSegmentInspection(),
			"zpa_application_segment_browser_access":       resourceApplicationSegmentBrowserAccess(),
			"zpa_ba_certificate":                           resourceBaCertificate(),
			"zpa_cloud_browser_isolation_certificate":      resourceCBICertificates(),
			"zpa_cloud_browser_isolation_external_profile": resourceCBIExternalProfile(),
			"zpa_cloud_browser_isolation_banner":           resourceCBIBanners(),
			"zpa_emergency_access_user":                    resourceEmergencyAccess(),
			"zpa_segment_group":                            resourceSegmentGroup(),
			"zpa_server_group":                             resourceServerGroup(),
			"zpa_policy_access_rule_reorder":               resourcePolicyAccessRuleReorder(),
			"zpa_policy_access_rule":                       resourcePolicyAccessRule(),
			"zpa_policy_browser_protection_rule":           resourcePolicyBrowserProtectionRule(),
			"zpa_policy_inspection_rule":                   resourcePolicyInspectionRule(),
			"zpa_policy_timeout_rule":                      resourcePolicyTimeoutRule(),
			"zpa_policy_forwarding_rule":                   resourcePolicyForwardingRule(),
			"zpa_policy_isolation_rule":                    resourcePolicyIsolationRule(),
			"zpa_policy_redirection_rule":                  resourcePolicyRedictionRule(),
			"zpa_policy_access_rule_v2":                    resourcePolicyAccessRuleV2(),
			"zpa_policy_isolation_rule_v2":                 resourcePolicyIsolationRuleV2(),
			"zpa_policy_inspection_rule_v2":                resourcePolicyInspectionRuleV2(),
			"zpa_policy_forwarding_rule_v2":                resourcePolicyForwardingRuleV2(),
			"zpa_policy_timeout_rule_v2":                   resourcePolicyTimeoutRuleV2(),
			"zpa_policy_credential_rule":                   resourcePolicyCredentialAccessRule(),
			"zpa_policy_capabilities_rule":                 resourcePolicyCapabilitiesAccessRule(),
			"zpa_policy_portal_access_rule":                resourcePolicyPortalAccessRule(),
			"zpa_provisioning_key":                         resourceProvisioningKey(),
			"zpa_service_edge_group":                       resourceServiceEdgeGroup(),
			"zpa_service_edge_assistant_schedule":          resourceServiceEdgeAssistantSchedule(),
			"zpa_lss_config_controller":                    resourceLSSConfigController(),
			"zpa_inspection_custom_controls":               resourceInspectionCustomControls(),
			"zpa_inspection_profile":                       resourceInspectionProfile(),
			"zpa_microtenant_controller":                   resourceMicrotenantController(),
			"zpa_pra_approval_controller":                  resourcePRAPrivilegedApprovalController(),
			"zpa_pra_portal_controller":                    resourcePRAPortalController(),
			"zpa_pra_credential_controller":                resourcePRACredentialController(),
			"zpa_pra_credential_pool":                      resourcePRACredentialPool(),
			"zpa_pra_console_controller":                   resourcePRAConsoleController(),
			"zpa_private_cloud_group":                      resourcePrivateCloudGroup(),
			"zpa_user_portal_controller":                   resourceUserPortalController(),
			"zpa_user_portal_link":                         resourceUserPortalLink(),
			"zpa_user_portal_aup":                          resourceUserPortalAUP(),
			"zpa_c2c_ip_ranges":                            resourceC2CIPRanges(),
			//"zpa_browser_protection":                       resourceBrowserProtection(),
			"zpa_zia_cloud_config": resourceZiaCloudConfig(),

			// The day I realized I was naming stuff wrong :'-(
			"zpa_browser_access": deprecateIncorrectNaming(resourceApplicationSegmentBrowserAccess(), zpaBrowserAccess),
		},
		DataSourcesMap: map[string]*schema.Resource{
			// terraform data source name: data source schema
			"zpa_application_server":                       dataSourceApplicationServer(),
			"zpa_application_segment":                      dataSourceApplicationSegment(),
			"zpa_application_segment_multimatch_bulk":      dataSourceApplicationSegmentMultimatchBulk(),
			"zpa_application_segment_pra":                  dataSourceApplicationSegmentPRA(),
			"zpa_application_segment_inspection":           dataSourceApplicationSegmentInspection(),
			"zpa_application_segment_browser_access":       dataSourceApplicationSegmentBrowserAccess(),
			"zpa_application_segment_by_type":              dataSourceApplicationSegmentByType(),
			"zpa_segment_group":                            dataSourceSegmentGroup(),
			"zpa_app_connector_group":                      dataSourceAppConnectorGroup(),
			"zpa_app_connector_controller":                 dataSourceAppConnectorController(),
			"zpa_app_connector_assistant_schedule":         dataSourceAppConnectorAssistantSchedule(),
			"zpa_ba_certificate":                           dataSourceBaCertificate(),
			"zpa_customer_version_profile":                 dataSourceCustomerVersionProfile(),
			"zpa_cloud_connector_group":                    dataSourceCloudConnectorGroup(),
			"zpa_branch_connector_group":                   dataSourceBranchConnectorGroup(),
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
			"zpa_service_edge_assistant_schedule":          dataSourceServiceEdgeAssistantSchedule(),
			"zpa_saml_attribute":                           dataSourceSamlAttribute(),
			"zpa_scim_groups":                              dataSourceScimGroup(),
			"zpa_scim_attribute_header":                    dataSourceScimAttributeHeader(),
			"zpa_server_group":                             dataSourceServerGroup(),
			"zpa_enrollment_cert":                          dataSourceEnrollmentCert(),
			"zpa_trusted_network":                          dataSourceTrustedNetwork(),
			"zpa_access_policy_platforms":                  dataSourceAccessPolicyPlatforms(),
			"zpa_access_policy_client_types":               dataSourceAccessPolicyClientTypes(),
			"zpa_risk_score_values":                        dataSourceRiskScoreValues(),
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
			"zpa_pra_credential_pool":                      dataSourcePRACredentialPool(),
			"zpa_pra_console_controller":                   dataSourcePRAConsoleController(),
			"zpa_private_cloud_group":                      dataSourcePrivateCloudGroup(),
			"zpa_private_cloud_controller":                 dataSourcePrivateCloudController(),
			"zpa_user_portal_controller":                   dataSourceUserPortalController(),
			"zpa_user_portal_link":                         dataSourceUserPortalLink(),
			"zpa_user_portal_aup":                          dataSourceUserPortalAUP(),
			"zpa_location_controller":                      dataSourceLocationController(),
			"zpa_location_group_controller":                dataSourceLocationGroupController(),
			"zpa_location_controller_summary":              dataSourceLocationControllerSummary(),
			"zpa_c2c_ip_ranges":                            dataSourceC2CIPRanges(),
			"zpa_extranet_resource_partner":                dataSourceExtranetResourcePartner(),
			"zpa_managed_browser_profile":                  dataSourceManagedBrowserProfile(),
			"zpa_browser_protection":                       dataSourceBrowserProtection(),
			"zpa_workload_tag_group":                       dataSourceWorkloadTagGroup(),
			"zpa_zia_cloud_config":                         dataSourceZiaCloudConfig(),
		},
	}

	p.ConfigureContextFunc = func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		r, err := providerConfigure(d, terraformVersion)
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

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	log.Printf("[INFO] Initializing Zscaler client")

	// Create a new configuration
	config := NewConfig(d)
	config.TerraformVersion = terraformVersion

	// Load the correct SDK client (prioritizing V3)
	if diags := config.loadClients(); diags.HasError() {
		return nil, diags
	}

	// Ensure the Client instance is returned
	client, err := config.Client()
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("failed to initialize client: %w", err))
	}

	return client, nil
}

func deprecateIncorrectNaming(d *schema.Resource, newResource string) *schema.Resource {
	d.DeprecationMessage = fmt.Sprintf("Resource is deprecated due to a correction in naming conventions, please use '%s' instead.", newResource)
	return d
}

func resourceFuncNoOp(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}
