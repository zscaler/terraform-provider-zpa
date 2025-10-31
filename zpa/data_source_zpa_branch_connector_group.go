package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/branch_connector_group"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

func dataSourceBranchConnectorGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBranchConnectorGroupRead,
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

func dataSourceBranchConnectorGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *common.CommonSummary
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for branch connector group %s\n", id)
		// Get all branch connector groups and find the one with matching ID
		allGroups, _, err := branch_connector_group.GetBranchConnectorGroupSummary(ctx, service)
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
		log.Printf("[INFO] Getting data for branch connector group name %s\n", name)
		// Get all branch connector groups and find the one with matching name
		allGroups, _, err := branch_connector_group.GetBranchConnectorGroupSummary(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, group := range allGroups {
			if group.Name == name {
				resp = &group
				break
			}
		}
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("enabled", resp.Enabled)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any branch connector group with name '%s' or id '%s'", name, id))
	}

	return nil
}
