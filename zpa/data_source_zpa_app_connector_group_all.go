package zpa

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/appconnectorgroup"
)

func dataSourceAppConnectorGroupAll() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConnectorGroupAllRead,
		Schema: map[string]*schema.Schema{
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: appConnectorGroupSchema(),
				},
			},
		},
	}
}
func dataSourceConnectorGroupAllRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	list, _, err := zClient.appconnectorgroup.GetAll()
	if err != nil {
		return err
	}
	d.SetId("app-connector-group-list")
	_ = d.Set("list", flattenAppConnectorGroupList(list))

	return nil
}

func flattenAppConnectorGroupList(list []appconnectorgroup.AppConnectorGroup) []interface{} {
	appConnectorGroup := make([]interface{}, len(list))
	for i, item := range list {
		appConnectorGroup[i] = map[string]interface{}{
			"id":                       item.ID,
			"city_country":             item.CityCountry,
			"country_code":             item.CountryCode,
			"creation_time":            item.CreationTime,
			"description":              item.Description,
			"dns_query_type":           item.DNSQueryType,
			"enabled":                  item.Enabled,
			"geo_location_id":          item.GeoLocationID,
			"latitude":                 item.Latitude,
			"location":                 item.Location,
			"longitude":                item.Longitude,
			"modifiedby":               item.ModifiedBy,
			"modified_time":            item.ModifiedTime,
			"name":                     item.Name,
			"override_version_profile": item.OverrideVersionProfile,
			"lss_app_connector_group":  item.LSSAppConnectorGroup,
			"upgrade_day":              item.UpgradeDay,
			"upgrade_time_in_secs":     item.UpgradeTimeInSecs,
			"version_profile_id":       item.VersionProfileID,
			"version_profile_name":     item.VersionProfileName,
			"connectors":               flattenConnectors(&item),
			"server_groups":            flattenServerGroups(&item),
		}
	}
	return appConnectorGroup
}
