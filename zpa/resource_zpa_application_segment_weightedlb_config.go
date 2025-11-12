package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
)

func resourceApplicationSegmentWeightedLBConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationSegmentWeightedLBConfigCreate,
		ReadContext:   resourceApplicationSegmentWeightedLBConfigRead,
		UpdateContext: resourceApplicationSegmentWeightedLBConfigUpdate,
		DeleteContext: resourceApplicationSegmentWeightedLBConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"application_id", "application_name"},
				Description:  "Application segment identifier to manage. Either application_id or application_name must be provided.",
			},
			"application_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"application_id", "application_name"},
				Description:  "Application segment name to manage. Either application_id or application_name must be provided.",
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional microtenant identifier.",
			},
			"weighted_load_balancing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable or disable weighted load balancing for the application segment.",
			},
			"application_to_server_group_mappings": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Application to server group mapping details and weights.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Server group mapping identifier.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Server group name.",
						},
						"passive": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether the server group is passive.",
						},
						"weight": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Assigned weight for the server group.",
						},
					},
				},
			},
		},
	}
}

func resourceApplicationSegmentWeightedLBConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	applicationID, applicationName, diags := resolveApplicationSegmentIdentity(ctx, service, d)
	if diags.HasError() {
		return diags
	}

	if diags := updateWeightedLBConfig(ctx, service, applicationID, d); diags.HasError() {
		return diags
	}

	d.SetId(applicationID)
	if applicationName != "" {
		_ = d.Set("application_name", applicationName)
	}

	return resourceApplicationSegmentWeightedLBConfigRead(ctx, d, meta)
}

func resourceApplicationSegmentWeightedLBConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	applicationID := d.Id()
	if applicationID == "" {
		applicationID = GetString(d.Get("application_id"))
	}
	if applicationID == "" {
		return diag.FromErr(fmt.Errorf("application_id is not set in state"))
	}

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	config, _, err := applicationsegment.GetWeightedLoadBalancerConfig(ctx, service, applicationID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve weighted load balancer config for application %s: %w", applicationID, err))
	}
	if config == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] Retrieved weighted load balancer config for application %s", applicationID)

	if err := d.Set("application_id", applicationID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("weighted_load_balancing", config.WeightedLoadBalancing); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("application_to_server_group_mappings", flattenApplicationToServerGroupMappings(config.ApplicationToServerGroupMaps)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set application_to_server_group_mappings: %w", err))
	}

	if GetString(d.Get("application_name")) == "" {
		if app, _, err := applicationsegment.Get(ctx, service, applicationID); err == nil && app != nil {
			_ = d.Set("application_name", app.Name)
		}
	}

	d.SetId(applicationID)
	return nil
}

func resourceApplicationSegmentWeightedLBConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	applicationID := d.Id()
	if applicationID == "" {
		applicationID = GetString(d.Get("application_id"))
	}
	if applicationID == "" {
		return diag.FromErr(fmt.Errorf("application_id must be set to update weighted load balancer config"))
	}

	if diags := updateWeightedLBConfig(ctx, service, applicationID, d); diags.HasError() {
		return diags
	}

	return resourceApplicationSegmentWeightedLBConfigRead(ctx, d, meta)
}

func resourceApplicationSegmentWeightedLBConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	applicationID := d.Id()
	if applicationID == "" {
		applicationID = GetString(d.Get("application_id"))
	}
	if applicationID == "" {
		return nil
	}

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	payload := applicationsegment.WeightedLoadBalancerConfig{
		ApplicationID:         applicationID,
		WeightedLoadBalancing: false,
	}

	log.Printf("[INFO] Disabling weighted load balancer config for application %s", applicationID)
	if _, _, err := applicationsegment.UpdateWeightedLoadBalancerConfig(ctx, service, applicationID, payload); err != nil {
		return diag.FromErr(fmt.Errorf("failed to disable weighted load balancer config for application %s: %w", applicationID, err))
	}

	d.SetId("")
	return nil
}

func resolveApplicationSegmentIdentity(ctx context.Context, service *zscaler.Service, d *schema.ResourceData) (string, string, diag.Diagnostics) {
	applicationID := GetString(d.Get("application_id"))
	applicationName := GetString(d.Get("application_name"))

	if applicationID != "" {
		return applicationID, applicationName, nil
	}

	if applicationName == "" {
		return "", "", diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "Missing application identifier",
			Detail:   "Either application_id or application_name must be provided.",
		}}
	}

	app, _, err := applicationsegment.GetByName(ctx, service, applicationName)
	if err != nil {
		return "", "", diag.FromErr(fmt.Errorf("failed to find application segment named %s: %w", applicationName, err))
	}

	return app.ID, app.Name, nil
}

func updateWeightedLBConfig(ctx context.Context, service *zscaler.Service, applicationID string, d *schema.ResourceData) diag.Diagnostics {
	config := applicationsegment.WeightedLoadBalancerConfig{
		ApplicationID:         applicationID,
		WeightedLoadBalancing: d.Get("weighted_load_balancing").(bool),
	}

	if v, ok := d.GetOk("application_to_server_group_mappings"); ok {
		mappings, diags := expandApplicationToServerGroupMappings(ctx, service, v.([]interface{}))
		if diags.HasError() {
			return diags
		}
		config.ApplicationToServerGroupMaps = mappings
	}

	log.Printf("[INFO] Updating weighted load balancer config for application %s", applicationID)
	if _, _, err := applicationsegment.UpdateWeightedLoadBalancerConfig(ctx, service, applicationID, config); err != nil {
		return diag.FromErr(fmt.Errorf("failed to update weighted load balancer config for application %s: %w", applicationID, err))
	}

	return nil
}

func expandApplicationToServerGroupMappings(ctx context.Context, service *zscaler.Service, items []interface{}) ([]applicationsegment.ApplicationToServerGroupMapping, diag.Diagnostics) {
	result := make([]applicationsegment.ApplicationToServerGroupMapping, 0, len(items))

	for idx, item := range items {
		if item == nil {
			continue
		}
		data := item.(map[string]interface{})

		id := GetString(data["id"])
		name := GetString(data["name"])

		if id == "" && name == "" {
			return nil, diag.Diagnostics{{
				Severity: diag.Error,
				Summary:  "Missing server group identifier",
				Detail:   fmt.Sprintf("application_to_server_group_mappings[%d] must include either id or name", idx),
			}}
		}

		var resolvedName string
		if name != "" {
			resolvedName = name
		}
		var resolvedID string
		if id != "" {
			resolvedID = id
		}
		if resolvedID == "" && resolvedName != "" {
			group, _, err := servergroup.GetByName(ctx, service, resolvedName)
			if err != nil || group == nil {
				return nil, diag.FromErr(fmt.Errorf("failed to find server group named %s: %w", resolvedName, err))
			}
			resolvedID = group.ID
			resolvedName = group.Name
			data["id"] = resolvedID
			data["name"] = resolvedName
		}

		mapping := applicationsegment.ApplicationToServerGroupMapping{
			ID:      resolvedID,
			Name:    resolvedName,
			Passive: GetBool(data["passive"]),
			Weight:  "0",
		}

		if weightValue, ok := data["weight"]; ok {
			if weight := GetString(weightValue); weight != "" {
				mapping.Weight = weight
			}
		}

		result = append(result, mapping)
	}

	return result, nil
}
