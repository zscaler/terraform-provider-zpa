package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praapproval"
)

func praApprovalItem() praapproval.PrivilegedApproval {
	return praapproval.PrivilegedApproval{
		ID: "ap1", EmailIDs: []string{"jdoe@acme.com"}, Status: "ACTIVE",
		StartTime: "1700000000", EndTime: "1700086400",
		Applications: []praapproval.Applications{{ID: "pra1", Name: "rdp"}},
	}
}

func TestDataSourcePRAPrivilegedApprovalController_ByID(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("approval", "ap1"), common.SuccessResponse(praApprovalItem()))

	d := fx.resourceData(dataSourcePRAPrivilegedApprovalController(), "ap1", map[string]any{})
	fx.run(dataSourcePRAPrivilegedApprovalControllerRead, d)
	assertState(t, d, map[string]any{"status": "ACTIVE", "email_ids": []any{"jdoe@acme.com"}})
}

// email_ids is a set; the approval is found by the email address it is assigned to.
func TestDataSourcePRAPrivilegedApprovalController_ByEmail(t *testing.T) {
	fx := newUnitFixture(t)
	other := praApprovalItem()
	other.ID, other.EmailIDs = "ap0", []string{"someone@acme.com"}
	fx.on("GET", fx.path("approval"), common.SuccessResponse(common.ZPAList([]praapproval.PrivilegedApproval{other, praApprovalItem()})))

	d := fx.resourceData(dataSourcePRAPrivilegedApprovalController(), "", map[string]any{"email_ids": []any{"JDOE@acme.com"}})
	fx.run(dataSourcePRAPrivilegedApprovalControllerRead, d)
	if d.Id() != "ap1" {
		t.Errorf("id = %q, want ap1", d.Id())
	}
}
