package zpa

import (
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/serviceedgegroup"
)

func resourceServiceEdgeGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceEdgeGroupCreate,
		Read:   resourceServiceEdgeGroupRead,
		Update: resourceServiceEdgeGroupUpdate,
		Delete: resourceServiceEdgeGroupDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := zClient.serviceedgegroup.GetByName(id)
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Latitude for the Service Edge Group.",
				ValidateFunc: ValidateStringFloatBetween(-90, 90),
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Location for the Service Edge Group.",
			},
			"longitude": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Longitude for the Service Edge Group.",
				ValidateFunc: ValidateStringFloatBetween(-180.0, 180.0),
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"trusted_networks": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeList,
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
				Computed:    true,
				Description: "ID of the version profile.",
			},
		},
	}
}

func resourceServiceEdgeGroupCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	if err := validateAndSetProfileNameID(d); err != nil {
		return err
	}
	req := expandServiceEdgeGroup(d)
	log.Printf("[INFO] Creating zpa service edge group with request\n%+v\n", req)

	resp, _, err := zClient.serviceedgegroup.Create(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created service edge group request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceServiceEdgeGroupRead(d, m)
}

func resourceServiceEdgeGroupRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.serviceedgegroup.Get(d.Id())
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
	_ = d.Set("version_profile_visibility_scope", resp.VersionProfileVisibilityScope)
	_ = d.Set("trusted_networks", flattenTrustedNetworks(resp))
	_ = d.Set("service_edges", flattenServiceEdges(resp))
	return nil

}

func resourceServiceEdgeGroupUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	if err := validateAndSetProfileNameID(d); err != nil {
		return err
	}
	id := d.Id()
	log.Printf("[INFO] Updating service edge group ID: %v\n", id)
	req := expandServiceEdgeGroup(d)

	if _, err := zClient.serviceedgegroup.Update(id, &req); err != nil {
		return err
	}

	return resourceServiceEdgeGroupRead(d, m)
}

func resourceServiceEdgeGroupDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting service edge group ID: %v\n", d.Id())

	if _, err := zClient.serviceedgegroup.Delete(d.Id()); err != nil {
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
		network := trustedNetworksInterface.(*schema.Set)
		log.Printf("[INFO] trusted network data: %+v\n", network)
		var trustedNetworks []serviceedgegroup.TrustedNetworks
		for _, trustedNetwork := range network.List() {
			trustedNetwork, ok := trustedNetwork.(map[string]interface{})
			if ok {
				for _, id := range trustedNetwork["id"].([]interface{}) {
					trustedNetworks = append(trustedNetworks, serviceedgegroup.TrustedNetworks{
						ID: id.(string),
					})
				}
			}
		}
		return trustedNetworks
	}

	return []serviceedgegroup.TrustedNetworks{}
}
