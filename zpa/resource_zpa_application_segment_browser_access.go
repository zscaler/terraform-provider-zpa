package zpa

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentbrowseraccess"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
)

func resourceApplicationSegmentBrowserAccess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationSegmentBrowserAccessCreate,
		ReadContext:   resourceApplicationSegmentBrowserAccessRead,
		UpdateContext: resourceApplicationSegmentBrowserAccessUpdate,
		DeleteContext: resourceApplicationSegmentBrowserAccessDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			clientlessAppsRaw, ok := d.GetOk("clientless_apps")
			if !ok {
				return nil
			}

			clientlessApps, ok := clientlessAppsRaw.([]interface{})
			if !ok {
				return nil
			}

			for i, appRaw := range clientlessApps {
				app, ok := appRaw.(map[string]interface{})
				if !ok {
					continue
				}

				extLabel, hasExtLabel := app["ext_label"].(string)
				extDomain, hasExtDomain := app["ext_domain"].(string)
				certID, hasCertID := app["certificate_id"].(string)

				extFieldsSet := (hasExtLabel && extLabel != "") || (hasExtDomain && extDomain != "")
				certSet := hasCertID && certID != ""

				if extFieldsSet && certSet {
					return fmt.Errorf(
						"clientless_apps[%d]: 'certificate_id' cannot be set when either 'ext_label' or 'ext_domain' is configured",
						i,
					)
				}
			}

			return nil
		},

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				client := meta.(*Client)
				service := client.Service

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
					resp, _, err := applicationsegmentbrowseraccess.GetByName(ctx, service, id)
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
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
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
				Default:     "NEVER",
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
				Default:  "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{
					"DEFAULT",
					"SIEM",
				}, false),
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
			"health_check_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DEFAULT",
					"NONE",
				}, false),
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
			// "policy_style": {
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	Description: "Enable dual policy evaluation (resolve FQDN to Server IP and enforce policies based on Server IP and FQDN). false = NONE (disabled), true = DUAL_POLICY_EVAL (enabled).",
			// },
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
			"tcp_keep_alive": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"0", "1",
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
				Description: "Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the application.",
			},
			"clientless_apps": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"microtenant_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"allow_options": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If you want ZPA to forward unauthenticated HTTP preflight OPTIONS requests from the browser to the app.",
						},
						"application_port": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Port for the BA app.",
						},
						"application_protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol for the BA app.",
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP",
								"HTTPS",
							}, false),
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the BA certificate.",
						},
						"cname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the BA certificate.",
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Domain name or IP address of the BA app.",
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"trust_untrusted_cert": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether Use Untrusted Certificates is enabled or disabled for a BA app.",
						},
						"ext_label": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The domain prefix for the privileged portal URL. The supported string can include numbers, lower case characters, and only supports a hyphen (-).",
						},
						"ext_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The external domain name prefix of the Browser Access application that is used for Zscaler-managed certificates when creating a privileged portal.",
						},
					},
				},
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

func resourceApplicationSegmentBrowserAccessCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandBrowserAccess(ctx, d, zClient, "")

	if err := validateAppPorts(req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating browser access request\n%+v\n", req)

	if req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provide a valid segment group for the application segment")
		return diag.FromErr(fmt.Errorf("please provide a valid segment group for the application segment"))
	}

	browseraccess, _, err := applicationsegmentbrowseraccess.Create(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created browser access request. ID: %v\n", browseraccess.ID)
	d.SetId(browseraccess.ID)

	return resourceApplicationSegmentBrowserAccessRead(ctx, d, meta)
}

func resourceApplicationSegmentBrowserAccessRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := applicationsegmentbrowseraccess.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing browser access %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting browser access:\n%+v\n", resp)
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
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	_ = d.Set("select_connector_close_to_app", resp.SelectConnectorCloseToApp)
	_ = d.Set("use_in_dr_mode", resp.UseInDrMode)
	_ = d.Set("is_incomplete_dr_config", resp.IsIncompleteDRConfig)
	_ = d.Set("fqdn_dns_check", resp.FQDNDnsCheck)
	_ = d.Set("tcp_keep_alive", resp.TCPKeepAlive)
	_ = d.Set("is_cname_enabled", resp.IsCnameEnabled)
	_ = d.Set("icmp_access_type", resp.ICMPAccessType)
	_ = d.Set("health_reporting", resp.HealthReporting)
	// _ = d.Set("policy_style", PolicyStyleAPIToBool(resp.PolicyStyle))
	_ = d.Set("zpn_er_id", flattenCommonZPNERIDSimple(resp.ZPNERID))

	if err := d.Set("clientless_apps", flattenBaClientlessApps(resp)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read clientless apps %s", err))
	}

	if err := d.Set("server_groups", flattenCommonAppServerGroupSimple(resp.AppServerGroups)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read app server groups %s", err))
	}

	_ = d.Set("tcp_port_ranges", convertPortsToListString(resp.TCPAppPortRange))
	_ = d.Set("udp_port_ranges", convertPortsToListString(resp.UDPAppPortRange))

	if err := d.Set("tcp_port_range", flattenNetworkPorts(resp.TCPAppPortRange)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("udp_port_range", flattenNetworkPorts(resp.UDPAppPortRange)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceApplicationSegmentBrowserAccessUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating browser access ID: %v\n", id)

	// Step 1: Retrieve existing configuration to get app_id and clientless_apps.id
	existingSegment, _, err := applicationsegmentbrowseraccess.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// Step 2: Build the update payload
	req := expandBrowserAccess(ctx, d, zClient, "")

	// Step 3: Inject app_id and clientless_apps.id from the existing configuration
	req.ID = existingSegment.ID // Assign app_id to the parent application
	for i, clientlessApp := range req.ClientlessApps {
		if i < len(existingSegment.ClientlessApps) {
			clientlessApp.ID = existingSegment.ClientlessApps[i].ID // Existing clientless_app id
			clientlessApp.AppID = existingSegment.ID                // Assign parent app_id to clientless_app
			req.ClientlessApps[i] = clientlessApp
		}
	}

	// Validate the update request
	if err := validateAppPorts(req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return diag.FromErr(err)
	}

	// Step 4: Ensure segment_group_id is valid if changed
	if d.HasChange("segment_group_id") && req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provide a valid segment group for the browser access application segment")
		return diag.FromErr(fmt.Errorf("please provide a valid segment group for the browser access application segment"))
	}

	// Step 5: Perform the update
	if _, err := applicationsegmentbrowseraccess.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	// Step 6: Refresh the state after updating
	return resourceApplicationSegmentBrowserAccessRead(ctx, d, meta)
}

func resourceApplicationSegmentBrowserAccessDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Deleting browser access application with id %v\n", id)
	if _, err := applicationsegmentbrowseraccess.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[INFO] browser access application deleted successfully")
	return nil
}

func expandBrowserAccess(ctx context.Context, d *schema.ResourceData, zClient *Client, id string) applicationsegmentbrowseraccess.BrowserAccess {
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
	details := applicationsegmentbrowseraccess.BrowserAccess{
		ID:                        d.Id(),
		Name:                      d.Get("name").(string),
		SegmentGroupID:            d.Get("segment_group_id").(string),
		SegmentGroupName:          d.Get("segment_group_name").(string),
		BypassType:                d.Get("bypass_type").(string),
		ConfigSpace:               d.Get("config_space").(string),
		ICMPAccessType:            d.Get("icmp_access_type").(string),
		Description:               d.Get("description").(string),
		MicroTenantID:             d.Get("microtenant_id").(string),
		DomainNames:               SetToStringList(d, "domain_names"),
		HealthCheckType:           d.Get("health_check_type").(string),
		HealthReporting:           d.Get("health_reporting").(string),
		TCPKeepAlive:              d.Get("tcp_keep_alive").(string),
		DoubleEncrypt:             d.Get("double_encrypt").(bool),
		Enabled:                   d.Get("enabled").(bool),
		PassiveHealthEnabled:      d.Get("passive_health_enabled").(bool),
		IPAnchored:                d.Get("ip_anchored").(bool),
		IsCnameEnabled:            d.Get("is_cname_enabled").(bool),
		SelectConnectorCloseToApp: d.Get("select_connector_close_to_app").(bool),
		UseInDrMode:               d.Get("use_in_dr_mode").(bool),
		IsIncompleteDRConfig:      d.Get("is_incomplete_dr_config").(bool),
		FQDNDnsCheck:              d.Get("fqdn_dns_check").(bool),
		// PolicyStyle:               PolicyStyleBoolToAPIString(GetBool(d.Get("policy_style"))),
		ZPNERID: extranet,
		AppServerGroups: func() []servergroup.ServerGroup {
			groups := expandCommonServerGroups(d)
			if groups == nil {
				return []servergroup.ServerGroup{}
			}
			return groups
		}(),
		ClientlessApps: expandClientlessApps(d),

		TCPAppPortRange: []common.NetworkPorts{},
		UDPAppPortRange: []common.NetworkPorts{},
	}
	remoteTCPAppPortRanges := []string{}
	remoteUDPAppPortRanges := []string{}
	if service != nil && id != "" {
		resource, _, err := applicationsegmentbrowseraccess.Get(ctx, service, id)
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
	// Handle specific changes, to be sure we're updating the correct fields
	if d.HasChange("name") {
		details.Name = d.Get("name").(string)
	}

	if d.HasChange("clientless_apps") {
		details.ClientlessApps = expandClientlessApps(d)
	}

	return details
}

func expandClientlessApps(d *schema.ResourceData) []applicationsegmentbrowseraccess.ClientlessApps {
	clientlessInterface, ok := d.GetOk("clientless_apps")
	if ok {
		clientless := clientlessInterface.([]interface{})
		log.Printf("[INFO] clientless apps data: %+v\n", clientless)
		var clientlessApps []applicationsegmentbrowseraccess.ClientlessApps
		for _, clientlessApp := range clientless {
			clientlessApp, ok := clientlessApp.(map[string]interface{})
			if ok {
				clientlessApps = append(clientlessApps, applicationsegmentbrowseraccess.ClientlessApps{
					ID:                  clientlessApp["id"].(string),
					AppID:               clientlessApp["app_id"].(string),
					AllowOptions:        clientlessApp["allow_options"].(bool),
					ApplicationPort:     clientlessApp["application_port"].(string),
					ApplicationProtocol: clientlessApp["application_protocol"].(string),
					CertificateID:       clientlessApp["certificate_id"].(string),
					Description:         clientlessApp["description"].(string),
					Domain:              clientlessApp["domain"].(string),
					Enabled:             clientlessApp["enabled"].(bool),
					Name:                clientlessApp["name"].(string),
					MicroTenantID:       clientlessApp["microtenant_id"].(string),
					TrustUntrustedCert:  clientlessApp["trust_untrusted_cert"].(bool),
					ExtLabel:            clientlessApp["ext_label"].(string),
					ExtDomain:           clientlessApp["ext_domain"].(string),
				})
			}
		}
		return clientlessApps
	}

	return []applicationsegmentbrowseraccess.ClientlessApps{}
}

func flattenBaClientlessApps(clientlessApp *applicationsegmentbrowseraccess.BrowserAccess) []interface{} {
	clientlessApps := make([]interface{}, len(clientlessApp.ClientlessApps))
	for i, clientlessApp := range clientlessApp.ClientlessApps {
		clientlessApps[i] = map[string]interface{}{
			"id":                   clientlessApp.ID,
			"app_id":               clientlessApp.AppID,
			"name":                 clientlessApp.Name,
			"description":          clientlessApp.Description,
			"cname":                clientlessApp.Cname,
			"microtenant_id":       clientlessApp.MicroTenantID,
			"allow_options":        clientlessApp.AllowOptions,
			"application_port":     clientlessApp.ApplicationPort,
			"application_protocol": clientlessApp.ApplicationProtocol,
			"certificate_id":       clientlessApp.CertificateID,
			"domain":               clientlessApp.Domain,
			"enabled":              clientlessApp.Enabled,
			"trust_untrusted_cert": clientlessApp.TrustUntrustedCert,
			"ext_label":            clientlessApp.ExtLabel,
			"ext_domain":           clientlessApp.ExtDomain,
		}
	}

	return clientlessApps
}
