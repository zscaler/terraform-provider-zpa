package zpa

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/microtenants"
)

func resourceMicrotenant() *schema.Resource {
	return &schema.Resource{
		Create: resourceMicrotenantCreate,
		Read:   resourceMicrotenantRead,
		Update: resourceMicrotenantUpdate,
		Delete: resourceMicrotenantDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := zClient.microtenants.GetByName(id)
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
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the microtenant.",
				Required:    true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"criteria_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"criteria_attribute_values": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"priority": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"custom_role": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceMicrotenantCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandMicroTenant(d)
	log.Printf("[INFO] Creating microtenant with request\n%+v\n", req)

	microTenant, _, err := zClient.microtenants.Create(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created microtenant request. ID: %v\n", microTenant)

	d.SetId(microTenant.ID)
	return resourceMicrotenantRead(d, m)

}

func resourceMicrotenantRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.microtenants.Get(d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing microtenant %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting microtenant:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("criteria_attribute", resp.CriteriaAttribute)
	_ = d.Set("criteria_attribute_values", resp.CriteriaAttributeValues)
	_ = d.Set("roles", flattenMicrotenantRolesSimple(resp.Roles))
	return nil
}

func flattenMicrotenantRolesSimple(apps []microtenants.Roles) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(apps))
	for i, app := range apps {
		ids[i] = app.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func resourceMicrotenantUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating microtenant ID: %v\n", id)
	req := expandMicroTenant(d)

	if _, _, err := zClient.microtenants.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := zClient.microtenants.Update(id, &req); err != nil {
		return err
	}

	return resourceMicrotenantRead(d, m)
}

func resourceMicrotenantDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting microtenant ID: %v\n", d.Id())

	if _, err := zClient.microtenants.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] microtenant deleted")
	return nil
}

func expandMicroTenant(d *schema.ResourceData) microtenants.MicroTenant {
	microTenants := microtenants.MicroTenant{
		ID:                      d.Id(),
		Name:                    d.Get("name").(string),
		Description:             d.Get("description").(string),
		Enabled:                 d.Get("enabled").(bool),
		CriteriaAttribute:       d.Get("criteria_attribute").(string),
		CriteriaAttributeValues: SetToStringSlice(d.Get("criteria_attribute_values").(*schema.Set)),
		Roles:                   expandMicroTenantRoles(d),
	}
	return microTenants
}

func expandMicroTenantRoles(d *schema.ResourceData) []microtenants.Roles {
	microtenantsInterface, ok := d.GetOk("roles")
	if ok {
		microTenant := microtenantsInterface.(*schema.Set)
		log.Printf("[INFO] server group application data: %+v\n", microTenant)
		var microTenants []microtenants.Roles
		for _, serverGroupApp := range microTenant.List() {
			microTenant, ok := serverGroupApp.(map[string]interface{})
			if ok {
				for _, id := range microTenant["id"].([]interface{}) {
					microTenants = append(microTenants, microtenants.Roles{
						ID: id.(string),
					})
				}
			}
		}
		return microTenants
	}

	return []microtenants.Roles{}
}
