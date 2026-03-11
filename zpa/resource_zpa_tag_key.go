package zpa

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	tag_key_controller "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_key"
)

func resourceTagKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTagKeyCreate,
		ReadContext:   resourceTagKeyRead,
		UpdateContext: resourceTagKeyUpdate,
		DeleteContext: resourceTagKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				client := meta.(*Client)
				service := client.Service

				microTenantID := GetString(d.Get("microtenant_id"))
				if microTenantID != "" {
					service = service.WithMicroTenant(microTenantID)
				}

				// Import format: namespace_id/tag_key_id_or_name
				parts := strings.SplitN(d.Id(), "/", 2)
				if len(parts) != 2 {
					return nil, fmt.Errorf("unexpected import format: expected namespace_id/tag_key_id_or_name, got %s", d.Id())
				}
				namespaceID := parts[0]
				tagKeyIdentifier := parts[1]

				_ = d.Set("namespace_id", namespaceID)

				resp, _, err := tag_key_controller.Get(ctx, service, namespaceID, tagKeyIdentifier)
				if err == nil {
					d.SetId(resp.ID)
					_ = d.Set("id", resp.ID)
				} else {
					resp, _, err = tag_key_controller.GetByName(ctx, service, namespaceID, tagKeyIdentifier)
					if err == nil {
						d.SetId(resp.ID)
						_ = d.Set("id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, fmt.Errorf("could not find tag key with identifier '%s' in namespace '%s': %v", tagKeyIdentifier, namespaceID, err)
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
				Description: "Name of the tag key",
			},
			"namespace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the tag namespace this key belongs to",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the tag key",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether this tag key is enabled",
			},
			"origin": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "CUSTOM",
				Description: "Origin of the tag key. Defaults to CUSTOM",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "STATIC",
				Description: "Type of the tag key. Defaults to STATIC",
			},
			"tag_values": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tag values associated with this tag key",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the tag value",
						},
					},
				},
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Microtenant ID",
			},
		},
	}
}

func resourceTagKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	namespaceID := d.Get("namespace_id").(string)
	req := expandTagKey(d)
	log.Printf("[INFO] Creating zpa tag key with request\n%+v\n", req)

	resp, _, err := tag_key_controller.Create(ctx, service, namespaceID, req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zpa tag key. ID: %v\n", resp.ID)
	d.SetId(resp.ID)

	return resourceTagKeyRead(ctx, d, meta)
}

func resourceTagKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	namespaceID := d.Get("namespace_id").(string)
	resp, _, err := tag_key_controller.Get(ctx, service, namespaceID, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing zpa_tag_key %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zpa tag key:\n%+v\n", resp)
	_ = d.Set("name", resp.Name)
	_ = d.Set("namespace_id", resp.NamespaceID)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("origin", resp.Origin)
	_ = d.Set("type", resp.Type)
	_ = d.Set("tag_values", flattenTagValues(resp.TagValues))
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	return nil
}

func resourceTagKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)
	log.Printf("[INFO] Updating zpa tag key ID: %v\n", id)
	req := expandTagKey(d)

	if _, _, err := tag_key_controller.Get(ctx, service, namespaceID, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if _, err := tag_key_controller.Update(ctx, service, namespaceID, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceTagKeyRead(ctx, d, meta)
}

func resourceTagKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	namespaceID := d.Get("namespace_id").(string)
	log.Printf("[INFO] Deleting zpa tag key with id %v\n", d.Id())

	if _, err := tag_key_controller.Delete(ctx, service, namespaceID, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandTagKey(d *schema.ResourceData) tag_key_controller.TagKey {
	tagKey := tag_key_controller.TagKey{
		ID:          d.Get("id").(string),
		Name:        d.Get("name").(string),
		NamespaceID: d.Get("namespace_id").(string),
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(bool),
		Origin:      d.Get("origin").(string),
		Type:        d.Get("type").(string),
	}

	if v, ok := d.GetOk("tag_values"); ok {
		tagKey.TagValues = expandTagKeyTagValues(v.([]interface{}))
	} else {
		tagKey.TagValues = []tag_key_controller.TagValue{}
	}

	return tagKey
}

func expandTagKeyTagValues(values []interface{}) []tag_key_controller.TagValue {
	result := make([]tag_key_controller.TagValue, 0, len(values))
	for _, val := range values {
		item := val.(map[string]interface{})
		entry := tag_key_controller.TagValue{
			Name: item["name"].(string),
		}
		if v, ok := item["id"]; ok && v.(string) != "" {
			entry.ID = v.(string)
		}
		result = append(result, entry)
	}
	return result
}

func flattenTagValues(tagValues []tag_key_controller.TagValue) []interface{} {
	if len(tagValues) == 0 {
		return nil
	}
	result := make([]interface{}, len(tagValues))
	for i, tv := range tagValues {
		result[i] = map[string]interface{}{
			"id":   tv.ID,
			"name": tv.Name,
		}
	}
	return result
}
