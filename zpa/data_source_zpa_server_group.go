package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/servergroup"
)

func dataSourceServerGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerGroupRead,
		Schema: map[string]*schema.Schema{
			"applications": {
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
					},
				},
			},
			"app_connector_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"connectors": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
									"fingerprint": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"issued_cert_id": {
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
									"upgrade_attempt": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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
										Optional: true,
									},
								},
							},
						},
						"siem_app_connector_group": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"upgrade_time_in_secs": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"upgrade_day": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_profile_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
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
				Optional: true,
			},
			"ip_anchored": {
				Type:     schema.TypeBool,
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
				Optional: true,
			},
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_server_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
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
		},
	}
}

func dataSourceServerGroupRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *servergroup.ServerGroup
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for server group  %s\n", id)
		res, _, err := zClient.servergroup.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for server group name %s\n", name)
		res, _, err := zClient.servergroup.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("config_space", resp.ConfigSpace)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("description", resp.Description)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("dynamic_discovery", resp.DynamicDiscovery)
		_ = d.Set("ip_anchored", resp.IpAnchored)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)

		if err := d.Set("applications", flattenServerGroupApplications(resp.Applications)); err != nil {
			return err
		}

		if err := d.Set("app_connector_groups", flattenAppConnectorGroups(resp.AppConnectorGroups)); err != nil {
			return err
		}

		if err := d.Set("servers", flattenServers(resp.Servers)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("couldn't find any server group with name '%s' or id '%s'", name, id)
	}

	return nil
}
func flattenServerGroupApplications(applications []servergroup.Applications) []interface{} {
	serverGroupApplications := make([]interface{}, len(applications))
	for i, srvApplication := range applications {
		serverGroupApplications[i] = map[string]interface{}{
			"id":   srvApplication.ID,
			"name": srvApplication.Name,
		}
	}

	return serverGroupApplications
}

func flattenAppConnectorGroups(appConnectorGroup []servergroup.AppConnectorGroups) []interface{} {
	appConnectorGroups := make([]interface{}, len(appConnectorGroup))
	for i, appConnectorGroup := range appConnectorGroup {
		appConnectorGroups[i] = map[string]interface{}{
			"city_country":             appConnectorGroup.Citycountry,
			"country_code":             appConnectorGroup.CountryCode,
			"creation_time":            appConnectorGroup.CreationTime,
			"description":              appConnectorGroup.Description,
			"dns_query_type":           appConnectorGroup.DnsqueryType,
			"enabled":                  appConnectorGroup.Enabled,
			"geolocation_id":           appConnectorGroup.GeolocationID,
			"id":                       appConnectorGroup.ID,
			"latitude":                 appConnectorGroup.Latitude,
			"location":                 appConnectorGroup.Location,
			"longitude":                appConnectorGroup.Longitude,
			"modifiedby":               appConnectorGroup.ModifiedBy,
			"modified_time":            appConnectorGroup.ModifiedTime,
			"name":                     appConnectorGroup.Name,
			"siem_app_connector_group": appConnectorGroup.SiemAppconnectorGroup,
			"upgrade_day":              appConnectorGroup.UpgradeDay,
			"upgrade_time_in_secs":     appConnectorGroup.UpgradeTimeinSecs,
			"version_profile_id":       appConnectorGroup.VersionProfileID,
			"server_groups":            flattenAppConnectorServerGroups(appConnectorGroup),
			"connectors":               flattenAppConnectors(appConnectorGroup),
		}
	}

	return appConnectorGroups
}

func flattenAppConnectorServerGroups(serverGroup servergroup.AppConnectorGroups) []interface{} {
	serverGroups := make([]interface{}, len(serverGroup.AppServerGroups))
	for i, serverGroup := range serverGroup.AppServerGroups {
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

func flattenAppConnectors(connector servergroup.AppConnectorGroups) []interface{} {
	appConnectors := make([]interface{}, len(connector.Connectors))
	for i, appConnector := range connector.Connectors {
		appConnectors[i] = map[string]interface{}{
			"creation_time": appConnector.CreationTime,
			"description":   appConnector.Description,
			"enabled":       appConnector.Enabled,
			"id":            appConnector.ID,
			"modifiedby":    appConnector.ModifiedBy,
			"modified_time": appConnector.ModifiedTime,
			"name":          appConnector.Name,
		}
	}

	return appConnectors
}

func flattenServers(applicationServer []servergroup.ApplicationServer) []interface{} {
	applicationServers := make([]interface{}, len(applicationServer))
	for i, appServerItem := range applicationServer {
		applicationServers[i] = map[string]interface{}{
			"address":              appServerItem.Address,
			"app_server_group_ids": appServerItem.AppServerGroupIds,
			"config_space":         appServerItem.ConfigSpace,
			"creation_time":        appServerItem.CreationTime,
			"description":          appServerItem.Description,
			"enabled":              appServerItem.Enabled,
			"id":                   appServerItem.ID,
			"modifiedby":           appServerItem.ModifiedBy,
			"modified_time":        appServerItem.ModifiedTime,
			"name":                 appServerItem.Name,
		}
	}
	return applicationServers
}
