package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgecontroller"
)

func dataSourceServiceEdgeController() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServiceEdgeControllerRead,
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
				Optional: true,
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
	}
}

func dataSourceServiceEdgeControllerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var resp *serviceedgecontroller.ServiceEdgeController
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for service edge controller %s\n", id)
		res, _, err := serviceedgecontroller.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for service edge controller name %s\n", name)
		res, _, err := serviceedgecontroller.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	// Ensure resp is not nil before proceeding
	if resp == nil {
		return diag.FromErr(fmt.Errorf("couldn't find any service edge controller with name '%s' or id '%s'", name, id))
	}

	// Set the values in Terraform schema
	d.SetId(resp.ID)
	_ = d.Set("application_start_time", resp.ApplicationStartTime)
	_ = d.Set("service_edge_group_id", resp.ServiceEdgeGroupID)
	_ = d.Set("service_edge_group_name", resp.ServiceEdgeGroupName)
	_ = d.Set("control_channel_status", resp.ControlChannelStatus)
	_ = d.Set("creation_time", resp.CreationTime)
	_ = d.Set("ctrl_broker_name", resp.CtrlBrokerName)
	_ = d.Set("current_version", resp.CurrentVersion)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("expected_upgrade_time", resp.ExpectedUpgradeTime)
	_ = d.Set("expected_version", resp.ExpectedVersion)
	_ = d.Set("fingerprint", resp.Fingerprint)
	_ = d.Set("ip_acl", resp.IPACL)
	_ = d.Set("issued_cert_id", resp.IssuedCertID)
	_ = d.Set("last_broker_connect_time", resp.LastBrokerConnectTime)
	_ = d.Set("last_broker_connect_time_duration", resp.LastBrokerConnectTimeDuration)
	_ = d.Set("last_broker_disconnect_time", resp.LastBrokerDisconnectTime)
	_ = d.Set("last_broker_disconnect_time_duration", resp.LastBrokerDisconnectTimeDuration)
	_ = d.Set("last_upgrade_time", resp.LastUpgradeTime)
	_ = d.Set("latitude", resp.Latitude)
	_ = d.Set("location", resp.Location)
	_ = d.Set("longitude", resp.Longitude)
	_ = d.Set("listen_ips", resp.ListenIPs)
	_ = d.Set("modified_by", resp.ModifiedBy)
	_ = d.Set("modified_time", resp.ModifiedTime)
	_ = d.Set("name", resp.Name)
	_ = d.Set("provisioning_key_id", resp.ProvisioningKeyID)
	_ = d.Set("provisioning_key_name", resp.ProvisioningKeyName)
	_ = d.Set("platform", resp.Platform)
	_ = d.Set("previous_version", resp.PreviousVersion)
	_ = d.Set("private_ip", resp.PrivateIP)
	_ = d.Set("public_ip", resp.PublicIP)
	_ = d.Set("publish_ips", resp.PublishIPs)
	_ = d.Set("publish_ipv6", resp.PublishIPv6)
	_ = d.Set("runtime_os", resp.RuntimeOS)
	_ = d.Set("sarge_version", resp.SargeVersion)
	_ = d.Set("enrollment_cert", resp.EnrollmentCert)
	_ = d.Set("upgrade_attempt", resp.UpgradeAttempt)
	_ = d.Set("upgrade_status", resp.UpgradeStatus)
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	_ = d.Set("microtenant_name", resp.MicroTenantName)

	// Check if the PrivateBrokerVersion has a valid ID before trying to flatten it
	if resp.PrivateBrokerVersion.ID != "" {
		_ = d.Set("private_broker_version", flattenPrivateBrokerVersion(&resp.PrivateBrokerVersion))
	} else {
		log.Printf("[WARN] PrivateBrokerVersion is empty for service edge controller %s\n", resp.ID)
	}

	return nil
}

func flattenPrivateBrokerVersion(privateBrokerVersion *serviceedgecontroller.PrivateBrokerVersion) []interface{} {
	if privateBrokerVersion == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	// Flattening the basic fields of PrivateBrokerVersion
	result["id"] = privateBrokerVersion.ID
	result["application_start_time"] = privateBrokerVersion.ApplicationStartTime
	result["broker_id"] = privateBrokerVersion.BrokerId
	result["creation_time"] = privateBrokerVersion.CreationTime
	result["current_version"] = privateBrokerVersion.CurrentVersion
	result["disable_auto_update"] = privateBrokerVersion.DisableAutoUpdate
	result["last_connect_time"] = privateBrokerVersion.LastConnectTime
	result["last_disconnect_time"] = privateBrokerVersion.LastDisconnectTime
	result["last_upgraded_time"] = privateBrokerVersion.LastUpgradedTime
	result["lone_warrior"] = privateBrokerVersion.LoneWarrior
	result["modified_by"] = privateBrokerVersion.ModifiedBy
	result["modified_time"] = privateBrokerVersion.ModifiedTime
	result["platform"] = privateBrokerVersion.Platform
	result["platform_detail"] = privateBrokerVersion.PlatformDetail
	result["previous_version"] = privateBrokerVersion.PreviousVersion
	result["service_edge_group_id"] = privateBrokerVersion.ServiceEdgeGroupID
	result["private_ip"] = privateBrokerVersion.PrivateIP
	result["public_ip"] = privateBrokerVersion.PublicIP
	result["restart_instructions"] = privateBrokerVersion.RestartInstructions
	result["restart_time_in_sec"] = privateBrokerVersion.RestartTimeInSec
	result["runtime_os"] = privateBrokerVersion.RuntimeOS
	result["sarge_version"] = privateBrokerVersion.SargeVersion
	result["system_start_time"] = privateBrokerVersion.SystemStartTime
	result["tunnel_id"] = privateBrokerVersion.TunnelId
	result["upgrade_attempt"] = privateBrokerVersion.UpgradeAttempt
	result["upgrade_status"] = privateBrokerVersion.UpgradeStatus
	result["upgrade_now_once"] = privateBrokerVersion.UpgradeNowOnce

	// Flatten the ZPNSubModuleUpgrade list
	result["zpn_sub_module_upgrade"] = flattenZPNSubModuleUpgrade(privateBrokerVersion.ZPNSubModuleUpgradeList)

	return []interface{}{result}
}

// Helper function to flatten ZPNSubModuleUpgrade
func flattenZPNSubModuleUpgrade(zpnSubModules []common.ZPNSubModuleUpgrade) []interface{} {
	if len(zpnSubModules) == 0 {
		return []interface{}{}
	}

	flattened := make([]interface{}, len(zpnSubModules))
	for i, subModule := range zpnSubModules {
		flattened[i] = map[string]interface{}{
			"id":               subModule.ID,
			"creation_time":    subModule.CreationTime,
			"current_version":  subModule.CurrentVersion,
			"entity_gid":       subModule.EntityGid,
			"entity_type":      subModule.EntityType,
			"expected_version": subModule.ExpectedVersion,
			"modified_by":      subModule.ModifiedBy,
			"modified_time":    subModule.ModifiedTime,
			"previous_version": subModule.PreviousVersion,
			"role":             subModule.Role,
			"upgrade_status":   subModule.UpgradeStatus,
			"upgrade_time":     subModule.UpgradeTime,
		}
	}

	return flattened
}
