package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/cloudbrowserisolation/cbiregions"
)

func dataSourceCBIRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCBIRegionsRead,
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

func dataSourceCBIRegionsRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.CBIRegions

	var resp *cbiregions.CBIRegions
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for cbi regions name %s\n", name)
		res, _, err := cbiregions.GetByName(service, name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)

	} else {
		return fmt.Errorf("couldn't find any cbi regions with name '%s'", name)
	}

	return nil
}
