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
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.MicroTenants

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := microtenants.GetByName(service, id)
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the Microtenant for the ZPA tenant.",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the microtenant.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Microtenant.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether or not the Microtenant is enabled.",
			},
			"criteria_attribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The criteria attribute for the Microtenant. The supported value is AuthDomain.",
			},
			"criteria_attribute_values": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The value for the criteria attribute. This is the valid authentication domains configured for a customer (e.g., ExampleAuthDomain.com).",
			},
			"privileged_approvals_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if Privileged Approvals is enabled (true) for the Microtenant. This allows approval-based access even if no Authentication Domain is selected.",
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

func resourceMicrotenantCreate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.MicroTenants

	req := expandMicroTenant(d)
	log.Printf("[INFO] Creating microtenant with request\n%+v\n", req)

	microTenant, _, err := microtenants.Create(service, req)
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
	return resourceMicrotenantRead(d, meta)
}

func resourceMicrotenantRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.MicroTenants

	resp, _, err := microtenants.Get(service, d.Id())
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
	_ = d.Set("privileged_approvals_enabled", resp.PrivilegedApprovalsEnabled)

	if resp.UserResource != nil {
		userList := flattenMicroTenantUser(resp.UserResource)
		_ = d.Set("user", userList)
		log.Printf("[DEBUG] Flattened User: %s", userList)
	}

	return nil
}

func resourceMicrotenantUpdate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.MicroTenants

	id := d.Id()
	log.Printf("[INFO] Updating microtenant ID: %v\n", id)
	req := expandMicroTenant(d)

	if _, _, err := microtenants.Get(service, id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := microtenants.Update(service, id, &req); err != nil {
		return err
	}

	return resourceMicrotenantRead(d, meta)
}

func resourceMicrotenantDelete(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.MicroTenants

	log.Printf("[INFO] Deleting microtenant ID: %v\n", d.Id())

	if _, err := microtenants.Delete(service, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] microtenant deleted")
	return nil
}

func expandMicroTenant(d *schema.ResourceData) microtenants.MicroTenant {
	microTenants := microtenants.MicroTenant{
		ID:                         d.Id(),
		Name:                       d.Get("name").(string),
		Description:                d.Get("description").(string),
		Enabled:                    d.Get("enabled").(bool),
		CriteriaAttribute:          d.Get("criteria_attribute").(string),
		CriteriaAttributeValues:    SetToStringSlice(d.Get("criteria_attribute_values").(*schema.Set)),
		PrivilegedApprovalsEnabled: d.Get("privileged_approvals_enabled").(bool),
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
