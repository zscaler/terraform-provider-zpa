package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/browseraccess"
)

func resourceApplicationSegmentBrowserAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationSegmentBrowserAccessCreate,
		Read:   resourceApplicationSegmentBrowserAccessRead,
		Update: resourceApplicationSegmentBrowserAccessUpdate,
		Delete: resourceApplicationSegmentBrowserAccessDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				client := meta.(*Client)
				service := client.AppServerController

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
					resp, _, err := browseraccess.GetByName(service, id)
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
				Default:  "DEFAULT",
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
							ForceNew:    true,
							Description: "Port for the BA app.",
						},
						"application_protocol": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Protocol for the BA app.",
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP",
								"HTTPS",
							}, false),
						},
						"certificate_id": {
							Type:        schema.TypeString,
							ForceNew:    true,
							Optional:    true,
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
					},
				},
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

func resourceApplicationSegmentBrowserAccessCreate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.BrowserAccess

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandBrowserAccess(d, zClient)

	if err := validateAppPorts(req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return err
	}

	log.Printf("[INFO] Creating browser access request\n%+v\n", req)

	if req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provide a valid segment group for the application segment")
		return fmt.Errorf("please provide a valid segment group for the application segment")
	}

	browseraccess, _, err := browseraccess.Create(service, req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Created browser access request. ID: %v\n", browseraccess.ID)
	d.SetId(browseraccess.ID)

	return resourceApplicationSegmentBrowserAccessRead(d, meta)
}

func resourceApplicationSegmentBrowserAccessRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.BrowserAccess

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := browseraccess.Get(service, d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing browser access %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return err
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
	_ = d.Set("tcp_keep_alive", resp.TCPKeepAlive)
	_ = d.Set("is_cname_enabled", resp.IsCnameEnabled)
	_ = d.Set("icmp_access_type", resp.ICMPAccessType)
	_ = d.Set("health_reporting", resp.HealthReporting)

	if err := d.Set("clientless_apps", flattenBaClientlessApps(resp)); err != nil {
		return fmt.Errorf("failed to read clientless apps %s", err)
	}

	if err := d.Set("server_groups", flattenCommonAppServerGroups(resp.AppServerGroups)); err != nil {
		return fmt.Errorf("failed to read app server groups %s", err)
	}

	if err := d.Set("tcp_port_ranges", resp.TCPPortRanges); err != nil {
		return err
	}
	if err := d.Set("udp_port_ranges", resp.UDPPortRanges); err != nil {
		return err
	}

	if err := d.Set("tcp_port_range", flattenNetworkPorts(resp.TCPAppPortRange)); err != nil {
		return err
	}

	if err := d.Set("udp_port_range", flattenNetworkPorts(resp.UDPAppPortRange)); err != nil {
		return err
	}

	return nil
}

func resourceApplicationSegmentBrowserAccessUpdate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.BrowserAccess

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating browser access ID: %v\n", id)
	req := expandBrowserAccess(d, zClient)

	if err := validateAppPorts(req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return err
	}

	if d.HasChange("segment_group_id") && req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provide a valid segment group for the browser access application segment")
		return fmt.Errorf("please provide a valid segment group for the browser access application segment")
	}

	if _, _, err := browseraccess.Get(service, id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := browseraccess.Update(service, id, &req); err != nil {
		return err
	}

	return resourceApplicationSegmentBrowserAccessRead(d, meta)
}

func resourceApplicationSegmentBrowserAccessDelete(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)

	service := zClient.BrowserAccess.WithMicroTenant(GetString(d.Get("microtenant_id")))

	id := d.Id()
	log.Printf("[INFO] Deleting browser access application with id %v\n", id)
	if _, err := browseraccess.Delete(service, id); err != nil {
		return err
	}

	return nil
}

func expandBrowserAccess(d *schema.ResourceData, zClient *Client) browseraccess.BrowserAccess {
	// Capture MicroTenantID for child tenant scenario
	microTenantID := GetString(d.Get("microtenant_id"))

	details := browseraccess.BrowserAccess{
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
		AppServerGroups:           expandCommonServerGroups(d),
		ClientlessApps:            expandClientlessApps(d),

		TCPPortRanges: ListToStringSlice(d.Get("tcp_port_ranges").([]interface{})),
		UDPPortRanges: ListToStringSlice(d.Get("udp_port_ranges").([]interface{})),
		// Use the helper for structured port ranges
		TCPAppPortRange: expandAppSegmentNetworkPorts(d, "tcp_port_range"),
		UDPAppPortRange: expandAppSegmentNetworkPorts(d, "udp_port_range"),
	}

	// // Optionally handle any further details if needed
	// if zClient != nil && id != "" {
	// 	microTenantID := GetString(d.Get("microtenant_id"))
	// 	service := zClient.BrowserAccess
	// 	if microTenantID != "" {
	// 		service = service.WithMicroTenant(microTenantID)
	// 	}
	// }
	// Apply microtenant context if available
	if microTenantID != "" {
		zClient.BrowserAccess = zClient.BrowserAccess.WithMicroTenant(microTenantID)
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

func expandClientlessApps(d *schema.ResourceData) []browseraccess.ClientlessApps {
	clientlessInterface, ok := d.GetOk("clientless_apps")
	if ok {
		clientless := clientlessInterface.([]interface{})
		log.Printf("[INFO] clientless apps data: %+v\n", clientless)
		var clientlessApps []browseraccess.ClientlessApps
		for _, clientlessApp := range clientless {
			clientlessApp, ok := clientlessApp.(map[string]interface{})
			if ok {
				clientlessApps = append(clientlessApps, browseraccess.ClientlessApps{
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
				})
			}
		}
		return clientlessApps
	}

	return []browseraccess.ClientlessApps{}
}

func flattenBaClientlessApps(clientlessApp *browseraccess.BrowserAccess) []interface{} {
	clientlessApps := make([]interface{}, len(clientlessApp.ClientlessApps))
	for i, clientlessApp := range clientlessApp.ClientlessApps {
		clientlessApps[i] = map[string]interface{}{
			"id":                   clientlessApp.ID,
			"name":                 clientlessApp.Name,
			"microtenant_id":       clientlessApp.MicroTenantID,
			"allow_options":        clientlessApp.AllowOptions,
			"application_port":     clientlessApp.ApplicationPort,
			"application_protocol": clientlessApp.ApplicationProtocol,
			"certificate_id":       clientlessApp.CertificateID,
			"description":          clientlessApp.Description,
			"domain":               clientlessApp.Domain,
			"enabled":              clientlessApp.Enabled,
			"trust_untrusted_cert": clientlessApp.TrustUntrustedCert,
		}
	}

	return clientlessApps
}
