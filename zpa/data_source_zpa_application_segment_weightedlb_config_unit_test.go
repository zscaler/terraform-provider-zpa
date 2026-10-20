package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

func TestDataSourceApplicationSegmentWeightedLBConfig(t *testing.T) {
	cases := map[string]map[string]any{
		"by application id":   {"application_id": "2001"},
		"by application name": {"application_name": "crm"},
	}
	for name, config := range cases {
		t.Run(name, func(t *testing.T) {
			fx := newUnitFixture(t)
			mockWeightedLB(fx)
			d := fx.resourceData(dataSourceApplicationSegmentWeightedLBConfig(), "", config)
			fx.run(dataSourceApplicationSegmentWeightedLBConfigRead, d)
			if d.Id() != "2001" {
				t.Errorf("id = %q, want 2001", d.Id())
			}
			assertState(t, d, map[string]any{
				"weighted_load_balancing": true,
				"application_to_server_group_mappings": []any{
					map[string]any{"id": "sg1", "weight": "10"},
					map[string]any{"id": "sg2", "passive": true},
				},
			})
		})
	}

	t.Run("RequiresApplication", func(t *testing.T) {
		fx := newUnitFixture(t)
		d := fx.resourceData(dataSourceApplicationSegmentWeightedLBConfig(), "", map[string]any{})
		if msg := fx.runErr(dataSourceApplicationSegmentWeightedLBConfigRead, d); !strings.Contains(msg, "must be provided") {
			t.Errorf("got %q", msg)
		}
	})

	t.Run("APIError", func(t *testing.T) {
		fx := newUnitFixture(t)
		fx.on("GET", weightedLBPath, apiError())
		d := fx.resourceData(dataSourceApplicationSegmentWeightedLBConfig(), "", map[string]any{"application_id": "2001"})
		if msg := fx.runErr(dataSourceApplicationSegmentWeightedLBConfigRead, d); !strings.Contains(msg, apiErrorMarker) {
			t.Errorf("expected the API error to be returned, got %q", msg)
		}
	})

	// in a microtenant the ID is scoped to it
	t.Run("MicrotenantID", func(t *testing.T) {
		fx := newUnitFixture(t)
		fx.on("GET", weightedLBPath, common.SuccessResponse(weightedLBResponse()))
		d := fx.resourceData(dataSourceApplicationSegmentWeightedLBConfig(), "", map[string]any{"application_id": "2001", "microtenant_id": "mt1"})
		fx.run(dataSourceApplicationSegmentWeightedLBConfigRead, d)
		if d.Id() != "mt1:2001" {
			t.Errorf("id = %q, want mt1:2001", d.Id())
		}
	})
}
