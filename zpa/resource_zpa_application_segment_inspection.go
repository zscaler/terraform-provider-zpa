package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegmentinspection"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/segmentgroup"
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
			"tcp_keep_alive": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"0", "1",
				}, false),
			},
			"common_apps_dto": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"apps_config": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"app_types": {
										Type:     schema.TypeSet,
										Optional: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"application_port": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"application_protocol": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										ValidateFunc: validation.StringInSlice([]string{
											"HTTP",
											"HTTPS",
										}, false),
									},
									"certificate_id": {
										Type:     schema.TypeString,
										ForceNew: true,
										Optional: true,
									},
									"domain": {
										Type:     schema.TypeString,
										ForceNew: true,
										Optional: true,
									},
									"trust_untrusted_cert": {
										Type:     schema.TypeBool,
										ForceNew: true,
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
	}
}

func resourceApplicationSegmentInspectionCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandInspectionApplicationSegment(d, zClient, "")

	if err := validateAppPorts(req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return err
	}
	if err := validateProtocolAndCertID(d); err != nil {
		return err
	}
	log.Printf("[INFO] Creating application segment request\n%+v\n", req)
	if req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provide a valid segment group for the application segment")
		return fmt.Errorf("please provide a valid segment group for the application segment")
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
	_ = d.Set("icmp_access_type", resp.ICMPAccessType)
	_ = d.Set("select_connector_close_to_app", resp.SelectConnectorCloseToApp)
	_ = d.Set("use_in_dr_mode", resp.UseInDrMode)
	_ = d.Set("is_incomplete_dr_config", resp.IsIncompleteDRConfig)
	_ = d.Set("tcp_keep_alive", resp.TCPKeepAlive)
	_ = d.Set("ip_anchored", resp.IPAnchored)
	_ = d.Set("health_reporting", resp.HealthReporting)
	_ = d.Set("tcp_port_ranges", convertPortsToListString(resp.TCPAppPortRange))
	_ = d.Set("udp_port_ranges", convertPortsToListString(resp.UDPAppPortRange))
	_ = d.Set("server_groups", flattenInspectionAppServerGroupsSimple(resp.AppServerGroups))

	if err := d.Set("common_apps_dto", flattenInspectionCommonAppsDto(d, resp.InspectionAppDto)); err != nil {
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

func flattenInspectionAppServerGroupsSimple(serverGroup []applicationsegmentinspection.AppServerGroups) []interface{} {
	ids := make([]interface{}, len(serverGroup))
	for i, group := range serverGroup {
		ids[i] = group.ID
	}
	return ids
}

func resourceApplicationSegmentInspectionUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating inspection application segment ID: %v\n", id)
	req := expandInspectionApplicationSegment(d, zClient, id)

	if err := validateAppPorts(req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return err
	}

	if d.HasChange("segment_group_id") && req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provide a valid segment group for the inspection application segment")
		return fmt.Errorf("please provide a valid segment group for the inspection application segment")
	}
	if err := validateProtocolAndCertID(d); err != nil {
		return err
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
		Name:                      d.Get("name").(string),
		SegmentGroupID:            d.Get("segment_group_id").(string),
		SegmentGroupName:          d.Get("segment_group_name").(string),
		BypassType:                d.Get("bypass_type").(string),
		ConfigSpace:               d.Get("config_space").(string),
		ICMPAccessType:            d.Get("icmp_access_type").(string),
		Description:               d.Get("description").(string),
		HealthReporting:           d.Get("health_reporting").(string),
		HealthCheckType:           d.Get("health_check_type").(string),
		PassiveHealthEnabled:      d.Get("passive_health_enabled").(bool),
		DoubleEncrypt:             d.Get("double_encrypt").(bool),
		Enabled:                   d.Get("enabled").(bool),
		IPAnchored:                d.Get("ip_anchored").(bool),
		IsCnameEnabled:            d.Get("is_cname_enabled").(bool),
		SelectConnectorCloseToApp: d.Get("select_connector_close_to_app").(bool),
		UseInDrMode:               d.Get("use_in_dr_mode").(bool),
		TCPKeepAlive:              d.Get("tcp_keep_alive").(string),
		IsIncompleteDRConfig:      d.Get("is_incomplete_dr_config").(bool),
		DomainNames:               expandStringInSlice(d, "domain_names"),
		TCPAppPortRange:           []common.NetworkPorts{},
		UDPAppPortRange:           []common.NetworkPorts{},
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
	appsConfigSet, ok := appsConfigInterface.(*schema.Set)
	if !ok {
		return result
	}
	for _, appconf := range appsConfigSet.List() {
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
				Name:                commonAppConfig["name"].(string),
				Description:         commonAppConfig["description"].(string),
				Enabled:             commonAppConfig["enabled"].(bool),
				ApplicationPort:     commonAppConfig["application_port"].(string),
				ApplicationProtocol: commonAppConfig["application_protocol"].(string),
				CertificateID:       commonAppConfig["certificate_id"].(string),
				Domain:              commonAppConfig["domain"].(string),
				TrustUntrustedCert:  commonAppConfig["trust_untrusted_cert"].(bool),
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
			if ok {
				for _, id := range appServerGroup["id"].(*schema.Set).List() {
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

func flattenInspectionCommonAppsDto(d *schema.ResourceData, apps []applicationsegmentinspection.InspectionAppDto) []interface{} {
	commonApp := make([]interface{}, 1)
	commonApp[0] = map[string]interface{}{
		"apps_config": flattenInspectionAppsConfig(d, apps),
	}
	return commonApp
}

func flattenInspectionAppsConfig(d *schema.ResourceData, appConfigs []applicationsegmentinspection.InspectionAppDto) []interface{} {
	cApp := expandInspectionCommonAppsDto(d)

	appConfig := make([]interface{}, len(appConfigs))
	for i, val := range appConfigs {
		appTypes := []string{}
		for _, a := range cApp.AppsConfig {
			if a.Name == val.Name {
				appTypes = a.AppTypes
			}
		}
		appConfig[i] = map[string]interface{}{
			"name":                 val.Name,
			"description":          val.Description,
			"enabled":              val.Enabled,
			"domain":               val.Domain,
			"application_port":     val.ApplicationPort,
			"certificate_id":       val.CertificateID,
			"application_protocol": val.ApplicationProtocol,
			"trust_untrusted_cert": val.TrustUntrustedCert,
			"app_types":            appTypes,
		}
	}
	return appConfig
}

func validateProtocolAndCertID(d *schema.ResourceData) error {
	commonAppsDto, ok := d.GetOk("common_apps_dto")
	if !ok || len(commonAppsDto.(*schema.Set).List()) == 0 {
		return nil // or handle it as per your logic
	}

	appsConfig := commonAppsDto.(*schema.Set).List()[0].(map[string]interface{})["apps_config"].(*schema.Set).List()
	for _, config := range appsConfig {
		appConfig := config.(map[string]interface{})
		protocol := appConfig["application_protocol"].(string)
		certID := appConfig["certificate_id"].(string)

		if protocol == "HTTP" && certID != "" {
			return fmt.Errorf("certificate ID should not be set when application protocol is HTTP")
		}
	}
	return nil
}
