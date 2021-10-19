package zpa

/*
import (
	"fmt"
	"log"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/provisioningkey"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProvisioningKey() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceProvisioningKeyRead,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"app_connector_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_connector_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"expiration_in_epoch_sec": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_acl": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_usage": {
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
			"provisioning_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enrollment_cert_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enrollment_cert_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ui_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"usage_count": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zcomponent_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zcomponent_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceProvisioningKeyRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *provisioningkey.ProvisioningKey
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data provisining key %s\n", id)
		res, _, err := zClient.provisioningkey.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for provisining key name %s\n", name)
		res, _, err := zClient.provisioningkey.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("app_connector_group_id", resp.AppConnectorGroupID)
		_ = d.Set("app_connector_group_name", resp.AppConnectorGroupName)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("expiration_in_epoch_sec", resp.ExpirationInEpochSec)
		_ = d.Set("ip_acl", resp.IPACL)
		_ = d.Set("ip_acl", resp.IPACL)
		_ = d.Set("max_usage", resp.MaxUsage)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("provisioning_key", resp.ProvisioningKey)
		_ = d.Set("enrollment_cert_id", resp.EnrollmentCertID)
		_ = d.Set("enrollment_cert_name", resp.EnrollmentCertName)
		_ = d.Set("ui_config", resp.UIConfig)
		_ = d.Set("usage_count", resp.UsageCount)
		_ = d.Set("zcomponent_id", resp.ZcomponentID)
		_ = d.Set("zcomponent_name", resp.ZcomponentName)

	} else {
		return fmt.Errorf("couldn't find any provisining key with name '%s' or id '%s'", name, id)
	}

	return nil
}
*/
