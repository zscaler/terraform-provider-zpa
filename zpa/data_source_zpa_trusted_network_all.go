package zpa

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/trustednetwork"
)

func dataSourceTrustedNetworkAll() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTrustedNetworkAllRead,
		Schema: map[string]*schema.Schema{
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: trustedNetworkSchema(),
				},
			},
		},
	}
}

func dataSourceTrustedNetworkAllRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	list, _, err := zClient.trustednetwork.GetAll()
	if err != nil {
		return err
	}

	d.SetId("trusted-network-list")
	_ = d.Set("list", flattenTrustedNetworkList(list))
	return nil
}

func flattenTrustedNetworkList(list []trustednetwork.TrustedNetwork) []interface{} {
	keys := make([]interface{}, len(list))
	for i, item := range list {
		keys[i] = map[string]interface{}{
			"id":                 item.ID,
			"creation_time":      item.CreationTime,
			"domain":             item.Domain,
			"modifiedby":         item.ModifiedBy,
			"modified_time":      item.ModifiedTime,
			"name":               item.Name,
			"network_id":         item.NetworkID,
			"zscaler_cloud":      item.ZscalerCloud,
			"master_customer_id": item.MasterCustomerID,
		}
	}
	return keys
}
