package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/extranet_resource"
)

func dataSourceExtranetResourcePartner() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExtranetResourcePartnerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceExtranetResourcePartnerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *common.CommonSummary
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for extranet resource partner %s\n", id)
		// Get all extranet resource partners and find the one with matching ID
		allPartners, _, err := extranet_resource.GetExtranetResourcePartner(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, partner := range allPartners {
			if partner.ID == id {
				resp = &partner
				break
			}
		}
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for extranet resource partner%s\n", name)
		res, _, err := extranet_resource.GetExtranetResourcePartnerByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("enabled", resp.Enabled)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any extranet resource partner with name '%s' or id '%s'", name, id))
	}

	return nil
}
