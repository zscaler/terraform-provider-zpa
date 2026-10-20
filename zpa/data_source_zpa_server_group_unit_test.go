package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
)

func TestDataSourceServerGroup(t *testing.T) {
	match := servergroup.ServerGroup{ID: "sg1", Name: "CRM servers", Enabled: true, DynamicDiscovery: true}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceServerGroup(),
		listPath:  v1Path("serverGroup"),
		itemPath:  v1Path("serverGroup", "sg1"),
		items:     []servergroup.ServerGroup{{ID: "other", Name: "CRM servers (old)"}, match},
		item:      match,
		name:      "crm servers",
		id:        "sg1",
		wantState: map[string]any{"dynamic_discovery": true, "enabled": true},
	})
}
