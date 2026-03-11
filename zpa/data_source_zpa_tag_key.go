package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tag_key_controller "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_key"
)

func dataSourceTagKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTagKeyRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the tag namespace this key belongs to",
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"customer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_values": {
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

func dataSourceTagKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	namespaceID := d.Get("namespace_id").(string)

	var resp *tag_key_controller.TagKey
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for tag key %s in namespace %s\n", id, namespaceID)
		res, _, err := tag_key_controller.Get(ctx, service, namespaceID, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for tag key name %s in namespace %s\n", name, namespaceID)
		res, _, err := tag_key_controller.GetByName(ctx, service, namespaceID, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("namespace_id", resp.NamespaceID)
		_ = d.Set("description", resp.Description)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("customer_id", resp.CustomerID)
		_ = d.Set("origin", resp.Origin)
		_ = d.Set("type", resp.Type)
		_ = d.Set("tag_values", flattenTagValues(resp.TagValues))
		_ = d.Set("microtenant_id", resp.MicroTenantID)
		_ = d.Set("microtenant_name", resp.MicroTenantName)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any tag key with name '%s' or id '%s' in namespace '%s'", name, id, namespaceID))
	}

	return nil
}
