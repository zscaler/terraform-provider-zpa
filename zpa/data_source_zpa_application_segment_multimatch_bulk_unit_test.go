package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
)

// The data source reports the segments that cannot use multimatch for the given domains.
func TestDataSourceApplicationSegmentMultimatchBulk(t *testing.T) {
	fx := newUnitFixture(t)
	path := fx.path("application", "multimatchUnsupportedReferences")
	fx.on("POST", path, common.SuccessResponse([]applicationsegment.MultiMatchUnsupportedReferencesResponse{
		{ID: "2001", AppSegmentName: "crm", Domains: []string{"crm.acme.com"}, TCPPorts: []string{"443"}, MatchStyle: "EXCLUSIVE"},
	}))

	d := fx.resourceData(dataSourceApplicationSegmentMultimatchBulk(), "", map[string]any{"domain_names": []any{"crm.acme.com"}})
	fx.run(dataSourceApplicationSegmentMultimatchBulkRead, d)
	if d.Id() == "" {
		t.Error("expected an ID derived from the domain names")
	}
	assertState(t, d, map[string]any{"unsupported_references": []any{map[string]any{
		"id": "2001", "app_segment_name": "crm", "match_style": "EXCLUSIVE",
	}}})
	reqs := fx.api.Server.AllRequests()
	if len(reqs) != 1 || !strings.Contains(string(reqs[0].Body), "crm.acme.com") {
		t.Errorf("the domain names must be sent in the request body, got %v", reqs)
	}

	t.Run("APIError", func(t *testing.T) {
		fx := newUnitFixture(t)
		fx.on("POST", path, apiError())
		d := fx.resourceData(dataSourceApplicationSegmentMultimatchBulk(), "", map[string]any{"domain_names": []any{"crm.acme.com"}})
		if msg := fx.runErr(dataSourceApplicationSegmentMultimatchBulkRead, d); !strings.Contains(msg, apiErrorMarker) {
			t.Errorf("expected the API error to be returned, got %q", msg)
		}
	})
}
