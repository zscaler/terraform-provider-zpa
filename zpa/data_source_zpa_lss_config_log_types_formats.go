package zpa

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLSSLotTypeFormats() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceLSSLotTypeFormatsRead,
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
			"zpn_audit_log": {
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
			"zpn_http_trans_log": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceLSSLotTypeFormatsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	log.Printf("[INFO] Getting data for LSS Log Types Format set\n")

	resp, _, err := zClient.lssconfigcontroller.GetFormats()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Getting LSS Log Types Format:\n%+v\n", resp)
	// d.SetId(resp.ID)
	_ = d.Set("zpn_auth_log", resp.ZPNAuthLog)
	_ = d.Set("zpn_ast_auth_log", resp.ZPNAstAuthLog)
	_ = d.Set("zpn_audit_log", resp.ZPNAuditLog)
	_ = d.Set("zpn_trans_log", resp.ZPNTransLog)
	_ = d.Set("zpn_http_trans_log", resp.ZPNHTTPTransLog)

	return nil
}
