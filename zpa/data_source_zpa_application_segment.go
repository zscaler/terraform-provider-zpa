package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApplicationSegment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationSegmentRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
			"health_check_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_reporting": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"select_connector_close_to_app": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"use_in_dr_mode": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_incomplete_dr_config": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of the server group IDs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microtenant_name": {
				Type:     schema.TypeString,
				Optional: true,
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
			"tcp_port_range": resourceNetworkPortsSchema("tcp port range"),
			"udp_port_range": resourceNetworkPortsSchema("udp port range"),
		},
	}
}

func dataSourceApplicationSegmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var resp *applicationsegment.ApplicationSegmentResource
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for application segment %s\n", id)
		res, _, err := applicationsegment.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data for application segment name %s\n", name)
		res, _, err := applicationsegment.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
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
		_ = d.Set("health_check_type", resp.HealthCheckType)
		_ = d.Set("health_reporting", resp.HealthReporting)
		_ = d.Set("select_connector_close_to_app", resp.SelectConnectorCloseToApp)
		_ = d.Set("use_in_dr_mode", resp.UseInDrMode)
		_ = d.Set("is_incomplete_dr_config", resp.IsIncompleteDRConfig)
		_ = d.Set("ip_anchored", resp.IpAnchored)
		_ = d.Set("is_cname_enabled", resp.IsCnameEnabled)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("passive_health_enabled", resp.PassiveHealthEnabled)
		_ = d.Set("microtenant_id", resp.MicroTenantID)
		_ = d.Set("microtenant_name", resp.MicroTenantName)

		if err := d.Set("server_groups", flattenCommonAppServerGroups(resp.ServerGroups)); err != nil {
			return diag.FromErr(fmt.Errorf("failed to read app server groups %s", err))
		}

		if err := d.Set("tcp_port_ranges", resp.TCPPortRanges); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("udp_port_ranges", resp.UDPPortRanges); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("tcp_port_range", flattenNetworkPorts(resp.TCPAppPortRange)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("udp_port_range", flattenNetworkPorts(resp.UDPAppPortRange)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any application segment with name '%s' or id '%s'", name, id))
	}

	return nil
}
