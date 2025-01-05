package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbizpaprofile"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCBIZPAProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCBIZPAProfilesRead,
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

func dataSourceCBIZPAProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *cbizpaprofile.ZPAProfiles
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for cbi zpa profile name %s\n", name)
		res, _, err := cbizpaprofile.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
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
		return diag.FromErr(fmt.Errorf("couldn't find any cbi zpa profile with name '%s'", name))
	}

	return nil
}
