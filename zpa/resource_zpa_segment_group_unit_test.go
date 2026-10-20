package zpa

import (
	"context"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
)

func TestResourceSegmentGroup_Create(t *testing.T) {
	fx := newUnitFixture(t)
	created := segmentgroup.SegmentGroup{ID: "1001", Name: "Apps", Description: "Business apps", Enabled: true}
	fx.on("POST", fx.path("segmentGroup"), common.CreatedResponse(created))
	fx.on("GET", fx.path("segmentGroup", "1001"), common.SuccessResponse(created))

	d := fx.resourceData(resourceSegmentGroup(), "", map[string]any{
		"name": "Apps", "description": "Business apps", "enabled": true,
	})
	if diags := resourceSegmentGroupCreate(context.Background(), d, fx.client); diags.HasError() {
		t.Fatalf("create failed: %v", diags)
	}

	if d.Id() != "1001" {
		t.Errorf("id = %q, want 1001", d.Id())
	}
	bodies := fx.requestBodies("POST", fx.path("segmentGroup"))
	if len(bodies) != 1 {
		t.Fatalf("expected 1 create request, got %d", len(bodies))
	}
	if bodies[0]["name"] != "Apps" || bodies[0]["enabled"] != true {
		t.Errorf("unexpected create body: %v", bodies[0])
	}
	if d.Get("description") != "Business apps" {
		t.Errorf("description = %q after read", d.Get("description"))
	}
}

// Updates must send only the group's own settings through the v2 endpoint, never
// the application list: sending the full list fails with payload.size.exceeded
// on large groups (fixed in 4.4.14).
func TestResourceSegmentGroup_UpdateSendsOnlyGroupSettings(t *testing.T) {
	fx := newUnitFixture(t)
	apps := make([]segmentgroup.Application, 0, 4000)
	for i := 0; i < 4000; i++ {
		apps = append(apps, segmentgroup.Application{ID: "app"})
	}
	current := segmentgroup.SegmentGroup{ID: "1001", Name: "Apps", Description: "new", Enabled: true, Applications: apps}
	fx.on("GET", fx.path("segmentGroup", "1001"), common.SuccessResponse(current))
	fx.on("PUT", fx.pathV2("segmentGroup", "1001"), common.NoContentResponse())

	d := fx.resourceData(resourceSegmentGroup(), "1001", map[string]any{
		"name": "Apps", "description": "new", "enabled": true,
	})
	if diags := resourceSegmentGroupUpdate(context.Background(), d, fx.client); diags.HasError() {
		t.Fatalf("update failed: %v", diags)
	}

	bodies := fx.requestBodies("PUT", fx.pathV2("segmentGroup", "1001"))
	if len(bodies) != 1 {
		t.Fatalf("expected 1 v2 update request, got %d", len(bodies))
	}
	for _, field := range []string{"applications", "addedApps", "deletedApps"} {
		if _, ok := bodies[0][field]; ok {
			t.Errorf("update body must not contain %q: %v", field, bodies[0])
		}
	}
	if bodies[0]["description"] != "new" {
		t.Errorf("description not sent: %v", bodies[0])
	}
	if n := fx.calls("PUT", fx.path("segmentGroup", "1001")); n != 0 {
		t.Errorf("the v1 update endpoint must not be used, got %d calls", n)
	}
}

// Membership is owned by zpa_application_segment.segment_group_id; adding or
// removing applications on an existing group is rejected at plan time.
func TestResourceSegmentGroup_PlanRejectsMembershipChange(t *testing.T) {
	fx := newUnitFixture(t)
	state := map[string]any{
		"name": "Apps", "enabled": true,
		"applications": []any{map[string]any{"id": "a1"}},
	}
	config := map[string]any{
		"name": "Apps", "enabled": true,
		"applications": []any{map[string]any{"id": "a1"}, map[string]any{"id": "a2"}},
	}
	err := fx.diff(resourceSegmentGroup(), "1001", state, config)
	if err == nil || !strings.Contains(err.Error(), "cannot be changed from zpa_segment_group") {
		t.Fatalf("expected the membership change to be rejected, got %v", err)
	}
}

func TestResourceSegmentGroup_PlanAllowsSameMembersInAnyOrder(t *testing.T) {
	fx := newUnitFixture(t)
	state := map[string]any{
		"name": "Apps", "enabled": true,
		"applications": []any{map[string]any{"id": "a1"}, map[string]any{"id": "a2"}},
	}
	config := map[string]any{
		"name": "Apps renamed", "enabled": true,
		"applications": []any{map[string]any{"id": "a2"}, map[string]any{"id": "a1"}},
	}
	if err := fx.diff(resourceSegmentGroup(), "1001", state, config); err != nil {
		t.Fatalf("listing the current members in another order must be allowed: %v", err)
	}
}

// Refresh reads from the paginated list endpoint and maps the API object to state.
func TestResourceSegmentGroup_ReadFromList(t *testing.T) {
	fx := newUnitFixture(t)
	groups := []segmentgroup.SegmentGroup{
		{ID: "1001", Name: "Apps", Description: "from list", Enabled: true,
			Applications: []segmentgroup.Application{{ID: "a1"}}},
		{ID: "1002", Name: "Other", Enabled: false},
	}
	fx.on("GET", fx.path("segmentGroup"), common.SuccessResponse(common.ZPAList(groups)))

	d := fx.resourceData(resourceSegmentGroup(), "1001", map[string]any{"name": "Apps"})
	if diags := resourceSegmentGroupRead(context.Background(), d, fx.client); diags.HasError() {
		t.Fatalf("read failed: %v", diags)
	}
	if d.Get("description") != "from list" || d.Get("enabled") != true {
		t.Errorf("unexpected state: description=%v enabled=%v", d.Get("description"), d.Get("enabled"))
	}
	if n := fx.calls("GET", fx.path("segmentGroup", "1001")); n != 0 {
		t.Errorf("a group found in the list must not be read by ID, got %d calls", n)
	}
}

func TestResourceSegmentGroup_ReadNotFoundRemovesFromState(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("segmentGroup"), common.SuccessResponse(common.ZPAList([]segmentgroup.SegmentGroup{})))
	fx.on("GET", fx.path("segmentGroup", "1001"), common.NotFoundResponse())

	d := fx.resourceData(resourceSegmentGroup(), "1001", map[string]any{"name": "Apps"})
	if diags := resourceSegmentGroupRead(context.Background(), d, fx.client); diags.HasError() {
		t.Fatalf("read failed: %v", diags)
	}
	if d.Id() != "" {
		t.Errorf("a deleted group must be removed from state, id = %q", d.Id())
	}
}

func TestResourceSegmentGroup_Delete(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("DELETE", fx.path("segmentGroup", "1001"), common.NoContentResponse())

	d := fx.resourceData(resourceSegmentGroup(), "1001", map[string]any{"name": "Apps"})
	if diags := resourceSegmentGroupDelete(context.Background(), d, fx.client); diags.HasError() {
		t.Fatalf("delete failed: %v", diags)
	}
	if n := fx.calls("DELETE", fx.path("segmentGroup", "1001")); n != 1 {
		t.Errorf("expected 1 delete request, got %d", n)
	}
}

// detachPolicyTypes are the policy types a segment group or application
// segment is detached from on delete.
var detachPolicyTypes = []string{"ACCESS_POLICY", "TIMEOUT_POLICY", "SIEM_POLICY", "CLIENT_FORWARDING_POLICY", "INSPECTION_POLICY"}

func TestResourceSegmentGroup_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource:   resourceSegmentGroup(),
		config:     map[string]any{"name": "Apps", "description": "Business apps", "enabled": true},
		id:         "1001",
		response:   segmentgroup.SegmentGroup{ID: "1001", Name: "Apps", Description: "Business apps", Enabled: true},
		createPath: v1Path("segmentGroup"),
		itemPath:   v1Path("segmentGroup", "1001"),
		updatePath: v2Path("segmentGroup", "1001"),
		listPath:   v1Path("segmentGroup"),
		setup: func(fx *unitFixture) {
			for _, policyType := range detachPolicyTypes {
				fx.mockPolicySet(policyType)
			}
		},
		updateConfig: func(c map[string]any) { c["description"] = "Updated" },
		wantBody:     map[string]any{"enabled": true},
		wantState:    map[string]any{"description": "Business apps"},
	})
}

// Deleting a segment group removes it from v1 rules (one operand per value) and
// v2 rules (a values list), drops conditions left empty, and leaves other rules alone.
func TestResourceSegmentGroup_DeleteDetachesFromRules(t *testing.T) {
	fx := newUnitFixture(t)
	for _, policyType := range detachPolicyTypes {
		fx.mockPolicySet(policyType)
	}
	ps := "ps-ACCESS_POLICY"
	group := func(id string) map[string]any {
		return map[string]any{"objectType": "APP_GROUP", "lhs": "id", "rhs": id, "values": []string{id}}
	}
	rules := []map[string]any{
		{"id": "r1", "name": "shared", "conditions": []any{
			map[string]any{"operator": "OR", "operands": []any{group("1001"), group("1002")}},
		}},
		{"id": "r2", "name": "only this group", "conditions": []any{
			map[string]any{"operator": "OR", "operands": []any{group("1001")}},
		}},
		{"id": "r3", "name": "unrelated", "conditions": []any{
			map[string]any{"operator": "OR", "operands": []any{group("1002")}},
		}},
	}
	fx.on("GET", fx.path("policySet", "rules", "policyType", "ACCESS_POLICY"), common.SuccessResponse(common.ZPAList(rules)))
	for _, r := range rules {
		id := r["id"].(string)
		fx.on("GET", fx.path("policySet", ps, "rule", id), common.SuccessResponse(r))
		fx.on("PUT", fx.path("policySet", ps, "rule", id), common.NoContentResponse())
		fx.on("PUT", fx.pathV2("policySet", ps, "rule", id), common.NoContentResponse())
	}
	fx.on("DELETE", fx.path("segmentGroup", "1001"), common.NoContentResponse())

	d := fx.resourceData(resourceSegmentGroup(), "1001", map[string]any{"name": "Apps"})
	fx.run(resourceSegmentGroupDelete, d)

	for _, path := range []string{fx.path("policySet", ps, "rule", "r1"), fx.pathV2("policySet", ps, "rule", "r1")} {
		for _, body := range fx.requestBodies("PUT", path) {
			ops := body["conditions"].([]any)[0].(map[string]any)["operands"].([]any)
			for _, op := range ops {
				if strings.Contains(mustJSON(t, op), "1001") {
					t.Errorf("%s still references the deleted group: %v", path, op)
				}
			}
		}
	}
	if n := fx.calls("PUT", fx.path("policySet", ps, "rule", "r1")) + fx.calls("PUT", fx.pathV2("policySet", ps, "rule", "r1")); n == 0 {
		t.Error("the shared rule was not updated")
	}
	for _, path := range []string{fx.path("policySet", ps, "rule", "r2"), fx.pathV2("policySet", ps, "rule", "r2")} {
		for _, body := range fx.requestBodies("PUT", path) {
			if conds, _ := body["conditions"].([]any); len(conds) != 0 {
				t.Errorf("%s: a condition left without operands must be dropped, got %v", path, conds)
			}
		}
	}
	if n := fx.calls("PUT", fx.path("policySet", ps, "rule", "r3")) + fx.calls("PUT", fx.pathV2("policySet", ps, "rule", "r3")); n != 0 {
		t.Errorf("a rule that does not use the group must not be updated, got %d", n)
	}
	if fx.calls("DELETE", fx.path("segmentGroup", "1001")) != 1 {
		t.Error("the segment group itself was not deleted")
	}
}
