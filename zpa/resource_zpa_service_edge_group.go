package zpa

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgecontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgegroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/trustednetwork"
)

func resourceServiceEdgeGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceEdgeGroupCreate,
		ReadContext:   resourceServiceEdgeGroupRead,
		UpdateContext: resourceServiceEdgeGroupUpdate,
		DeleteContext: resourceServiceEdgeGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				microTenantID := GetString(d.Get("microtenant_id"))
				if microTenantID != "" {
					service = service.WithMicroTenant(microTenantID)
				}

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := serviceedgegroup.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(resp.ID)
						_ = d.Set("id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Service Edge Group.",
			},
			"city_country": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"country_code": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateCountryCode,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Service Edge Group.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether this Service Edge Group is enabled or not.",
			},
			"is_public": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable or disable public access for the Service Edge Group.",
			},
			"latitude": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     ValidateLatitude,
				DiffSuppressFunc: DiffSuppressFuncCoordinate,
				Description:      "Latitude for the Service Edge Group.",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Location for the Service Edge Group.",
			},
			"longitude": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     ValidateLongitude,
				DiffSuppressFunc: DiffSuppressFuncCoordinate,
				Description:      "Longitude for the Service Edge Group.",
			},
			"override_version_profile": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the default version profile of the App Connector Group is applied or overridden.",
			},
			"use_in_dr_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			// "service_edges": {
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	MaxItems: 1,
			// Description: "WARNING: Service edge membership is managed externally. " +
			// 	"Specifying this field will enforce exact membership matching.",
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"id": {
			// 				Type:     schema.TypeSet,
			// 				Required: true,
			// 				Elem:     &schema.Schema{Type: schema.TypeString},
			// 			},
			// 		},
			// 	},
			// },
			"service_edges": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true, // This is key
				MaxItems: 1,
				Description: "WARNING: Service edge membership is managed externally. " +
					"Specifying this field will enforce exact membership matching." +
					"This field will be deprecated in future versions",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// If the field isn't set in config (regardless of static or dynamic block)
					if _, ok := d.GetOk("service_edges"); !ok {
						return true
					}
					return false
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"trusted_networks": {
				Type:     schema.TypeList,
				Optional: true,
				// MaxItems:    1,
				Description: "List of trusted network IDs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"upgrade_day": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "SUNDAY",
				Description: "Service Edges in this group will attempt to update to a newer version of the software during this specified day.",
			},
			"upgrade_time_in_secs": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "66600",
				Description: "Service Edges in this group will attempt to update to a newer version of the software during this specified time.",
			},
			"version_profile_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the version profile.",
			},
			"version_profile_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the version profile. To learn more, see Version Profile Use Cases. This value is required, if the value for overrideVersionProfile is set to true",
				ValidateFunc: validation.StringInSlice([]string{
					"Default", "Previous Default",
					"New Release", "Default - el8",
					"New Release - el8", "Previous Default - el8",
				}, false),
			},
			"version_profile_visibility_scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the version profile.",
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"grace_distance_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If enabled, allows ZPA Private Service Edge Groups within the specified distance to be prioritized over a closer ZPA Public Service Edge.",
			},
			"grace_distance_value": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"grace_distance_enabled", "grace_distance_value_unit"},
				Description:  "Indicates the maximum distance in miles or kilometers to ZPA Private Service Edge groups that would override a ZPA Public Service Edge",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Normalize the old and new values by converting them to float64 and formatting them as strings with one decimal place
					oldVal, errOld := strconv.ParseFloat(old, 64)
					if errOld != nil {
						return false // If the old value can't be parsed as float, don't suppress the diff
					}
					newVal, errNew := strconv.ParseFloat(new, 64)
					if errNew != nil {
						return false // If the new value can't be parsed as float, don't suppress the diff
					}
					// Return true if the normalized old and new values are equal
					return fmt.Sprintf("%.1f", oldVal) == fmt.Sprintf("%.1f", newVal)
				},
			},

			"grace_distance_value_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"grace_distance_enabled", "grace_distance_value"},
				Description:  "Indicates the grace distance unit of measure in miles or kilometers. This value is only required if grace_distance_value is set to true",
				ValidateFunc: validation.StringInSlice([]string{
					"MILES", "KMS",
				}, false),
			},
		},
	}
}

func resourceServiceEdgeGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Ensure version_profile_id is set if version_profile_name is provided
	if err := validateAndSetProfileNameID(ctx, d, service); err != nil {
		return diag.FromErr(err)
	}

	req := expandServiceEdgeGroup(d)
	log.Printf("[INFO] Creating zpa service edge group with request\n%+v\n", req)

	resp, _, err := serviceedgegroup.Create(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created service edge group request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceServiceEdgeGroupRead(ctx, d, meta)
}

func resourceServiceEdgeGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := serviceedgegroup.Get(ctx, service, d.Id())
	if err != nil {
		if err.(*errorx.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing service edge group %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting service edge group:\n%+v\n", resp)
	d.SetId(resp.ID)
	isPublic, _ := strconv.ParseBool(resp.IsPublic)
	_ = d.Set("name", resp.Name)
	_ = d.Set("city_country", resp.CityCountry)
	_ = d.Set("country_code", resp.CountryCode)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("is_public", isPublic)
	_ = d.Set("latitude", resp.Latitude)
	_ = d.Set("longitude", resp.Longitude)
	_ = d.Set("location", resp.Location)
	_ = d.Set("upgrade_day", resp.UpgradeDay)
	_ = d.Set("use_in_dr_mode", resp.UseInDrMode)
	_ = d.Set("upgrade_time_in_secs", resp.UpgradeTimeInSecs)
	_ = d.Set("override_version_profile", resp.OverrideVersionProfile)
	_ = d.Set("version_profile_id", resp.VersionProfileID)
	_ = d.Set("version_profile_name", resp.VersionProfileName)
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	_ = d.Set("version_profile_visibility_scope", resp.VersionProfileVisibilityScope)
	log.Printf("[DEBUG] Set grace_distance_enabled to: %v", resp.GraceDistanceEnabled)
	_ = d.Set("grace_distance_enabled", resp.GraceDistanceEnabled)
	if resp.GraceDistanceEnabled {
		_ = d.Set("grace_distance_value", resp.GraceDistanceValue)
		_ = d.Set("grace_distance_value_unit", resp.GraceDistanceValueUnit)
	}
	_ = d.Set("trusted_networks", flattenAppTrustedNetworksSimple(resp.TrustedNetworks))
	// Only set service_edges if they were previously configured or exist in the diff
	if _, ok := d.GetOk("service_edges"); ok {
		_ = d.Set("service_edges", flattenServiceEdgeSimple(resp.ServiceEdges))
	} else {
		// Explicitly remove from state if not configured
		_ = d.Set("service_edges", nil)
	}

	return nil
}

func resourceServiceEdgeGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Ensure version_profile_id is set if version_profile_name is provided
	if err := validateAndSetProfileNameID(ctx, d, service); err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	log.Printf("[INFO] Updating app connector group ID: %v\n", id)
	req := expandServiceEdgeGroup(d)

	if _, _, err := serviceedgegroup.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := serviceedgegroup.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceServiceEdgeGroupRead(ctx, d, meta)
}

func resourceServiceEdgeGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}
	log.Printf("[INFO] Deleting service edge group ID: %v\n", d.Id())

	if _, err := serviceedgegroup.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] service edge group deleted")
	return nil
}

func expandServiceEdgeGroup(d *schema.ResourceData) serviceedgegroup.ServiceEdgeGroup {
	serviceEdgeGroup := serviceedgegroup.ServiceEdgeGroup{
		ID:                            d.Get("id").(string),
		Name:                          d.Get("name").(string),
		CityCountry:                   d.Get("city_country").(string),
		CountryCode:                   d.Get("country_code").(string),
		Description:                   d.Get("description").(string),
		Enabled:                       d.Get("enabled").(bool),
		IsPublic:                      strings.ToUpper(strconv.FormatBool(d.Get("is_public").(bool))),
		Latitude:                      d.Get("latitude").(string),
		Location:                      d.Get("location").(string),
		Longitude:                     d.Get("longitude").(string),
		UpgradeDay:                    d.Get("upgrade_day").(string),
		UseInDrMode:                   d.Get("use_in_dr_mode").(bool),
		UpgradeTimeInSecs:             d.Get("upgrade_time_in_secs").(string),
		VersionProfileID:              d.Get("version_profile_id").(string),
		VersionProfileName:            d.Get("version_profile_name").(string),
		VersionProfileVisibilityScope: d.Get("version_profile_visibility_scope").(string),
		OverrideVersionProfile:        d.Get("override_version_profile").(bool),
		MicroTenantID:                 d.Get("microtenant_id").(string),
		GraceDistanceEnabled:          d.Get("grace_distance_enabled").(bool),
		GraceDistanceValue:            d.Get("grace_distance_value").(string),
		GraceDistanceValueUnit:        d.Get("grace_distance_value_unit").(string),
		ServiceEdges:                  expandServiceEdges(d),
		TrustedNetworks:               expandTrustedNetworks(d),
	}

	return serviceEdgeGroup
}

/*
func expandServiceEdges(d *schema.ResourceData) []serviceedgecontroller.ServiceEdgeController {
	raw, ok := d.GetOk("service_edges")
	if !ok || raw == nil {
		return nil
	}

	blocks := raw.([]interface{})
	if len(blocks) == 0 {
		return nil
	}

	block, ok := blocks[0].(map[string]interface{})
	if !ok {
		return nil
	}

	idRaw, ok := block["id"]
	if !ok || idRaw == nil {
		return nil
	}

	idSet, ok := idRaw.(*schema.Set)
	if !ok {
		return nil
	}

	var edges []serviceedgecontroller.ServiceEdgeController
	for _, id := range idSet.List() {
		edges = append(edges, serviceedgecontroller.ServiceEdgeController{
			ID: id.(string),
		})
	}
	return edges
}
*/

func expandServiceEdges(d *schema.ResourceData) []serviceedgecontroller.ServiceEdgeController {
	// Return nil if service_edges isn't in config
	if _, ok := d.GetOk("service_edges"); !ok {
		return nil
	}

	raw := d.Get("service_edges")
	if raw == nil {
		return nil
	}

	blocks := raw.([]interface{})
	if len(blocks) == 0 || blocks[0] == nil {
		return nil
	}

	block, ok := blocks[0].(map[string]interface{})
	if !ok {
		return nil
	}

	idRaw, ok := block["id"]
	if !ok || idRaw == nil {
		return nil
	}

	idSet, ok := idRaw.(*schema.Set)
	if !ok || idSet.Len() == 0 {
		return nil
	}

	var edges []serviceedgecontroller.ServiceEdgeController
	for _, id := range idSet.List() {
		edges = append(edges, serviceedgecontroller.ServiceEdgeController{
			ID: id.(string),
		})
	}
	return edges
}

func expandTrustedNetworks(d *schema.ResourceData) []trustednetwork.TrustedNetwork {
	raw, ok := d.GetOk("trusted_networks")
	if !ok || raw == nil {
		return nil
	}

	blocks := raw.([]interface{})
	if len(blocks) == 0 {
		return nil
	}

	block, ok := blocks[0].(map[string]interface{})
	if !ok {
		return nil
	}

	idRaw, ok := block["id"]
	if !ok || idRaw == nil {
		return nil
	}

	idSet, ok := idRaw.(*schema.Set)
	if !ok {
		return nil
	}

	var networks []trustednetwork.TrustedNetwork
	for _, id := range idSet.List() {
		networks = append(networks, trustednetwork.TrustedNetwork{
			ID: id.(string),
		})
	}
	return networks
}

// func flattenServiceEdgeSimple(serviceEdges []serviceedgecontroller.ServiceEdgeController) []interface{} {
// 	ids := make([]interface{}, len(serviceEdges))
// 	for i, edge := range serviceEdges {
// 		ids[i] = edge.ID
// 	}

// 	return []interface{}{
// 		map[string]interface{}{
// 			"id": schema.NewSet(schema.HashString, ids),
// 		},
// 	}
// }

// func flattenServiceEdgeSimple(serviceEdges []serviceedgecontroller.ServiceEdgeController) []interface{} {
// 	if len(serviceEdges) == 0 {
// 		return nil
// 	}

// 	ids := make([]interface{}, len(serviceEdges))
// 	for i, edge := range serviceEdges {
// 		ids[i] = edge.ID
// 	}

// 	return []interface{}{
// 		map[string]interface{}{
// 			"id": schema.NewSet(schema.HashString, ids),
// 		},
// 	}
// }

func flattenServiceEdgeSimple(serviceEdges []serviceedgecontroller.ServiceEdgeController) []interface{} {
	if len(serviceEdges) == 0 {
		return nil
	}

	ids := make([]interface{}, len(serviceEdges))
	for i, edge := range serviceEdges {
		ids[i] = edge.ID
	}

	return []interface{}{
		map[string]interface{}{
			"id": schema.NewSet(schema.HashString, ids),
		},
	}
}

// func flattenAppTrustedNetworksSimple(trustedNetworks []trustednetwork.TrustedNetwork) []interface{} {
// 	ids := make([]interface{}, len(trustedNetworks))
// 	for i, network := range trustedNetworks {
// 		ids[i] = network.ID
// 	}

// 	return []interface{}{
// 		map[string]interface{}{
// 			"id": schema.NewSet(schema.HashString, ids),
// 		},
// 	}
// }

func flattenAppTrustedNetworksSimple(trustedNetworks []trustednetwork.TrustedNetwork) []interface{} {
	if len(trustedNetworks) == 0 {
		return nil // ✅ don't emit anything if no IDs
	}

	ids := make([]interface{}, len(trustedNetworks))
	for i, network := range trustedNetworks {
		ids[i] = network.ID
	}

	return []interface{}{
		map[string]interface{}{
			"id": schema.NewSet(schema.HashString, ids),
		},
	}
}
