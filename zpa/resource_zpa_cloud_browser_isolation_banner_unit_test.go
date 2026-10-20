package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbibannercontroller"
)

func TestResourceCBIBanner_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourceCBIBanners(),
		config: map[string]any{
			"name":               "Corp Banner",
			"primary_color":      "#0076BE",
			"text_color":         "#FFFFFF",
			"notification_title": "Heads up",
			"notification_text":  "Isolated browsing",
			"banner":             true,
			"persist":            true,
		},
		id: "b1",
		response: cbibannercontroller.CBIBannerController{
			ID: "b1", Name: "Corp Banner", PrimaryColor: "#0076BE", TextColor: "#FFFFFF",
			NotificationTitle: "Heads up", NotificationText: "Isolated browsing", Banner: true, Persist: true,
		},
		// create posts to /banner, the object lives under /banners
		createPath: cbiPath("banner"),
		itemPath:   cbiPath("banners", "b1"),
		wantBody: map[string]any{
			"name": "Corp Banner", "primaryColor": "#0076BE", "textColor": "#FFFFFF",
			"notificationTitle": "Heads up", "banner": true, "persist": true,
		},
		wantState: map[string]any{"primary_color": "#0076BE", "persist": true},
	})
}
