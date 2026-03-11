package zpa

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_group"
)

func resourceTagGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTagGroupCreate,
		ReadContext:   resourceTagGroupRead,
		UpdateContext: resourceTagGroupUpdate,
		DeleteContext: resourceTagGroupDelete,
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
					resp, _, err := tag_group.GetByName(ctx, service, id)
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
				Description: "Name of the tag group",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the tag group",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of tag value IDs associated with this tag group",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Microtenant ID",
			},
		},
	}
}

func resourceTagGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandTagGroup(d)
	log.Printf("[INFO] Creating zpa tag group with request\n%+v\n", req)

	resp, _, err := tag_group.Create(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zpa tag group. ID: %v\n", resp.ID)
	d.SetId(resp.ID)

	return resourceTagGroupRead(ctx, d, meta)
}

func resourceTagGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := tag_group.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing zpa_tag_group %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zpa tag group:\n%+v\n", resp)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("tags", flattenTagGroupTags(resp.Tags))
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	return nil
}

func resourceTagGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating zpa tag group ID: %v\n", id)
	req := expandTagGroup(d)

	if _, _, err := tag_group.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if _, err := tag_group.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceTagGroupRead(ctx, d, meta)
}

func resourceTagGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting zpa tag group with id %v\n", d.Id())

	if _, err := tag_group.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandTagGroup(d *schema.ResourceData) tag_group.TagGroup {
	tg := tag_group.TagGroup{
		ID:          d.Get("id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if v, ok := d.GetOk("tags"); ok {
		tg.Tags = expandTagGroupTags(v.(*schema.Set).List())
	} else {
		tg.Tags = []tag_group.Tag{}
	}

	return tg
}

func expandTagGroupTags(tagValueIDs []interface{}) []tag_group.Tag {
	result := make([]tag_group.Tag, 0, len(tagValueIDs))
	for _, id := range tagValueIDs {
		result = append(result, tag_group.Tag{
			TagValue: &tag_group.TagValue{
				ID: id.(string),
			},
		})
	}
	return result
}

func flattenTagGroupTags(tags []tag_group.Tag) []interface{} {
	if len(tags) == 0 {
		return nil
	}
	result := make([]interface{}, 0, len(tags))
	for _, t := range tags {
		if t.TagValue != nil {
			result = append(result, t.TagValue.ID)
		}
	}
	return result
}
