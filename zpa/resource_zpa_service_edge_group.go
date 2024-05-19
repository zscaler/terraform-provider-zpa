package zpa

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/serviceedgegroup"
)

func resourceServiceEdgeGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceEdgeGroupCreate,
		Read:   resourceServiceEdgeGroupRead,
		Update: resourceServiceEdgeGroupUpdate,
		Delete: resourceServiceEdgeGroupDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				service := m.(*Client).serviceedgegroup.WithMicroTenant(GetString(d.Get("microtenant_id")))

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := service.GetByName(id)
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
			"service_edges": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"trusted_networks": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
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
				ValidateFunc: validation.StringInSlice([]string{
					"0", "1", "2",
				}, false),
			},
			"version_profile_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the version profile.",
				ValidateFunc: validation.StringInSlice([]string{
					"Default", "Previous Default", "New Release",
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

func resourceServiceEdgeGroupCreate(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).serviceedgegroup.WithMicroTenant(GetString(d.Get("microtenant_id")))

	if err := validateAndSetProfileNameID(d); err != nil {
		return err
	}
	req := expandServiceEdgeGroup(d)
	log.Printf("[INFO] Creating zpa service edge group with request\n%+v\n", req)

	resp, _, err := service.Create(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created service edge group request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceServiceEdgeGroupRead(d, m)
}

func resourceServiceEdgeGroupRead(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).serviceedgegroup.WithMicroTenant(GetString(d.Get("microtenant_id")))

	resp, _, err := service.Get(d.Id())
	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing service edge group %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
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
	_ = d.Set("service_edges", flattenServiceEdgeSimple(resp.ServiceEdges))
	return nil
}

func resourceServiceEdgeGroupUpdate(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).serviceedgegroup.WithMicroTenant(GetString(d.Get("microtenant_id")))

	if err := validateAndSetProfileNameID(d); err != nil {
		return err
	}
	id := d.Id()
	log.Printf("[INFO] Updating service edge group ID: %v\n", id)
	req := expandServiceEdgeGroup(d)

	if _, _, err := service.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := service.Update(id, &req); err != nil {
		return err
	}

	return resourceServiceEdgeGroupRead(d, m)
}

func resourceServiceEdgeGroupDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.serviceedgegroup.WithMicroTenant(GetString(d.Get("microtenant_id")))

	log.Printf("[INFO] Deleting service edge group ID: %v\n", d.Id())

	if _, err := service.Delete(d.Id()); err != nil {
		return err
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

func expandServiceEdges(d *schema.ResourceData) []serviceedgegroup.ServiceEdges {
	serviceEdgesGroupInterface, ok := d.GetOk("service_edges")
	if ok {
		serviceEdge := serviceEdgesGroupInterface.(*schema.Set)
		log.Printf("[INFO] service edges data: %+v\n", serviceEdge)
		var serviceEdgesGroups []serviceedgegroup.ServiceEdges
		for _, serviceEdgesGroup := range serviceEdge.List() {
			serviceEdgesGroup, ok := serviceEdgesGroup.(map[string]interface{})
			if ok {
				for _, id := range serviceEdgesGroup["id"].([]interface{}) {
					serviceEdgesGroups = append(serviceEdgesGroups, serviceedgegroup.ServiceEdges{
						ID: id.(string),
					})
				}
			}
		}
		return serviceEdgesGroups
	}

	return []serviceedgegroup.ServiceEdges{}
}

func expandTrustedNetworks(d *schema.ResourceData) []serviceedgegroup.TrustedNetworks {
	trustedNetworksInterface, ok := d.GetOk("trusted_networks")
	if ok {
		trusteNetwork := trustedNetworksInterface.(*schema.Set)
		log.Printf("[INFO] trusted network data: %+v\n", trusteNetwork)
		var trusteNetworks []serviceedgegroup.TrustedNetworks
		for _, trusteNetwork := range trusteNetwork.List() {
			trusteNetwork, _ := trusteNetwork.(map[string]interface{})
			if trusteNetwork != nil {
				for _, id := range trusteNetwork["id"].([]interface{}) {
					trusteNetworks = append(trusteNetworks, serviceedgegroup.TrustedNetworks{
						ID: id.(string),
					})
				}
			}
		}
		return trusteNetworks
	}

	return []serviceedgegroup.TrustedNetworks{}
}

func flattenAppTrustedNetworksSimple(trustedNetworks []serviceedgegroup.TrustedNetworks) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(trustedNetworks))
	for i, networks := range trustedNetworks {
		ids[i] = networks.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func flattenServiceEdgeSimple(serviceEdges []serviceedgegroup.ServiceEdges) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(serviceEdges))
	for i, edges := range serviceEdges {
		ids[i] = edges.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}
