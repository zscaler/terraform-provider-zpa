package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praconsole"
)

func TestDataSourcePRAConsoleController(t *testing.T) {
	match := praconsole.PRAConsole{ID: "con1", Name: "ssh_console", Enabled: true, IconText: "ssh"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourcePRAConsoleController(),
		listPath:  v1Path("praConsole"),
		itemPath:  v1Path("praConsole", "con1"),
		items:     []praconsole.PRAConsole{{ID: "other", Name: "ssh_console-old"}, match},
		item:      match,
		name:      "ssh_console",
		id:        "con1",
		wantState: map[string]any{"enabled": true, "icon_text": "ssh"},
	})
}
