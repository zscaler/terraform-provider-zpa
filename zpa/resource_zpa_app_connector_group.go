package zpa

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAppConnectorGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppConnectorGroupCreate,
		ReadContext:   resourceAppConnectorGroupRead,
		UpdateContext: resourceAppConnectorGroupUpdate,
		DeleteContext: resourceAppConnectorGroupDelete,
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
					resp, _, err := appconnectorgroup.GetByName(ctx, service, id)
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the App Connector Group",
			},
			"city_country": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"country_code": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateCountryCode,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the App Connector Group",
			},
			"dns_query_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "IPV4_IPV6",
				Description: "Whether to enable IPv4 or IPv6, or both, for DNS resolution of all applications in the App Connector Group",
				ValidateFunc: validation.StringInSlice([]string{
					"IPV4_IPV6",
					"IPV4",
					"IPV6",
				}, false),
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether this App Connector Group is enabled or not",
			},
			"latitude": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     ValidateLatitude,
				DiffSuppressFunc: DiffSuppressFuncCoordinate,
				Description:      "Latitude of the App Connector Group. Integer or decimal. With values in the range of -90 to 90",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Location of the App Connector Group",
			},
			"longitude": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     ValidateLongitude,
				DiffSuppressFunc: DiffSuppressFuncCoordinate,
				Description:      "Longitude of the App Connector Group. Integer or decimal. With values in the range of -180 to 180",
			},
			"lss_app_connector_group": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether or not the App Connector Group is configured for the Log Streaming Service (LSS)",
			},
			"tcp_quick_ack_app": {
				Description: "Whether TCP Quick Acknowledgement is enabled or disabled for the application. The tcpQuickAckApp, tcpQuickAckAssistant, and tcpQuickAckReadAssistant fields must all share the same value.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"tcp_quick_ack_assistant": {
				Description: "Whether TCP Quick Acknowledgement is enabled or disabled for the application. The tcpQuickAckApp, tcpQuickAckAssistant, and tcpQuickAckReadAssistant fields must all share the same value.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"tcp_quick_ack_read_assistant": {
				Description: "Whether TCP Quick Acknowledgement is enabled or disabled for the application. The tcpQuickAckApp, tcpQuickAckAssistant, and tcpQuickAckReadAssistant fields must all share the same value.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"use_in_dr_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"pra_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"waf_disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"upgrade_day": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "SUNDAY",
				Description: "App Connectors in this group will attempt to update to a newer version of the software during this specified day. List of valid days (i.e., Sunday, Monday)",
			},
			"upgrade_time_in_secs": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "66600",
				Description: "App Connectors in this group will attempt to update to a newer version of the software during this specified time. Integer in seconds (i.e., -66600). The integer should be greater than or equal to 0 and less than 86400, in 15 minute intervals",
			},
			"override_version_profile": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the default version profile of the App Connector Group is applied or overridden. Supported values: true, false",
			},
			"version_profile_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the version profile. To learn more, see Version Profile Use Cases. This value is required, if the value for overrideVersionProfile is set to true",
				ValidateFunc: validation.StringInSlice([]string{
					"Default", "Previous Default", "New Release",
				}, false),
			},
			"version_profile_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the version profile. To learn more, see Version Profile Use Cases. This value is required, if the value for overrideVersionProfile is set to true",
				ValidateFunc: validation.StringInSlice([]string{
					"0", "1", "2",
				}, false),
			},
		},
	}
}

func resourceAppConnectorGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	if err := validateAndSetProfileNameID(d); err != nil {
		return diag.FromErr(err)
	}
	req := expandAppConnectorGroup(d)
	log.Printf("[INFO] Creating zpa app connector group with request\n%+v\n", req)

	if err := validateTCPQuickAck(req); err != nil {
		return diag.FromErr(err)
	}

	resp, _, err := appconnectorgroup.Create(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created app connector group request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceAppConnectorGroupRead(ctx, d, meta)
}

func resourceAppConnectorGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := appconnectorgroup.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing app connector group %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting app connector group:\n%+v\n", resp)
	_ = d.Set("name", resp.Name)
	_ = d.Set("city_country", resp.CityCountry)
	_ = d.Set("country_code", resp.CountryCode)
	_ = d.Set("description", resp.Description)
	_ = d.Set("dns_query_type", resp.DNSQueryType)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("latitude", resp.Latitude)
	_ = d.Set("longitude", resp.Longitude)
	_ = d.Set("location", resp.Location)
	_ = d.Set("lss_app_connector_group", resp.LSSAppConnectorGroup)
	_ = d.Set("tcp_quick_ack_app", resp.TCPQuickAckApp)
	_ = d.Set("tcp_quick_ack_assistant", resp.TCPQuickAckAssistant)
	_ = d.Set("tcp_quick_ack_read_assistant", resp.TCPQuickAckReadAssistant)
	_ = d.Set("use_in_dr_mode", resp.UseInDrMode)
	_ = d.Set("upgrade_day", resp.UpgradeDay)
	_ = d.Set("upgrade_time_in_secs", resp.UpgradeTimeInSecs)
	_ = d.Set("override_version_profile", resp.OverrideVersionProfile)
	_ = d.Set("pra_enabled", resp.PRAEnabled)
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	_ = d.Set("version_profile_name", resp.VersionProfileName)
	_ = d.Set("version_profile_id", resp.VersionProfileID)
	_ = d.Set("waf_disabled", resp.WAFDisabled)
	return nil
}

func resourceAppConnectorGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	if err := validateAndSetProfileNameID(d); err != nil {
		return diag.FromErr(err)
	}
	id := d.Id()
	log.Printf("[INFO] Updating app connector group ID: %v\n", id)
	req := expandAppConnectorGroup(d)

	if err := validateTCPQuickAck(req); err != nil {
		return diag.FromErr(err)
	}

	if _, _, err := appconnectorgroup.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := appconnectorgroup.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceAppConnectorGroupRead(ctx, d, meta)
}

func resourceAppConnectorGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Use MicroTenant if available
	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting app connector group ID: %v\n", d.Id())

	// Detach app connector group from all access policy rules
	if err := detachAppConnectorGroupFromAllAccessPolicyRules(ctx, d, service); err != nil {
		return diag.FromErr(err)
	}

	// Call Delete with context and necessary parameters
	if _, err := appconnectorgroup.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[INFO] App connector group deleted successfully")
	return nil
}

func expandAppConnectorGroup(d *schema.ResourceData) appconnectorgroup.AppConnectorGroup {
	appConnectorGroup := appconnectorgroup.AppConnectorGroup{
		ID:                       d.Get("id").(string),
		Name:                     d.Get("name").(string),
		CityCountry:              d.Get("city_country").(string),
		CountryCode:              d.Get("country_code").(string),
		Description:              d.Get("description").(string),
		DNSQueryType:             d.Get("dns_query_type").(string),
		Enabled:                  d.Get("enabled").(bool),
		Latitude:                 d.Get("latitude").(string),
		Longitude:                d.Get("longitude").(string),
		Location:                 d.Get("location").(string),
		LSSAppConnectorGroup:     d.Get("lss_app_connector_group").(bool),
		TCPQuickAckApp:           d.Get("tcp_quick_ack_app").(bool),
		TCPQuickAckAssistant:     d.Get("tcp_quick_ack_assistant").(bool),
		TCPQuickAckReadAssistant: d.Get("tcp_quick_ack_read_assistant").(bool),
		UseInDrMode:              d.Get("use_in_dr_mode").(bool),
		UpgradeDay:               d.Get("upgrade_day").(string),
		UpgradeTimeInSecs:        d.Get("upgrade_time_in_secs").(string),
		OverrideVersionProfile:   d.Get("override_version_profile").(bool),
		PRAEnabled:               d.Get("pra_enabled").(bool),
		WAFDisabled:              d.Get("waf_disabled").(bool),
		MicroTenantID:            d.Get("microtenant_id").(string),
		VersionProfileID:         d.Get("version_profile_id").(string),
		VersionProfileName:       d.Get("version_profile_name").(string),
	}
	return appConnectorGroup
}

func detachAppConnectorGroupFromAllAccessPolicyRules(ctx context.Context, d *schema.ResourceData, service *zscaler.Service) error {
	policyRulesDetchLock.Lock()
	defer policyRulesDetchLock.Unlock()

	// Retrieve access policy set
	accessPolicySet, _, err := policysetcontroller.GetByPolicyType(ctx, service, "ACCESS_POLICY")
	if err != nil {
		return fmt.Errorf("failed to get access policy set: %w", err)
	}

	// Retrieve all access policy rules
	rules, _, err := policysetcontroller.GetAllByType(ctx, service, "ACCESS_POLICY")
	if err != nil {
		return fmt.Errorf("failed to get access policy rules: %w", err)
	}

	// Iterate over rules and detach the app connector group
	for _, rule := range rules {
		var updatedGroups []appconnectorgroup.AppConnectorGroup
		changed := false

		for _, group := range rule.AppConnectorGroups {
			if group.ID == d.Id() {
				changed = true
				continue
			}
			updatedGroups = append(updatedGroups, appconnectorgroup.AppConnectorGroup{ID: group.ID})
		}

		// If the rule was modified, update it
		if changed {
			rule.AppConnectorGroups = updatedGroups
			if _, err := policysetcontroller.UpdateRule(ctx, service, accessPolicySet.ID, rule.ID, &rule); err != nil {
				log.Printf("[WARN] Failed to update policy rule %s: %v", rule.ID, err)
				// Continue processing other rules despite errors
			}
		}
	}

	return nil
}

func validateTCPQuickAck(tcp appconnectorgroup.AppConnectorGroup) error {
	if tcp.TCPQuickAckApp != tcp.TCPQuickAckAssistant {
		return fmt.Errorf("the values of tcpQuickAck related flags need to be consistent")
	}
	if tcp.TCPQuickAckApp != tcp.TCPQuickAckReadAssistant {
		return fmt.Errorf("the values of tcpQuickAck related flags need to be consistent")
	}
	if tcp.TCPQuickAckAssistant != tcp.TCPQuickAckReadAssistant {
		return fmt.Errorf("the values of tcpQuickAck related flags need to be consistent")
	}
	return nil
}
