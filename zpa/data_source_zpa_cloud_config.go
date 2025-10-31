package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/custom_config_controller"
)

func dataSourceZiaCloudConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceZiaCloudConfigRead,
		Schema: map[string]*schema.Schema{
			"zia_cloud_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zia_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceZiaCloudConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	log.Printf("[INFO] Getting ZIA cloud config")
	resp, _, err := custom_config_controller.GetZIACloudConfig(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp != nil {
		// Since there's no ID, use a static identifier
		d.SetId("zia_cloud_config")
		_ = d.Set("zia_cloud_domain", resp.ZIACloudDomain)
		_ = d.Set("zia_username", resp.ZIAUsername)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't retrieve ZIA cloud config"))
	}

	return nil
}
