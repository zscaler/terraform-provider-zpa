package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
)

var (
	_ resource.Resource                = &ApplicationSegmentWeightedLBConfigResource{}
	_ resource.ResourceWithConfigure   = &ApplicationSegmentWeightedLBConfigResource{}
	_ resource.ResourceWithImportState = &ApplicationSegmentWeightedLBConfigResource{}
)

func NewApplicationSegmentWeightedLBConfigResource() resource.Resource {
	return &ApplicationSegmentWeightedLBConfigResource{}
}

type ApplicationSegmentWeightedLBConfigResource struct {
	client *client.Client
}

type ApplicationSegmentWeightedLBConfigModel struct {
	ID                               types.String                           `tfsdk:"id"`
	ApplicationID                    types.String                           `tfsdk:"application_id"`
	ApplicationName                  types.String                           `tfsdk:"application_name"`
	MicrotenantID                    types.String                           `tfsdk:"microtenant_id"`
	WeightedLoadBalancing            types.Bool                             `tfsdk:"weighted_load_balancing"`
	ApplicationToServerGroupMappings []ApplicationToServerGroupMappingModel `tfsdk:"application_to_server_group_mappings"`
}

type ApplicationToServerGroupMappingModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Passive types.Bool   `tfsdk:"passive"`
	Weight  types.String `tfsdk:"weight"`
}

func (r *ApplicationSegmentWeightedLBConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_segment_weightedlb_config"
}

func (r *ApplicationSegmentWeightedLBConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages weighted load balancer configuration for an application segment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"application_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Application segment identifier to manage. Either application_id or application_name must be provided.",
			},
			"application_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Application segment name to manage. Either application_id or application_name must be provided.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Description: "Optional microtenant identifier.",
			},
			"weighted_load_balancing": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable or disable weighted load balancing for the application segment.",
			},
		},
		Blocks: map[string]schema.Block{
			"application_to_server_group_mappings": schema.ListNestedBlock{
				Description: "Application to server group mapping details and weights.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Server group mapping identifier.",
						},
						"name": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Server group name.",
						},
						"passive": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Whether the server group is passive.",
						},
						"weight": schema.StringAttribute{
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

func (r *ApplicationSegmentWeightedLBConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cl, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = cl
}

func (r *ApplicationSegmentWeightedLBConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan ApplicationSegmentWeightedLBConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	applicationID, applicationName, diags := r.resolveApplicationSegmentIdentity(ctx, service, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.updateWeightedLBConfig(ctx, service, applicationID, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(applicationID)
	if applicationName != "" {
		plan.ApplicationName = types.StringValue(applicationName)
	}
	plan.ApplicationID = types.StringValue(applicationID)

	state, readDiags := r.readWeightedLBConfig(ctx, service, applicationID, plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = plan.ID
	state.ApplicationID = plan.ApplicationID
	state.ApplicationName = plan.ApplicationName

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationSegmentWeightedLBConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state ApplicationSegmentWeightedLBConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	applicationID := state.ID.ValueString()
	if applicationID == "" {
		applicationID = helpers.StringValue(state.ApplicationID)
	}
	if applicationID == "" {
		resp.Diagnostics.AddError("Validation Error", "application_id is not set in state")
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	newState, diags := r.readWeightedLBConfig(ctx, service, applicationID, state.MicrotenantID)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	newState.ID = types.StringValue(applicationID)
	newState.ApplicationID = types.StringValue(applicationID)
	if !state.ApplicationName.IsNull() && !state.ApplicationName.IsUnknown() {
		newState.ApplicationName = state.ApplicationName
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ApplicationSegmentWeightedLBConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan ApplicationSegmentWeightedLBConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	applicationID := plan.ID.ValueString()
	if applicationID == "" {
		applicationID = helpers.StringValue(plan.ApplicationID)
	}
	if applicationID == "" {
		resp.Diagnostics.AddError("Validation Error", "application_id must be set to update weighted load balancer config")
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	diags := r.updateWeightedLBConfig(ctx, service, applicationID, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, readDiags := r.readWeightedLBConfig(ctx, service, applicationID, plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = plan.ID
	state.ApplicationID = types.StringValue(applicationID)
	if !plan.ApplicationName.IsNull() && !plan.ApplicationName.IsUnknown() {
		state.ApplicationName = plan.ApplicationName
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationSegmentWeightedLBConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state ApplicationSegmentWeightedLBConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	applicationID := state.ID.ValueString()
	if applicationID == "" {
		applicationID = helpers.StringValue(state.ApplicationID)
	}
	if applicationID == "" {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	payload := applicationsegment.WeightedLoadBalancerConfig{
		ApplicationID:         applicationID,
		WeightedLoadBalancing: false,
	}

	if _, _, err := applicationsegment.UpdateWeightedLoadBalancerConfig(ctx, service, applicationID, payload); err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to disable weighted load balancer config for application %s: %v", applicationID, err),
		)
		return
	}
}

func (r *ApplicationSegmentWeightedLBConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before importing weighted LB config.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the application segment ID.")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), types.StringValue(id))...)
}

func (r *ApplicationSegmentWeightedLBConfigResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	id := helpers.StringValue(microtenantID)
	if id != "" {
		service = service.WithMicroTenant(id)
	}
	return service
}

func (r *ApplicationSegmentWeightedLBConfigResource) resolveApplicationSegmentIdentity(ctx context.Context, service *zscaler.Service, plan *ApplicationSegmentWeightedLBConfigModel) (string, string, diag.Diagnostics) {
	applicationID := helpers.StringValue(plan.ApplicationID)
	applicationName := helpers.StringValue(plan.ApplicationName)

	if applicationID != "" {
		return applicationID, applicationName, nil
	}

	if applicationName == "" {
		return "", "", diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Missing application identifier",
				"Either application_id or application_name must be provided.",
			),
		}
	}

	app, _, err := applicationsegment.GetByName(ctx, service, applicationName)
	if err != nil {
		return "", "", diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Client Error",
				fmt.Sprintf("Failed to find application segment named %s: %v", applicationName, err),
			),
		}
	}

	return app.ID, app.Name, nil
}

func (r *ApplicationSegmentWeightedLBConfigResource) updateWeightedLBConfig(ctx context.Context, service *zscaler.Service, applicationID string, plan *ApplicationSegmentWeightedLBConfigModel) diag.Diagnostics {
	var diags diag.Diagnostics

	config := applicationsegment.WeightedLoadBalancerConfig{
		ApplicationID:         applicationID,
		WeightedLoadBalancing: helpers.BoolValue(plan.WeightedLoadBalancing, false),
	}

	if len(plan.ApplicationToServerGroupMappings) > 0 {
		mappings, mappingDiags := r.expandApplicationToServerGroupMappings(ctx, service, plan.ApplicationToServerGroupMappings)
		diags.Append(mappingDiags...)
		if diags.HasError() {
			return diags
		}
		config.ApplicationToServerGroupMaps = mappings
	}

	if _, _, err := applicationsegment.UpdateWeightedLoadBalancerConfig(ctx, service, applicationID, config); err != nil {
		diags.AddError(
			"Client Error",
			fmt.Sprintf("Failed to update weighted load balancer config for application %s: %v", applicationID, err),
		)
		return diags
	}

	return diags
}

func (r *ApplicationSegmentWeightedLBConfigResource) readWeightedLBConfig(ctx context.Context, service *zscaler.Service, applicationID string, microtenantID types.String) (ApplicationSegmentWeightedLBConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	config, _, err := applicationsegment.GetWeightedLoadBalancerConfig(ctx, service, applicationID)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return ApplicationSegmentWeightedLBConfigModel{}, diag.Diagnostics{
				diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Weighted load balancer config not found for application %s", applicationID)),
			}
		}
		return ApplicationSegmentWeightedLBConfigModel{}, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Client Error",
				fmt.Sprintf("Failed to retrieve weighted load balancer config for application %s: %v", applicationID, err),
			),
		}
	}

	if config == nil {
		return ApplicationSegmentWeightedLBConfigModel{}, diag.Diagnostics{
			diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("No weighted load balancer config returned for application %s", applicationID)),
		}
	}

	state := ApplicationSegmentWeightedLBConfigModel{
		WeightedLoadBalancing: types.BoolValue(config.WeightedLoadBalancing),
		MicrotenantID:         microtenantID,
	}

	mappings, mappingDiags := r.flattenApplicationToServerGroupMappings(ctx, config.ApplicationToServerGroupMaps)
	diags.Append(mappingDiags...)
	state.ApplicationToServerGroupMappings = mappings

	if app, _, err := applicationsegment.Get(ctx, service, applicationID); err == nil && app != nil {
		state.ApplicationName = helpers.StringValueOrNull(app.Name)
	}

	return state, diags
}

func (r *ApplicationSegmentWeightedLBConfigResource) expandApplicationToServerGroupMappings(ctx context.Context, service *zscaler.Service, models []ApplicationToServerGroupMappingModel) ([]applicationsegment.ApplicationToServerGroupMapping, diag.Diagnostics) {
	var diags diag.Diagnostics
	result := make([]applicationsegment.ApplicationToServerGroupMapping, 0, len(models))

	for idx, model := range models {
		id := helpers.StringValue(model.ID)
		name := helpers.StringValue(model.Name)

		if id == "" && name == "" {
			diags.AddError(
				"Missing server group identifier",
				fmt.Sprintf("application_to_server_group_mappings[%d] must include either id or name", idx),
			)
			continue
		}

		var resolvedName string
		var resolvedID string

		if id != "" {
			resolvedID = id
		}
		if name != "" {
			resolvedName = name
		}

		if resolvedID == "" && resolvedName != "" {
			group, _, err := servergroup.GetByName(ctx, service, resolvedName)
			if err != nil || group == nil {
				diags.AddError(
					"Client Error",
					fmt.Sprintf("Failed to find server group named %s: %v", resolvedName, err),
				)
				continue
			}
			resolvedID = group.ID
			resolvedName = group.Name
		}

		mapping := applicationsegment.ApplicationToServerGroupMapping{
			ID:      resolvedID,
			Name:    resolvedName,
			Passive: helpers.BoolValue(model.Passive, false),
			Weight:  "0",
		}

		if !model.Weight.IsNull() && !model.Weight.IsUnknown() {
			weight := helpers.StringValue(model.Weight)
			if weight != "" {
				mapping.Weight = weight
			}
		}

		result = append(result, mapping)
	}

	return result, diags
}

func (r *ApplicationSegmentWeightedLBConfigResource) flattenApplicationToServerGroupMappings(ctx context.Context, mappings []applicationsegment.ApplicationToServerGroupMapping) ([]ApplicationToServerGroupMappingModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(mappings) == 0 {
		return nil, diags
	}

	result := make([]ApplicationToServerGroupMappingModel, 0, len(mappings))
	for _, mapping := range mappings {
		result = append(result, ApplicationToServerGroupMappingModel{
			ID:      helpers.StringValueOrNull(mapping.ID),
			Name:    helpers.StringValueOrNull(mapping.Name),
			Passive: types.BoolValue(mapping.Passive),
			Weight:  helpers.StringValueOrNull(mapping.Weight),
		})
	}

	return result, diags
}
