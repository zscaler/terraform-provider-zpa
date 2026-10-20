package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentinspection"
	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
)

func inspectionSegmentConfig() map[string]any {
	return map[string]any{
		"name":             "jenkins.acme.com",
		"enabled":          true,
		"domain_names":     []any{"jenkins.acme.com"},
		"segment_group_id": "sg1",
		"tcp_port_ranges":  []any{"443", "443"},
		"tcp_port_range":   []any{map[string]any{"from": "443", "to": "443"}},
		"common_apps_dto": []any{map[string]any{"apps_config": []any{map[string]any{
			"domain":               "jenkins.acme.com",
			"application_port":     "443",
			"application_protocol": "HTTPS",
			"certificate_id":       "c1",
			"app_types":            []any{"INSPECT"},
		}}}},
	}
}

func inspectionSegmentResponse() applicationsegmentinspection.AppSegmentInspection {
	return applicationsegmentinspection.AppSegmentInspection{
		ID: "ai1", Name: "jenkins.acme.com", Enabled: true, SegmentGroupID: "sg1",
		DomainNames:     []string{"jenkins.acme.com"},
		TCPAppPortRange: []zpacommon.NetworkPorts{{From: "443", To: "443"}},
		InspectionAppDto: []applicationsegmentinspection.InspectionAppDto{{
			ID: "ia1", AppID: "ai1", Name: "jenkins.acme.com", Domain: "jenkins.acme.com",
			ApplicationPort: "443", ApplicationProtocol: "HTTPS", CertificateID: "c1",
		}},
	}
}

func TestResourceApplicationSegmentInspection_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource:   resourceApplicationSegmentInspection(),
		config:     inspectionSegmentConfig(),
		id:         "ai1",
		response:   inspectionSegmentResponse(),
		createPath: v1Path("application"),
		itemPath:   v1Path("application", "ai1"),
		setup: func(fx *unitFixture) {
			fx.on("GET", fx.path("segmentGroup", "sg1"), common.SuccessResponse(segmentgroup.SegmentGroup{ID: "sg1"}))
			fx.on("PUT", fx.path("segmentGroup", "sg1"), common.NoContentResponse())
		},
		wantBody: map[string]any{
			"segmentGroupId": "sg1", "domainNames": []any{"jenkins.acme.com"},
			"commonAppsDto": map[string]any{"appsConfig": []any{map[string]any{
				"domain": "jenkins.acme.com", "name": "jenkins.acme.com", "applicationProtocol": "HTTPS",
				"applicationPort": "443", "appTypes": []any{"INSPECT"},
			}}},
		},
		wantState: map[string]any{
			"common_apps_dto": []any{map[string]any{"apps_config": []any{map[string]any{
				"domain": "jenkins.acme.com", "app_id": "ai1", "inspect_app_id": "ia1",
			}}}},
		},
	})
}

// With Auto App Protection the API reports the protocol as DYNAMIC; the
// configured protocol stays in state so the configuration does not drift.
func TestResourceApplicationSegmentInspection_DynamicProtocolKeepsConfiguredValue(t *testing.T) {
	fx := newUnitFixture(t)
	resp := inspectionSegmentResponse()
	resp.AutoAppProtectEnabled = true
	resp.InspectionAppDto[0].ApplicationProtocol = "DYNAMIC"
	resp.InspectionAppDto[0].ApplicationPort = ""
	fx.on("GET", fx.path("application", "ai1"), common.SuccessResponse(resp))

	d := fx.resourceData(resourceApplicationSegmentInspection(), "ai1", inspectionSegmentConfig())
	fx.run(resourceApplicationSegmentInspectionRead, d)
	assertState(t, d, map[string]any{
		"common_apps_dto": []any{map[string]any{"apps_config": []any{map[string]any{
			"application_protocol": "HTTPS", "application_port": "443",
		}}}},
	})
}

// Deleting the segment first removes it from its segment group.
func TestResourceApplicationSegmentInspection_DeleteDetachesFromSegmentGroup(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("segmentGroup", "sg1"), common.SuccessResponse(segmentgroup.SegmentGroup{
		ID: "sg1", Applications: []segmentgroup.Application{{ID: "ai1"}, {ID: "app2"}},
	}))
	fx.on("PUT", fx.path("segmentGroup", "sg1"), common.NoContentResponse())
	fx.on("DELETE", fx.path("application", "ai1"), common.NoContentResponse())

	d := fx.resourceData(resourceApplicationSegmentInspection(), "ai1", inspectionSegmentConfig())
	fx.run(resourceApplicationSegmentInspectionDelete, d)

	assertSubset(t, "segment group update", fx.body("PUT", fx.path("segmentGroup", "sg1")),
		map[string]any{"applications": []any{map[string]any{"id": "app2"}}})
	if fx.calls("DELETE", fx.path("application", "ai1")) != 1 {
		t.Error("the segment itself was not deleted")
	}
}
