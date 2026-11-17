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

// Some ZPA clouds return the SCIM group id as a string (issue #705); the
// lookup works either way, and large numbers such as idpId keep every digit.
func TestDataSourceScimGroup_StringID(t *testing.T) {
	fx := newUnitFixture(t)
	mockIdp(fx)
	fx.on("GET", common.ZPAUserConfigPath(common.TestCustomerID, "scimgroup", "idpId", "idp1"), common.SuccessResponse(`{"totalPages":1,"totalCount":1,"list":[
		{"id":"799580","name":"A000","idpGroupId":"428lik90rms4ih","idpId":72060623063613444,"iamIdpId":"410lik57tas4db","creationTime":1794849806,"modifiedTime":1794849806}
	]}`))

	d := fx.resourceData(dataSourceScimGroup(), "", map[string]any{"idp_name": "Okta_Users", "name": "A000"})
	fx.run(dataSourceScimGroupRead, d)
	if d.Id() != "799580" {
		t.Errorf("id = %q, want 799580", d.Id())
	}
	// The payload has no idpName, so idp_name comes from the IdP lookup.
	assertState(t, d, map[string]any{
		"idp_group_id": "428lik90rms4ih", "idp_id": 72060623063613444, "creation_time": 1794849806,
		"idp_name": "Okta_Users", "iam_idp_id": "410lik57tas4db",
	})
}

// Lookup by id when the API returns the id as a string (issue #705).
func TestDataSourceScimGroup_ByID_StringID(t *testing.T) {
	fx := newUnitFixture(t)
	mockIdp(fx)
	fx.on("GET", common.ZPAUserConfigPath(common.TestCustomerID, "scimgroup", "799580"), common.SuccessResponse(
		`{"id":"799580","name":"A000","idpGroupId":"428lik90rms4ih","idpId":72060623063613444,"iamIdpId":"410lik57tas4db"}`))

	d := fx.resourceData(dataSourceScimGroup(), "799580", map[string]any{"idp_name": "Okta_Users"})
	fx.run(dataSourceScimGroupRead, d)
	if d.Id() != "799580" {
		t.Errorf("id = %q, want 799580", d.Id())
	}
	assertState(t, d, map[string]any{"name": "A000", "idp_id": 72060623063613444, "iam_idp_id": "410lik57tas4db"})
}

// Issue #705 payload with two groups named A000 from different IAM IdPs:
// iam_idp_id and iam_idp_name each select the right group.
func TestDataSourceScimGroup_StringID_ByIamIdp(t *testing.T) {
	body := `{"totalPages":1,"totalCount":2,"list":[
		{"id":"799580","name":"A000","idpGroupId":"428lik90rms4ih","idpId":72060623063613444,"iamIdpId":"410lik57tas4db","iamIdpName":"BD_Okta_Users"},
		{"id":"799581","name":"A000","idpGroupId":"429lik90rms4ij","idpId":72060623063613444,"iamIdpId":"411lik57tas4dc","iamIdpName":"BD_Entra_Users"}
	]}`
	for _, tc := range []struct {
		name   string
		config map[string]any
		wantID string
	}{
		{"iam_idp_id", map[string]any{"idp_name": "Okta_Users", "name": "A000", "iam_idp_id": "410lik57tas4db"}, "799580"},
		{"iam_idp_name", map[string]any{"idp_name": "Okta_Users", "name": "A000", "iam_idp_name": "BD_Entra_Users"}, "799581"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fx := newUnitFixture(t)
			mockIdp(fx)
			fx.on("GET", common.ZPAUserConfigPath(common.TestCustomerID, "scimgroup", "idpId", "idp1"), common.SuccessResponse(body))

			d := fx.resourceData(dataSourceScimGroup(), "", tc.config)
			fx.run(dataSourceScimGroupRead, d)
			if d.Id() != tc.wantID {
				t.Errorf("id = %q, want %s", d.Id(), tc.wantID)
			}
		})
	}
}
