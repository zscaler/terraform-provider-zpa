package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/workload_tag_group"
)

func dataSourceWorkloadTagGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkloadTagGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceWorkloadTagGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *common.CommonSummary
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for workload tag group %s\n", id)
		// Get all workload tag groups and find the one with matching ID
		allGroups, _, err := workload_tag_group.GetWorkloadTagGroup(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, group := range allGroups {
			if group.ID == id {
				resp = &group
				break
			}
		}
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for workload tag group name %s\n", name)
		// Use GetByName for direct name lookup
		res, _, err := workload_tag_group.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("enabled", resp.Enabled)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any workload tag group with name '%s' or id '%s'", name, id))
	}

	return nil
}
