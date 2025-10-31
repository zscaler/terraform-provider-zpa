package zpa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

func dataSourceRiskScoreValues() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRiskScoreValuesRead,
		Schema: map[string]*schema.Schema{
			"exclude_unknown": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceRiskScoreValuesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	excludeUnknown := d.Get("exclude_unknown").(bool)

	values, _, err := policysetcontrollerv2.GetRiskScoreValues(ctx, service, &excludeUnknown)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("risk_score_values")
	_ = d.Set("values", values)

	return nil
}
