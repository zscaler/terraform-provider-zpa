package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_controller"
)

func dataSourcePrivateCloudController() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateCloudControllerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
			"expected_sarge_version": {
				Type:     schema.TypeString,
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
			"ip_acl": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"last_os_upgrade_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_sarge_upgrade_time": {
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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"longitude": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_last_sync_time": {
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
			"provisioning_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_key_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_upgrade_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"os_upgrade_status": {
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
			"platform_version": {
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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"restriction_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"runtime": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sarge_upgrade_attempt": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sarge_upgrade_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sarge_version": {
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
			"shard_last_sync_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enrollment_cert": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     schema.TypeString,
			},
			"private_cloud_controller_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_cloud_controller_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_cloud_controller_version": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     schema.TypeString,
			},
			"site_sp_dns_name": {
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
			"userdb_last_sync_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_sub_module_upgrade_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"zscaler_managed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourcePrivateCloudControllerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var resp *private_cloud_controller.PrivateCloudController
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for private cloud controller %s\n", id)
		res, _, err := private_cloud_controller.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for private cloud controller name %s\n", name)
		res, _, err := private_cloud_controller.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.Id)
		_ = d.Set("name", resp.Name)
		_ = d.Set("application_start_time", resp.ApplicationStartTime)
		_ = d.Set("control_channel_status", resp.ControlChannelStatus)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("ctrl_broker_name", resp.CtrlBrokerName)
		_ = d.Set("current_version", resp.CurrentVersion)
		_ = d.Set("description", resp.Description)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("expected_sarge_version", resp.ExpectedSargeVersion)
		_ = d.Set("expected_upgrade_time", resp.ExpectedUpgradeTime)
		_ = d.Set("expected_version", resp.ExpectedVersion)
		_ = d.Set("fingerprint", resp.Fingerprint)
		_ = d.Set("ip_acl", resp.IpAcl)
		_ = d.Set("issued_cert_id", resp.IssuedCertId)
		_ = d.Set("last_broker_connect_time", resp.LastBrokerConnectTime)
		_ = d.Set("last_broker_connect_time_duration", resp.LastBrokerConnectTimeDuration)
		_ = d.Set("last_broker_disconnect_time", resp.LastBrokerDisconnectTime)
		_ = d.Set("last_broker_disconnect_time_duration", resp.LastBrokerDisconnectTimeDuration)
		_ = d.Set("last_os_upgrade_time", resp.LastOsUpgradeTime)
		_ = d.Set("last_sarge_upgrade_time", resp.LastSargeUpgradeTime)
		_ = d.Set("last_upgrade_time", resp.LastUpgradeTime)
		_ = d.Set("latitude", resp.Latitude)
		_ = d.Set("listen_ips", resp.ListenIps)
		_ = d.Set("location", resp.Location)
		_ = d.Set("longitude", resp.Longitude)
		_ = d.Set("master_last_sync_time", resp.MasterLastSyncTime)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("provisioning_key_id", resp.ProvisioningKeyId)
		_ = d.Set("provisioning_key_name", resp.ProvisioningKeyName)
		_ = d.Set("os_upgrade_enabled", resp.OsUpgradeEnabled)
		_ = d.Set("os_upgrade_status", resp.OsUpgradeStatus)
		_ = d.Set("platform", resp.Platform)
		_ = d.Set("platform_detail", resp.PlatformDetail)
		_ = d.Set("platform_version", resp.PlatformVersion)
		_ = d.Set("previous_version", resp.PreviousVersion)
		_ = d.Set("private_ip", resp.PrivateIp)
		_ = d.Set("public_ip", resp.PublicIp)
		_ = d.Set("publish_ips", resp.PublishIps)
		_ = d.Set("read_only", resp.ReadOnly)
		_ = d.Set("restriction_type", resp.RestrictionType)
		_ = d.Set("runtime", resp.Runtime)
		_ = d.Set("sarge_upgrade_attempt", resp.SargeUpgradeAttempt)
		_ = d.Set("sarge_upgrade_status", resp.SargeUpgradeStatus)
		_ = d.Set("sarge_version", resp.SargeVersion)
		_ = d.Set("microtenant_id", resp.MicrotenantId)
		_ = d.Set("microtenant_name", resp.MicrotenantName)
		_ = d.Set("shard_last_sync_time", resp.ShardLastSyncTime)
		_ = d.Set("enrollment_cert", resp.EnrollmentCert)
		_ = d.Set("private_cloud_controller_group_id", resp.PrivateCloudControllerGroupId)
		_ = d.Set("private_cloud_controller_group_name", resp.PrivateCloudControllerGroupName)
		_ = d.Set("private_cloud_controller_version", resp.PrivateCloudControllerVersion)
		_ = d.Set("site_sp_dns_name", resp.SiteSpDnsName)
		_ = d.Set("upgrade_attempt", resp.UpgradeAttempt)
		_ = d.Set("upgrade_status", resp.UpgradeStatus)
		_ = d.Set("userdb_last_sync_time", resp.UserdbLastSyncTime)
		_ = d.Set("zpn_sub_module_upgrade_list", resp.ZpnSubModuleUpgradeList)
		_ = d.Set("zscaler_managed", resp.ZscalerManaged)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any private cloud controller with name '%s' or id '%s'", name, id))
	}

	return nil
}
