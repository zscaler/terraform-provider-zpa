package zpa

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/clienttypes"
)

func dataSourceAccessPolicyClientTypes() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceAccessPolicyClientTypesRead,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"zpn_client_type_exporter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_client_type_exporter_noauth": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_client_type_browser_isolation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_client_type_machine_tunnel": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_client_type_ip_anchoring": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_client_type_edge_connector": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_client_type_zapp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_client_type_slogger": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_client_type_branch_connector": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAccessPolicyClientTypesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.ClientTypes

	log.Printf("[INFO] Getting data for all client types set\n")

	resp, _, err := clienttypes.GetAllClientTypes(service)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Getting data for all client types:\n%+v\n", resp)
	d.SetId("client_types")
	_ = d.Set("zpn_client_type_exporter", resp.ZPNClientTypeExplorer)
	_ = d.Set("zpn_client_type_exporter_noauth", resp.ZPNClientTypeNoAuth)
	_ = d.Set("zpn_client_type_browser_isolation", resp.ZPNClientTypeBrowserIsolation)
	_ = d.Set("zpn_client_type_machine_tunnel", resp.ZPNClientTypeMachineTunnel)
	_ = d.Set("zpn_client_type_ip_anchoring", resp.ZPNClientTypeIPAnchoring)
	_ = d.Set("zpn_client_type_edge_connector", resp.ZPNClientTypeEdgeConnector)
	_ = d.Set("zpn_client_type_zapp", resp.ZPNClientTypeZAPP)
	_ = d.Set("zpn_client_type_slogger", resp.ZPNClientTypeSlogger)
	_ = d.Set("zpn_client_type_branch_connector", resp.ZPNClientTypeBranchConnector)

	return nil
}
