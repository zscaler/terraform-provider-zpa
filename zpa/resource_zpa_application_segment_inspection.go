package zpa

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/applicationsegmentinspection"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/segmentgroup"
)

func resourceApplicationSegmentInspection() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationSegmentInspectionCreate,
		Read:   resourceApplicationSegmentInspectionRead,
		Update: resourceApplicationSegmentInspectionUpdate,
		Delete: resourceApplicationSegmentInspectionDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					d.Set("id", id)
				} else {
					resp, _, err := zClient.applicationsegmentinspection.GetByName(id)
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
			"tcp_keep_alive": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"0", "1",
				}, false),
			},
			"common_apps_dto": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"apps_config": {
							Type:     schema.TypeSet,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_options": {
										Type:     schema.TypeBool,
										Computed: true,
										Optional: true,
									},
									"app_id": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"app_types": {
										Type:     schema.TypeSet,
										Computed: true,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"INSPECT",
											}, false),
										},
									},
									"application_port": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"application_protocol": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"HTTP",
											"HTTPS",
										}, false),
									},
									"certificate_id": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"certificate_name": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"cname": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"domain": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
										Optional: true,
									},
									"hidden": {
										Type:     schema.TypeBool,
										Computed: true,
										Optional: true,
									},
									"local_domain": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"portal": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"trust_untrusted_cert": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
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

func resourceApplicationSegmentInspectionCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandInspectionApplicationSegment(d, zClient, "")
	if err := checkForInspectionPortsOverlap(zClient, req); err != nil {
		return err
	}
	log.Printf("[INFO] Creating application segment request\n%+v\n", req)
	if req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provde a valid segment group for the application segment")
		return fmt.Errorf("please provde a valid segment group for the application segment")
	}

	resp, _, err := zClient.applicationsegmentinspection.Create(req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Created inspection application segment request. ID: %v\n", resp.ID)
	d.SetId(resp.ID)

	return resourceApplicationSegmentInspectionRead(d, m)
}

func resourceApplicationSegmentInspectionRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.applicationsegmentinspection.Get(d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing inspection application segment %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	log.Printf("[INFO] Getting inspection application segment:\n%+v\n", resp)
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
	_ = d.Set("icmp_access_type", resp.ICMPAccessType)
	_ = d.Set("ip_anchored", resp.IPAnchored)
	_ = d.Set("health_reporting", resp.HealthReporting)
	_ = d.Set("select_connector_close_to_app", resp.SelectConnectorCloseToApp)
	_ = d.Set("use_in_dr_mode", resp.UseInDrMode)
	_ = d.Set("is_incomplete_dr_config", resp.IsIncompleteDRConfig)
	_ = d.Set("is_cname_enabled", resp.IsCnameEnabled)
	_ = d.Set("tcp_keep_alive", resp.TCPKeepAlive)
	_ = d.Set("tcp_port_ranges", resp.TCPPortRanges)
	_ = d.Set("udp_port_ranges", resp.UDPPortRanges)
	_ = d.Set("server_groups", flattenInspectionAppServerGroupsSimple(resp))

	if err := d.Set("common_apps_dto", flattenInspectionCommonAppsDto(resp.InspectionAppDto)); err != nil {
		return fmt.Errorf("failed to read common application in application segment %s", err)
	}

	if err := d.Set("tcp_port_range", flattenNetworkPorts(resp.TCPAppPortRange)); err != nil {
		return err
	}

	if err := d.Set("tcp_port_range", flattenNetworkPorts(resp.UDPAppPortRange)); err != nil {
		return err
	}

	return nil

}

func flattenInspectionAppServerGroupsSimple(serverGroup *applicationsegmentinspection.AppSegmentInspection) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(serverGroup.AppServerGroups))
	for i, group := range serverGroup.AppServerGroups {
		ids[i] = group.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func resourceApplicationSegmentInspectionUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating inspection application segment ID: %v\n", id)
	req := expandInspectionApplicationSegment(d, zClient, id)

	if d.HasChange("segment_group_id") && req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provde a valid segment group for the inspection application segment")
		return fmt.Errorf("please provde a valid segment group for the inspection application segment")
	}

	if _, _, err := zClient.applicationsegmentinspection.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := zClient.applicationsegmentinspection.Update(id, &req); err != nil {
		return err
	}

	return resourceApplicationSegmentInspectionRead(d, m)
}

func resourceApplicationSegmentInspectionDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	id := d.Id()
	segmentGroupID, ok := d.GetOk("segment_group_id")
	if ok && segmentGroupID != nil {
		gID, ok := segmentGroupID.(string)
		if ok && gID != "" {
			// detach it from segment group first
			if err := detachInspectionPortalsFromGroup(zClient, id, gID); err != nil {
				return err
			}
		}
	}
	log.Printf("[INFO] Deleting inspection application segment with id %v\n", id)
	if _, err := zClient.applicationsegmentinspection.Delete(id); err != nil {
		return err
	}

	return nil
}

func detachInspectionPortalsFromGroup(client *Client, segmentID, segmentGroupID string) error {
	log.Printf("[INFO] Detaching inspection application segment  %s from segment group: %s\n", segmentID, segmentGroupID)
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

func expandInspectionApplicationSegment(d *schema.ResourceData, zClient *Client, id string) applicationsegmentinspection.AppSegmentInspection {
	details := applicationsegmentinspection.AppSegmentInspection{
		ID:                        d.Id(),
		SegmentGroupID:            d.Get("segment_group_id").(string),
		BypassType:                d.Get("bypass_type").(string),
		ConfigSpace:               d.Get("config_space").(string),
		ICMPAccessType:            d.Get("icmp_access_type").(string),
		Description:               d.Get("description").(string),
		HealthReporting:           d.Get("health_reporting").(string),
		HealthCheckType:           d.Get("health_check_type").(string),
		TCPKeepAlive:              d.Get("tcp_keep_alive").(string),
		DoubleEncrypt:             d.Get("double_encrypt").(bool),
		Enabled:                   d.Get("enabled").(bool),
		PassiveHealthEnabled:      d.Get("passive_health_enabled").(bool),
		IPAnchored:                d.Get("ip_anchored").(bool),
		IsCnameEnabled:            d.Get("is_cname_enabled").(bool),
		SelectConnectorCloseToApp: d.Get("select_connector_close_to_app").(bool),
		UseInDrMode:               d.Get("use_in_dr_mode").(bool),
		IsIncompleteDRConfig:      d.Get("is_incomplete_dr_config").(bool),
		DomainNames:               expandStringInSlice(d, "domain_names"),
		TCPPortRanges:             expandList(d.Get("tcp_port_ranges").([]interface{})),
		UDPPortRanges:             expandList(d.Get("udp_port_ranges").([]interface{})),
		AppServerGroups:           expandInspectionAppServerGroups(d),
		CommonAppsDto:             expandInspectionCommonAppsDto(d),
	}
	if d.HasChange("name") {
		details.Name = d.Get("name").(string)
	}
	if d.HasChange("segment_group_name") {
		details.SegmentGroupName = d.Get("segment_group_name").(string)
	}
	if d.HasChange("server_groups") {
		details.AppServerGroups = expandInspectionAppServerGroups(d)
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

func expandInspectionCommonAppsDto(d *schema.ResourceData) applicationsegmentinspection.CommonAppsDto {
	result := applicationsegmentinspection.CommonAppsDto{}
	appsConfigInterface, ok := d.GetOk("common_apps_dto")
	if !ok {
		return result
	}
	appsConfigList, ok := appsConfigInterface.([]interface{})
	if !ok {
		return result
	}
	for _, appconf := range appsConfigList {
		appConfMap, ok := appconf.(map[string]interface{})
		if !ok {
			return result
		}
		result.AppsConfig = expandInspectionAppsConfig(appConfMap["apps_config"])
	}
	return result
}

func expandInspectionAppsConfig(appsConfigInterface interface{}) []applicationsegmentinspection.AppsConfig {
	appsConfig, ok := appsConfigInterface.(*schema.Set)
	if !ok {
		return []applicationsegmentinspection.AppsConfig{}
	}
	log.Printf("[INFO] apps config data: %+v\n", appsConfig)
	var commonAppConfigDto []applicationsegmentinspection.AppsConfig
	for _, commonAppConfig := range appsConfig.List() {
		commonAppConfig, ok := commonAppConfig.(map[string]interface{})
		if ok {
			appTypesSet, ok := commonAppConfig["app_types"].(*schema.Set)
			if !ok {
				continue
			}
			appTypes := SetToStringSlice(appTypesSet)
			commonAppConfigDto = append(commonAppConfigDto, applicationsegmentinspection.AppsConfig{
				AllowOptions:        commonAppConfig["allow_options"].(bool),
				AppID:               commonAppConfig["app_id"].(string),
				ID:                  commonAppConfig["id"].(string),
				ApplicationPort:     commonAppConfig["application_port"].(string),
				ApplicationProtocol: commonAppConfig["application_protocol"].(string),
				CertificateID:       commonAppConfig["certificate_id"].(string),
				CertificateName:     commonAppConfig["certificate_name"].(string),
				Cname:               commonAppConfig["cname"].(string),
				Description:         commonAppConfig["description"].(string),
				Domain:              commonAppConfig["domain"].(string),
				Enabled:             commonAppConfig["enabled"].(bool),
				Hidden:              commonAppConfig["hidden"].(bool),
				LocalDomain:         commonAppConfig["local_domain"].(string),
				Name:                commonAppConfig["name"].(string),
				Portal:              commonAppConfig["portal"].(bool),
				AppTypes:            appTypes,
			})
		}
	}
	return commonAppConfigDto
}

func expandInspectionAppServerGroups(d *schema.ResourceData) []applicationsegmentinspection.AppServerGroups {
	serverGroupsInterface, ok := d.GetOk("server_groups")
	if ok {
		serverGroup := serverGroupsInterface.(*schema.Set)
		log.Printf("[INFO] app server groups data: %+v\n", serverGroup)
		var serverGroups []applicationsegmentinspection.AppServerGroups
		for _, appServerGroup := range serverGroup.List() {
			appServerGroup, _ := appServerGroup.(map[string]interface{})
			if appServerGroup != nil {
				for _, id := range appServerGroup["id"].([]interface{}) {
					serverGroups = append(serverGroups, applicationsegmentinspection.AppServerGroups{
						ID: id.(string),
					})
				}
			}
		}
		return serverGroups
	}

	return []applicationsegmentinspection.AppServerGroups{}
}

func flattenInspectionCommonAppsDto(apps []applicationsegmentinspection.InspectionAppDto) []interface{} {
	commonApp := make([]interface{}, 1)
	commonApp[0] = map[string]interface{}{
		"apps_config": flattenInspectionAppsConfig(apps),
	}
	return commonApp
}

func flattenInspectionAppsConfig(appConfigs []applicationsegmentinspection.InspectionAppDto) []interface{} {
	appConfig := make([]interface{}, len(appConfigs))
	for i, val := range appConfigs {
		appConfig[i] = map[string]interface{}{
			"name":                 val.Name,
			"id":                   val.ID,
			"app_id":               val.AppID,
			"application_port":     val.ApplicationPort,
			"application_protocol": val.ApplicationProtocol,
			"certificate_id":       val.CertificateID,
			"certificate_name":     val.CertificateName,
			"description":          val.Description,
			"domain":               val.Domain,
			"enabled":              val.Enabled,
			"hidden":               val.Hidden,
			"portal":               val.Portal,
		}
	}
	return appConfig
}

func checkForInspectionPortsOverlap(client *Client, app applicationsegmentinspection.AppSegmentInspection) error {
	time.Sleep(time.Second * time.Duration(rand.Intn(5)))
	apps, _, err := client.browseraccess.GetAll()
	if err != nil {
		return err
	}
	for _, app2 := range apps {
		if found, common := sliceHasCommon(app.DomainNames, app2.DomainNames); found && app2.ID != app.ID && app2.Name != app.Name {
			// check for udp ports
			if overlap, o1, o2 := InspectionPortOverlap(app.TCPPortRanges, app2.TCPPortRanges); overlap {
				return fmt.Errorf("found TCP overlapping ports: %v of application %s with %v of application %s (%s) with common domain name %s", o1, app.Name, o2, app2.Name, app2.ID, common)
			}
			if overlap, o1, o2 := InspectionPortOverlap(app.UDPPortRanges, app2.UDPPortRanges); overlap {
				return fmt.Errorf("found UDP overlapping ports: %v of application %s with %v of application %s (%s) with common domain name %s", o1, app.Name, o2, app2.Name, app2.ID, common)
			}
		}
	}
	return nil

}

func InspectionPortOverlap(s1, s2 []string) (bool, []string, []string) {
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
