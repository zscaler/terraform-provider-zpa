package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/applicationsegmentpra"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/client"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/segmentgroup"
)

func resourceAppSegmentPRA() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppSegmentPRACreate,
		Read:   resourceAppSegmentPRARead,
		Update: resourceAppSegmentPRAUpdate,
		Delete: resourceAppSegmentPRADelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					d.Set("id", id)
				} else {
					resp, _, err := zClient.applicationsegmentpra.GetByName(id)
					if err == nil {
						d.SetId(resp.ID)
						d.Set("id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"segment_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"segment_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bypass_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether users can bypass ZPA to access applications. Default: NEVER. Supported values: ALWAYS, NEVER, ON_NET. The value NEVER indicates the use of the client forwarding policy.",
				ValidateFunc: validation.StringInSlice([]string{
					"ALWAYS",
					"NEVER",
					"ON_NET",
				}, false),
			},
			"tcp_port_range": resourceAppSegmentPortRange("tcp port range"),
			"udp_port_range": resourceAppSegmentPortRange("udp port range"),

			"tcp_port_ranges": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "TCP port ranges used to access the app.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"udp_port_ranges": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "UDP port ranges used to access the app.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"config_space": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the application.",
			},
			"domain_names": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of domains and IPs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"double_encrypt": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether Double Encryption is enabled or disabled for the app.",
			},
			"health_check_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"passive_health_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"health_reporting": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Whether health reporting for the app is Continuous or On Access. Supported values: NONE, ON_ACCESS, CONTINUOUS.",
			},
			"icmp_access_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PING_TRACEROUTING",
					"PING",
					"NONE",
				}, false),
			},
			"ip_anchored": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_cname_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the application.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"sra_apps": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"application_port": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"application_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"RDP",
								"SSH",
							}, false),
						},
						"domain": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"app_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"hidden": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"portal": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"connection_security": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ANY",
								"NLA",
								"NLA_EXT",
								"TLS",
								"VM_CONNECT",
								"RDP",
							}, false),
						},
					},
				},
			},
			"server_groups": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "List of the server group IDs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeList,
							Required: true,
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

func resourceAppSegmentPRACreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandSRAApplicationSegment(d)
	log.Printf("[INFO] Creating sra application segment request\n%+v\n", req)

	if req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provde a valid segment group for the sra application segment")
		return fmt.Errorf("please provde a valid segment group for the sra application segment")
	}

	resp, _, err := zClient.applicationsegmentpra.Create(req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Created sra application segment request. ID: %v\n", resp.ID)
	d.SetId(resp.ID)

	return resourceAppSegmentPRARead(d, m)
}

func resourceAppSegmentPRARead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.applicationsegmentpra.Get(d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing sra application segment %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	log.Printf("[INFO] Getting sra application segment:\n%+v\n", resp)
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

	if err := d.Set("sra_apps", flattenSRAApps(resp)); err != nil {
		return fmt.Errorf("failed to read sra application segment %s", err)
	}

	if err := d.Set("server_groups", flattenSRAAppServerGroups(resp.AppServerGroups)); err != nil {
		return fmt.Errorf("failed to read app server groups %s", err)
	}

	if err := d.Set("tcp_port_range", flattenNetworkPorts(resp.TCPAppPortRange)); err != nil {
		return err
	}

	if err := d.Set("tcp_port_range", flattenNetworkPorts(resp.UDPAppPortRange)); err != nil {
		return err
	}

	return nil

}

func resourceAppSegmentPRAUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating sra application segment ID: %v\n", id)
	req := expandSRAApplicationSegment(d)

	if d.HasChange("segment_group_id") && req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provde a valid segment group for the sra application segment")
		return fmt.Errorf("please provde a valid segment group for the sra application segment")
	}

	if _, err := zClient.applicationsegmentpra.Update(id, &req); err != nil {
		return err
	}

	return resourceAppSegmentPRARead(d, m)
}

func resourceAppSegmentPRADelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	id := d.Id()
	segmentGroupID, ok := d.GetOk("segment_group_id")
	if ok && segmentGroupID != nil {
		gID, ok := segmentGroupID.(string)
		if ok && gID != "" {
			// detach it from segment group first
			if err := detachSraPortalsFromGroup(zClient, id, gID); err != nil {
				return err
			}
		}
	}
	log.Printf("[INFO] Deleting sra application segment with id %v\n", id)
	if _, err := zClient.applicationsegmentpra.Delete(id); err != nil {
		return err
	}

	return nil
}

func detachSraPortalsFromGroup(client *Client, segmentID, segmentGroupID string) error {
	log.Printf("[INFO] Detaching sra application segment  %s from segment group: %s\n", segmentID, segmentGroupID)
	segGroup, _, err := client.segmentgroup.Get(segmentGroupID)
	if err != nil {
		log.Printf("[error] Error while getting segment group id: %s", segmentGroupID)
		return err
	}
	adaptedApplications := []segmentgroup.Application{}
	for _, app := range segGroup.Applications {
		if app.ID != segmentID {
			adaptedApplications = append(adaptedApplications, app)
		}
	}
	segGroup.Applications = adaptedApplications
	_, err = client.segmentgroup.Update(segmentGroupID, segGroup)
	return err

}

func expandSRAApplicationSegment(d *schema.ResourceData) applicationsegmentpra.AppSegmentSraApps {
	details := applicationsegmentpra.AppSegmentSraApps{
		SegmentGroupID:       d.Get("segment_group_id").(string),
		BypassType:           d.Get("bypass_type").(string),
		ConfigSpace:          d.Get("config_space").(string),
		PassiveHealthEnabled: d.Get("passive_health_enabled").(bool),
		IcmpAccessType:       d.Get("icmp_access_type").(string),
		Description:          d.Get("description").(string),
		DoubleEncrypt:        d.Get("double_encrypt").(bool),
		Enabled:              d.Get("enabled").(bool),
		HealthReporting:      d.Get("health_reporting").(string),
		HealthCheckType:      d.Get("health_check_type").(string),
		IPAnchored:           d.Get("ip_anchored").(bool),
		IsCnameEnabled:       d.Get("is_cname_enabled").(bool),
		DomainNames:          expandStringInSlice(d, "domain_names"),
		TCPPortRanges:        expandList(d.Get("tcp_port_ranges").([]interface{})),
		UDPPortRanges:        expandList(d.Get("udp_port_ranges").([]interface{})),
		AppServerGroups:      expandSRAAppServerGroups(d),
	}
	if d.HasChange("name") {
		details.Name = d.Get("name").(string)
	}
	if d.HasChange("segment_group_name") {
		details.SegmentGroupName = d.Get("segment_group_name").(string)
	}
	if d.HasChange("server_groups") {
		details.AppServerGroups = expandSRAAppServerGroups(d)
	}
	if d.HasChange("sra_apps") {
		details.SraApps = expandSRAApps(d)
	}
	TCPAppPortRange := expandNetwokPorts(d, "tcp_port_range")
	if TCPAppPortRange != nil {
		details.TCPAppPortRange = TCPAppPortRange
	}
	UDPAppPortRange := expandNetwokPorts(d, "udp_port_range")
	if UDPAppPortRange != nil {
		details.UDPAppPortRange = UDPAppPortRange
	}
	if d.HasChange("udp_port_ranges") {
		details.UDPPortRanges = convertToListString(d.Get("udp_port_ranges"))
	}
	if d.HasChange("tcp_port_ranges") {
		details.TCPPortRanges = convertToListString(d.Get("tcp_port_ranges"))
	}
	return details
}

func expandSRAApps(d *schema.ResourceData) []applicationsegmentpra.SraApps {
	sraInterface, ok := d.GetOk("sra_apps")
	if ok {
		sra := sraInterface.([]interface{})
		log.Printf("[INFO] sra apps data: %+v\n", sra)
		var sraApps []applicationsegmentpra.SraApps
		for _, sraApp := range sra {
			sraApp, ok := sraApp.(map[string]interface{})
			if ok {
				sraApps = append(sraApps, applicationsegmentpra.SraApps{
					Name:                sraApp["name"].(string),
					Description:         sraApp["description"].(string),
					Enabled:             sraApp["enabled"].(bool),
					ApplicationPort:     sraApp["application_port"].(string),
					ApplicationProtocol: sraApp["application_protocol"].(string),
					Domain:              sraApp["domain"].(string),
					AppID:               sraApp["app_id"].(string),
					Hidden:              sraApp["hidden"].(bool),
					Portal:              sraApp["portal"].(bool),
					ConnectionSecurity:  sraApp["connection_security"].(string),
				})
			}
		}
		return sraApps
	}

	return []applicationsegmentpra.SraApps{}
}

func expandSRAAppServerGroups(d *schema.ResourceData) []applicationsegmentpra.AppServerGroups {
	serverGroupsInterface, ok := d.GetOk("server_groups")
	if ok {
		serverGroup := serverGroupsInterface.(*schema.Set)
		log.Printf("[INFO] app server groups data: %+v\n", serverGroup)
		var serverGroups []applicationsegmentpra.AppServerGroups
		for _, appServerGroup := range serverGroup.List() {
			appServerGroup, _ := appServerGroup.(map[string]interface{})
			if appServerGroup != nil {
				for _, id := range appServerGroup["id"].([]interface{}) {
					serverGroups = append(serverGroups, applicationsegmentpra.AppServerGroups{
						ID: id.(string),
					})
				}
			}
		}
		return serverGroups
	}

	return []applicationsegmentpra.AppServerGroups{}
}
