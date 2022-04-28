package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/appsegment_inspection"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/client"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/segmentgroup"
)

func resourceAppSegmentInspection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppSegmentInspectionCreate,
		Read:   resourceAppSegmentInspectionRead,
		Update: resourceAppSegmentInspectionUpdate,
		Delete: resourceAppSegmentInspectionDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					d.Set("id", id)
				} else {
					resp, _, err := zClient.appsegment_inspection.GetByName(id)
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
			"common_apps_dto": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"apps_config": {
							Type:     schema.TypeList,
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
									"app_types": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"INSPECT",
										}, false),
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
									"ba_app_id": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
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
									"inspect_app_id": {
										Type:     schema.TypeString,
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
									},
									"path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"portal": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"trust_untrusted_cert": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"inspection_apps": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							}, false),
						},
						"certificate_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"certificate_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
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

func resourceAppSegmentInspectionCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandAppSegmentInspection(d)
	log.Printf("[INFO] Creating application segment inspection request\n%+v\n", req)

	if req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provde a valid segment group for the application segment")
		return fmt.Errorf("please provde a valid segment group for the application segment")
	}

	appsegment_inspection, _, err := zClient.appsegment_inspection.Create(req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Created application segment inspection request. ID: %v\n", appsegment_inspection.ID)
	d.SetId(appsegment_inspection.ID)

	return resourceAppSegmentInspectionRead(d, m)
}

func resourceAppSegmentInspectionRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.appsegment_inspection.Get(d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing application segment inspection %s from state because it no longer exists in ZPA", d.Id())
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

	if err := d.Set("common_apps_dto", flattenCommonAppsDto(resp.CommonAppsDto)); err != nil {
		return fmt.Errorf("failed to read common apps %s", err)
	}

	if err := d.Set("inspection_apps", flattenInspectionApps(resp)); err != nil {
		return fmt.Errorf("failed to read inspection application segment %s", err)
	}

	if err := d.Set("server_groups", flattenInspectionAppServerGroups(resp.AppServerGroups)); err != nil {
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

func resourceAppSegmentInspectionUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating application segment inspection ID: %v\n", id)
	req := expandAppSegmentInspection(d)

	if d.HasChange("segment_group_id") && req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provde a valid segment group for the browser access application segment")
		return fmt.Errorf("please provde a valid segment group for the browser access application segment")
	}

	if _, err := zClient.appsegment_inspection.Update(id, &req); err != nil {
		return err
	}

	return resourceAppSegmentInspectionRead(d, m)
}

func resourceAppSegmentInspectionDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	id := d.Id()
	segmentGroupID, ok := d.GetOk("segment_group_id")
	if ok && segmentGroupID != nil {
		gID, ok := segmentGroupID.(string)
		if ok && gID != "" {
			// detach it from segment group first
			if err := detachAppSegmentInspectionFromGroup(zClient, id, gID); err != nil {
				return err
			}
		}
	}
	log.Printf("[INFO] Deleting application segment inspection with id %v\n", id)
	if _, err := zClient.appsegment_inspection.Delete(id); err != nil {
		return err
	}

	return nil
}

func detachAppSegmentInspectionFromGroup(client *Client, segmentID, segmentGroupID string) error {
	log.Printf("[INFO] Detaching application segment inspection  %s from segment group: %s\n", segmentID, segmentGroupID)
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

func expandAppSegmentInspection(d *schema.ResourceData) appsegment_inspection.AppSegmentInspection {
	details := appsegment_inspection.AppSegmentInspection{
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
		details.AppServerGroups = expandAppSegmentInspectionServerGroups(d)
	}
	if d.HasChange("inspection_apps") {
		details.InspectionApps = expandInspectionApps(d)
	}
	if d.HasChange("common_apps_dto") {
		details.CommonAppsDto = expandCommonAppsDto(d)
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

func expandCommonAppsDto(d *schema.ResourceData) appsegment_inspection.CommonAppsDto {
	appSegmentInspection := appsegment_inspection.CommonAppsDto{
		AppConfig: expandAppsConfig(d),
	}
	return appSegmentInspection
}

func expandAppsConfig(d *schema.ResourceData) []appsegment_inspection.AppConfig {
	appConfigInterface, ok := d.GetOk("apps_config")
	if ok {
		appConfig := appConfigInterface.([]interface{})
		log.Printf("[INFO] application segment inspection config data: %+v\n", appConfig)
		var appConfigs []appsegment_inspection.AppConfig
		for _, inspectionAppConfig := range appConfig {
			inspectionAppConfig, ok := inspectionAppConfig.(map[string]interface{})
			if ok {
				appConfigs = append(appConfigs, appsegment_inspection.AppConfig{
					AllowOptions:        inspectionAppConfig["allow_options"].(bool),
					AppID:               inspectionAppConfig["app_id"].(string),
					AppTypes:            inspectionAppConfig["app_types"].(string),
					ApplicationPort:     inspectionAppConfig["application_port"].(string),
					ApplicationProtocol: inspectionAppConfig["application_protocol"].(string),
					CertificateID:       inspectionAppConfig["certificate_id"].(string),
					CertificateName:     inspectionAppConfig["certificate_name"].(string),
					Cname:               inspectionAppConfig["cname"].(string),
					Description:         inspectionAppConfig["description"].(string),
					Domain:              inspectionAppConfig["domain"].(string),
					Enabled:             inspectionAppConfig["enabled"].(bool),
					Hidden:              inspectionAppConfig["hidden"].(bool),
					InspectAppId:        inspectionAppConfig["inspect_app_id"].(string),
					LocalDomain:         inspectionAppConfig["local_domain"].(string),
					Path:                inspectionAppConfig["path"].(string),
					TrustUntrustedCert:  inspectionAppConfig["trust_untrusted_cert"].(bool),
					Name:                inspectionAppConfig["name"].(string),
				})
			}
		}
		return appConfigs
	}

	return []appsegment_inspection.AppConfig{}
}

func expandInspectionApps(d *schema.ResourceData) []appsegment_inspection.InspectionApps {
	inspectionAppInterface, ok := d.GetOk("inspection_apps")
	if ok {
		inspection := inspectionAppInterface.([]interface{})
		log.Printf("[INFO] inspection apps data: %+v\n", inspection)
		var inspectionApps []appsegment_inspection.InspectionApps
		for _, inspectionApp := range inspection {
			inspectionApp, ok := inspectionApp.(map[string]interface{})
			if ok {
				inspectionApps = append(inspectionApps, appsegment_inspection.InspectionApps{
					AppID:               inspectionApp["app_id"].(string),
					ApplicationPort:     inspectionApp["application_port"].(string),
					ApplicationProtocol: inspectionApp["application_protocol"].(string),
					CertificateID:       inspectionApp["certificate_id"].(string),
					CertificateName:     inspectionApp["certificate_name"].(string),
					Description:         inspectionApp["description"].(string),
					Domain:              inspectionApp["domain"].(string),
					Enabled:             inspectionApp["enabled"].(bool),
					Name:                inspectionApp["name"].(string),
				})
			}
		}
		return inspectionApps
	}

	return []appsegment_inspection.InspectionApps{}
}

func expandAppSegmentInspectionServerGroups(d *schema.ResourceData) []appsegment_inspection.AppServerGroups {
	serverGroupsInterface, ok := d.GetOk("server_groups")
	if ok {
		serverGroup := serverGroupsInterface.(*schema.Set)
		log.Printf("[INFO] app server groups data: %+v\n", serverGroup)
		var serverGroups []appsegment_inspection.AppServerGroups
		for _, appServerGroup := range serverGroup.List() {
			appServerGroup, _ := appServerGroup.(map[string]interface{})
			if appServerGroup != nil {
				for _, id := range appServerGroup["id"].([]interface{}) {
					serverGroups = append(serverGroups, appsegment_inspection.AppServerGroups{
						ID: id.(string),
					})
				}
			}
		}
		return serverGroups
	}

	return []appsegment_inspection.AppServerGroups{}
}
