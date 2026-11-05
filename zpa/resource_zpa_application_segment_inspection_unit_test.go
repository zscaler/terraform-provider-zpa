package zpa

import (
	"fmt"
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

func inspectionApp(domain, port, cert string) map[string]any {
	return map[string]any{
		"domain": domain, "application_port": port, "application_protocol": "HTTPS",
		"certificate_id": cert, "app_types": []any{"INSPECT"},
	}
}

func inspectionConfigWith(apps ...map[string]any) map[string]any {
	config := inspectionSegmentConfig()
	domains, list := []any{}, []any{}
	for _, a := range apps {
		domains = append(domains, a["domain"])
		list = append(list, a)
	}
	config["domain_names"] = domains
	config["common_apps_dto"] = []any{map[string]any{"apps_config": list}}
	return config
}

var (
	jenkinsInspectApp = applicationsegmentinspection.InspectionAppDto{
		ID: "ia1", AppID: "ai1", Name: "jenkins.acme.com", Domain: "jenkins.acme.com",
		ApplicationPort: "443", ApplicationProtocol: "HTTPS", CertificateID: "c1",
	}
	gitInspectApp = applicationsegmentinspection.InspectionAppDto{
		ID: "ia2", AppID: "ai1", Name: "git.acme.com", Domain: "git.acme.com",
		ApplicationPort: "443", ApplicationProtocol: "HTTPS", CertificateID: "c1",
	}
)

// apps_config is a list in the configured order, so apps_config[N] can be
// referenced whatever order the API returns the apps in.
func TestResourceApplicationSegmentInspection_InspectAppIDIsIndexable(t *testing.T) {
	fx := newUnitFixture(t)
	resp := inspectionSegmentResponse()
	resp.InspectionAppDto = []applicationsegmentinspection.InspectionAppDto{gitInspectApp, jenkinsInspectApp}
	fx.on("GET", fx.path("application", "ai1"), common.SuccessResponse(resp))

	d := fx.resourceData(resourceApplicationSegmentInspection(), "ai1",
		inspectionConfigWith(inspectionApp("jenkins.acme.com", "443", "c1"), inspectionApp("git.acme.com", "443", "c1")))
	fx.run(resourceApplicationSegmentInspectionRead, d)

	for i, want := range []string{"ia1", "ia2"} {
		if got := d.Get(fmt.Sprintf("common_apps_dto.0.apps_config.%d.inspect_app_id", i)); got != want {
			t.Errorf("apps_config[%d].inspect_app_id = %v, want %s", i, got, want)
		}
	}
}

// Changing an existing app's settings plans an update and sends the configured
// values, with the app's IDs.
func TestResourceApplicationSegmentInspection_UpdateSendsChangedAppSettings(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("application", "ai1"), common.SuccessResponse(inspectionSegmentResponse()))
	fx.on("PUT", fx.path("application", "ai1"), common.NoContentResponse())

	changed := inspectionApp("jenkins.acme.com", "443", "c2")
	changed["trust_untrusted_cert"] = true
	d := fx.updateData(resourceApplicationSegmentInspection(), "ai1",
		inspectionConfigWith(inspectionApp("jenkins.acme.com", "443", "c1")), inspectionConfigWith(changed))
	if !d.HasChange("common_apps_dto") {
		t.Fatal("changing certificate_id and trust_untrusted_cert planned no update")
	}
	fx.run(resourceApplicationSegmentInspectionUpdate, d)

	assertSubset(t, "update request", fx.body("PUT", fx.path("application", "ai1")), map[string]any{
		"commonAppsDto": map[string]any{"appsConfig": []any{map[string]any{
			"domain": "jenkins.acme.com", "appId": "ai1", "inspectAppId": "ia1",
			"certificateId": "c2", "trustUntrustedCert": true,
		}}},
	})
}

// An app added to an existing segment plans its IDs as unknown (known after
// apply); an app that moves to another position keeps its own IDs.
func TestResourceApplicationSegmentInspection_PlanAddedAppIDsUnknown(t *testing.T) {
	state := inspectionConfigWith(inspectionApp("jenkins.acme.com", "443", "c1"))
	app := state["common_apps_dto"].([]any)[0].(map[string]any)["apps_config"].([]any)[0].(map[string]any)
	app["app_id"], app["inspect_app_id"], app["name"] = "ai1", "ia1", "jenkins.acme.com"

	jenkins, git := inspectionApp("jenkins.acme.com", "443", "c1"), inspectionApp("git.acme.com", "443", "c1")
	for name, order := range map[string][]map[string]any{
		"added after":  {jenkins, git},
		"added before": {git, jenkins},
	} {
		t.Run(name, func(t *testing.T) {
			for _, a := range plannedCommonApps(t, resourceApplicationSegmentInspection(), state, inspectionConfigWith(order...)) {
				id := a.GetAttr("inspect_app_id")
				switch a.GetAttr("domain").AsString() {
				case "jenkins.acme.com":
					if !id.IsKnown() || id.AsString() != "ia1" {
						t.Errorf("existing app: planned inspect_app_id = %#v, want ia1", id)
					}
				case "git.acme.com":
					if id.IsKnown() {
						t.Errorf("added app: planned inspect_app_id = %#v, want unknown", id)
					}
				}
			}
		})
	}
}
