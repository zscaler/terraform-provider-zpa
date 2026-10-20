package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/custom_config_controller"
)

func ziaCloudConfig() map[string]any {
	return map[string]any{
		"zia_cloud_domain":          "zscalerthree",
		"zia_username":              "zpa-integration@acme.com",
		"zia_password":              "s3cret",
		"zia_sandbox_api_token":     "sandbox-token",
		"zia_cloud_service_api_key": "api-key",
	}
}

var ziaCloudConfigResponse = custom_config_controller.ZIACloudConfig{
	ZIACloudDomain: "zscalerthree.net", ZIAUsername: "zpa-integration@acme.com",
}

// The tenant has one ZIA cloud configuration: create and update both POST it.
func TestResourceZiaCloudConfig(t *testing.T) {
	path := v1Path("config", "ziaCloudConfig")
	for _, step := range []struct {
		name string
		fn   crudFunc
		id   string
	}{
		{"Create", resourceZiaCloudConfigCreate, ""},
		{"Update", resourceZiaCloudConfigUpdate, "zia_cloud_config"},
	} {
		t.Run(step.name, func(t *testing.T) {
			fx := newUnitFixture(t)
			fx.on("POST", path, common.SuccessResponse(ziaCloudConfigResponse))
			fx.on("GET", path, common.SuccessResponse(ziaCloudConfigResponse))
			d := fx.resourceData(resourceZiaCloudConfig(), step.id, ziaCloudConfig())
			fx.run(step.fn, d)
			if d.Id() != "zia_cloud_config" {
				t.Errorf("id = %q", d.Id())
			}
			// the cloud domain is sent with its .net suffix
			assertSubset(t, "request", fx.body("POST", path), map[string]any{
				"ziaCloudDomain": "zscalerthree.net", "ziaUsername": "zpa-integration@acme.com",
			})
			assertState(t, d, map[string]any{"zia_cloud_domain": "zscalerthree.net"})
		})
	}

	t.Run("ReadAPIError", func(t *testing.T) {
		fx := newUnitFixture(t)
		fx.on("GET", path, apiError())
		d := fx.resourceData(resourceZiaCloudConfig(), "zia_cloud_config", ziaCloudConfig())
		if msg := fx.runErr(resourceZiaCloudConfigRead, d); !strings.Contains(msg, apiErrorMarker) {
			t.Errorf("expected the API error to be returned, got %q", msg)
		}
	})

	t.Run("ReadNotFoundRemovesFromState", func(t *testing.T) {
		fx := newUnitFixture(t)
		fx.on("GET", path, common.NotFoundResponse())
		d := fx.resourceData(resourceZiaCloudConfig(), "zia_cloud_config", ziaCloudConfig())
		fx.run(resourceZiaCloudConfigRead, d)
		if d.Id() != "" {
			t.Errorf("id = %q, want it removed from state", d.Id())
		}
	})

	t.Run("CreateAPIError", func(t *testing.T) {
		fx := newUnitFixture(t)
		fx.on("POST", path, apiError())
		d := fx.resourceData(resourceZiaCloudConfig(), "", ziaCloudConfig())
		if msg := fx.runErr(resourceZiaCloudConfigCreate, d); !strings.Contains(msg, apiErrorMarker) {
			t.Errorf("expected the API error to be returned, got %q", msg)
		}
	})

	t.Run("Import", func(t *testing.T) {
		fx := newUnitFixture(t)
		fx.on("GET", path, common.SuccessResponse(ziaCloudConfigResponse))
		d := resourceZiaCloudConfig().Data(nil)
		d.SetId("anything")
		got, err := resourceZiaCloudConfig().Importer.StateContext(t.Context(), d, fx.client)
		if err != nil || len(got) != 1 || got[0].Id() != "zia_cloud_config" {
			t.Errorf("import: err=%v ids=%v", err, importedIDs(got))
		}
	})

	t.Run("DeleteMakesNoAPICall", func(t *testing.T) {
		fx := newUnitFixture(t)
		d := fx.resourceData(resourceZiaCloudConfig(), "zia_cloud_config", ziaCloudConfig())
		fx.run(resourceZiaCloudConfig().DeleteContext, d)
		if n := len(fx.api.Server.AllRequests()); n != 0 {
			t.Errorf("expected no API calls, got %d", n)
		}
	})
}
