package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorgroup"
)

func dataSourceAppConnectorGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConnectorGroupRead,
		Schema: map[string]*schema.Schema{
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
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"geo_location_id": {
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
			"tcp_quick_ack_app": {
				Description: "Whether TCP Quick Acknowledgement is enabled or disabled for the application. The tcpQuickAckApp, tcpQuickAckAssistant, and tcpQuickAckReadAssistant fields must all share the same value.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"tcp_quick_ack_assistant": {
				Description: "Whether TCP Quick Acknowledgement is enabled or disabled for the application. The tcpQuickAckApp, tcpQuickAckAssistant, and tcpQuickAckReadAssistant fields must all share the same value.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"tcp_quick_ack_read_assistant": {
				Description: "Whether TCP Quick Acknowledgement is enabled or disabled for the application. The tcpQuickAckApp, tcpQuickAckAssistant, and tcpQuickAckReadAssistant fields must all share the same value.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"use_in_dr_mode": {
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
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microtenant_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceConnectorGroupRead(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).appconnectorgroup.WithMicroTenant(GetString(d.Get("microtenant_id")))

	var resp *appconnectorgroup.AppConnectorGroup
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for app connector group  %s\n", id)
		res, _, err := service.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for app connector group name %s\n", name)
		res, _, err := service.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
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
		_ = d.Set("tcp_quick_ack_app", resp.TCPQuickAckApp)
		_ = d.Set("tcp_quick_ack_assistant", resp.TCPQuickAckAssistant)
		_ = d.Set("tcp_quick_ack_read_assistant", resp.TCPQuickAckReadAssistant)
		_ = d.Set("use_in_dr_mode", resp.UseInDrMode)
		_ = d.Set("upgrade_day", resp.UpgradeDay)
		_ = d.Set("upgrade_time_in_secs", resp.UpgradeTimeInSecs)
		_ = d.Set("version_profile_id", resp.VersionProfileID)
		_ = d.Set("version_profile_name", resp.VersionProfileName)
		_ = d.Set("connectors", flattenConnectors(resp.Connectors))
		_ = d.Set("microtenant_id", resp.MicroTenantID)
		_ = d.Set("microtenant_name", resp.MicroTenantName)

		if err := d.Set("server_groups", flattenServerGroups(resp)); err != nil {
			return fmt.Errorf("failed to read server groups %s", err)
		}
	} else {
		return fmt.Errorf("couldn't find any app connector group with name '%s' or id '%s'", name, id)
	}

	return nil
}

func flattenConnectors(appConnector []appconnectorgroup.Connector) []interface{} {
	appConnectors := make([]interface{}, len(appConnector))
	for i, appConnectorItem := range appConnector {
		appConnectors[i] = map[string]interface{}{
			"application_start_time":               appConnectorItem.ApplicationStartTime,
			"appconnector_group_id":                appConnectorItem.AppConnectorGroupID,
			"appconnector_group_name":              appConnectorItem.AppConnectorGroupName,
			"control_channel_status":               appConnectorItem.ControlChannelStatus,
			"creation_time":                        appConnectorItem.CreationTime,
			"ctrl_broker_name":                     appConnectorItem.CtrlBrokerName,
			"current_version":                      appConnectorItem.CurrentVersion,
			"description":                          appConnectorItem.Description,
			"enabled":                              appConnectorItem.Enabled,
			"expected_upgrade_time":                appConnectorItem.ExpectedUpgradeTime,
			"expected_version":                     appConnectorItem.ExpectedVersion,
			"fingerprint":                          appConnectorItem.Fingerprint,
			"id":                                   appConnectorItem.ID,
			"ipacl":                                appConnectorItem.IPACL,
			"issued_cert_id":                       appConnectorItem.IssuedCertID,
			"last_broker_connect_time":             appConnectorItem.LastBrokerConnectTime,
			"last_broker_connect_time_duration":    appConnectorItem.LastBrokerConnectTimeDuration,
			"last_broker_disconnect_time":          appConnectorItem.LastBrokerDisconnectTime,
			"last_broker_disconnect_time_duration": appConnectorItem.LastBrokerDisconnectTimeDuration,
			"last_upgrade_time":                    appConnectorItem.LastUpgradeTime,
			"latitude":                             appConnectorItem.Latitude,
			"location":                             appConnectorItem.Location,
			"longitude":                            appConnectorItem.Longitude,
			"modifiedby":                           appConnectorItem.ModifiedBy,
			"modified_time":                        appConnectorItem.ModifiedTime,
			"name":                                 appConnectorItem.Name,
			"provisioning_key_id":                  appConnectorItem.ProvisioningKeyID,
			"provisioning_key_name":                appConnectorItem.ProvisioningKeyName,
			"platform":                             appConnectorItem.Platform,
			"previous_version":                     appConnectorItem.PreviousVersion,
			"private_ip":                           appConnectorItem.PrivateIP,
			"public_ip":                            appConnectorItem.PublicIP,
			"sarge_version":                        appConnectorItem.SargeVersion,
			"enrollment_cert":                      appConnectorItem.EnrollmentCert,
			"upgrade_attempt":                      appConnectorItem.UpgradeAttempt,
			"upgrade_status":                       appConnectorItem.UpgradeStatus,
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
