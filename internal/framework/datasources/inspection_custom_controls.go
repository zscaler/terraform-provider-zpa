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
	inspectioncustom "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_custom_controls"
)

var (
	_ datasource.DataSource              = &InspectionCustomControlsDataSource{}
	_ datasource.DataSourceWithConfigure = &InspectionCustomControlsDataSource{}
)

func NewInspectionCustomControlsDataSource() datasource.DataSource {
	return &InspectionCustomControlsDataSource{}
}

type InspectionCustomControlsDataSource struct {
	client *client.Client
}

type InspectionCustomControlsModel struct {
	ID                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	Action                           types.String `tfsdk:"action"`
	ActionValue                      types.String `tfsdk:"action_value"`
	ControlNumber                    types.String `tfsdk:"control_number"`
	ControlRuleJSON                  types.String `tfsdk:"control_rule_json"`
	ControlType                      types.String `tfsdk:"control_type"`
	CreationTime                     types.String `tfsdk:"creation_time"`
	DefaultAction                    types.String `tfsdk:"default_action"`
	DefaultActionValue               types.String `tfsdk:"default_action_value"`
	Description                      types.String `tfsdk:"description"`
	ModifiedBy                       types.String `tfsdk:"modifiedby"`
	ModifiedTime                     types.String `tfsdk:"modified_time"`
	ParanoiaLevel                    types.String `tfsdk:"paranoia_level"`
	ProtocolType                     types.String `tfsdk:"protocol_type"`
	Severity                         types.String `tfsdk:"severity"`
	Type                             types.String `tfsdk:"type"`
	Version                          types.String `tfsdk:"version"`
	Rules                            types.List   `tfsdk:"rules"`
	AssociatedInspectionProfileNames types.Set    `tfsdk:"associated_inspection_profile_names"`
}

func (d *InspectionCustomControlsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inspection_custom_controls"
}

func (d *InspectionCustomControlsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a custom inspection control by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the custom control.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the custom control.",
			},
			"action":               schema.StringAttribute{Computed: true},
			"action_value":         schema.StringAttribute{Computed: true},
			"control_number":       schema.StringAttribute{Computed: true},
			"control_rule_json":    schema.StringAttribute{Computed: true},
			"control_type":         schema.StringAttribute{Computed: true},
			"creation_time":        schema.StringAttribute{Computed: true},
			"default_action":       schema.StringAttribute{Computed: true},
			"default_action_value": schema.StringAttribute{Computed: true},
			"description":          schema.StringAttribute{Computed: true},
			"modifiedby":           schema.StringAttribute{Computed: true},
			"modified_time":        schema.StringAttribute{Computed: true},
			"paranoia_level":       schema.StringAttribute{Computed: true},
			"protocol_type":        schema.StringAttribute{Computed: true},
			"severity":             schema.StringAttribute{Computed: true},
			"type":                 schema.StringAttribute{Computed: true},
			"version":              schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
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
			"associated_inspection_profile_names": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *InspectionCustomControlsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *InspectionCustomControlsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data InspectionCustomControlsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())
	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided to read a custom control.")
		return
	}

	var (
		control *inspectioncustom.InspectionCustomControl
		err     error
	)

	if id != "" {
		tflog.Debug(ctx, "Retrieving custom inspection control by ID", map[string]any{"id": id})
		control, _, err = inspectioncustom.Get(ctx, d.client.Service, id)
	} else {
		tflog.Debug(ctx, "Retrieving custom inspection control by name", map[string]any{"name": name})
		control, _, err = inspectioncustom.GetByName(ctx, d.client.Service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read custom inspection control: %v", err))
		return
	}

	if control == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Custom inspection control with id %q or name %q was not found.", id, name))
		return
	}

	associated, assocDiags := flattenAssociatedProfileNames(ctx, control.AssociatedInspectionProfileNames)
	resp.Diagnostics.Append(assocDiags...)

	rules, rulesDiags := flattenInspectionCustomRules(ctx, control.Rules)
	resp.Diagnostics.Append(rulesDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(control.ID)
	data.Name = stringOrNull(control.Name)
	data.Action = stringOrNull(control.Action)
	data.ActionValue = stringOrNull(control.ActionValue)
	data.ControlNumber = stringOrNull(control.ControlNumber)
	data.ControlRuleJSON = stringOrNull(control.ControlRuleJson)
	data.ControlType = stringOrNull(control.ControlType)
	data.CreationTime = stringOrNull(control.CreationTime)
	data.DefaultAction = stringOrNull(control.DefaultAction)
	data.DefaultActionValue = stringOrNull(control.DefaultActionValue)
	data.Description = stringOrNull(control.Description)
	data.ModifiedBy = stringOrNull(control.ModifiedBy)
	data.ModifiedTime = stringOrNull(control.ModifiedTime)
	data.ParanoiaLevel = stringOrNull(control.ParanoiaLevel)
	data.ProtocolType = stringOrNull(control.ProtocolType)
	data.Severity = stringOrNull(control.Severity)
	data.Type = stringOrNull(control.Type)
	data.Version = stringOrNull(control.Version)
	data.Rules = rules
	data.AssociatedInspectionProfileNames = associated

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenInspectionCustomRules(ctx context.Context, rules []inspectioncustom.Rules) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"type":       types.StringType,
		"names":      types.SetType{ElemType: types.StringType},
		"conditions": types.ListType{ElemType: types.ObjectType{AttrTypes: ruleConditionAttrTypes()}},
	}

	if len(rules) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(rules))
	var diags diag.Diagnostics
	for _, rule := range rules {
		conditions, condDiags := flattenCustomRuleConditions(ctx, rule.Conditions)
		diags.Append(condDiags...)

		names := make([]string, 0, len(rule.Names))
		for _, name := range rule.Names {
			if trimmed := strings.TrimSpace(name); trimmed != "" {
				names = append(names, trimmed)
			}
		}
		namesSet, namesDiags := types.SetValueFrom(ctx, types.StringType, names)
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

func flattenCustomRuleConditions(ctx context.Context, conditions []inspectioncustom.Conditions) (types.List, diag.Diagnostics) {
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
