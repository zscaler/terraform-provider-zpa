package resources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_profile"
)

var (
	_ resource.Resource                = &InspectionProfileResource{}
	_ resource.ResourceWithConfigure   = &InspectionProfileResource{}
	_ resource.ResourceWithImportState = &InspectionProfileResource{}
)

func NewInspectionProfileResource() resource.Resource {
	return &InspectionProfileResource{}
}

type InspectionProfileResource struct {
	client *client.Client
}

type InspectionProfileModel struct {
	ID                                types.String                          `tfsdk:"id"`
	Name                              types.String                          `tfsdk:"name"`
	Description                       types.String                          `tfsdk:"description"`
	APIProfile                        types.Bool                            `tfsdk:"api_profile"`
	OverrideAction                    types.String                          `tfsdk:"override_action"`
	AssociateAllControls              types.Bool                            `tfsdk:"associate_all_controls"`
	ControlsInfo                      []InspectionProfileControlInfoModel   `tfsdk:"controls_info"`
	CustomControls                    []InspectionProfileCustomControlModel `tfsdk:"custom_controls"`
	GlobalControlActions              types.Set                             `tfsdk:"global_control_actions"`
	CommonGlobalOverrideActionsConfig types.Map                             `tfsdk:"common_global_override_actions_config"`
	ParanoiaLevel                     types.String                          `tfsdk:"paranoia_level"`
	PredefinedControls                []InspectionProfileCommonControlModel `tfsdk:"predefined_controls"`
	PredefinedAPIControls             []InspectionProfileCommonControlModel `tfsdk:"predefined_api_controls"`
	ThreatLabzControls                []InspectionProfileSimpleControlModel `tfsdk:"threat_labz_controls"`
	WebSocketControls                 []InspectionProfileSimpleControlModel `tfsdk:"websocket_controls"`
	PredefinedControlsVersion         types.String                          `tfsdk:"predefined_controls_version"`
	ZSDefinedControlChoice            types.String                          `tfsdk:"zs_defined_control_choice"`
}

type InspectionProfileControlInfoModel struct {
	ControlType types.String `tfsdk:"control_type"`
	Count       types.String `tfsdk:"count"`
}

type InspectionProfileCustomControlModel struct {
	ID          types.String `tfsdk:"id"`
	Action      types.String `tfsdk:"action"`
	ActionValue types.String `tfsdk:"action_value"`
}

type InspectionProfileCommonControlModel struct {
	ID          types.String `tfsdk:"id"`
	Action      types.String `tfsdk:"action"`
	ActionValue types.String `tfsdk:"action_value"`
}

type InspectionProfileSimpleControlModel struct {
	ID          types.String `tfsdk:"id"`
	Action      types.String `tfsdk:"action"`
	ActionValue types.String `tfsdk:"action_value"`
}

func (r *InspectionProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inspection_profile"
}

func (r *InspectionProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"name": schema.StringAttribute{
			Optional:      true,
			Computed:      true,
			Description:   "Name of the inspection profile.",
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"description": schema.StringAttribute{
			Optional:    true,
			Description: "Description of the inspection profile.",
		},
		"api_profile": schema.BoolAttribute{
			Optional: true,
			Computed: true,
		},
		"override_action": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("COMMON", "NONE", "SPECIFIC"),
			},
		},
		"associate_all_controls": schema.BoolAttribute{
			Optional:      true,
			Computed:      true,
			Default:       booldefault.StaticBool(false),
			PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			Description:   "When enabled, associates all predefined controls with the inspection profile after create/update.",
		},
		"global_control_actions": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
		},
		"common_global_override_actions_config": schema.MapAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
		},
		"paranoia_level": schema.StringAttribute{
			Optional: true,
		},
		"predefined_controls_version": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString("OWASP_CRS/3.3.0"),
		},
		"zs_defined_control_choice": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("ALL", "SPECIFIC"),
			},
		},
	}

	blocks := map[string]schema.Block{
		"controls_info": schema.ListNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"control_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"count": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
		"custom_controls": schema.SetNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required: true,
					},
					"action": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("PASS", "BLOCK", "REDIRECT"),
						},
					},
					"action_value": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
		"predefined_controls": schema.SetNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"id":           schema.StringAttribute{Optional: true},
					"action":       schema.StringAttribute{Optional: true},
					"action_value": schema.StringAttribute{Optional: true},
				},
			},
		},
		"predefined_api_controls": schema.SetNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"id":           schema.StringAttribute{Optional: true},
					"action":       schema.StringAttribute{Optional: true},
					"action_value": schema.StringAttribute{Optional: true},
				},
			},
		},
		"threat_labz_controls": schema.ListNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"id":           schema.StringAttribute{Optional: true},
					"action":       schema.StringAttribute{Optional: true},
					"action_value": schema.StringAttribute{Optional: true},
				},
			},
		},
		"websocket_controls": schema.ListNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"id":           schema.StringAttribute{Optional: true},
					"action":       schema.StringAttribute{Optional: true},
					"action_value": schema.StringAttribute{Optional: true},
				},
			},
		},
	}

	resp.Schema = schema.Schema{
		Description: "Manages ZPA inspection profiles.",
		Attributes:  attributes,
		Blocks:      blocks,
	}
}

func (r *InspectionProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *InspectionProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan InspectionProfileModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diags := r.expand(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := inspection_profile.Create(ctx, r.client.Service, payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create inspection profile: %v", err))
		return
	}

	resp.Diagnostics.Append(r.handleAssociateAllControls(ctx, created.ID, plan.AssociateAllControls)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, readDiags := r.readIntoState(ctx, created.ID, plan.AssociateAllControls)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *InspectionProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state InspectionProfileModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readIntoState(ctx, state.ID.ValueString(), state.AssociateAllControls)
	if diags.HasError() {
		for _, d := range diags {
			if d.Severity() == diag.SeverityError && strings.Contains(strings.ToLower(d.Detail()), "not found") {
				resp.State.RemoveResource(ctx)
				return
			}
		}
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *InspectionProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan InspectionProfileModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	payload, diags := r.expand(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := inspection_profile.Update(ctx, r.client.Service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update inspection profile: %v", err))
		return
	}

	resp.Diagnostics.Append(r.handleAssociateAllControls(ctx, plan.ID.ValueString(), plan.AssociateAllControls)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, readDiags := r.readIntoState(ctx, plan.ID.ValueString(), plan.AssociateAllControls)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *InspectionProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state InspectionProfileModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := inspection_profile.Delete(ctx, r.client.Service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete inspection profile: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *InspectionProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before importing resources.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the inspection profile ID or name.")
		return
	}

	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		profile, _, lookupErr := inspection_profile.GetByName(ctx, r.client.Service, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate inspection profile %q: %v", id, lookupErr))
			return
		}
		id = profile.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *InspectionProfileResource) readIntoState(ctx context.Context, id string, associateFlag types.Bool) (InspectionProfileModel, diag.Diagnostics) {
	profile, _, err := inspection_profile.Get(ctx, r.client.Service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return InspectionProfileModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Inspection profile %s not found", id))}
		}
		return InspectionProfileModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read inspection profile: %v", err))}
	}

	customControls := flattenInspectionProfileCustomControls(profile.CustomControls)
	predefinedControls := flattenInspectionProfileCommonControls(profile.PredefinedControls)
	predefinedAPIControls := flattenInspectionProfileCommonControls(profile.PredefinedAPIControls)
	threatControls := flattenInspectionProfileThreatLabzControls(profile.ThreatLabzControls)
	websocketControls := flattenInspectionProfileWebSocketControls(profile.WebSocketControls)
	controlsInfo := flattenInspectionProfileControlsInfo(profile.ControlInfoResource)

	globalActions, actionsDiags := types.SetValueFrom(ctx, types.StringType, profile.GlobalControlActions)
	if actionsDiags.HasError() {
		return InspectionProfileModel{}, actionsDiags
	}

	overrideConfig, mapDiags := mapToTypesMap(ctx, profile.CommonGlobalOverrideActionsConfig)
	if mapDiags.HasError() {
		return InspectionProfileModel{}, mapDiags
	}

	state := InspectionProfileModel{
		ID:                                helpers.StringValueOrNull(profile.ID),
		Name:                              helpers.StringValueOrNull(profile.Name),
		Description:                       helpers.StringValueOrNull(profile.Description),
		APIProfile:                        types.BoolValue(profile.APIProfile),
		OverrideAction:                    helpers.StringValueOrNull(profile.OverrideAction),
		AssociateAllControls:              associateFlag,
		ControlsInfo:                      controlsInfo,
		CustomControls:                    customControls,
		GlobalControlActions:              globalActions,
		CommonGlobalOverrideActionsConfig: overrideConfig,
		ParanoiaLevel:                     helpers.StringValueOrNull(profile.ParanoiaLevel),
		PredefinedControls:                predefinedControls,
		PredefinedAPIControls:             predefinedAPIControls,
		ThreatLabzControls:                threatControls,
		WebSocketControls:                 websocketControls,
		PredefinedControlsVersion:         helpers.StringValueOrNull(profile.PredefinedControlsVersion),
		ZSDefinedControlChoice:            helpers.StringValueOrNull(profile.ZSDefinedControlChoice),
	}

	return state, nil
}

func (r *InspectionProfileResource) expand(ctx context.Context, plan InspectionProfileModel) (inspection_profile.InspectionProfile, diag.Diagnostics) {
	var diags diag.Diagnostics

	globalActions, actionsDiags := helpers.SetValueToStringSlice(ctx, plan.GlobalControlActions)
	diags.Append(actionsDiags...)

	overrideConfig, mapDiags := mapFromTypesMap(ctx, plan.CommonGlobalOverrideActionsConfig)
	diags.Append(mapDiags...)

	customControls := expandInspectionProfileCustomControls(plan.CustomControls)
	predefinedControls := expandInspectionProfileCommonControls(plan.PredefinedControls)
	predefinedAPIControls := expandInspectionProfileCommonControls(plan.PredefinedAPIControls)
	threatControls := expandInspectionProfileThreatLabzControls(plan.ThreatLabzControls)
	websocketControls := expandInspectionProfileWebSocketControls(plan.WebSocketControls)
	controlsInfo := expandInspectionProfileControlsInfo(plan.ControlsInfo)

	payload := inspection_profile.InspectionProfile{
		ID:                                helpers.StringValue(plan.ID),
		Name:                              helpers.StringValue(plan.Name),
		Description:                       helpers.StringValue(plan.Description),
		APIProfile:                        helpers.BoolValue(plan.APIProfile, false),
		OverrideAction:                    helpers.StringValue(plan.OverrideAction),
		CommonGlobalOverrideActionsConfig: overrideConfig,
		GlobalControlActions:              globalActions,
		ParanoiaLevel:                     helpers.StringValue(plan.ParanoiaLevel),
		PredefinedControlsVersion:         helpers.StringValue(plan.PredefinedControlsVersion),
		ControlInfoResource:               controlsInfo,
		CustomControls:                    customControls,
		PredefinedAPIControls:             predefinedAPIControls,
		PredefinedControls:                predefinedControls,
		ThreatLabzControls:                threatControls,
		WebSocketControls:                 websocketControls,
		ZSDefinedControlChoice:            helpers.StringValue(plan.ZSDefinedControlChoice),
	}

	if payload.PredefinedControlsVersion == "" {
		payload.PredefinedControlsVersion = "OWASP_CRS/3.3.0"
	}

	diags.Append(validateInspectionProfile(payload)...)

	return payload, diags
}

func (r *InspectionProfileResource) handleAssociateAllControls(ctx context.Context, profileID string, flag types.Bool) diag.Diagnostics {
	if flag.IsNull() || flag.IsUnknown() || !flag.ValueBool() {
		return nil
	}

	profile, _, err := inspection_profile.Get(ctx, r.client.Service, profileID)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to fetch inspection profile for association: %v", err))}
	}

	if _, err := inspection_profile.PutAssociate(ctx, r.client.Service, profileID, profile); err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to associate predefined controls: %v", err))}
	}

	return nil
}

func validateInspectionProfile(profile inspection_profile.InspectionProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, control := range profile.CustomControls {
		if strings.EqualFold(control.Action, "REDIRECT") && strings.TrimSpace(control.ActionValue) == "" {
			diags.Append(diag.NewErrorDiagnostic("Validation Error", "when custom_controls.action is REDIRECT, action_value must be set"))
		}
	}

	for _, control := range profile.PredefinedControls {
		if strings.EqualFold(control.Action, "REDIRECT") && strings.TrimSpace(control.ActionValue) == "" {
			diags.Append(diag.NewErrorDiagnostic("Validation Error", "when predefined_controls.action is REDIRECT, action_value must be set"))
		}
	}

	return diags
}

func flattenInspectionProfileCustomControls(controls []inspection_profile.InspectionCustomControl) []InspectionProfileCustomControlModel {
	result := make([]InspectionProfileCustomControlModel, 0, len(controls))
	for _, control := range controls {
		result = append(result, InspectionProfileCustomControlModel{
			ID:          helpers.StringValueOrNull(control.ID),
			Action:      helpers.StringValueOrNull(control.Action),
			ActionValue: helpers.StringValueOrNull(control.ActionValue),
		})
	}
	return result
}

func expandInspectionProfileCustomControls(controls []InspectionProfileCustomControlModel) []inspection_profile.InspectionCustomControl {
	result := make([]inspection_profile.InspectionCustomControl, 0, len(controls))
	for _, control := range controls {
		result = append(result, inspection_profile.InspectionCustomControl{
			ID:          helpers.StringValue(control.ID),
			Action:      helpers.StringValue(control.Action),
			ActionValue: helpers.StringValue(control.ActionValue),
		})
	}
	return result
}

func flattenInspectionProfileCommonControls(controls []common.CustomCommonControls) []InspectionProfileCommonControlModel {
	result := make([]InspectionProfileCommonControlModel, 0, len(controls))
	for _, control := range controls {
		result = append(result, InspectionProfileCommonControlModel{
			ID:          helpers.StringValueOrNull(control.ID),
			Action:      helpers.StringValueOrNull(control.Action),
			ActionValue: helpers.StringValueOrNull(control.ActionValue),
		})
	}
	return result
}

func expandInspectionProfileCommonControls(controls []InspectionProfileCommonControlModel) []common.CustomCommonControls {
	result := make([]common.CustomCommonControls, 0, len(controls))
	for _, control := range controls {
		result = append(result, common.CustomCommonControls{
			ID:          helpers.StringValue(control.ID),
			Action:      helpers.StringValue(control.Action),
			ActionValue: helpers.StringValue(control.ActionValue),
		})
	}
	return result
}

func flattenInspectionProfileThreatLabzControls(controls []inspection_profile.ThreatLabzControls) []InspectionProfileSimpleControlModel {
	result := make([]InspectionProfileSimpleControlModel, 0, len(controls))
	for _, control := range controls {
		result = append(result, InspectionProfileSimpleControlModel{
			ID:          helpers.StringValueOrNull(control.ID),
			Action:      helpers.StringValueOrNull(control.Action),
			ActionValue: helpers.StringValueOrNull(control.ActionValue),
		})
	}
	return result
}

func flattenInspectionProfileWebSocketControls(controls []inspection_profile.WebSocketControls) []InspectionProfileSimpleControlModel {
	result := make([]InspectionProfileSimpleControlModel, 0, len(controls))
	for _, control := range controls {
		result = append(result, InspectionProfileSimpleControlModel{
			ID:          helpers.StringValueOrNull(control.ID),
			Action:      helpers.StringValueOrNull(control.Action),
			ActionValue: helpers.StringValueOrNull(control.ActionValue),
		})
	}
	return result
}

func expandInspectionProfileThreatLabzControls(controls []InspectionProfileSimpleControlModel) []inspection_profile.ThreatLabzControls {
	result := make([]inspection_profile.ThreatLabzControls, 0, len(controls))
	for _, control := range controls {
		result = append(result, inspection_profile.ThreatLabzControls{
			ID:          helpers.StringValue(control.ID),
			Action:      helpers.StringValue(control.Action),
			ActionValue: helpers.StringValue(control.ActionValue),
		})
	}
	return result
}

func expandInspectionProfileWebSocketControls(controls []InspectionProfileSimpleControlModel) []inspection_profile.WebSocketControls {
	result := make([]inspection_profile.WebSocketControls, 0, len(controls))
	for _, control := range controls {
		result = append(result, inspection_profile.WebSocketControls{
			ID:          helpers.StringValue(control.ID),
			Action:      helpers.StringValue(control.Action),
			ActionValue: helpers.StringValue(control.ActionValue),
		})
	}
	return result
}

func expandInspectionProfileControlsInfo(info []InspectionProfileControlInfoModel) []inspection_profile.ControlInfoResource {
	result := make([]inspection_profile.ControlInfoResource, 0, len(info))
	for _, item := range info {
		result = append(result, inspection_profile.ControlInfoResource{
			ControlType: helpers.StringValue(item.ControlType),
		})
	}
	return result
}

func flattenInspectionProfileControlsInfo(info []inspection_profile.ControlInfoResource) []InspectionProfileControlInfoModel {
	result := make([]InspectionProfileControlInfoModel, 0, len(info))
	for _, item := range info {
		result = append(result, InspectionProfileControlInfoModel{
			ControlType: helpers.StringValueOrNull(item.ControlType),
			Count:       helpers.StringValueOrNull(item.Count),
		})
	}
	return result
}

func mapFromTypesMap(ctx context.Context, value types.Map) (map[string]interface{}, diag.Diagnostics) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}

	data := make(map[string]string)
	diags := value.ElementsAs(ctx, &data, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make(map[string]interface{}, len(data))
	for k, v := range data {
		result[k] = v
	}
	return result, nil
}

func mapToTypesMap(ctx context.Context, value map[string]interface{}) (types.Map, diag.Diagnostics) {
	if len(value) == 0 {
		return types.MapNull(types.StringType), nil
	}

	data := make(map[string]string, len(value))
	for k, v := range value {
		data[k] = fmt.Sprint(v)
	}

	mapValue, diags := types.MapValueFrom(ctx, types.StringType, data)
	return mapValue, diags
}
