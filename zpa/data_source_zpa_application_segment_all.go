package zpa

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/applicationsegment"
)

func dataSourceApplicationSegmentAll() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApplicationSegmentAllRead,
		Schema: map[string]*schema.Schema{
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: appSegmentSchema(),
				},
			},
		},
	}
}

func dataSourceApplicationSegmentAllRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	list, _, err := zClient.applicationsegment.GetAll()
	log.Printf("[INFO] got %d apps\n", len(list))
	if err != nil {
		return err
	}
	d.SetId("app-segment-list")
	_ = d.Set("list", flattenAppSegmentList(list))

	return nil
}

func flattenAppSegmentList(list []applicationsegment.ApplicationSegmentResource) []interface{} {
	appSegments := make([]interface{}, len(list))
	for i, item := range list {
		appSegments[i] = map[string]interface{}{
			"id":                     item.ID,
			"segment_group_id":       item.SegmentGroupID,
			"segment_group_name":     item.SegmentGroupName,
			"bypass_type":            item.BypassType,
			"config_space":           item.ConfigSpace,
			"creation_time":          item.CreationTime,
			"description":            item.Description,
			"domain_names":           item.DomainNames,
			"double_encrypt":         item.DoubleEncrypt,
			"enabled":                item.Enabled,
			"health_checktype":       item.HealthCheckType,
			"health_reporting":       item.HealthReporting,
			"ip_anchored":            item.IpAnchored,
			"is_cname_enabled":       item.IsCnameEnabled,
			"modifiedby":             item.ModifiedBy,
			"modified_time":          item.ModifiedTime,
			"name":                   item.Name,
			"passive_health_enabled": item.PassiveHealthEnabled,
			"tcp_port_ranges":        item.TCPPortRanges,
			"udp_port_ranges":        item.UDPPortRanges,
			"clientless_apps":        flattenClientlessApps(&item),
			"server_groups":          flattenAppServerGroups(&item),
		}
	}
	return appSegments
}
