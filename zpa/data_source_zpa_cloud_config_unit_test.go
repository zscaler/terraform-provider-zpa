package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

func TestDataSourceZiaCloudConfig(t *testing.T) {
	path := v1Path("config", "ziaCloudConfig")
	fx := newUnitFixture(t)
	fx.on("GET", path, common.SuccessResponse(ziaCloudConfigResponse))
	d := fx.resourceData(dataSourceZiaCloudConfig(), "", map[string]any{})
	fx.run(dataSourceZiaCloudConfigRead, d)
	assertState(t, d, map[string]any{"zia_cloud_domain": "zscalerthree.net", "zia_username": "zpa-integration@acme.com"})

	t.Run("APIError", func(t *testing.T) {
		fx := newUnitFixture(t)
		fx.on("GET", path, apiError())
		d := fx.resourceData(dataSourceZiaCloudConfig(), "", map[string]any{})
		if msg := fx.runErr(dataSourceZiaCloudConfigRead, d); !strings.Contains(msg, apiErrorMarker) {
			t.Errorf("expected the API error to be returned, got %q", msg)
		}
	})
}
