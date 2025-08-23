package zpa

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_group"
)

func resourcePrivateCloudGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateCloudGroupCreate,
		ReadContext:   resourcePrivateCloudGroupRead,
		UpdateContext: resourcePrivateCloudGroupUpdate,
		DeleteContext: resourcePrivateCloudGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				client := meta.(*Client)
				service := client.Service

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
					resp, _, err := private_cloud_group.GetByName(ctx, service, id)
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
				Description: "Name of the Private Cloud Group",
			},
			"city_country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "City and country of the Private Cloud Group",
			},
			"country_code": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateCountryCode,
				Description:  "Country code of the Private Cloud Group",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Private Cloud Group",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether this Private Cloud Group is enabled or not",
			},
			"is_public": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Whether the Private Cloud Group is public",
			},
			"latitude": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     ValidateLatitude,
				DiffSuppressFunc: DiffSuppressFuncCoordinate,
				Description:      "Latitude of the Private Cloud Group. Integer or decimal. With values in the range of -90 to 90",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location of the Private Cloud Group",
			},
			"longitude": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     ValidateLongitude,
				DiffSuppressFunc: DiffSuppressFuncCoordinate,
				Description:      "Longitude of the Private Cloud Group. Integer or decimal. With values in the range of -180 to 180",
			},
			"override_version_profile": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the default version profile of the Private Cloud Group is applied or overridden",
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Microtenant ID for the Private Cloud Group",
			},
			"site_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Site ID for the Private Cloud Group",
			},
			"upgrade_day": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Private Cloud Controllers in this group will attempt to update to a newer version of the software during this specified day",
				ValidateFunc: validation.StringInSlice([]string{
					"SUNDAY", "MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY",
				}, false),
			},
			"upgrade_time_in_secs": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Private Cloud Controllers in this group will attempt to update to a newer version of the software during this specified time. Integer in seconds (i.e., -66600). The integer should be greater than or equal to 0 and less than 86400, in 15 minute intervals",
			},
			"version_profile_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the version profile for the Private Cloud Group",
			},
		},
	}
}

func resourcePrivateCloudGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandPrivateCloudGroup(d)
	log.Printf("[INFO] Creating zpa private cloud group with request\n%+v\n", req)

	resp, _, err := private_cloud_group.Create(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created private cloud group request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourcePrivateCloudGroupRead(ctx, d, meta)
}

func resourcePrivateCloudGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := private_cloud_group.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing private cloud group %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting private cloud group:\n%+v\n", resp)
	_ = d.Set("name", resp.Name)
	_ = d.Set("city_country", resp.CityCountry)
	_ = d.Set("country_code", resp.CountryCode)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("is_public", resp.IsPublic)
	_ = d.Set("latitude", resp.Latitude)
	_ = d.Set("location", resp.Location)
	_ = d.Set("longitude", resp.Longitude)
	_ = d.Set("override_version_profile", resp.OverrideVersionProfile)
	_ = d.Set("microtenant_id", resp.MicrotenantID)
	_ = d.Set("site_id", resp.SiteID)
	_ = d.Set("upgrade_day", resp.UpgradeDay)
	_ = d.Set("upgrade_time_in_secs", resp.UpgradeTimeInSecs)
	_ = d.Set("version_profile_id", resp.VersionProfileID)
	return nil
}

func resourcePrivateCloudGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating private cloud group ID: %v\n", id)
	req := expandPrivateCloudGroup(d)

	if _, _, err := private_cloud_group.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if _, err := private_cloud_group.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourcePrivateCloudGroupRead(ctx, d, meta)
}

func resourcePrivateCloudGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting private cloud group with id %v\n", d.Id())

	if _, err := private_cloud_group.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandPrivateCloudGroup(d *schema.ResourceData) private_cloud_group.PrivateCloudGroup {
	return private_cloud_group.PrivateCloudGroup{
		ID:                     d.Get("id").(string),
		Name:                   d.Get("name").(string),
		CityCountry:            d.Get("city_country").(string),
		CountryCode:            d.Get("country_code").(string),
		Description:            d.Get("description").(string),
		Enabled:                d.Get("enabled").(bool),
		IsPublic:               d.Get("is_public").(string),
		Latitude:               d.Get("latitude").(string),
		Location:               d.Get("location").(string),
		Longitude:              d.Get("longitude").(string),
		OverrideVersionProfile: d.Get("override_version_profile").(bool),
		MicrotenantID:          d.Get("microtenant_id").(string),
		SiteID:                 d.Get("site_id").(string),
		UpgradeDay:             d.Get("upgrade_day").(string),
		UpgradeTimeInSecs:      d.Get("upgrade_time_in_secs").(string),
		VersionProfileID:       d.Get("version_profile_id").(string),
	}
}
