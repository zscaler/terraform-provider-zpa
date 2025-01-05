package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/trustednetwork"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTrustedNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrustedNetworkRead,
		Schema: map[string]*schema.Schema{
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modifiedby": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zscaler_cloud": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTrustedNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var resp *trustednetwork.TrustedNetwork
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for trusted network %s\n", id)
		res, _, err := trustednetwork.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err) // Wrap error using diag.FromErr
		}
		resp = res

	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for trusted network name %s\n", name)
		res, _, err := trustednetwork.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("domain", resp.Domain)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("network_id", resp.NetworkID)
		_ = d.Set("zscaler_cloud", resp.ZscalerCloud)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any trusted network with name '%s' or id '%s'", name, id))
	}

	return nil
}
