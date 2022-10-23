package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/appconnectorcontroller"
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
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of IDs for bulk deleting the Connectors",
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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}

// https://help.zscaler.com/zpa/connector-controller#/mgmtconfig/v1/admin/customers/{customerId}/connector/bulkDelete-post
func resourceAppConnectorControllerCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	if d.HasChange("ids") {
		i, ok := d.GetOk("ids")
		if ok {
			idsI, ok := i.([]interface{})
			if ok && len(idsI) > 0 {
				ids := []string{}
				for _, id := range idsI {
					ids = append(ids, id.(string))
				}
				_, err := zClient.appconnectorcontroller.BulkDelete(ids)
				if err != nil {
					return fmt.Errorf("failed bulk deleting application controllers (%#v):%v", ids, err)
				}
			}
		}

	}

	return nil
}

// https://help.zscaler.com/zpa/connector-controller#/mgmtconfig/v1/admin/customers/{customerId}/connector/{connectorId}-get
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

	return nil

}

func expandAppConnectorController(d *schema.ResourceData) appconnectorcontroller.AppConnector {
	return appconnectorcontroller.AppConnector{
		ID:                    d.Id(),
		ApplicationStartTime:  d.Get("application_start_time").(string),
		AppConnectorGroupID:   d.Get("app_connector_group_id").(string),
		AppConnectorGroupName: d.Get("app_connector_group_name").(string),
		Description:           d.Get("description").(string),
		Enabled:               d.Get("enabled").(bool),
		Fingerprint:           d.Get("fingerprint").(string),
		IPACL:                 d.Get("ip_acl").(string),
		IssuedCertID:          d.Get("issued_cert_id").(string),
		Latitude:              d.Get("latitude").(string),
		Location:              d.Get("location").(string),
		Longitude:             d.Get("longitude").(string),
		Name:                  d.Get("name").(string),
		ProvisioningKeyID:     d.Get("provisioning_key_id").(string),
		ProvisioningKeyName:   d.Get("provisioning_key_name").(string),
		Platform:              d.Get("platform").(string),
		PrivateIP:             d.Get("private_ip").(string),
		PublicIP:              d.Get("public_ip").(string),
		SargeVersion:          d.Get("sarge_version").(string),
		EnrollmentCert:        d.Get("enrollment_cert").(map[string]interface{}),
	}
}

// https://help.zscaler.com/zpa/connector-controller#/mgmtconfig/v1/admin/customers/{customerId}/connector/{connectorId}-put
func resourceAppConnectorControllerUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	app := expandAppConnectorController(d)
	_, _, err := zClient.appconnectorcontroller.Update(app.ID, app)
	if err != nil {
		return fmt.Errorf("failed updating app connector controller :%v", err)
	}
	return resourceAppConnectorControllerRead(d, m)
}

// https://help.zscaler.com/zpa/connector-controller#/mgmtconfig/v1/admin/customers/{customerId}/connector/{connectorId}-delete
func resourceAppConnectorControllerDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	_, err := zClient.appconnectorcontroller.Delete(d.Id())
	if err != nil {
		return fmt.Errorf("failed deleting app connector controller :%v", err)
	}
	d.SetId("")
	return nil
}
