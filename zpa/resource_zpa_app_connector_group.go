package zpa

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/appconnectorgroup"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/client"
)

func resourceAppConnectorGroup() *schema.Resource {
	return &schema.Resource{
		Create:   resourceAppConnectorGroupCreate,
		Read:     resourceAppConnectorGroupRead,
		Update:   resourceAppConnectorGroupUpdate,
		Delete:   resourceAppConnectorGroupDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"city_country": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"country_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_query_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"latitude": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"longitude": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lss_app_connector_group": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"upgrade_day": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"upgrade_time_in_secs": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"override_version_profile": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"version_profile_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_profile_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAppConnectorGroupCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandAppConnectorGroup(d)
	log.Printf("[INFO] Creating zpa app connector group with request\n%+v\n", req)

	resp, _, err := zClient.appconnectorgroup.Create(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created app connector group request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceAppConnectorGroupRead(d, m)
}

func resourceAppConnectorGroupRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.appconnectorgroup.Get(d.Id())
	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing app connector group %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting application server:\n%+v\n", resp)
	_ = d.Set("name", resp.Name)
	_ = d.Set("city_country", resp.CityCountry)
	_ = d.Set("country_code", resp.CountryCode)
	_ = d.Set("description", resp.Description)
	_ = d.Set("dns_query_type", resp.DNSQueryType)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("latitude", resp.Latitude)
	_ = d.Set("longitude", resp.Longitude)
	_ = d.Set("location", resp.Location)
	_ = d.Set("lss_app_connector_group", resp.LSSAppConnectorGroup)
	_ = d.Set("upgrade_day", resp.UpgradeDay)
	_ = d.Set("upgrade_time_in_secs", resp.UpgradeTimeInSecs)
	_ = d.Set("override_version_profile", resp.OverrideVersionProfile)
	_ = d.Set("version_profile_id", resp.VersionProfileID)
	_ = d.Set("version_profile_name", resp.VersionProfileName)
	return nil

}

func resourceAppConnectorGroupUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating app connector group ID: %v\n", id)
	req := expandAppConnectorGroup(d)

	if _, err := zClient.appconnectorgroup.Update(id, &req); err != nil {
		return err
	}

	return resourceAppConnectorGroupRead(d, m)
}

func resourceAppConnectorGroupDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting app connector groupID: %v\n", d.Id())

	if _, err := zClient.appconnectorgroup.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] app connector group deleted")
	return nil
}

func expandAppConnectorGroup(d *schema.ResourceData) appconnectorgroup.AppConnectorGroup {
	appConnectorGroup := appconnectorgroup.AppConnectorGroup{
		ID:                     d.Get("id").(string),
		Name:                   d.Get("name").(string),
		CityCountry:            d.Get("city_country").(string),
		CountryCode:            d.Get("country_code").(string),
		Description:            d.Get("description").(string),
		DNSQueryType:           d.Get("dns_query_type").(string),
		Enabled:                d.Get("enabled").(bool),
		Latitude:               d.Get("latitude").(string),
		Longitude:              d.Get("longitude").(string),
		Location:               d.Get("location").(string),
		LSSAppConnectorGroup:   d.Get("lss_app_connector_group").(bool),
		UpgradeDay:             d.Get("upgrade_day").(string),
		UpgradeTimeInSecs:      d.Get("upgrade_time_in_secs").(string),
		OverrideVersionProfile: d.Get("override_version_profile").(bool),
		VersionProfileID:       d.Get("version_profile_id").(string),
		VersionProfileName:     d.Get("version_profile_name").(string),
	}
	return appConnectorGroup
}
