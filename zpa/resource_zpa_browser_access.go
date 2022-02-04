package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/browseraccess"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/client"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/segmentgroup"
)

func resourceBrowserAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceBrowserAccessCreate,
		Read:   resourceBrowserAccessRead,
		Update: resourceBrowserAccessUpdate,
		Delete: resourceBrowserAccessDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					d.Set("id", id)
				} else {
					resp, _, err := zClient.browseraccess.GetByName(id)
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
				Default:     "NEVER",
				Description: "Indicates whether users can bypass ZPA to access applications. Default: NEVER. Supported values: ALWAYS, NEVER, ON_NET. The value NEVER indicates the use of the client forwarding policy.",
				ValidateFunc: validation.StringInSlice([]string{
					"ALWAYS",
					"NEVER",
					"ON_NET",
				}, false),
			},
			"tcp_port_range": resourceNetworkPortsSchema("tcp port range"),
			"udp_port_range": resourceNetworkPortsSchema("udp port range"),

			"tcp_port_ranges": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Deprecated:  "The tcp_port_ranges and udp_port_ranges fields are deprecated and replaced with tcp_port_range and udp_port_range.",
				Description: "TCP port ranges used to access the app.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"udp_port_ranges": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Deprecated:  "The tcp_port_ranges and udp_port_ranges fields are deprecated and replaced with tcp_port_range and udp_port_range.",
				Description: "UDP port ranges used to access the app.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"config_space": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DEFAULT",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the application.",
			},
			"domain_names": {
				Type:        schema.TypeList,
				Optional:    true,
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
			"ip_anchored": {
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
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_options": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_port": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"application_protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP",
								"HTTPS",
								"FTP",
								"RDP",
							}, false),
						},
						"certificate_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"certificate_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"cname": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"hidden": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"trust_untrusted_cert": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
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

func resourceBrowserAccessCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandBrowserAccess(d)
	log.Printf("[INFO] Creating browser access request\n%+v\n", req)

	if req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provde a valid segment group for the application segment")
		return fmt.Errorf("please provde a valid segment group for the application segment")
	}

	browseraccess, _, err := zClient.browseraccess.Create(req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Created browser access request. ID: %v\n", browseraccess.ID)
	d.SetId(browseraccess.ID)

	return resourceBrowserAccessRead(d, m)
}

func resourceBrowserAccessRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.browseraccess.Get(d.Id())
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

	if err := d.Set("udp_port_range", flattenNetworkPorts(resp.UDPAppPortRange)); err != nil {
		return err
	}

	return nil

}

func resourceBrowserAccessUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating browser access ID: %v\n", id)
	req := expandBrowserAccess(d)

	if d.HasChange("segment_group_id") && req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provde a valid segment group for the browser access application segment")
		return fmt.Errorf("please provde a valid segment group for the browser access application segment")
	}

	if _, err := zClient.browseraccess.Update(id, &req); err != nil {
		return err
	}

	return resourceBrowserAccessRead(d, m)
}

func resourceBrowserAccessDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	id := d.Id()
	segmentGroupID, ok := d.GetOk("segment_group_id")
	if ok && segmentGroupID != nil {
		gID, ok := segmentGroupID.(string)
		if ok && gID != "" {
			// detach it from segment group first
			if err := detachBrowserAccessFromGroup(zClient, id, gID); err != nil {
				return err
			}
		}
	}
	log.Printf("[INFO] Deleting browser access application with id %v\n", id)
	if _, err := zClient.browseraccess.Delete(id); err != nil {
		return err
	}

	return nil
}

func detachBrowserAccessFromGroup(client *Client, segmentID, segmentGroupID string) error {
	log.Printf("[INFO] Detaching browser access  %s from segment group: %s\n", segmentID, segmentGroupID)
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

func expandBrowserAccess(d *schema.ResourceData) browseraccess.BrowserAccess {
	details := browseraccess.BrowserAccess{
		BypassType:      d.Get("bypass_type").(string),
		Description:     d.Get("description").(string),
		DoubleEncrypt:   d.Get("double_encrypt").(bool),
		Enabled:         d.Get("enabled").(bool),
		HealthReporting: d.Get("health_reporting").(string),
		IPAnchored:      d.Get("ip_anchored").(bool),
		IsCnameEnabled:  d.Get("is_cname_enabled").(bool),
		DomainNames:     expandStringInSlice(d, "domain_names"),
		SegmentGroupID:  d.Get("segment_group_id").(string),
	}
	if d.HasChange("name") {
		details.Name = d.Get("name").(string)
	}
	if d.HasChange("segment_group_name") {
		details.SegmentGroupName = d.Get("segment_group_name").(string)
	}
	if d.HasChange("server_groups") {
		details.AppServerGroups = expandClientlessAppServerGroups(d)
	}
	if d.HasChange("clientless_apps") {
		details.ClientlessApps = expandClientlessApps(d)
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
					AppID:               clientlessApp["app_id"].(string),
					ApplicationPort:     clientlessApp["application_port"].(string),
					ApplicationProtocol: clientlessApp["application_protocol"].(string),
					CertificateID:       clientlessApp["certificate_id"].(string),
					CertificateName:     clientlessApp["certificate_name"].(string),
					Cname:               clientlessApp["cname"].(string),
					Description:         clientlessApp["description"].(string),
					Domain:              clientlessApp["domain"].(string),
					Enabled:             clientlessApp["enabled"].(bool),
					Hidden:              clientlessApp["hidden"].(bool),
					LocalDomain:         clientlessApp["local_domain"].(string),
					Name:                clientlessApp["name"].(string),
					Path:                clientlessApp["path"].(string),
					TrustUntrustedCert:  clientlessApp["trust_untrusted_cert"].(bool),
				})
			}
		}
		return clientlessApps
	}

	return []browseraccess.ClientlessApps{}
}

func expandClientlessAppServerGroups(d *schema.ResourceData) []browseraccess.AppServerGroups {
	serverGroupsInterface, ok := d.GetOk("server_groups")
	if ok {
		serverGroup := serverGroupsInterface.(*schema.Set)
		log.Printf("[INFO] app server groups data: %+v\n", serverGroup)
		var serverGroups []browseraccess.AppServerGroups
		for _, appServerGroup := range serverGroup.List() {
			appServerGroup, _ := appServerGroup.(map[string]interface{})
			if appServerGroup != nil {
				for _, id := range appServerGroup["id"].([]interface{}) {
					serverGroups = append(serverGroups, browseraccess.AppServerGroups{
						ID: id.(string),
					})
				}
			}
		}
		return serverGroups
	}

	return []browseraccess.AppServerGroups{}
}

func flattenBaClientlessApps(clientlessApp *browseraccess.BrowserAccess) []interface{} {
	clientlessApps := make([]interface{}, len(clientlessApp.ClientlessApps))
	for i, clientlessApp := range clientlessApp.ClientlessApps {
		clientlessApps[i] = map[string]interface{}{
			"allow_options":        clientlessApp.AllowOptions,
			"app_id":               clientlessApp.AppID,
			"application_port":     clientlessApp.ApplicationPort,
			"application_protocol": clientlessApp.ApplicationProtocol,
			"certificate_id":       clientlessApp.CertificateID,
			"certificate_name":     clientlessApp.CertificateName,
			"cname":                clientlessApp.Cname,
			"description":          clientlessApp.Description,
			"domain":               clientlessApp.Domain,
			"enabled":              clientlessApp.Enabled,
			"hidden":               clientlessApp.Hidden,
			"id":                   clientlessApp.ID,
			"local_domain":         clientlessApp.LocalDomain,
			"name":                 clientlessApp.Name,
			"path":                 clientlessApp.Path,
			"trust_untrusted_cert": clientlessApp.TrustUntrustedCert,
		}
	}

	return clientlessApps
}

func flattenClientlessAppServerGroups(appServerGroup []browseraccess.AppServerGroups) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(appServerGroup))
	for i, serverGroup := range appServerGroup {
		ids[i] = serverGroup.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}
