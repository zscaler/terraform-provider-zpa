package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/trustednetwork"
)

func dataSourceTrustedNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTrustedNetworkRead,
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

func dataSourceTrustedNetworkRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.TrustedNetwork

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}
	var resp *trustednetwork.TrustedNetwork
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for trusted network %s\n", id)
		res, _, err := trustednetwork.Get(service, id)
		if err != nil {
			return err
		}
		resp = res

	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for trusted network name %s\n", name)
		res, _, err := trustednetwork.GetByName(service, name)
		if err != nil {
			return err
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
		return fmt.Errorf("couldn't find any trusted network with name '%s' or id '%s'", name, id)
	}

	return nil
}
