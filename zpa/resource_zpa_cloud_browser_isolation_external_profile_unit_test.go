package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbiprofilecontroller"
)

func cbiProfileConfig() map[string]any {
	return map[string]any{
		"name":            "Isolation profile",
		"banner_id":       "b1",
		"region_ids":      []any{"r1", "r2"},
		"certificate_ids": []any{"cc1"},
		"user_experience": []any{map[string]any{"session_persistence": true, "browser_in_browser": true}},
		"security_controls": []any{map[string]any{
			"copy_paste": "all", "upload_download": "none", "document_viewer": true, "allow_printing": true,
		}},
	}
}

func cbiProfileResponse() cbiprofilecontroller.IsolationProfile {
	return cbiprofilecontroller.IsolationProfile{
		ID: "p1", Name: "Isolation profile", BannerID: "b1",
		RegionIDs: []string{"r1", "r2"}, CertificateIDs: []string{"cc1"},
		UserExperience:   &cbiprofilecontroller.UserExperience{SessionPersistence: true, BrowserInBrowser: true},
		SecurityControls: &cbiprofilecontroller.SecurityControls{CopyPaste: "all", UploadDownload: "none", DocumentViewer: true, AllowPrinting: true},
	}
}

func TestResourceCBIExternalProfile_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource:   resourceCBIExternalProfile(),
		config:     cbiProfileConfig(),
		id:         "p1",
		response:   cbiProfileResponse(),
		createPath: cbiPath("profiles"),
		itemPath:   cbiPath("profiles", "p1"),
		// import looks the profile up by name or ID
		importLookup: true,
		wantBody: map[string]any{
			"name": "Isolation profile", "bannerId": "b1",
			"regionIds": []any{"r1", "r2"}, "certificateIds": []any{"cc1"},
			"userExperience":   map[string]any{"sessionPersistence": true, "browserInBrowser": true},
			"securityControls": map[string]any{"copyPaste": "all", "uploadDownload": "none", "documentViewer": true},
		},
		wantState: map[string]any{
			"region_ids":        []any{"r1", "r2"},
			"security_controls": []any{map[string]any{"copy_paste": "all", "upload_download": "none"}},
		},
	})
}

// Create sends only the ID lists; the region, certificate and banner objects are left out.
func TestResourceCBIExternalProfile_CreateSendsIDsOnly(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("POST", cbiPath("profiles"), common.CreatedResponse(cbiProfileResponse()))
	fx.on("GET", cbiPath("profiles", "p1"), common.SuccessResponse(cbiProfileResponse()))

	d := fx.resourceData(resourceCBIExternalProfile(), "", cbiProfileConfig())
	fx.run(resourceCBIExternalProfileCreate, d)

	body := fx.body("POST", cbiPath("profiles"))
	for _, key := range []string{"regions", "certificates", "banner"} {
		if _, ok := body[key]; ok {
			t.Errorf("create request must not include %q: %v", key, body[key])
		}
	}
}

// When the API returns region and certificate objects without the ID lists,
// the IDs are taken from the objects.
func TestResourceCBIExternalProfile_ReadIDsFromObjects(t *testing.T) {
	fx := newUnitFixture(t)
	resp := cbiProfileResponse()
	resp.RegionIDs, resp.CertificateIDs, resp.BannerID = nil, nil, ""
	resp.Regions = []cbiprofilecontroller.Regions{{ID: "r1"}, {ID: "r2"}}
	resp.Certificates = []cbiprofilecontroller.Certificates{{ID: "cc1"}}
	resp.Banner = &cbiprofilecontroller.Banner{ID: "b1"}
	fx.on("GET", cbiPath("profiles", "p1"), common.SuccessResponse(resp))

	d := fx.resourceData(resourceCBIExternalProfile(), "p1", cbiProfileConfig())
	fx.run(resourceCBIExternalProfileRead, d)
	assertState(t, d, map[string]any{
		"region_ids": []any{"r1", "r2"}, "certificate_ids": []any{"cc1"}, "banner_id": "b1",
	})
}

func TestResourceCBIExternalProfile_RequiresTwoRegions(t *testing.T) {
	fx := newUnitFixture(t)
	config := cbiProfileConfig()
	config["region_ids"] = []any{"r1"}

	d := fx.resourceData(resourceCBIExternalProfile(), "", config)
	if msg := fx.runErr(resourceCBIExternalProfileCreate, d); !strings.Contains(msg, "at least 2 items") {
		t.Fatalf("expected a single region to be rejected, got %q", msg)
	}
	if n := len(fx.api.Server.AllRequests()); n != 0 {
		t.Errorf("no API call may be made, got %d", n)
	}
}
