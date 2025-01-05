package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbiprofilecontroller"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCBIExternalProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCBIExternalProfileRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"banner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"href": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"certificate_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "This field defines the list of certificate IDs.",
			},
			"user_experience": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zgpu": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"forward_to_zia": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"organization_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cloud_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pac_file_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"browser_in_browser": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"persist_isolation_bar": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"translate": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"session_persistence": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"debug_mode": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"file_password": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"security_controls": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"copy_paste": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"upload_download": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"document_viewer": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"local_render": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"allow_printing": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"restrict_keystrokes": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"deep_link": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"applications": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"flattened_pdf": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"watermark": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"show_user_id": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"show_timestamp": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"show_message": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"message": {
										Type:     schema.TypeString,
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

func dataSourceCBIExternalProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *cbiprofilecontroller.IsolationProfile
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for CBI external profile ID: %s\n", id)
		res, _, err := cbiprofilecontroller.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data for CBI external profile name: %s\n", name)
		res, _, err := cbiprofilecontroller.GetByNameOrID(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		log.Printf("[INFO] CBI external profile response: %+v\n", resp)
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("is_default", resp.IsDefault)
		_ = d.Set("href", resp.Href)
		_ = d.Set("regions", flattenRegions(resp))
		if resp.SecurityControls != nil {
			_ = d.Set("security_controls", flattenSecurityControls(resp.SecurityControls))
		}
		if resp.UserExperience != nil {
			_ = d.Set("user_experience", flattenUserExperience(resp.UserExperience))
		}
		if resp.DebugMode != nil {
			log.Printf("[INFO] Setting debug_mode: %+v\n", resp.DebugMode)
			_ = d.Set("debug_mode", flattenDebugMode(resp.DebugMode))
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any CBI external profile with name '%s' or id '%s'", name, id))
	}

	return nil
}

func flattenRegions(cbiIsolationProfile *cbiprofilecontroller.IsolationProfile) []interface{} {
	regions := make([]interface{}, len(cbiIsolationProfile.Regions))
	for i, regionItem := range cbiIsolationProfile.Regions {
		regions[i] = map[string]interface{}{
			"id":   regionItem.ID,
			"name": regionItem.Name,
		}
	}

	return regions
}

func flattenSecurityControls(securityControls *cbiprofilecontroller.SecurityControls) []interface{} {
	if securityControls == nil {
		return []interface{}{}
	}

	deepLink := []interface{}{}
	if securityControls.DeepLink != nil {
		deepLink = append(deepLink, map[string]interface{}{
			"enabled":      securityControls.DeepLink.Enabled,
			"applications": securityControls.DeepLink.Applications,
		})
	}

	watermark := []interface{}{}
	if securityControls.Watermark != nil {
		watermark = append(watermark, map[string]interface{}{
			"enabled":        securityControls.Watermark.Enabled,
			"show_user_id":   securityControls.Watermark.ShowUserID,
			"show_timestamp": securityControls.Watermark.ShowTimestamp,
			"show_message":   securityControls.Watermark.ShowMessage,
			"message":        securityControls.Watermark.Message,
		})
	}

	return []interface{}{
		map[string]interface{}{
			"copy_paste":          securityControls.CopyPaste,
			"upload_download":     securityControls.UploadDownload,
			"document_viewer":     securityControls.DocumentViewer,
			"local_render":        securityControls.LocalRender,
			"allow_printing":      securityControls.AllowPrinting,
			"restrict_keystrokes": securityControls.RestrictKeystrokes,
			"deep_link":           deepLink,
			"flattened_pdf":       securityControls.FlattenedPdf,
			"watermark":           watermark,
		},
	}
}

func flattenUserExperience(userExperience *cbiprofilecontroller.UserExperience) []interface{} {
	if userExperience == nil {
		return []interface{}{}
	}

	forwardToZia := []interface{}{}
	if userExperience.ForwardToZia != nil {
		forwardToZia = append(forwardToZia, map[string]interface{}{
			"enabled":         userExperience.ForwardToZia.Enabled,
			"organization_id": userExperience.ForwardToZia.OrganizationID,
			"cloud_name":      userExperience.ForwardToZia.CloudName,
			"pac_file_url":    userExperience.ForwardToZia.PacFileUrl,
		})
	}

	return []interface{}{
		map[string]interface{}{
			"zgpu":                  userExperience.ZGPU,
			"forward_to_zia":        forwardToZia,
			"browser_in_browser":    userExperience.BrowserInBrowser,
			"persist_isolation_bar": userExperience.PersistIsolationBar,
			"translate":             userExperience.Translate,
			"session_persistence":   userExperience.SessionPersistence,
		},
	}
}

func flattenDebugMode(debugMode *cbiprofilecontroller.DebugMode) []interface{} {
	if debugMode == nil {
		log.Printf("[INFO] No debug mode data found")
		return []interface{}{}
	}

	log.Printf("[INFO] Flattening debug mode data: %+v\n", debugMode)
	return []interface{}{
		map[string]interface{}{
			"allowed":       debugMode.Allowed,
			"file_password": debugMode.FilePassword,
		},
	}
}
