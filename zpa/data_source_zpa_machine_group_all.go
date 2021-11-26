package zpa

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/machinegroup"
)

func dataSourceMachineGroupAll() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMachineGroupAllRead,
		Schema: map[string]*schema.Schema{
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: machineGroupSchema(),
				},
			},
		},
	}
}

func dataSourceMachineGroupAllRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	list, _, err := zClient.machinegroup.GetAll()
	if err != nil {
		return err
	}

	d.SetId("machine-group-list")
	_ = d.Set("list", flattenMachineGroupList(list))
	return nil
}

func flattenMachineGroupList(list []machinegroup.MachineGroup) []interface{} {
	keys := make([]interface{}, len(list))
	for i, item := range list {
		keys[i] = map[string]interface{}{
			"id":            item.ID,
			"creation_time": item.CreationTime,
			"description":   item.Description,
			"enabled":       item.Enabled,
			"modifiedby":    item.ModifiedBy,
			"modified_time": item.ModifiedTime,
			"name":          item.Name,
			"machines":      flattenMachines(&item),
		}
	}
	return keys
}
