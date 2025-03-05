package zpa

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/customerversionprofile"
)

func dataSourceCustomerVersionProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomerVersionProfileRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
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
			"number_of_assistants": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"number_of_customers": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"number_of_private_brokers": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"number_of_site_controllers": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"number_of_updated_assistants": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"number_of_updated_private_brokers": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"number_of_updated_site_controllers": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"upgrade_priority": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"visibility_scope": {
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"exclude_constellation": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_partner": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"custom_scope_request_customer_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"add_customer_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"delete_customer_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
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
		},
	}
}

func dataSourceCustomerVersionProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *customerversionprofile.CustomerVersionProfile
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for customer version profile name %s\n", name)
		res, _, err := customerversionprofile.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("customer_id", resp.CustomerID)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("upgrade_priority", resp.UpgradePriority)
		_ = d.Set("visibility_scope", resp.VisibilityScope)

		if err := d.Set("custom_scope_customer_ids", flattenScopeCustomerIDs(resp.CustomScopeCustomerIDs)); err != nil {
			return diag.FromErr(fmt.Errorf("failed to read custom scope customer ids %s", err))
		}

		if err := d.Set("custom_scope_request_customer_ids", flattenScopeRequestCustomerIDs(resp.CustomScopeRequestCustomerIDs)); err != nil {
			return diag.FromErr(fmt.Errorf("failed to read custom scope request customer ids %s", err))
		}
		if err := d.Set("versions", flattenVersions(resp.Versions)); err != nil {
			return diag.FromErr(fmt.Errorf("failed to read versions %s", err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any customer version profilee with name '%s'", name))
	}

	return nil
}

func flattenScopeCustomerIDs(scopeCustomerID []customerversionprofile.CustomScopeCustomerIDs) []interface{} {
	scopeCustomerIDs := make([]interface{}, len(scopeCustomerID))
	for i, val := range scopeCustomerID {
		scopeCustomerIDs[i] = map[string]interface{}{
			"name":                  val.Name,
			"customer_id":           val.CustomerID, // Ensure customer ID is a string
			"is_partner":            val.IsPartner,  // Convert boolean to string
			"exclude_constellation": val.ExcludeConstellation,
		}
	}
	return scopeCustomerIDs
}

func flattenScopeRequestCustomerIDs(requestIDs customerversionprofile.CustomScopeRequestCustomerIDs) []interface{} {
	if requestIDs.AddCustomerIDs == "" && requestIDs.DeletecustomerIDs == "" {
		return nil // Return nil if both are empty to prevent unnecessary block creation
	}

	return []interface{}{
		map[string]interface{}{
			"add_customer_ids":    strings.Split(requestIDs.AddCustomerIDs, ","),    // Convert comma-separated string to list
			"delete_customer_ids": strings.Split(requestIDs.DeletecustomerIDs, ","), // Convert comma-separated string to list
		},
	}
}

func flattenVersions(version []customerversionprofile.Versions) []interface{} {
	if len(version) == 0 {
		return nil // Avoids unnecessary empty slice allocation
	}

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
