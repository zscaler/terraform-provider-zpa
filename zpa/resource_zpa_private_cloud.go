package zpa

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud"
)

func resourcePrivateCloud() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateCloudCreate,
		ReadContext:   resourcePrivateCloudRead,
		UpdateContext: resourcePrivateCloudUpdate,
		DeleteContext: resourcePrivateCloudDelete,
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
					resp, _, err := private_cloud.GetByName(ctx, service, id)
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
				Description: "Name of the Private Cloud Controller",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Private Cloud Controller",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether this Private Cloud Controller is enabled or not",
			},
			"re_enroll_period": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The re-enrollment period for the Private Cloud Controller",
			},
			"fire_drill_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether fire drill is enabled for the Private Cloud Controller",
			},
			"sitec_preferred": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the Site Controller is preferred",
			},
			"remote_lss": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether remote Log Streaming Service (LSS) is enabled",
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Microtenant ID for the Private Cloud Controller",
			},
			"assistant_groups_ids":      commonSummarySchema("The list of Assistant (App Connector) Group IDs associated with the Private Cloud Controller"),
			"site_controller_group_ids": commonSummarySchema("The list of Site Controller Group IDs associated with the Private Cloud Controller"),
			"siem_ids":                  commonSummarySchema("The list of SIEM IDs associated with the Private Cloud Controller"),
			"private_broker_group_ids":  commonSummarySchema("The list of Private Broker Group IDs associated with the Private Cloud Controller"),
			"zpn_fire_drill_site": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The fire drill site configuration for the Private Cloud Controller",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"microtenant_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"fire_drill_interval": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"fire_drill_interval_time_unit": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

// commonSummarySchema returns an optional/computed schema for an associated
// group reference. It follows the same convention used by zpa_server_group's
// app_connector_groups: a single block whose "id" is a set of string IDs. The
// Private Cloud POST/PUT API only accepts the ID for each associated group.
func commonSummarySchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: description,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func resourcePrivateCloudCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandPrivateCloud(d)

	// The zpn_fire_drill_site block is only valid when fire drill is enabled.
	if req.ZPNFireDrillSite != nil && !req.FireDrillEnabled {
		return diag.Errorf("zpn_fire_drill_site requires fire_drill_enabled to be set to true")
	}

	// The API rejects fire_drill_enabled=true on initial creation; it only
	// accepts it on a subsequent update. The zpn_fire_drill_site block is also
	// only valid once fire drill is enabled. Therefore both are stripped from
	// the create request and applied with a follow-up update below.
	fireDrillRequested := req.FireDrillEnabled
	fireDrillSite := req.ZPNFireDrillSite
	req.FireDrillEnabled = false
	req.ZPNFireDrillSite = nil

	log.Printf("[INFO] Creating zpa private cloud with request\n%+v\n", req)

	resp, _, err := private_cloud.Create(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created private cloud request. ID: %v\n", resp)
	d.SetId(resp.ID)

	// If fire_drill_enabled was requested, enable it (and apply the fire drill
	// site configuration) via a follow-up update, since the create API does not
	// accept these as set.
	if fireDrillRequested {
		log.Printf("[INFO] Enabling fire_drill_enabled for private cloud ID: %s via follow-up update", resp.ID)
		updateReq := req
		updateReq.ID = resp.ID
		updateReq.FireDrillEnabled = true
		updateReq.ZPNFireDrillSite = fireDrillSite
		if _, err := private_cloud.Update(ctx, service, resp.ID, &updateReq); err != nil {
			return diag.Errorf("private cloud created successfully (ID: %s), but enabling fire_drill_enabled failed: %s. You can retry by running terraform apply again.", resp.ID, err)
		}
	}

	return resourcePrivateCloudRead(ctx, d, meta)
}

func resourcePrivateCloudRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := private_cloud.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing private cloud %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting private cloud:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("re_enroll_period", resp.ReEnrollPeriod)
	_ = d.Set("fire_drill_enabled", resp.FireDrillEnabled)
	_ = d.Set("sitec_preferred", resp.SitecPreferred)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("remote_lss", resp.RemoteLss)
	_ = d.Set("assistant_groups_ids", flattenCommonSummaryIDs(resp.AssistantGroupsIDs))
	_ = d.Set("site_controller_group_ids", flattenCommonSummaryIDs(resp.SiteControllerGroupIDs))
	_ = d.Set("siem_ids", flattenCommonSummaryIDs(resp.SiemIDs))
	_ = d.Set("private_broker_group_ids", flattenCommonSummaryIDs(resp.PrivateBrokerGroupIDs))
	_ = d.Set("zpn_fire_drill_site", flattenZPNFireDrillSite(resp.ZPNFireDrillSite))
	return nil
}

func resourcePrivateCloudUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating private cloud ID: %v\n", id)
	req := expandPrivateCloud(d)

	// The zpn_fire_drill_site block is only valid when fire drill is enabled.
	if req.ZPNFireDrillSite != nil && !req.FireDrillEnabled {
		return diag.Errorf("zpn_fire_drill_site requires fire_drill_enabled to be set to true")
	}

	if _, _, err := private_cloud.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if _, err := private_cloud.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourcePrivateCloudRead(ctx, d, meta)
}

func resourcePrivateCloudDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting private cloud with id %v\n", d.Id())

	if _, err := private_cloud.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandPrivateCloud(d *schema.ResourceData) private_cloud.PrivateCloudController {
	return private_cloud.PrivateCloudController{
		ID:                     d.Get("id").(string),
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		Enabled:                d.Get("enabled").(bool),
		ReEnrollPeriod:         d.Get("re_enroll_period").(string),
		FireDrillEnabled:       d.Get("fire_drill_enabled").(bool),
		SitecPreferred:         d.Get("sitec_preferred").(bool),
		RemoteLss:              d.Get("remote_lss").(bool),
		AssistantGroupsIDs:     expandCommonSummaryIDs(d, "assistant_groups_ids"),
		SiteControllerGroupIDs: expandCommonSummaryIDs(d, "site_controller_group_ids"),
		SiemIDs:                expandCommonSummaryIDs(d, "siem_ids"),
		PrivateBrokerGroupIDs:  expandCommonSummaryIDs(d, "private_broker_group_ids"),
		ZPNFireDrillSite:       expandZPNFireDrillSite(d.Get("zpn_fire_drill_site").([]interface{})),
	}
}

// expandCommonSummaryIDs builds a list of common summary objects from the
// resource configuration. It follows the same single-block + set-of-string-IDs
// convention used by zpa_server_group's app_connector_groups. The Private Cloud
// POST/PUT API only accepts the ID for each associated group, so only the id is
// sent.
func expandCommonSummaryIDs(d *schema.ResourceData, key string) []common.CommonSummary {
	raw, ok := d.GetOk(key)
	if !ok || raw == nil {
		return nil
	}

	blocks := raw.([]interface{})
	if len(blocks) == 0 || blocks[0] == nil {
		return nil
	}

	block, ok := blocks[0].(map[string]interface{})
	if !ok {
		return nil
	}

	idRaw, ok := block["id"]
	if !ok || idRaw == nil {
		return nil
	}

	idSet, ok := idRaw.(*schema.Set)
	if !ok || idSet.Len() == 0 {
		return nil
	}

	var out []common.CommonSummary
	for _, id := range idSet.List() {
		out = append(out, common.CommonSummary{
			ID: id.(string),
		})
	}
	return out
}

// flattenCommonSummaryIDs writes back a single block whose "id" is the set of
// associated group IDs, matching the schema used by the resource.
func flattenCommonSummaryIDs(list []common.CommonSummary) []interface{} {
	if len(list) == 0 {
		return []interface{}{}
	}
	ids := make([]interface{}, len(list))
	for i, item := range list {
		ids[i] = item.ID
	}
	return []interface{}{
		map[string]interface{}{
			"id": schema.NewSet(schema.HashString, ids),
		},
	}
}

func expandZPNFireDrillSite(raw []interface{}) *private_cloud.ZPNFireDrillSite {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return nil
	}
	return &private_cloud.ZPNFireDrillSite{
		ID:                        m["id"].(string),
		MicrotenantID:             m["microtenant_id"].(string),
		FireDrillInterval:         m["fire_drill_interval"].(string),
		FireDrillIntervalTimeUnit: m["fire_drill_interval_time_unit"].(string),
	}
}

func flattenZPNFireDrillSite(site *private_cloud.ZPNFireDrillSite) []interface{} {
	if site == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"id":                            site.ID,
			"microtenant_id":                site.MicrotenantID,
			"fire_drill_interval":           site.FireDrillInterval,
			"fire_drill_interval_time_unit": site.FireDrillIntervalTimeUnit,
		},
	}
}
