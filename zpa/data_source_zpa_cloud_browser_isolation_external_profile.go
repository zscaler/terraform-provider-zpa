package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/cloudbrowserisolation/cbiprofilecontroller"
)

func dataSourceCBIExternalProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCBIExternalProfileRead,
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
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"href": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"regions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of Cloud Browser Isolation Regions",
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
			"security_controls": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"copy_paste": {
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
						"upload_download": {
							Type:     schema.TypeString,
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
					},
				},
			},
			"user_experience": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"session_persistence": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"browser_in_browser": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCBIExternalProfileRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *cbiprofilecontroller.IsolationProfile
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for cbi external profile %s\n", id)
		res, _, err := zClient.cbiprofilecontroller.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data for cbi external profile name %s\n", name)
		res, _, err := zClient.cbiprofilecontroller.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
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

	} else {
		return fmt.Errorf("couldn't find any cbi external profile with name '%s' or id '%s'", name, id)
	}

	return nil
}

func flattenRegions(cbiIsolationProfile *cbiprofilecontroller.IsolationProfile) []interface{} {
	profiles := make([]interface{}, len(cbiIsolationProfile.Regions))
	for i, profileItem := range cbiIsolationProfile.Regions {
		profiles[i] = map[string]interface{}{
			"id":   profileItem.ID,
			"name": profileItem.Name,
		}
	}

	return profiles
}

func flattenSecurityControls(securityControls *cbiprofilecontroller.SecurityControls) []map[string]interface{} {
	result := make([]map[string]interface{}, 1)
	result[0] = make(map[string]interface{})
	result[0]["copy_paste"] = securityControls.CopyPaste
	result[0]["document_viewer"] = securityControls.DocumentViewer
	result[0]["local_render"] = securityControls.LocalRender
	result[0]["upload_download"] = securityControls.UploadDownload
	result[0]["allow_printing"] = securityControls.AllowPrinting
	return result
}

func flattenUserExperience(experience *cbiprofilecontroller.UserExperience) interface{} {
	return []map[string]interface{}{
		{
			"session_persistence": experience.SessionPersistence,
			"browser_in_browser":  experience.BrowserInBrowser,
		},
	}
}
