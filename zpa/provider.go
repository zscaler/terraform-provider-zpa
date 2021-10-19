package zpa

import (
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFunc("ZPA_CLIENT_ID"),
				Description: "zpa client id",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: envDefaultFunc("ZPA_CLIENT_SECRET"),
				Description: "zpa client secret",
			},
			"customerid": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: envDefaultFunc("ZPA_CUSTOMER_ID"),
				Description: "zpa customer id",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			/*
			   terraform resource name: resource schema
			   resource formation: provider-resourcename-subresource
			*/
			"zpa_app_connector_group":    resourceAppConnectorGroup(),
			"zpa_application_server":     resourceApplicationServer(),
			"zpa_application_segment":    resourceApplicationSegment(),
			"zpa_server_group":           resourceServerGroup(),
			"zpa_segment_group":          resourceSegmentGroup(),
			"zpa_browser_access":         resourceBrowserAccess(),
			"zpa_policy_access_rule":     resourcePolicyAccessRule(),
			"zpa_policy_timeout_rule":    resourcePolicyTimeoutRule(),
			"zpa_policy_forwarding_rule": resourcePolicyForwardingRule(),
			//"zpa_provisioning_key":       resourceProvisioningKey(),
			"zpa_service_edge_group": resourceServiceEdgeGroup(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			// terraform date source name: data source schema
			"zpa_posture_profile":          dataSourcePostureProfile(),
			"zpa_trusted_network":          dataSourceTrustedNetwork(),
			"zpa_saml_attribute":           dataSourceSamlAttribute(),
			"zpa_scim_groups":              dataSourceScimGroup(),
			"zpa_scim_attribute_header":    dataSourceScimAttributeHeader(),
			"zpa_ba_certificate":           dataSourceBaCertificate(),
			"zpa_machine_group":            dataSourceMachineGroup(),
			"zpa_application_segment":      dataSourceApplicationSegment(),
			"zpa_application_server":       dataSourceApplicationServer(),
			"zpa_server_group":             dataSourceServerGroup(),
			"zpa_cloud_connector_group":    dataSourceCloudConnectorGroup(),
			"zpa_app_connector_group":      dataSourceAppConnectorGroup(),
			"zpa_segment_group":            dataSourceSegmentGroup(),
			"zpa_idp_controller":           dataSourceIdpController(),
			"zpa_global_access_policy":     dataSourceGlobalAccessPolicy(),
			"zpa_global_policy_timeout":    dataSourceGlobalPolicyTimeout(),
			"zpa_global_policy_forwarding": dataSourceGlobalPolicyForwarding(),
			//"zpa_provisioning_key":         dataSourceProvisioningKey(),
			"zpa_service_edge_group": dataSourceServiceEdgeGroup(),
		},
		ConfigureFunc: zscalerConfigure,
	}
}

func zscalerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[INFO] Initializing ZPA client")
	config := Config{
		ClientID:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
		CustomerID:   d.Get("customerid").(string),
	}

	return config.Client()
}

func envDefaultFunc(k string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(k); v != "" {
			return v, nil
		}

		return nil, nil
	}
}
