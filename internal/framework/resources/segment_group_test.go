// Copyright (c) SecurityGeekIO, Inc.
// SPDX-License-Identifier: MPL-2.0

package resources_test

import (
	"context"
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
)

func TestAccSegmentGroup_basic(t *testing.T) {
	acctest.PreCheck(t)
	var segmentGroup segmentgroup.SegmentGroup

	rName := sdkacctest.RandString(8)
	resourceName := "zpa_segment_group.test"
	zpaClient := acctest.TestClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckSegmentGroupDestroy(zpaClient),
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentGroupConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSegmentGroupExists(zpaClient, resourceName, &segmentGroup),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-testacc-segment-group-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", "Test segment group for acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testAccSegmentGroupConfig_updated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSegmentGroupExists(zpaClient, resourceName, &segmentGroup),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-testacc-segment-group-updated-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSegmentGroupExists(zClient *client.Client, resourceName string, segmentGroup *segmentgroup.SegmentGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("segment group not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("segment group ID is not set")
		}

		ctx := context.Background()
		sg, _, err := segmentgroup.Get(ctx, zClient.Service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to retrieve segment group %s: %w", rs.Primary.ID, err)
		}

		*segmentGroup = *sg
		return nil
	}
}

func testAccCheckSegmentGroupDestroy(zClient *client.Client) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ctx := context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "zpa_segment_group" || rs.Primary.ID == "" {
				continue
			}

			_, _, err := segmentgroup.Get(ctx, zClient.Service, rs.Primary.ID)
			if err == nil {
				if _, delErr := segmentgroup.Delete(ctx, zClient.Service, rs.Primary.ID); delErr != nil {
					if respErr, ok := delErr.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
						continue
					}
					return fmt.Errorf("segment group %s still exists and failed to delete: %w", rs.Primary.ID, delErr)
				}
				continue
			}
			if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
				continue
			}
			return fmt.Errorf("error checking segment group %s destruction: %w", rs.Primary.ID, err)
		}
		return nil
	}
}

func testAccSegmentGroupConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "zpa_segment_group" "test" {
  name        = "tf-testacc-segment-group-%[1]s"
  description = "Test segment group for acceptance testing"
  enabled     = true
}
`, rName)
}

func testAccSegmentGroupConfig_updated(rName string) string {
	return fmt.Sprintf(`
resource "zpa_segment_group" "test" {
  name        = "tf-testacc-segment-group-updated-%[1]s"
  description = "Updated description for acceptance testing"
  enabled     = false
}
`, rName)
}
