package zpa

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

func appSegmentConfig() map[string]any {
	return map[string]any{
		"name":             "crm",
		"enabled":          true,
		"domain_names":     []any{"crm.example.com"},
		"segment_group_id": "1001",
		"tcp_port_range":   []any{map[string]any{"from": "8443", "to": "8443"}},
		"server_groups":    []any{map[string]any{"id": []any{"3001"}}},
		"tags": []any{
			map[string]any{"tag_key_id": "k1", "tag_value_id": "v1"},
		},
	}
}

func appSegmentResponse() applicationsegment.ApplicationSegmentResource {
	return applicationsegment.ApplicationSegmentResource{
		ID: "2001", Name: "crm", Enabled: true, SegmentGroupID: "1001",
		DomainNames: []string{"crm.example.com"}, TCPPortRanges: []string{"8443", "8443"},
	}
}

func TestResourceApplicationSegment_CreateRequest(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("POST", fx.path("application"), common.CreatedResponse(appSegmentResponse()))
	fx.on("GET", fx.path("application", "2001"), common.SuccessResponse(appSegmentResponse()))

	d := fx.resourceData(resourceApplicationSegment(), "", appSegmentConfig())
	if diags := resourceApplicationSegmentCreate(context.Background(), d, fx.client); diags.HasError() {
		t.Fatalf("create failed: %v", diags)
	}
	if d.Id() != "2001" {
		t.Errorf("id = %q, want 2001", d.Id())
	}

	bodies := fx.requestBodies("POST", fx.path("application"))
	if len(bodies) != 1 {
		t.Fatalf("expected 1 create request, got %d", len(bodies))
	}
	body := bodies[0]
	if body["segmentGroupId"] != "1001" {
		t.Errorf("segmentGroupId = %v", body["segmentGroupId"])
	}
	if got := body["tcpPortRanges"]; !reflect.DeepEqual(got, []any{"8443", "8443"}) {
		t.Errorf("tcpPortRanges = %v, want [8443 8443]", got)
	}
	groups, _ := body["serverGroups"].([]any)
	if len(groups) != 1 || groups[0].(map[string]any)["id"] != "3001" {
		t.Errorf("serverGroups = %v", body["serverGroups"])
	}
	// tags are sent as tagKey.id / tagValue.id, the shape the ZPA API accepts.
	tags, _ := body["tags"].([]any)
	if len(tags) != 1 {
		t.Fatalf("tags = %v, want one tag", body["tags"])
	}
	tag := tags[0].(map[string]any)
	if tag["tagKey"].(map[string]any)["id"] != "k1" || tag["tagValue"].(map[string]any)["id"] != "v1" {
		t.Errorf("unexpected tag: %v", tag)
	}
}

func TestResourceApplicationSegment_CreateWithoutSegmentGroupFailsBeforeAPI(t *testing.T) {
	fx := newUnitFixture(t)
	config := appSegmentConfig()
	delete(config, "segment_group_id")

	d := fx.resourceData(resourceApplicationSegment(), "", config)
	diags := resourceApplicationSegmentCreate(context.Background(), d, fx.client)
	if !diags.HasError() {
		t.Fatal("expected an error without segment_group_id")
	}
	if n := fx.calls("POST", fx.path("application")); n != 0 {
		t.Errorf("no create request may be sent without a segment group, got %d", n)
	}
}

func TestResourceApplicationSegment_UpdateUsesPut(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("application", "2001"), common.SuccessResponse(appSegmentResponse()))
	fx.on("PUT", fx.path("application", "2001"), common.NoContentResponse())

	config := appSegmentConfig()
	config["domain_names"] = []any{"crm.example.com", "crm2.example.com"}
	d := fx.resourceData(resourceApplicationSegment(), "2001", config)
	if diags := resourceApplicationSegmentUpdate(context.Background(), d, fx.client); diags.HasError() {
		t.Fatalf("update failed: %v", diags)
	}

	bodies := fx.requestBodies("PUT", fx.path("application", "2001"))
	if len(bodies) != 1 {
		t.Fatalf("expected 1 update request, got %d", len(bodies))
	}
	if got := bodies[0]["domainNames"]; !reflect.DeepEqual(got, []any{"crm.example.com", "crm2.example.com"}) {
		t.Errorf("domainNames = %v", got)
	}
}

func TestResourceApplicationSegment_ReadNotFoundRemovesFromState(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("application"), common.SuccessResponse(common.ZPAList([]applicationsegment.ApplicationSegmentResource{})))
	fx.on("GET", fx.path("application", "2001"), common.NotFoundResponse())

	d := fx.resourceData(resourceApplicationSegment(), "2001", appSegmentConfig())
	if diags := resourceApplicationSegmentRead(context.Background(), d, fx.client); diags.HasError() {
		t.Fatalf("read failed: %v", diags)
	}
	if d.Id() != "" {
		t.Errorf("a deleted segment must be removed from state, id = %q", d.Id())
	}
}

func TestResourceApplicationSegment_ReadMapsPortsAndDomains(t *testing.T) {
	fx := newUnitFixture(t)
	resp := appSegmentResponse()
	resp.DomainNames = []string{"crm.example.com", "api.example.com"}
	fx.on("GET", fx.path("application"), common.SuccessResponse(common.ZPAList([]applicationsegment.ApplicationSegmentResource{resp})))

	d := fx.resourceData(resourceApplicationSegment(), "2001", map[string]any{"name": "crm"})
	if diags := resourceApplicationSegmentRead(context.Background(), d, fx.client); diags.HasError() {
		t.Fatalf("read failed: %v", diags)
	}
	if got := d.Get("domain_names").(*schema.Set).Len(); got != 2 {
		t.Errorf("domain_names has %d entries, want 2", got)
	}
	if d.Get("segment_group_id") != "1001" {
		t.Errorf("segment_group_id = %v", d.Get("segment_group_id"))
	}
}

func TestResourceApplicationSegment_Delete(t *testing.T) {
	fx := newUnitFixture(t)
	for _, policyType := range []string{"ACCESS_POLICY", "TIMEOUT_POLICY", "SIEM_POLICY", "CLIENT_FORWARDING_POLICY", "INSPECTION_POLICY"} {
		fx.mockPolicySet(policyType)
	}
	fx.on("DELETE", fx.path("application", "2001"), common.NoContentResponse())

	d := fx.resourceData(resourceApplicationSegment(), "2001", appSegmentConfig())
	fx.run(resourceApplicationSegmentDelete, d)
	if n := fx.calls("DELETE", fx.path("application", "2001")); n != 1 {
		t.Errorf("expected 1 delete request, got %d", n)
	}
}

func TestResourceApplicationSegment_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource:   resourceApplicationSegment(),
		config:     appSegmentConfig(),
		id:         "2001",
		response:   appSegmentResponse(),
		createPath: v1Path("application"),
		itemPath:   v1Path("application", "2001"),
		listPath:   v1Path("application"),
		setup: func(fx *unitFixture) {
			for _, policyType := range detachPolicyTypes {
				fx.mockPolicySet(policyType)
			}
		},
		updateConfig: func(c map[string]any) { c["description"] = "Updated" },
		wantBody:     map[string]any{"segmentGroupId": "1001", "domainNames": []any{"crm.example.com"}},
		wantState:    map[string]any{"segment_group_id": "1001"},
	})
}

// Deleting a segment removes it from the policy rules that reference it; a
// condition left without operands is dropped, and other rules are not touched.
func TestResourceApplicationSegment_DeleteDetachesFromRules(t *testing.T) {
	fx := newUnitFixture(t)
	for _, policyType := range detachPolicyTypes {
		fx.mockPolicySet(policyType)
	}
	app := func(id string) policysetcontroller.Operands {
		return policysetcontroller.Operands{ObjectType: "APP", LHS: "id", RHS: id}
	}
	fx.on("GET", fx.path("policySet", "rules", "policyType", "ACCESS_POLICY"), common.SuccessResponse(common.ZPAList([]policysetcontroller.PolicyRule{
		{ID: "r1", Name: "shared", Conditions: []policysetcontroller.Conditions{
			{Operator: "OR", Operands: []policysetcontroller.Operands{app("2001"), app("2002")}},
		}},
		{ID: "r2", Name: "only this segment", Conditions: []policysetcontroller.Conditions{
			{Operator: "OR", Operands: []policysetcontroller.Operands{app("2001")}},
			{Operator: "OR", Operands: []policysetcontroller.Operands{{ObjectType: "CLIENT_TYPE", LHS: "id", RHS: "zpn_client_type_zapp"}}},
		}},
		{ID: "r3", Name: "unrelated", Conditions: []policysetcontroller.Conditions{
			{Operator: "OR", Operands: []policysetcontroller.Operands{app("2002")}},
		}},
	})))
	for _, id := range []string{"r1", "r2", "r3"} {
		fx.on("PUT", fx.path("policySet", "ps-ACCESS_POLICY", "rule", id), common.NoContentResponse())
	}
	fx.on("DELETE", fx.path("application", "2001"), common.NoContentResponse())

	d := fx.resourceData(resourceApplicationSegment(), "2001", appSegmentConfig())
	fx.run(resourceApplicationSegmentDelete, d)

	assertSubset(t, "shared rule", fx.body("PUT", fx.path("policySet", "ps-ACCESS_POLICY", "rule", "r1")), map[string]any{
		"conditions": []any{map[string]any{"operands": []any{map[string]any{"rhs": "2002"}}}},
	})
	assertSubset(t, "rule that only had this segment", fx.body("PUT", fx.path("policySet", "ps-ACCESS_POLICY", "rule", "r2")), map[string]any{
		"conditions": []any{map[string]any{"operands": []any{map[string]any{"objectType": "CLIENT_TYPE"}}}},
	})
	if n := fx.calls("PUT", fx.path("policySet", "ps-ACCESS_POLICY", "rule", "r3")); n != 0 {
		t.Errorf("a rule that does not use the segment must not be updated, got %d", n)
	}
	if fx.calls("DELETE", fx.path("application", "2001")) != 1 {
		t.Error("the segment itself was not deleted")
	}
}
