package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/customerversionprofile"
)

func dataSourceCustomerVersionProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCustomerVersionProfileRead,
		Schema: map[string]*schema.Schema{
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_scope_customer_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"customer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"exclude_constellation": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			// Acceptance tests returning panic: custom_scope_request_customer_ids.add_customer_ids: can only set full list
			// "custom_scope_request_customer_ids": {
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"add_customer_ids": {
			// 				Type:     schema.TypeString,
			// 				Computed: true,
			// 			},
			// 			"delete_customer_ids": {
			// 				Type:     schema.TypeString,
			// 				Computed: true,
			// 			},
			// 		},
			// 	},
			// },
			"customer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"upgrade_priority": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"customer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
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
						"platform": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"restart_after_uptime_in_days": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_profile_gid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"visibility_scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCustomerVersionProfileRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *customerversionprofile.CustomerVersionProfile
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for customer version profile %s\n", id)
		res, _, err := zClient.customerversionprofile.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data for customer version profile name %s\n", name)
		res, _, err := zClient.customerversionprofile.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("customer_id", resp.CustomerID)
		_ = d.Set("description", resp.Description)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("upgrade_priority", resp.UpgradePriority)
		_ = d.Set("visibility_scope", resp.VisibilityScope)
		// Acceptance tests returning panic: custom_scope_request_customer_ids.add_customer_ids: can only set full list
		// _ = d.Set("custom_scope_request_customer_ids.add_customer_ids", resp.CustomScopeRequestCustomerIDs.AddCustomerIDs)
		// _ = d.Set("custom_scope_request_customer_ids.delete_customer_ids", resp.CustomScopeRequestCustomerIDs.DeletecustomerIDs)

		if err := d.Set("custom_scope_customer_ids", flattenCustomerIDName(resp.CustomScopeCustomerIDs)); err != nil {
			return fmt.Errorf("failed to read custom scope customer ids %s", err)
		}

		if err := d.Set("versions", flattenVersions(resp.Versions)); err != nil {
			return fmt.Errorf("failed to read versions %s", err)
		}
	} else {
		return fmt.Errorf("couldn't find any customer version profile with name '%s' or id '%s'", name, id)
	}

	return nil
}

func flattenCustomerIDName(scopeCustomerID []customerversionprofile.CustomScopeCustomerIDs) []interface{} {
	scopeCustomerIDs := make([]interface{}, len(scopeCustomerID))
	for i, val := range scopeCustomerID {
		scopeCustomerIDs[i] = map[string]interface{}{
			"customer_id":           val.CustomerID,
			"exclude_constellation": val.ExcludeConstellation,
			"name":                  val.Name,
		}
	}
	return scopeCustomerIDs
}

func flattenVersions(version []customerversionprofile.Versions) []interface{} {
	versions := make([]interface{}, len(version))
	for i, val := range version {
		versions[i] = map[string]interface{}{
			"creation_time":                val.CreationTime,
			"customer_id":                  val.CustomerID,
			"id":                           val.ID,
			"modified_by":                  val.ModifiedBy,
			"modified_time":                val.ModifiedTime,
			"platform":                     val.Platform,
			"restart_after_uptime_in_days": val.RestartAfterUptimeInDays,
			"role":                         val.Role,
			"version":                      val.Version,
			"version_profile_gid":          val.VersionProfileGID,
		}
	}
	return versions
}
