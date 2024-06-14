package zpa

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/platforms"
)

func dataSourceAccessPolicyPlatforms() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceAccessPolicyPlatformsRead,
		Importer: &schema.ResourceImporter{},

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

func dataSourceAccessPolicyPlatformsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.Platforms

	log.Printf("[INFO] Getting data for all platforms set\n")

	resp, _, err := platforms.GetAllPlatforms(service)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Getting data for all platforms:\n%+v\n", resp)
	d.SetId("platforms")
	_ = d.Set("linux", resp.Linux)
	_ = d.Set("android", resp.Android)
	_ = d.Set("windows", resp.Windows)
	_ = d.Set("ios", resp.IOS)
	_ = d.Set("mac", resp.MacOS)

	return nil
}
