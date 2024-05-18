package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/cloudbrowserisolation/isolationprofile"
)

func dataSourceIsolationProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIsolationProfileRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
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
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifiedby": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isolation_profile_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isolation_tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isolation_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIsolationProfileRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *isolationprofile.IsolationProfile
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for isolation profile name %s\n", name)
		res, _, err := zClient.isolationprofile.GetByName(name)
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
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("isolation_profile_id", resp.IsolationProfileID)
		_ = d.Set("isolation_tenant_id", resp.IsolationTenantID)
		_ = d.Set("isolation_url", resp.IsolationURL)
	} else {
		return fmt.Errorf("couldn't find any isolation profile with name '%s'", name)
	}

	return nil
}
