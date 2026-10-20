package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/clienttypes"
)

func TestDataSourceAccessPolicyClientTypes(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("clientTypes"), common.SuccessResponse(clienttypes.ClientTypes{
		ZPNClientTypeExplorer: "Web Browser", ZPNClientTypeZAPP: "Client Connector", ZPNClientTypeVDI: "Client Connector for VDI",
	}))

	d := fx.resourceData(dataSourceAccessPolicyClientTypes(), "", map[string]any{})
	fx.run(dataSourceAccessPolicyClientTypesRead, d)
	if d.Id() != "client_types" {
		t.Errorf("id = %q, want client_types", d.Id())
	}
	assertState(t, d, map[string]any{
		"zpn_client_type_exporter": "Web Browser", "zpn_client_type_zapp": "Client Connector",
		"zpn_client_type_vdi": "Client Connector for VDI",
	})
}
