package zpa

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentinspection"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceApplicationSegmentInspection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationSegmentInspectionCreate,
		ReadContext:   resourceApplicationSegmentInspectionRead,
		UpdateContext: resourceApplicationSegmentInspectionUpdate,
		DeleteContext: resourceApplicationSegmentInspectionDelete,
		CustomizeDiff: customizeDiffApplicationSegmentInspection,
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
					_ = d.Set("id", id)
				} else {
					resp, _, err := applicationsegmentinspection.GetByName(ctx, service, id)
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
			"adp_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates if Active Directory Inspection is enabled or not for the application. This allows the application segment's traffic to be inspected by Active Directory (AD) Protection.",
			},
			"auto_app_protect_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If autoAppProtectEnabled is set to true, this field indicates if the application segment’s traffic is inspected by AppProtection.",
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
			"tcp_protocols": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "TCP port ranges used to access the app.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"udp_protocols": {
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"apps_config": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"app_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"inspect_app_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"app_types": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"application_port": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"application_protocol": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
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
									"domain": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"trust_untrusted_cert": {
										Type:     schema.TypeBool,
										Computed: true,
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

func resourceApplicationSegmentInspectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandInspectionApplicationSegment(ctx, d, zClient, "")

	if err := validateAppPorts(req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating application segment request\n%+v\n", req)
	if req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provide a valid segment group for the application segment")
		return diag.FromErr(fmt.Errorf("please provide a valid segment group for the application segment"))
	}

	resp, _, err := applicationsegmentinspection.Create(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created inspection application segment request. ID: %v\n", resp.ID)
	d.SetId(resp.ID)

	return resourceApplicationSegmentInspectionRead(ctx, d, meta)
}

func resourceApplicationSegmentInspectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, _, err := applicationsegmentinspection.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing inspection application segment %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting sra application segment:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("adp_enabled", resp.AdpEnabled)
	_ = d.Set("auto_app_protect_enabled", resp.AutoAppProtectEnabled)
	_ = d.Set("segment_group_id", resp.SegmentGroupID)
	_ = d.Set("bypass_type", resp.BypassType)
	_ = d.Set("bypass_on_reauth", resp.BypassOnReauth)
	_ = d.Set("config_space", resp.ConfigSpace)
	_ = d.Set("domain_names", resp.DomainNames)
	_ = d.Set("description", resp.Description)
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
	_ = d.Set("tcp_protocols", resp.TCPProtocols)
	_ = d.Set("udp_protocols", resp.UDPProtocols)
	_ = d.Set("server_groups", flattenCommonAppServerGroups(resp.AppServerGroups))

	// Map inspect to common_apps_dto.apps_config for state management
	if err := mapInspectAppsToCommonApps(d, resp.InspectionAppDto); err != nil {
		return diag.FromErr(err)
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

func resourceApplicationSegmentInspectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating inspection application segment ID: %v\n", id)

	// Retrieve the current resource to get app_id and pra_app_id
	resp, _, err := applicationsegmentinspection.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error retrieving application segment: %v", err))
	}

	// Extract app_id and inspect_app_id from praApps and set in common_apps_dto in state
	if err := setInspectionAppIDsInCommonAppsDto(d, resp.InspectionAppDto); err != nil {
		return diag.FromErr(fmt.Errorf("error setting app_id and inspect_app_id in common_apps_dto: %v", err))
	}

	req := expandInspectionApplicationSegment(ctx, d, zClient, "")

	if err := validateAppPorts(req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return diag.FromErr(err)
	}

	_, err = applicationsegmentinspection.Update(ctx, service, id, &req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating application segment inspection: %v", err))
	}

	return resourceApplicationSegmentInspectionRead(ctx, d, meta)
}

func resourceApplicationSegmentInspectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id := d.Id()
	segmentGroupID, ok := d.GetOk("segment_group_id")
	if ok && segmentGroupID != nil {
		gID, ok := segmentGroupID.(string)
		if ok && gID != "" {
			// detach it from segment group first
			if err := detachSegmentGroup(ctx, zClient, id, gID); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	log.Printf("[INFO] Deleting inspection application segment with id %v\n", id)
	if _, err := applicationsegmentinspection.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandInspectionApplicationSegment(ctx context.Context, d *schema.ResourceData, zClient *Client, id string) applicationsegmentinspection.AppSegmentInspection {
	microTenantID := GetString(d.Get("microtenant_id"))
	service := zClient.Service
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	details := applicationsegmentinspection.AppSegmentInspection{
		ID:                        d.Id(),
		Name:                      d.Get("name").(string),
		SegmentGroupID:            d.Get("segment_group_id").(string),
		BypassType:                d.Get("bypass_type").(string),
		BypassOnReauth:            d.Get("bypass_on_reauth").(bool),
		ConfigSpace:               d.Get("config_space").(string),
		ICMPAccessType:            d.Get("icmp_access_type").(string),
		Description:               d.Get("description").(string),
		HealthReporting:           d.Get("health_reporting").(string),
		HealthCheckType:           d.Get("health_check_type").(string),
		AdpEnabled:                d.Get("adp_enabled").(bool),
		AutoAppProtectEnabled:     d.Get("auto_app_protect_enabled").(bool),
		PassiveHealthEnabled:      d.Get("passive_health_enabled").(bool),
		Enabled:                   d.Get("enabled").(bool),
		DoubleEncrypt:             d.Get("double_encrypt").(bool),
		IPAnchored:                d.Get("ip_anchored").(bool),
		IsCnameEnabled:            d.Get("is_cname_enabled").(bool),
		SelectConnectorCloseToApp: d.Get("select_connector_close_to_app").(bool),
		UseInDrMode:               d.Get("use_in_dr_mode").(bool),
		TCPKeepAlive:              d.Get("tcp_keep_alive").(string),
		IsIncompleteDRConfig:      d.Get("is_incomplete_dr_config").(bool),
		DomainNames:               expandStringInSlice(d, "domain_names"),
		TCPProtocols:              expandStringInSlice(d, "tcp_protocols"),
		UDPProtocols:              expandStringInSlice(d, "udp_protocols"),
		AppServerGroups:           expandCommonServerGroups(d),
		CommonAppsDto:             expandInspectionCommonAppsDto(d),

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

	if d.HasChange("name") {
		details.Name = d.Get("name").(string)
	}
	if d.HasChange("server_groups") {
		details.AppServerGroups = expandCommonServerGroups(d)
	}

	return details
}

func expandInspectionCommonAppsDto(d *schema.ResourceData) applicationsegmentinspection.CommonAppsDto {
	result := applicationsegmentinspection.CommonAppsDto{}
	if commonAppsInterface, ok := d.GetOk("common_apps_dto"); ok {
		commonAppsList := commonAppsInterface.(*schema.Set).List()
		if len(commonAppsList) > 0 {
			commonAppMap := commonAppsList[0].(map[string]interface{})
			result.AppsConfig = expandInspectionAppsConfig(commonAppMap["apps_config"])
		}
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
		appConfigMap, ok := commonAppConfig.(map[string]interface{})
		if ok {
			// Automatically set `name` to match `domain` to prevent drift
			appConfigMap["name"] = appConfigMap["domain"].(string)

			appTypesSet, ok := appConfigMap["app_types"].(*schema.Set)
			if !ok {
				continue
			}
			appTypes := SetToStringSlice(appTypesSet)

			// Retrieve protocols as a slice of strings
			// protocolsSet, ok := appConfigMap["protocols"].(*schema.Set)
			// var protocols []string
			// if ok {
			// 	protocols = SetToStringSlice(protocolsSet)
			// }

			commonAppConfigDto = append(commonAppConfigDto, applicationsegmentinspection.AppsConfig{
				AppID:               appConfigMap["app_id"].(string),
				InspectAppID:        appConfigMap["inspect_app_id"].(string),
				Name:                appConfigMap["name"].(string),
				Description:         appConfigMap["description"].(string),
				Enabled:             appConfigMap["enabled"].(bool),
				ApplicationPort:     appConfigMap["application_port"].(string),
				ApplicationProtocol: appConfigMap["application_protocol"].(string),
				CertificateID:       appConfigMap["certificate_id"].(string),
				Domain:              appConfigMap["domain"].(string),
				TrustUntrustedCert:  appConfigMap["trust_untrusted_cert"].(bool),
				// Protocols:           protocols, // Set protocols here
				AppTypes: appTypes,
			})
		}
	}
	return commonAppConfigDto
}

func mapInspectAppsToCommonApps(d *schema.ResourceData, inspectionApps []applicationsegmentinspection.InspectionAppDto) error {
	// If the API returned any Inspection Apps, map them to common_apps_dto.apps_config
	if len(inspectionApps) == 0 {
		return nil
	}

	// Create a single common_apps_dto with multiple apps_config blocks
	commonAppsConfig := make([]interface{}, len(inspectionApps))
	for i, app := range inspectionApps {
		commonAppMap := map[string]interface{}{
			"app_id":               app.AppID, // Populate app_id from InspectionAppDto
			"app_types":            []interface{}{"INSPECT"},
			"application_protocol": app.ApplicationProtocol,
			"application_port":     app.ApplicationPort,
			"certificate_id":       app.CertificateID,
			"description":          app.Description,
			"domain":               app.Domain,
			"enabled":              app.Enabled,
			"name":                 app.Name,
			// "protocols":            app.Protocols,
			"trust_untrusted_cert": app.TrustUntrustedCert,
		}
		// Only set inspect_app_id if it's present in the response
		if app.ID != "" {
			commonAppMap["inspect_app_id"] = app.ID // Populate inspect_app_id from InspectAppID
		}
		commonAppsConfig[i] = commonAppMap
	}

	// Wrap commonAppsConfig in the common_apps_dto block
	commonAppsDto := []interface{}{
		map[string]interface{}{
			"apps_config": commonAppsConfig,
		},
	}

	// Set common_apps_dto in the resource data
	if err := d.Set("common_apps_dto", commonAppsDto); err != nil {
		return fmt.Errorf("failed to set common_apps_dto: %s", err)
	}
	return nil
}

func setInspectionAppIDsInCommonAppsDto(d *schema.ResourceData, inspectionApps []applicationsegmentinspection.InspectionAppDto) error {
	if len(inspectionApps) == 0 {
		return nil
	}

	// Extract app_id and inspect_app_id from the first Inspect app in the list
	appID := inspectionApps[0].AppID
	inspectAppID := inspectionApps[0].ID

	// Update the common_apps_dto with extracted app_id and inspect_app_id values
	commonAppsDto := d.Get("common_apps_dto").(*schema.Set).List()
	if len(commonAppsDto) == 0 {
		return fmt.Errorf("common_apps_dto block is missing")
	}

	// Update the first entry in commonAppsDto.appsConfig with app_id and inspect_app_id
	commonAppConfig := commonAppsDto[0].(map[string]interface{})
	appsConfig := commonAppConfig["apps_config"].(*schema.Set).List()

	if len(appsConfig) > 0 {
		appConfig := appsConfig[0].(map[string]interface{})
		appConfig["app_id"] = appID
		appConfig["inspect_app_id"] = inspectAppID
	}

	// Write the updated config back to the resource data
	if err := d.Set("common_apps_dto", commonAppsDto); err != nil {
		return fmt.Errorf("failed to set common_apps_dto: %v", err)
	}

	return nil
}

func customizeDiffApplicationSegmentInspection(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// Validation for adp_enabled and auto_app_protect_enabled attributes
	adpEnabled := d.Get("adp_enabled").(bool)
	autoAppProtectEnabled := d.Get("auto_app_protect_enabled").(bool)
	if adpEnabled && autoAppProtectEnabled {
		return fmt.Errorf("if 'adp_enabled' is set to true, 'auto_app_protect_enabled' cannot be true")
	}

	// Validation for common_apps_dto.apps_config fields
	commonAppsDto, ok := d.GetOk("common_apps_dto")
	if !ok || len(commonAppsDto.(*schema.Set).List()) == 0 {
		return nil // If there's no common_apps_dto, skip further validation
	}

	appsConfig := commonAppsDto.(*schema.Set).List()[0].(map[string]interface{})["apps_config"].(*schema.Set).List()
	for _, config := range appsConfig {
		appConfig := config.(map[string]interface{})
		protocol := appConfig["application_protocol"].(string)
		certID, hasCertID := appConfig["certificate_id"]

		// Check if protocol is HTTP and certificate ID is set
		if protocol == "HTTP" && hasCertID && certID.(string) != "" {
			return fmt.Errorf("certificate ID should not be set when 'application_protocol' is HTTP")
		}
	}

	return nil
}
