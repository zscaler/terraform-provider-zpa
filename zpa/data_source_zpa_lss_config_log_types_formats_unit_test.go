package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/lssconfigcontroller"
)

func TestDataSourceLSSLogTypeFormats(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", lssAdminPath+"/logType/formats", common.SuccessResponse(lssconfigcontroller.LSSFormats{
		Csv: "%s{LogTimestamp:time}", Tsv: "%s{LogTimestamp:time}\t", Json: `{"LogTimestamp":%j{LogTimestamp:time}}`,
	}))

	d := fx.resourceData(dataSourceLSSLogTypeFormats(), "", map[string]any{"log_type": "zpn_trans_log"})
	fx.run(dataSourceLSSLogTypeFormatsRead, d)
	if d.Id() != "lss_log_types_zpn_trans_log" {
		t.Errorf("id = %q", d.Id())
	}
	assertState(t, d, map[string]any{"csv": "%s{LogTimestamp:time}"})
	reqs := fx.api.Server.AllRequests()
	if len(reqs) != 1 || !strings.Contains(reqs[0].Query, "logType=zpn_trans_log") {
		t.Errorf("the log type must be passed as a query parameter, got %v", reqs)
	}
}
