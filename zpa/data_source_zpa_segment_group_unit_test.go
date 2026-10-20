package zpa

import (
	"context"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
)

func TestDataSourceSegmentGroup_ByName(t *testing.T) {
	fx := newUnitFixture(t)
	groups := []segmentgroup.SegmentGroup{
		{ID: "1001", Name: "Apps", Description: "Business apps", Enabled: true},
		{ID: "1002", Name: "Apps Legacy", Enabled: false},
	}
	fx.on("GET", fx.path("segmentGroup"), common.SuccessResponse(common.ZPAList(groups)))

	d := fx.resourceData(dataSourceSegmentGroup(), "", map[string]any{"name": "apps"})
	if diags := dataSourceSegmentGroupRead(context.Background(), d, fx.client); diags.HasError() {
		t.Fatalf("read failed: %v", diags)
	}
	// The name lookup is case-insensitive and must pick the exact match, not the first result.
	if d.Id() != "1001" || d.Get("description") != "Business apps" {
		t.Errorf("got id=%q description=%v, want 1001 / Business apps", d.Id(), d.Get("description"))
	}
	reqs := fx.api.Server.AllRequests()
	if len(reqs) == 0 || !strings.Contains(reqs[0].Query, "search=name%2BEQ%2Bapps") {
		t.Errorf("the lookup must pass the name as a search filter (name EQ apps), got %v", reqs)
	}
}

func TestDataSourceSegmentGroup_ByID(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("segmentGroup", "1001"),
		common.SuccessResponse(segmentgroup.SegmentGroup{ID: "1001", Name: "Apps", Enabled: true}))

	// For data sources the "id" argument is Terraform's resource ID.
	d := fx.resourceData(dataSourceSegmentGroup(), "1001", map[string]any{})
	if diags := dataSourceSegmentGroupRead(context.Background(), d, fx.client); diags.HasError() {
		t.Fatalf("read failed: %v", diags)
	}
	if d.Get("name") != "Apps" {
		t.Errorf("name = %v, want Apps", d.Get("name"))
	}
}

func TestDataSourceSegmentGroup_NameNotFound(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("segmentGroup"), common.SuccessResponse(common.ZPAList([]segmentgroup.SegmentGroup{})))

	d := fx.resourceData(dataSourceSegmentGroup(), "", map[string]any{"name": "missing"})
	if diags := dataSourceSegmentGroupRead(context.Background(), d, fx.client); !diags.HasError() {
		t.Fatal("expected an error for a segment group that does not exist")
	}
}
