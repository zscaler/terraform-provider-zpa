package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/browseraccess"
)

func dataSourceBrowserAccess() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBrowserAccessRead,
		Schema: map[string]*schema.Schema{
			"id": {
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
			"tcp_port_range": resourceNetworkPortsSchema("tcp port range"),
			"udp_port_range": resourceNetworkPortsSchema("udp port range"),

			"tcp_port_ranges": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "TCP port ranges used to access the app.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"udp_port_ranges": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "UDP port ranges used to access the app.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"config_space": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the application.",
			},
			"domain_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of domains and IPs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"double_encrypt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether Double Encryption is enabled or disabled for the app.",
			},
			"health_check_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"passive_health_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"health_reporting": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether health reporting for the app is Continuous or On Access. Supported values: NONE, ON_ACCESS, CONTINUOUS.",
			},
			"ip_anchored": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_cname_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the application.",
			},
			"clientless_apps": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_options": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"hidden": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trust_untrusted_cert": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
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
		},
	}
}

func dataSourceBrowserAccessRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	var resp *browseraccess.BrowserAccess
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for browser access application %s\n", id)
		res, _, err := zClient.browseraccess.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data for browser access application name %s\n", name)
		res, _, err := zClient.browseraccess.GetByName(name)
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
		_ = d.Set("domain_names", resp.DomainNames)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("passive_health_enabled", resp.PassiveHealthEnabled)
		_ = d.Set("double_encrypt", resp.DoubleEncrypt)
		_ = d.Set("health_check_type", resp.HealthCheckType)
		_ = d.Set("is_cname_enabled", resp.IsCnameEnabled)
		_ = d.Set("ip_anchored", resp.IPAnchored)
		_ = d.Set("health_reporting", resp.HealthReporting)
		_ = d.Set("tcp_port_ranges", resp.TCPPortRanges)
		_ = d.Set("udp_port_ranges", resp.UDPPortRanges)

		if err := d.Set("clientless_apps", flattenBaClientlessApps(resp)); err != nil {
			return fmt.Errorf("failed to read clientless apps %s", err)
		}

		if err := d.Set("server_groups", flattenClientlessAppServerGroups(resp.AppServerGroups)); err != nil {
			return fmt.Errorf("failed to read app server groups %s", err)
		}

		if err := d.Set("tcp_port_range", flattenNetworkPorts(resp.TCPAppPortRange)); err != nil {
			return err
		}

		if err := d.Set("tcp_port_range", flattenNetworkPorts(resp.UDPAppPortRange)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("couldn't find any browser access application with name '%s' or id '%s'", name, id)
	}

	return nil

}
