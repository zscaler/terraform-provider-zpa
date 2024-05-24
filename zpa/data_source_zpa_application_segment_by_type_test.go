package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
)

func TestAccDataSourceApplicationSegmentByType_Basic(t *testing.T) {
	// Generate random suffixes
	_, _, resourceNameSuffix := method.GenerateRandomSourcesTypeAndName("zpa_application_segment_pra")
	_, _, domainNameSuffix := method.GenerateRandomSourcesTypeAndName("zpa_application_segment_inspection")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceApplicationSegmentByTypeConfig_basic(resourceNameSuffix, domainNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceApplicationSegmentByTypeCheck("data.zpa_application_segment_by_type.pra"),
					testAccDataSourceApplicationSegmentByTypeCheck("data.zpa_application_segment_by_type.inspect"),
					testAccDataSourceApplicationSegmentByTypeCheck("data.zpa_application_segment_by_type.ba"),
				),
			},
		},
	})
}

func testAccDataSourceApplicationSegmentByTypeCheck(application_type string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(application_type, "application_type"),
		resource.TestCheckResourceAttrSet(application_type, "name"),
	)
}

func testAccCheckDataSourceApplicationSegmentByTypeConfig_basic(resourceNameSuffix, domainNameSuffix string) string {
	return fmt.Sprintf(`
resource "zpa_segment_group" "this" {
  name                   = "tf-acc-test-%s"
  description            = "tf-acc-test-%s"
  enabled                = true
}

resource "zpa_application_segment_pra" "this" {
  name             = "tf-acc-test-%s-1"
  description      = "tf-acc-test-%s-1"
  enabled          = true
  health_reporting = "ON_ACCESS"
  bypass_type      = "NEVER"
  is_cname_enabled = true
  tcp_port_ranges  = ["2222", "2222"]
  domain_names     = ["tests-%s.example.com"]
  segment_group_id = zpa_segment_group.this.id
  common_apps_dto {
    apps_config {
      name                 = "%s-app"
      domain               = "tests-%s.example.com"
      application_protocol = "SSH"
      application_port     = "2222"
      enabled = true
      app_types = [ "SECURE_REMOTE_ACCESS" ]
    }
  }
}

data "zpa_ba_certificate" "jenkins" {
  name = "jenkins.bd-hashicorp.com"
}

resource "zpa_application_segment_inspection" "this" {
  name             = "tf-acc-test-%s-2"
  description      = "tf-acc-test-%s-2"
  enabled          = true
  health_reporting = "ON_ACCESS"
  bypass_type      = "NEVER"
  is_cname_enabled = true
  tcp_port_ranges  = ["4444", "4444"]
  domain_names     = ["tests-%s.example.com"]
  segment_group_id = zpa_segment_group.this.id
  common_apps_dto {
    apps_config {
      name                 = "%s-app"
      domain               = "tests-%s.example.com"
      application_protocol = "HTTPS"
      application_port     = "4444"
      certificate_id       = data.zpa_ba_certificate.jenkins.id
      enabled              = true
      app_types            = [ "INSPECT" ]
    }
  }
}

resource "zpa_application_segment_browser_access" "this" {
    name                      = "tf-acc-test-%s-3"
    description               = "tf-acc-test-%s-3"
    enabled                   = true
    health_reporting          = "ON_ACCESS"
    bypass_type               = "NEVER"
    tcp_port_ranges           = ["4445", "4445"]
    domain_names              = ["tests-%s.example.com"]
    segment_group_id          = zpa_segment_group.this.id

    clientless_apps {
        name                  = "%s-app"
		    enabled               = true
		    domain                = "tests-%s.example.com"
        application_protocol  = "HTTPS"
        application_port      = "4445"
        certificate_id        = data.zpa_ba_certificate.jenkins.id
        trust_untrusted_cert  = true
    }
}

data "zpa_application_segment_by_type" "pra" {
  application_type = "SECURE_REMOTE_ACCESS"
  depends_on = [zpa_segment_group.this, zpa_application_segment_pra.this]
}

data "zpa_application_segment_by_type" "inspect" {
  application_type = "INSPECT"
  depends_on = [zpa_segment_group.this, zpa_application_segment_inspection.this]
}

data "zpa_application_segment_by_type" "ba" {
	application_type = "BROWSER_ACCESS"
	depends_on = [zpa_segment_group.this, zpa_application_segment_browser_access.this]
  }
`, resourceNameSuffix, resourceNameSuffix, resourceNameSuffix, resourceNameSuffix, domainNameSuffix, resourceNameSuffix, domainNameSuffix,
		resourceNameSuffix, resourceNameSuffix, domainNameSuffix, resourceNameSuffix, domainNameSuffix,
		resourceNameSuffix, resourceNameSuffix, domainNameSuffix, resourceNameSuffix, domainNameSuffix,
	)
}
