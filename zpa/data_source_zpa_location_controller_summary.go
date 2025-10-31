package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/location_controller"
)

func dataSourceLocationControllerSummary() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLocationControllerSummaryRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the location.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the location.",
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceLocationControllerSummaryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *common.CommonSummary
	id, idOk := d.Get("id").(string)
	name, nameOk := d.Get("name").(string)

	// Ensure at least one of id or name is provided
	if (!idOk || id == "") && (!nameOk || name == "") {
		log.Printf("[INFO] Either location ID or name is required\n")
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	// Search by ID
	if idOk && id != "" {
		log.Printf("[INFO] Getting data for location controller ID %s\n", id)
		allLocations, _, err := location_controller.GetLocationSummary(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}
		// Find the location with matching ID
		for i := range allLocations {
			if allLocations[i].ID == id {
				resp = &allLocations[i]
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("location with ID '%s' not found", id))
		}
	} else if nameOk && name != "" {
		// Search by name
		log.Printf("[INFO] Getting data for location controller name %s\n", name)
		res, _, err := location_controller.GetLocationSummaryByName(ctx, service, name)
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
		return diag.FromErr(fmt.Errorf("couldn't find any location with name '%s' or id '%s'", name, id))
	}

	return nil
}
