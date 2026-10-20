package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/lssconfigcontroller"
)

// The LSS lookup endpoints are not customer-scoped: /zpa/mgmtconfig/v2/admin/lssConfig/...
const lssAdminPath = "/zpa/mgmtconfig/v2/admin/lssConfig"

func TestDataSourceLSSClientTypes(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", lssAdminPath+"/clientTypes", common.SuccessResponse(lssconfigcontroller.LSSClientTypes{
		ZPNClientTypeExporter: "Web Browser", ZPNClientTypeZAPP: "Client Connector",
	}))

	d := fx.resourceData(dataSourceLSSClientTypes(), "", map[string]any{})
	fx.run(dataSourceLSSClientTypesRead, d)
	assertState(t, d, map[string]any{"zpn_client_type_exporter": "Web Browser", "zpn_client_type_zapp": "Client Connector"})
}
