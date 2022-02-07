package zpa

import (
	"encoding/json"
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
			"zpn_sys_auth_log": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func toMapString(v map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for key, val := range v {
		data, err := json.MarshalIndent(&val, "", " ")
		if err != nil {
			log.Printf("[ERROR] MarshalIndent failed %v\n", err)
			continue
		}
		result[key] = string(data)
	}

	return result
}
func dataSourceLSSStatusCodesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	log.Printf("[INFO] Getting data for LSS Status Codes set\n")

	resp, _, err := zClient.lssconfigcontroller.GetStatusCodes()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Getting LSS Status Codes:\n%+v\n", resp)
	d.SetId("lss_status_codes")
	_ = d.Set("zpn_auth_log", toMapString(resp.ZPNAstAuthLog))
	_ = d.Set("zpn_ast_auth_log", toMapString(resp.ZPNAstAuthLog))
	_ = d.Set("zpn_trans_log", toMapString(resp.ZPNTransLog))
	_ = d.Set("zpn_sys_auth_log", toMapString(resp.ZPNSysAuthLog))

	return nil
}
