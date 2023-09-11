package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCBIRegions_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceCBIRegionsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceCBIRegionsCheck("data.zpa_cloud_browser_isolation_region.singapore"),
					testAccDataSourceCBIRegionsCheck("data.zpa_cloud_browser_isolation_region.washington"),
					testAccDataSourceCBIRegionsCheck("data.zpa_cloud_browser_isolation_region.portland"),
					testAccDataSourceCBIRegionsCheck("data.zpa_cloud_browser_isolation_region.london"),
					testAccDataSourceCBIRegionsCheck("data.zpa_cloud_browser_isolation_region.frankfurt"),
					testAccDataSourceCBIRegionsCheck("data.zpa_cloud_browser_isolation_region.sydney"),
					testAccDataSourceCBIRegionsCheck("data.zpa_cloud_browser_isolation_region.mumbai"),
					testAccDataSourceCBIRegionsCheck("data.zpa_cloud_browser_isolation_region.tokyo"),
					testAccDataSourceCBIRegionsCheck("data.zpa_cloud_browser_isolation_region.ohio"),
					testAccDataSourceCBIRegionsCheck("data.zpa_cloud_browser_isolation_region.paris"),
					testAccDataSourceCBIRegionsCheck("data.zpa_cloud_browser_isolation_region.hyderabad"),
					testAccDataSourceCBIRegionsCheck("data.zpa_cloud_browser_isolation_region.hongkong"),
				),
			},
		},
	})
}

func testAccDataSourceCBIRegionsCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceCBIRegionsConfig_basic = `

data "zpa_cloud_browser_isolation_region" "washington" {
    name = "Washington"
}

data "zpa_cloud_browser_isolation_region" "singapore" {
    name = "Singapore"
}

data "zpa_cloud_browser_isolation_region" "portland" {
    name = "Portland Oregon"
}

data "zpa_cloud_browser_isolation_region" "london" {
    name = "London"
}

data "zpa_cloud_browser_isolation_region" "frankfurt" {
    name = "Frankfurt"
}

data "zpa_cloud_browser_isolation_region" "sydney" {
    name = "Sydney"
}

data "zpa_cloud_browser_isolation_region" "mumbai" {
    name = "Mumbai"
}

data "zpa_cloud_browser_isolation_region" "tokyo" {
    name = "Tokyo"
}

data "zpa_cloud_browser_isolation_region" "ohio" {
    name = "Ohio"
}

data "zpa_cloud_browser_isolation_region" "paris" {
    name = "Paris"
}

data "zpa_cloud_browser_isolation_region" "hyderabad" {
    name = "Hyderabad"
}

data "zpa_cloud_browser_isolation_region" "hongkong" {
    name = "HongKong"
}
`
