package zpa

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontroller"
)

var policyRulesDetchLock sync.Mutex

func resourceApplicationSegment() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationSegmentCreate,
		Read:   resourceApplicationSegmentRead,
		Update: resourceApplicationSegmentUpdate,
		Delete: resourceApplicationSegmentDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				client := meta.(*Client)
				service := client.ApplicationSegment

				microTenantID := GetString(d.Get("microtenant_id"))
				if microTenantID != "" {
					service = service.WithMicroTenant(microTenantID)
				}

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := applicationsegment.GetByName(service, id)
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
			"match_style": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"EXCLUSIVE",
					"INCLUSIVE",
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
			"ip_anchored": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"select_connector_close_to_app": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
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
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "List of the server group IDs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeSet,
							Optional: true,
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

func resourceApplicationSegmentCreate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.ApplicationSegment

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandApplicationSegmentRequest(d, zClient, "")

	if err := validateAppPorts(req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return err
	}

	log.Printf("[INFO] Creating application segment request\n%+v\n", req)
	if req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provide a valid segment group for the application segment")
		return fmt.Errorf("please provide a valid segment group for the application segment")
	}

	resp, _, err := applicationsegment.Create(service, req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Created application segment request. ID: %v\n", resp.ID)
	d.SetId(resp.ID)

	return resourceApplicationSegmentRead(d, meta)
}

func resourceApplicationSegmentRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.ApplicationSegment

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := applicationsegment.Get(service, d.Id())
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
	_ = d.Set("match_style", resp.MatchStyle)
	_ = d.Set("microtenant_id", resp.MicroTenantID)
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
	_ = d.Set("server_groups", flattenAppServerGroupsSimple(resp.ServerGroups))

	if err := d.Set("tcp_port_range", flattenNetworkPorts(resp.TCPAppPortRange)); err != nil {
		return err
	}

	if err := d.Set("udp_port_range", flattenNetworkPorts(resp.UDPAppPortRange)); err != nil {
		return err
	}

	return nil
}

func flattenAppServerGroupsSimple(serverGroups []applicationsegment.AppServerGroups) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(serverGroups))
	for i, group := range serverGroups {
		ids[i] = group.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func resourceApplicationSegmentUpdate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	microTenantID := GetString(d.Get("microtenant_id"))

	id := d.Id()
	log.Printf("[INFO] Updating application segment ID: %v\n", id)
	req := expandApplicationSegmentRequest(d, zClient, id)

	if err := validateAppPorts(req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return err
	}

	if d.HasChange("segment_group_id") && req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provide a valid segment group for the application segment")
		return fmt.Errorf("please provide a valid segment group for the application segment")
	}

	service := zClient.ApplicationSegment
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	if _, _, err := applicationsegment.Get(service, id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := applicationsegment.Update(service, id, req); err != nil {
		return err
	}

	return resourceApplicationSegmentRead(d, meta)
}

func detachAppsFromAllPolicyRules(id string, policySetControllerService *services.Service) {
	policyRulesDetchLock.Lock()
	defer policyRulesDetchLock.Unlock()
	var rules []policysetcontroller.PolicyRule
	types := []string{"ACCESS_POLICY", "TIMEOUT_POLICY", "SIEM_POLICY", "CLIENT_FORWARDING_POLICY", "INSPECTION_POLICY"}
	for _, t := range types {
		policySet, _, err := policysetcontroller.GetByPolicyType(policySetControllerService, t)
		if err != nil {
			continue
		}
		r, _, err := policysetcontroller.GetAllByType(policySetControllerService, t)
		if err != nil {
			continue
		}
		for _, rule := range r {
			rule.PolicySetID = policySet.ID
			rules = append(rules, rule)
		}
	}
	log.Printf("[INFO] detachAppsFromAllPolicyRules Updating policy rules, len:%d \n", len(rules))
	for _, rr := range rules {
		rule := rr
		changed := false
		for i, condition := range rr.Conditions {
			operands := []policysetcontroller.Operands{}
			for _, op := range condition.Operands {
				if op.ObjectType == "APP" && op.LHS == "id" && op.RHS == id {
					changed = true
					continue
				}
				operands = append(operands, op)
			}
			rule.Conditions[i].Operands = operands
		}
		if len(rule.Conditions) == 0 {
			rule.Conditions = []policysetcontroller.Conditions{}
		}
		if changed {
			if _, err := policysetcontroller.UpdateRule(policySetControllerService, rule.PolicySetID, rule.ID, &rule); err != nil {
				continue
			}
		}
	}
}

func resourceApplicationSegmentDelete(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.ApplicationSegment.WithMicroTenant(GetString(d.Get("microtenant_id")))
	policySetControllerService := zClient.PolicySetController.WithMicroTenant(GetString(d.Get("microtenant_id")))
	id := d.Id()
	log.Printf("[INFO] Deleting application segment with id %v\n", id)
	detachAppsFromAllPolicyRules(id, policySetControllerService)
	if _, err := applicationsegment.Delete(service, id); err != nil {
		return err
	}

	return nil
}

func expandApplicationSegmentRequest(d *schema.ResourceData, client *Client, id string) applicationsegment.ApplicationSegmentResource {
	microTenantID := GetString(d.Get("microtenant_id"))
	service := client.ApplicationSegment
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

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
		MatchStyle:                d.Get("match_style").(string),
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
	if service != nil && id != "" {
		resource, _, err := applicationsegment.Get(service, id)
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
	serverGroupsInterface, ok := d.GetOk("server_groups")
	if ok {
		serverGroup := serverGroupsInterface.(*schema.Set)
		log.Printf("[INFO] app server groups data: %+v\n", serverGroup)
		var serverGroups []applicationsegment.AppServerGroups
		for _, appServerGroup := range serverGroup.List() {
			appServerGroup, _ := appServerGroup.(map[string]interface{})
			if ok {
				for _, id := range appServerGroup["id"].(*schema.Set).List() {
					serverGroups = append(serverGroups, applicationsegment.AppServerGroups{
						ID: id.(string),
					})
				}
			}
		}
		return serverGroups
	}

	return []applicationsegment.AppServerGroups{}
}
