package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/provisioningkey"
)

func resourceProvisioningKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceProvisioningKeyCreate,
		Read:   resourceProvisioningKeyRead,
		Update: resourceProvisioningKeyUpdate,
		Delete: resourceProvisioningKeyDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				client := meta.(*Client)
				service := client.ProvisioningKey

				microTenantID := GetString(d.Get("microtenant_id"))
				if microTenantID != "" {
					service = service.WithMicroTenant(microTenantID)
				}

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				_, associationTypeSet := d.GetOk("association_type")
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
					if !associationTypeSet {
						_, assoc_type, _, err := provisioningkey.GetByIDAllAssociations(service, id)
						if err != nil {
							return []*schema.ResourceData{d}, err
						} else {
							_ = d.Set("association_type", assoc_type)
						}
					}
				} else {
					resp, assoc_type, _, err := provisioningkey.GetByNameAllAssociations(service, id)
					if err == nil {
						d.SetId(resp.ID)
						_ = d.Set("id", resp.ID)
						if !associationTypeSet {
							_ = d.Set("association_type", assoc_type)
						}
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
				Optional: true,
				Computed: true,
			},
			"app_connector_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_connector_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Read only property. Applicable only for GET calls, ignored in PUT/POST calls.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the provisioning key is enabled or not. Supported values: true, false",
			},
			"max_usage": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The maximum number of instances where this provisioning key can be used for enrolling an App Connector or Service Edge.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the provisioning key.",
			},
			"enrollment_cert_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the enrollment certificate that can be used for this provisioning key.",
			},
			"ui_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"usage_count": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The provisioning key utilization count.",
			},
			"zcomponent_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the existing App Connector or Service Edge Group.",
			},
			"zcomponent_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Read only property. Applicable only for GET calls, ignored in PUT/POST calls.",
			},
			"provisioning_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "read only field. Ignored in PUT/POST calls.",
			},
			"association_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Specifies the provisioning key type for App Connectors or ZPA Private Service Edges. The supported values are CONNECTOR_GRP and SERVICE_EDGE_GRP.",
				ValidateFunc: validation.StringInSlice(provisioningkey.ProvisioningKeyAssociationTypes, false),
			},
			"ip_acl": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func getAssociationType(d *schema.ResourceData) (string, bool) {
	val, ok := d.GetOk("association_type")
	if !ok {
		return "", ok
	}
	value, ok := val.(string)
	return value, ok
}

func resourceProvisioningKeyCreate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.ProvisioningKey

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	associationType, ok := getAssociationType(d)
	if !ok {
		return fmt.Errorf("associationType is required")
	}
	req := expandProvisioningKey(d)
	log.Printf("[INFO] Creating zpa provisining key with request\n%+v\n", req)

	resp, _, err := provisioningkey.Create(service, associationType, &req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created provisining key  request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceProvisioningKeyRead(d, meta)
}

func resourceProvisioningKeyRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.ProvisioningKey

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	associationType, ok := getAssociationType(d)
	if !ok {
		return fmt.Errorf("associationType is required")
	}
	resp, _, err := provisioningkey.Get(service, associationType, d.Id())
	if err != nil {
		if obj, ok := err.(*client.ErrorResponse); ok && obj.IsObjectNotFound() {
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
	_ = d.Set("ui_config", resp.UIConfig)
	_ = d.Set("usage_count", resp.UsageCount)
	_ = d.Set("zcomponent_id", resp.ZcomponentID)
	_ = d.Set("zcomponent_name", resp.ZcomponentName)
	_ = d.Set("ip_acl", resp.IPACL)
	_ = d.Set("provisioning_key", resp.ProvisioningKey)
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	return nil
}

func resourceProvisioningKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.ProvisioningKey

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	associationType, ok := getAssociationType(d)
	if !ok {
		return fmt.Errorf("associationType is required")
	}
	id := d.Id()
	log.Printf("[INFO] Updating provisining key ID: %v\n", id)
	req := expandProvisioningKey(d)
	if _, _, err := provisioningkey.Get(service, associationType, id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := provisioningkey.Update(service, associationType, id, &req); err != nil {
		return err
	}

	return resourceProvisioningKeyRead(d, meta)
}

func resourceProvisioningKeyDelete(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.ProvisioningKey

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	associationType, ok := getAssociationType(d)
	if !ok {
		return fmt.Errorf("associationType is required")
	}
	log.Printf("[INFO] Deleting provisining key  ID: %v\n", d.Id())

	if _, err := provisioningkey.Delete(service, associationType, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] provisining key  deleted")
	return nil
}

func expandProvisioningKey(d *schema.ResourceData) provisioningkey.ProvisioningKey {
	provisioningKey := provisioningkey.ProvisioningKey{
		ID:                    d.Id(),
		AppConnectorGroupID:   d.Get("app_connector_group_id").(string),
		AppConnectorGroupName: d.Get("app_connector_group_name").(string),
		Enabled:               d.Get("enabled").(bool),
		MaxUsage:              d.Get("max_usage").(string),
		Name:                  d.Get("name").(string),
		EnrollmentCertID:      d.Get("enrollment_cert_id").(string),
		UIConfig:              d.Get("ui_config").(string),
		UsageCount:            d.Get("usage_count").(string),
		ZcomponentID:          d.Get("zcomponent_id").(string),
		ZcomponentName:        d.Get("zcomponent_name").(string),
		MicroTenantID:         d.Get("microtenant_id").(string),
		IPACL:                 SetToStringList(d, "ip_acl"),
		ProvisioningKey:       d.Get("provisioning_key").(string),
	}
	return provisioningKey
}
