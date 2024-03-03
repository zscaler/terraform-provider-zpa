package zpa

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/browseraccess"
)

func TestAccResourceApplicationSegmentBrowserAccessBasic(t *testing.T) {
	var browserAccess browseraccess.BrowserAccess
	browserAccessTypeAndName, _, browserAccessGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAApplicationSegmentBrowserAccess)
	rDomain := acctest.RandomWithPrefix("tf-acc-test")
	rDescription := acctest.RandomWithPrefix("tf-acc-test")
	updatedDescription := acctest.RandomWithPrefix("tf-acc-test-updated") // New name for update test

	// serverGroupTypeAndName, _, serverGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAServerGroup)
	// serverGroupHCL := testAccCheckServerGroupConfigure(serverGroupTypeAndName, serverGroupGeneratedName, "", "", "", "", variable.ServerGroupEnabled, variable.ServerGroupDynamicDiscovery)

	segmentGroupTypeAndName, _, segmentGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPASegmentGroup)
	segmentGroupHCL := testAccCheckSegmentGroupConfigure(segmentGroupTypeAndName, segmentGroupGeneratedName, variable.SegmentGroupDescription, variable.SegmentGroupEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationSegmentBrowserAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplicationSegmentBrowserAccessConfigure(browserAccessTypeAndName, browserAccessGeneratedName, browserAccessGeneratedName, rDescription, segmentGroupHCL, segmentGroupTypeAndName, variable.BrowserAccessEnabled, rDomain, variable.BrowserAccessCnameEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationSegmentBrowserAccessExists(browserAccessTypeAndName, &browserAccess),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "name", "tf-acc-test-"+browserAccessGeneratedName),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "description", rDescription),
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
				Config: testAccCheckApplicationSegmentBrowserAccessConfigure(browserAccessTypeAndName, browserAccessGeneratedName, browserAccessGeneratedName, updatedDescription, segmentGroupHCL, segmentGroupTypeAndName, variable.BrowserAccessEnabled, rDomain, variable.BrowserAccessCnameEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationSegmentBrowserAccessExists(browserAccessTypeAndName, &browserAccess),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "name", "tf-acc-test-"+browserAccessGeneratedName),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "description", updatedDescription),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "enabled", strconv.FormatBool(variable.BrowserAccessEnabled)),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "is_cname_enabled", strconv.FormatBool(variable.BrowserAccessCnameEnabled)),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "bypass_type", "NEVER"),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "health_reporting", "ON_ACCESS"),
					resource.TestCheckResourceAttrSet(browserAccessTypeAndName, "segment_group_id"),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "clientless_apps.#", "1"),
					resource.TestCheckResourceAttr(browserAccessTypeAndName, "tcp_port_range.#", "1"),
				),
			},
			// Import test
			{
				ResourceName:      browserAccessTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckApplicationSegmentBrowserAccessDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAApplicationSegmentBrowserAccess {
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

func testAccCheckApplicationSegmentBrowserAccessExists(resource string, segment *browseraccess.BrowserAccess) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedSegment, _, err := apiClient.browseraccess.Get(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*segment = *receivedSegment

		return nil
	}
}

func testAccCheckApplicationSegmentBrowserAccessConfigure(resourceTypeAndName, generatedName, name, description, segmentGroupHCL, segmentGroupTypeAndName string, enabled bool, rDomain string, cnameEnabled bool) string {
	port := strconv.Itoa(acctest.RandIntRange(4001, 5001))
	return fmt.Sprintf(`

// application segment browser access resource
%s

// application segment browser access resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		segmentGroupHCL,
		// serverGroupHCL,
		getBrowserAccessResourceHCL(generatedName, name, description, segmentGroupTypeAndName, enabled, rDomain, cnameEnabled, port),

		// data source variables
		resourcetype.ZPAApplicationSegmentBrowserAccess,
		generatedName,
		resourceTypeAndName,
	)
}

func getBrowserAccessResourceHCL(generatedName, name, description, segmentGroupTypeAndName string, enabled bool, rDomain string, cnameEnabled bool, port string) string {
	return fmt.Sprintf(`

data "zpa_ba_certificate" "jenkins" {
	name = "jenkins.bd-hashicorp.com"
}

resource "%s" "%s" {
	name = "tf-acc-test-%s"
	description = "%s"
	enabled = "%s"
	is_cname_enabled = "%s"
	select_connector_close_to_app = true
	health_reporting = "ON_ACCESS"
	bypass_type = "NEVER"
	tcp_keep_alive = "1"
	tcp_port_range {
		from = "%s"
		to = "%s"
	}
	domain_names = ["%s.bd-hashicorp.com"]
	segment_group_id = "${%s.id}"
	clientless_apps {
		name                 = "%s.bd-hashicorp.com"
		application_protocol = "HTTPS"
		application_port     = "%s"
		certificate_id       = data.zpa_ba_certificate.jenkins.id
		trust_untrusted_cert = true
		enabled              = true
		domain               = "%s.bd-hashicorp.com"
	}
	server_groups {
		id = []
	}
	depends_on = [ %s ]
}
`,

		// resource variables
		resourcetype.ZPAApplicationSegmentBrowserAccess,
		generatedName,
		name,
		description,
		strconv.FormatBool(enabled),
		strconv.FormatBool(cnameEnabled),
		port,
		port,
		rDomain,
		segmentGroupTypeAndName,
		rDomain,
		port,
		rDomain,
		segmentGroupTypeAndName,
	)
}
