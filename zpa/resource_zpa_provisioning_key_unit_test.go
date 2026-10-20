package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/provisioningkey"
)

func TestResourceProvisioningKey_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourceProvisioningKey(),
		config: map[string]any{
			"name":               "connector-key",
			"association_type":   "CONNECTOR_GRP",
			"enabled":            true,
			"max_usage":          "10",
			"enrollment_cert_id": "ec1",
			"zcomponent_id":      "acg1",
		},
		id: "pk1",
		response: provisioningkey.ProvisioningKey{
			ID: "pk1", Name: "connector-key", Enabled: true, MaxUsage: "10", EnrollmentCertID: "ec1",
			ZcomponentID: "acg1", ProvisioningKey: "3|api.private.zscaler.com|abc",
		},
		// the association type is part of the path
		createPath: v1Path("associationType", "CONNECTOR_GRP", "provisioningKey"),
		itemPath:   v1Path("associationType", "CONNECTOR_GRP", "provisioningKey", "pk1"),
		wantBody: map[string]any{
			"name": "connector-key", "maxUsage": "10", "enrollmentCertId": "ec1", "zcomponentId": "acg1",
		},
		wantState: map[string]any{"provisioning_key": "3|api.private.zscaler.com|abc", "zcomponent_id": "acg1"},
	})
}
