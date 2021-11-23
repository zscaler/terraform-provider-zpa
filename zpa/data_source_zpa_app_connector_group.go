package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/appconnectorgroup"
)

func appConnectorGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connectors": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"application_start_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"appconnector_group_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"appconnector_group_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"control_channel_status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"creation_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"ctrl_broker_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"current_version": {
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
					"expected_upgrade_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"expected_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"fingerprint": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"ipacl": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"issued_cert_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"last_broker_connect_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"last_broker_connect_time_duration": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"last_broker_disconnect_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"last_broker_disconnect_time_duration": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"last_upgrade_time": {
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
					"modifiedby": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"modified_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"provisioning_key_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"provisioning_key_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"platform": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"previous_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"public_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"sarge_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"enrollment_cert": {
						Type:     schema.TypeMap,
						Elem:     schema.TypeString,
						Computed: true,
					},
					"upgrade_attempt": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"upgrade_status": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"city_country": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"country_code": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"creation_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dns_query_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"geo_location_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"geolocation_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"id": {
			Type:     schema.TypeString,
			Optional: true,
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
		"modifiedby": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"modified_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"override_version_profile": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"server_groups": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"config_space": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"creation_time": {
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
					"id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"dynamic_discovery": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"modifiedby": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"modified_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"lss_app_connector_group": {
			Type:     schema.TypeBool,
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
		"version_profile_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"version_profile_visibility_scope": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceAppConnectorGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConnectorGroupRead,
		Schema: MergeSchema(appConnectorGroupSchema(),
			map[string]*schema.Schema{
				"list": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: appConnectorGroupSchema(),
					},
				},
			}),
	}
}

func dataSourceConnectorGroupRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *appconnectorgroup.AppConnectorGroup
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for app connector group  %s\n", id)
		res, _, err := zClient.appconnectorgroup.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for app connector group name %s\n", name)
		res, _, err := zClient.appconnectorgroup.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("id", resp.ID)
		_ = d.Set("city_country", resp.CityCountry)
		_ = d.Set("country_code", resp.CountryCode)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("description", resp.Description)
		_ = d.Set("dns_query_type", resp.DNSQueryType)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("geo_location_id", resp.GeoLocationID)
		_ = d.Set("latitude", resp.Latitude)
		_ = d.Set("location", resp.Location)
		_ = d.Set("longitude", resp.Longitude)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("override_version_profile", resp.OverrideVersionProfile)
		_ = d.Set("lss_app_connector_group", resp.LSSAppConnectorGroup)
		_ = d.Set("upgrade_day", resp.UpgradeDay)
		_ = d.Set("upgrade_time_in_secs", resp.UpgradeTimeInSecs)
		_ = d.Set("version_profile_id", resp.VersionProfileID)
		_ = d.Set("version_profile_name", resp.VersionProfileName)
		_ = d.Set("connectors", flattenConnectors(resp))
		_ = d.Set("list", flattenAppConnectorGroupList([]appconnectorgroup.AppConnectorGroup{*resp}))
		if err := d.Set("server_groups", flattenServerGroups(resp)); err != nil {
			return fmt.Errorf("failed to read server groups %s", err)
		}
	} else if id != "" || name != "" {
		return fmt.Errorf("couldn't find any app connector group with name '%s' or id '%s'", name, id)
	} else {
		// get all
		list, _, err := zClient.appconnectorgroup.GetAll()
		if err != nil {
			return err
		}
		d.SetId("app-connector-group-list")
		_ = d.Set("list", flattenAppConnectorGroupList(list))
	}

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

func flattenConnectors(appConnector *appconnectorgroup.AppConnectorGroup) []interface{} {
	appConnectors := make([]interface{}, len(appConnector.Connectors))
	for i, appConnector := range appConnector.Connectors {
		appConnectors[i] = map[string]interface{}{
			"application_start_time":               appConnector.ApplicationStartTime,
			"appconnector_group_id":                appConnector.AppConnectorGroupID,
			"appconnector_group_name":              appConnector.AppConnectorGroupName,
			"control_channel_status":               appConnector.ControlChannelStatus,
			"creation_time":                        appConnector.CreationTime,
			"ctrl_broker_name":                     appConnector.CtrlBrokerName,
			"current_version":                      appConnector.CurrentVersion,
			"description":                          appConnector.Description,
			"enabled":                              appConnector.Enabled,
			"expected_upgrade_time":                appConnector.ExpectedUpgradeTime,
			"expected_version":                     appConnector.ExpectedVersion,
			"fingerprint":                          appConnector.Fingerprint,
			"id":                                   appConnector.ID,
			"ipacl":                                appConnector.IPACL,
			"issued_cert_id":                       appConnector.IssuedCertID,
			"last_broker_connect_time":             appConnector.LastBrokerConnectTime,
			"last_broker_connect_time_duration":    appConnector.LastBrokerConnectTimeDuration,
			"last_broker_disconnect_time":          appConnector.LastBrokerDisconnectTime,
			"last_broker_disconnect_time_duration": appConnector.LastBrokerDisconnectTimeDuration,
			"last_upgrade_time":                    appConnector.LastUpgradeTime,
			"latitude":                             appConnector.Latitude,
			"location":                             appConnector.Location,
			"longitude":                            appConnector.Longitude,
			"modifiedby":                           appConnector.ModifiedBy,
			"modified_time":                        appConnector.ModifiedTime,
			"name":                                 appConnector.Name,
			"provisioning_key_id":                  appConnector.ProvisioningKeyID,
			"provisioning_key_name":                appConnector.ProvisioningKeyName,
			"platform":                             appConnector.Platform,
			"previous_version":                     appConnector.PreviousVersion,
			"private_ip":                           appConnector.PrivateIP,
			"public_ip":                            appConnector.PublicIP,
			"sarge_version":                        appConnector.SargeVersion,
			"enrollment_cert":                      appConnector.EnrollmentCert,
			"upgrade_attempt":                      appConnector.UpgradeAttempt,
			"upgrade_status":                       appConnector.UpgradeStatus,
		}
	}

	return appConnectors
}

func flattenServerGroups(serverGroup *appconnectorgroup.AppConnectorGroup) []interface{} {
	serverGroups := make([]interface{}, len(serverGroup.AppServerGroup))
	for i, serverGroup := range serverGroup.AppServerGroup {
		serverGroups[i] = map[string]interface{}{
			"config_space":      serverGroup.ConfigSpace,
			"creation_time":     serverGroup.CreationTime,
			"description":       serverGroup.Description,
			"enabled":           serverGroup.Enabled,
			"id":                serverGroup.ID,
			"dynamic_discovery": serverGroup.DynamicDiscovery,
			"modifiedby":        serverGroup.ModifiedBy,
			"modified_time":     serverGroup.ModifiedTime,
			"name":              serverGroup.Name,
		}
	}

	return serverGroups
}
