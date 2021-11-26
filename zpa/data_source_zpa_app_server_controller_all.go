package zpa

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/appservercontroller"
)

func dataSourceApplicationServerAll() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApplicationServerAllRead,
		Schema: map[string]*schema.Schema{
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: applicationServerSchema(),
				},
			},
		},
	}
}

func dataSourceApplicationServerAllRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	list, _, err := zClient.appservercontroller.GetAll()
	if err != nil {
		return err
	}

	d.SetId("app-server-controller-list")
	_ = d.Set("list", flattenApplicationServer(list))
	return nil

}

func flattenApplicationServer(list []appservercontroller.ApplicationServer) []interface{} {
	appServer := make([]interface{}, len(list))
	for i, item := range list {
		appServer[i] = map[string]interface{}{
			"id":                   item.ID,
			"address":              item.Address,
			"creation_time":        item.CreationTime,
			"app_server_group_ids": item.AppServerGroupIds,
			"modifiedby":           item.ModifiedBy,
			"modified_time":        item.ModifiedTime,
			"name":                 item.Name,
			"config_space":         item.ConfigSpace,
			"description":          item.Description,
		}
	}
	return appServer
}
