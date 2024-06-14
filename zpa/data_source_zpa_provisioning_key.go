package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/provisioningkey"
)

func dataSourceProvisioningKey() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceProvisioningKeyRead,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microtenant_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceProvisioningKeyRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.ProvisioningKey

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	associationType, ok := getAssociationType(d)
	if !ok {
		return fmt.Errorf("associationType is required")
	}
	var resp *provisioningkey.ProvisioningKey
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data provisioning key %s\n", id)
		res, _, err := provisioningkey.Get(service, associationType, id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for provisioning key name %s\n", name)
		res, _, err := provisioningkey.GetByName(service, associationType, name)
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
		_ = d.Set("microtenant_id", resp.MicroTenantID)
		_ = d.Set("microtenant_name", resp.MicroTenantName)
	} else {
		return fmt.Errorf("couldn't find any provisioning key with name '%s' or id '%s'", name, id)
	}
	return nil
}
