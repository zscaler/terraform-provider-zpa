package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbiregions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCBIRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCBIRegionsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceCBIRegionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *cbiregions.CBIRegions
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for cbi regions name %s\n", name)
		res, _, err := cbiregions.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any cbi regions with name '%s'", name))
	}

	return nil
}
