package zpa

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
)

func resourceUserPortalController() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserPortalControllerCreate,
		ReadContext:   resourceUserPortalControllerRead,
		UpdateContext: resourceUserPortalControllerUpdate,
		DeleteContext: resourceUserPortalControllerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				client := meta.(*Client)
				service := client.Service

				microTenantID := GetString(d.Get("microtenant_id"))
				if microTenantID != "" {
					service = service.WithMicroTenant(microTenantID)
				}

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := portal_controller.GetByName(ctx, service, id)
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the User Portal Controller",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Certificate ID for the User Portal Controller",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the User Portal Controller",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Domain for the User Portal Controller",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether this User Portal Controller is enabled or not",
			},
			"ext_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "External domain for the User Portal Controller",
			},
			"ext_domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "External domain name for the User Portal Controller",
			},
			"ext_domain_translation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "External domain translation for the User Portal Controller",
			},
			"ext_label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "External label for the User Portal Controller",
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Microtenant ID for the User Portal Controller",
			},
			"user_notification": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User notification message for the User Portal Controller",
			},
			"user_notification_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether user notifications are enabled for the User Portal Controller",
			},
		},
	}
}

func resourceUserPortalControllerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandUserPortalController(d)
	log.Printf("[INFO] Creating zpa user portal controller with request\n%+v\n", req)

	resp, _, err := portal_controller.Create(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created user portal controller request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceUserPortalControllerRead(ctx, d, meta)
}

func resourceUserPortalControllerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := portal_controller.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing user portal controller %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting user portal controller:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("certificate_id", resp.CertificateId)
	_ = d.Set("description", resp.Description)
	_ = d.Set("domain", resp.Domain)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("ext_domain", resp.ExtDomain)
	_ = d.Set("ext_domain_name", resp.ExtDomainName)
	_ = d.Set("ext_domain_translation", resp.ExtDomainTranslation)
	_ = d.Set("ext_label", resp.ExtLabel)
	_ = d.Set("microtenant_id", resp.MicrotenantId)
	_ = d.Set("user_notification", resp.UserNotification)
	_ = d.Set("user_notification_enabled", resp.UserNotificationEnabled)
	return nil
}

func resourceUserPortalControllerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating user portal controller ID: %v\n", id)
	req := expandUserPortalController(d)

	if _, _, err := portal_controller.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if _, err := portal_controller.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceUserPortalControllerRead(ctx, d, meta)
}

func resourceUserPortalControllerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting user portal controller with id %v\n", d.Id())

	if _, err := portal_controller.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandUserPortalController(d *schema.ResourceData) portal_controller.UserPortalController {
	return portal_controller.UserPortalController{
		ID:                      d.Get("id").(string),
		Name:                    d.Get("name").(string),
		CertificateId:           d.Get("certificate_id").(string),
		Description:             d.Get("description").(string),
		Domain:                  d.Get("domain").(string),
		Enabled:                 d.Get("enabled").(bool),
		ExtDomain:               d.Get("ext_domain").(string),
		ExtDomainName:           d.Get("ext_domain_name").(string),
		ExtDomainTranslation:    d.Get("ext_domain_translation").(string),
		ExtLabel:                d.Get("ext_label").(string),
		MicrotenantId:           d.Get("microtenant_id").(string),
		UserNotification:        d.Get("user_notification").(string),
		UserNotificationEnabled: d.Get("user_notification_enabled").(bool),
	}
}
