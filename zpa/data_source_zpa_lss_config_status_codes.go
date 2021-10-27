package zpa

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLSSStatusCodes() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceLSSStatusCodesRead,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"zpn_auth_log": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"zpn_ast_auth_log": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"zpn_trans_log": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceLSSStatusCodesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	log.Printf("[INFO] Getting data for LSS Status Codes set\n")

	resp, _, err := zClient.lssconfigcontroller.GetStatusCodes()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Getting LSS Status Codes:\n%+v\n", resp)
	// d.SetId(resp.ID)
	_ = d.Set("zpn_auth_log", resp.ZPNAstAuthLog)
	_ = d.Set("zpn_ast_auth_log", resp.ZPNAstAuthLog)
	_ = d.Set("zpn_trans_log", resp.ZPNTransLog)

	return nil
}
