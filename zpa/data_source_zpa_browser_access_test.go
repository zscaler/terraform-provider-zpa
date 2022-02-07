package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceBrowserAccess_ByIdAndName(t *testing.T) {
	rName := acctest.RandString(15)
	rDesc := acctest.RandString(15)
	sgrName := acctest.RandString(15)
	sgrDesc := acctest.RandString(15)
	resourceName := "data.zpa_browser_access.by_id"
	resourceName2 := "data.zpa_browser_access.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBrowserAccessByID(rName, rDesc, sgrName, sgrDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceBrowserAccess(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName2, "name", rName),
					resource.TestCheckResourceAttr(resourceName2, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName2, "enabled", "true"),
				),
				PreventPostDestroyRefresh: true,
			},
		},
	})
}

func testAccDataSourceBrowserAccessByID(rName, rDesc, sgName, sgDesc string) string {
	return fmt.Sprintf(`
	resource "zpa_browser_access" "test_browser_access" {
		name             = "%s"
		description      = "%s"
		enabled          = true
		health_reporting = "ON_ACCESS"
		bypass_type      = "NEVER"
		is_cname_enabled = true
		tcp_port_ranges = ["80", "80"]
		domain_names     = ["jenkins.securitygeek.io"]
		segment_group_id = zpa_segment_group.test_app_group.id
		clientless_apps {
			name                 = "jenkins.securitygeek.io"
			application_protocol = "HTTP"
			application_port     = "80"
			certificate_id       = data.zpa_ba_certificate.test_ba_cert.id
			trust_untrusted_cert = true
			enabled              = true
			domain               = "jenkins.securitygeek.io"
		  }
		  server_groups {
			id = []
		  }
	}

	resource "zpa_segment_group" "test_app_group" {
		name = "%s"
		description = "%s"
		enabled = true
	}


	data "zpa_ba_certificate" "test_ba_cert" {
		name = "jenkins.securitygeek.io"
	}

	data "zpa_browser_access" "by_name" {
		name = zpa_browser_access.test_browser_access.name
	}
	data "zpa_browser_access" "by_id" {
		id = zpa_browser_access.test_browser_access.id
	}
	`, rName, rDesc, sgName, sgName)
}

func testAccDataSourceBrowserAccess(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}
		return nil
	}
}
