package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/lssconfigcontroller"
)

func dataSourceLSSLogTypeFormats() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceLSSLogTypeFormatsRead,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"log_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"zpn_ast_comprehensive_stats",
					"zpn_auth_log",
					"zpn_pbroker_comprehensive_stats",
					"zpn_ast_auth_log",
					"zpn_audit_log",
					"zpn_trans_log",
					"zpn_http_trans_log",
					"zpn_waf_http_exchanges_log",
					"zpn_sys_auth_log",
				}, false),
			},
			"tsv": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"csv": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getLogType(d *schema.ResourceData) (string, bool) {
	val, ok := d.GetOk("log_type")
	if !ok {
		return "", ok
	}
	value, ok := val.(string)
	return value, ok
}

func dataSourceLSSLogTypeFormatsRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.LSSConfigController

	log.Printf("[INFO] Getting data for LSS Log Types Format set\n")
	logType, ok := getLogType(d)
	if !ok {
		return fmt.Errorf("[ERROR] log type is required")
	}
	resp, _, err := lssconfigcontroller.GetFormats(service, logType)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Getting LSS Log Types Format:\n%+v\n", resp)
	d.SetId("lss_log_types_" + logType)
	_ = d.Set("tsv", resp.Tsv)
	_ = d.Set("csv", resp.Csv)
	_ = d.Set("json", resp.Json)

	return nil
}
