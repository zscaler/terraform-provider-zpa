package zpa

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/policysetcontroller"
)

func resourceApplicationSegment() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationSegmentCreate,
		Read:   resourceApplicationSegmentRead,
		Update: resourceApplicationSegmentUpdate,
		Delete: resourceApplicationSegmentDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := zClient.applicationsegment.GetByName(id)
					if err == nil {
						d.SetId(resp.ID)
						_ = d.Set("id", resp.ID)
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
				Optional: true,
				Computed: true,
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
				Description: "Indicates whether users can bypass ZPA to access applications.",
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
				Description: "TCP port ranges used to access the app.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"udp_port_ranges": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "UDP port ranges used to access the app.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"config_space": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DEFAULT",
					"SIEM",
				}, false),
				Default: "DEFAULT",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the application.",
			},
			"domain_names": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "List of domains and IPs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"double_encrypt": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether Double Encryption is enabled or disabled for the app.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether this application is enabled or not.",
			},
			"health_check_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{
					"DEFAULT",
					"NONE",
				}, false),
			},
			"health_reporting": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "NONE",
				Description: "Whether health reporting for the app is Continuous or On Access. Supported values: NONE, ON_ACCESS, CONTINUOUS.",
				ValidateFunc: validation.StringInSlice([]string{
					"NONE",
					"ON_ACCESS",
					"CONTINUOUS",
				}, false),
			},
			"icmp_access_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NONE",
				ValidateFunc: validation.StringInSlice([]string{
					"PING_TRACEROUTING",
					"PING",
					"NONE",
				}, false),
			},
			"default_idle_timeout": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_anchored": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// Need to implement validation.
			// When attribute is set to true only tcp_port_range and tcp_port_ranges are allowed
			// udp_port_range and udp_port_ranges are not allowed
			// Return error message: "The protocol configuration for the application is invalid. App Connector Closest to App supports only TCP applications."
			// Add validation to utils.go as it will be used across all application segments.
			"select_connector_close_to_app": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"use_in_dr_mode": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_incomplete_dr_config": {
				Type:     schema.TypeBool,
				Optional: true,
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
			"passive_health_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"tcp_keep_alive": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"0", "1",
				}, false),
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

func resourceApplicationSegmentCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandApplicationSegmentRequest(d, zClient, "")

	if err := validateAppPorts(zClient, req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return err
	}

	if err := checkForPortsOverlap(zClient, req); err != nil {
		return err
	}
	log.Printf("[INFO] Creating application segment request\n%+v\n", req)
	if req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provde a valid segment group for the application segment")
		return fmt.Errorf("please provde a valid segment group for the application segment")
	}
	resp, _, err := zClient.applicationsegment.Create(req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Created application segment request. ID: %v\n", resp.ID)
	d.SetId(resp.ID)

	return resourceApplicationSegmentRead(d, m)
}

func resourceApplicationSegmentRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.applicationsegment.Get(d.Id())

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing application segment %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Reading application segment and settings states: %+v\n", resp)
	_ = d.Set("id", resp.ID)
	_ = d.Set("segment_group_id", resp.SegmentGroupID)
	_ = d.Set("segment_group_name", resp.SegmentGroupName)
	_ = d.Set("bypass_type", resp.BypassType)
	_ = d.Set("config_space", resp.ConfigSpace)
	_ = d.Set("description", resp.Description)
	_ = d.Set("domain_names", resp.DomainNames)
	_ = d.Set("double_encrypt", resp.DoubleEncrypt)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("health_check_type", resp.HealthCheckType)
	_ = d.Set("health_reporting", resp.HealthReporting)
	_ = d.Set("icmp_access_type", resp.IcmpAccessType)
	_ = d.Set("ip_anchored", resp.IpAnchored)
	_ = d.Set("select_connector_close_to_app", resp.SelectConnectorCloseToApp)
	_ = d.Set("use_in_dr_mode", resp.UseInDrMode)
	_ = d.Set("is_incomplete_dr_config", resp.IsIncompleteDRConfig)
	_ = d.Set("is_cname_enabled", resp.IsCnameEnabled)
	_ = d.Set("tcp_keep_alive", resp.TCPKeepAlive)
	_ = d.Set("name", resp.Name)
	_ = d.Set("passive_health_enabled", resp.PassiveHealthEnabled)
	_ = d.Set("ip_anchored", resp.IpAnchored)
	_ = d.Set("tcp_port_ranges", convertPortsToListString(resp.TCPAppPortRange))
	_ = d.Set("udp_port_ranges", convertPortsToListString(resp.UDPAppPortRange))
	_ = d.Set("server_groups", flattenAppServerGroupsSimple(resp))

	if err := d.Set("tcp_port_range", flattenNetworkPorts(resp.TCPAppPortRange)); err != nil {
		return err
	}

	if err := d.Set("udp_port_range", flattenNetworkPorts(resp.UDPAppPortRange)); err != nil {
		return err
	}

	return nil
}
func flattenAppServerGroupsSimple(serverGroup *applicationsegment.ApplicationSegmentResource) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(serverGroup.ServerGroups))
	for i, group := range serverGroup.ServerGroups {
		ids[i] = group.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func resourceApplicationSegmentUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating application segment ID: %v\n", id)
	req := expandApplicationSegmentRequest(d, zClient, id)

	if err := validateAppPorts(zClient, req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return err
	}

	if err := checkForPortsOverlap(zClient, req); err != nil {
		return err
	}
	if d.HasChange("segment_group_id") && req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provide a valid segment group for the application segment")
		return fmt.Errorf("please provide a valid segment group for the application segment")
	}

	if _, _, err := zClient.applicationsegment.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := zClient.applicationsegment.Update(id, req); err != nil {
		return err
	}

	return resourceApplicationSegmentRead(d, m)
}

func detachApplicationSegmentFromAllPolicyRules(id string, zClient *Client) {
	var rules []policysetcontroller.PolicyRule
	types := []string{"ACCESS_POLICY", "TIMEOUT_POLICY", "SIEM_POLICY", "CLIENT_FORWARDING_POLICY", "INSPECTION_POLICY"}
	for _, t := range types {
		policySet, _, err := zClient.policysetcontroller.GetByPolicyType(t)
		if err != nil {
			continue
		}
		r, _, err := zClient.policysetcontroller.GetAllByType(t)
		if err != nil {
			continue
		}
		for _, rule := range r {
			rule.PolicySetID = policySet.ID
			rules = append(rules, rule)
		}
	}
	for _, rule := range rules {
		for i, condition := range rule.Conditions {
			var operands []policysetcontroller.Operands
			for _, op := range condition.Operands {
				if op.ObjectType == "APP" && op.LHS == "id" && op.RHS == id {
					continue
				}
				operands = append(operands, op)
			}
			rule.Conditions[i].Operands = operands
		}
		if _, err := zClient.policysetcontroller.Update(rule.PolicySetID, rule.ID, &rule); err != nil {
			continue
		}
	}
}

func resourceApplicationSegmentDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	id := d.Id()
	log.Printf("[INFO] Deleting application segment with id %v\n", id)

	detachApplicationSegmentFromAllPolicyRules(id, zClient)

	if _, err := zClient.applicationsegment.Delete(id); err != nil {
		return err
	}

	return nil
}

func expandApplicationSegmentRequest(d *schema.ResourceData, zClient *Client, id string) applicationsegment.ApplicationSegmentResource {
	details := applicationsegment.ApplicationSegmentResource{
		ID:                        d.Id(),
		Name:                      d.Get("name").(string),
		SegmentGroupID:            d.Get("segment_group_id").(string),
		SegmentGroupName:          d.Get("segment_group_name").(string),
		BypassType:                d.Get("bypass_type").(string),
		ConfigSpace:               d.Get("config_space").(string),
		IcmpAccessType:            d.Get("icmp_access_type").(string),
		Description:               d.Get("description").(string),
		DomainNames:               SetToStringList(d, "domain_names"),
		HealthCheckType:           d.Get("health_check_type").(string),
		HealthReporting:           d.Get("health_reporting").(string),
		TCPKeepAlive:              d.Get("tcp_keep_alive").(string),
		PassiveHealthEnabled:      d.Get("passive_health_enabled").(bool),
		DoubleEncrypt:             d.Get("double_encrypt").(bool),
		Enabled:                   d.Get("enabled").(bool),
		IpAnchored:                d.Get("ip_anchored").(bool),
		IsCnameEnabled:            d.Get("is_cname_enabled").(bool),
		SelectConnectorCloseToApp: d.Get("select_connector_close_to_app").(bool),
		UseInDrMode:               d.Get("use_in_dr_mode").(bool),
		IsIncompleteDRConfig:      d.Get("is_incomplete_dr_config").(bool),

		ServerGroups:    expandAppServerGroups(d),
		TCPAppPortRange: []common.NetworkPorts{},
		UDPAppPortRange: []common.NetworkPorts{},
	}
	remoteTCPAppPortRanges := []string{}
	remoteUDPAppPortRanges := []string{}
	if zClient != nil && id != "" {
		resource, _, err := zClient.applicationsegment.Get(id)
		if err == nil {
			remoteTCPAppPortRanges = resource.TCPPortRanges
			remoteUDPAppPortRanges = resource.UDPPortRanges
		}
	}
	TCPAppPortRange := expandAppSegmentNetwokPorts(d, "tcp_port_range")
	TCPAppPortRanges := convertToPortRange(d.Get("tcp_port_ranges").([]interface{}))
	if isSameSlice(TCPAppPortRange, TCPAppPortRanges) || isSameSlice(TCPAppPortRange, remoteTCPAppPortRanges) {
		details.TCPPortRanges = TCPAppPortRanges
	} else {
		details.TCPPortRanges = TCPAppPortRange
	}

	UDPAppPortRange := expandAppSegmentNetwokPorts(d, "udp_port_range")
	UDPAppPortRanges := convertToPortRange(d.Get("udp_port_ranges").([]interface{}))
	if isSameSlice(UDPAppPortRange, UDPAppPortRanges) || isSameSlice(UDPAppPortRange, remoteUDPAppPortRanges) {
		details.UDPPortRanges = UDPAppPortRanges
	} else {
		details.UDPPortRanges = UDPAppPortRange
	}

	if details.TCPPortRanges == nil {
		details.TCPPortRanges = []string{}
	}
	if details.UDPPortRanges == nil {
		details.UDPPortRanges = []string{}
	}
	return details
}

func expandAppServerGroups(d *schema.ResourceData) []applicationsegment.AppServerGroups {
	appServerGroupsInterface, ok := d.GetOk("server_groups")
	if ok {
		appServer := appServerGroupsInterface.(*schema.Set)
		log.Printf("[INFO] app server groups data: %+v\n", appServer)
		var appServerGroups []applicationsegment.AppServerGroups
		for _, appServerGroup := range appServer.List() {
			appServerGroup, _ := appServerGroup.(map[string]interface{})
			if appServerGroup != nil {
				for _, id := range appServerGroup["id"].([]interface{}) {
					appServerGroups = append(appServerGroups, applicationsegment.AppServerGroups{
						ID: id.(string),
					})
				}
			}
		}
		return appServerGroups
	}

	return []applicationsegment.AppServerGroups{}
}

func validateAppPorts(client *Client, selectConnectorCloseToApp bool, udpAppPortRange []common.NetworkPorts, udpPortRanges []string) error {
	if selectConnectorCloseToApp {
		if len(udpAppPortRange) > 0 || len(udpPortRanges) > 0 {
			return errors.New("the protocol configuration for the application is invalid. App Connector Closest to App supports only TCP applications")
		}
	}
	return nil

}

func checkForPortsOverlap(client *Client, app applicationsegment.ApplicationSegmentResource) error {
	time.Sleep(time.Second * time.Duration(rand.Intn(5)))
	apps, _, err := client.applicationsegment.GetAll()
	if err != nil {
		return err
	}
	for _, app2 := range apps {
		if found, common := sliceHasCommon(app.DomainNames, app2.DomainNames); found && app2.ID != app.ID && app2.Name != app.Name {
			// check for udp ports
			if overlap, o1, o2 := portOverlap(app.TCPPortRanges, app2.TCPPortRanges); overlap {
				return fmt.Errorf("found TCP overlapping ports: %v of application %s with %v of application %s (%s) with common domain name %s", o1, app.Name, o2, app2.Name, app2.ID, common)
			}
			if overlap, o1, o2 := portOverlap(app.UDPPortRanges, app2.UDPPortRanges); overlap {
				return fmt.Errorf("found UDP overlapping ports: %v of application %s with %v of application %s (%s) with common domain name %s", o1, app.Name, o2, app2.Name, app2.ID, common)
			}
		}
	}
	return nil

}

func portOverlap(s1, s2 []string) (bool, []string, []string) {
	for i1 := 0; i1 < len(s1); i1 += 2 {
		port1Start, _ := strconv.Atoi(s1[i1])
		port1End, _ := strconv.Atoi(s1[i1+1])
		port1Start, port1End = int(math.Min(float64(port1Start), float64(port1End))), int(math.Max(float64(port1Start), float64(port1End)))
		for i2 := 0; i2 < len(s2); i2 += 2 {
			port2Start, _ := strconv.Atoi(s2[i2])
			port2End, _ := strconv.Atoi(s2[i2+1])
			port2Start, port2End = int(math.Min(float64(port2Start), float64(port2End))), int(math.Max(float64(port2Start), float64(port2End)))
			if port1Start == port2Start || port1End == port2End || port1Start == port2End || port2Start == port1End {
				return true, s1[i1 : i1+2], s2[i2 : i2+2]
			}
			if port1Start < port2Start && port1End > port2Start {
				return true, s1[i1 : i1+2], s2[i2 : i2+2]
			}
			if port1End < port2End && port1End > port2Start {
				return true, s1[i1 : i1+2], s2[i2 : i2+2]
			}
			if port2Start < port1Start && port2End > port1Start {
				return true, s1[i1 : i1+2], s2[i2 : i2+2]
			}
			if port2End < port1End && port2End > port1Start {
				return true, s1[i1 : i1+2], s2[i2 : i2+2]
			}
		}
	}
	return false, nil, nil
}
