package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/provisioningkey"
)

func TestDataSourceProvisioningKey(t *testing.T) {
	match := provisioningkey.ProvisioningKey{ID: "pk1", Name: "connector-key", MaxUsage: "10", ZcomponentID: "acg1", ProvisioningKey: "3|api.private.zscaler.com|abc"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceProvisioningKey(),
		listPath:  v1Path("associationType", "CONNECTOR_GRP", "provisioningKey"),
		itemPath:  v1Path("associationType", "CONNECTOR_GRP", "provisioningKey", "pk1"),
		items:     []provisioningkey.ProvisioningKey{{ID: "pk0", Name: "connector-key-old"}, match},
		item:      match,
		name:      "Connector-Key",
		id:        "pk1",
		extra:     map[string]any{"association_type": "CONNECTOR_GRP"},
		wantState: map[string]any{"max_usage": "10", "zcomponent_id": "acg1"},
	})
}

func TestDataSourceProvisioningKey_RequiresAssociationType(t *testing.T) {
	fx := newUnitFixture(t)
	d := fx.resourceData(dataSourceProvisioningKey(), "", map[string]any{"name": "connector-key"})
	if msg := fx.runErr(dataSourceProvisioningKeyRead, d); !strings.Contains(msg, "associationType is required") {
		t.Fatalf("got %q, want an association type error", msg)
	}
}
