package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/machinegroup"
)

func TestDataSourceMachineGroup(t *testing.T) {
	match := machinegroup.MachineGroup{ID: "mg1", Name: "Corp Machines", Enabled: true, Machines: []machinegroup.Machines{{ID: "m1", Name: "laptop-01", MachineGroupID: "mg1"}}}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceMachineGroup(),
		listPath:  v1Path("machineGroup"),
		itemPath:  v1Path("machineGroup", "mg1"),
		items:     []machinegroup.MachineGroup{{ID: "other", Name: "Corp Machines (old)"}, match},
		item:      match,
		name:      "corp machines",
		id:        "mg1",
		wantState: map[string]any{"enabled": true, "machines": []any{map[string]any{"id": "m1", "name": "laptop-01", "machine_group_id": "mg1"}}},
	})
}
