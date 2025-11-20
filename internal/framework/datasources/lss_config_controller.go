package datasources

import (
	"context"
	"fmt"
	"html"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/lssconfigcontroller"
)

var (
	_ datasource.DataSource              = &LSSConfigControllerDataSource{}
	_ datasource.DataSourceWithConfigure = &LSSConfigControllerDataSource{}
)

func NewLSSConfigControllerDataSource() datasource.DataSource {
	return &LSSConfigControllerDataSource{}
}

type LSSConfigControllerDataSource struct {
	client *client.Client
}

type LSSConfigControllerModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Config          types.List   `tfsdk:"config"`
	ConnectorGroups types.List   `tfsdk:"connector_groups"`
	PolicyRule      types.List   `tfsdk:"policy_rule"`
}

func (d *LSSConfigControllerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lss_config_controller"
}

func (d *LSSConfigControllerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves an LSS configuration by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the LSS configuration.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the LSS configuration.",
			},
		},
		Blocks: map[string]schema.Block{
			"config": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"audit_message": schema.StringAttribute{Computed: true},
						"description":   schema.StringAttribute{Computed: true},
						"enabled":       schema.BoolAttribute{Computed: true},
						"filter": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"format":           schema.StringAttribute{Computed: true},
						"id":               schema.StringAttribute{Computed: true},
						"name":             schema.StringAttribute{Computed: true},
						"lss_host":         schema.StringAttribute{Computed: true},
						"lss_port":         schema.StringAttribute{Computed: true},
						"source_log_type":  schema.StringAttribute{Computed: true},
						"use_tls":          schema.BoolAttribute{Computed: true},
						"microtenant_id":   schema.StringAttribute{Computed: true},
						"microtenant_name": schema.StringAttribute{Computed: true},
					},
				},
			},
			"connector_groups": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
					},
				},
			},
			"policy_rule": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"action":                      schema.StringAttribute{Computed: true},
						"action_id":                   schema.StringAttribute{Computed: true},
						"bypass_default_rule":         schema.BoolAttribute{Computed: true},
						"creation_time":               schema.StringAttribute{Computed: true},
						"custom_msg":                  schema.StringAttribute{Computed: true},
						"default_rule":                schema.BoolAttribute{Computed: true},
						"description":                 schema.StringAttribute{Computed: true},
						"id":                          schema.StringAttribute{Computed: true},
						"isolation_default_rule":      schema.BoolAttribute{Computed: true},
						"modified_by":                 schema.StringAttribute{Computed: true},
						"modified_time":               schema.StringAttribute{Computed: true},
						"name":                        schema.StringAttribute{Computed: true},
						"operator":                    schema.StringAttribute{Computed: true},
						"policy_set_id":               schema.StringAttribute{Computed: true},
						"policy_type":                 schema.StringAttribute{Computed: true},
						"priority":                    schema.StringAttribute{Computed: true},
						"reauth_default_rule":         schema.BoolAttribute{Computed: true},
						"reauth_idle_timeout":         schema.StringAttribute{Computed: true},
						"reauth_timeout":              schema.StringAttribute{Computed: true},
						"rule_order":                  schema.StringAttribute{Computed: true},
						"lss_default_rule":            schema.BoolAttribute{Computed: true},
						"zpn_cbi_profile_id":          schema.StringAttribute{Computed: true},
						"zpn_inspection_profile_id":   schema.StringAttribute{Computed: true},
						"zpn_inspection_profile_name": schema.StringAttribute{Computed: true},
						"microtenant_id":              schema.StringAttribute{Computed: true},
						"microtenant_name":            schema.StringAttribute{Computed: true},
					},
					Blocks: map[string]schema.Block{
						"conditions": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"creation_time": schema.StringAttribute{Computed: true},
									"id":            schema.StringAttribute{Computed: true},
									"modified_by":   schema.StringAttribute{Computed: true},
									"modified_time": schema.StringAttribute{Computed: true},
									"operator":      schema.StringAttribute{Computed: true},
									"negated":       schema.BoolAttribute{Computed: true},
								},
								Blocks: map[string]schema.Block{
									"operands": schema.ListNestedBlock{
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"creation_time": schema.StringAttribute{Computed: true},
												"id":            schema.StringAttribute{Computed: true},
												"idp_id":        schema.StringAttribute{Computed: true},
												"lhs":           schema.StringAttribute{Computed: true},
												"modified_by":   schema.StringAttribute{Computed: true},
												"modified_time": schema.StringAttribute{Computed: true},
												"name":          schema.StringAttribute{Computed: true},
												"object_type":   schema.StringAttribute{Computed: true},
												"rhs":           schema.StringAttribute{Computed: true},
												"operator":      schema.StringAttribute{Computed: true},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *LSSConfigControllerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LSSConfigControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data LSSConfigControllerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())
	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided.")
		return
	}

	var (
		resource *lssconfigcontroller.LSSResource
		err      error
	)

	if id != "" {
		tflog.Debug(ctx, "Retrieving LSS config controller by ID", map[string]any{"id": id})
		resource, _, err = lssconfigcontroller.Get(ctx, d.client.Service, id)
	} else {
		tflog.Debug(ctx, "Retrieving LSS config controller by name", map[string]any{"name": name})
		resource, _, err = lssconfigcontroller.GetByName(ctx, d.client.Service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read LSS config controller: %v", err))
		return
	}

	if resource == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("LSS config controller with id %q or name %q was not found.", id, name))
		return
	}

	configList, cfgDiags := flattenLSSConfig(ctx, resource.LSSConfig)
	resp.Diagnostics.Append(cfgDiags...)
	connectorGroups, cgDiags := flattenLSSConnectorGroups(ctx, resource.ConnectorGroups)
	resp.Diagnostics.Append(cgDiags...)
	policyRule, prDiags := flattenLSSPolicyRule(ctx, resource.PolicyRule)
	resp.Diagnostics.Append(prDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(resource.ID)
	if resource.LSSConfig != nil && strings.TrimSpace(resource.LSSConfig.Name) != "" {
		data.Name = types.StringValue(resource.LSSConfig.Name)
	} else if name != "" {
		data.Name = types.StringValue(name)
	} else {
		data.Name = types.StringNull()
	}
	data.Config = configList
	data.ConnectorGroups = connectorGroups
	data.PolicyRule = policyRule

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenLSSConfig(ctx context.Context, cfg *lssconfigcontroller.LSSConfig) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"audit_message":    types.StringType,
		"description":      types.StringType,
		"enabled":          types.BoolType,
		"filter":           types.ListType{ElemType: types.StringType},
		"format":           types.StringType,
		"id":               types.StringType,
		"name":             types.StringType,
		"lss_host":         types.StringType,
		"lss_port":         types.StringType,
		"source_log_type":  types.StringType,
		"use_tls":          types.BoolType,
		"microtenant_id":   types.StringType,
		"microtenant_name": types.StringType,
	}

	if cfg == nil {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	filterList, diags := types.ListValueFrom(ctx, types.StringType, cfg.Filter)

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"audit_message":    stringOrNull(html.UnescapeString(cfg.AuditMessage)),
		"description":      stringOrNull(cfg.Description),
		"enabled":          types.BoolValue(cfg.Enabled),
		"filter":           filterList,
		"format":           stringOrNull(html.UnescapeString(cfg.Format)),
		"id":               stringOrNull(cfg.ID),
		"name":             stringOrNull(cfg.Name),
		"lss_host":         stringOrNull(cfg.LSSHost),
		"lss_port":         stringOrNull(cfg.LSSPort),
		"source_log_type":  stringOrNull(cfg.SourceLogType),
		"use_tls":          types.BoolValue(cfg.UseTLS),
		"microtenant_id":   stringOrNull(cfg.MicroTenantID),
		"microtenant_name": stringOrNull(cfg.MicroTenantName),
	})
	diags.Append(objDiags...)

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	diags.Append(listDiags...)
	return list, diags
}

func flattenLSSConnectorGroups(ctx context.Context, groups []lssconfigcontroller.ConnectorGroups) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}

	if len(groups) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(groups))
	var diags diag.Diagnostics
	for _, g := range groups {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":   stringOrNull(g.ID),
			"name": stringOrNull(g.Name),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}

func flattenLSSPolicyRule(ctx context.Context, rule *lssconfigcontroller.PolicyRule) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"action":                      types.StringType,
		"action_id":                   types.StringType,
		"bypass_default_rule":         types.BoolType,
		"creation_time":               types.StringType,
		"custom_msg":                  types.StringType,
		"default_rule":                types.BoolType,
		"description":                 types.StringType,
		"id":                          types.StringType,
		"isolation_default_rule":      types.BoolType,
		"modified_by":                 types.StringType,
		"modified_time":               types.StringType,
		"name":                        types.StringType,
		"operator":                    types.StringType,
		"policy_set_id":               types.StringType,
		"policy_type":                 types.StringType,
		"priority":                    types.StringType,
		"reauth_default_rule":         types.BoolType,
		"reauth_idle_timeout":         types.StringType,
		"reauth_timeout":              types.StringType,
		"rule_order":                  types.StringType,
		"lss_default_rule":            types.BoolType,
		"zpn_cbi_profile_id":          types.StringType,
		"zpn_inspection_profile_id":   types.StringType,
		"zpn_inspection_profile_name": types.StringType,
		"microtenant_id":              types.StringType,
		"microtenant_name":            types.StringType,
		"conditions":                  types.ListType{ElemType: types.ObjectType{AttrTypes: lssConditionAttrTypes()}},
	}

	if rule == nil {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	conditions, condDiags := flattenLSSConditions(ctx, rule.Conditions)
	var diags diag.Diagnostics
	diags.Append(condDiags...)

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"action":                      stringOrNull(rule.Action),
		"action_id":                   stringOrNull(rule.ActionID),
		"bypass_default_rule":         types.BoolValue(rule.BypassDefaultRule),
		"creation_time":               stringOrNull(rule.CreationTime),
		"custom_msg":                  stringOrNull(rule.CustomMsg),
		"default_rule":                types.BoolValue(rule.DefaultRule),
		"description":                 stringOrNull(rule.Description),
		"id":                          stringOrNull(rule.ID),
		"isolation_default_rule":      types.BoolValue(rule.IsolationDefaultRule),
		"modified_by":                 stringOrNull(rule.ModifiedBy),
		"modified_time":               stringOrNull(rule.ModifiedTime),
		"name":                        stringOrNull(rule.Name),
		"operator":                    stringOrNull(rule.Operator),
		"policy_set_id":               stringOrNull(rule.PolicySetID),
		"policy_type":                 stringOrNull(rule.PolicyType),
		"priority":                    stringOrNull(rule.Priority),
		"reauth_default_rule":         types.BoolValue(rule.ReauthDefaultRule),
		"reauth_idle_timeout":         stringOrNull(rule.ReauthIdleTimeout),
		"reauth_timeout":              stringOrNull(rule.ReauthTimeout),
		"rule_order":                  stringOrNull(rule.RuleOrder),
		"lss_default_rule":            types.BoolValue(rule.LssDefaultRule),
		"zpn_cbi_profile_id":          stringOrNull(rule.ZpnCbiProfileID),
		"zpn_inspection_profile_id":   stringOrNull(rule.ZpnInspectionProfileID),
		"zpn_inspection_profile_name": stringOrNull(rule.ZpnInspectionProfileName),
		"microtenant_id":              stringOrNull(rule.MicroTenantID),
		"microtenant_name":            stringOrNull(rule.MicroTenantName),
		"conditions":                  conditions,
	})
	diags.Append(objDiags...)

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	diags.Append(listDiags...)
	return list, diags
}

func lssConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"creation_time": types.StringType,
		"id":            types.StringType,
		"modified_by":   types.StringType,
		"modified_time": types.StringType,
		"operator":      types.StringType,
		"negated":       types.BoolType,
		"operands":      types.ListType{ElemType: types.ObjectType{AttrTypes: lssOperandAttrTypes()}},
	}
}

func lssOperandAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"creation_time": types.StringType,
		"id":            types.StringType,
		"idp_id":        types.StringType,
		"lhs":           types.StringType,
		"modified_by":   types.StringType,
		"modified_time": types.StringType,
		"name":          types.StringType,
		"object_type":   types.StringType,
		"rhs":           types.StringType,
		"operator":      types.StringType,
	}
}

func flattenLSSConditions(ctx context.Context, conditions []lssconfigcontroller.Conditions) (types.List, diag.Diagnostics) {
	attrTypes := lssConditionAttrTypes()

	if len(conditions) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(conditions))
	var diags diag.Diagnostics

	for _, condition := range conditions {
		operands, opDiags := flattenLSSOperands(ctx, condition.Operands)
		diags.Append(opDiags...)

		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"creation_time": stringOrNull(condition.CreationTime),
			"id":            stringOrNull(condition.ID),
			"modified_by":   stringOrNull(condition.ModifiedBy),
			"modified_time": stringOrNull(condition.ModifiedTime),
			"operator":      stringOrNull(condition.Operator),
			"negated":       types.BoolValue(condition.Negated),
			"operands":      operands,
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}

func flattenLSSOperands(ctx context.Context, operands *[]lssconfigcontroller.Operands) (types.List, diag.Diagnostics) {
	attrTypes := lssOperandAttrTypes()

	if operands == nil || len(*operands) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(*operands))
	var diags diag.Diagnostics

	for _, operand := range *operands {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"creation_time": stringOrNull(operand.CreationTime),
			"id":            stringOrNull(operand.ID),
			"idp_id":        stringOrNull(operand.IdpID),
			"lhs":           stringOrNull(operand.LHS),
			"modified_by":   stringOrNull(operand.ModifiedBy),
			"modified_time": stringOrNull(operand.ModifiedTime),
			"name":          stringOrNull(operand.Name),
			"object_type":   stringOrNull(operand.ObjectType),
			"rhs":           stringOrNull(operand.RHS),
			"operator":      types.StringNull(), // Not returned by API, but in SDKv2 schema
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}

func stringOrNull(value string) types.String {
	if strings.TrimSpace(value) == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}
