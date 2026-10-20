package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbibannercontroller"
)

func TestDataSourceCBIBanner(t *testing.T) {
	match := cbibannercontroller.CBIBannerController{ID: "b1", Name: "Corp Banner", PrimaryColor: "#0076BE", Banner: true}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceCBIBanners(),
		listPath:  cbiPath("banners"),
		plainList: true,
		itemPath:  cbiPath("banners", "b1"),
		items:     []cbibannercontroller.CBIBannerController{{ID: "b0", Name: "Corp Banner Old"}, match},
		item:      match,
		name:      "corp banner",
		id:        "b1",
		wantState: map[string]any{"primary_color": "#0076BE", "banner": true},
	})
}
