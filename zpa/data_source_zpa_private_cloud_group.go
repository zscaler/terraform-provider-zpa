package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_group"
)

func dataSourcePrivateCloudGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateCloudGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"city_country": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"country_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"geo_location_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latitude": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"longitude": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"override_version_profile": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"restriction_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microtenant_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"site_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"upgrade_day": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"upgrade_time_in_secs": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_profile_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zscaler_managed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourcePrivateCloudGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var resp *private_cloud_group.PrivateCloudGroup
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for private cloud group %s\n", id)
		res, _, err := private_cloud_group.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for private cloud group name %s\n", name)
		res, _, err := private_cloud_group.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("city_country", resp.CityCountry)
		_ = d.Set("country_code", resp.CountryCode)
		_ = d.Set("description", resp.Description)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("geo_location_id", resp.GeoLocationID)
		_ = d.Set("is_public", resp.IsPublic)
		_ = d.Set("latitude", resp.Latitude)
		_ = d.Set("location", resp.Location)
		_ = d.Set("longitude", resp.Longitude)
		_ = d.Set("name", resp.Name)
		_ = d.Set("override_version_profile", resp.OverrideVersionProfile)
		_ = d.Set("read_only", resp.ReadOnly)
		_ = d.Set("restriction_type", resp.RestrictionType)
		_ = d.Set("microtenant_id", resp.MicrotenantID)
		_ = d.Set("microtenant_name", resp.MicrotenantName)
		_ = d.Set("site_id", resp.SiteID)
		_ = d.Set("site_name", resp.SiteName)
		_ = d.Set("upgrade_day", resp.UpgradeDay)
		_ = d.Set("upgrade_time_in_secs", resp.UpgradeTimeInSecs)
		_ = d.Set("version_profile_id", resp.VersionProfileID)
		_ = d.Set("zscaler_managed", resp.ZscalerManaged)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any private cloud group with name '%s' or id '%s'", name, id))
	}

	return nil
}
