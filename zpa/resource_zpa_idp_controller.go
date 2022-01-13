package zpa

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/client"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/idpcontroller"
)

func resourceIdpController() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdpControllerCreate,
		Read:   resourceIdpControllerRead,
		Update: resourceIdpControllerUpdate,
		Delete: resourceIdpControllerDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					d.Set("id", id)
				} else {
					resp, _, err := zClient.servergroup.GetByName(id)
					if err == nil {
						d.SetId(resp.ID)
						d.Set("id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"admin_metadata": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sp_base_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sp_entity_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sp_metadata_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sp_post_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"certificates": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cname": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"certificate": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"serial_no": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"valid_from_in_sec": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"valid_to_in_sec": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"admin_sp_signing_cert_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_provision": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disable_saml_based_policy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"domain_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_scim_based_policy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Default value if null is True",
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"idp_entity_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"login_name_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"login_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reauth_on_user_update": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"redirect_binding": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"zpa_saml_request": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scim_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"scim_service_provider_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scim_shared_secret": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scim_shared_secret_exists": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sign_saml_request": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sso_type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"use_custom_sp_metadata": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"user_metadata": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sp_base_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sp_entity_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sp_metadata_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sp_post_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"user_sp_signing_cert_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceIdpControllerCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandIdpController(d)
	log.Printf("[INFO] Creating IdP Controller with request\n%+v\n", req)

	idpcontroller, _, err := zClient.idpcontroller.Create(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created IdP Controller request. ID: %v\n", idpcontroller)

	d.SetId(idpcontroller.ID)
	return resourceIdpControllerRead(d, m)

}

func resourceIdpControllerRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.idpcontroller.Get(d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing IdP Controller %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting Idp Controller:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("admin_sp_signing_cert_id", resp.AdminSpSigningCertID)
	_ = d.Set("auto_provision", resp.AutoProvision)
	_ = d.Set("creation_time", resp.CreationTime)
	_ = d.Set("description", resp.Description)
	_ = d.Set("disable_saml_based_policy", resp.DisableSamlBasedPolicy)
	_ = d.Set("domain_list", resp.DomainList)
	_ = d.Set("enable_scim_based_policy", resp.EnableScimBasedPolicy)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("idp_entity_id", resp.IdpEntityID)
	_ = d.Set("login_name_attribute", resp.LoginNameAttribute)
	_ = d.Set("login_url", resp.LoginURL)
	_ = d.Set("modifiedby", resp.ModifiedBy)
	_ = d.Set("modified_time", resp.ModifiedTime)
	_ = d.Set("name", resp.Name)
	_ = d.Set("reauth_on_user_update", resp.ReauthOnUserUpdate)
	_ = d.Set("redirect_binding", resp.RedirectBinding)
	_ = d.Set("scim_enabled", resp.ScimEnabled)
	_ = d.Set("zpa_saml_request", resp.ZPASAMLRequest)
	_ = d.Set("scim_service_provider_endpoint", resp.ScimServiceProviderEndpoint)
	_ = d.Set("scim_shared_secret", resp.ScimSharedSecret)
	_ = d.Set("scim_shared_secret_exists", resp.ScimSharedSecretExists)
	_ = d.Set("sign_saml_request", resp.SignSamlRequest)
	_ = d.Set("sso_type", resp.SsoType)
	_ = d.Set("use_custom_sp_metadata", resp.UseCustomSpMetadata)
	_ = d.Set("user_sp_signing_cert_id", resp.UserSpSigningCertID)
	// if resp.UserMetadata != nil {
	// 	_ = d.Set("user_metadata", flattenUserMeta(resp.UserMetadata))
	// }
	// if resp.AdminMetadata != nil {
	// 	_ = d.Set("admin_metadata", flattenAdminMeta(resp.AdminMetadata))
	// }

	return nil
}

func resourceIdpControllerUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating Idp Controller ID: %v\n", id)
	req := expandIdpController(d)

	if _, err := zClient.idpcontroller.Update(id, &req); err != nil {
		return err
	}

	return resourceIdpControllerRead(d, m)
}

func resourceIdpControllerDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting Idp Controller ID: %v\n", d.Id())

	if _, err := zClient.idpcontroller.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] Idp Controller deleted")
	return nil
}

func expandIdpController(d *schema.ResourceData) idpcontroller.IdpController {
	idpController := idpcontroller.IdpController{
		AdminSpSigningCertID:        d.Get("admin_sp_signing_cert_id").(string),
		UserSpSigningCertID:         d.Get("user_sp_signing_cert_id").(string),
		AutoProvision:               d.Get("auto_provision").(string),
		Description:                 d.Get("description").(string),
		DisableSamlBasedPolicy:      d.Get("disable_saml_based_policy").(bool),
		DomainList:                  ListToStringSlice(d.Get("domain_list").([]interface{})), // Need to check format
		EnableScimBasedPolicy:       d.Get("enable_scim_based_policy").(bool),
		Enabled:                     d.Get("enabled").(bool),
		Name:                        d.Get("name").(string),
		IdpEntityID:                 d.Get("idp_entity_id").(string),
		LoginNameAttribute:          d.Get("login_name_attribute").(string),
		LoginURL:                    d.Get("login_url").(string),
		ReauthOnUserUpdate:          d.Get("reauth_on_user_update").(bool),
		RedirectBinding:             d.Get("redirect_binding").(bool),
		ScimEnabled:                 d.Get("scim_enabled").(bool),
		ZPASAMLRequest:              d.Get("zpa_saml_request").(string),
		ScimServiceProviderEndpoint: d.Get("scim_service_provider_endpoint").(string),
		ScimSharedSecretExists:      d.Get("scim_shared_secret_exists").(bool),
		ScimSharedSecret:            d.Get("scim_shared_secret").(string),
		SignSamlRequest:             d.Get("sign_saml_request").(string),
		SsoType:                     ListToStringSlice(d.Get("sso_type").([]interface{})), // Need to check format
		UseCustomSpMetadata:         d.Get("use_custom_sp_metadata").(bool),
		// AdminMetadata:               expandAdminMetaData(d),
		// UserMetadata:                expandUserMetaData(d),
	}
	return idpController
}

/*
func expandAdminMetaData(d *schema.ResourceData) *idpcontroller.AdminMetadata {
	adminInterface, ok := d.GetOk("admin_meta_data")
	if ok {
		adminList := adminInterface.([]interface{})
		if len(adminList) == 0 {
			return nil
		}
		metadata, _ := adminList[0].(map[string]interface{})
		return &idpcontroller.AdminMetadata{
			CertificateURL: metadata["certificate_url"].(string),
			SpEntityID:     metadata["sp_entity_id"].(string),
			SpMetadataURL:  metadata["sp_metadata_url"].(string),
			SpPostURL:      metadata["sp_post_url"].(string),
		}
	}
	return nil
}

func expandUserMetaData(d *schema.ResourceData) *idpcontroller.UserMetadata {
	userInterface, ok := d.GetOk("user_metadata")
	if ok {
		userList := userInterface.([]interface{})
		if len(userList) == 0 {
			return nil
		}
		metadata, _ := userList[0].(map[string]interface{})
		return &idpcontroller.UserMetadata{
			CertificateURL: metadata["certificate_url"].(string),
			SpBaseURL:      metadata["Sp_base_url"].(string),
			SpEntityID:     metadata["Sp_entity_id"].(string),
			SpMetadataURL:  metadata["sp_metadata_url"].(string),
			SpPostURL:      metadata["sp_post_url"].(string),
		}
	}
	return nil
}
*/
