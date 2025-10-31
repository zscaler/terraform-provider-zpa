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

func dataSourceLocationGroupController() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLocationGroupControllerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the location group (same as location_group_id).",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the location group.",
			},
			"location_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the location within the ziaLocations block to search for.",
			},
			"zia_er_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the extranet resource partner.",
			},
			"location_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the location group.",
			},
			"location_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the location group.",
			},
			"zia_locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceLocationGroupControllerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Get required parameters
	locationName, locationNameOk := d.Get("location_name").(string)
	ziaErName, ziaErNameOk := d.Get("zia_er_name").(string)

	// Ensure both location_name and zia_er_name are provided
	if (!locationNameOk && !ziaErNameOk) || (locationName == "" && ziaErName == "") {
		log.Printf("[INFO] Both location name and extranet resource name are required\n")
		return diag.FromErr(fmt.Errorf("both 'location_name' and 'zia_er_name' are required"))
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

	// Step 2: Use zpnErID to get location groups (API returns a list)
	log.Printf("[INFO] Getting location group using zpnErID: %s\n", zpnErID)
	locationGroups, _, err := location_controller.GetLocationGroupExtranetResource(ctx, service, zpnErID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error fetching location groups with zpnErID '%s': %w", zpnErID, err))
	}

	if len(locationGroups) == 0 {
		return diag.FromErr(fmt.Errorf("no location groups found for extranet resource '%s'", ziaErName))
	}

	// Step 3: Search for the specific location by name in ziaLocations across all location groups
	var foundLocation *common.CommonSummary
	var locationGroupResp *common.LocationGroupDTO

	for i := range locationGroups {
		for j := range locationGroups[i].ZiaLocations {
			if locationGroups[i].ZiaLocations[j].Name == locationName {
				foundLocation = &locationGroups[i].ZiaLocations[j]
				locationGroupResp = &locationGroups[i]
				break
			}
		}
		if foundLocation != nil {
			break
		}
	}

	if foundLocation == nil {
		return diag.FromErr(fmt.Errorf("location with name '%s' not found in location groups for extranet resource '%s'", locationName, ziaErName))
	}

	// Step 4: Set the data - ID is the LOCATION GROUP ID
	d.SetId(locationGroupResp.ID)             // Location Group ID
	_ = d.Set("name", locationGroupResp.Name) // Location Group Name
	_ = d.Set("location_group_id", locationGroupResp.ID)
	_ = d.Set("location_group_name", locationGroupResp.Name)
	_ = d.Set("zia_locations", flattenZiaLocations(locationGroupResp.ZiaLocations))

	return nil
}

// flattenZiaLocations flattens the ZiaLocations slice into a format suitable for Terraform
func flattenZiaLocations(ziaLocations []common.CommonSummary) []map[string]interface{} {
	if ziaLocations == nil {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(ziaLocations))
	for _, location := range ziaLocations {
		item := map[string]interface{}{
			"id":      location.ID,
			"name":    location.Name,
			"enabled": location.Enabled,
		}
		result = append(result, item)
	}

	return result
}
