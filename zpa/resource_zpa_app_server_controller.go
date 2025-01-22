package zpa

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"
)

func resourceApplicationServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationServerCreate,
		ReadContext:   resourceApplicationServerRead,
		UpdateContext: resourceApplicationServerUpdate,
		DeleteContext: resourceApplicationServerDelete,
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
					resp, _, err := appservercontroller.GetByName(ctx, service, id)
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "This field defines the name of the server.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "This field defines the description of the server.",
			},
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This field defines the domain or IP address of the server.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "This field defines the status of the server.",
			},
			// App Server Group ID can only be attached if Dynamic Server Discovery in Server Group is False
			"app_server_group_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "This field defines the list of server groups IDs.",
				Optional:    true,
			},
			"config_space": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{
					"DEFAULT",
					"SIEM",
				}, false),
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceApplicationServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandAppServerRequest(d)
	log.Printf("[INFO] Creating zpa application server with request\n%+v\n", req)

	resp, _, err := appservercontroller.Create(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created application server request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceApplicationServerRead(ctx, d, meta)
}

func resourceApplicationServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := appservercontroller.Get(ctx, service, d.Id())
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing application server %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting application server:\n%+v\n", resp)
	_ = d.Set("address", resp.Address)
	_ = d.Set("app_server_group_ids", resp.AppServerGroupIds)
	_ = d.Set("config_space", resp.ConfigSpace)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("name", resp.Name)
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	return nil
}

func resourceApplicationServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating application server ID: %v\n", id)
	req := expandAppServerRequest(d)

	if _, _, err := appservercontroller.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := appservercontroller.Update(ctx, service, id, req); err != nil {
		return diag.FromErr(err)
	}

	return resourceApplicationServerRead(ctx, d, meta)
}

func resourceApplicationServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Call Delete with context and necessary parameters
	if _, err := appservercontroller.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[INFO] application server deleted successfully")
	return nil
}

func expandAppServerRequest(d *schema.ResourceData) appservercontroller.ApplicationServer {
	applicationServer := appservercontroller.ApplicationServer{
		ID:                d.Id(),
		Address:           d.Get("address").(string),
		ConfigSpace:       d.Get("config_space").(string),
		AppServerGroupIds: SetToStringSlice(d.Get("app_server_group_ids").(*schema.Set)),
		Description:       d.Get("description").(string),
		Enabled:           d.Get("enabled").(bool),
		Name:              d.Get("name").(string),
		MicroTenantID:     d.Get("microtenant_id").(string),
	}
	return applicationServer
}
