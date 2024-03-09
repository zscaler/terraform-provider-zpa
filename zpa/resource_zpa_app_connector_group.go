package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontroller"
)

func resourceAppConnectorGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppConnectorGroupCreate,
		Read:   resourceAppConnectorGroupRead,
		Update: resourceAppConnectorGroupUpdate,
		Delete: resourceAppConnectorGroupDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				service := m.(*Client).appconnectorgroup.WithMicroTenant(GetString(d.Get("microtenant_id")))

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := service.GetByName(id)
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
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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

func resourceAppConnectorGroupCreate(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).appconnectorgroup.WithMicroTenant(GetString(d.Get("microtenant_id")))

	if err := validateAndSetProfileNameID(d); err != nil {
		return err
	}
	req := expandAppConnectorGroup(d)
	log.Printf("[INFO] Creating zpa app connector group with request\n%+v\n", req)

	if err := validateTCPQuickAck(req); err != nil {
		return err
	}

	resp, _, err := service.Create(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created app connector group request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceAppConnectorGroupRead(d, m)
}

func resourceAppConnectorGroupRead(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).appconnectorgroup.WithMicroTenant(GetString(d.Get("microtenant_id")))

	resp, _, err := service.Get(d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing app connector group %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting application server:\n%+v\n", resp)
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

func resourceAppConnectorGroupUpdate(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).appconnectorgroup.WithMicroTenant(GetString(d.Get("microtenant_id")))

	if err := validateAndSetProfileNameID(d); err != nil {
		return err
	}
	id := d.Id()
	log.Printf("[INFO] Updating app connector group ID: %v\n", id)
	req := expandAppConnectorGroup(d)

	if err := validateTCPQuickAck(req); err != nil {
		return err
	}

	if _, _, err := service.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := service.Update(id, &req); err != nil {
		return err
	}

	return resourceAppConnectorGroupRead(d, m)
}

func detachAppConnectorGroupFromAllAccessPolicyRules(d *schema.ResourceData, policySetControllerService *policysetcontroller.Service) {
	policyRulesDetchLock.Lock()
	defer policyRulesDetchLock.Unlock()
	accessPolicySet, _, err := policySetControllerService.GetByPolicyType("ACCESS_POLICY")
	if err != nil {
		return
	}
	rules, _, err := policySetControllerService.GetAllByType("ACCESS_POLICY")
	if err != nil {
		return
	}
	for _, rule := range rules {
		ids := []policysetcontroller.AppConnectorGroups{}
		changed := false
		for _, app := range rule.AppConnectorGroups {
			if app.ID == d.Id() {
				changed = true
				continue
			}
			ids = append(ids, policysetcontroller.AppConnectorGroups{
				ID: app.ID,
			})
		}
		rule.AppConnectorGroups = ids
		if changed {
			if _, err := policySetControllerService.WithMicroTenant(GetString(d.Get("microtenant_id"))).UpdateRule(accessPolicySet.ID, rule.ID, &rule); err != nil {
				continue
			}
		}
	}
}

func resourceAppConnectorGroupDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	policySetControllerService := m.(*Client).policysetcontroller.WithMicroTenant(GetString(d.Get("microtenant_id")))
	service := zClient.appconnectorgroup.WithMicroTenant(GetString(d.Get("microtenant_id")))
	log.Printf("[INFO] Deleting app connector groupID: %v\n", d.Id())

	// detach app connector group from all access policy rules
	detachAppConnectorGroupFromAllAccessPolicyRules(d, policySetControllerService)

	if _, err := service.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] app connector group deleted")
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
