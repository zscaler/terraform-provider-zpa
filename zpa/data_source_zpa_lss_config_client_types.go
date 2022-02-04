package zpa

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLSSClientTypes() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceLSSClientTypesRead,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"zpn_client_type_exporter": {
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
		},
	}
}

func dataSourceLSSClientTypesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	log.Printf("[INFO] Getting data for global policy set\n")

	resp, _, err := zClient.lssconfigcontroller.GetClientTypes()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Getting Policy Set Global Rules:\n%+v\n", resp)
	d.SetId("lss_client_types")
	_ = d.Set("zpn_client_type_exporter", resp.ZPNClientTypeExporter)
	_ = d.Set("zpn_client_type_machine_tunnel", resp.ZPNClientTypeMachineTunnel)
	_ = d.Set("zpn_client_type_ip_anchoring", resp.ZPNClientTypeIPAnchoring)
	_ = d.Set("zpn_client_type_edge_connector", resp.ZPNClientTypeEdgeConnector)
	_ = d.Set("zpn_client_type_zapp", resp.ZPNClientTypeZAPP)
	_ = d.Set("zpn_client_type_slogger", resp.ZPNClientTypeSlogger)

	return nil
}
