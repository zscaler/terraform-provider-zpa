package zpa

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/cloudbrowserisolation/cbiprofilecontroller"
)

func resourceCBIExternalProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceCBIExternalProfileCreate,
		Read:   resourceCBIExternalProfileRead,
		Update: resourceCBIExternalProfileUpdate,
		Delete: resourceCBIExternalProfileDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := zClient.cbiprofilecontroller.GetByName(id)
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
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"banner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "This field defines the list of server groups IDs.",
			},
			"certificate_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "This field defines the list of server groups IDs.",
			},
			"user_experience": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"session_persistence": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"browser_in_browser": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"security_controls": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"copy_paste": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"none",
								"all",
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
						"upload_download": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"none",
								"all",
							}, false),
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
					},
				},
			},
		},
	}
}

func resourceCBIExternalProfileCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandCBIExternalProfile(d)
	req.Regions = nil
	req.Certificates = nil
	req.Banner = nil
	log.Printf("[INFO] Creating cbi external profile with request\n%+v\n", req)
	cbiProfile, _, err := zClient.cbiprofilecontroller.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created cbi external profile request. ID: %v\n", cbiProfile)

	d.SetId(cbiProfile.ID)
	return resourceCBIExternalProfileRead(d, m)

}

func resourceCBIExternalProfileRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.cbiprofilecontroller.Get(d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing cbi profile %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
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
	return nil
}

func flattenRegionsSimple(regions *cbiprofilecontroller.IsolationProfile) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(regions.Regions))
	for i, group := range regions.Regions {
		ids[i] = group.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func flattenCertificatesSimple(certificates *cbiprofilecontroller.IsolationProfile) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(certificates.Certificates))
	for i, group := range certificates.Certificates {
		ids[i] = group.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func resourceCBIExternalProfileUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating cbi profile ID: %v\n", id)
	req := expandCBIExternalProfile(d)

	if _, _, err := zClient.cbiprofilecontroller.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := zClient.cbiprofilecontroller.Update(id, &req); err != nil {
		return err
	}

	return resourceCBIExternalProfileRead(d, m)
}

func resourceCBIExternalProfileDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting cbi profile ID: %v\n", d.Id())

	if _, err := zClient.cbiprofilecontroller.Delete(d.Id()); err != nil {
		return err
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
	profiles, ok := profileObj.(*schema.Set)
	if !ok {
		return nil
	}
	if len(profiles.List()) > 0 {
		profileObj := profiles.List()[0]
		profile, ok := profileObj.(map[string]interface{})
		if !ok {
			return nil
		}
		return &cbiprofilecontroller.UserExperience{
			SessionPersistence: profile["session_persistence"].(bool),
			BrowserInBrowser:   profile["browser_in_browser"].(bool),
		}
	}
	return nil
}

func expandCBISecurityControls(d *schema.ResourceData) *cbiprofilecontroller.SecurityControls {
	profileObj, ok := d.GetOk("security_controls")
	if !ok {
		return nil
	}
	profiles, ok := profileObj.(*schema.Set)
	if !ok {
		return nil
	}
	if len(profiles.List()) > 0 {
		profileObj := profiles.List()[0]
		profile, ok := profileObj.(map[string]interface{})
		if !ok {
			return nil
		}
		return &cbiprofilecontroller.SecurityControls{
			CopyPaste:          profile["copy_paste"].(string),
			DocumentViewer:     profile["document_viewer"].(bool),
			LocalRender:        profile["local_render"].(bool),
			UploadDownload:     profile["upload_download"].(string),
			AllowPrinting:      profile["allow_printing"].(bool),
			RestrictKeystrokes: profile["restrict_keystrokes"].(bool),
		}
	}
	return nil
}

func flattenUserExperience(experience *cbiprofilecontroller.UserExperience) interface{} {
	return []map[string]interface{}{
		{
			"session_persistence": experience.SessionPersistence,
			"browser_in_browser":  experience.BrowserInBrowser,
		},
	}
}
