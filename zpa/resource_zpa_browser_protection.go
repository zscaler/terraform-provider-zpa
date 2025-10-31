package zpa

/*
import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/browser_protection"
)

func resourceBrowserProtection() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceBrowserProtectionRead,
		CreateContext: resourceBrowserProtectionCreate,
		UpdateContext: resourceBrowserProtectionUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				diags := resourceBrowserProtectionRead(ctx, d, meta)
				if diags.HasError() {
					return nil, fmt.Errorf("failed to read atp malware policy import: %s", diags[0].Summary)
				}
				d.SetId("browser_protection")
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"criteria_flags_mask": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_csp": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"criteria": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"finger_print_criteria": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"collect_location": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"fingerprint_timeout": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"browser": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"browser_eng": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"browser_eng_ver": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"browser_name": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"browser_version": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"canvas": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"flash_ver": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"fp_usr_agent_str": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"is_cookie": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"is_local_storage": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"is_sess_storage": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"ja3": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"mime": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"plugin": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"silverlight_ver": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
									"location": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"lat": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"lon": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
									"system": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"avail_screen_resolution": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"cpu_arch": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"curr_screen_resolution": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"font": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"java_ver": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"mobile_dev_type": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"monitor_mobile": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"os_name": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"os_version": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"sys_lang": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"tz": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"usr_lang": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceBrowserProtectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandBrowserProtection(d)
	log.Printf("[INFO] Creating browser protection profile: %s\n", req.Name)

	// Since there's no Create function in the SDK, we'll use UpdateBrowserProtectionProfile
	// This assumes the profile already exists and we're updating it
	_, err := browser_protection.UpdateBrowserProtectionProfile(ctx, service, req.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(req.ID)

	return resourceBrowserProtectionRead(ctx, d, meta)
}

func resourceBrowserProtectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	profileID := d.Id()
	log.Printf("[INFO] Reading browser protection profile: %s\n", profileID)

	// Get all browser protection profiles and find the one with matching ID
	allProfiles, _, err := browser_protection.GetBrowserProtectionProfile(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	var resp *browser_protection.BrowserProtection
	for _, profile := range allProfiles {
		if profile.ID == profileID {
			resp = &profile
			break
		}
	}

	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("criteria_flags_mask", resp.CriteriaFlagsMask)
		_ = d.Set("default_csp", resp.DefaultCSP)
		_ = d.Set("criteria", flattenCriteria(resp.Criteria))
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find browser protection profile with id '%s'", profileID))
	}

	return nil
}

func resourceBrowserProtectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	profileID := d.Id()
	log.Printf("[INFO] Updating browser protection profile: %s\n", profileID)

	// Check if the profile still exists before updating
	allProfiles, _, err := browser_protection.GetBrowserProtectionProfile(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	var profileExists bool
	for _, profile := range allProfiles {
		if profile.ID == profileID {
			profileExists = true
			break
		}
	}

	if !profileExists {
		log.Printf("[WARN] Browser protection profile %s no longer exists, removing from state", profileID)
		d.SetId("")
		return nil
	}

	// Since there's no Update function in the SDK, we'll use UpdateBrowserProtectionProfile
	// This sets the profile as active
	_, err = browser_protection.UpdateBrowserProtectionProfile(ctx, service, profileID)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceBrowserProtectionRead(ctx, d, meta)
}

func expandBrowserProtection(d *schema.ResourceData) browser_protection.BrowserProtection {
	result := browser_protection.BrowserProtection{
		ID:                d.Id(),
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		CriteriaFlagsMask: d.Get("criteria_flags_mask").(string),
		DefaultCSP:        d.Get("default_csp").(bool),
		Criteria:          expandCriteria(d.Get("criteria")),
	}
	return result
}

// expandCriteria expands the criteria from Terraform schema to SDK struct
func expandCriteria(criteria interface{}) browser_protection.Criteria {
	if criteria == nil {
		return browser_protection.Criteria{}
	}

	criteriaList := criteria.([]interface{})
	if len(criteriaList) == 0 {
		return browser_protection.Criteria{}
	}

	criteriaMap := criteriaList[0].(map[string]interface{})

	return browser_protection.Criteria{
		FingerPrintCriteria: expandFingerPrintCriteria(criteriaMap["finger_print_criteria"]),
	}
}

// expandFingerPrintCriteria expands the fingerprint criteria
func expandFingerPrintCriteria(fpc interface{}) browser_protection.FingerPrintCriteria {
	if fpc == nil {
		return browser_protection.FingerPrintCriteria{}
	}

	fpcList := fpc.([]interface{})
	if len(fpcList) == 0 {
		return browser_protection.FingerPrintCriteria{}
	}

	fpcMap := fpcList[0].(map[string]interface{})

	return browser_protection.FingerPrintCriteria{
		CollectLocation:    fpcMap["collect_location"].(bool),
		FingerprintTimeout: fpcMap["fingerprint_timeout"].(string),
		Browser:            expandBrowserCriteria(fpcMap["browser"]),
		Location:           expandLocationCriteria(fpcMap["location"]),
		System:             expandSystemCriteria(fpcMap["system"]),
	}
}

// expandBrowserCriteria expands the browser criteria
func expandBrowserCriteria(bc interface{}) browser_protection.BrowserCriteria {
	if bc == nil {
		return browser_protection.BrowserCriteria{}
	}

	bcList := bc.([]interface{})
	if len(bcList) == 0 {
		return browser_protection.BrowserCriteria{}
	}

	bcMap := bcList[0].(map[string]interface{})

	return browser_protection.BrowserCriteria{
		BrowserEng:     bcMap["browser_eng"].(bool),
		BrowserEngVer:  bcMap["browser_eng_ver"].(bool),
		BrowserName:    bcMap["browser_name"].(bool),
		BrowserVersion: bcMap["browser_version"].(bool),
		Canvas:         bcMap["canvas"].(bool),
		FlashVer:       bcMap["flash_ver"].(bool),
		FpUsrAgentStr:  bcMap["fp_usr_agent_str"].(bool),
		IsCookie:       bcMap["is_cookie"].(bool),
		IsLocalStorage: bcMap["is_local_storage"].(bool),
		IsSessStorage:  bcMap["is_sess_storage"].(bool),
		Ja3:            bcMap["ja3"].(bool),
		Mime:           bcMap["mime"].(bool),
		Plugin:         bcMap["plugin"].(bool),
		SilverlightVer: bcMap["silverlight_ver"].(bool),
	}
}

// expandLocationCriteria expands the location criteria
func expandLocationCriteria(lc interface{}) browser_protection.LocationCriteria {
	if lc == nil {
		return browser_protection.LocationCriteria{}
	}

	lcList := lc.([]interface{})
	if len(lcList) == 0 {
		return browser_protection.LocationCriteria{}
	}

	lcMap := lcList[0].(map[string]interface{})

	return browser_protection.LocationCriteria{
		Lat: lcMap["lat"].(bool),
		Lon: lcMap["lon"].(bool),
	}
}

// expandSystemCriteria expands the system criteria
func expandSystemCriteria(sc interface{}) browser_protection.SystemCriteria {
	if sc == nil {
		return browser_protection.SystemCriteria{}
	}

	scList := sc.([]interface{})
	if len(scList) == 0 {
		return browser_protection.SystemCriteria{}
	}

	scMap := scList[0].(map[string]interface{})

	return browser_protection.SystemCriteria{
		AvailScreenResolution: scMap["avail_screen_resolution"].(bool),
		CPUArch:               scMap["cpu_arch"].(bool),
		CurrScreenResolution:  scMap["curr_screen_resolution"].(bool),
		Font:                  scMap["font"].(bool),
		JavaVer:               scMap["java_ver"].(bool),
		MobileDevType:         scMap["mobile_dev_type"].(bool),
		MonitorMobile:         scMap["monitor_mobile"].(bool),
		OSName:                scMap["os_name"].(bool),
		OSVersion:             scMap["os_version"].(bool),
		SysLang:               scMap["sys_lang"].(bool),
		Tz:                    scMap["tz"].(bool),
		UsrLang:               scMap["usr_lang"].(bool),
	}
}
*/
