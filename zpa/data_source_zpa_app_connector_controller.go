package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/appconnectorcontroller"
)

func dataSourceAppConnectorController() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppConnectorControllerRead,
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
	}
}

func dataSourceAppConnectorControllerRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *appconnectorcontroller.AppConnector
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for app connector  %s\n", id)
		res, _, err := zClient.appconnectorcontroller.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for app connector name %s\n", name)
		res, _, err := zClient.appconnectorcontroller.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("application_start_time", resp.ApplicationStartTime)
		_ = d.Set("app_connector_group_id", resp.AppConnectorGroupID)
		_ = d.Set("app_connector_group_name", resp.AppConnectorGroupName)
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
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("provisioning_key_id", resp.ProvisioningKeyID)
		_ = d.Set("provisioning_key_name", resp.ProvisioningKeyName)
		_ = d.Set("platform", resp.Platform)
		_ = d.Set("previous_version", resp.PreviousVersion)
		_ = d.Set("private_ip", resp.PrivateIP)
		_ = d.Set("public_ip", resp.PublicIP)
		_ = d.Set("sarge_version", resp.SargeVersion)
		_ = d.Set("enrollment_cert", resp.EnrollmentCert)
		_ = d.Set("upgrade_attempt", resp.UpgradeAttempt)
		_ = d.Set("upgrade_status", resp.UpgradeStatus)

	} else {
		return fmt.Errorf("couldn't find any app connector with name '%s' or id '%s'", name, id)
	}

	return nil
}
