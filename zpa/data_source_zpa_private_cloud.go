package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud"
)

func dataSourcePrivateCloud() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateCloudRead,
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
			"re_enroll_period": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fire_drill_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sitec_preferred": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"remote_lss": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"zscaler_managed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"microtenant_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assistant_groups_ids":       dataCommonSummarySchema(),
			"site_controller_group_ids":  dataCommonSummarySchema(),
			"siem_ids":                   dataCommonSummarySchema(),
			"private_exporter_group_ids": dataCommonSummarySchema(),
			"private_broker_group_ids":   dataCommonSummarySchema(),
			"zpn_fire_drill_site": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"microtenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"microtenant_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fire_drill_interval": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fire_drill_interval_time_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// dataCommonSummarySchema returns a computed list schema for a list of
// common summary objects (id/name/enabled).
func dataCommonSummarySchema() *schema.Schema {
	return &schema.Schema{
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
	}
}

func dataSourcePrivateCloudRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var resp *private_cloud.PrivateCloudController
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for private cloud %s\n", id)
		res, _, err := private_cloud.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data for private cloud name %s\n", name)
		res, _, err := private_cloud.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("re_enroll_period", resp.ReEnrollPeriod)
		_ = d.Set("fire_drill_enabled", resp.FireDrillEnabled)
		_ = d.Set("sitec_preferred", resp.SitecPreferred)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("remote_lss", resp.RemoteLss)
		_ = d.Set("read_only", resp.ReadOnly)
		_ = d.Set("zscaler_managed", resp.ZscalerManaged)
		_ = d.Set("microtenant_name", resp.MicrotenantName)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("assistant_groups_ids", flattenCommonSummaryData(resp.AssistantGroupsIDs))
		_ = d.Set("site_controller_group_ids", flattenCommonSummaryData(resp.SiteControllerGroupIDs))
		_ = d.Set("siem_ids", flattenCommonSummaryData(resp.SiemIDs))
		_ = d.Set("private_exporter_group_ids", flattenCommonSummaryData(resp.PrivateExporterGroupIDs))
		_ = d.Set("private_broker_group_ids", flattenCommonSummaryData(resp.PrivateBrokerGroupIDs))
		_ = d.Set("zpn_fire_drill_site", flattenZPNFireDrillSiteData(resp.ZPNFireDrillSite))
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any private cloud with name '%s' or id '%s'", name, id))
	}

	return nil
}

// flattenCommonSummaryData flattens a list of common summary objects with all
// of their attributes (id/name/enabled). The data source exposes the full set
// of fields since they are all returned in the API response.
func flattenCommonSummaryData(list []common.CommonSummary) []interface{} {
	if len(list) == 0 {
		return []interface{}{}
	}
	out := make([]interface{}, len(list))
	for i, item := range list {
		out[i] = map[string]interface{}{
			"id":      item.ID,
			"name":    item.Name,
			"enabled": item.Enabled,
		}
	}
	return out
}

// flattenZPNFireDrillSiteData flattens the fire drill site object with all of
// its attributes. The data source exposes the full set of fields since they
// are all returned in the API response.
func flattenZPNFireDrillSiteData(site *private_cloud.ZPNFireDrillSite) []interface{} {
	if site == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"id":                            site.ID,
			"microtenant_id":                site.MicrotenantID,
			"microtenant_name":              site.MicrotenantName,
			"fire_drill_interval":           site.FireDrillInterval,
			"fire_drill_interval_time_unit": site.FireDrillIntervalTimeUnit,
			"creation_time":                 site.CreationTime,
			"modified_by":                   site.ModifiedBy,
			"modified_time":                 site.ModifiedTime,
		},
	}
}
