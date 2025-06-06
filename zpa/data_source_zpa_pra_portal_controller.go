package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praportal"
)

func dataSourcePRAPortalController() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePRAPortalControllerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the privileged portal",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the privileged portal",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the privileged portal",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not the privileged portal is enabled",
			},
			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The canonical name (CNAME DNS records) associated with the privileged portal",
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The domain of the privileged portal",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the certificate",
			},
			"certificate_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the certificate",
			},
			"user_notification": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The notification message displayed in the banner of the privileged portallink, if enabled",
			},
			"user_notification_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the Notification Banner is enabled (true) or disabled (false)",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the privileged portal is created",
			},
			"modified_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the tenant who modified the privileged portal",
			},
			"modified_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the privileged portal is modified",
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.",
			},
			"microtenant_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Microtenant",
			},
			"ext_label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The domain prefix for the privileged portal URL. The supported string can include numbers, lower case characters, and only supports a hyphen (-).",
			},
			"ext_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The external domain name prefix of the Browser Access application that is used for Zscaler-managed certificates when creating a privileged portal.",
			},
			"ext_domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The domain suffix for the privileged portal URL. This field must be one of the customer's authentication domains.",
			},
			"ext_domain_translation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The translation of the external domain name prefix of the Browser Access application that is used for Zscaler-managed certificates when creating a privileged portal.",
			},
			"user_portal_gid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the user portal.",
			},
			"user_portal_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the user portal.",
			},
			"getc_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the user portal.",
			},
		},
	}
}

func dataSourcePRAPortalControllerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var resp *praportal.PRAPortal
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for pra portal controller %s\n", id)
		res, _, err := praportal.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err) // Wrap error using diag.FromErr
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data for pra portal controller name %s\n", name)
		res, _, err := praportal.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err) // Wrap error using diag.FromErr
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("cname", resp.CName)
		_ = d.Set("domain", resp.Domain)
		_ = d.Set("certificate_id", resp.CertificateID)
		_ = d.Set("certificate_name", resp.CertificateName)
		_ = d.Set("user_notification", resp.UserNotification)
		_ = d.Set("user_notification_enabled", resp.UserNotificationEnabled)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("microtenant_id", resp.MicroTenantID)
		_ = d.Set("microtenant_name", resp.MicroTenantName)
		_ = d.Set("ext_label", resp.ExtLabel)
		_ = d.Set("ext_domain", resp.ExtDomain)
		_ = d.Set("ext_domain_name", resp.ExtDomainName)
		_ = d.Set("ext_domain_translation", resp.ExtDomainTranslation)
		_ = d.Set("user_portal_gid", resp.UserPortalGid)
		_ = d.Set("getc_name", resp.GetcName)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any pra portal controller with name '%s' or id '%s'", name, id))
	}

	return nil
}
