package zpa

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/microtenants"
)

func resourceMicrotenantController() *schema.Resource {
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
			"user": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"microtenant_id": {
							Type:     schema.TypeString,
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
	if microTenant.UserResource != nil {
		userList := flattenMicroTenantUser(microTenant.UserResource)
		_ = d.Set("user", userList)
		log.Printf("[DEBUG] Flattened User: %s", userList)
	}
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

	// log.Printf("[INFO] Getting microtenant:\n%+v\n", resp)
	log.Printf("[DEBUG] API Response: \n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("criteria_attribute", resp.CriteriaAttribute)
	_ = d.Set("criteria_attribute_values", resp.CriteriaAttributeValues)

	if resp.UserResource != nil {
		userList := flattenMicroTenantUser(resp.UserResource)
		_ = d.Set("user", userList)
		log.Printf("[DEBUG] Flattened User: %s", userList)
	}

	return nil
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
	}
	return microTenants
}

// flattenUserResource converts a UserResource struct into a map[string]interface{} that's suitable for use with Terraform.
func flattenMicroTenantUser(user *microtenants.UserResource) []map[string]interface{} {
	// Log the received user data
	log.Printf("[DEBUG] Received user data to flatten: %+v", user)

	if user == nil {
		return nil
	}

	flattenedData := []map[string]interface{}{{
		"display_name":   user.DisplayName,
		"username":       user.Username,
		"password":       user.Password,
		"microtenant_id": user.MicrotenantID,
	}}

	// Log the flattened data
	log.Printf("[DEBUG] Flattened user data: %+v", flattenedData)

	return flattenedData
}
