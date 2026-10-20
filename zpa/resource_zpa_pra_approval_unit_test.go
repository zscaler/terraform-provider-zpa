package zpa

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praapproval"
)

// praApprovalWindow returns RFC 1123 start and end times relative to now:
// the provider validates the window against the current time.
func praApprovalWindow(startIn, length time.Duration) (start, end time.Time) {
	start = time.Now().UTC().Add(startIn).Truncate(time.Second)
	return start, start.Add(length)
}

func praApprovalConfig(start, end time.Time) map[string]any {
	return map[string]any{
		"email_ids":    []any{"jdoe@acme.com"},
		"start_time":   start.Format(time.RFC1123),
		"end_time":     end.Format(time.RFC1123),
		"status":       "FUTURE",
		"applications": []any{map[string]any{"id": []any{"pra1"}}},
	}
}

func TestResourcePRAPrivilegedApprovalController_CRUD(t *testing.T) {
	start, end := praApprovalWindow(time.Hour, 48*time.Hour)
	startEpoch, endEpoch := fmt.Sprint(start.Unix()), fmt.Sprint(end.Unix())
	testCRUD(t, crudSpec{
		resource: resourcePRAPrivilegedApprovalController(),
		config:   praApprovalConfig(start, end),
		id:       "ap1",
		response: praapproval.PrivilegedApproval{
			ID: "ap1", EmailIDs: []string{"jdoe@acme.com"}, StartTime: startEpoch, EndTime: endEpoch, Status: "FUTURE",
			Applications: []praapproval.Applications{{ID: "pra1"}},
		},
		createPath: v1Path("approval"),
		itemPath:   v1Path("approval", "ap1"),
		// RFC 1123 times in HCL are sent as epoch seconds
		wantBody: map[string]any{
			"emailIds": []any{"jdoe@acme.com"}, "startTime": startEpoch, "endTime": endEpoch,
			"applications": []any{map[string]any{"id": "pra1"}},
		},
		wantState: map[string]any{"email_ids": []any{"jdoe@acme.com"}, "status": "FUTURE"},
	})
}

func TestResourcePRAPrivilegedApprovalController_InvalidWindowRejected(t *testing.T) {
	cases := map[string]struct {
		startIn, length time.Duration
		wantErr         string
	}{
		"start more than an hour ago": {-2 * time.Hour, time.Hour, "cannot be more than 1 hour in the past"},
		"window longer than a year":   {time.Hour, 366 * 24 * time.Hour, "max range of 1 year"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fx := newUnitFixture(t)
			start, end := praApprovalWindow(tc.startIn, tc.length)
			d := fx.resourceData(resourcePRAPrivilegedApprovalController(), "", praApprovalConfig(start, end))
			if msg := fx.runErr(resourcePRAPrivilegedApprovalControllerCreate, d); !strings.Contains(msg, tc.wantErr) {
				t.Fatalf("got %q, want an error containing %q", msg, tc.wantErr)
			}
			if n := len(fx.api.Server.AllRequests()); n != 0 {
				t.Errorf("no API call may be made, got %d", n)
			}
		})
	}
}
