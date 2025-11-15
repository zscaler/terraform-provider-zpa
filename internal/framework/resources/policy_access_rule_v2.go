package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/fabiotavarespr/iso3166"
	stringvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	fwrschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	fwplanmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	fwstringplanmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	fwvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

var (
	_ resource.Resource                = &PolicyAccessRuleV2Resource{}
	_ resource.ResourceWithConfigure   = &PolicyAccessRuleV2Resource{}
	_ resource.ResourceWithImportState = &PolicyAccessRuleV2Resource{}
)

func NewPolicyAccessRuleV2Resource() resource.Resource {
	return &PolicyAccessRuleV2Resource{}
}

type PolicyAccessRuleV2Resource struct {
	client *client.Client
}

type PolicyAccessRuleV2Model struct {
	ID                 types.String                 `tfsdk:"id"`
	Name               types.String                 `tfsdk:"name"`
	Description        types.String                 `tfsdk:"description"`
	Action             types.String                 `tfsdk:"action"`
	Operator           types.String                 `tfsdk:"operator"`
	PolicySetID        types.String                 `tfsdk:"policy_set_id"`
	CustomMsg          types.String                 `tfsdk:"custom_msg"`
	AppServerGroups    types.List                   `tfsdk:"app_server_groups"`
	AppConnectorGroups types.List                   `tfsdk:"app_connector_groups"`
	ExtranetEnabled    types.Bool                   `tfsdk:"extranet_enabled"`
	ExtranetDTO        types.List                   `tfsdk:"extranet_dto"`
	Conditions         []PolicyAccessConditionModel `tfsdk:"conditions"`
	MicrotenantID      types.String                 `tfsdk:"microtenant_id"`
}

type PolicyAccessConditionModel struct {
	ID       types.String               `tfsdk:"id"`
	Operator types.String               `tfsdk:"operator"`
	Operands []PolicyAccessOperandModel `tfsdk:"operands"`
}

type PolicyAccessOperandModel struct {
	Values      types.Set    `tfsdk:"values"`
	ObjectType  types.String `tfsdk:"object_type"`
	EntryValues types.Set    `tfsdk:"entry_values"`
}

type PolicyAccessEntryValueModel struct {
	RHS types.String `tfsdk:"rhs"`
	LHS types.String `tfsdk:"lhs"`
}

type PolicyAccessExtranetModel struct {
	ZPNErID          types.String `tfsdk:"zpn_er_id"`
	LocationDTO      types.List   `tfsdk:"location_dto"`
	LocationGroupDTO types.List   `tfsdk:"location_group_dto"`
}

type PolicyAccessExtranetLocationModel struct {
	ID types.String `tfsdk:"id"`
}

var (
	entryValueAttrTypes = map[string]attr.Type{
		"lhs": types.StringType,
		"rhs": types.StringType,
	}
	extranetLocationAttrTypes = map[string]attr.Type{
		"id": types.StringType,
	}
	extranetAttrTypes = map[string]attr.Type{
		"zpn_er_id":          types.StringType,
		"location_dto":       types.ListType{ElemType: types.ObjectType{AttrTypes: extranetLocationAttrTypes}},
		"location_group_dto": types.ListType{ElemType: types.ObjectType{AttrTypes: extranetLocationAttrTypes}},
	}
)

func (r *PolicyAccessRuleV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_access_rule_v2"
}

func (r *PolicyAccessRuleV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := map[string]fwrschema.Attribute{
		"id": fwrschema.StringAttribute{
			Computed: true,
			PlanModifiers: []fwplanmodifier.String{
				fwstringplanmodifier.UseStateForUnknown(),
			},
		},
		"name":        fwrschema.StringAttribute{Required: true},
		"description": fwrschema.StringAttribute{Optional: true},
		"action": fwrschema.StringAttribute{
			Optional: true,
			Validators: []fwvalidator.String{
				stringvalidator.OneOf("ALLOW", "DENY", "REQUIRE_APPROVAL"),
			},
		},
		"operator": fwrschema.StringAttribute{
			Optional: true,
			Computed: true,
			Validators: []fwvalidator.String{
				stringvalidator.OneOf("AND", "OR"),
			},
		},
		"policy_set_id": fwrschema.StringAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []fwplanmodifier.String{
				fwstringplanmodifier.UseStateForUnknown(),
			},
		},
		"custom_msg": fwrschema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"microtenant_id": fwrschema.StringAttribute{
			Optional: true,
			PlanModifiers: []fwplanmodifier.String{
				fwstringplanmodifier.UseStateForUnknown(),
			},
		},
		"extranet_enabled": fwrschema.BoolAttribute{
			Optional: true,
		},
	}

	objectTypes := []string{
		"APP",
		"APP_GROUP",
		"LOCATION",
		"IDP",
		"SAML",
		"SCIM",
		"SCIM_GROUP",
		"CLIENT_TYPE",
		"POSTURE",
		"TRUSTED_NETWORK",
		"BRANCH_CONNECTOR_GROUP",
		"EDGE_CONNECTOR_GROUP",
		"MACHINE_GRP",
		"COUNTRY_CODE",
		"PLATFORM",
		"RISK_FACTOR_TYPE",
		"CHROME_ENTERPRISE",
		"CHROME_POSTURE_PROFILE",
		"USER_PORTAL",
	}

	resp.Schema = fwrschema.Schema{
		Description: "Manages ZPA Access Policy rules (v2).",
		Attributes:  attributes,
		Blocks: map[string]fwrschema.Block{
			"conditions": fwrschema.SetNestedBlock{
				NestedObject: fwrschema.NestedBlockObject{
					Attributes: map[string]fwrschema.Attribute{
						"id": fwrschema.StringAttribute{
							Computed: true,
						},
						"operator": fwrschema.StringAttribute{
							Optional: true,
							Computed: true,
							Validators: []fwvalidator.String{
								stringvalidator.OneOf("AND", "OR"),
							},
						},
					},
					Blocks: map[string]fwrschema.Block{
						"operands": fwrschema.SetNestedBlock{
							NestedObject: fwrschema.NestedBlockObject{
								Attributes: map[string]fwrschema.Attribute{
									"values": fwrschema.SetAttribute{
										ElementType: types.StringType,
										Optional:    true,
									},
									"object_type": fwrschema.StringAttribute{
										Optional: true,
										Validators: []fwvalidator.String{
											stringvalidator.OneOf(objectTypes...),
										},
									},
								},
								Blocks: map[string]fwrschema.Block{
									"entry_values": fwrschema.SetNestedBlock{
										NestedObject: fwrschema.NestedBlockObject{
											Attributes: map[string]fwrschema.Attribute{
												"rhs": fwrschema.StringAttribute{Optional: true},
												"lhs": fwrschema.StringAttribute{Optional: true},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"app_server_groups": fwrschema.ListNestedBlock{
				NestedObject: fwrschema.NestedBlockObject{
					Attributes: map[string]fwrschema.Attribute{
						"id": fwrschema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
			"app_connector_groups": fwrschema.ListNestedBlock{
				NestedObject: fwrschema.NestedBlockObject{
					Attributes: map[string]fwrschema.Attribute{
						"id": fwrschema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
			"extranet_dto": fwrschema.ListNestedBlock{
				NestedObject: fwrschema.NestedBlockObject{
					Attributes: map[string]fwrschema.Attribute{
						"zpn_er_id": fwrschema.StringAttribute{
							Optional: true,
						},
					},
					Blocks: map[string]fwrschema.Block{
						"location_dto": fwrschema.ListNestedBlock{
							NestedObject: fwrschema.NestedBlockObject{
								Attributes: map[string]fwrschema.Attribute{
									"id": fwrschema.StringAttribute{
										Required: true,
									},
								},
							},
						},
						"location_group_dto": fwrschema.ListNestedBlock{
							NestedObject: fwrschema.NestedBlockObject{
								Attributes: map[string]fwrschema.Attribute{
									"id": fwrschema.StringAttribute{
										Required: true,
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

func (r *PolicyAccessRuleV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got %T", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *PolicyAccessRuleV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy access rules.")
		return
	}

	var plan PolicyAccessRuleV2Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateObjectTypeUniquenessV2(ctx, plan.Conditions)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(validatePolicyRuleConditionsV2(ctx, plan.Conditions)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)
	helperClient := helpers.NewHelperClient(r.client)

	policySetID := helpers.StringValue(plan.PolicySetID)
	microTenantID := helpers.StringValue(plan.MicrotenantID)
	if policySetID == "" {
		var err error
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeAccess, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	request, diags := expandPolicyAccessRuleV2(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := policysetcontrollerv2.CreateRule(ctx, service, request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create policy access rule: %v", err))
		return
	}

	state, readDiags := r.readPolicyAccessRuleV2(ctx, service, policySetID, created.ID, plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PolicyAccessRuleV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy access rules.")
		return
	}

	var state PolicyAccessRuleV2Model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)
	helperClient := helpers.NewHelperClient(r.client)

	policySetID := helpers.StringValue(state.PolicySetID)
	microTenantID := helpers.StringValue(state.MicrotenantID)
	if policySetID == "" {
		var err error
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeAccess, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	newState, diags := r.readPolicyAccessRuleV2(ctx, service, policySetID, state.ID.ValueString(), state.MicrotenantID)
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

	newState.PolicySetID = types.StringValue(policySetID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *PolicyAccessRuleV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy access rules.")
		return
	}

	var plan PolicyAccessRuleV2Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	resp.Diagnostics.Append(validateObjectTypeUniquenessV2(ctx, plan.Conditions)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(validatePolicyRuleConditionsV2(ctx, plan.Conditions)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)
	helperClient := helpers.NewHelperClient(r.client)

	policySetID := helpers.StringValue(plan.PolicySetID)
	microTenantID := helpers.StringValue(plan.MicrotenantID)
	if policySetID == "" {
		var err error
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeAccess, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	request, diags := expandPolicyAccessRuleV2(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := policysetcontrollerv2.UpdateRule(ctx, service, policySetID, plan.ID.ValueString(), request); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update policy access rule: %v", err))
		return
	}

	newState, readDiags := r.readPolicyAccessRuleV2(ctx, service, policySetID, plan.ID.ValueString(), plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *PolicyAccessRuleV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy access rules.")
		return
	}

	var state PolicyAccessRuleV2Model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)
	helperClient := helpers.NewHelperClient(r.client)

	policySetID := helpers.StringValue(state.PolicySetID)
	microTenantID := helpers.StringValue(state.MicrotenantID)
	if policySetID == "" {
		var err error
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeAccess, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	if _, err := policysetcontrollerv2.Delete(ctx, service, policySetID, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete policy access rule: %v", err))
		return
	}
}

func (r *PolicyAccessRuleV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before importing policy access rules.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the policy access rule ID or name.")
		return
	}

	if _, err := fmt.Sscan(id, new(int64)); err != nil {
		rule, _, lookupErr := policysetcontrollerv2.GetByNameAndTypes(ctx, r.client.Service, []string{helpers.PolicyTypeAccess, helpers.PolicyTypeGlobal}, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate policy access rule %q: %v", id, lookupErr))
			return
		}
		id = rule.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *PolicyAccessRuleV2Resource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && microtenantID.ValueString() != "" {
		service = service.WithMicroTenant(microtenantID.ValueString())
	}
	return service
}

func (r *PolicyAccessRuleV2Resource) readPolicyAccessRuleV2(ctx context.Context, service *zscaler.Service, policySetID, ruleID string, microTenantID types.String) (PolicyAccessRuleV2Model, diag.Diagnostics) {
	ruleResource, _, err := policysetcontrollerv2.GetPolicyRule(ctx, service, policySetID, ruleID)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return PolicyAccessRuleV2Model{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Policy access rule %s not found", ruleID))}
		}
		return PolicyAccessRuleV2Model{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read policy access rule: %v", err))}
	}

	rule := helpers.ConvertV1ResponseToV2Request(*ruleResource)
	return flattenPolicyAccessRuleV2(ctx, rule, policySetID, microTenantID)
}

func expandPolicyAccessRuleV2(ctx context.Context, model *PolicyAccessRuleV2Model, policySetID string) (*policysetcontrollerv2.PolicyRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	serverGroups, sgDiags := helpers.ExpandServerGroups(ctx, model.AppServerGroups)
	diags.Append(sgDiags...)

	connectorGroups, cgDiags := helpers.ExpandAppConnectorGroups(ctx, model.AppConnectorGroups)
	diags.Append(cgDiags...)

	conditions, condDiags := expandPolicyRuleConditionsV2(ctx, model.Conditions)
	diags.Append(condDiags...)

	extranet, extranetDiags := expandExtranetDTOModel(ctx, model.ExtranetDTO)
	diags.Append(extranetDiags...)

	if diags.HasError() {
		return nil, diags
	}

	rule := &policysetcontrollerv2.PolicyRule{
		ID:                 helpers.StringValue(model.ID),
		Name:               helpers.StringValue(model.Name),
		Description:        helpers.StringValue(model.Description),
		Action:             helpers.StringValue(model.Action),
		Operator:           helpers.StringValue(model.Operator),
		PolicySetID:        policySetID,
		CustomMsg:          helpers.StringValue(model.CustomMsg),
		MicroTenantID:      helpers.StringValue(model.MicrotenantID),
		ExtranetEnabled:    helpers.BoolValue(model.ExtranetEnabled, false),
		AppServerGroups:    serverGroups,
		AppConnectorGroups: connectorGroups,
		Conditions:         conditions,
	}

	if extranet != nil {
		rule.ExtranetDTO = *extranet
	}

	return rule, diags
}

func flattenPolicyAccessRuleV2(ctx context.Context, rule policysetcontrollerv2.PolicyRule, policySetID string, microTenantID types.String) (PolicyAccessRuleV2Model, diag.Diagnostics) {
	var diags diag.Diagnostics

	serverGroups, sgDiags := helpers.FlattenServerGroups(ctx, rule.AppServerGroups)
	diags.Append(sgDiags...)

	connectorGroups, cgDiags := helpers.FlattenAppConnectorGroups(ctx, rule.AppConnectorGroups)
	diags.Append(cgDiags...)

	conditions, condDiags := flattenPolicyRuleConditionsV2(ctx, rule.Conditions)
	diags.Append(condDiags...)

	extranet, extranetDiags := flattenExtranetDTOModel(ctx, rule.ExtranetDTO)
	diags.Append(extranetDiags...)

	model := PolicyAccessRuleV2Model{
		ID:                 types.StringValue(rule.ID),
		Name:               types.StringValue(rule.Name),
		Description:        types.StringValue(rule.Description),
		Action:             types.StringValue(rule.Action),
		Operator:           types.StringValue(rule.Operator),
		PolicySetID:        types.StringValue(policySetID),
		CustomMsg:          types.StringValue(rule.CustomMsg),
		AppServerGroups:    serverGroups,
		AppConnectorGroups: connectorGroups,
		ExtranetEnabled:    types.BoolValue(rule.ExtranetEnabled),
		ExtranetDTO:        extranet,
		Conditions:         conditions,
		MicrotenantID:      microTenantID,
	}

	return model, diags
}

func expandPolicyRuleConditionsV2(ctx context.Context, models []PolicyAccessConditionModel) ([]policysetcontrollerv2.PolicyRuleResourceConditions, diag.Diagnostics) {
	result := make([]policysetcontrollerv2.PolicyRuleResourceConditions, 0, len(models))
	var diags diag.Diagnostics

	for _, condition := range models {
		operandModels := make([]policysetcontrollerv2.PolicyRuleResourceOperands, 0, len(condition.Operands))

		for _, operand := range condition.Operands {
			values, valueDiags := helpers.SetValueToStringSlice(ctx, operand.Values)
			diags.Append(valueDiags...)
			if diags.HasError() {
				return nil, diags
			}

			entryModels, entryDiags := extractEntryValueModels(ctx, operand.EntryValues)
			diags.Append(entryDiags...)
			if diags.HasError() {
				return nil, diags
			}

			entryValues := make([]policysetcontrollerv2.OperandsResourceLHSRHSValue, 0, len(entryModels))
			for _, entry := range entryModels {
				entryValues = append(entryValues, policysetcontrollerv2.OperandsResourceLHSRHSValue{
					LHS: helpers.StringValue(entry.LHS),
					RHS: helpers.StringValue(entry.RHS),
				})
			}

			operandModels = append(operandModels, policysetcontrollerv2.PolicyRuleResourceOperands{
				ObjectType:        helpers.StringValue(operand.ObjectType),
				Values:            values,
				EntryValuesLHSRHS: entryValues,
			})
		}

		result = append(result, policysetcontrollerv2.PolicyRuleResourceConditions{
			ID:       helpers.StringValue(condition.ID),
			Operator: helpers.StringValue(condition.Operator),
			Operands: operandModels,
		})
	}

	return result, diags
}

func flattenPolicyRuleConditionsV2(ctx context.Context, conditions []policysetcontrollerv2.PolicyRuleResourceConditions) ([]PolicyAccessConditionModel, diag.Diagnostics) {
	result := make([]PolicyAccessConditionModel, 0, len(conditions))
	var diags diag.Diagnostics

	for _, condition := range conditions {
		operands := make([]PolicyAccessOperandModel, 0, len(condition.Operands))

		for _, operand := range condition.Operands {
			valuesSet := types.SetNull(types.StringType)
			if len(operand.Values) > 0 {
				setValue, setDiags := types.SetValueFrom(ctx, types.StringType, operand.Values)
				diags.Append(setDiags...)
				valuesSet = setValue
			}

			entrySet := types.SetNull(types.ObjectType{AttrTypes: entryValueAttrTypes})
			if len(operand.EntryValuesLHSRHS) > 0 {
				entryValues := make([]attr.Value, 0, len(operand.EntryValuesLHSRHS))
				for _, entry := range operand.EntryValuesLHSRHS {
					objValue, objDiags := types.ObjectValue(entryValueAttrTypes, map[string]attr.Value{
						"lhs": helpers.StringValueOrNull(entry.LHS),
						"rhs": helpers.StringValueOrNull(entry.RHS),
					})
					diags.Append(objDiags...)
					entryValues = append(entryValues, objValue)
				}
				setValue, setDiags := types.SetValue(types.ObjectType{AttrTypes: entryValueAttrTypes}, entryValues)
				diags.Append(setDiags...)
				entrySet = setValue
			}

			operands = append(operands, PolicyAccessOperandModel{
				ObjectType:  types.StringValue(operand.ObjectType),
				Values:      valuesSet,
				EntryValues: entrySet,
			})
		}

		result = append(result, PolicyAccessConditionModel{
			ID:       types.StringValue(condition.ID),
			Operator: types.StringValue(condition.Operator),
			Operands: operands,
		})
	}

	return result, diags
}

func expandExtranetDTOModel(ctx context.Context, list types.List) (*common.ExtranetDTO, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, diag.Diagnostics{}
	}

	var models []PolicyAccessExtranetModel
	var diags diag.Diagnostics
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() || len(models) == 0 {
		return nil, diags
	}

	model := models[0]
	dto := &common.ExtranetDTO{
		ZpnErID: helpers.StringValue(model.ZPNErID),
	}

	if !model.LocationDTO.IsNull() && !model.LocationDTO.IsUnknown() {
		var locations []PolicyAccessExtranetLocationModel
		diags.Append(model.LocationDTO.ElementsAs(ctx, &locations, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for _, loc := range locations {
			id := helpers.StringValue(loc.ID)
			if id != "" {
				dto.LocationDTO = append(dto.LocationDTO, common.LocationDTO{ID: id})
			}
		}
	}

	if !model.LocationGroupDTO.IsNull() && !model.LocationGroupDTO.IsUnknown() {
		var groups []PolicyAccessExtranetLocationModel
		diags.Append(model.LocationGroupDTO.ElementsAs(ctx, &groups, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for _, group := range groups {
			id := helpers.StringValue(group.ID)
			if id != "" {
				dto.LocationGroupDTO = append(dto.LocationGroupDTO, common.LocationGroupDTO{ID: id})
			}
		}
	}

	return dto, diags
}

func flattenExtranetDTOModel(ctx context.Context, dto common.ExtranetDTO) (types.List, diag.Diagnostics) {
	if dto.ZpnErID == "" && len(dto.LocationDTO) == 0 && len(dto.LocationGroupDTO) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: extranetAttrTypes}), diag.Diagnostics{}
	}

	var diags diag.Diagnostics

	locationValues := make([]attr.Value, 0, len(dto.LocationDTO))
	for _, loc := range dto.LocationDTO {
		objValue, objDiags := types.ObjectValue(extranetLocationAttrTypes, map[string]attr.Value{
			"id": helpers.StringValueOrNull(loc.ID),
		})
		diags.Append(objDiags...)
		locationValues = append(locationValues, objValue)
	}
	locationList, listDiags := types.ListValue(types.ObjectType{AttrTypes: extranetLocationAttrTypes}, locationValues)
	diags.Append(listDiags...)

	groupValues := make([]attr.Value, 0, len(dto.LocationGroupDTO))
	for _, group := range dto.LocationGroupDTO {
		objValue, objDiags := types.ObjectValue(extranetLocationAttrTypes, map[string]attr.Value{
			"id": helpers.StringValueOrNull(group.ID),
		})
		diags.Append(objDiags...)
		groupValues = append(groupValues, objValue)
	}
	groupList, groupDiags := types.ListValue(types.ObjectType{AttrTypes: extranetLocationAttrTypes}, groupValues)
	diags.Append(groupDiags...)

	objValue, objDiags := types.ObjectValue(extranetAttrTypes, map[string]attr.Value{
		"zpn_er_id":          helpers.StringValueOrNull(dto.ZpnErID),
		"location_dto":       locationList,
		"location_group_dto": groupList,
	})
	diags.Append(objDiags...)

	listValue, listValueDiags := types.ListValue(types.ObjectType{AttrTypes: extranetAttrTypes}, []attr.Value{objValue})
	diags.Append(listValueDiags...)
	return listValue, diags
}

func validateObjectTypeUniquenessV2(ctx context.Context, conditions []PolicyAccessConditionModel) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, condition := range conditions {
		seen := make(map[string]struct{})
		for _, operand := range condition.Operands {
			objectType := strings.TrimSpace(helpers.StringValue(operand.ObjectType))
			if objectType == "" {
				continue
			}
			if _, exists := seen[objectType]; exists {
				diags.AddError(
					"Invalid policy condition",
					fmt.Sprintf("object_type '%s' can only be specified once in the operands block. Please aggregate all entry_values under the same object_type", objectType),
				)
				return diags
			}
			seen[objectType] = struct{}{}
		}
	}

	return diags
}

func validatePolicyRuleConditionsV2(ctx context.Context, conditions []PolicyAccessConditionModel) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, condition := range conditions {
		for _, operand := range condition.Operands {
			objectType := strings.TrimSpace(helpers.StringValue(operand.ObjectType))
			if objectType == "" {
				continue
			}

			values, valueDiags := helpers.SetValueToStringSlice(ctx, operand.Values)
			diags.Append(valueDiags...)
			if diags.HasError() {
				return diags
			}

			entryModels, entryDiags := extractEntryValueModels(ctx, operand.EntryValues)
			diags.Append(entryDiags...)
			if diags.HasError() {
				return diags
			}

			switch objectType {
			case "APP":
				if len(values) == 0 {
					diags.AddError("Invalid policy condition", "An Application Segment ID must be provided when object_type = APP")
					return diags
				}
			case "APP_GROUP":
				if len(values) == 0 {
					diags.AddError("Invalid policy condition", "A Segment Group ID must be provided when object_type = APP_GROUP")
					return diags
				}
			case "MACHINE_GRP":
				if len(values) == 0 {
					diags.AddError("Invalid policy condition", "A Machine Group ID must be provided when object_type = MACHINE_GRP")
					return diags
				}
			case "LOCATION":
				if len(values) == 0 {
					diags.AddError("Invalid policy condition", "A Location ID must be provided when object_type = LOCATION")
					return diags
				}
			case "EDGE_CONNECTOR_GROUP":
				if len(values) == 0 {
					diags.AddError("Invalid policy condition", "An Edge Connector Group ID must be provided when object_type = EDGE_CONNECTOR_GROUP")
					return diags
				}
			case "BRANCH_CONNECTOR_GROUP":
				if len(values) == 0 {
					diags.AddError("Invalid policy condition", "A Branch Connector Group ID must be provided when object_type = BRANCH_CONNECTOR_GROUP")
					return diags
				}
			case "USER_PORTAL":
				if len(values) == 0 {
					diags.AddError("Invalid policy condition", "A User Portal ID must be provided when object_type = USER_PORTAL")
					return diags
				}
			case "CLIENT_TYPE":
				if len(values) == 0 {
					diags.AddError("Invalid policy condition", fmt.Sprintf("Please provide one of the valid Client Types: %v", validClientTypes))
					return diags
				}
				for _, v := range values {
					if !containsString(validClientTypes, v) {
						diags.AddError("Invalid policy condition", fmt.Sprintf("Invalid Client Type '%s'. Please provide one of the valid Client Types: %v", v, validClientTypes))
						return diags
					}
				}
			case "PLATFORM":
				if len(entryModels) == 0 {
					diags.AddError("Invalid policy condition", fmt.Sprintf("Please provide one of the valid platform types: %v", validPlatformTypes))
					return diags
				}
				for _, entry := range entryModels {
					lhs := strings.TrimSpace(helpers.StringValue(entry.LHS))
					if !containsString(validPlatformTypes, lhs) {
						diags.AddError("Invalid policy condition", fmt.Sprintf("Please provide one of the valid platform types: %v", validPlatformTypes))
						return diags
					}
				}
			case "RISK_FACTOR_TYPE":
				if len(entryModels) == 0 {
					diags.AddError("Invalid policy condition", fmt.Sprintf("Please provide valid risk factor values: %v", validRiskScores))
					return diags
				}
				for _, entry := range entryModels {
					lhs := strings.TrimSpace(helpers.StringValue(entry.LHS))
					rhs := strings.TrimSpace(helpers.StringValue(entry.RHS))
					if lhs != "ZIA" {
						diags.AddError("Invalid policy condition", "LHS must be 'ZIA' for RISK_FACTOR_TYPE")
						return diags
					}
					if !containsString(validRiskScores, rhs) {
						diags.AddError("Invalid policy condition", "RHS must be one of 'UNKNOWN', 'LOW', 'MEDIUM', 'HIGH', 'CRITICAL' for RISK_FACTOR_TYPE")
						return diags
					}
				}
			case "POSTURE":
				if len(entryModels) == 0 {
					diags.AddError("Invalid policy condition", "Please provide a valid Posture UDID")
					return diags
				}
				for _, entry := range entryModels {
					lhs := strings.TrimSpace(helpers.StringValue(entry.LHS))
					rhs := strings.TrimSpace(helpers.StringValue(entry.RHS))
					if lhs == "" {
						diags.AddError("Invalid policy condition", "LHS must be a valid Posture UDID and cannot be empty for POSTURE object_type")
						return diags
					}
					if rhs != "true" && rhs != "false" {
						diags.AddError("Invalid policy condition", "rhs value must be 'true' or 'false' for POSTURE object_type")
						return diags
					}
				}
			case "TRUSTED_NETWORK":
				if len(entryModels) == 0 {
					diags.AddError("Invalid policy condition", "Please provide a valid Network ID")
					return diags
				}
				for _, entry := range entryModels {
					lhs := strings.TrimSpace(helpers.StringValue(entry.LHS))
					rhs := strings.TrimSpace(helpers.StringValue(entry.RHS))
					if lhs == "" {
						diags.AddError("Invalid policy condition", "LHS must be a valid Network ID and cannot be empty for TRUSTED_NETWORK object_type")
						return diags
					}
					if rhs != "true" && rhs != "false" {
						diags.AddError("Invalid policy condition", "rhs value must be 'true' or 'false' for TRUSTED_NETWORK object_type")
						return diags
					}
				}
			case "COUNTRY_CODE":
				if len(entryModels) == 0 {
					diags.AddError("Invalid policy condition", "Please provide a valid country code in 'entry_values'")
					return diags
				}
				var invalidCodes []string
				for _, entry := range entryModels {
					lhs := strings.TrimSpace(helpers.StringValue(entry.LHS))
					rhs := strings.TrimSpace(helpers.StringValue(entry.RHS))
					if lhs == "" || !iso3166.ExistsIso3166ByAlpha2Code(lhs) {
						invalidCodes = append(invalidCodes, lhs)
					}
					if rhs != "true" {
						diags.AddError("Invalid policy condition", "rhs value must be 'true' for COUNTRY_CODE object_type")
						return diags
					}
				}
				if len(invalidCodes) > 0 {
					diags.AddError("Invalid policy condition", fmt.Sprintf("'%s' is not a valid ISO-3166 Alpha-2 country code. Please visit the following site for reference: https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes", strings.Join(invalidCodes, "', '")))
					return diags
				}
			case "SAML":
				if len(entryModels) == 0 {
					diags.AddError("Invalid policy condition", "entry_values must be provided for SAML object_type")
					return diags
				}
				for _, entry := range entryModels {
					lhs := strings.TrimSpace(helpers.StringValue(entry.LHS))
					rhs := strings.TrimSpace(helpers.StringValue(entry.RHS))
					if lhs == "" {
						diags.AddError("Invalid policy condition", "LHS must be a valid SAML attribute ID and cannot be empty for SAML object_type")
						return diags
					}
					if rhs == "" {
						diags.AddError("Invalid policy condition", "RHS must be a valid string and cannot be empty for SAML object_type")
						return diags
					}
				}
			case "SCIM":
				if len(entryModels) == 0 {
					diags.AddError("Invalid policy condition", "entry_values must be provided for SCIM object_type")
					return diags
				}
				for _, entry := range entryModels {
					lhs := strings.TrimSpace(helpers.StringValue(entry.LHS))
					rhs := strings.TrimSpace(helpers.StringValue(entry.RHS))
					if lhs == "" {
						diags.AddError("Invalid policy condition", "LHS must be a valid IdP ID and cannot be empty for SCIM object_type")
						return diags
					}
					if rhs == "" {
						diags.AddError("Invalid policy condition", "RHS must be a valid string and cannot be empty for SCIM object_type")
						return diags
					}
				}
			case "SCIM_GROUP":
				if len(entryModels) == 0 {
					diags.AddError("Invalid policy condition", "entry_values must be provided for SCIM_GROUP object_type")
					return diags
				}
				for _, entry := range entryModels {
					lhs := strings.TrimSpace(helpers.StringValue(entry.LHS))
					rhs := strings.TrimSpace(helpers.StringValue(entry.RHS))
					if lhs == "" {
						diags.AddError("Invalid policy condition", "LHS must be a valid IdP ID and cannot be empty for SCIM_GROUP object_type")
						return diags
					}
					if rhs == "" {
						diags.AddError("Invalid policy condition", "RHS must be a valid SCIM group ID and cannot be empty for SCIM_GROUP object_type")
						return diags
					}
				}
			case "CHROME_POSTURE_PROFILE":
				if len(values) == 0 {
					diags.AddError("Invalid policy condition", "A Chrome Posture Profile ID must be provided when object_type = CHROME_POSTURE_PROFILE")
					return diags
				}
			case "CHROME_ENTERPRISE":
				if len(entryModels) == 0 {
					diags.AddError("Invalid policy condition", "entry_values must be provided for CHROME_ENTERPRISE object_type")
					return diags
				}
				for _, entry := range entryModels {
					lhs := strings.TrimSpace(helpers.StringValue(entry.LHS))
					rhs := strings.TrimSpace(helpers.StringValue(entry.RHS))
					if lhs != "managed" {
						diags.AddError("Invalid policy condition", "LHS must be 'managed' for CHROME_ENTERPRISE object_type")
						return diags
					}
					if rhs != "true" && rhs != "false" {
						diags.AddError("Invalid policy condition", "rhs value must be 'true' or 'false' for CHROME_ENTERPRISE object_type")
						return diags
					}
				}
			}
		}
	}

	return diags
}

func extractEntryValueModels(ctx context.Context, set types.Set) ([]PolicyAccessEntryValueModel, diag.Diagnostics) {
	if set.IsNull() || set.IsUnknown() {
		return nil, diag.Diagnostics{}
	}

	var models []PolicyAccessEntryValueModel
	var diags diag.Diagnostics
	diags.Append(set.ElementsAs(ctx, &models, false)...)
	return models, diags
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

var (
	validClientTypes = []string{
		"zpn_client_type_exporter",
		"zpn_client_type_exporter_noauth",
		"zpn_client_type_machine_tunnel",
		"zpn_client_type_edge_connector",
		"zpn_client_type_zia_inspection",
		"zpn_client_type_vdi",
		"zpn_client_type_zapp",
		"zpn_client_type_slogger",
		"zpn_client_type_browser_isolation",
		"zpn_client_type_ip_anchoring",
		"zpn_client_type_zapp_partner",
		"zpn_client_type_branch_connector",
	}
	validPlatformTypes = []string{"mac", "linux", "ios", "windows", "android"}
	validRiskScores    = []string{"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL"}
)
