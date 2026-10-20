package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimattributeheader"
)

// SCIM attributes are looked up within an IdP; their values come from the userconfig API.
func TestDataSourceScimAttributeHeader(t *testing.T) {
	fx := newUnitFixture(t)
	mockIdp(fx)
	attrs := []scimattributeheader.ScimAttributeHeader{
		{ID: "attr0", Name: "costCenter", IdpID: "idp1"},
		{ID: "attr1", Name: "department", IdpID: "idp1", DataType: "String", MultiValued: false},
	}
	fx.on("GET", fx.path("idp", "idp1", "scimattribute"), common.SuccessResponse(common.ZPAList(attrs)))
	fx.on("GET", common.ZPAUserConfigPath(common.TestCustomerID, "scimattribute", "idpId", "idp1", "attributeId", "attr1"),
		common.SuccessResponse(common.ZPAList([]string{"Engineering", "Sales"})))

	d := fx.resourceData(dataSourceScimAttributeHeader(), "", map[string]any{"idp_name": "Okta_Users", "name": "department"})
	fx.run(dataSourceScimAttributeHeaderRead, d)
	if d.Id() != "attr1" {
		t.Errorf("id = %q, want attr1", d.Id())
	}
	assertState(t, d, map[string]any{"data_type": "String", "values": []any{"Engineering", "Sales"}})
}
