package zpa

import (
	"fmt"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/provisioningkey"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceProvisioningKeyAll() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceProvisioningKeyAllRead,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"association_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the provisioning key type for App Connectors or ZPA Private Service Edges. The supported values are CONNECTOR_GRP and SERVICE_EDGE_GRP.",
				ValidateFunc: validation.StringInSlice([]string{
					"CONNECTOR_GRP", "SERVICE_EDGE_GRP",
				}, false),
			},
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: provisioningKeySchema(),
				},
			},
		},
	}
}

func dataSourceProvisioningKeyAllRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	associationType, ok := getAssociationType(d)
	if !ok {
		return fmt.Errorf("associationType is required")
	}
	list, _, err := zClient.provisioningkey.GetAll(associationType)
	if err != nil {
		return err
	}
	d.SetId("provisionning-key-list")
	_ = d.Set("list", flattenProvisionningKeyList(list))
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
