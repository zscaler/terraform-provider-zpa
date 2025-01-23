package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgegroup"
)

func dataSourceServiceEdgeGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServiceEdgeGroupRead,
		Importer:    &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the Service Edge Group.",
			},
			"alt_cloud": {
				Type:     schema.TypeString,
				Computed: true,
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the Service Edge Group.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this Service Edge Group is enabled or not.",
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Latitude for the Service Edge Group.",
			},
			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Location for the Service Edge Group.",
			},
			"longitude": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Longitude for the Service Edge Group.",
			},
			"override_version_profile": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the default version profile of the App Connector Group is applied or overridden.",
			},
			"modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"use_in_dr_mode": {
				Type:     schema.TypeBool,
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service Edges in this group will attempt to update to a newer version of the software during this specified day.",
			},
			"upgrade_time_in_secs": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service Edges in this group will attempt to update to a newer version of the software during this specified time.",
			},
			"version_profile_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the version profile. To learn more",
			},
			"version_profile_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the version profile. To learn more",
			},
			"version_profile_visibility_scope": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the version profile. To learn more",
			},
			"grace_distance_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If enabled, allows ZPA Private Service Edge Groups within the specified distance to be prioritized over a closer ZPA Public Service Edge.",
			},
			"grace_distance_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the maximum distance in miles or kilometers to ZPA Private Service Edge groups that would override a ZPA Public Service Edge",
			},
			"grace_distance_value_unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the grace distance unit of measure in miles or kilometers. This value is only required if grace_distance_value is set to true",
			},
			"service_edges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_edge_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_edge_group_name": {
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
						"ip_acl": {
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
						"modified_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listen_ips": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
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
						"publish_ips": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"publish_ipv6": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"sarge_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime_os": {
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
						"microtenant_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"microtenant_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"private_broker_version": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"application_start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"broker_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"creation_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"current_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disable_auto_update": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"last_connect_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"last_disconnect_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"last_upgraded_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"lone_warrior": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"modified_by": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"modified_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"platform": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"platform_detail": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"previous_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service_edge_group_id": {
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
									"restart_instructions": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"restart_time_in_sec": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"runtime_os": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sarge_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"system_start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tunnel_id": {
										Type:     schema.TypeString,
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
									"upgrade_now_once": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"zpn_sub_module_upgrade": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"creation_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"current_version": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"entity_gid": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"entity_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"expected_version": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"modified_by": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"modified_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"previous_version": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"role": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"upgrade_status": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"upgrade_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"trusted_networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_customer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_by": {
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
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zscaler_cloud": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceServiceEdgeGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var resp *serviceedgegroup.ServiceEdgeGroup
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for service edge group %s\n", id)
		res, _, err := serviceedgegroup.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for service edge group name %s\n", name)
		res, _, err := serviceedgegroup.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("city_country", resp.CityCountry)
		_ = d.Set("country_code", resp.CountryCode)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("description", resp.Description)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("alt_cloud", resp.AltCloud)
		_ = d.Set("geo_location_id", resp.GeoLocationID)
		_ = d.Set("is_public", resp.IsPublic)
		_ = d.Set("latitude", resp.Latitude)
		_ = d.Set("location", resp.Location)
		_ = d.Set("longitude", resp.Longitude)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("upgrade_day", resp.UpgradeDay)
		_ = d.Set("use_in_dr_mode", resp.UseInDrMode)
		_ = d.Set("site_id", resp.SiteID)
		_ = d.Set("site_name", resp.SiteName)
		_ = d.Set("upgrade_time_in_secs", resp.UpgradeTimeInSecs)
		_ = d.Set("override_version_profile", resp.OverrideVersionProfile)
		_ = d.Set("version_profile_id", resp.VersionProfileID)
		_ = d.Set("version_profile_name", resp.VersionProfileName)
		_ = d.Set("version_profile_visibility_scope", resp.VersionProfileVisibilityScope)
		_ = d.Set("grace_distance_enabled", resp.GraceDistanceEnabled)
		_ = d.Set("grace_distance_value", resp.GraceDistanceValue)
		_ = d.Set("grace_distance_value_unit", resp.GraceDistanceValueUnit)
		_ = d.Set("trusted_networks", flattenTrustedNetworks(resp))
		_ = d.Set("service_edges", flattenServiceEdges(resp))

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any service edge group with name '%s' or id '%s'", name, id))
	}

	return nil
}

func flattenServiceEdges(serviceEdgeGroup *serviceedgegroup.ServiceEdgeGroup) []interface{} {
	// Access the ServiceEdges directly from the ServiceEdgeGroup struct
	serviceEdges := make([]interface{}, len(serviceEdgeGroup.ServiceEdges))
	for i, serviceEdge := range serviceEdgeGroup.ServiceEdges {
		serviceEdges[i] = map[string]interface{}{
			"application_start_time":               serviceEdge.ApplicationStartTime,
			"control_channel_status":               serviceEdge.ControlChannelStatus,
			"creation_time":                        serviceEdge.CreationTime,
			"ctrl_broker_name":                     serviceEdge.CtrlBrokerName,
			"current_version":                      serviceEdge.CurrentVersion,
			"description":                          serviceEdge.Description,
			"enabled":                              serviceEdge.Enabled,
			"expected_upgrade_time":                serviceEdge.ExpectedUpgradeTime,
			"expected_version":                     serviceEdge.ExpectedVersion,
			"fingerprint":                          serviceEdge.Fingerprint,
			"id":                                   serviceEdge.ID,
			"ip_acl":                               serviceEdge.IPACL,
			"issued_cert_id":                       serviceEdge.IssuedCertID,
			"last_broker_connect_time":             serviceEdge.LastBrokerConnectTime,
			"last_broker_connect_time_duration":    serviceEdge.LastBrokerConnectTimeDuration,
			"last_broker_disconnect_time":          serviceEdge.LastBrokerDisconnectTime,
			"last_broker_disconnect_time_duration": serviceEdge.LastBrokerDisconnectTimeDuration,
			"last_upgrade_time":                    serviceEdge.LastUpgradeTime,
			"latitude":                             serviceEdge.Latitude,
			"location":                             serviceEdge.Location,
			"longitude":                            serviceEdge.Longitude,
			"listen_ips":                           serviceEdge.ListenIPs,
			"modified_by":                          serviceEdge.ModifiedBy,
			"modified_time":                        serviceEdge.ModifiedTime,
			"name":                                 serviceEdge.Name,
			"provisioning_key_id":                  serviceEdge.ProvisioningKeyID,
			"provisioning_key_name":                serviceEdge.ProvisioningKeyName,
			"platform":                             serviceEdge.Platform,
			"previous_version":                     serviceEdge.PreviousVersion,
			"service_edge_group_id":                serviceEdge.ServiceEdgeGroupID,
			"service_edge_group_name":              serviceEdge.ServiceEdgeGroupName,
			"private_ip":                           serviceEdge.PrivateIP,
			"public_ip":                            serviceEdge.PublicIP,
			"publish_ips":                          serviceEdge.PublishIPs,
			"sarge_version":                        serviceEdge.SargeVersion,
			"enrollment_cert":                      serviceEdge.EnrollmentCert,
			"upgrade_attempt":                      serviceEdge.UpgradeAttempt,
			"upgrade_status":                       serviceEdge.UpgradeStatus,
			"private_broker_version":               flattenPrivateBrokerVersion(&serviceEdge.PrivateBrokerVersion),
		}
	}
	return serviceEdges
}

func flattenTrustedNetworks(trustedNetwork *serviceedgegroup.ServiceEdgeGroup) []interface{} {
	trustedNetworks := make([]interface{}, len(trustedNetwork.TrustedNetworks))
	for i, val := range trustedNetwork.TrustedNetworks {
		trustedNetworks[i] = map[string]interface{}{
			"creation_time":      val.CreationTime,
			"domain":             val.Domain,
			"id":                 val.ID,
			"master_customer_id": val.MasterCustomerID,
			"modified_by":        val.ModifiedBy,
			"modified_time":      val.ModifiedTime,
			"name":               val.Name,
			"network_id":         val.NetworkID,
			"zscaler_cloud":      val.ZscalerCloud,
		}
	}

	return trustedNetworks
}
