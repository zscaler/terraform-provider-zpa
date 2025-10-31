package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/aup"
)

func dataSourceUserPortalAUP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserPortalAUPRead,
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
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"aup": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"phone_num": {
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
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserPortalAUPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var resp *aup.UserPortalAup
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for user portal controller %s\n", id)
		res, _, err := aup.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for user portal controller name %s\n", name)
		res, _, err := aup.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		d.Set("name", resp.Name)
		d.Set("description", resp.Description)
		d.Set("enabled", resp.Enabled)
		d.Set("aup", resp.Aup)
		d.Set("email", resp.Email)
		d.Set("phone_num", resp.PhoneNum)
		d.Set("microtenant_id", resp.MicrotenantID)
		d.Set("microtenant_name", resp.MicrotenantName)
		d.Set("creation_time", resp.CreationTime)
		d.Set("modified_time", resp.ModifiedTime)
		d.Set("modified_by", resp.ModifiedBy)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any user portal aup with name '%s' or id '%s'", name, id))
	}

	return nil
}
