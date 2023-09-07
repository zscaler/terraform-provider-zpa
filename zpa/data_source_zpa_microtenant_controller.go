package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/microtenants"
)

func dataSourceMicrotenantController() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMicrotenantControllerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"criteria_attribute": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"criteria_attribute_values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"operator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"priority": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"roles": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_role": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"user": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"comments": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"customer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eula": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"force_pwd_change": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"group_ids": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_locked": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"language_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_login_disabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"phone_number": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"pin_session": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_id": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"microtenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"microtenant_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timezone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tmp_password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"token_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"two_factor_auth_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"two_factor_auth_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modifiedby": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMicrotenantControllerRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *microtenants.MicroTenant
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for microtenant %s\n", id)
		res, _, err := zClient.microtenants.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}

	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for microtenant name %s\n", name)
		res, _, err := zClient.microtenants.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("criteria_attribute", resp.CriteriaAttribute)
		_ = d.Set("criteria_attribute_values", resp.CriteriaAttributeValues)
		_ = d.Set("operator", resp.Operator)
		_ = d.Set("priority", resp.Priority)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("roles", flattenMicrotenantRoles(resp))

		if resp.UserResource != nil {
			_ = d.Set("user", flattenMicroTenantUserResource(resp.UserResource))
		}
	} else {
		return fmt.Errorf("couldn't find any microtenants with name '%s' or id '%s'", name, id)
	}

	return nil
}

func flattenMicrotenantRoles(microTenant *microtenants.MicroTenant) []interface{} {
	microtenants := make([]interface{}, len(microTenant.Roles))
	for i, microtenantItem := range microTenant.Roles {
		microtenants[i] = map[string]interface{}{
			"id":          microtenantItem.ID,
			"name":        microtenantItem.Name,
			"custom_role": microtenantItem.CustomRole,
		}
	}

	return microtenants
}

func flattenMicroTenantUserResource(userResource *microtenants.UserResource) []map[string]interface{} {
	result := make([]map[string]interface{}, 1)
	result[0] = make(map[string]interface{})
	result[0]["id"] = userResource.ID
	result[0]["name"] = userResource.Name
	result[0]["description"] = userResource.Description
	result[0]["enabled"] = userResource.Enabled
	result[0]["comments"] = userResource.Comments
	result[0]["customer_id"] = userResource.CustomerID
	result[0]["display_name"] = userResource.DisplayName
	result[0]["email"] = userResource.Email
	result[0]["eula"] = userResource.Eula
	result[0]["force_pwd_change"] = userResource.ForcePwdChange
	result[0]["group_ids"] = userResource.GroupIDs
	result[0]["is_enabled"] = userResource.IsEnabled
	result[0]["is_locked"] = userResource.IsLocked
	result[0]["language_code"] = userResource.LanguageCode
	result[0]["local_login_disabled"] = userResource.LocalLoginDisabled
	result[0]["password"] = userResource.Password
	result[0]["phone_number"] = userResource.PhoneNumber
	result[0]["pin_session"] = userResource.PinSession
	result[0]["role_id"] = userResource.RoleID
	result[0]["microtenant_id"] = userResource.MicrotenantID
	result[0]["microtenant_name"] = userResource.MicrotenantName
	result[0]["time_zone"] = userResource.Timezone
	result[0]["tmp_password"] = userResource.TmpPassword
	result[0]["token_id"] = userResource.TokenID
	result[0]["two_factor_auth_enabled"] = userResource.TwoFactorAuthEnabled
	result[0]["two_factor_auth_type"] = userResource.TwoFactorAuthType
	result[0]["username"] = userResource.Username
	result[0]["creation_time"] = userResource.CreationTime
	result[0]["modified_by"] = userResource.ModifiedBy
	result[0]["modified_time"] = userResource.ModifiedTime
	return result
}

