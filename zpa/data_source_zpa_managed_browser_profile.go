package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/managed_browser"
)

func dataSourceManagedBrowserProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceManagedBrowserProfileRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"browser_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_id": {
				Type:     schema.TypeString,
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
			"chrome_posture_profile": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"browser_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"crowd_strike_agent": {
							Type:     schema.TypeBool,
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
					},
				},
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
		},
	}
}

func dataSourceManagedBrowserProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *managed_browser.ManagedBrowserProfile
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for managed browser profile %s\n", id)
		// Get all managed browser profiles and find the one with matching ID
		allProfiles, _, err := managed_browser.GetAll(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, profile := range allProfiles {
			if profile.ID == id {
				resp = &profile
				break
			}
		}
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for managed browser profile name %s\n", name)
		// Use GetByName for direct name lookup
		res, _, err := managed_browser.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("browser_type", resp.BrowserType)
		_ = d.Set("customer_id", resp.CustomerID)
		_ = d.Set("microtenant_id", resp.MicrotenantID)
		_ = d.Set("microtenant_name", resp.MicrotenantName)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("chrome_posture_profile", flattenChromePostureProfile(resp.ChromePostureProfile))
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any managed browser profile with name '%s' or id '%s'", name, id))
	}

	return nil
}

// flattenChromePostureProfile flattens the ChromePostureProfile struct
func flattenChromePostureProfile(profile managed_browser.ChromePostureProfile) []interface{} {
	result := map[string]interface{}{
		"id":                 profile.ID,
		"browser_type":       profile.BrowserType,
		"crowd_strike_agent": profile.CrowdStrikeAgent,
		"creation_time":      profile.CreationTime,
		"modified_by":        profile.ModifiedBy,
		"modified_time":      profile.ModifiedTime,
	}

	return []interface{}{result}
}
