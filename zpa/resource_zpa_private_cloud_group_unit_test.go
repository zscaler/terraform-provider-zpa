package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_group"
)

func TestResourcePrivateCloudGroup_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourcePrivateCloudGroup(),
		config: map[string]any{
			"name":                     "SJC Private Cloud",
			"enabled":                  true,
			"city_country":             "San Jose, US",
			"country_code":             "US",
			"latitude":                 "37.3382082",
			"longitude":                "-121.8863286",
			"location":                 "San Jose, CA, USA",
			"upgrade_day":              "SUNDAY",
			"upgrade_time_in_secs":     "66600",
			"override_version_profile": false,
			"enrollment_cert_id":       "ec1",
		},
		id: "pcg1",
		response: private_cloud_group.PrivateCloudGroup{
			ID: "pcg1", Name: "SJC Private Cloud", Enabled: true, CityCountry: "San Jose, US", CountryCode: "US",
			Latitude: "37.3382082", Longitude: "-121.8863286", Location: "San Jose, CA, USA",
			UpgradeDay: "SUNDAY", UpgradeTimeInSecs: "66600", EnrollmentCertID: "ec1",
		},
		createPath: v1Path("privateCloudControllerGroup"),
		itemPath:   v1Path("privateCloudControllerGroup", "pcg1"),
		wantBody: map[string]any{
			"name": "SJC Private Cloud", "countryCode": "US", "latitude": "37.3382082", "enrollmentCertId": "ec1",
		},
		wantState: map[string]any{"country_code": "US", "location": "San Jose, CA, USA"},
	})
}
