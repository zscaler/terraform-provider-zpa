package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/lssconfigcontroller"
)

func dataSourceLSSLogTypeFormats() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLSSLogTypeFormatsRead,
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
					"zpn_smb_inspection_log",
					"zpn_auth_log_1id",
					"zpn_sitec_auth_log",
					"zpn_sitec_comprehensive_stats",
					"zpn_ldap_inspection_log",
					"zms_flow_log",
					"zpn_krb_inspection_log",
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

func dataSourceLSSLogTypeFormatsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	log.Printf("[INFO] Getting data for LSS Log Types Format set\n")
	logType, ok := getLogType(d)
	if !ok {
		return diag.FromErr(fmt.Errorf("[ERROR] log type is required"))
	}
	resp, _, err := lssconfigcontroller.GetFormats(ctx, service, logType)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting LSS Log Types Format:\n%+v\n", resp)
	d.SetId("lss_log_types_" + logType)
	_ = d.Set("tsv", resp.Tsv)
	_ = d.Set("csv", resp.Csv)
	_ = d.Set("json", resp.Json)

	return nil
}
