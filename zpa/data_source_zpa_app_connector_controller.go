package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorcontroller"
)

func dataSourceAppConnectorController() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppConnectorControllerRead,
		Schema: map[string]*schema.Schema{
			"application_start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_connector_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_connector_group_name": {
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
			"platform_detail": {
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
			"runtime_os": {
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
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microtenant_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"assistant_version": {
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
						"app_connector_group_id": {
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
						"ctrl_channel_status": {
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
						"expected_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_broker_connect_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_broker_disconnect_time": {
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
						"latitude": {
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
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
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
						"mtunnel_id": {
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
					},
				},
			},
			"zpn_sub_module_upgrade_list": {
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
	}
}

func dataSourceAppConnectorControllerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var resp *appconnectorcontroller.AppConnector
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for app connector  %s\n", id)
		res, _, err := appconnectorcontroller.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err) // Wrap error using diag.FromErr
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for app connector name %s\n", name)
		res, _, err := appconnectorcontroller.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err) // Wrap error using diag.FromErr
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("id", resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("application_start_time", resp.ApplicationStartTime)
		_ = d.Set("app_connector_group_id", resp.AppConnectorGroupID)
		_ = d.Set("app_connector_group_name", resp.AppConnectorGroupName)
		_ = d.Set("control_channel_status", resp.ControlChannelStatus)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("ctrl_broker_name", resp.CtrlBrokerName)
		_ = d.Set("current_version", resp.CurrentVersion)
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
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("provisioning_key_id", resp.ProvisioningKeyID)
		_ = d.Set("provisioning_key_name", resp.ProvisioningKeyName)
		_ = d.Set("platform", resp.Platform)
		_ = d.Set("platform_detail", resp.PlatformDetail)
		_ = d.Set("previous_version", resp.PreviousVersion)
		_ = d.Set("private_ip", resp.PrivateIP)
		_ = d.Set("public_ip", resp.PublicIP)
		_ = d.Set("runtime_os", resp.RuntimeOS)
		_ = d.Set("sarge_version", resp.SargeVersion)
		_ = d.Set("enrollment_cert", resp.EnrollmentCert)
		_ = d.Set("upgrade_attempt", resp.UpgradeAttempt)
		_ = d.Set("upgrade_status", resp.UpgradeStatus)
		_ = d.Set("microtenant_id", resp.MicroTenantID)
		_ = d.Set("microtenant_name", resp.MicroTenantName)

		if err := d.Set("zpn_sub_module_upgrade_list", flattenZPNSubModuleUpgrade(resp.ZPNSubModuleUpgrade)); err != nil {
			return diag.FromErr(fmt.Errorf("failed to read app server groups %s", err))
		}

		if err := d.Set("assistant_version", flattenAssistantVersion(&resp.AssistantVersion)); err != nil {
			return diag.FromErr(fmt.Errorf("failed to read app server groups %s", err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any app connector with name '%s' or id '%s'", name, id))
	}

	return nil
}

func flattenAssistantVersion(assistantVersion *appconnectorcontroller.AssistantVersion) []interface{} {
	if assistantVersion == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	// Flattening the basic fields of PrivateBrokerVersion
	result["id"] = assistantVersion.ID
	result["application_start_time"] = assistantVersion.ApplicationStartTime
	result["app_connector_group_id"] = assistantVersion.AppConnectorGroupID
	result["broker_id"] = assistantVersion.BrokerId
	result["creation_time"] = assistantVersion.CreationTime
	result["ctrl_channel_status"] = assistantVersion.CtrlChannelStatus
	result["current_version"] = assistantVersion.CurrentVersion
	result["disable_auto_update"] = assistantVersion.DisableAutoUpdate
	result["expected_version"] = assistantVersion.ExpectedVersion
	result["last_broker_connect_time"] = assistantVersion.LastBrokerConnectTime
	result["last_broker_disconnect_time"] = assistantVersion.LastBrokerDisconnectTime
	result["last_upgraded_time"] = assistantVersion.LastUpgradedTime
	result["latitude"] = assistantVersion.Latitude
	result["lone_warrior"] = assistantVersion.LoneWarrior
	result["longitude"] = assistantVersion.Longitude
	result["modified_by"] = assistantVersion.ModifiedBy
	result["modified_time"] = assistantVersion.ModifiedTime
	result["mtunnel_id"] = assistantVersion.MtunnelID
	result["platform"] = assistantVersion.Platform
	result["platform_detail"] = assistantVersion.PlatformDetail
	result["previous_version"] = assistantVersion.PreviousVersion
	result["private_ip"] = assistantVersion.PrivateIP
	result["public_ip"] = assistantVersion.PublicIP
	result["restart_time_in_sec"] = assistantVersion.RestartTimeInSec
	result["runtime_os"] = assistantVersion.RuntimeOS
	result["sarge_version"] = assistantVersion.SargeVersion
	result["system_start_time"] = assistantVersion.SystemStartTime
	result["upgrade_attempt"] = assistantVersion.UpgradeAttempt
	result["upgrade_status"] = assistantVersion.UpgradeStatus
	result["upgrade_now_once"] = assistantVersion.UpgradeNowOnce

	return []interface{}{result}
}
