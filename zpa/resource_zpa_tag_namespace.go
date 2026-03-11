package zpa

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_namespace"
)

func resourceTagNamespace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTagNamespaceCreate,
		ReadContext:   resourceTagNamespaceRead,
		UpdateContext: resourceTagNamespaceUpdate,
		DeleteContext: resourceTagNamespaceDelete,
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
					_ = d.Set("id", id)
				} else {
					resp, _, err := tag_namespace.GetByName(ctx, service, id)
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
				Description: "Name of the tag namespace",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the tag namespace",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether this tag namespace is enabled",
			},
			"origin": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "CUSTOM",
				Description: "Origin of the tag namespace. Defaults to CUSTOM",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "STATIC",
				Description: "Type of the tag namespace. Defaults to STATIC",
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Microtenant ID",
			},
		},
	}
}

func resourceTagNamespaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandTagNamespace(d)
	log.Printf("[INFO] Creating zpa tag namespace with request\n%+v\n", req)

	resp, _, err := tag_namespace.Create(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zpa tag namespace. ID: %v\n", resp.ID)
	d.SetId(resp.ID)

	return resourceTagNamespaceRead(ctx, d, meta)
}

func resourceTagNamespaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := tag_namespace.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing zpa_tag_namespace %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zpa tag namespace:\n%+v\n", resp)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("origin", resp.Origin)
	_ = d.Set("type", resp.Type)
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	return nil
}

func resourceTagNamespaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating zpa tag namespace ID: %v\n", id)
	req := expandTagNamespace(d)

	if _, _, err := tag_namespace.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if _, err := tag_namespace.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceTagNamespaceRead(ctx, d, meta)
}

func resourceTagNamespaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting zpa tag namespace with id %v\n", d.Id())

	if _, err := tag_namespace.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandTagNamespace(d *schema.ResourceData) tag_namespace.Namespace {
	return tag_namespace.Namespace{
		ID:            d.Get("id").(string),
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Enabled:       d.Get("enabled").(bool),
		Origin:        d.Get("origin").(string),
		Type:          d.Get("type").(string),
		MicroTenantID: d.Get("microtenant_id").(string),
	}
}
