package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/serviceedgegroup"
)

func dataSourceServiceEdgeGroup() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceServiceEdgeGroupRead,
		Importer: &schema.ResourceImporter{},

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
			"service_edges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_start_time": {
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
						"listen_ips": {
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
						"service_edge_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_edge_group_name": {
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
			"trusted_networks": {
				Type:     schema.TypeSet,
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
		},
	}
}

func dataSourceServiceEdgeGroupRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *serviceedgegroup.ServiceEdgeGroup
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for service edge group %s\n", id)
		res, _, err := zClient.serviceedgegroup.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for service edge group name %s\n", name)
		res, _, err := zClient.serviceedgegroup.GetByName(name)
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
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("geo_location_id", resp.GeoLocationID)
		_ = d.Set("is_public", resp.IsPublic)
		_ = d.Set("latitude", resp.Latitude)
		_ = d.Set("location", resp.Location)
		_ = d.Set("longitude", resp.Longitude)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("upgrade_day", resp.UpgradeDay)
		_ = d.Set("upgrade_time_in_secs", resp.UpgradeTimeInSecs)
		_ = d.Set("override_version_profile", resp.OverrideVersionProfile)
		_ = d.Set("version_profile_id", resp.VersionProfileID)
		_ = d.Set("version_profile_name", resp.VersionProfileName)
		_ = d.Set("version_profile_visibility_scope", resp.VersionProfileVisibilityScope)
		_ = d.Set("trusted_networks", flattenTrustedNetworks(resp))
		_ = d.Set("service_edges", flattenServiceEdges(resp))

	} else {
		return fmt.Errorf("couldn't find any service edge group with name '%s' or id '%s'", name, id)
	}

	return nil
}

func flattenServiceEdges(serviceEdge *serviceedgegroup.ServiceEdgeGroup) []interface{} {
	serviceEdges := make([]interface{}, len(serviceEdge.ServiceEdges))
	for i, val := range serviceEdge.ServiceEdges {
		serviceEdges[i] = map[string]interface{}{
			"application_start_time":               val.ApplicationStartTime,
			"control_channel_status":               val.ControlChannelStatus,
			"creation_time":                        val.CreationTime,
			"ctrl_broker_name":                     val.CtrlBrokerName,
			"current_version":                      val.CurrentVersion,
			"description":                          val.Description,
			"enabled":                              val.Enabled,
			"expected_upgrade_time":                val.ExpectedUpgradeTime,
			"expected_version":                     val.ExpectedVersion,
			"fingerprint":                          val.Fingerprint,
			"id":                                   val.ID,
			"ipacl":                                val.IPACL,
			"issued_cert_id":                       val.IssuedCertID,
			"last_broker_connect_time":             val.LastBrokerConnectTime,
			"last_broker_connect_time_duration":    val.LastBrokerConnectTimeDuration,
			"last_broker_disconnect_time":          val.LastBrokerDisconnectTime,
			"last_broker_disconnect_time_duration": val.LastBrokerDisconnectTimeDuration,
			"last_upgrade_time":                    val.LastUpgradeTime,
			"latitude":                             val.Latitude,
			"location":                             val.Location,
			"longitude":                            val.Longitude,
			"listen_ips":                           val.ListenIPs,
			"modifiedby":                           val.ModifiedBy,
			"modified_time":                        val.ModifiedTime,
			"name":                                 val.Name,
			"provisioning_key_id":                  val.ProvisioningKeyID,
			"provisioning_key_name":                val.ProvisioningKeyName,
			"platform":                             val.Platform,
			"previous_version":                     val.PreviousVersion,
			"service_edge_group_id":                val.ServiceEdgeGroupID,
			"service_edge_group_name":              val.ServiceEdgeGroupName,
			"private_ip":                           val.PrivateIP,
			"public_ip":                            val.PublicIP,
			"publish_ips":                          val.PublishIPs,
			"sarge_version":                        val.SargeVersion,
			"enrollment_cert":                      val.EnrollmentCert,
			"upgrade_attempt":                      val.UpgradeAttempt,
			"upgrade_status":                       val.UpgradeStatus,
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
			"modifiedby":         val.ModifiedBy,
			"modified_time":      val.ModifiedTime,
			"name":               val.Name,
			"network_id":         val.NetworkID,
			"zscaler_cloud":      val.ZscalerCloud,
		}
	}

	return trustedNetworks
}
