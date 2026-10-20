package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbiregions"
)

var cbiRegions = []cbiregions.CBIRegions{{ID: "r0", Name: "Europe"}, {ID: "r1", Name: "Washington"}}

func TestDataSourceCBIRegion(t *testing.T) {
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceCBIRegions(),
		listPath:  cbiPath("regions"),
		plainList: true,
		items:     cbiRegions,
		name:      "washington",
		id:        "r1",
	})
}

// There is no single-region endpoint: the ID lookup searches the list.
func TestDataSourceCBIRegion_ByID(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", cbiPath("regions"), common.SuccessResponse(cbiRegions))

	d := fx.resourceData(dataSourceCBIRegions(), "r1", map[string]any{})
	fx.run(dataSourceCBIRegions().ReadContext, d)
	assertState(t, d, map[string]any{"name": "Washington"})
}
