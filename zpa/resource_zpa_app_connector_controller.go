package zpa

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/zscaler/zscaler-sdk-go/zpa"
)

func resourceAppConnectorController() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppConnectorControllerCreate,
		Read:   resourceAppConnectorControllerRead,
		Update: resourceAppConnectorControllerUpdate,
		Delete: resourceAppConnectorControllerDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := zClient.appconnectorcontroller.GetByName(id)
					if err == nil {
						d.SetId(resp.ID)
						_ = d.Set("id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}

				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of App Connector Controllers", // I think we can make this a typeSet as it does not need to be an ordered list
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"application_start_time": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"app_connector_group_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"app_connector_group_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			// "control_channel_status": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			// "creation_time": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			// "ctrl_broker_name": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			// "current_version": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			// "expected_upgrade_time": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			// "expected_version": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"ip_acl": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"issued_cert_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			// "last_broker_connect_time": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			// "last_broker_connect_time_duration": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			// "last_broker_disconnect_time": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			// "last_broker_disconnect_time_duration": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			// "last_upgrade_time": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			"latitude": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"longitude": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			// "modifiedby": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			// "modified_time": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provisioning_key_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"provisioning_key_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			// "previous_version": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"sarge_version": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"enrollment_cert": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			// "upgrade_attempt": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			// "upgrade_status": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
		},
	}
}

func resourceAppConnectorControllerCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
}

func resourceAppConnectorControllerRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.appconnectorcontroller.Get(d.Id())
	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing app connector controller %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting application server:\n%+v\n", resp)
	_ = d.Set("application_start_time", resp.ApplicationStartTime)
	_ = d.Set("app_connector_group_id", resp.AppConnectorGroupID)
	_ = d.Set("app_connector_group_name", resp.AppConnectorGroupName)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("fingerprint", resp.Fingerprint)
	_ = d.Set("ip_acl", resp.IPACL)
	_ = d.Set("issued_cert_id", resp.IssuedCertID)
	_ = d.Set("latitude", resp.Latitude)
	_ = d.Set("location", resp.Location)
	_ = d.Set("longitude", resp.Longitude)
	_ = d.Set("name", resp.Name)
	_ = d.Set("provisioning_key_id", resp.ProvisioningKeyID)
	_ = d.Set("provisioning_key_name", resp.ProvisioningKeyName)
	_ = d.Set("platform", resp.Platform)
	_ = d.Set("private_ip", resp.PrivateIP)
	_ = d.Set("public_ip", resp.PublicIP)
	_ = d.Set("sarge_version", resp.SargeVersion)
	_ = d.Set("enrollment_cert", resp.EnrollmentCert)
	// _ = d.Set("control_channel_status", resp.ControlChannelStatus)
	// _ = d.Set("creation_time", resp.CreationTime)
	// _ = d.Set("ctrl_broker_name", resp.CtrlBrokerName)
	// _ = d.Set("current_version", resp.CurrentVersion)
	// _ = d.Set("expected_upgrade_time", resp.ExpectedUpgradeTime)
	// _ = d.Set("expected_version", resp.ExpectedVersion)
	// _ = d.Set("modifiedby", resp.ModifiedBy)
	// _ = d.Set("modified_time", resp.ModifiedTime)
	// _ = d.Set("previous_version", resp.PreviousVersion)
	// _ = d.Set("last_broker_connect_time", resp.LastBrokerConnectTime)
	// _ = d.Set("last_broker_connect_time_duration", resp.LastBrokerConnectTimeDuration)
	// _ = d.Set("last_broker_disconnect_time", resp.LastBrokerDisconnectTime)
	// _ = d.Set("last_broker_disconnect_time_duration", resp.LastBrokerDisconnectTimeDuration)
	// _ = d.Set("last_upgrade_time", resp.LastUpgradeTime)
	// _ = d.Set("upgrade_attempt", resp.UpgradeAttempt)
	// _ = d.Set("upgrade_status", resp.UpgradeStatus)

	return nil

}

func resourceAppConnectorControllerUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
}

func resourceAppConnectorControllerDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting app connector controller ID: %v\n", d.Id())

	if _, err := zClient.appconnectorcontroller.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] app connector controller deleted")
	return nil
}
