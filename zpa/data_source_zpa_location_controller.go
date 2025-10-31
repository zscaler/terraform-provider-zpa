package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/extranet_resource"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/location_controller"
)

func dataSourceLocationController() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLocationControllerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the location.",
			},
			"zia_er_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the extranet resource partner.",
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceLocationControllerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Get required parameters
	locationName, locationNameOk := d.Get("name").(string)
	ziaErName, ziaErNameOk := d.Get("zia_er_name").(string)

	// Ensure both name and zia_er_name are provided
	if (!locationNameOk && !ziaErNameOk) || (locationName == "" && ziaErName == "") {
		log.Printf("[INFO] Both location name and extranet resource name are required\n")
		return diag.FromErr(fmt.Errorf("both 'name' and 'zia_er_name' are required"))
	}

	// Step 1: Get the extranet resource to obtain zpnErID
	log.Printf("[INFO] Getting extranet resource by name: %s\n", ziaErName)
	extranetResp, _, err := extranet_resource.GetExtranetResourcePartnerByName(ctx, service, ziaErName)
	if err != nil || extranetResp == nil {
		log.Printf("[INFO] Couldn't find extranet resource by name: %s\n", ziaErName)
		return diag.FromErr(fmt.Errorf("error fetching extranet resource by name '%s': %w", ziaErName, err))
	}

	zpnErID := extranetResp.ID
	log.Printf("[INFO] Found extranet resource ID: %s\n", zpnErID)

	// Step 2: Use zpnErID to get locations (API returns a list)
	log.Printf("[INFO] Getting locations using zpnErID: %s\n", zpnErID)
	locations, _, err := location_controller.GetLocationExtranetResource(ctx, service, zpnErID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error fetching locations with zpnErID '%s': %w", zpnErID, err))
	}

	if len(locations) == 0 {
		return diag.FromErr(fmt.Errorf("no locations found for extranet resource '%s'", ziaErName))
	}

	// Step 3: Search for the specific location by name
	var foundLocation *common.CommonSummary
	for i := range locations {
		if locations[i].Name == locationName {
			foundLocation = &locations[i]
			break
		}
	}

	if foundLocation == nil {
		return diag.FromErr(fmt.Errorf("location with name '%s' not found for extranet resource '%s'", locationName, ziaErName))
	}

	// Step 4: Set the data
	d.SetId(foundLocation.ID)
	_ = d.Set("name", foundLocation.Name)
	_ = d.Set("enabled", foundLocation.Enabled)

	return nil
}
