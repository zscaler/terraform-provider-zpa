package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredentialpool"
)

// The credential pool API lives under the waap-pra-config base path.
func credentialPoolPath(parts ...string) string {
	return common.ZPAWaapPRAPath(common.TestCustomerID, "credential-pool", parts...)
}

func TestResourcePRACredentialPool_CRUD(t *testing.T) {
	pool := pracredentialpool.CredentialPool{
		ID: "cp1", Name: "Linux admins", CredentialType: "USERNAME_PASSWORD",
		PRACredentials: []zpacommon.CommonIDName{{ID: "cr1"}, {ID: "cr2"}},
	}
	testCRUD(t, crudSpec{
		resource: resourcePRACredentialPool(),
		config: map[string]any{
			"name":            "Linux admins",
			"credential_type": "USERNAME_PASSWORD",
			"credentials":     []any{map[string]any{"id": []any{"cr1", "cr2"}}},
		},
		id:         "cp1",
		response:   pool,
		createPath: credentialPoolPath(),
		itemPath:   credentialPoolPath("cp1"),
		setup: func(fx *unitFixture) {
			fx.mockPolicySet("CREDENTIAL_POLICY")
			// the create response has no ID: the SDK finds the new pool by name in the list
			fx.on("GET", credentialPoolPath(), common.SuccessResponse(common.ZPAList([]pracredentialpool.CredentialPool{
				{ID: "cp0", Name: "Linux admins (old)"}, pool,
			})))
		},
		// credentials { id = [...] } is sent as credentials[].id
		wantBody: map[string]any{
			"name": "Linux admins", "credentialType": "USERNAME_PASSWORD",
			"credentials": []any{map[string]any{"id": "cr1"}, map[string]any{"id": "cr2"}},
		},
		wantState: map[string]any{
			"credentials": []any{map[string]any{"id": []any{"cr1", "cr2"}}},
		},
	})
}
