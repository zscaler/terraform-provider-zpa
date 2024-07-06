package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/cloudbrowserisolation/cbizpaprofile"
)

func dataSourceCBIZPAProfiles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCBIZPAProfilesRead,
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
			"cbi_tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cbi_profile_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cbi_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCBIZPAProfilesRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.CBIZpaProfile

	var resp *cbizpaprofile.ZPAProfiles

	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for cbi zpa profile name %s\n", name)
		res, _, err := cbizpaprofile.GetByName(service, name)
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
		_ = d.Set("cbi_tenant_id", resp.CBITenantID)
		_ = d.Set("cbi_profile_id", resp.CBIProfileID)
		_ = d.Set("cbi_url", resp.CBIURL)

	} else {
		return fmt.Errorf("couldn't find any cbi zpa profile with name '%s'", name)
	}

	return nil
}
