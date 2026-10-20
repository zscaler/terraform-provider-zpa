package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/platforms"
)

func TestDataSourceAccessPolicyPlatforms(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("platform"), common.SuccessResponse(platforms.Platforms{
		Linux: "Linux", Android: "Android", Windows: "Windows", IOS: "iOS", MacOS: "Mac",
	}))

	d := fx.resourceData(dataSourceAccessPolicyPlatforms(), "", map[string]any{})
	fx.run(dataSourceAccessPolicyPlatformsRead, d)
	assertState(t, d, map[string]any{"linux": "Linux", "ios": "iOS", "mac": "Mac"})
}
