package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
)

var multimatchPath = v1Path("application", "bulkUpdateMultiMatch")

// mockMultimatch serves two segments; refs are the match styles the references endpoint reports.
func mockMultimatch(fx *unitFixture, styles map[string]string) {
	var refs []applicationsegment.MultiMatchUnsupportedReferencesResponse
	for _, id := range []string{"2001", "2002"} {
		fx.on("GET", fx.path("application", id), common.SuccessResponse(applicationsegment.ApplicationSegmentResource{
			ID: id, DomainNames: []string{id + ".acme.com"}, MatchStyle: styles[id],
		}))
		refs = append(refs, applicationsegment.MultiMatchUnsupportedReferencesResponse{ID: id, MatchStyle: styles[id]})
	}
	fx.on("POST", fx.path("application", "multimatchUnsupportedReferences"), common.SuccessResponse(refs))
	fx.on("PUT", multimatchPath, common.NoContentResponse())
}

func multimatchConfig(style string) map[string]any {
	return map[string]any{"application_ids": []any{"2001", "2002"}, "match_style": style}
}

// Application IDs are sent as numbers; the resource ID is derived from IDs and match style.
func TestResourceApplicationSegmentMultimatchBulk_Create(t *testing.T) {
	fx := newUnitFixture(t)
	mockMultimatch(fx, map[string]string{"2001": "INCLUSIVE", "2002": "INCLUSIVE"})

	d := fx.resourceData(resourceApplicationSegmentMultimatchBulk(), "", multimatchConfig("INCLUSIVE"))
	fx.run(resourceApplicationSegmentMultimatchBulkCreate, d)
	assertSubset(t, "bulk request", fx.body("PUT", multimatchPath), map[string]any{
		"applicationIds": []any{2001.0, 2002.0}, "matchStyle": "INCLUSIVE",
	})
	if d.Id() != generateBulkMultimatchID([]string{"2002", "2001"}, "INCLUSIVE") {
		t.Errorf("id %q is not derived from the IDs (in any order) and match style", d.Id())
	}
}

func TestResourceApplicationSegmentMultimatchBulk_NonNumericIDRejected(t *testing.T) {
	fx := newUnitFixture(t)
	config := multimatchConfig("INCLUSIVE")
	config["application_ids"] = []any{"abc"}
	d := fx.resourceData(resourceApplicationSegmentMultimatchBulk(), "", config)
	if msg := fx.runErr(resourceApplicationSegmentMultimatchBulkCreate, d); !strings.Contains(msg, "failed to convert application_id") {
		t.Fatalf("got %q", msg)
	}
}

// A segment whose match style was changed outside Terraform shows up in state,
// whether it is reported by the references endpoint or read per segment.
func TestResourceApplicationSegmentMultimatchBulk_ReadDetectsDrift(t *testing.T) {
	for name, references := range map[string]bool{"from references": true, "from segments": false} {
		t.Run(name, func(t *testing.T) {
			fx := newUnitFixture(t)
			mockMultimatch(fx, map[string]string{"2001": "INCLUSIVE", "2002": "EXCLUSIVE"})
			if !references {
				fx.on("POST", fx.path("application", "multimatchUnsupportedReferences"), apiError())
			}
			d := fx.resourceData(resourceApplicationSegmentMultimatchBulk(), "id", multimatchConfig("INCLUSIVE"))
			fx.run(resourceApplicationSegmentMultimatchBulkRead, d)
			assertState(t, d, map[string]any{"match_style": "EXCLUSIVE"})
		})
	}
}

func TestResourceApplicationSegmentMultimatchBulk_ReadInSync(t *testing.T) {
	fx := newUnitFixture(t)
	mockMultimatch(fx, map[string]string{"2001": "INCLUSIVE", "2002": "INCLUSIVE"})
	d := fx.resourceData(resourceApplicationSegmentMultimatchBulk(), "id", multimatchConfig("INCLUSIVE"))
	fx.run(resourceApplicationSegmentMultimatchBulkRead, d)
	assertState(t, d, map[string]any{"match_style": "INCLUSIVE"})
}

func TestResourceApplicationSegmentMultimatchBulk_Update(t *testing.T) {
	fx := newUnitFixture(t)
	mockMultimatch(fx, map[string]string{"2001": "EXCLUSIVE", "2002": "EXCLUSIVE"})
	d := fx.updateData(resourceApplicationSegmentMultimatchBulk(), "old", multimatchConfig("INCLUSIVE"), multimatchConfig("EXCLUSIVE"))
	fx.run(resourceApplicationSegmentMultimatchBulkUpdate, d)
	assertSubset(t, "bulk request", fx.body("PUT", multimatchPath), map[string]any{"matchStyle": "EXCLUSIVE"})
	if d.Id() == "old" {
		t.Error("the ID must follow the new match style")
	}
}

func TestResourceApplicationSegmentMultimatchBulk_APIError(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("PUT", multimatchPath, apiError())
	d := fx.resourceData(resourceApplicationSegmentMultimatchBulk(), "", multimatchConfig("INCLUSIVE"))
	if msg := fx.runErr(resourceApplicationSegmentMultimatchBulkCreate, d); !strings.Contains(msg, apiErrorMarker) {
		t.Errorf("expected the API error to be returned, got %q", msg)
	}
	if d.Id() != "" {
		t.Errorf("a failed create must not store an ID, got %q", d.Id())
	}
}
