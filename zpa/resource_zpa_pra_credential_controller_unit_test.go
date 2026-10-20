package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredential"
)

func praCredentialConfig() map[string]any {
	return map[string]any{
		"name":            "John Doe",
		"description":     "Created with Terraform",
		"credential_type": "USERNAME_PASSWORD",
		"user_domain":     "acme.com",
		"username":        "jdoe",
		"password":        "s3cret",
	}
}

func TestResourcePRACredentialController_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourcePRACredentialController(),
		config:   praCredentialConfig(),
		id:       "cr1",
		response: pracredential.Credential{
			ID: "cr1", Name: "John Doe", Description: "Created with Terraform",
			CredentialType: "USERNAME_PASSWORD", UserDomain: "acme.com", UserName: "jdoe",
		},
		createPath: v1Path("credential"),
		itemPath:   v1Path("credential", "cr1"),
		setup: func(fx *unitFixture) {
			fx.mockPolicySet("CREDENTIAL_POLICY")
		},
		updateConfig: func(c map[string]any) { c["description"] = "Rotated" },
		// the mock API redacts secrets in recorded bodies, so the password value is not compared
		wantBody: map[string]any{
			"credentialType": "USERNAME_PASSWORD", "userName": "jdoe", "userDomain": "acme.com",
		},
		wantState: map[string]any{"credential_type": "USERNAME_PASSWORD", "username": "jdoe"},
	})
}

func TestResourcePRACredentialController_CredentialTypeChangeRejected(t *testing.T) {
	fx := newUnitFixture(t)
	config := praCredentialConfig()
	config["credential_type"] = "PASSWORD"

	d := fx.updateData(resourcePRACredentialController(), "cr1", praCredentialConfig(), config)
	if msg := fx.runErr(resourcePRACredentialControllerUpdate, d); !strings.Contains(msg, "changing 'credential_type'") {
		t.Fatalf("expected a credential_type change to be rejected, got %q", msg)
	}
	if n := len(fx.api.Server.AllRequests()); n != 0 {
		t.Errorf("no API call may be made for a rejected update, got %d", n)
	}
}
