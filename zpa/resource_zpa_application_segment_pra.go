package zpa

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentpra"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

func resourceApplicationSegmentPRA() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationSegmentPRACreate,
		ReadContext:   resourceApplicationSegmentPRARead,
		UpdateContext: resourceApplicationSegmentPRAUpdate,
		DeleteContext: resourceApplicationSegmentPRADelete,
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
					_ = d.Set("id", id)
				} else {
					resp, _, err := applicationsegmentpra.GetByName(ctx, service, id)
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
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"segment_group_id": {
				Type:     schema.TypeString,
				Required: true,
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
									"pra_app_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
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
											"RDP", "SSH", "VNC",
										}, false),
									},
									"connection_security": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ValidateFunc: validation.StringInSlice([]string{
											"ANY", "NLA", "NLA_EXT", "TLS", "VM_CONNECT", "RDP",
										}, false),
									},
									"domain": {
										Type:     schema.TypeString,
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

func resourceApplicationSegmentPRACreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandSRAApplicationSegment(ctx, d, zClient, "")

	if err := validateAppPorts(req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating application segment request\n%+v\n", req)
	resp, _, err := applicationsegmentpra.Create(ctx, service, req)
	if err != nil {
		log.Printf("[ERROR] Failed to create application segment: %s", err)
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created application segment request. ID: %v\n", resp.ID)
	d.SetId(resp.ID)

	return resourceApplicationSegmentPRARead(ctx, d, meta)
}

func resourceApplicationSegmentPRARead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := applicationsegmentpra.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing sra application segment %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting sra application segment:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("segment_group_id", resp.SegmentGroupID)
	_ = d.Set("bypass_type", resp.BypassType)
	_ = d.Set("bypass_on_reauth", resp.BypassOnReauth)
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
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	_ = d.Set("select_connector_close_to_app", resp.SelectConnectorCloseToApp)
	_ = d.Set("use_in_dr_mode", resp.UseInDrMode)
	_ = d.Set("is_incomplete_dr_config", resp.IsIncompleteDRConfig)
	_ = d.Set("tcp_keep_alive", resp.TCPKeepAlive)
	_ = d.Set("ip_anchored", resp.IpAnchored)
	_ = d.Set("health_reporting", resp.HealthReporting)
	_ = d.Set("server_groups", flattenCommonAppServerGroups(resp.ServerGroups))

	// Map pra_apps to common_apps_dto.apps_config for state management
	if err := mapPRAAppsToCommonApps(d, resp.PRAApps); err != nil {
		return diag.FromErr(err)
	}

	// Map pra_apps back to common_apps_dto.apps_config
	// if err := mapPRAAppsToCommonApps(d, resp.PRAApps); err != nil {
	// 	return err
	// }

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

func resourceApplicationSegmentPRAUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating PRA application segment ID: %v\n", id)

	// Retrieve the current resource to get app_id and pra_app_id
	resp, _, err := applicationsegmentpra.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error retrieving application segment: %v", err))
	}

	// Extract app_id and pra_app_id from praApps and set in common_apps_dto in state
	if err := setAppIDsInCommonAppsDto(d, resp.PRAApps); err != nil {
		return diag.FromErr(fmt.Errorf("error setting app_id and pra_app_id in common_apps_dto: %v", err))
	}

	// Prepare the request payload for the update
	req := expandSRAApplicationSegment(ctx, d, zClient, "")

	if err := validateAppPorts(req.SelectConnectorCloseToApp, req.UDPAppPortRange, req.UDPPortRanges); err != nil {
		return diag.FromErr(err)
	}

	// Perform the update
	_, err = applicationsegmentpra.Update(ctx, service, id, &req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating application segment: %v", err))
	}

	// Refresh the state after the update to ensure correctness
	return resourceApplicationSegmentPRARead(ctx, d, meta)
}

func resourceApplicationSegmentPRADelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Use MicroTenant if available
	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// service := zClient.V3.ApplicationSegmentPRA.WithMicroTenant(GetString(d.Get("microtenant_id")))
	// policySetControllerService := zClient.V3.PolicySetController.WithMicroTenant(GetString(d.Get("microtenant_id")))

	log.Printf("[INFO] Deleting application segment pra with id %v\n", d.Id())
	detachAppConnectorGroupFromAllAccessPolicyRules(ctx, d, service)

	if _, err := applicationsegmentpra.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[INFO] Application segment pra deleted successfully")
	return nil
}

func expandSRAApplicationSegment(ctx context.Context, d *schema.ResourceData, zClient *Client, id string) applicationsegmentpra.AppSegmentPRA {
	microTenantID := GetString(d.Get("microtenant_id"))
	service := zClient.Service // Unified service interface
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	details := applicationsegmentpra.AppSegmentPRA{
		ID:                        d.Id(),
		SegmentGroupID:            d.Get("segment_group_id").(string),
		BypassType:                d.Get("bypass_type").(string),
		BypassOnReauth:            d.Get("bypass_on_reauth").(bool),
		ConfigSpace:               d.Get("config_space").(string),
		IcmpAccessType:            d.Get("icmp_access_type").(string),
		Description:               d.Get("description").(string),
		HealthReporting:           d.Get("health_reporting").(string),
		HealthCheckType:           d.Get("health_check_type").(string),
		PassiveHealthEnabled:      d.Get("passive_health_enabled").(bool),
		DoubleEncrypt:             d.Get("double_encrypt").(bool),
		Enabled:                   d.Get("enabled").(bool),
		IpAnchored:                d.Get("ip_anchored").(bool),
		MicroTenantID:             d.Get("microtenant_id").(string),
		IsCnameEnabled:            d.Get("is_cname_enabled").(bool),
		SelectConnectorCloseToApp: d.Get("select_connector_close_to_app").(bool),
		UseInDrMode:               d.Get("use_in_dr_mode").(bool),
		TCPKeepAlive:              d.Get("tcp_keep_alive").(string),
		IsIncompleteDRConfig:      d.Get("is_incomplete_dr_config").(bool),
		DomainNames:               SetToStringList(d, "domain_names"),
		ServerGroups:              expandCommonServerGroups(d),
		CommonAppsDto:             expandCommonAppsDto(d),

		TCPAppPortRange: []common.NetworkPorts{},
		UDPAppPortRange: []common.NetworkPorts{},
	}
	remoteTCPAppPortRanges := []string{}
	remoteUDPAppPortRanges := []string{}
	if service != nil && id != "" {
		resource, _, err := applicationsegmentpra.Get(ctx, service, id)
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
		details.ServerGroups = expandCommonServerGroups(d)
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
		appConfigMap, ok := commonAppConfig.(map[string]interface{})
		if ok {
			// Automatically set `name` to match `domain` to prevent drift
			appConfigMap["name"] = appConfigMap["domain"].(string)

			appTypesSet, ok := appConfigMap["app_types"].(*schema.Set)
			if !ok {
				continue
			}
			appTypes := SetToStringSlice(appTypesSet)
			commonAppConfigDto = append(commonAppConfigDto, applicationsegmentpra.AppsConfig{
				AppID:               appConfigMap["app_id"].(string),
				PRAAppID:            appConfigMap["pra_app_id"].(string),
				Name:                appConfigMap["name"].(string),
				Enabled:             appConfigMap["enabled"].(bool),
				Domain:              appConfigMap["domain"].(string),
				ApplicationPort:     appConfigMap["application_port"].(string),
				ApplicationProtocol: appConfigMap["application_protocol"].(string),
				ConnectionSecurity:  appConfigMap["connection_security"].(string),
				AppTypes:            appTypes,
			})
		}
	}
	return commonAppConfigDto
}

func customizeDiffApplicationSegmentPRA(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// Clear any diffs related to pra_apps to prevent it from being included in the plan or state
	if d.HasChange("pra_apps") {
		d.Clear("pra_apps")
	}

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

func mapPRAAppsToCommonApps(d *schema.ResourceData, praApps []applicationsegmentpra.PRAApps) error {
	// If the API returned any PRA Apps, map them to common_apps_dto.apps_config
	if len(praApps) == 0 {
		return nil
	}

	// Create a single common_apps_dto with multiple apps_config blocks
	commonAppsConfig := make([]interface{}, len(praApps))
	for i, app := range praApps {
		commonAppMap := map[string]interface{}{
			"name":                 app.Name,
			"domain":               app.Domain,
			"application_protocol": app.ApplicationProtocol,
			"application_port":     app.ApplicationPort,
			"enabled":              app.Enabled,
			"app_types":            []interface{}{"SECURE_REMOTE_ACCESS"},
			"app_id":               app.AppID, // Populate app_id from praApps
			"connection_security":  app.ConnectionSecurity,
		}
		// Only set pra_app_id if it's present in the response
		if app.ID != "" {
			commonAppMap["pra_app_id"] = app.ID // Populate pra_app_id from praApps
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

func setAppIDsInCommonAppsDto(d *schema.ResourceData, praApps []applicationsegmentpra.PRAApps) error {
	if len(praApps) == 0 {
		return nil
	}

	// Extract app_id and pra_app_id from the first PRA app in the list
	appID := praApps[0].AppID
	praAppID := praApps[0].ID

	// Update the common_apps_dto with extracted app_id and pra_app_id values
	commonAppsDto := d.Get("common_apps_dto").(*schema.Set).List()
	if len(commonAppsDto) == 0 {
		return fmt.Errorf("common_apps_dto block is missing")
	}

	// Update the first entry in commonAppsDto.appsConfig with app_id and pra_app_id
	commonAppConfig := commonAppsDto[0].(map[string]interface{})
	appsConfig := commonAppConfig["apps_config"].(*schema.Set).List()

	if len(appsConfig) > 0 {
		appConfig := appsConfig[0].(map[string]interface{})
		appConfig["app_id"] = appID
		appConfig["pra_app_id"] = praAppID
	}

	// Write the updated config back to the resource data
	if err := d.Set("common_apps_dto", commonAppsDto); err != nil {
		return fmt.Errorf("failed to set common_apps_dto: %v", err)
	}

	return nil
}

/*
func flattenPRAApps(apps []applicationsegmentpra.PRAApps) []interface{} {
	if len(apps) == 0 {
		return []interface{}{}
	}

	appsConfig := make([]interface{}, len(apps))
	for i, app := range apps {
		appConfigMap := map[string]interface{}{
			"id":                   app.ID,
			"name":                 app.Name,
			"enabled":              app.Enabled,
			"application_port":     app.ApplicationPort,
			"application_protocol": app.ApplicationProtocol,
			"domain":               app.Domain,
			"app_id":               app.AppID,
			"hidden":               app.Hidden,
			"connection_security":  app.ConnectionSecurity,
			"microtenant_id":       app.MicroTenantID,
		}
		appsConfig[i] = appConfigMap
	}
	return appsConfig
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
*/

/*
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
*/

/*
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
*/
