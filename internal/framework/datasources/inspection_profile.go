package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	inspectionprofile "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_profile"
)

var (
	_ datasource.DataSource              = &InspectionProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &InspectionProfileDataSource{}
)

func NewInspectionProfileDataSource() datasource.DataSource {
	return &InspectionProfileDataSource{}
}

type InspectionProfileDataSource struct {
	client *client.Client
}

type InspectionProfileModel struct {
	ID                           types.String `tfsdk:"id"`
	Name                         types.String `tfsdk:"name"`
	Description                  types.String `tfsdk:"description"`
	APIProfile                   types.Bool   `tfsdk:"api_profile"`
	OverrideAction               types.String `tfsdk:"override_action"`
	CommonGlobalOverrideActions  types.Map    `tfsdk:"common_global_override_actions_config"`
	CreationTime                 types.String `tfsdk:"creation_time"`
	ZSDefinedControlChoice       types.String `tfsdk:"zs_defined_control_choice"`
	GlobalControlActions         types.List   `tfsdk:"global_control_actions"`
	IncarnationNumber            types.String `tfsdk:"incarnation_number"`
	ModifiedBy                   types.String `tfsdk:"modified_by"`
	ModifiedTime                 types.String `tfsdk:"modified_time"`
	ParanoiaLevel                types.String `tfsdk:"paranoia_level"`
	PredefinedControlsVersion    types.String `tfsdk:"predefined_controls_version"`
	CheckControlDeploymentStatus types.Bool   `tfsdk:"check_control_deployment_status"`
	ControlsInfo                 types.List   `tfsdk:"controls_info"`
	CustomControls               types.List   `tfsdk:"custom_controls"`
	PredefinedControls           types.List   `tfsdk:"predefined_controls"`
	WebSocketControls            types.List   `tfsdk:"web_socket_controls"`
}

func (d *InspectionProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inspection_profile"
}

func (d *InspectionProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves an inspection profile by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the inspection profile.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the inspection profile.",
			},
			"description":     schema.StringAttribute{Computed: true},
			"api_profile":     schema.BoolAttribute{Computed: true},
			"override_action": schema.StringAttribute{Computed: true},
			"common_global_override_actions_config": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"creation_time":             schema.StringAttribute{Computed: true},
			"zs_defined_control_choice": schema.StringAttribute{Computed: true},
			"global_control_actions": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"incarnation_number":              schema.StringAttribute{Computed: true},
			"modified_by":                     schema.StringAttribute{Computed: true},
			"modified_time":                   schema.StringAttribute{Computed: true},
			"paranoia_level":                  schema.StringAttribute{Computed: true},
			"predefined_controls_version":     schema.StringAttribute{Computed: true},
			"check_control_deployment_status": schema.BoolAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			"controls_info": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"control_type": schema.StringAttribute{Computed: true},
						"count":        schema.StringAttribute{Computed: true},
					},
				},
			},
			"custom_controls": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"action":               schema.StringAttribute{Computed: true},
						"action_value":         schema.StringAttribute{Computed: true},
						"control_number":       schema.StringAttribute{Computed: true},
						"control_rule_json":    schema.StringAttribute{Computed: true},
						"control_type":         schema.StringAttribute{Computed: true},
						"creation_time":        schema.StringAttribute{Computed: true},
						"default_action":       schema.StringAttribute{Computed: true},
						"default_action_value": schema.StringAttribute{Computed: true},
						"description":          schema.StringAttribute{Computed: true},
						"id":                   schema.StringAttribute{Computed: true},
						"modified_by":          schema.StringAttribute{Computed: true},
						"modified_time":        schema.StringAttribute{Computed: true},
						"name":                 schema.StringAttribute{Computed: true},
						"paranoia_level":       schema.StringAttribute{Computed: true},
						"protocol_type":        schema.StringAttribute{Computed: true},
						"severity":             schema.StringAttribute{Computed: true},
						"type":                 schema.StringAttribute{Computed: true},
						"version":              schema.StringAttribute{Computed: true},
					},
					Blocks: map[string]schema.Block{
						"associated_inspection_profile_names": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id":   schema.StringAttribute{Computed: true},
									"name": schema.StringAttribute{Computed: true},
								},
							},
						},
						"rules": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{Computed: true},
									"names": schema.SetAttribute{
										ElementType: types.StringType,
										Computed:    true,
									},
								},
								Blocks: map[string]schema.Block{
									"conditions": schema.ListNestedBlock{
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"lhs": schema.StringAttribute{Computed: true},
												"op":  schema.StringAttribute{Computed: true},
												"rhs": schema.StringAttribute{Computed: true},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"predefined_controls": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"action":               schema.StringAttribute{Computed: true},
						"action_value":         schema.StringAttribute{Computed: true},
						"attachment":           schema.StringAttribute{Computed: true},
						"control_group":        schema.StringAttribute{Computed: true},
						"control_number":       schema.StringAttribute{Computed: true},
						"control_type":         schema.StringAttribute{Computed: true},
						"creation_time":        schema.StringAttribute{Computed: true},
						"default_action":       schema.StringAttribute{Computed: true},
						"default_action_value": schema.StringAttribute{Computed: true},
						"description":          schema.StringAttribute{Computed: true},
						"id":                   schema.StringAttribute{Computed: true},
						"modified_by":          schema.StringAttribute{Computed: true},
						"modified_time":        schema.StringAttribute{Computed: true},
						"name":                 schema.StringAttribute{Computed: true},
						"paranoia_level":       schema.StringAttribute{Computed: true},
						"protocol_type":        schema.StringAttribute{Computed: true},
						"severity":             schema.StringAttribute{Computed: true},
						"version":              schema.StringAttribute{Computed: true},
					},
					Blocks: map[string]schema.Block{
						"associated_inspection_profile_names": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id":   schema.StringAttribute{Computed: true},
									"name": schema.StringAttribute{Computed: true},
								},
							},
						},
					},
				},
			},
			"web_socket_controls": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"action":                    schema.StringAttribute{Computed: true},
						"action_value":              schema.StringAttribute{Computed: true},
						"control_number":            schema.StringAttribute{Computed: true},
						"control_type":              schema.StringAttribute{Computed: true},
						"creation_time":             schema.StringAttribute{Computed: true},
						"default_action":            schema.StringAttribute{Computed: true},
						"default_action_value":      schema.StringAttribute{Computed: true},
						"description":               schema.StringAttribute{Computed: true},
						"id":                        schema.StringAttribute{Computed: true},
						"modified_by":               schema.StringAttribute{Computed: true},
						"modified_time":             schema.StringAttribute{Computed: true},
						"name":                      schema.StringAttribute{Computed: true},
						"paranoia_level":            schema.StringAttribute{Computed: true},
						"severity":                  schema.StringAttribute{Computed: true},
						"version":                   schema.StringAttribute{Computed: true},
						"zs_defined_control_choice": schema.StringAttribute{Computed: true},
					},
					Blocks: map[string]schema.Block{
						"associated_inspection_profile_names": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id":   schema.StringAttribute{Computed: true},
									"name": schema.StringAttribute{Computed: true},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *InspectionProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cl, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = cl
}

func (d *InspectionProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data InspectionProfileModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())
	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided to read an inspection profile.")
		return
	}

	var (
		profile *inspectionprofile.InspectionProfile
		err     error
	)

	if id != "" {
		tflog.Debug(ctx, "Retrieving inspection profile by ID", map[string]any{"id": id})
		profile, _, err = inspectionprofile.Get(ctx, d.client.Service, id)
	} else {
		tflog.Debug(ctx, "Retrieving inspection profile by name", map[string]any{"name": name})
		profile, _, err = inspectionprofile.GetByName(ctx, d.client.Service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read inspection profile: %v", err))
		return
	}

	if profile == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Inspection profile with id %q or name %q was not found.", id, name))
		return
	}

	mapValue, mapDiags := mapInterfaceToStringMap(ctx, profile.CommonGlobalOverrideActionsConfig)
	resp.Diagnostics.Append(mapDiags...)

	controlsInfo, controlsDiags := flattenControlsInfo(ctx, profile.ControlInfoResource)
	resp.Diagnostics.Append(controlsDiags...)

	customControls, customDiags := flattenInspectionCustomControls(ctx, profile.CustomControls)
	resp.Diagnostics.Append(customDiags...)

	predefinedControls, predefinedDiags := flattenCommonControls(ctx, profile.PredefinedControls)
	resp.Diagnostics.Append(predefinedDiags...)

	webSocketControls, wsDiags := flattenWebSocketControls(ctx, profile.WebSocketControls)
	resp.Diagnostics.Append(wsDiags...)

	globalControlActions, listDiags := types.ListValueFrom(ctx, types.StringType, profile.GlobalControlActions)
	resp.Diagnostics.Append(listDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(profile.ID)
	data.Name = types.StringValue(profile.Name)
	data.Description = types.StringValue(profile.Description)
	data.APIProfile = types.BoolValue(profile.APIProfile)
	data.OverrideAction = stringOrNull(profile.OverrideAction)
	data.CommonGlobalOverrideActions = mapValue
	data.CreationTime = types.StringValue(profile.CreationTime)
	data.ZSDefinedControlChoice = stringOrNull(profile.ZSDefinedControlChoice)
	data.GlobalControlActions = globalControlActions
	data.IncarnationNumber = stringOrNull(profile.IncarnationNumber)
	data.ModifiedBy = stringOrNull(profile.ModifiedBy)
	data.ModifiedTime = stringOrNull(profile.ModifiedTime)
	data.ParanoiaLevel = stringOrNull(profile.ParanoiaLevel)
	data.PredefinedControlsVersion = stringOrNull(profile.PredefinedControlsVersion)
	data.CheckControlDeploymentStatus = types.BoolValue(profile.CheckControlDeploymentStatus)
	data.ControlsInfo = controlsInfo
	data.CustomControls = customControls
	data.PredefinedControls = predefinedControls
	data.WebSocketControls = webSocketControls

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenControlsInfo(ctx context.Context, controls []inspectionprofile.ControlInfoResource) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"control_type": types.StringType,
		"count":        types.StringType,
	}

	if len(controls) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(controls))
	var diags diag.Diagnostics
	for _, control := range controls {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"control_type": stringOrNull(control.ControlType),
			"count":        stringOrNull(control.Count),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}

func flattenInspectionCustomControls(ctx context.Context, controls []inspectionprofile.InspectionCustomControl) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"action":                              types.StringType,
		"action_value":                        types.StringType,
		"control_number":                      types.StringType,
		"control_rule_json":                   types.StringType,
		"control_type":                        types.StringType,
		"creation_time":                       types.StringType,
		"default_action":                      types.StringType,
		"default_action_value":                types.StringType,
		"description":                         types.StringType,
		"id":                                  types.StringType,
		"modified_by":                         types.StringType,
		"modified_time":                       types.StringType,
		"name":                                types.StringType,
		"paranoia_level":                      types.StringType,
		"protocol_type":                       types.StringType,
		"severity":                            types.StringType,
		"type":                                types.StringType,
		"version":                             types.StringType,
		"associated_inspection_profile_names": types.ListType{ElemType: types.ObjectType{AttrTypes: associatedProfileAttrTypes()}},
		"rules":                               types.ListType{ElemType: types.ObjectType{AttrTypes: inspectionRuleAttrTypes()}},
	}

	if len(controls) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(controls))
	var diags diag.Diagnostics
	for _, control := range controls {
		associated, assocDiags := flattenAssociatedProfileNamesAsList(ctx, control.AssociatedInspectionProfileNames)
		diags.Append(assocDiags...)
		rules, rulesDiags := flattenInspectionRules(ctx, control.Rules)
		diags.Append(rulesDiags...)

		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"action":                              stringOrNull(control.Action),
			"action_value":                        stringOrNull(control.ActionValue),
			"control_number":                      stringOrNull(control.ControlNumber),
			"control_rule_json":                   stringOrNull(control.ControlRuleJson),
			"control_type":                        stringOrNull(control.ControlType),
			"creation_time":                       stringOrNull(control.CreationTime),
			"default_action":                      stringOrNull(control.DefaultAction),
			"default_action_value":                stringOrNull(control.DefaultActionValue),
			"description":                         stringOrNull(control.Description),
			"id":                                  stringOrNull(control.ID),
			"modified_by":                         stringOrNull(control.ModifiedBy),
			"modified_time":                       stringOrNull(control.ModifiedTime),
			"name":                                stringOrNull(control.Name),
			"paranoia_level":                      stringOrNull(control.ParanoiaLevel),
			"protocol_type":                       stringOrNull(control.ProtocolType),
			"severity":                            stringOrNull(control.Severity),
			"type":                                stringOrNull(control.Type),
			"version":                             stringOrNull(control.Version),
			"associated_inspection_profile_names": associated,
			"rules":                               rules,
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}

func flattenInspectionRules(ctx context.Context, rules []common.Rules) (types.List, diag.Diagnostics) {
	attrTypes := inspectionRuleAttrTypes()

	if len(rules) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(rules))
	var diags diag.Diagnostics
	for _, rule := range rules {
		conditions, condDiags := flattenRuleConditions(ctx, rule.Conditions)
		diags.Append(condDiags...)

		namesSlice := make([]string, 0)
		if trimmed := strings.TrimSpace(rule.Names); trimmed != "" {
			for _, part := range strings.Split(trimmed, ",") {
				if s := strings.TrimSpace(part); s != "" {
					namesSlice = append(namesSlice, s)
				}
			}
			if len(namesSlice) == 0 {
				namesSlice = append(namesSlice, trimmed)
			}
		}

		namesSet, namesDiags := types.SetValueFrom(ctx, types.StringType, namesSlice)
		diags.Append(namesDiags...)

		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"type":       stringOrNull(rule.Type),
			"names":      namesSet,
			"conditions": conditions,
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}

func inspectionRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":       types.StringType,
		"names":      types.SetType{ElemType: types.StringType},
		"conditions": types.ListType{ElemType: types.ObjectType{AttrTypes: ruleConditionAttrTypes()}},
	}
}

func flattenRuleConditions(ctx context.Context, conditions []common.Conditions) (types.List, diag.Diagnostics) {
	attrTypes := ruleConditionAttrTypes()

	if len(conditions) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(conditions))
	var diags diag.Diagnostics
	for _, condition := range conditions {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"lhs": stringOrNull(condition.LHS),
			"op":  stringOrNull(condition.OP),
			"rhs": stringOrNull(condition.RHS),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}

func ruleConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"lhs": types.StringType,
		"op":  types.StringType,
		"rhs": types.StringType,
	}
}

func flattenCommonControls(ctx context.Context, controls []common.CustomCommonControls) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"action":                              types.StringType,
		"action_value":                        types.StringType,
		"attachment":                          types.StringType,
		"control_group":                       types.StringType,
		"control_number":                      types.StringType,
		"control_type":                        types.StringType,
		"creation_time":                       types.StringType,
		"default_action":                      types.StringType,
		"default_action_value":                types.StringType,
		"description":                         types.StringType,
		"id":                                  types.StringType,
		"modified_by":                         types.StringType,
		"modified_time":                       types.StringType,
		"name":                                types.StringType,
		"paranoia_level":                      types.StringType,
		"protocol_type":                       types.StringType,
		"severity":                            types.StringType,
		"version":                             types.StringType,
		"associated_inspection_profile_names": types.ListType{ElemType: types.ObjectType{AttrTypes: associatedProfileAttrTypes()}},
	}

	if len(controls) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(controls))
	var diags diag.Diagnostics
	for _, control := range controls {
		associated, assocDiags := flattenAssociatedProfileNamesAsList(ctx, control.AssociatedInspectionProfileNames)
		diags.Append(assocDiags...)

		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"action":                              stringOrNull(control.Action),
			"action_value":                        stringOrNull(control.ActionValue),
			"attachment":                          stringOrNull(control.Attachment),
			"control_group":                       stringOrNull(control.ControlGroup),
			"control_number":                      stringOrNull(control.ControlNumber),
			"control_type":                        stringOrNull(control.ControlType),
			"creation_time":                       stringOrNull(control.CreationTime),
			"default_action":                      stringOrNull(control.DefaultAction),
			"default_action_value":                stringOrNull(control.DefaultActionValue),
			"description":                         stringOrNull(control.Description),
			"id":                                  stringOrNull(control.ID),
			"modified_by":                         stringOrNull(control.ModifiedBy),
			"modified_time":                       stringOrNull(control.ModifiedTime),
			"name":                                stringOrNull(control.Name),
			"paranoia_level":                      stringOrNull(control.ParanoiaLevel),
			"protocol_type":                       stringOrNull(control.ProtocolType),
			"severity":                            stringOrNull(control.Severity),
			"version":                             stringOrNull(control.Version),
			"associated_inspection_profile_names": associated,
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}

func flattenWebSocketControls(ctx context.Context, controls []inspectionprofile.WebSocketControls) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"action":                              types.StringType,
		"action_value":                        types.StringType,
		"control_number":                      types.StringType,
		"control_type":                        types.StringType,
		"creation_time":                       types.StringType,
		"default_action":                      types.StringType,
		"default_action_value":                types.StringType,
		"description":                         types.StringType,
		"id":                                  types.StringType,
		"modified_by":                         types.StringType,
		"modified_time":                       types.StringType,
		"name":                                types.StringType,
		"paranoia_level":                      types.StringType,
		"severity":                            types.StringType,
		"version":                             types.StringType,
		"zs_defined_control_choice":           types.StringType,
		"associated_inspection_profile_names": types.ListType{ElemType: types.ObjectType{AttrTypes: associatedProfileAttrTypes()}},
	}

	if len(controls) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(controls))
	var diags diag.Diagnostics
	for _, control := range controls {
		associated, assocDiags := flattenAssociatedProfileNamesAsList(ctx, control.AssociatedInspectionProfileNames)
		diags.Append(assocDiags...)

		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"action":                              stringOrNull(control.Action),
			"action_value":                        stringOrNull(control.ActionValue),
			"control_number":                      stringOrNull(control.ControlNumber),
			"control_type":                        stringOrNull(control.ControlType),
			"creation_time":                       stringOrNull(control.CreationTime),
			"default_action":                      stringOrNull(control.DefaultAction),
			"default_action_value":                stringOrNull(control.DefaultActionValue),
			"description":                         stringOrNull(control.Description),
			"id":                                  stringOrNull(control.ID),
			"modified_by":                         stringOrNull(control.ModifiedBy),
			"modified_time":                       stringOrNull(control.ModifiedTime),
			"name":                                stringOrNull(control.Name),
			"paranoia_level":                      stringOrNull(control.ParanoiaLevel),
			"severity":                            stringOrNull(control.Severity),
			"version":                             stringOrNull(control.Version),
			"zs_defined_control_choice":           stringOrNull(control.ZSDefinedControlChoice),
			"associated_inspection_profile_names": associated,
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}

func associatedProfileAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
}

func flattenAssociatedProfileNames(ctx context.Context, names []common.AssociatedProfileNames) (types.Set, diag.Diagnostics) {
	attrTypes := associatedProfileAttrTypes()

	if len(names) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(names))
	var diags diag.Diagnostics
	for _, name := range names {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":   stringOrNull(name.ID),
			"name": stringOrNull(name.Name),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	set, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(setDiags...)
	return set, diags
}

func flattenAssociatedProfileNamesAsList(ctx context.Context, names []common.AssociatedProfileNames) (types.List, diag.Diagnostics) {
	attrTypes := associatedProfileAttrTypes()

	if len(names) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(names))
	var diags diag.Diagnostics
	for _, name := range names {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":   stringOrNull(name.ID),
			"name": stringOrNull(name.Name),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}
