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
			"skip_credentials_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "Skip credentials validation and API client initialization entirely. " +
					"Intended for configurations where every zpa_* resource and data source is conditionally disabled " +
					"(e.g., count = 0) and no API call will ever be made, such as multi-environment deployments where " +
					"Zscaler is not present in every environment. Any resource or data source that does attempt an API " +
					"call will fail with an explanatory error. Can also be sourced from the " +
					"ZSCALER_SKIP_CREDENTIALS_VALIDATION environment variable.",
			},
			"http_proxy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alternate HTTP proxy of scheme://hostname or scheme://hostname:port format",
			},
			"backoff": {
				Type:        schema.TypeBool,
				Optional:    true,
				Deprecated:  "This attribute has never had any effect and will be removed in a future major release. Remove it from the provider block. Retry back-off is handled automatically by the provider.",
				Description: "Deprecated and ignored. Retry back-off is handled automatically and this attribute has no effect.",
			},
			"min_wait_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Deprecated:  "This attribute no longer has any effect and will be removed in a future major release. Remove it from the provider block. Retry back-off is handled automatically: the provider honours the Retry-After interval returned by the API and grows the wait between retries on its own.",
				Description: "Deprecated and ignored. Previously set the minimum wait between retries of a rate-limited request. Retry back-off is now handled automatically and this attribute has no effect.",
			},
			"max_wait_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Deprecated:  "This attribute no longer has any effect and will be removed in a future major release. Remove it from the provider block. Retry back-off is handled automatically: the provider honours the Retry-After interval returned by the API and grows the wait between retries on its own.",
				Description: "Deprecated and ignored. Previously set the maximum wait between retries of a rate-limited request. Retry back-off is now handled automatically and this attribute has no effect.",
			},
			"max_retries": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intAtMost(100),
				Description:      "Maximum number of times a rate-limited or transiently failing request is retried before the operation fails. Default and maximum: 100. Does not affect request throughput; lowering it only makes runs fail sooner under sustained rate limiting.",
			},
			"parallelism": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "This attribute no longer has any effect and will be removed in a future major release. Remove it from the provider block. API rate limits are handled automatically: the provider honours the Retry-After header returned on a 429 response and retries transparently.",
				Description: "Deprecated and ignored. Previously limited the number of concurrent API requests. " +
					"Rate limiting is now handled automatically and this attribute has no effect.",
			},
			"request_timeout": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(0, 300),
				Description:      "Timeout in seconds for a single HTTP request to the Zscaler API. Default: 240. Accepted values: 0-300, where 0 selects the SDK's built-in 60-second timeout. Does not affect request throughput.",
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
			"zpa_application_segment_weightedlb_config":    resourceApplicationSegmentWeightedLBConfig(),
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
			"zpa_private_cloud":                            resourcePrivateCloud(),
			"zpa_user_portal_controller":                   resourceUserPortalController(),
			"zpa_user_portal_link":                         resourceUserPortalLink(),
			"zpa_user_portal_aup":                          resourceUserPortalAUP(),
			"zpa_c2c_ip_ranges":                            resourceC2CIPRanges(),
			//"zpa_browser_protection":                       resourceBrowserProtection(),
			"zpa_tag_namespace":    resourceTagNamespace(),
			"zpa_tag_key":          resourceTagKey(),
			"zpa_tag_group":        resourceTagGroup(),
			"zpa_zia_cloud_config": resourceZiaCloudConfig(),

			// The day I realized I was naming stuff wrong :'-(
			"zpa_browser_access": deprecateIncorrectNaming(resourceApplicationSegmentBrowserAccess(), zpaBrowserAccess),
		},
		DataSourcesMap: map[string]*schema.Resource{
			// terraform data source name: data source schema
			"zpa_application_server":                       dataSourceApplicationServer(),
			"zpa_application_segment":                      dataSourceApplicationSegment(),
			"zpa_application_segment_multimatch_bulk":      dataSourceApplicationSegmentMultimatchBulk(),
			"zpa_application_segment_weightedlb_config":    dataSourceApplicationSegmentWeightedLBConfig(),
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
			"zpa_private_cloud":                            dataSourcePrivateCloud(),
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
			"zpa_tag_namespace":                            dataSourceTagNamespace(),
			"zpa_tag_key":                                  dataSourceTagKey(),
			"zpa_tag_group":                                dataSourceTagGroup(),
			"zpa_workload_tag_group":                       dataSourceWorkloadTagGroup(),
			"zpa_zia_cloud_config":                         dataSourceZiaCloudConfig(),
		},
	}

	// Guard every resource and data source against the inert client returned
	// when skip_credentials_validation is enabled, so an accidental API call
	// yields a descriptive error instead of a nil-pointer panic.
	for _, r := range p.ResourcesMap {
		guardResourceAgainstInertClient(r)
	}
	for _, ds := range p.DataSourcesMap {
		guardResourceAgainstInertClient(ds)
	}

	p.ConfigureContextFunc = func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		r, diags := providerConfigure(d, terraformVersion)
		if diags.HasError() {
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "failed configuring the provider",
					Detail:        fmt.Sprintf("error:%v", diags),
					AttributePath: cty.Path{},
				},
			}
		}
		// Pass through non-error diagnostics (e.g. the
		// skip_credentials_validation warning).
		return r, diags
	}

	return p
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	log.Printf("[INFO] Initializing Zscaler client")

	// Create a new configuration
	config := NewConfig(d)
	config.TerraformVersion = terraformVersion

	// Skip mode: never construct the SDK client. The OneAPI SDK authenticates
	// inside its constructor, so building a client without valid credentials
	// would fail at configure time even when no resource will ever call the
	// API. Return an inert client instead; the CRUD guards installed in
	// ZPAProvider turn any actual API use into a descriptive error.
	if config.skipCredentialsValidation {
		log.Printf("[WARN] skip_credentials_validation is enabled; ZPA API client was not initialized")
		return &Client{
				skipCredentialsValidation: true,
				policySetIDCache:          make(map[string]string),
			}, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "ZPA credentials were not validated",
					Detail: "skip_credentials_validation is enabled, so the ZPA API client was not initialized. " +
						"Any zpa_* resource or data source that attempts an API call will fail. " +
						"This mode is intended for configurations where all ZPA resources are conditionally disabled (e.g., count = 0).",
				},
			}
	}

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

// inertClientDiag is the error returned when a resource or data source is
// evaluated while the provider is in skip_credentials_validation mode.
func inertClientDiag() diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ZPA provider was configured with skip_credentials_validation",
			Detail: "This resource or data source attempted a ZPA API call, but the provider was configured with " +
				"skip_credentials_validation = true (or ZSCALER_SKIP_CREDENTIALS_VALIDATION=true), so no API client exists. " +
				"Either provide valid credentials and remove skip_credentials_validation, or ensure every zpa_* resource " +
				"and data source is disabled (e.g., count = 0) in this configuration.",
		},
	}
}

// isInertClient reports whether the provider meta is the inert client created
// in skip_credentials_validation mode.
func isInertClient(meta interface{}) bool {
	client, ok := meta.(*Client)
	return ok && client.skipCredentialsValidation
}

// guardResourceAgainstInertClient wraps a resource's CRUD and import
// functions so that, in skip_credentials_validation mode, they return a
// descriptive error instead of dereferencing the nil SDK service.
func guardResourceAgainstInertClient(r *schema.Resource) {
	wrap := func(f func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
		if f == nil {
			return nil
		}
		return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			if isInertClient(meta) {
				return inertClientDiag()
			}
			return f(ctx, d, meta)
		}
	}

	r.CreateContext = wrap(r.CreateContext)
	r.ReadContext = wrap(r.ReadContext)
	r.UpdateContext = wrap(r.UpdateContext)
	r.DeleteContext = wrap(r.DeleteContext)

	if r.Importer != nil && r.Importer.StateContext != nil {
		importer := r.Importer.StateContext
		r.Importer.StateContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			if isInertClient(meta) {
				return nil, fmt.Errorf("cannot import: the ZPA provider was configured with skip_credentials_validation, so no API client exists")
			}
			return importer(ctx, d, meta)
		}
	}
}

func deprecateIncorrectNaming(d *schema.Resource, newResource string) *schema.Resource {
	d.DeprecationMessage = fmt.Sprintf("Resource is deprecated due to a correction in naming conventions, please use '%s' instead.", newResource)
	return d
}

func resourceFuncNoOp(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}
