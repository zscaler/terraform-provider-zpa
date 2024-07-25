package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/segmentgroup"
)

func dataSourceSegmentGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSegmentGroupRead,
		Schema: map[string]*schema.Schema{
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"domain_name": {
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
						"health_check_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_anchored": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"log_features": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"modified_by": {
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
									"modified_by": {
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
						"tcp_ports_in": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"tcp_ports_out": {
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
				},
			},
			// "config_space": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
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
				Optional: true,
			},
			"modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microtenant_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceSegmentGroupRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.SegmentGroup

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var resp *segmentgroup.SegmentGroup
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for segment group %s\n", id)
		res, _, err := segmentgroup.Get(service, id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for segment group name %s\n", name)
		res, _, err := segmentgroup.GetByName(service, name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(resp.ID)
		// _ = d.Set("config_space", resp.ConfigSpace)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("description", resp.Description)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("microtenant_id", resp.MicroTenantID)
		_ = d.Set("microtenant_name", resp.MicroTenantName)

		if err := d.Set("applications", flattenSegmentGroupApplications(resp)); err != nil {
			return fmt.Errorf("failed to read applications %s", err)
		}
	} else {
		return fmt.Errorf("couldn't find any segment group with name '%s' or id '%s'", name, id)
	}

	return nil
}

func flattenSegmentGroupApplications(segmentGroup *segmentgroup.SegmentGroup) []interface{} {
	segmentGroupApplications := make([]interface{}, len(segmentGroup.Applications))
	for i, segmentGroupApplication := range segmentGroup.Applications {
		segmentGroupApplications[i] = map[string]interface{}{
			// "bypass_type":            segmentGroupApplication.BypassType,
			// "config_space":           segmentGroupApplication.ConfigSpace,
			// "creation_time":          segmentGroupApplication.CreationTime,
			// "default_idle_timeout":   segmentGroupApplication.DefaultIdleTimeout,
			// "default_max_age":        segmentGroupApplication.DefaultMaxAge,
			// "description":            segmentGroupApplication.Description,
			// "domain_name":            segmentGroupApplication.DomainName,
			// "domain_names":           segmentGroupApplication.DomainNames,
			// "double_encrypt":         segmentGroupApplication.DoubleEncrypt,
			// "enabled":                segmentGroupApplication.Enabled,
			// "health_check_type":      segmentGroupApplication.HealthCheckType,
			// "ip_anchored":            segmentGroupApplication.IPAnchored,
			// "modified_by":            segmentGroupApplication.ModifiedBy,
			// "modified_time":          segmentGroupApplication.ModifiedTime,
			"name": segmentGroupApplication.Name,
			"id":   segmentGroupApplication.ID,
			// "passive_health_enabled": segmentGroupApplication.PassiveHealthEnabled,
			// "tcp_port_ranges":        segmentGroupApplication.TCPPortRanges,
			// "tcp_ports_in":           segmentGroupApplication.TCPPortsIn,
			// "tcp_ports_out":          segmentGroupApplication.TCPPortsOut,
			"server_groups": flattenAppServerGroup(segmentGroupApplication),
		}
	}

	return segmentGroupApplications
}

func flattenAppServerGroup(segmentGroup segmentgroup.Application) []interface{} {
	segmentServerGroups := make([]interface{}, len(segmentGroup.ServerGroup))
	for i, segmentServerGroup := range segmentGroup.ServerGroup {
		segmentServerGroups[i] = map[string]interface{}{
			"config_space":  segmentServerGroup.ConfigSpace,
			"creation_time": segmentServerGroup.CreationTime,
			"description":   segmentServerGroup.Description,
			"enabled":       segmentServerGroup.Enabled,
			"id":            segmentServerGroup.ID,
			"modified_by":   segmentServerGroup.ModifiedBy,
			"modified_time": segmentServerGroup.ModifiedTime,
			"name":          segmentServerGroup.Name,
		}
	}

	return segmentServerGroups
}
