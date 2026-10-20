package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

func TestDataSourceRiskScoreValues(t *testing.T) {
	path := v2Path("riskScoreValues")
	fx := newUnitFixture(t)
	fx.on("GET", path, common.SuccessResponse([]string{"LOW", "MEDIUM", "HIGH", "CRITICAL"}))

	d := fx.resourceData(dataSourceRiskScoreValues(), "", map[string]any{"exclude_unknown": true})
	fx.run(dataSourceRiskScoreValuesRead, d)
	assertState(t, d, map[string]any{"values": []any{"LOW", "MEDIUM", "HIGH", "CRITICAL"}})
	if reqs := fx.api.Server.AllRequests(); len(reqs) != 1 || !strings.Contains(reqs[0].Query, "excludeUnknown=true") {
		t.Errorf("exclude_unknown must be passed as a query parameter, got %v", reqs)
	}

	t.Run("APIError", func(t *testing.T) {
		fx := newUnitFixture(t)
		fx.on("GET", path, apiError())
		d := fx.resourceData(dataSourceRiskScoreValues(), "", map[string]any{})
		if msg := fx.runErr(dataSourceRiskScoreValuesRead, d); !strings.Contains(msg, apiErrorMarker) {
			t.Errorf("expected the API error to be returned, got %q", msg)
		}
	})
}
