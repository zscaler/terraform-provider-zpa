package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/browser_protection"
)

func dataSourceBrowserProtection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBrowserProtectionRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_csp": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"criteria_flags_mask": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"criteria": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"finger_print_criteria": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"collect_location": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"fingerprint_timeout": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"browser": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"browser_eng": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"browser_eng_ver": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"browser_name": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"browser_version": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"canvas": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"flash_ver": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"fp_usr_agent_str": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"is_cookie": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"is_local_storage": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"is_sess_storage": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"ja3": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"mime": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"plugin": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"silverlight_ver": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"location": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"lat": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"lon": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"system": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"avail_screen_resolution": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"cpu_arch": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"curr_screen_resolution": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"font": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"java_ver": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"mobile_dev_type": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"monitor_mobile": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"os_name": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"os_version": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"sys_lang": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"tz": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"usr_lang": {
													Type:     schema.TypeBool,
													Computed: true,
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

func dataSourceBrowserProtectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *browser_protection.BrowserProtection

	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for browser protection profile name %s\n", name)
		res, _, err := browser_protection.GetBrowserProtectionProfileByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	} else {
		// If no name is provided, get the first profile (typically the default/active one)
		log.Printf("[INFO] Getting default browser protection profile")
		allProfiles, _, err := browser_protection.GetBrowserProtectionProfile(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(allProfiles) > 0 {
			resp = &allProfiles[0]
		}
	}

	if resp != nil {
		// Use actual ID if available, otherwise generate a short ID for the default profile
		var profileID string
		if resp.ID != "" {
			profileID = resp.ID
		} else {
			// Default profile "Zs Recommended profile" doesn't have an ID, so generate one
			profileID = generateShortID(resp.Name)
		}

		d.SetId(profileID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("default_csp", resp.DefaultCSP)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("criteria_flags_mask", resp.CriteriaFlagsMask)
		_ = d.Set("criteria", flattenCriteria(resp.Criteria))
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any browser protection profile with name '%s'", name))
	}

	return nil
}

// flattenCriteria flattens the Criteria struct into a format suitable for Terraform
func flattenCriteria(criteria browser_protection.Criteria) []map[string]interface{} {
	result := map[string]interface{}{
		"finger_print_criteria": flattenFingerPrintCriteria(criteria.FingerPrintCriteria),
	}

	return []map[string]interface{}{result}
}

// flattenFingerPrintCriteria flattens the FingerPrintCriteria struct
func flattenFingerPrintCriteria(fpc browser_protection.FingerPrintCriteria) []map[string]interface{} {
	result := map[string]interface{}{
		"collect_location":    fpc.CollectLocation,
		"fingerprint_timeout": fpc.FingerprintTimeout,
		"browser":             flattenBrowserCriteria(fpc.Browser),
		"location":            flattenLocationCriteria(fpc.Location),
		"system":              flattenSystemCriteria(fpc.System),
	}

	return []map[string]interface{}{result}
}

// flattenBrowserCriteria flattens the BrowserCriteria struct
func flattenBrowserCriteria(bc browser_protection.BrowserCriteria) []map[string]interface{} {
	result := map[string]interface{}{
		"browser_eng":      bc.BrowserEng,
		"browser_eng_ver":  bc.BrowserEngVer,
		"browser_name":     bc.BrowserName,
		"browser_version":  bc.BrowserVersion,
		"canvas":           bc.Canvas,
		"flash_ver":        bc.FlashVer,
		"fp_usr_agent_str": bc.FpUsrAgentStr,
		"is_cookie":        bc.IsCookie,
		"is_local_storage": bc.IsLocalStorage,
		"is_sess_storage":  bc.IsSessStorage,
		"ja3":              bc.Ja3,
		"mime":             bc.Mime,
		"plugin":           bc.Plugin,
		"silverlight_ver":  bc.SilverlightVer,
	}

	return []map[string]interface{}{result}
}

// flattenLocationCriteria flattens the LocationCriteria struct
func flattenLocationCriteria(lc browser_protection.LocationCriteria) []map[string]interface{} {
	result := map[string]interface{}{
		"lat": lc.Lat,
		"lon": lc.Lon,
	}

	return []map[string]interface{}{result}
}

// flattenSystemCriteria flattens the SystemCriteria struct
func flattenSystemCriteria(sc browser_protection.SystemCriteria) []map[string]interface{} {
	result := map[string]interface{}{
		"avail_screen_resolution": sc.AvailScreenResolution,
		"cpu_arch":                sc.CPUArch,
		"curr_screen_resolution":  sc.CurrScreenResolution,
		"font":                    sc.Font,
		"java_ver":                sc.JavaVer,
		"mobile_dev_type":         sc.MobileDevType,
		"monitor_mobile":          sc.MonitorMobile,
		"os_name":                 sc.OSName,
		"os_version":              sc.OSVersion,
		"sys_lang":                sc.SysLang,
		"tz":                      sc.Tz,
		"usr_lang":                sc.UsrLang,
	}

	return []map[string]interface{}{result}
}
