package zpa

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegmentpra"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

func resourceApplicationSegmentPRA() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationSegmentPRACreate,
		Read:   resourceApplicationSegmentPRARead,
		Update: resourceApplicationSegmentPRAUpdate,
		Delete: resourceApplicationSegmentPRADelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				client := meta.(*Client)
				service := client.ApplicationSegmentPRA

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					d.Set("id", id)
				} else {
					resp, _, err := applicationsegmentpra.GetByName(service, id)
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the application.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"segment_group_id": {
				Type:     schema.TypeString,
				Required: true,
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
				Computed:    true,
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
			"select_connector_close_to_app": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"use_in_dr_mode": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"is_incomplete_dr_config": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"is_cname_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors.",
			},
			"tcp_keep_alive": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"0", "1",
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
			"common_apps_dto": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"apps_config": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"app_types": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"application_port": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"application_protocol": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"RDP", "SSH", "VNC",
										}, false),
									},
									"connection_security": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"ANY", "NLA", "NLA_EXT", "TLS", "VM_CONNECT", "RDP",
										}, false),
									},
									"domain": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
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
		CustomizeDiff: customizeDiffApplicationSegmentPRA,
	}
}

func resourceApplicationSegmentPRACreate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.ApplicationSegmentPRA

	req := expandSRAApplicationSegment(d, zClient, "")
	if err := checkForPRAPortsOverlap(zClient, req); err != nil {
		return err
	}

	if err := validateAppPorts(req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return err
	}

	log.Printf("[INFO] Creating application segment request\n%+v\n", req)
	resp, _, err := applicationsegmentpra.Create(service, req)
	if err != nil {
		log.Printf("[ERROR] Failed to create application segment: %s", err)
		return err
	}

	log.Printf("[INFO] Created application segment request. ID: %v\n", resp.ID)
	d.SetId(resp.ID)

	// Introduce a brief delay to allow the backend to fully process the creation
	//lintignore:R018
	time.Sleep(5 * time.Second)

	// Explicitly call GET using the ID to fetch the latest resource state
	_, _, err = applicationsegmentpra.Get(service, resp.ID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch application segment after creation: %s", err)
		return err
	}

	// Now, update Terraform state with the latest fetched details
	return resourceApplicationSegmentPRARead(d, meta)
}

func resourceApplicationSegmentPRARead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.ApplicationSegmentPRA

	resp, _, err := applicationsegmentpra.Get(service, d.Id())
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
	_ = d.Set("icmp_access_type", resp.IcmpAccessType)
	_ = d.Set("match_style", resp.MatchStyle)
	_ = d.Set("select_connector_close_to_app", resp.SelectConnectorCloseToApp)
	_ = d.Set("use_in_dr_mode", resp.UseInDrMode)
	_ = d.Set("is_incomplete_dr_config", resp.IsIncompleteDRConfig)
	_ = d.Set("tcp_keep_alive", resp.TCPKeepAlive)
	_ = d.Set("ip_anchored", resp.IpAnchored)
	_ = d.Set("health_reporting", resp.HealthReporting)
	_ = d.Set("tcp_port_ranges", convertPortsToListString(resp.TCPAppPortRange))
	_ = d.Set("udp_port_ranges", convertPortsToListString(resp.UDPAppPortRange))
	_ = d.Set("server_groups", flattenPRAAppServerGroupsSimple(resp.ServerGroups))

	if err := d.Set("common_apps_dto", flattenCommonAppsDto(resp.PRAApps)); err != nil {
		return fmt.Errorf("failed to read common application in application segment %s", err)
	}

	if err := d.Set("tcp_port_range", flattenNetworkPorts(resp.TCPAppPortRange)); err != nil {
		return err
	}

	if err := d.Set("udp_port_range", flattenNetworkPorts(resp.UDPAppPortRange)); err != nil {
		return err
	}
	return nil
}

func flattenPRAAppServerGroupsSimple(serverGroup []applicationsegmentpra.AppServerGroups) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(serverGroup))
	for i, group := range serverGroup {
		ids[i] = group.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func resourceApplicationSegmentPRAUpdate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.ApplicationSegmentPRA

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating pra application segment ID: %v\n", id)
	req := expandSRAApplicationSegment(d, zClient, id)

	if err := checkForPRAPortsOverlap(zClient, req); err != nil {
		return err
	}

	if err := validateAppPorts(req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return err
	}

	if _, _, err := applicationsegmentpra.Get(service, id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	// Perform the update
	_, err := applicationsegmentpra.Update(service, id, &req)
	if err != nil {
		return err
	}

	// Introduce a brief delay to allow the backend to fully process the update
	//lintignore:R018
	time.Sleep(5 * time.Second)

	// Fetch the latest resource state after the update
	_, _, err = applicationsegmentpra.Get(service, id)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch application segment after update: %s", err)
		return err
	}

	// Now, update Terraform state with the latest fetched details
	return resourceApplicationSegmentPRARead(d, meta)
}

func resourceApplicationSegmentPRADelete(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.ApplicationSegmentPRA

	id := d.Id()
	segmentGroupID, ok := d.GetOk("segment_group_id")
	if ok && segmentGroupID != nil {
		gID, ok := segmentGroupID.(string)
		if ok && gID != "" {
			// detach it from segment group first
			if err := detachSegmentGroup(zClient, id, gID); err != nil {
				return err
			}
		}
	}
	log.Printf("[INFO] Deleting pra application segment with id %v\n", id)
	if _, err := applicationsegmentpra.Delete(service, id); err != nil {
		return err
	}

	return nil
}

func expandSRAApplicationSegment(d *schema.ResourceData, zClient *Client, id string) applicationsegmentpra.AppSegmentPRA {
	details := applicationsegmentpra.AppSegmentPRA{
		ID:                        d.Id(),
		SegmentGroupID:            d.Get("segment_group_id").(string),
		BypassType:                d.Get("bypass_type").(string),
		ConfigSpace:               d.Get("config_space").(string),
		IcmpAccessType:            d.Get("icmp_access_type").(string),
		Description:               d.Get("description").(string),
		HealthReporting:           d.Get("health_reporting").(string),
		HealthCheckType:           d.Get("health_check_type").(string),
		PassiveHealthEnabled:      d.Get("passive_health_enabled").(bool),
		DoubleEncrypt:             d.Get("double_encrypt").(bool),
		Enabled:                   d.Get("enabled").(bool),
		IpAnchored:                d.Get("ip_anchored").(bool),
		IsCnameEnabled:            d.Get("is_cname_enabled").(bool),
		SelectConnectorCloseToApp: d.Get("select_connector_close_to_app").(bool),
		UseInDrMode:               d.Get("use_in_dr_mode").(bool),
		TCPKeepAlive:              d.Get("tcp_keep_alive").(string),
		MatchStyle:                d.Get("match_style").(string),
		IsIncompleteDRConfig:      d.Get("is_incomplete_dr_config").(bool),
		DomainNames:               SetToStringList(d, "domain_names"),
		TCPAppPortRange:           []common.NetworkPorts{},
		UDPAppPortRange:           []common.NetworkPorts{},
		ServerGroups:              expandPRAAppServerGroups(d),
		CommonAppsDto:             expandCommonAppsDto(d),
	}

	if d.HasChange("name") {
		details.Name = d.Get("name").(string)
	}
	if d.HasChange("server_groups") {
		details.ServerGroups = expandPRAAppServerGroups(d)
	}

	remoteTCPAppPortRanges := []string{}
	remoteUDPAppPortRanges := []string{}
	if zClient != nil && id != "" {
		microTenantID := GetString(d.Get("microtenant_id"))
		service := zClient.ApplicationSegmentPRA
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		resource, _, err := applicationsegmentpra.Get(service, id)
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

func expandCommonAppsDto(d *schema.ResourceData) applicationsegmentpra.CommonAppsDto {
	result := applicationsegmentpra.CommonAppsDto{}
	if commonAppsInterface, ok := d.GetOk("common_apps_dto"); ok {
		commonAppsList := commonAppsInterface.(*schema.Set).List()
		if len(commonAppsList) > 0 {
			commonAppMap := commonAppsList[0].(map[string]interface{})
			result.AppsConfig = expandAppsConfig(commonAppMap["apps_config"])
		}
	}
	return result
}

func expandAppsConfig(appsConfigInterface interface{}) []applicationsegmentpra.AppsConfig {
	appsConfig, ok := appsConfigInterface.(*schema.Set)
	if !ok {
		return []applicationsegmentpra.AppsConfig{}
	}
	log.Printf("[INFO] apps config data: %+v\n", appsConfig)
	var commonAppConfigDto []applicationsegmentpra.AppsConfig
	for _, commonAppConfig := range appsConfig.List() {
		commonAppConfig, ok := commonAppConfig.(map[string]interface{})
		if ok {
			appTypesSet, ok := commonAppConfig["app_types"].(*schema.Set)
			if !ok {
				continue
			}
			appTypes := SetToStringSlice(appTypesSet)
			commonAppConfigDto = append(commonAppConfigDto, applicationsegmentpra.AppsConfig{
				ID:                  commonAppConfig["id"].(string),
				Name:                commonAppConfig["name"].(string),
				Enabled:             commonAppConfig["enabled"].(bool),
				Domain:              commonAppConfig["domain"].(string),
				ApplicationPort:     commonAppConfig["application_port"].(string),
				ApplicationProtocol: commonAppConfig["application_protocol"].(string),
				ConnectionSecurity:  commonAppConfig["connection_security"].(string),
				AppTypes:            appTypes,
			})
		}
	}
	return commonAppConfigDto
}

func expandPRAAppServerGroups(d *schema.ResourceData) []applicationsegmentpra.AppServerGroups {
	serverGroupsInterface, ok := d.GetOk("server_groups")
	if ok {
		serverGroup := serverGroupsInterface.(*schema.Set)
		log.Printf("[INFO] app server groups data: %+v\n", serverGroup)
		var serverGroups []applicationsegmentpra.AppServerGroups
		for _, appServerGroup := range serverGroup.List() {
			appServerGroup, _ := appServerGroup.(map[string]interface{})
			if ok {
				for _, id := range appServerGroup["id"].(*schema.Set).List() {
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

func flattenCommonAppsDto(apps []applicationsegmentpra.PRAApps) []interface{} {
	commonAppsDto := make([]interface{}, 1)
	appsConfig := make([]interface{}, len(apps))
	for i, app := range apps {
		appTypes := []string{}
		if app.ApplicationProtocol == "RDP" || app.ApplicationProtocol == "SSH" || app.ApplicationProtocol == "VNC" {
			appTypes = append(appTypes, "SECURE_REMOTE_ACCESS")
		}
		appConfigMap := map[string]interface{}{
			"id":                   app.ID,
			"name":                 app.Name,
			"enabled":              app.Enabled,
			"domain":               app.Domain,
			"application_port":     app.ApplicationPort,
			"application_protocol": app.ApplicationProtocol,
			"app_types":            appTypes,
		}
		if app.ApplicationProtocol == "RDP" {
			appConfigMap["connection_security"] = app.ConnectionSecurity
		}
		appsConfig[i] = appConfigMap
	}
	commonAppsDto[0] = map[string]interface{}{
		"apps_config": appsConfig,
	}
	return commonAppsDto
}

func checkForPRAPortsOverlap(client *Client, app applicationsegmentpra.AppSegmentPRA) error {
	//lintignore:R018
	time.Sleep(time.Second * time.Duration(rand.Intn(5)))

	microTenantID := app.MicroTenantID
	service := client.ApplicationSegmentPRA
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	apps, _, err := applicationsegmentpra.GetAll(service)
	if err != nil {
		return err
	}
	for _, app2 := range apps {
		if found, common := sliceHasCommon(app.DomainNames, app2.DomainNames); found && app2.ID != app.ID && app2.Name != app.Name {
			// check for TCP ports
			if overlap, o1, o2 := PRAPortOverlap(app.TCPPortRanges, app2.TCPPortRanges); overlap {
				return fmt.Errorf("found TCP overlapping ports: %v of application %s with %v of application %s (%s) with common domain name %s", o1, app.Name, o2, app2.Name, app2.ID, common)
			}
			// check for UDP ports
			if overlap, o1, o2 := PRAPortOverlap(app.UDPPortRanges, app2.UDPPortRanges); overlap {
				return fmt.Errorf("found UDP overlapping ports: %v of application %s with %v of application %s (%s) with common domain name %s", o1, app.Name, o2, app2.Name, app2.ID, common)
			}
		}
	}
	return nil
}

func PRAPortOverlap(s1, s2 []string) (bool, []string, []string) {
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

func customizeDiffApplicationSegmentPRA(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	commonAppsDto := d.Get("common_apps_dto").(*schema.Set).List()

	for _, dto := range commonAppsDto {
		dtoMap := dto.(map[string]interface{})
		appsConfig := dtoMap["apps_config"].(*schema.Set).List()

		for _, appConfig := range appsConfig {
			appConfigMap := appConfig.(map[string]interface{})
			appProtocol := appConfigMap["application_protocol"].(string)
			connSecurity, connSecurityExists := appConfigMap["connection_security"]

			if appProtocol == "RDP" {
				if !connSecurityExists || connSecurity.(string) == "" {
					return errors.New("connection_security is required when application_protocol is RDP")
				}
			} else {
				if connSecurityExists && connSecurity.(string) != "" {
					return errors.New("connection_security can only be set when application_protocol is RDP")
				}
			}
		}
	}
	return nil
}
