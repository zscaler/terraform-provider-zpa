package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredentialpool"
)

func TestDataSourcePRACredentialPool(t *testing.T) {
	match := pracredentialpool.CredentialPool{ID: "cp1", Name: "Linux admins", CredentialType: "USERNAME_PASSWORD", CredentialMappingCount: "2"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourcePRACredentialPool(),
		listPath:  credentialPoolPath(),
		itemPath:  credentialPoolPath("cp1"),
		items:     []pracredentialpool.CredentialPool{{ID: "other", Name: "Linux admins-old"}, match},
		item:      match,
		name:      "linux admins",
		id:        "cp1",
		wantState: map[string]any{"credential_type": "USERNAME_PASSWORD", "credential_mapping_count": "2"},
	})
}
