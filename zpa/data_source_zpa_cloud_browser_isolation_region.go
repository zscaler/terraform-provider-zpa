package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbiregions"
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
	var searchCriteria string

	// Check if searching by ID first
	id, idOk := d.Get("id").(string)
	if idOk && id != "" {
		log.Printf("[INFO] Getting CBI region by id: %s\n", id)
		searchCriteria = fmt.Sprintf("id=%s", id)

		// Get all CBI regions and find the one with matching ID
		allRegions, _, err := cbiregions.GetAll(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}

		for _, region := range allRegions {
			if region.ID == id {
				resp = &region
				break
			}
		}
	}

	// Check if searching by name (only if ID search didn't find anything)
	name, nameOk := d.Get("name").(string)
	if resp == nil && nameOk && name != "" {
		log.Printf("[INFO] Getting CBI region by name: %s\n", name)
		searchCriteria = fmt.Sprintf("name=%s", name)

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
		return diag.FromErr(fmt.Errorf("couldn't find any CBI region with %s", searchCriteria))
	}

	return nil
}
