package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/samlattribute"
)

// mockIdp registers IdP idp1 ("Okta_Users") for lookups by ID and by name.
func mockIdp(fx *unitFixture) {
	idp := idpcontroller.IdpController{ID: "idp1", Name: "Okta_Users", CreationTime: "1700000000", ModifiedTime: "1700000000"}
	fx.on("GET", fx.path("idp", "idp1"), common.SuccessResponse(idp))
	fx.on("GET", fx.pathV2("idp"), common.SuccessResponse(common.ZPAList([]idpcontroller.IdpController{idp})))
}

var samlAttributes = []samlattribute.SamlAttribute{
	{ID: "sa0", Name: "Email_Okta_Users_Old", IdpID: "idp1"},
	{ID: "sa1", Name: "Email_Okta_Users", SamlName: "email", IdpID: "idp1", IdpName: "Okta_Users", UserAttribute: true},
}

// SAML attributes are looked up within an IdP, given by name or ID.
func TestDataSourceSamlAttribute(t *testing.T) {
	cases := map[string]map[string]any{
		"by name in idp by name": {"idp_name": "Okta_Users", "name": "Email_Okta_Users"},
		"by name in idp by id":   {"idp_id": "idp1", "name": "Email_Okta_Users"},
	}
	for name, config := range cases {
		t.Run(name, func(t *testing.T) {
			fx := newUnitFixture(t)
			mockIdp(fx)
			fx.on("GET", fx.pathV2("samlAttribute", "idp", "idp1"), common.SuccessResponse(common.ZPAList(samlAttributes)))

			d := fx.resourceData(dataSourceSamlAttribute(), "", config)
			fx.run(dataSourceSamlAttributeRead, d)
			if d.Id() != "sa1" {
				t.Errorf("id = %q, want sa1", d.Id())
			}
			assertState(t, d, map[string]any{"saml_name": "email", "user_attribute": true})
		})
	}
}

func TestDataSourceSamlAttribute_ByID(t *testing.T) {
	fx := newUnitFixture(t)
	mockIdp(fx)
	fx.on("GET", fx.pathV2("samlAttribute", "idp", "idp1"), common.SuccessResponse(common.ZPAList(samlAttributes)))

	d := fx.resourceData(dataSourceSamlAttribute(), "sa1", map[string]any{"idp_id": "idp1"})
	fx.run(dataSourceSamlAttributeRead, d)
	assertState(t, d, map[string]any{"name": "Email_Okta_Users"})
}

func TestDataSourceSamlAttribute_RequiresIdp(t *testing.T) {
	fx := newUnitFixture(t)
	d := fx.resourceData(dataSourceSamlAttribute(), "", map[string]any{"name": "Email_Okta_Users"})
	if msg := fx.runErr(dataSourceSamlAttributeRead, d); !strings.Contains(msg, "idp name or id is required") {
		t.Fatalf("got %q, want an IdP required error", msg)
	}
}
