package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/customerversionprofile"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgegroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/trustednetwork"
)

func serviceEdgeGroupConfig() map[string]any {
	return map[string]any{
		"name":                     "SJC Service Edges",
		"enabled":                  true,
		"is_public":                true,
		"latitude":                 "37.3382082",
		"longitude":                "-121.8863286",
		"location":                 "San Jose, CA, USA",
		"upgrade_day":              "SUNDAY",
		"upgrade_time_in_secs":     "66600",
		"override_version_profile": true,
		"version_profile_id":       "0",
		"enrollment_cert_id":       "ec1",
		"trusted_networks":         []any{map[string]any{"id": []any{"tn1"}}},
	}
}

func serviceEdgeGroupResponse() serviceedgegroup.ServiceEdgeGroup {
	return serviceedgegroup.ServiceEdgeGroup{
		ID: "seg1", Name: "SJC Service Edges", Enabled: true, IsPublic: "TRUE",
		Latitude: "37.3382082", Longitude: "-121.8863286", Location: "San Jose, CA, USA",
		UpgradeDay: "SUNDAY", UpgradeTimeInSecs: "66600", OverrideVersionProfile: true, VersionProfileID: "0",
		EnrollmentCertID: "ec1", TrustedNetworks: []trustednetwork.TrustedNetwork{{ID: "tn1"}},
	}
}

func TestResourceServiceEdgeGroup_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource:   resourceServiceEdgeGroup(),
		config:     serviceEdgeGroupConfig(),
		id:         "seg1",
		response:   serviceEdgeGroupResponse(),
		createPath: v1Path("serviceEdgeGroup"),
		itemPath:   v1Path("serviceEdgeGroup", "seg1"),
		// is_public is a boolean in HCL and "TRUE"/"FALSE" in the API
		wantBody: map[string]any{
			"name": "SJC Service Edges", "isPublic": "TRUE", "upgradeDay": "SUNDAY", "enrollmentCertId": "ec1",
			"trustedNetworks": []any{map[string]any{"id": "tn1"}},
		},
		wantState: map[string]any{
			"is_public":        true,
			"trusted_networks": []any{map[string]any{"id": []any{"tn1"}}},
		},
	})
}

// version_profile_name is resolved to version_profile_id when the profile is overridden.
func TestResourceServiceEdgeGroup_VersionProfileNameResolved(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("visible", "versionProfiles"), common.SuccessResponse(common.ZPAList([]customerversionprofile.CustomerVersionProfile{
		{ID: "0", Name: "Default"}, {ID: "2", Name: "New Release"},
	})))
	fx.on("POST", fx.path("serviceEdgeGroup"), common.CreatedResponse(serviceEdgeGroupResponse()))
	fx.on("GET", fx.path("serviceEdgeGroup", "seg1"), common.SuccessResponse(serviceEdgeGroupResponse()))

	config := serviceEdgeGroupConfig()
	delete(config, "version_profile_id")
	config["version_profile_name"] = "New Release"
	d := fx.resourceData(resourceServiceEdgeGroup(), "", config)
	fx.run(resourceServiceEdgeGroupCreate, d)

	assertSubset(t, "create request", fx.body("POST", fx.path("serviceEdgeGroup")), map[string]any{"versionProfileId": "2"})
}

func TestResourceServiceEdgeGroup_OverrideWithoutVersionProfileRejected(t *testing.T) {
	fx := newUnitFixture(t)
	config := serviceEdgeGroupConfig()
	delete(config, "version_profile_id")

	d := fx.resourceData(resourceServiceEdgeGroup(), "", config)
	if msg := fx.runErr(resourceServiceEdgeGroupCreate, d); !strings.Contains(msg, "version_profile_name must be provided") {
		t.Fatalf("expected an overridden version profile without a name to be rejected, got %q", msg)
	}
	if n := fx.calls("POST", fx.path("serviceEdgeGroup")); n != 0 {
		t.Errorf("no create request may be sent, got %d", n)
	}
}

// Without enrollment_cert_id the "Service Edge" enrollment certificate is used.
func TestResourceServiceEdgeGroup_EnrollmentCertAutoResolved(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.pathV2("enrollmentCert"), common.SuccessResponse(common.ZPAList([]map[string]any{
		{"id": "ec1", "name": "Connector"}, {"id": "ec2", "name": "Service Edge"},
	})))
	fx.on("POST", fx.path("serviceEdgeGroup"), common.CreatedResponse(serviceEdgeGroupResponse()))
	fx.on("GET", fx.path("serviceEdgeGroup", "seg1"), common.SuccessResponse(serviceEdgeGroupResponse()))

	config := serviceEdgeGroupConfig()
	delete(config, "enrollment_cert_id")
	d := fx.resourceData(resourceServiceEdgeGroup(), "", config)
	fx.run(resourceServiceEdgeGroupCreate, d)

	assertSubset(t, "create request", fx.body("POST", fx.path("serviceEdgeGroup")), map[string]any{"enrollmentCertId": "ec2"})
}
