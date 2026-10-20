package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimgroup"
)

// SCIM groups are looked up by name within an IdP; their IDs are numeric.
func TestDataSourceScimGroup(t *testing.T) {
	fx := newUnitFixture(t)
	mockIdp(fx)
	fx.on("GET", common.ZPAUserConfigPath(common.TestCustomerID, "scimgroup", "idpId", "idp1"), common.SuccessResponse(common.ZPAList([]scimgroup.ScimGroup{
		{ID: 100, Name: "Engineering Contractors", IdpID: 1},
		{ID: 101, Name: "Engineering", IdpID: 1, IdpName: "Okta_Users", IdpGroupID: "g-eng"},
	})))

	d := fx.resourceData(dataSourceScimGroup(), "", map[string]any{"idp_name": "Okta_Users", "name": "engineering"})
	fx.run(dataSourceScimGroupRead, d)
	if d.Id() != "101" {
		t.Errorf("id = %q, want 101", d.Id())
	}
	assertState(t, d, map[string]any{"idp_group_id": "g-eng", "idp_name": "Okta_Users"})
}

func TestDataSourceScimGroup_ByID(t *testing.T) {
	fx := newUnitFixture(t)
	mockIdp(fx)
	fx.on("GET", common.ZPAUserConfigPath(common.TestCustomerID, "scimgroup", "101"),
		common.SuccessResponse(scimgroup.ScimGroup{ID: 101, Name: "Engineering", IdpGroupID: "g-eng"}))

	d := fx.resourceData(dataSourceScimGroup(), "101", map[string]any{"idp_name": "Okta_Users"})
	fx.run(dataSourceScimGroupRead, d)
	assertState(t, d, map[string]any{"name": "Engineering"})
}

// idp_id is numeric: a lookup by IdP ID alone must work.
func TestDataSourceScimGroup_ByIdpID(t *testing.T) {
	fx := newUnitFixture(t)
	idp := idpcontroller.IdpController{ID: "216196257331281920", Name: "Okta_Users", CreationTime: "1700000000", ModifiedTime: "1700000000"}
	fx.on("GET", fx.path("idp", idp.ID), common.SuccessResponse(idp))
	fx.on("GET", common.ZPAUserConfigPath(common.TestCustomerID, "scimgroup", "idpId", idp.ID), common.SuccessResponse(common.ZPAList([]scimgroup.ScimGroup{
		{ID: 101, Name: "Engineering", IdpGroupID: "g-eng"},
	})))

	d := fx.resourceData(dataSourceScimGroup(), "", map[string]any{"idp_id": 216196257331281920, "name": "Engineering"})
	fx.run(dataSourceScimGroupRead, d)
	if d.Id() != "101" {
		t.Errorf("id = %q, want 101", d.Id())
	}
}
