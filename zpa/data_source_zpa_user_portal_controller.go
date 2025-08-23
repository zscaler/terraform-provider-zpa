package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
)

func dataSourceUserPortalController() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserPortalControllerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_name": {
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
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ext_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_domain_translation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"getc_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_data": {
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
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microtenant_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_notification": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_notification_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"managed_by_zs": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceUserPortalControllerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var resp *portal_controller.UserPortalController
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for user portal controller %s\n", id)
		res, _, err := portal_controller.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for user portal controller name %s\n", name)
		res, _, err := portal_controller.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("certificate_id", resp.CertificateId)
		_ = d.Set("certificate_name", resp.CertificateName)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("description", resp.Description)
		_ = d.Set("domain", resp.Domain)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("ext_domain", resp.ExtDomain)
		_ = d.Set("ext_domain_name", resp.ExtDomainName)
		_ = d.Set("ext_domain_translation", resp.ExtDomainTranslation)
		_ = d.Set("ext_label", resp.ExtLabel)
		_ = d.Set("getc_name", resp.GetcName)
		_ = d.Set("image_data", resp.ImageData)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("microtenant_id", resp.MicrotenantId)
		_ = d.Set("microtenant_name", resp.MicrotenantName)
		_ = d.Set("user_notification", resp.UserNotification)
		_ = d.Set("user_notification_enabled", resp.UserNotificationEnabled)
		_ = d.Set("managed_by_zs", resp.ManagedByZS)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any user portal controller with name '%s' or id '%s'", name, id))
	}

	return nil
}
