package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredential"
)

func TestDataSourcePRACredentialController(t *testing.T) {
	match := pracredential.Credential{ID: "cr1", Name: "John Doe", CredentialType: "USERNAME_PASSWORD", UserName: "jdoe", UserDomain: "acme.com"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourcePRACredentialController(),
		listPath:  v1Path("credential"),
		itemPath:  v1Path("credential", "cr1"),
		items:     []pracredential.Credential{{ID: "other", Name: "John Doe-old"}, match},
		item:      match,
		name:      "john doe",
		id:        "cr1",
		wantState: map[string]any{"credential_type": "USERNAME_PASSWORD", "username": "jdoe"},
	})
}
