package zpa

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/enrollmentcert"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

func appConnectorGroupConfig() map[string]any {
	return map[string]any{
		"name":                     "SJC Connectors",
		"enabled":                  true,
		"country_code":             "US",
		"city_country":             "San Jose, US",
		"latitude":                 "37.3382082",
		"longitude":                "-121.8863286",
		"location":                 "San Jose, CA, USA",
		"upgrade_day":              "SUNDAY",
		"upgrade_time_in_secs":     "66600",
		"override_version_profile": false,
		"dns_query_type":           "IPV4_IPV6",
		"enrollment_cert_id":       "ec1",
	}
}

func appConnectorGroupResponse() appconnectorgroup.AppConnectorGroup {
	return appconnectorgroup.AppConnectorGroup{
		ID: "acg1", Name: "SJC Connectors", Enabled: true, CountryCode: "US", CityCountry: "San Jose, US",
		Latitude: "37.3382082", Longitude: "-121.8863286", Location: "San Jose, CA, USA",
		UpgradeDay: "SUNDAY", UpgradeTimeInSecs: "66600", DNSQueryType: "IPV4_IPV6", EnrollmentCertID: "ec1",
	}
}

func TestResourceAppConnectorGroup_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource:   resourceAppConnectorGroup(),
		config:     appConnectorGroupConfig(),
		id:         "acg1",
		response:   appConnectorGroupResponse(),
		createPath: v1Path("appConnectorGroup"),
		itemPath:   v1Path("appConnectorGroup", "acg1"),
		listPath:   v1Path("appConnectorGroup"),
		setup: func(fx *unitFixture) {
			fx.mockPolicySet("ACCESS_POLICY")
		},
		wantBody: map[string]any{
			"name": "SJC Connectors", "latitude": "37.3382082", "longitude": "-121.8863286",
			"upgradeDay": "SUNDAY", "dnsQueryType": "IPV4_IPV6", "enrollmentCertId": "ec1",
		},
		wantState: map[string]any{"location": "San Jose, CA, USA", "enrollment_cert_id": "ec1"},
	})
}

// Groups are always read from the list: the single-group GET returns a wrong
// lssAppConnectorGroup value.
func TestResourceAppConnectorGroup_ReadUsesListNotGet(t *testing.T) {
	fx := newUnitFixture(t)
	listed := appConnectorGroupResponse()
	listed.LSSAppConnectorGroup = true
	fx.on("GET", fx.path("appConnectorGroup"), common.SuccessResponse(common.ZPAList([]any{listed})))
	fx.on("GET", fx.path("appConnectorGroup", "acg1"), common.SuccessResponse(appConnectorGroupResponse()))

	d := fx.resourceData(resourceAppConnectorGroup(), "acg1", appConnectorGroupConfig())
	fx.run(resourceAppConnectorGroupRead, d)
	assertState(t, d, map[string]any{"lss_app_connector_group": true})
	if n := fx.calls("GET", fx.path("appConnectorGroup", "acg1")); n != 0 {
		t.Errorf("the single-group GET must not be used, got %d calls", n)
	}
}

// Without enrollment_cert_id the "Connector" enrollment certificate is looked up and used.
func TestResourceAppConnectorGroup_EnrollmentCertAutoResolved(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.pathV2("enrollmentCert"), common.SuccessResponse(common.ZPAList([]enrollmentcert.EnrollmentCert{
		{ID: "ec0", Name: "Connector Legacy"}, {ID: "ec9", Name: "Connector"},
	})))
	fx.on("POST", fx.path("appConnectorGroup"), common.CreatedResponse(appConnectorGroupResponse()))
	fx.on("GET", fx.path("appConnectorGroup"), common.SuccessResponse(common.ZPAList([]any{appConnectorGroupResponse()})))

	config := appConnectorGroupConfig()
	delete(config, "enrollment_cert_id")
	d := fx.resourceData(resourceAppConnectorGroup(), "", config)
	fx.run(resourceAppConnectorGroupCreate, d)

	assertSubset(t, "create request", fx.body("POST", fx.path("appConnectorGroup")), map[string]any{"enrollmentCertId": "ec9"})
}

func TestResourceAppConnectorGroup_TCPQuickAckMustBeConsistent(t *testing.T) {
	fx := newUnitFixture(t)
	config := appConnectorGroupConfig()
	config["tcp_quick_ack_app"] = true
	config["tcp_quick_ack_assistant"] = false

	d := fx.resourceData(resourceAppConnectorGroup(), "", config)
	if msg := fx.runErr(resourceAppConnectorGroupCreate, d); !strings.Contains(msg, "tcpQuickAck") {
		t.Fatalf("expected inconsistent tcp_quick_ack flags to be rejected, got %q", msg)
	}
	if n := fx.calls("POST", fx.path("appConnectorGroup")); n != 0 {
		t.Errorf("no create request may be sent, got %d", n)
	}
}

// user_codes enroll the deployed connectors into the new group.
func TestResourceAppConnectorGroup_UserCodesVerified(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("POST", fx.path("appConnectorGroup"), common.CreatedResponse(appConnectorGroupResponse()))
	fx.on("GET", fx.path("appConnectorGroup"), common.SuccessResponse(common.ZPAList([]any{appConnectorGroupResponse()})))
	fx.on("POST", fx.path("CONNECTOR_GRP", "usercodes"), common.SuccessResponse(map[string]any{"componentGroupId": "acg1"}))

	config := appConnectorGroupConfig()
	config["user_codes"] = []any{"ABC123"}
	d := fx.resourceData(resourceAppConnectorGroup(), "", config)
	fx.run(resourceAppConnectorGroupCreate, d)

	assertSubset(t, "user code request", fx.body("POST", fx.path("CONNECTOR_GRP", "usercodes")),
		map[string]any{"componentGroupId": "acg1", "user_codes": []any{"ABC123"}})
}

// Deleting a group first removes it from the access policy rules that use it.
func TestResourceAppConnectorGroup_DeleteDetachesFromAccessRules(t *testing.T) {
	fx := newUnitFixture(t)
	ps := fx.mockPolicySet("ACCESS_POLICY")
	rule := policysetcontroller.PolicyRule{
		ID: "r1", Name: "Allow", PolicySetID: ps,
		AppConnectorGroups: []appconnectorgroup.AppConnectorGroup{{ID: "acg1"}, {ID: "acg2"}},
	}
	fx.on("GET", fx.path("policySet", "rules", "policyType", "ACCESS_POLICY"),
		common.SuccessResponse(common.ZPAList([]policysetcontroller.PolicyRule{rule})))
	fx.on("GET", fx.path("policySet", ps, "rule", "r1"), common.SuccessResponse(rule))
	fx.on("PUT", fx.path("policySet", ps, "rule", "r1"), common.NoContentResponse())
	fx.on("PUT", fx.pathV2("policySet", ps, "rule", "r1"), common.NoContentResponse())
	fx.on("DELETE", fx.path("appConnectorGroup", "acg1"), common.NoContentResponse())

	d := fx.resourceData(resourceAppConnectorGroup(), "acg1", appConnectorGroupConfig())
	fx.run(resourceAppConnectorGroupDelete, d)

	assertSubset(t, "rule update", fx.body("PUT", fx.path("policySet", ps, "rule", "r1")),
		map[string]any{"appConnectorGroups": []any{map[string]any{"id": "acg2"}}})
	// v2 rules are updated through the v2 API
	assertSubset(t, "v2 rule update", fx.body("PUT", fx.pathV2("policySet", ps, "rule", "r1")),
		map[string]any{"appConnectorGroups": []any{map[string]any{"id": "acg2"}}})
	if fx.calls("DELETE", fx.path("appConnectorGroup", "acg1")) != 1 {
		t.Error("the group itself was not deleted")
	}
}

// Coordinates and the country code are checked at plan time, for the app
// connector, service edge and private cloud groups alike.
func TestConnectorGroupLocationValidation(t *testing.T) {
	resources := map[string]func() map[string]any{
		"app connector group": appConnectorGroupConfig,
		"service edge group":  serviceEdgeGroupConfig,
	}
	cases := map[string]struct {
		field, value, wantErr string
	}{
		"latitude":               {"latitude", "37.3", ""},
		"latitude out of range":  {"latitude", "95", "latitude must be between -90 and 90"},
		"latitude not a number":  {"latitude", "north", "latitude must be between -90 and 90"},
		"longitude out of range": {"longitude", "-200", "longitude must be between -180 and 180"},
		"longitude not a number": {"longitude", "west", "longitude must be between -180 and 180"},
		"country code":           {"country_code", "BR", ""},
		"country code invalid":   {"country_code", "XX", "XX"},
	}
	for resourceName, config := range resources {
		r := map[string]func() *schema.Resource{
			"app connector group": resourceAppConnectorGroup,
			"service edge group":  resourceServiceEdgeGroup,
		}[resourceName]()
		for name, tc := range cases {
			t.Run(resourceName+"/"+name, func(t *testing.T) {
				c := config()
				c[tc.field] = tc.value
				got := validate(r, c)
				if tc.wantErr == "" && got != "" || tc.wantErr != "" && !strings.Contains(got, tc.wantErr) {
					t.Errorf("got %q, want %q", got, tc.wantErr)
				}
			})
		}
	}
}
