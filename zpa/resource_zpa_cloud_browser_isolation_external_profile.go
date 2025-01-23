package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbiprofilecontroller"
)

func resourceCBIExternalProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCBIExternalProfileCreate,
		ReadContext:   resourceCBIExternalProfileRead,
		UpdateContext: resourceCBIExternalProfileUpdate,
		DeleteContext: resourceCBIExternalProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				resp, _, err := cbiprofilecontroller.GetByNameOrID(ctx, service, id)
				if err != nil {
					return nil, err
				}
				d.SetId(resp.ID)
				_ = d.Set("id", resp.ID)
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
			"banner_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "This field defines the list of region IDs.",
			},
			"certificate_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "This field defines the list of certificate IDs.",
			},
			"user_experience": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zgpu": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"forward_to_zia": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"organization_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"cloud_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"pac_file_url": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"browser_in_browser": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"persist_isolation_bar": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"translate": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"session_persistence": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"debug_mode": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"file_password": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"security_controls": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"copy_paste": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"none",
								"all",
							}, false),
						},
						"upload_download": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"none",
								"all",
								"upstream",
							}, false),
						},

						"document_viewer": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"local_render": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"allow_printing": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"restrict_keystrokes": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"deep_link": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"applications": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"flattened_pdf": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"watermark": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"show_user_id": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"show_timestamp": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"show_message": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"message": {
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
		},
	}
}

func resourceCBIExternalProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Validate the region_ids length
	regionIds := d.Get("region_ids").(*schema.Set).List()
	if len(regionIds) < 2 {
		return diag.FromErr(fmt.Errorf("expected region_ids to contain at least 2 items, got %d", len(regionIds)))
	}

	zClient := meta.(*Client)
	service := zClient.Service

	req := expandCBIExternalProfile(d)
	req.Regions = nil
	req.Certificates = nil
	req.Banner = nil
	log.Printf("[INFO] Creating cbi external profile with request\n%+v\n", req)
	cbiProfile, _, err := cbiprofilecontroller.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created cbi external profile request. ID: %v\n", cbiProfile)

	d.SetId(cbiProfile.ID)
	return resourceCBIExternalProfileRead(ctx, d, meta)
}

func resourceCBIExternalProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, _, err := cbiprofilecontroller.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing cbi profile %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}
	if resp.BannerID == "" && resp.Banner != nil && resp.Banner.ID != "" {
		resp.BannerID = resp.Banner.ID
	}

	if len(resp.CertificateIDs) == 0 {
		for _, c := range resp.Certificates {
			resp.CertificateIDs = append(resp.CertificateIDs, c.ID)
		}
	}

	if len(resp.RegionIDs) == 0 {
		for _, r := range resp.Regions {
			resp.RegionIDs = append(resp.RegionIDs, r.ID)
		}
	}
	log.Printf("[INFO] Getting cbi profile:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("banner_id", resp.BannerID)
	_ = d.Set("region_ids", resp.RegionIDs)
	_ = d.Set("certificate_ids", resp.CertificateIDs)

	if resp.SecurityControls != nil {
		_ = d.Set("security_controls", flattenSecurityControls(resp.SecurityControls))
	}

	if resp.UserExperience != nil {
		_ = d.Set("user_experience", flattenUserExperience(resp.UserExperience))
	}
	log.Printf("[INFO] Setting debug_mode: %+v\n", resp.DebugMode)

	if resp.DebugMode != nil {
		log.Printf("[INFO] Setting debug_mode: %+v\n", resp.DebugMode)
		_ = d.Set("debug_mode", flattenDebugMode(resp.DebugMode))
	}

	return nil
}

func resourceCBIExternalProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Validate the region_ids length
	regionIds := d.Get("region_ids").(*schema.Set).List()
	if len(regionIds) < 2 {
		return diag.FromErr(fmt.Errorf("expected region_ids to contain at least 2 items, got %d", len(regionIds)))
	}

	zClient := meta.(*Client)
	service := zClient.Service

	id := d.Id()
	log.Printf("[INFO] Updating cbi profile ID: %v\n", id)
	req := expandCBIExternalProfile(d)

	if _, _, err := cbiprofilecontroller.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := cbiprofilecontroller.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceCBIExternalProfileRead(ctx, d, meta)
}

func resourceCBIExternalProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	log.Printf("[INFO] Deleting cbi profile ID: %v\n", d.Id())

	if _, err := cbiprofilecontroller.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] cbi profile deleted")
	return nil
}

func expandCBIExternalProfile(d *schema.ResourceData) cbiprofilecontroller.IsolationProfile {
	cbiProfile := cbiprofilecontroller.IsolationProfile{
		ID:          d.Id(),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		BannerID:    d.Get("banner_id").(string),
		Banner: &cbiprofilecontroller.Banner{
			ID: d.Get("banner_id").(string),
		},
		RegionIDs:        SetToStringSlice(d.Get("region_ids").(*schema.Set)),
		CertificateIDs:   SetToStringSlice(d.Get("certificate_ids").(*schema.Set)),
		UserExperience:   expandCBIUserExperience(d),
		SecurityControls: expandCBISecurityControls(d),
		DebugMode:        expandCBIDebugMode(d),
	}
	profile := expandCBIUserExperience(d)
	if profile != nil {
		cbiProfile.UserExperience = profile
	}
	for _, id := range cbiProfile.RegionIDs {
		cbiProfile.Regions = append(cbiProfile.Regions, cbiprofilecontroller.Regions{ID: id})
	}
	for _, id := range cbiProfile.CertificateIDs {
		cbiProfile.Certificates = append(cbiProfile.Certificates, cbiprofilecontroller.Certificates{ID: id})
	}
	return cbiProfile
}

func expandCBIUserExperience(d *schema.ResourceData) *cbiprofilecontroller.UserExperience {
	profileObj, ok := d.GetOk("user_experience")
	if !ok {
		return nil
	}
	profiles := profileObj.([]interface{})
	if len(profiles) > 0 {
		profile := profiles[0].(map[string]interface{})

		forwardToZiaObj, forwardToZiaExists := profile["forward_to_zia"].([]interface{})
		var forwardToZia *cbiprofilecontroller.ForwardToZia
		if forwardToZiaExists && len(forwardToZiaObj) > 0 {
			forwardToZiaData := forwardToZiaObj[0].(map[string]interface{})
			forwardToZia = &cbiprofilecontroller.ForwardToZia{
				Enabled:        forwardToZiaData["enabled"].(bool),
				OrganizationID: forwardToZiaData["organization_id"].(string),
				CloudName:      forwardToZiaData["cloud_name"].(string),
				PacFileUrl:     forwardToZiaData["pac_file_url"].(string),
			}
		}

		return &cbiprofilecontroller.UserExperience{
			ZGPU:                profile["zgpu"].(bool),
			ForwardToZia:        forwardToZia,
			BrowserInBrowser:    profile["browser_in_browser"].(bool),
			PersistIsolationBar: profile["persist_isolation_bar"].(bool),
			Translate:           profile["translate"].(bool),
			SessionPersistence:  profile["session_persistence"].(bool),
		}
	}
	return nil
}

func expandCBISecurityControls(d *schema.ResourceData) *cbiprofilecontroller.SecurityControls {
	profileObj, ok := d.GetOk("security_controls")
	if !ok {
		return nil
	}
	profiles := profileObj.([]interface{})
	if len(profiles) > 0 {
		profile := profiles[0].(map[string]interface{})

		deepLinkObj, deepLinkExists := profile["deep_link"].([]interface{})
		var deepLink *cbiprofilecontroller.DeepLink
		if deepLinkExists && len(deepLinkObj) > 0 {
			deepLinkData := deepLinkObj[0].(map[string]interface{})
			deepLink = &cbiprofilecontroller.DeepLink{
				Enabled:      deepLinkData["enabled"].(bool),
				Applications: SetToStringSlice(deepLinkData["applications"].(*schema.Set)),
			}
		}

		watermarkObj, watermarkExists := profile["watermark"].([]interface{})
		var watermark *cbiprofilecontroller.Watermark
		if watermarkExists && len(watermarkObj) > 0 {
			watermarkData := watermarkObj[0].(map[string]interface{})
			watermark = &cbiprofilecontroller.Watermark{
				Enabled:       watermarkData["enabled"].(bool),
				ShowUserID:    watermarkData["show_user_id"].(bool),
				ShowTimestamp: watermarkData["show_timestamp"].(bool),
				ShowMessage:   watermarkData["show_message"].(bool),
				Message:       watermarkData["message"].(string),
			}
		}

		return &cbiprofilecontroller.SecurityControls{
			CopyPaste:          profile["copy_paste"].(string),
			UploadDownload:     profile["upload_download"].(string),
			DocumentViewer:     profile["document_viewer"].(bool),
			LocalRender:        profile["local_render"].(bool),
			AllowPrinting:      profile["allow_printing"].(bool),
			RestrictKeystrokes: profile["restrict_keystrokes"].(bool),
			DeepLink:           deepLink,
			FlattenedPdf:       profile["flattened_pdf"].(bool),
			Watermark:          watermark,
		}
	}
	return nil
}

func expandCBIDebugMode(d *schema.ResourceData) *cbiprofilecontroller.DebugMode {
	profileObj, ok := d.GetOk("debug_mode")
	if !ok {
		return nil
	}
	profiles := profileObj.([]interface{})
	if len(profiles) > 0 {
		profile := profiles[0].(map[string]interface{})

		return &cbiprofilecontroller.DebugMode{
			Allowed:      profile["allowed"].(bool),
			FilePassword: profile["file_password"].(string),
		}
	}
	return nil
}
