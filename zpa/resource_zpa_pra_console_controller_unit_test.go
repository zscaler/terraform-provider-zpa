package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praconsole"
)

func TestResourcePRAConsoleController_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourcePRAConsoleController(),
		config: map[string]any{
			"name":            "ssh_console",
			"enabled":         true,
			"icon_text":       "",
			"pra_application": []any{map[string]any{"id": "pa1"}},
			"pra_portals":     []any{map[string]any{"id": []any{"pp1"}}},
		},
		id: "con1",
		response: praconsole.PRAConsole{
			ID: "con1", Name: "ssh_console", Enabled: true,
			PRAApplication: praconsole.PRAApplication{ID: "pa1"},
			PRAPortals:     []praconsole.PRAPortals{{ID: "pp1"}},
		},
		createPath: v1Path("praConsole"),
		itemPath:   v1Path("praConsole", "con1"),
		setup: func(fx *unitFixture) {
			fx.mockPolicySet("CREDENTIAL_POLICY")
		},
		wantBody: map[string]any{
			"name":           "ssh_console",
			"praApplication": map[string]any{"id": "pa1"},
			"praPortals":     []any{map[string]any{"id": "pp1"}},
		},
		wantState: map[string]any{
			"pra_application": []any{map[string]any{"id": "pa1"}},
			"pra_portals":     []any{map[string]any{"id": []any{"pp1"}}},
		},
	})
}

// Deleting a PRA console removes it from the credential policy rules that use
// it, whether the rule lists consoles one per operand (as the API returns them)
// or as a values list.
func TestResourcePRAConsoleController_DeleteDetachesFromCredentialRules(t *testing.T) {
	v1Shape := []any{
		map[string]any{"objectType": "CONSOLE", "lhs": "id", "rhs": "con1"},
		map[string]any{"objectType": "CONSOLE", "lhs": "id", "rhs": "con2"},
	}
	valuesShape := []any{map[string]any{"objectType": "CONSOLE", "lhs": "id", "values": []string{"con1", "con2"}}}
	for name, operands := range map[string][]any{"one operand per console": v1Shape, "values list": valuesShape} {
		t.Run(name, func(t *testing.T) {
			fx := newUnitFixture(t)
			ps := fx.mockPolicySet("CREDENTIAL_POLICY")
			fx.on("GET", fx.path("policySet", "rules", "policyType", "CREDENTIAL_POLICY"), common.SuccessResponse(common.ZPAList([]map[string]any{
				{"id": "r1", "name": "shared", "conditions": []any{map[string]any{"operator": "OR", "operands": operands}}},
				{"id": "r2", "name": "unrelated", "conditions": []any{map[string]any{"operator": "OR", "operands": []any{
					map[string]any{"objectType": "CONSOLE", "lhs": "id", "rhs": "con2"},
				}}}},
			})))
			fx.on("PUT", fx.pathV2("policySet", ps, "rule", "r1"), common.NoContentResponse())
			fx.on("PUT", fx.pathV2("policySet", ps, "rule", "r2"), common.NoContentResponse())
			fx.on("DELETE", fx.path("praConsole", "con1"), common.NoContentResponse())

			d := fx.resourceData(resourcePRAConsoleController(), "con1", map[string]any{"name": "ssh_console"})
			fx.run(resourcePRAConsoleControllerDelete, d)

			assertSubset(t, "shared rule", fx.body("PUT", fx.pathV2("policySet", ps, "rule", "r1")), map[string]any{
				"conditions": []any{map[string]any{"operands": []any{map[string]any{"objectType": "CONSOLE", "values": []any{"con2"}}}}},
			})
			if n := fx.calls("PUT", fx.pathV2("policySet", ps, "rule", "r2")); n != 0 {
				t.Errorf("a rule that does not use the console must not be updated, got %d", n)
			}
			if fx.calls("DELETE", fx.path("praConsole", "con1")) != 1 {
				t.Error("the console itself was not deleted")
			}
		})
	}
}
