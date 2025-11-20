// Copyright (c) SecurityGeekIO, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasources_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

func TestAccSegmentGroupDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandString(8)
	resourceName := "zpa_segment_group.test"
	dataSourceName := "data.zpa_segment_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil, // Data source doesn't create resources
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentGroupDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enabled", resourceName, "enabled"),
				),
			},
		},
	})
}

func TestAccSegmentGroupDataSource_byName(t *testing.T) {
	rName := sdkacctest.RandString(8)
	resourceName := "zpa_segment_group.test"
	dataSourceName := "data.zpa_segment_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil, // Data source doesn't create resources
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentGroupDataSourceConfig_byName(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enabled", resourceName, "enabled"),
				),
			},
		},
	})
}

func TestAccSegmentGroupDataSource_byID(t *testing.T) {
	rName := sdkacctest.RandString(8)
	resourceName := "zpa_segment_group.test"
	dataSourceName := "data.zpa_segment_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil, // Data source doesn't create resources
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentGroupDataSourceConfig_byID(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enabled", resourceName, "enabled"),
				),
			},
		},
	})
}

func testAccSegmentGroupDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "zpa_segment_group" "test" {
  name        = "tf-testacc-segment-group-%[1]s"
  description = "Test segment group for acceptance testing"
  enabled     = true
}

data "zpa_segment_group" "test" {
  id = zpa_segment_group.test.id
}
`, rName)
}

func testAccSegmentGroupDataSourceConfig_byName(rName string) string {
	return fmt.Sprintf(`
resource "zpa_segment_group" "test" {
  name        = "tf-testacc-segment-group-%[1]s"
  description = "Test segment group for acceptance testing"
  enabled     = true
}

data "zpa_segment_group" "test" {
  name = zpa_segment_group.test.name
}
`, rName)
}

func testAccSegmentGroupDataSourceConfig_byID(rName string) string {
	return fmt.Sprintf(`
resource "zpa_segment_group" "test" {
  name        = "tf-testacc-segment-group-%[1]s"
  description = "Test segment group for acceptance testing"
  enabled     = true
}

data "zpa_segment_group" "test" {
  id = zpa_segment_group.test.id
}
`, rName)
}
