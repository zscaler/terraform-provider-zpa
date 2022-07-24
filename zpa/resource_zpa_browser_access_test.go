package zpa

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/browseraccess"
)

func TestAccResourceBrowserAccessBasic(t *testing.T) {
	var browserAccess browseraccess.BrowserAccess
	browserAccessTypeAndName, _, browserAccessGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPABrowserAccess)
	// rPort := acctest.RandIntRange(1000, 9999)

	serverGroupTypeAndName, _, serverGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAServerGroup)
	serverGroupHCL := testAccCheckServerGroupConfigure(serverGroupTypeAndName, serverGroupGeneratedName, "", "", "", "", variable.ServerGroupEnabled, variable.ServerGroupDynamicDiscovery)

	segmentGroupTypeAndName, _, segmentGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPASegmentGroup)
	segmentGroupHCL := testAccCheckSegmentGroupConfigure(segmentGroupTypeAndName, segmentGroupGeneratedName, variable.SegmentGroupDescription, variable.SegmentGroupEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBrowserAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckBrowserAccessConfigure(browserAccessTypeAndName, browserAccessGeneratedName, browserAccessGeneratedName, browserAccessGeneratedName, segmentGroupHCL, segmentGroupTypeAndName, serverGroupHCL, serverGroupTypeAndName, variable.BrowserAccessEnabled, variable.BrowserAccessCnameEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBrowserAccessExists(browserAccessTypeAndName, &browserAccess),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "name", "tf-acc-test-"+browserAccessGeneratedName),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "description", "tf-acc-test-"+browserAccessGeneratedName),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "enabled", strconv.FormatBool(variable.BrowserAccessEnabled)),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "is_cname_enabled", strconv.FormatBool(variable.BrowserAccessCnameEnabled)),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "bypass_type", "NEVER"),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "health_reporting", "ON_ACCESS"),
					resource.TestCheckResourceAttrSet(browserAccessTypeAndName, "segment_group_id"),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "clientless_apps.#", "1"),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "tcp_port_range.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckBrowserAccessConfigure(browserAccessTypeAndName, browserAccessGeneratedName, browserAccessGeneratedName, browserAccessGeneratedName, segmentGroupHCL, segmentGroupTypeAndName, serverGroupHCL, serverGroupTypeAndName, variable.BrowserAccessEnabled, variable.BrowserAccessCnameEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBrowserAccessExists(browserAccessTypeAndName, &browserAccess),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "name", "tf-acc-test-"+browserAccessGeneratedName),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "description", "tf-acc-test-"+browserAccessGeneratedName),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "enabled", strconv.FormatBool(variable.BrowserAccessEnabled)),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "is_cname_enabled", strconv.FormatBool(variable.BrowserAccessCnameEnabled)),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "bypass_type", "NEVER"),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "health_reporting", "ON_ACCESS"),
					resource.TestCheckResourceAttrSet(browserAccessTypeAndName, "segment_group_id"),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "clientless_apps.#", "1"),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "tcp_port_range.#", "1"),
				),
			},
		},
	})
}

func testAccCheckBrowserAccessDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPABrowserAccess {
			continue
		}

		_, _, err := client.browseraccess.GetByName(rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Broser Access still exists")
		}

		return nil
	}
	return nil
}

func testAccCheckBrowserAccessExists(resource string, segment *browseraccess.BrowserAccess) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Broser Access Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Broser Access ID is set")
		}
		client := testAccProvider.Meta().(*Client)
		resp, _, err := client.browseraccess.GetByName(rs.Primary.Attributes["name"])
		if err != nil {
			return err
		}
		if resp.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("name Not found in created attributes")
		}
		if resp.Description != rs.Primary.Attributes["description"] {
			return fmt.Errorf("description Not found in created attributes")
		}
		return nil
	}
}

func testAccCheckBrowserAccessConfigure(resourceTypeAndName, generatedName, name, description, segmentGroupHCL, segmentGroupTypeAndName, serverGroupHCL, serverGroupTypeAndName string, enabled, cnameEnabled bool) string {
	return fmt.Sprintf(`

// segment group resource
%s

// application segment resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		segmentGroupHCL,
		// serverGroupHCL,
		getBrowserAccessResourceHCL(generatedName, name, description, segmentGroupTypeAndName, serverGroupTypeAndName, enabled, cnameEnabled),

		// data source variables
		resourcetype.ZPABrowserAccess,
		generatedName,
		resourceTypeAndName,
	)
}

func getBrowserAccessResourceHCL(generatedName, name, description, segmentGroupTypeAndName, serverGroupTypeAndName string, enabled, cnameEnabled bool) string {
	return fmt.Sprintf(`

data "zpa_ba_certificate" "jenkins" {
	name = "jenkins.bd-hashicorp.com"
}

resource "%s" "%s" {
	name = "tf-acc-test-%s"
	description = "tf-acc-test-%s"
	enabled = "%s"
	is_cname_enabled = "%s"
	health_reporting = "ON_ACCESS"
	bypass_type = "NEVER"
	tcp_port_range {
		from = "443"
		to = "443"
	}
	domain_names = ["jenkins.bd-hashicorp.com"]
	segment_group_id = "${%s.id}"
	clientless_apps {
		name                 = "jenkins.bd-hashicorp.com"
		application_protocol = "HTTPS"
		application_port     = "443"
		certificate_id       = data.zpa_ba_certificate.jenkins.id
		trust_untrusted_cert = true
		enabled              = true
		domain               = "jenkins.bd-hashicorp.com"
	}
	server_groups {
		id = []
	}
	depends_on = [ %s ]
}
`,

		// resource variables
		resourcetype.ZPABrowserAccess,
		generatedName,
		generatedName,
		generatedName,
		strconv.FormatBool(enabled),
		strconv.FormatBool(cnameEnabled),
		segmentGroupTypeAndName,
		// serverGroupTypeAndName,
		segmentGroupTypeAndName,
		// serverGroupTypeAndName,
	)
}
