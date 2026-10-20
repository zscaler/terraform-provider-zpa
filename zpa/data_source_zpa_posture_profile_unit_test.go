package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/postureprofile"
)

func TestDataSourcePostureProfile(t *testing.T) {
	match := postureprofile.PostureProfile{ID: "pp1", Name: "CrowdStrike_ZPA_Pre-ZTA", Domain: "acme.com", PostureudID: "udid-1"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourcePostureProfile(),
		listPath:  v2Path("posture"),
		itemPath:  v1Path("posture", "pp1"),
		items:     []postureprofile.PostureProfile{{ID: "other", Name: "CrowdStrike_ZPA_Pre-ZTA-old"}, match},
		item:      match,
		name:      "crowdstrike_zpa_pre-zta",
		id:        "pp1",
		wantState: map[string]any{"domain": "acme.com", "posture_udid": "udid-1"},
	})
}
