package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
)

func dataSourceApplicationSegmentWeightedLBConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationSegmentWeightedLBConfigRead,
		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Application segment identifier to query. One of application_id or application_name must be provided.",
			},
			"application_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Application segment name to query. One of application_id or application_name must be provided.",
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional microtenant identifier.",
			},
			"weighted_load_balancing": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if weighted load balancing is enabled for the application segment.",
			},
			"application_to_server_group_mappings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Application to server group mapping details and weights.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Server group mapping identifier.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Server group name.",
						},
						"passive": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the server group is passive.",
						},
						"weight": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Assigned weight for the server group.",
						},
					},
				},
			},
		},
	}
}

func dataSourceApplicationSegmentWeightedLBConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	applicationID := GetString(d.Get("application_id"))
	applicationName := GetString(d.Get("application_name"))

	switch {
	case applicationID != "":
		log.Printf("[INFO] Retrieving weighted LB config for application segment ID %s", applicationID)
	case applicationName != "":
		log.Printf("[INFO] Resolving application segment name %s for weighted LB config", applicationName)
		app, _, err := applicationsegment.GetByName(ctx, service, applicationName)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to find application segment named %s: %w", applicationName, err))
		}
		applicationID = app.ID
		_ = d.Set("application_id", applicationID)
		_ = d.Set("application_name", app.Name)
	default:
		return diag.FromErr(fmt.Errorf("either application_id or application_name must be provided"))
	}

	config, _, err := applicationsegment.GetWeightedLoadBalancerConfig(ctx, service, applicationID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve weighted load balancer config for application %s: %w", applicationID, err))
	}
	if config == nil {
		return diag.FromErr(fmt.Errorf("no weighted load balancer config returned for application %s", applicationID))
	}

	log.Printf("[INFO] Retrieved weighted load balancer config for application %s", applicationID)

	_ = d.Set("weighted_load_balancing", config.WeightedLoadBalancing)
	if err := d.Set("application_to_server_group_mappings", flattenApplicationToServerGroupMappings(config.ApplicationToServerGroupMaps)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set application_to_server_group_mappings: %w", err))
	}

	if err := d.Set("application_id", applicationID); err != nil {
		return diag.FromErr(err)
	}
	if applicationName != "" {
		_ = d.Set("application_name", applicationName)
	}

	stateID := applicationID
	if microTenantID != "" {
		stateID = fmt.Sprintf("%s:%s", microTenantID, applicationID)
	}
	d.SetId(stateID)

	return nil
}

func flattenApplicationToServerGroupMappings(mappings []applicationsegment.ApplicationToServerGroupMapping) []map[string]interface{} {
	if len(mappings) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(mappings))
	for i, mapping := range mappings {
		result[i] = map[string]interface{}{
			"id":      mapping.ID,
			"name":    mapping.Name,
			"passive": mapping.Passive,
			"weight":  mapping.Weight,
		}
	}
	return result
}
