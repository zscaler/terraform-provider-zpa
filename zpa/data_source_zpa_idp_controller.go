package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/idpcontroller"
)

func dataSourceIdpController() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIdpControllerRead,
		Schema: map[string]*schema.Schema{
			"admin_metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sp_base_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sp_entity_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sp_metadata_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sp_post_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"admin_sp_signing_cert_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_provision": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disable_saml_based_policy": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"domain_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_scim_based_policy": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Default value if null is True",
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"idp_entity_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"login_name_attribute": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"login_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"login_hint": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"force_auth": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enable_arbitrary_auth_domains": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifiedby": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reauth_on_user_update": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"redirect_binding": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"scim_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"scim_service_provider_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scim_shared_secret_exists": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sign_saml_request": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sso_type": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"use_custom_sp_metadata": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"user_metadata": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sp_base_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sp_entity_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sp_metadata_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sp_post_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"user_sp_signing_cert_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIdpControllerRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.IDPController

	var resp *idpcontroller.IdpController
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for idp controller %s\n", id)
		res, _, err := idpcontroller.Get(service, id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for idp controller name %s\n", name)
		res, _, err := idpcontroller.GetByName(service, name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("admin_sp_signing_cert_id", resp.AdminSpSigningCertID)
		_ = d.Set("auto_provision", resp.AutoProvision)
		_ = d.Set("description", resp.Description)
		_ = d.Set("disable_saml_based_policy", resp.DisableSamlBasedPolicy)
		_ = d.Set("domain_list", resp.Domainlist)
		_ = d.Set("enable_scim_based_policy", resp.EnableScimBasedPolicy)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("idp_entity_id", resp.IdpEntityID)
		_ = d.Set("login_name_attribute", resp.LoginNameAttribute)
		_ = d.Set("login_url", resp.LoginURL)
		_ = d.Set("login_hint", resp.LoginHint)
		_ = d.Set("force_auth", resp.ForceAuth)
		_ = d.Set("enable_arbitrary_auth_domains", resp.EnableArbitraryAuthDomains)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("name", resp.Name)
		_ = d.Set("reauth_on_user_update", resp.ReauthOnUserUpdate)
		_ = d.Set("redirect_binding", resp.RedirectBinding)
		_ = d.Set("scim_enabled", resp.ScimEnabled)
		_ = d.Set("scim_service_provider_endpoint", resp.ScimServiceProviderEndpoint)
		_ = d.Set("scim_shared_secret_exists", resp.ScimSharedSecretExists)
		_ = d.Set("sign_saml_request", resp.SignSamlRequest)
		_ = d.Set("sso_type", resp.SsoType)
		_ = d.Set("use_custom_sp_metadata", resp.UseCustomSpMetadata)
		_ = d.Set("user_sp_signing_cert_id", resp.UserSpSigningCertID)
		if resp.UserMetadata != nil {
			_ = d.Set("user_metadata", flattenUserMeta(resp.UserMetadata))
		}
		if resp.AdminMetadata != nil {
			_ = d.Set("admin_metadata", flattenAdminMeta(resp.AdminMetadata))
		}

		// Set epoch attributes explicitly
		creationTime, err := epochToRFC1123(resp.CreationTime, false)
		if err != nil {
			return fmt.Errorf("error formatting creation_time: %s", err)
		}
		if err := d.Set("creation_time", creationTime); err != nil {
			return fmt.Errorf("error setting creation_time: %s", err)
		}

		modifiedTime, err := epochToRFC1123(resp.ModifiedTime, false)
		if err != nil {
			return fmt.Errorf("error formatting modified_time: %s", err)
		}
		if err := d.Set("modified_time", modifiedTime); err != nil {
			return fmt.Errorf("error setting modified_time: %s", err)
		}
	} else {
		return fmt.Errorf("couldn't find any idp controller with name '%s' or id '%s'", name, id)
	}
	return nil
}

func flattenAdminMeta(metaData *idpcontroller.AdminMetadata) []map[string]interface{} {
	result := make([]map[string]interface{}, 1)
	result[0] = make(map[string]interface{})
	result[0]["certificate_url"] = metaData.CertificateURL
	result[0]["sp_base_url"] = metaData.SpBaseURL
	result[0]["sp_entity_id"] = metaData.SpEntityID
	result[0]["sp_metadata_url"] = metaData.SpMetadataURL
	result[0]["sp_post_url"] = metaData.SpPostURL
	return result
}

func flattenUserMeta(metaData *idpcontroller.UserMetadata) []map[string]interface{} {
	result := make([]map[string]interface{}, 1)
	result[0] = make(map[string]interface{})
	result[0]["certificate_url"] = metaData.CertificateURL
	result[0]["sp_base_url"] = metaData.SpBaseURL
	result[0]["sp_metadata_url"] = metaData.SpMetadataURL
	result[0]["sp_post_url"] = metaData.SpPostURL
	return result
}
