package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/applicationsegment"
)

func dataSourceApplicationSegment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApplicationSegmentRead,
		Schema: map[string]*schema.Schema{
			"segment_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"segment_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bypass_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_space": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_idle_timeout": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_max_age": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"double_encrypt": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"health_checktype": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_reporting": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_anchored": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_cname_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"modifiedby": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"passive_health_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"server_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_space": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dynamic_discovery": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"modifiedby": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tcp_port_ranges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"udp_port_ranges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceApplicationSegmentRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	var resp *applicationsegment.ApplicationSegmentResource
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for server group %s\n", id)
		res, _, err := zClient.applicationsegment.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data for server group name %s\n", name)
		res, _, err := zClient.applicationsegment.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("segment_group_id", resp.SegmentGroupID)
		_ = d.Set("segment_group_name", resp.SegmentGroupName)
		_ = d.Set("bypass_type", resp.BypassType)
		_ = d.Set("config_space", resp.ConfigSpace)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("description", resp.Description)
		_ = d.Set("domain_names", resp.DomainNames)
		_ = d.Set("double_encrypt", resp.DoubleEncrypt)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("health_checktype", resp.HealthCheckType)
		_ = d.Set("health_reporting", resp.HealthReporting)
		_ = d.Set("ip_anchored", resp.IpAnchored)
		_ = d.Set("is_cname_enabled", resp.IsCnameEnabled)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("passive_health_enabled", resp.PassiveHealthEnabled)
		_ = d.Set("tcp_port_ranges", resp.TCPPortRanges)
		_ = d.Set("udp_port_ranges", resp.UDPPortRanges)

		if err := d.Set("server_groups", flattenAppServerGroups(resp)); err != nil {
			return fmt.Errorf("failed to read app server groups %s", err)
		}
	} else {
		return fmt.Errorf("couldn't find any application segment with name '%s' or id '%s'", name, id)
	}

	return nil

}

func flattenAppServerGroups(serverGroup *applicationsegment.ApplicationSegmentResource) []interface{} {
	serverGroups := make([]interface{}, len(serverGroup.ServerGroups))
	for i, val := range serverGroup.ServerGroups {
		serverGroups[i] = map[string]interface{}{
			"name":              val.Name,
			"id":                val.ID,
			"config_space":      val.ConfigSpace,
			"creation_time":     val.CreationTime,
			"description":       val.Description,
			"enabled":           val.Enabled,
			"dynamic_discovery": val.DynamicDiscovery,
			"modifiedby":        val.ModifiedBy,
			"modified_time":     val.ModifiedTime,
		}
	}

	return serverGroups
}
