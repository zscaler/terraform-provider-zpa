package zpa

import (
	"fmt"
	"log"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/provisioningkey"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func provisiningKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Whether the provisioning key is enabled or not. Supported values: true, false",
		},
		"expiration_in_epoch_sec": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_acl": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
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
			Computed: true,
		},
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"provisioning_key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"enrollment_cert_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of the enrollment certificate that can be used for this provisioning key.",
		},
		"enrollment_cert_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Read only property. Applicable only for GET calls, ignored in PUT/POST calls.",
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
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of the existing App Connector or Service Edge Group.",
		},
		"zcomponent_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Read only property. Applicable only for GET calls, ignored in PUT/POST calls.",
		},
		"association_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specifies the provisioning key type for App Connectors or ZPA Private Service Edges. The supported values are CONNECTOR_GRP and SERVICE_EDGE_GRP.",
			ValidateFunc: validation.StringInSlice([]string{
				"CONNECTOR_GRP", "SERVICE_EDGE_GRP",
			}, false),
		},
	}
}
func dataSourceProvisioningKey() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceProvisioningKeyRead,
		Importer: &schema.ResourceImporter{},

		Schema: MergeSchema(
			provisiningKeySchema(),
			map[string]*schema.Schema{
				"list": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: provisiningKeySchema(),
					},
				},
			}),
	}
}

func dataSourceProvisioningKeyRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	associationType, ok := getAssociationType(d)
	if !ok {
		return fmt.Errorf("associationType is required")
	}
	var resp *provisioningkey.ProvisioningKey
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data provisining key %s\n", id)
		res, _, err := zClient.provisioningkey.Get(associationType, id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for provisining key name %s\n", name)
		res, _, err := zClient.provisioningkey.GetByName(associationType, name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("id", resp.ID)
		_ = d.Set("app_connector_group_id", resp.AppConnectorGroupID)
		_ = d.Set("app_connector_group_name", resp.AppConnectorGroupName)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("expiration_in_epoch_sec", resp.ExpirationInEpochSec)
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
		_ = d.Set("list", flattenProvisionningKeyList([]provisioningkey.ProvisioningKey{*resp}))
	} else if id != "" || name != "" {
		return fmt.Errorf("couldn't find any provisining key with name '%s' or id '%s'", name, id)
	} else {
		// get the list
		list, _, err := zClient.provisioningkey.GetAll(associationType)
		if err != nil {
			return err
		}
		d.SetId("provisionning-key-list")
		_ = d.Set("list", flattenProvisionningKeyList(list))
	}
	return nil
}

func flattenProvisionningKeyList(list []provisioningkey.ProvisioningKey) []interface{} {
	keys := make([]interface{}, len(list))
	for i, item := range list {
		keys[i] = map[string]interface{}{
			"id":                       item.ID,
			"app_connector_group_id":   item.AppConnectorGroupID,
			"app_connector_group_name": item.AppConnectorGroupName,
			"creation_time":            item.CreationTime,
			"enabled":                  item.Enabled,
			"expiration_in_epoch_sec":  item.ExpirationInEpochSec,
			"ip_acl":                   item.IPACL,
			"max_usage":                item.MaxUsage,
			"modifiedby":               item.ModifiedBy,
			"modified_time":            item.ModifiedTime,
			"name":                     item.Name,
			"provisioning_key":         item.ProvisioningKey,
			"enrollment_cert_id":       item.EnrollmentCertID,
			"enrollment_cert_name":     item.EnrollmentCertName,
			"ui_config":                item.UIConfig,
			"usage_count":              item.UsageCount,
			"zcomponent_id":            item.ZcomponentID,
			"zcomponent_name":          item.ZcomponentName,
		}
	}
	return keys
}
