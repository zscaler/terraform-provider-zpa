package zpa

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"zpa_client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "zpa client id",
			},
			"zpa_client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "zpa client secret",
			},
			"zpa_customer_id": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "zpa customer id",
			},
			"zpa_cloud": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Cloud to use PRODUCTION, BETA, GOV",
				ValidateFunc: validation.StringInSlice([]string{"PRODUCTION", "BETA", "GOV"}, true),
				Default:      "PRODUCTION",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			/*
			   terraform resource name: resource schema
			   resource formation: provider-resourcename-subresource
			*/
			"zpa_app_connector_group":            resourceAppConnectorGroup(),
			"zpa_application_server":             resourceApplicationServer(),
			"zpa_application_segment":            resourceApplicationSegment(),
			"zpa_application_segment_pra":        resourceApplicationSegmentPRA(),
			"zpa_application_segment_inspection": resourceApplicationSegmentInspection(),
			"zpa_segment_group":                  resourceSegmentGroup(),
			"zpa_server_group":                   resourceServerGroup(),
			"zpa_browser_access":                 resourceBrowserAccess(),
			"zpa_policy_access_rule":             resourcePolicyAccessRule(),
			"zpa_policy_inspection_rule":         resourcePolicyInspectionRule(),
			"zpa_policy_timeout_rule":            resourcePolicyTimeoutRule(),
			"zpa_policy_forwarding_rule":         resourcePolicyForwardingRule(),
			"zpa_provisioning_key":               resourceProvisioningKey(),
			"zpa_service_edge_group":             resourceServiceEdgeGroup(),
			"zpa_lss_config_controller":          resourceLSSConfigController(),
			"zpa_inspection_custom_controls":     resourceInspectionCustomControls(),
			"zpa_inspection_profile":             resourceInspectionProfile(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			// terraform data source name: data source schema
			"zpa_application_server":                 dataSourceApplicationServer(),
			"zpa_application_segment":                dataSourceApplicationSegment(),
			"zpa_application_segment_pra":            dataSourceApplicationSegmentPRA(),
			"zpa_application_segment_inspection":     dataSourceApplicationSegmentInspection(),
			"zpa_browser_access":                     dataSourceBrowserAccess(),
			"zpa_segment_group":                      dataSourceSegmentGroup(),
			"zpa_app_connector_group":                dataSourceAppConnectorGroup(),
			"zpa_app_connector_controller":           dataSourceAppConnectorController(),
			"zpa_ba_certificate":                     dataSourceBaCertificate(),
			"zpa_customer_version_profile":           dataSourceCustomerVersionProfile(),
			"zpa_cloud_connector_group":              dataSourceCloudConnectorGroup(),
			"zpa_idp_controller":                     dataSourceIdpController(),
			"zpa_machine_group":                      dataSourceMachineGroup(),
			"zpa_provisioning_key":                   dataSourceProvisioningKey(),
			"zpa_policy_type":                        dataSourcePolicyType(),
			"zpa_posture_profile":                    dataSourcePostureProfile(),
			"zpa_service_edge_group":                 dataSourceServiceEdgeGroup(),
			"zpa_service_edge_controller":            dataSourceServiceEdgeController(),
			"zpa_saml_attribute":                     dataSourceSamlAttribute(),
			"zpa_scim_groups":                        dataSourceScimGroup(),
			"zpa_scim_attribute_header":              dataSourceScimAttributeHeader(),
			"zpa_server_group":                       dataSourceServerGroup(),
			"zpa_enrollment_cert":                    dataSourceEnrollmentCert(),
			"zpa_trusted_network":                    dataSourceTrustedNetwork(),
			"zpa_lss_config_controller":              dataSourceLSSConfigController(),
			"zpa_lss_config_client_types":            dataSourceLSSClientTypes(),
			"zpa_lss_config_status_codes":            dataSourceLSSStatusCodes(),
			"zpa_lss_config_log_type_formats":        dataSourceLSSLogTypeFormats(),
			"zpa_inspection_predefined_controls":     dataSourceInspectionPredefinedControls(),
			"zpa_inspection_all_predefined_controls": dataSourceInspectionAllPredefinedControls(),
			"zpa_inspection_custom_controls":         dataSourceInspectionCustomControls(),
			"zpa_inspection_profile":                 dataSourceInspectionProfile(),
		},
		ConfigureFunc: zscalerConfigure,
	}
}

func zscalerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[INFO] Initializing ZPA client")
	config := Config{
		ClientID:     d.Get("zpa_client_id").(string),
		ClientSecret: d.Get("zpa_client_secret").(string),
		CustomerID:   d.Get("zpa_customer_id").(string),
		BaseURL:      d.Get("zpa_cloud").(string),
	}

	return config.Client()
}
