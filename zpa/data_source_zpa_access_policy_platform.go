package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/platforms"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccessPolicyPlatforms() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccessPolicyPlatformsRead,
		Importer:    &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"linux": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"android": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"windows": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ios": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mac": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAccessPolicyPlatformsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	log.Printf("[INFO] Getting data for all platforms set\n")

	resp, _, err := platforms.GetAllPlatforms(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting data for all platforms:\n%+v\n", resp)
	d.SetId("platforms")
	if err := d.Set("linux", resp.Linux); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set linux: %v", err))
	}
	if err := d.Set("android", resp.Android); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set android: %v", err))
	}
	if err := d.Set("windows", resp.Windows); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set windows: %v", err))
	}
	if err := d.Set("ios", resp.IOS); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set ios: %v", err))
	}
	if err := d.Set("mac", resp.MacOS); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set mac: %v", err))
	}

	return nil
}
