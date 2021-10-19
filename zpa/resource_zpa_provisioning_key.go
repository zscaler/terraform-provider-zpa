package zpa

import (
	"log"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/provisioningkey"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/client"
)

func resourceProvisioningKey() *schema.Resource {
	return &schema.Resource{
		Create:   resourceProvisioningKeyCreate,
		Read:     resourceProvisioningKeyRead,
		Update:   resourceProvisioningKeyUpdate,
		Delete:   resourceProvisioningKeyDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"app_connector_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_connector_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"max_usage": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enrollment_cert_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enrollment_cert_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ui_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"usage_count": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zcomponent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zcomponent_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceProvisioningKeyCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandProvisioningKey(d)
	log.Printf("[INFO] Creating zpa provisining key with request\n%+v\n", req)

	resp, _, err := zClient.provisioningkey.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created provisining key  request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceProvisioningKeyRead(d, m)
}

func resourceProvisioningKeyRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.provisioningkey.Get(d.Id())
	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing provisining key %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting provisining key:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("app_connector_group_id", resp.AppConnectorGroupID)
	_ = d.Set("app_connector_group_name", resp.AppConnectorGroupName)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("max_usage", resp.MaxUsage)
	_ = d.Set("name", resp.Name)
	_ = d.Set("enrollment_cert_id", resp.EnrollmentCertID)
	_ = d.Set("enrollment_cert_name", resp.EnrollmentCertName)
	_ = d.Set("ui_config", resp.UIConfig)
	_ = d.Set("usage_count", resp.UsageCount)
	_ = d.Set("zcomponent_id", resp.ZcomponentID)
	_ = d.Set("zcomponent_name", resp.ZcomponentName)
	return nil

}

func resourceProvisioningKeyUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating provisining key ID: %v\n", id)
	req := expandProvisioningKey(d)

	if _, err := zClient.provisioningkey.Update(id, &req); err != nil {
		return err
	}

	return resourceProvisioningKeyRead(d, m)
}

func resourceProvisioningKeyDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting provisining key  ID: %v\n", d.Id())

	if _, err := zClient.provisioningkey.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] provisining key  deleted")
	return nil
}

func expandProvisioningKey(d *schema.ResourceData) provisioningkey.ProvisioningKey {
	provisioningKey := provisioningkey.ProvisioningKey{
		AppConnectorGroupID:   d.Get("app_connector_group_id").(string),
		AppConnectorGroupName: d.Get("app_connector_group_name").(string),
		Enabled:               d.Get("enabled").(bool),
		MaxUsage:              d.Get("max_usage").(string),
		Name:                  d.Get("name").(string),
		EnrollmentCertID:      d.Get("enrollment_cert_id").(string),
		EnrollmentCertName:    d.Get("enrollment_cert_name").(string),
		UIConfig:              d.Get("ui_config").(string),
		UsageCount:            d.Get("usage_count").(string),
		ZcomponentID:          d.Get("zcomponent_id").(string),
		ZcomponentName:        d.Get("zcomponent_name").(string),
	}
	return provisioningKey
}
