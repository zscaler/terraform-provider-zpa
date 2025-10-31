package zpa

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/aup"
)

func resourceUserPortalAUP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserPortalAUPCreate,
		ReadContext:   resourceUserPortalAUPRead,
		UpdateContext: resourceUserPortalAUPUpdate,
		DeleteContext: resourceUserPortalAUPDelete,
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
					resp, _, err := aup.GetByName(ctx, service, id)
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
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"aup": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"phone_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceUserPortalAUPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandUserPortalAUP(d)
	log.Printf("[INFO] Creating zpa user portal aup with request\n%+v\n", req)

	resp, _, err := aup.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created user portal aup request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceUserPortalAUPRead(ctx, d, meta)
}

func resourceUserPortalAUPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := aup.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing user portal aup %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting user portal aup:\n%+v\n", resp)
	d.SetId(resp.ID)
	d.Set("name", resp.Name)
	d.Set("description", resp.Description)
	d.Set("enabled", resp.Enabled)
	d.Set("aup", resp.Aup)
	d.Set("email", resp.Email)
	d.Set("phone_num", resp.PhoneNum)
	d.Set("microtenant_id", resp.MicrotenantID)
	return nil
}

func resourceUserPortalAUPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating user portal aup ID: %v\n", id)
	req := expandUserPortalAUP(d)

	if _, _, err := aup.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if _, err := aup.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceUserPortalAUPRead(ctx, d, meta)
}

func resourceUserPortalAUPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting user portal aup with id %v\n", d.Id())

	if _, err := aup.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandUserPortalAUP(d *schema.ResourceData) aup.UserPortalAup {
	return aup.UserPortalAup{
		ID:            d.Get("id").(string),
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Enabled:       d.Get("enabled").(bool),
		Aup:           d.Get("aup").(string),
		Email:         d.Get("email").(string),
		PhoneNum:      d.Get("phone_num").(string),
		MicrotenantID: d.Get("microtenant_id").(string),
	}
}
