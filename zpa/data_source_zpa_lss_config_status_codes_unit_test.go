package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

// Each log type gets its own status codes; the user and App Connector auth codes must not be mixed up.
func TestDataSourceLSSStatusCodes(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", lssAdminPath+"/statusCodes", common.SuccessResponse(map[string]any{
		"zpn_auth_log":     map[string]any{"ZPN_STATUS_AUTHENTICATED": map[string]any{"name": "user authenticated"}},
		"zpn_ast_auth_log": map[string]any{"ZPN_STATUS_AUTH_FAILED": map[string]any{"name": "connector auth failed"}},
		"zpn_trans_log":    map[string]any{"BRK_MT_SETUP_TIMEOUT": map[string]any{"name": "timeout"}},
		"zpn_sys_auth_log": map[string]any{"ZPN_STATUS_DISCONNECTED": map[string]any{"name": "disconnected"}},
	}))

	d := fx.resourceData(dataSourceLSSStatusCodes(), "", map[string]any{})
	fx.run(dataSourceLSSStatusCodesRead, d)
	for key, code := range map[string]string{
		"zpn_auth_log": "ZPN_STATUS_AUTHENTICATED", "zpn_ast_auth_log": "ZPN_STATUS_AUTH_FAILED",
		"zpn_trans_log": "BRK_MT_SETUP_TIMEOUT", "zpn_sys_auth_log": "ZPN_STATUS_DISCONNECTED",
	} {
		codes, _ := d.Get(key).(map[string]any)
		if _, ok := codes[code]; !ok || len(codes) != 1 {
			t.Errorf("%s = %v, want only %s", key, codes, code)
		}
	}
}
