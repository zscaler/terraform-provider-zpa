package zpa

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment_share"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
)

var policyRulesDetchLock sync.Mutex

func resourceApplicationSegment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationSegmentCreate,
		ReadContext:   resourceApplicationSegmentRead,
		UpdateContext: resourceApplicationSegmentUpdate,
		DeleteContext: resourceApplicationSegmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

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
					resp, _, err := applicationsegment.GetByName(ctx, service, id)
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
			"bypass_on_reauth": {
				Type:     schema.TypeBool,
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
			"share_to_microtenants": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Share the Application Segment to microtenants",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether this application is enabled or not.",
			},
			"inspect_traffic_with_zia": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if Inspect Traffic with ZIA is enabled for the application.",
			},
			"health_check_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"fqdn_dns_check": {
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
			"api_protection_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, designates the application segment for API traffic inspection",
			},
			"server_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"zpn_er_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceApplicationSegmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandApplicationSegmentRequest(ctx, d, zClient, "")

	if err := validateAppPorts(req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating application segment request\n%+v\n", req)
	if req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provide a valid segment group for the application segment")
		return diag.FromErr(fmt.Errorf("please provide a valid segment group for the application segment"))
	}

	resp, _, err := applicationsegment.Create(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created application segment request. ID: %v\n", resp.ID)
	d.SetId(resp.ID)

	// 🔽 ADD THIS BLOCK RIGHT AFTER d.SetId(...)
	shareTo := SetToStringList(d, "share_to_microtenants")
	if len(shareTo) > 0 {
		log.Printf("[INFO] Sharing application segment %s to microtenants: %v", resp.ID, shareTo)
		shareReq := applicationsegment_share.AppSegmentSharedToMicrotenant{
			ApplicationID:       resp.ID,
			ShareToMicrotenants: shareTo,
			MicroTenantID:       microTenantID,
		}
		if _, err := applicationsegment_share.AppSegmentMicrotenantShare(ctx, service, resp.ID, shareReq); err != nil {
			return diag.FromErr(fmt.Errorf("failed to share application segment to microtenants: %w", err))
		}
	}

	return resourceApplicationSegmentRead(ctx, d, meta)
}

func resourceApplicationSegmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := applicationsegment.Get(ctx, service, d.Id())
	if err != nil {
		if err.(*errorx.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing application segment %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading application segment and settings states: %+v\n", resp)
	_ = d.Set("id", resp.ID)
	_ = d.Set("segment_group_id", resp.SegmentGroupID)
	_ = d.Set("segment_group_name", resp.SegmentGroupName)
	_ = d.Set("bypass_type", resp.BypassType)
	_ = d.Set("bypass_on_reauth", resp.BypassOnReauth)
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
	_ = d.Set("inspect_traffic_with_zia", resp.InspectTrafficWithZia)
	_ = d.Set("name", resp.Name)
	_ = d.Set("passive_health_enabled", resp.PassiveHealthEnabled)
	_ = d.Set("fqdn_dns_check", resp.FQDNDnsCheck)
	_ = d.Set("api_protection_enabled", resp.APIProtectionEnabled)
	_ = d.Set("tcp_port_ranges", convertPortsToListString(resp.TCPAppPortRange))
	_ = d.Set("udp_port_ranges", convertPortsToListString(resp.UDPAppPortRange))
	_ = d.Set("server_groups", flattenCommonAppServerGroupSimple(resp.ServerGroups))
	_ = d.Set("zpn_er_id", flattenCommonZPNERIDSimple(resp.ZPNERID))

	shareTo := []string{}
	if len(resp.SharedMicrotenantDetails.SharedToMicrotenants) > 0 {
		for _, smt := range resp.SharedMicrotenantDetails.SharedToMicrotenants {
			shareTo = append(shareTo, smt.ID)
		}
	}
	_ = d.Set("share_to_microtenants", shareTo)

	if err := d.Set("tcp_port_range", flattenNetworkPorts(resp.TCPAppPortRange)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("udp_port_range", flattenNetworkPorts(resp.UDPAppPortRange)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceApplicationSegmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating application segment ID: %v\n", id)
	req := expandApplicationSegmentRequest(ctx, d, zClient, id)

	if err := validateAppPorts(req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("segment_group_id") && req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provide a valid segment group for the application segment")
		return diag.FromErr(fmt.Errorf("please provide a valid segment group for the application segment"))
	}

	if _, _, err := applicationsegment.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := applicationsegment.Update(ctx, service, id, req); err != nil {
		return diag.FromErr(err)
	}

	// Share if needed
	shareTo := SetToStringList(d, "share_to_microtenants")
	if len(shareTo) > 0 {
		log.Printf("[INFO] Sharing updated application segment %s to microtenants: %v", id, shareTo)
		shareReq := applicationsegment_share.AppSegmentSharedToMicrotenant{
			ApplicationID:       id,
			ShareToMicrotenants: shareTo,
			MicroTenantID:       microTenantID,
		}
		if _, err := applicationsegment_share.AppSegmentMicrotenantShare(ctx, service, id, shareReq); err != nil {
			return diag.FromErr(fmt.Errorf("failed to share updated application segment to microtenants: %w", err))
		}
	}

	return resourceApplicationSegmentRead(ctx, d, meta)
}

func resourceApplicationSegmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Use MicroTenant if available
	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting application segment with id %v\n", d.Id())

	// Pass d.Id() as a string to the detachAppsFromAllPolicyRules function
	detachAppsFromAllPolicyRules(ctx, d.Id(), service)

	if _, err := applicationsegment.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[INFO] Application segment deleted successfully")
	return nil
}

func expandApplicationSegmentRequest(ctx context.Context, d *schema.ResourceData, zClient *Client, id string) applicationsegment.ApplicationSegmentResource {
	microTenantID := GetString(d.Get("microtenant_id"))
	service := zClient.Service // Unified service interface
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Expand zpn_er_id (extranet resource ID)
	var extranet *common.ZPNERID
	if v, ok := d.GetOk("zpn_er_id"); ok {
		if items := v.([]interface{}); len(items) > 0 {
			m := items[0].(map[string]interface{})
			if idSet, ok := m["id"].(*schema.Set); ok && idSet.Len() > 0 {
				ids := idSet.List()
				if len(ids) > 0 {
					if id, ok := ids[0].(string); ok && id != "" {
						extranet = &common.ZPNERID{ID: id}
					}
				}
			}
		}
	}

	details := applicationsegment.ApplicationSegmentResource{
		ID:                        d.Id(),
		Name:                      d.Get("name").(string),
		SegmentGroupID:            d.Get("segment_group_id").(string),
		SegmentGroupName:          d.Get("segment_group_name").(string),
		BypassType:                d.Get("bypass_type").(string),
		BypassOnReauth:            d.Get("bypass_on_reauth").(bool),
		ConfigSpace:               d.Get("config_space").(string),
		IcmpAccessType:            d.Get("icmp_access_type").(string),
		Description:               d.Get("description").(string),
		DomainNames:               SetToStringList(d, "domain_names"),
		HealthCheckType:           d.Get("health_check_type").(string),
		MatchStyle:                d.Get("match_style").(string),
		HealthReporting:           d.Get("health_reporting").(string),
		TCPKeepAlive:              d.Get("tcp_keep_alive").(string),
		MicroTenantID:             d.Get("microtenant_id").(string),
		ShareToMicrotenants:       SetToStringList(d, "share_to_microtenants"),
		PassiveHealthEnabled:      d.Get("passive_health_enabled").(bool),
		InspectTrafficWithZia:     d.Get("inspect_traffic_with_zia").(bool),
		DoubleEncrypt:             d.Get("double_encrypt").(bool),
		Enabled:                   d.Get("enabled").(bool),
		IpAnchored:                d.Get("ip_anchored").(bool),
		IsCnameEnabled:            d.Get("is_cname_enabled").(bool),
		SelectConnectorCloseToApp: d.Get("select_connector_close_to_app").(bool),
		UseInDrMode:               d.Get("use_in_dr_mode").(bool),
		IsIncompleteDRConfig:      d.Get("is_incomplete_dr_config").(bool),
		FQDNDnsCheck:              d.Get("fqdn_dns_check").(bool),
		APIProtectionEnabled:      d.Get("api_protection_enabled").(bool),
		ZPNERID:                   extranet,
		ServerGroups: func() []servergroup.ServerGroup {
			groups := expandCommonServerGroups(d)
			if groups == nil {
				return []servergroup.ServerGroup{}
			}
			return groups
		}(),
		// ServerGroups:    expandCommonServerGroups(d),
		TCPAppPortRange: []common.NetworkPorts{},
		UDPAppPortRange: []common.NetworkPorts{},
	}
	remoteTCPAppPortRanges := []string{}
	remoteUDPAppPortRanges := []string{}
	if service != nil && id != "" {
		resource, _, err := applicationsegment.Get(ctx, service, id)
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

func detachAppsFromAllPolicyRules(ctx context.Context, id string, policySetControllerService *zscaler.Service) {
	policyRulesDetchLock.Lock()
	defer policyRulesDetchLock.Unlock()
	var rules []policysetcontroller.PolicyRule
	types := []string{"ACCESS_POLICY", "TIMEOUT_POLICY", "SIEM_POLICY", "CLIENT_FORWARDING_POLICY", "INSPECTION_POLICY"}
	for _, t := range types {
		policySet, _, err := policysetcontroller.GetByPolicyType(ctx, policySetControllerService, t)
		if err != nil {
			continue
		}
		r, _, err := policysetcontroller.GetAllByType(ctx, policySetControllerService, t)
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
			if _, err := policysetcontroller.UpdateRule(ctx, policySetControllerService, rule.PolicySetID, rule.ID, &rule); err != nil {
				continue
			}
		}
	}
}
