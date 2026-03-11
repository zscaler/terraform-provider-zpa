package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_group"
)

func dataSourceTagGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTagGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"tag_key": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"tag_value": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microtenant_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTagGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var resp *tag_group.TagGroup
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for tag group %s\n", id)
		res, _, err := tag_group.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for tag group name %s\n", name)
		res, _, err := tag_group.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("tags", flattenTagGroupTagsDataSource(resp.Tags))
		_ = d.Set("microtenant_id", resp.MicroTenantID)
		_ = d.Set("microtenant_name", resp.MicroTenantName)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any tag group with name '%s' or id '%s'", name, id))
	}

	return nil
}

func flattenTagGroupTagsDataSource(tags []tag_group.Tag) []interface{} {
	if len(tags) == 0 {
		return nil
	}
	result := make([]interface{}, len(tags))
	for i, t := range tags {
		m := map[string]interface{}{
			"origin":    t.Origin,
			"namespace": flattenTagGroupNamespace(t.Namespace),
			"tag_key":   flattenTagGroupTagKey(t.TagKey),
			"tag_value": flattenTagGroupTagValue(t.TagValue),
		}
		result[i] = m
	}
	return result
}

func flattenTagGroupNamespace(ns *tag_group.TagNamespace) []interface{} {
	if ns == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id":      ns.ID,
			"name":    ns.Name,
			"enabled": ns.Enabled,
		},
	}
}

func flattenTagGroupTagKey(tk *tag_group.TagKey) []interface{} {
	if tk == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id":      tk.ID,
			"name":    tk.Name,
			"enabled": tk.Enabled,
		},
	}
}

func flattenTagGroupTagValue(tv *tag_group.TagValue) []interface{} {
	if tv == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id":   tv.ID,
			"name": tv.Name,
		},
	}
}
