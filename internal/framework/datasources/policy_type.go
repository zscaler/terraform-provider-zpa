package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

var (
	_ datasource.DataSource              = &PolicyTypeDataSource{}
	_ datasource.DataSourceWithConfigure = &PolicyTypeDataSource{}
)

func NewPolicyTypeDataSource() datasource.DataSource {
	return &PolicyTypeDataSource{}
}

type PolicyTypeDataSource struct {
	client *client.Client
}

type PolicyTypeModel struct {
	ID              types.String `tfsdk:"id"`
	CreationTime    types.String `tfsdk:"creation_time"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
	ModifiedTime    types.String `tfsdk:"modified_time"`
	Name            types.String `tfsdk:"name"`
	Sorted          types.Bool   `tfsdk:"sorted"`
	PolicyType      types.String `tfsdk:"policy_type"`
	MicroTenantID   types.String `tfsdk:"microtenant_id"`
	MicroTenantName types.String `tfsdk:"microtenant_name"`
	Rules           types.List   `tfsdk:"rules"`
}

func (d *PolicyTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_type"
}

func (d *PolicyTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA policy set by policy type.",
		Attributes: map[string]schema.Attribute{
			"policy_type": schema.StringAttribute{
				Optional:    true,
				Description: "Policy type to retrieve. Defaults to GLOBAL_POLICY when unspecified.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"ACCESS_POLICY", "GLOBAL_POLICY", "TIMEOUT_POLICY", "REAUTH_POLICY",
						"CLIENT_FORWARDING_POLICY", "BYPASS_POLICY", "ISOLATION_POLICY",
						"INSPECTION_POLICY", "SIEM_POLICY", "CREDENTIAL_POLICY",
						"CAPABILITIES_POLICY", "REDIRECTION_POLICY",
					),
				},
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID used to scope the lookup.",
			},
			"id":               schema.StringAttribute{Computed: true},
			"creation_time":    schema.StringAttribute{Computed: true},
			"description":      schema.StringAttribute{Computed: true},
			"enabled":          schema.BoolAttribute{Computed: true},
			"modified_by":      schema.StringAttribute{Computed: true},
			"modified_time":    schema.StringAttribute{Computed: true},
			"name":             schema.StringAttribute{Computed: true},
			"sorted":           schema.BoolAttribute{Computed: true},
			"microtenant_name": schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			"rules": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":                        schema.StringAttribute{Computed: true},
						"name":                      schema.StringAttribute{Computed: true},
						"description":               schema.StringAttribute{Computed: true},
						"action":                    schema.StringAttribute{Computed: true},
						"action_id":                 schema.StringAttribute{Computed: true},
						"bypass_default_rule":       schema.BoolAttribute{Computed: true},
						"creation_time":             schema.StringAttribute{Computed: true},
						"custom_msg":                schema.StringAttribute{Computed: true},
						"modified_by":               schema.StringAttribute{Computed: true},
						"modified_time":             schema.StringAttribute{Computed: true},
						"operator":                  schema.StringAttribute{Computed: true},
						"policy_set_id":             schema.StringAttribute{Computed: true},
						"policy_type":               schema.StringAttribute{Computed: true},
						"priority":                  schema.StringAttribute{Computed: true},
						"reauth_default_rule":       schema.BoolAttribute{Computed: true},
						"reauth_idle_timeout":       schema.StringAttribute{Computed: true},
						"reauth_timeout":            schema.StringAttribute{Computed: true},
						"rule_order":                schema.StringAttribute{Computed: true},
						"zpn_cbi_profile_id":        schema.StringAttribute{Computed: true},
						"zpn_isolation_profile_id":  schema.StringAttribute{Computed: true},
						"zpn_inspection_profile_id": schema.StringAttribute{Computed: true},
						"microtenant_id":            schema.StringAttribute{Computed: true},
						"microtenant_name":          schema.StringAttribute{Computed: true},
					},
					Blocks: map[string]schema.Block{
						"conditions": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"creation_time":  schema.StringAttribute{Computed: true},
									"id":             schema.StringAttribute{Computed: true},
									"modified_by":    schema.StringAttribute{Computed: true},
									"modified_time":  schema.StringAttribute{Computed: true},
									"operator":       schema.StringAttribute{Computed: true},
									"microtenant_id": schema.StringAttribute{Computed: true},
								},
								Blocks: map[string]schema.Block{
									"operands": schema.ListNestedBlock{
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"creation_time":  schema.StringAttribute{Computed: true},
												"id":             schema.StringAttribute{Computed: true},
												"idp_id":         schema.StringAttribute{Computed: true},
												"lhs":            schema.StringAttribute{Computed: true},
												"modified_by":    schema.StringAttribute{Computed: true},
												"modified_time":  schema.StringAttribute{Computed: true},
												"name":           schema.StringAttribute{Computed: true},
												"object_type":    schema.StringAttribute{Computed: true},
												"rhs":            schema.StringAttribute{Computed: true},
												"microtenant_id": schema.StringAttribute{Computed: true},
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

func (d *PolicyTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PolicyTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data PolicyTypeModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := d.client.Service
	if !data.MicroTenantID.IsNull() && !data.MicroTenantID.IsUnknown() {
		if microID := strings.TrimSpace(data.MicroTenantID.ValueString()); microID != "" {
			service = service.WithMicroTenant(microID)
			data.MicroTenantID = types.StringValue(microID)
		}
	}

	policyType := "GLOBAL_POLICY"
	if !data.PolicyType.IsNull() && !data.PolicyType.IsUnknown() {
		policyType = strings.TrimSpace(data.PolicyType.ValueString())
	}

	tflog.Debug(ctx, "Retrieving policy set by type", map[string]any{"policy_type": policyType})
	policySet, _, err := policysetcontroller.GetByPolicyType(ctx, service, policyType)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read policy type %q: %v", policyType, err))
		return
	}

	if policySet == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Policy type %q not found", policyType))
		return
	}

	state, diags := flattenPolicySet(ctx, policySet)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.MicroTenantID.IsNull() && strings.TrimSpace(data.MicroTenantID.ValueString()) != "" {
		state.MicroTenantID = data.MicroTenantID
	} else if policySet.MicroTenantID != "" {
		state.MicroTenantID = types.StringValue(policySet.MicroTenantID)
	} else {
		state.MicroTenantID = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenPolicySet(ctx context.Context, ps *policysetcontroller.PolicySet) (PolicyTypeModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	rules, ruleDiags := flattenPolicyRules(ctx, ps.Rules)
	diags.Append(ruleDiags...)

	state := PolicyTypeModel{
		ID:              types.StringValue(ps.ID),
		CreationTime:    types.StringValue(ps.CreationTime),
		Description:     types.StringValue(ps.Description),
		Enabled:         types.BoolValue(ps.Enabled),
		ModifiedBy:      types.StringValue(ps.ModifiedBy),
		ModifiedTime:    types.StringValue(ps.ModifiedTime),
		Name:            types.StringValue(ps.Name),
		Sorted:          types.BoolValue(ps.Sorted),
		PolicyType:      types.StringValue(ps.PolicyType),
		MicroTenantName: types.StringValue(ps.MicroTenantName),
		Rules:           rules,
	}

	return state, diags
}

func flattenPolicyRules(ctx context.Context, rules []policysetcontroller.PolicyRule) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	operandAttrTypes := map[string]attr.Type{
		"creation_time":  types.StringType,
		"id":             types.StringType,
		"idp_id":         types.StringType,
		"lhs":            types.StringType,
		"modified_by":    types.StringType,
		"modified_time":  types.StringType,
		"name":           types.StringType,
		"object_type":    types.StringType,
		"rhs":            types.StringType,
		"microtenant_id": types.StringType,
	}

	conditionAttrTypes := map[string]attr.Type{
		"creation_time":  types.StringType,
		"id":             types.StringType,
		"modified_by":    types.StringType,
		"modified_time":  types.StringType,
		"operator":       types.StringType,
		"microtenant_id": types.StringType,
		"operands":       types.ListType{ElemType: types.ObjectType{AttrTypes: operandAttrTypes}},
	}

	ruleAttrTypes := map[string]attr.Type{
		"id":                        types.StringType,
		"name":                      types.StringType,
		"description":               types.StringType,
		"action":                    types.StringType,
		"action_id":                 types.StringType,
		"bypass_default_rule":       types.BoolType,
		"creation_time":             types.StringType,
		"custom_msg":                types.StringType,
		"modified_by":               types.StringType,
		"modified_time":             types.StringType,
		"operator":                  types.StringType,
		"policy_set_id":             types.StringType,
		"policy_type":               types.StringType,
		"priority":                  types.StringType,
		"reauth_default_rule":       types.BoolType,
		"reauth_idle_timeout":       types.StringType,
		"reauth_timeout":            types.StringType,
		"rule_order":                types.StringType,
		"zpn_cbi_profile_id":        types.StringType,
		"zpn_isolation_profile_id":  types.StringType,
		"zpn_inspection_profile_id": types.StringType,
		"microtenant_id":            types.StringType,
		"microtenant_name":          types.StringType,
		"conditions":                types.ListType{ElemType: types.ObjectType{AttrTypes: conditionAttrTypes}},
	}

	if len(rules) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: ruleAttrTypes}), diags
	}

	values := make([]attr.Value, 0, len(rules))
	for _, rule := range rules {
		conditions, condDiags := flattenPolicyConditions(ctx, rule.Conditions, conditionAttrTypes, operandAttrTypes)
		diags.Append(condDiags...)

		obj, objDiags := types.ObjectValue(ruleAttrTypes, map[string]attr.Value{
			"id":                        types.StringValue(rule.ID),
			"name":                      types.StringValue(rule.Name),
			"description":               types.StringValue(rule.Description),
			"action":                    types.StringValue(rule.Action),
			"action_id":                 types.StringValue(rule.ActionID),
			"bypass_default_rule":       types.BoolValue(rule.BypassDefaultRule),
			"creation_time":             types.StringValue(rule.CreationTime),
			"custom_msg":                types.StringValue(rule.CustomMsg),
			"modified_by":               types.StringValue(rule.ModifiedBy),
			"modified_time":             types.StringValue(rule.ModifiedTime),
			"operator":                  types.StringValue(rule.Operator),
			"policy_set_id":             types.StringValue(rule.PolicySetID),
			"policy_type":               types.StringValue(rule.PolicyType),
			"priority":                  types.StringValue(rule.Priority),
			"reauth_default_rule":       types.BoolValue(rule.ReauthDefaultRule),
			"reauth_idle_timeout":       types.StringValue(rule.ReauthIdleTimeout),
			"reauth_timeout":            types.StringValue(rule.ReauthTimeout),
			"rule_order":                types.StringValue(rule.RuleOrder),
			"zpn_cbi_profile_id":        types.StringValue(rule.ZpnCbiProfileID),
			"zpn_isolation_profile_id":  types.StringValue(rule.ZpnIsolationProfileID),
			"zpn_inspection_profile_id": types.StringValue(rule.ZpnInspectionProfileID),
			"microtenant_id":            types.StringValue(rule.MicroTenantID),
			"microtenant_name":          types.StringValue(rule.MicroTenantName),
			"conditions":                conditions,
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: ruleAttrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}

func flattenPolicyConditions(ctx context.Context, conditions []policysetcontroller.Conditions, conditionAttrTypes, operandAttrTypes map[string]attr.Type) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(conditions) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: conditionAttrTypes}), diags
	}

	values := make([]attr.Value, 0, len(conditions))
	for _, condition := range conditions {
		operands, opDiags := flattenPolicyOperands(ctx, condition.Operands, operandAttrTypes)
		diags.Append(opDiags...)

		obj, objDiags := types.ObjectValue(conditionAttrTypes, map[string]attr.Value{
			"creation_time":  types.StringValue(condition.CreationTime),
			"id":             types.StringValue(condition.ID),
			"modified_by":    types.StringValue(condition.ModifiedBy),
			"modified_time":  types.StringValue(condition.ModifiedTime),
			"operator":       types.StringValue(condition.Operator),
			"microtenant_id": types.StringValue(condition.MicroTenantID),
			"operands":       operands,
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: conditionAttrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}

func flattenPolicyOperands(ctx context.Context, operands []policysetcontroller.Operands, operandAttrTypes map[string]attr.Type) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(operands) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: operandAttrTypes}), diags
	}

	values := make([]attr.Value, 0, len(operands))
	for _, operand := range operands {
		obj, objDiags := types.ObjectValue(operandAttrTypes, map[string]attr.Value{
			"creation_time":  types.StringValue(operand.CreationTime),
			"id":             types.StringValue(operand.ID),
			"idp_id":         types.StringValue(operand.IdpID),
			"lhs":            types.StringValue(operand.LHS),
			"modified_by":    types.StringValue(operand.ModifiedBy),
			"modified_time":  types.StringValue(operand.ModifiedTime),
			"name":           types.StringValue(operand.Name),
			"object_type":    types.StringValue(operand.ObjectType),
			"rhs":            types.StringValue(operand.RHS),
			"microtenant_id": types.StringValue(operand.MicroTenantID),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: operandAttrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}
