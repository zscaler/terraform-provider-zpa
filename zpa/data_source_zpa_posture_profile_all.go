package zpa

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/postureprofile"
)

func dataSourcePostureProfileAll() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePostureProfileAllRead,
		Schema: map[string]*schema.Schema{
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: postureProfileSchema(),
				},
			},
		},
	}
}

func dataSourcePostureProfileAllRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	list, _, err := zClient.postureprofile.GetAll()
	if err != nil {
		return err
	}
	d.SetId("posture-profile-list")
	_ = d.Set("list", flattenPostureProfileList(list))
	return nil
}
func flattenPostureProfileList(list []postureprofile.PostureProfile) []interface{} {
	keys := make([]interface{}, len(list))
	for i, item := range list {
		keys[i] = map[string]interface{}{
			"id":                  item.ID,
			"creation_time":       item.CreationTime,
			"domain":              item.Domain,
			"master_customer_id":  item.MasterCustomerID,
			"modifiedby":          item.ModifiedBy,
			"modified_time":       item.ModifiedTime,
			"name":                item.Name,
			"posture_udid":        item.PostureudID,
			"zscaler_cloud":       item.ZscalerCloud,
			"zscaler_customer_id": item.ZscalerCustomerID,
		}
	}
	return keys
}
