package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbibannercontroller"
)

func dataSourceCBIBanners() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCBIBannersRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_color": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"text_color": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notification_title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notification_text": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"logo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"banner": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceCBIBannersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *cbibannercontroller.CBIBannerController
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for cbi banner %s\n", id)
		res, _, err := cbibannercontroller.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data cbi banner name %s\n", name)
		res, _, err := cbibannercontroller.GetByNameOrID(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("primary_color", resp.PrimaryColor)
		_ = d.Set("text_color", resp.TextColor)
		_ = d.Set("notification_title", resp.NotificationTitle)
		_ = d.Set("notification_text", resp.NotificationText)
		_ = d.Set("logo", resp.Logo)
		_ = d.Set("banner", resp.Banner)
		_ = d.Set("is_default", resp.IsDefault)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any cbi banner with name '%s' or id '%s'", name, id))
	}

	return nil
}
