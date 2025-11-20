package resources

import (
	"context"
	"fmt"
	"html"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/lssconfigcontroller"
)

var (
	_ resource.Resource                = &LSSConfigControllerResource{}
	_ resource.ResourceWithConfigure   = &LSSConfigControllerResource{}
	_ resource.ResourceWithImportState = &LSSConfigControllerResource{}
)

func NewLSSConfigControllerResource() resource.Resource {
	return &LSSConfigControllerResource{}
}

type LSSConfigControllerResource struct {
	client *client.Client
}

type LSSConfigControllerModel struct {
	ID                 types.String             `tfsdk:"id"`
	PolicyRuleID       types.String             `tfsdk:"policy_rule_id"`
	PolicyRuleResource *PolicyRuleResourceModel `tfsdk:"policy_rule_resource"`
	ConnectorGroups    []ConnectorGroupModel    `tfsdk:"connector_groups"`
	Config             *LSSConfigBlockModel     `tfsdk:"config"`
}

type ConnectorGroupModel struct {
	IDs types.List `tfsdk:"id"`
}

type LSSConfigBlockModel struct {
	AuditMessage  types.String `tfsdk:"audit_message"`
	Description   types.String `tfsdk:"description"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Filter        types.Set    `tfsdk:"filter"`
	Format        types.String `tfsdk:"format"`
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	LSSHost       types.String `tfsdk:"lss_host"`
	LSSPort       types.String `tfsdk:"lss_port"`
	SourceLogType types.String `tfsdk:"source_log_type"`
	UseTLS        types.Bool   `tfsdk:"use_tls"`
}

type PolicyRuleResourceModel struct {
	Action                 types.String               `tfsdk:"action"`
	ActionID               types.String               `tfsdk:"action_id"`
	BypassDefaultRule      types.Bool                 `tfsdk:"bypass_default_rule"`
	CustomMsg              types.String               `tfsdk:"custom_msg"`
	DefaultRule            types.Bool                 `tfsdk:"default_rule"`
	Description            types.String               `tfsdk:"description"`
	ID                     types.String               `tfsdk:"id"`
	Name                   types.String               `tfsdk:"name"`
	Operator               types.String               `tfsdk:"operator"`
	PolicySetID            types.String               `tfsdk:"policy_set_id"`
	PolicyType             types.String               `tfsdk:"policy_type"`
	Priority               types.String               `tfsdk:"priority"`
	ReauthDefaultRule      types.Bool                 `tfsdk:"reauth_default_rule"`
	ReauthIdleTimeout      types.String               `tfsdk:"reauth_idle_timeout"`
	ReauthTimeout          types.String               `tfsdk:"reauth_timeout"`
	ZPNIsolationProfileID  types.String               `tfsdk:"zpn_isolation_profile_id"`
	ZPNCBIProfileID        types.String               `tfsdk:"zpn_cbi_profile_id"`
	ZPNInspectionProfileID types.String               `tfsdk:"zpn_inspection_profile_id"`
	RuleOrder              types.String               `tfsdk:"rule_order"`
	MicrotenantID          types.String               `tfsdk:"microtenant_id"`
	LSSDefaultRule         types.Bool                 `tfsdk:"lss_default_rule"`
	Conditions             []PolicyRuleConditionModel `tfsdk:"conditions"`
}

type PolicyRuleConditionModel struct {
	Operator types.String             `tfsdk:"operator"`
	Operands []PolicyRuleOperandModel `tfsdk:"operands"`
}

type PolicyRuleOperandModel struct {
	Values      types.Set                   `tfsdk:"values"`
	ObjectType  types.String                `tfsdk:"object_type"`
	EntryValues []PolicyRuleEntryValueModel `tfsdk:"entry_values"`
}

type PolicyRuleEntryValueModel struct {
	LHS types.String `tfsdk:"lhs"`
	RHS types.String `tfsdk:"rhs"`
}

func (r *LSSConfigControllerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lss_config_controller"
}

func (r *LSSConfigControllerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages ZPA LSS configuration controller.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"policy_rule_id": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			// policy_rule_resource: TypeList with MaxItems: 1 in SDKv2
			// Using SingleNestedBlock for single block syntax support
			"policy_rule_resource": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"action": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("LOG"),
						},
					},
					"action_id": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"bypass_default_rule": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
					"custom_msg": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"default_rule": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
					"description": schema.StringAttribute{
						Optional: true,
					},
					"id": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"name": schema.StringAttribute{
						Required: true,
					},
					"operator": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Validators: []validator.String{
							stringvalidator.OneOf("AND", "OR"),
						},
					},
					"policy_set_id": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"policy_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"priority": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"reauth_default_rule": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
					"reauth_idle_timeout": schema.StringAttribute{
						Optional: true,
					},
					"reauth_timeout": schema.StringAttribute{
						Optional: true,
					},
					"zpn_isolation_profile_id": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"zpn_cbi_profile_id": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"zpn_inspection_profile_id": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"rule_order": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"microtenant_id": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"lss_default_rule": schema.BoolAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					// conditions: TypeList in SDKv2
					// Using ListNestedBlock for block syntax support
					"conditions": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"operator": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("AND", "OR"),
									},
								},
							},
							Blocks: map[string]schema.Block{
								// operands: TypeList in SDKv2
								// Using ListNestedBlock for block syntax support
								"operands": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"values": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"object_type": schema.StringAttribute{
												Required: true,
												Validators: []validator.String{
													stringvalidator.OneOf(
														"APP",
														"APP_GROUP",
														"CLIENT_TYPE",
														"IDP",
														"SCIM",
														"SCIM_GROUP",
														"SAML",
													),
												},
											},
										},
										Blocks: map[string]schema.Block{
											// entry_values: TypeList in SDKv2
											// Using ListNestedBlock for block syntax support
											"entry_values": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"lhs": schema.StringAttribute{
															Optional: true,
														},
														"rhs": schema.StringAttribute{
															Optional: true,
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
				},
			},
			// connector_groups: TypeSet in SDKv2, id is TypeList (Optional)
			// Using SetNestedBlock for block syntax support
			"connector_groups": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
			// config: TypeList with MaxItems: 1 in SDKv2
			// Using SingleNestedBlock for single block syntax support
			"config": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"audit_message": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"description": schema.StringAttribute{
						Optional: true,
					},
					"enabled": schema.BoolAttribute{
						Optional:      true,
						Computed:      true,
						Default:       booldefault.StaticBool(true),
						PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
					},
					"filter": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"format": schema.StringAttribute{
						Required: true,
					},
					"id": schema.StringAttribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"lss_host": schema.StringAttribute{
						Required: true,
					},
					"lss_port": schema.StringAttribute{
						Required: true,
					},
					"source_log_type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"zpn_trans_log",
								"zpn_auth_log",
								"zpn_ast_auth_log",
								"zpn_http_trans_log",
								"zpn_audit_log",
								"zpn_ast_comprehensive_stats",
								"zpn_sys_auth_log",
								"zpn_waf_http_exchanges_log",
								"zpn_pbroker_comprehensive_stats",
							),
						},
					},
					"use_tls": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
				},
			},
		},
	}
}

func (r *LSSConfigControllerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LSSConfigControllerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan LSSConfigControllerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diags := expandLSSConfigController(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := lssconfigcontroller.Create(ctx, r.client.Service, &payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create LSS config controller: %v", err))
		return
	}

	// Pass plan to preserve policy_rule_resource (matching SDKv2 behavior)
	state, readDiags := r.readIntoState(ctx, created.ID, &plan)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *LSSConfigControllerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state LSSConfigControllerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Pass existing state to preserve policy_rule_resource (matching SDKv2 behavior)
	newState, diags := r.readIntoState(ctx, state.ID.ValueString(), &state)
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

func (r *LSSConfigControllerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan LSSConfigControllerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	payload, diags := expandLSSConfigController(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := lssconfigcontroller.Update(ctx, r.client.Service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update LSS config controller: %v", err))
		return
	}

	// Pass plan to preserve policy_rule_resource (matching SDKv2 behavior)
	state, readDiags := r.readIntoState(ctx, plan.ID.ValueString(), &plan)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *LSSConfigControllerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state LSSConfigControllerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := lssconfigcontroller.Delete(ctx, r.client.Service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete LSS config controller: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *LSSConfigControllerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before importing resources.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the LSS config controller ID or name.")
		return
	}

	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		resource, _, lookupErr := lssconfigcontroller.GetByName(ctx, r.client.Service, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate LSS config controller %q: %v", id, lookupErr))
			return
		}
		id = resource.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *LSSConfigControllerResource) readIntoState(ctx context.Context, id string, existingState *LSSConfigControllerModel) (LSSConfigControllerModel, diag.Diagnostics) {
	resource, _, err := lssconfigcontroller.Get(ctx, r.client.Service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return LSSConfigControllerModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("LSS config controller %s not found", id))}
		}
		return LSSConfigControllerModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read LSS config controller: %v", err))}
	}

	var diags diag.Diagnostics

	configModel, configDiags := flattenLSSConfig(ctx, resource.LSSConfig)
	diags.Append(configDiags...)

	connectorModels, connectorDiags := flattenConnectorGroups(ctx, resource.ConnectorGroups)
	diags.Append(connectorDiags...)

	// Preserve policy_rule_resource from existing state/plan (matching SDKv2 behavior)
	// SDKv2 doesn't set policy_rule_resource in Read, so it's preserved from state
	// We need to convert all unknown values to null to satisfy Terraform's requirements
	var policyRuleResourceModel *PolicyRuleResourceModel
	apiPolicyRuleResource, apiDiags := flattenPolicyRuleResource(ctx, resource.PolicyRuleResource)
	diags.Append(apiDiags...)

	if existingState != nil && existingState.PolicyRuleResource != nil {
		// Convert unknown values to null in preserved model, then merge with API response
		preserved := convertUnknownToNullPolicyRuleResource(existingState.PolicyRuleResource)
		policyRuleResourceModel = mergePolicyRuleResource(preserved, apiPolicyRuleResource)
	} else {
		// Only flatten from API if not present in existing state
		policyRuleResourceModel = apiPolicyRuleResource
	}

	model := LSSConfigControllerModel{
		ID:                 helpers.StringValueOrNull(resource.ID),
		PolicyRuleID:       types.StringNull(),
		PolicyRuleResource: policyRuleResourceModel,
		ConnectorGroups:    connectorModels,
		Config:             configModel,
	}

	if resource.PolicyRule != nil {
		model.PolicyRuleID = helpers.StringValueOrNull(resource.PolicyRule.ID)
	}

	return model, diags
}

func expandLSSConfigController(ctx context.Context, model *LSSConfigControllerModel) (lssconfigcontroller.LSSResource, diag.Diagnostics) {
	var diags diag.Diagnostics

	if model.Config == nil {
		diags.AddError("Validation Error", "config block must be provided")
		return lssconfigcontroller.LSSResource{}, diags
	}

	config, configDiags := expandLSSConfigBlock(ctx, model.Config)
	diags.Append(configDiags...)

	// Set config ID to top-level ID (required by API for updates, matching SDKv2 behavior)
	// SDKv2 sets config.ID from d.Get("id") at line 560 in resource_zpa_lss_config_controller.go
	topLevelID := helpers.StringValue(model.ID)
	if config != nil && topLevelID != "" {
		config.ID = topLevelID
	}

	connectorGroups, connectorDiags := expandConnectorGroups(ctx, model.ConnectorGroups)
	diags.Append(connectorDiags...)

	policyResource, policyDiags := expandPolicyRuleResource(ctx, model.PolicyRuleResource)
	diags.Append(policyDiags...)

	result := lssconfigcontroller.LSSResource{
		ID:                 topLevelID,
		LSSConfig:          config,
		ConnectorGroups:    connectorGroups,
		PolicyRuleResource: policyResource,
	}

	if config != nil && config.SourceLogType != "" {
		sourceLogType := config.SourceLogType
		if policyResource != nil && policyResource.Conditions != nil {
			for _, condition := range policyResource.Conditions {
				if condition.Operands != nil {
					operands := *condition.Operands
					for _, operand := range operands {
						diags.Append(helpers.ValidateLSSConfigControllerFilters(sourceLogType, operand.ObjectType, "", operand.Values, operands)...)
					}
				}
			}
		}
		for _, filter := range config.Filter {
			diags.Append(helpers.ValidateLSSConfigControllerFilters(sourceLogType, "", filter, nil, nil)...)
		}
	}

	return result, diags
}

func expandLSSConfigBlock(ctx context.Context, model *LSSConfigBlockModel) (*lssconfigcontroller.LSSConfig, diag.Diagnostics) {
	if model == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	filter, filterDiags := helpers.SetValueToStringSlice(ctx, model.Filter)
	diags.Append(filterDiags...)

	config := &lssconfigcontroller.LSSConfig{
		ID:            helpers.StringValue(model.ID),
		Name:          helpers.StringValue(model.Name),
		Description:   helpers.StringValue(model.Description),
		Enabled:       helpers.BoolValue(model.Enabled, true),
		Filter:        filter,
		Format:        helpers.StringValue(model.Format),
		AuditMessage:  helpers.StringValue(model.AuditMessage),
		LSSHost:       helpers.StringValue(model.LSSHost),
		LSSPort:       helpers.StringValue(model.LSSPort),
		SourceLogType: helpers.StringValue(model.SourceLogType),
		UseTLS:        helpers.BoolValue(model.UseTLS, false),
	}

	return config, diags
}

func expandConnectorGroups(ctx context.Context, models []ConnectorGroupModel) ([]lssconfigcontroller.ConnectorGroups, diag.Diagnostics) {
	var diags diag.Diagnostics
	var result []lssconfigcontroller.ConnectorGroups

	for _, model := range models {
		ids, idsDiags := helpers.ListValueToStringSlice(ctx, model.IDs)
		diags.Append(idsDiags...)
		for _, id := range ids {
			if strings.TrimSpace(id) == "" {
				continue
			}
			result = append(result, lssconfigcontroller.ConnectorGroups{ID: id})
		}
	}

	return result, diags
}

func expandPolicyRuleResource(ctx context.Context, model *PolicyRuleResourceModel) (*lssconfigcontroller.PolicyRuleResource, diag.Diagnostics) {
	if model == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	resource := &lssconfigcontroller.PolicyRuleResource{
		ID:                     helpers.StringValue(model.ID),
		Name:                   helpers.StringValue(model.Name),
		Description:            helpers.StringValue(model.Description),
		Action:                 helpers.StringValue(model.Action),
		ActionID:               helpers.StringValue(model.ActionID),
		CustomMsg:              helpers.StringValue(model.CustomMsg),
		Operator:               helpers.StringValue(model.Operator),
		PolicySetID:            helpers.StringValue(model.PolicySetID),
		PolicyType:             helpers.StringValue(model.PolicyType),
		Priority:               helpers.StringValue(model.Priority),
		ReauthIdleTimeout:      helpers.StringValue(model.ReauthIdleTimeout),
		ReauthTimeout:          helpers.StringValue(model.ReauthTimeout),
		RuleOrder:              helpers.StringValue(model.RuleOrder),
		ZpnCbiProfileID:        helpers.StringValue(model.ZPNCBIProfileID),
		ZpnInspectionProfileID: helpers.StringValue(model.ZPNInspectionProfileID),
		MicroTenantID:          helpers.StringValue(model.MicrotenantID),
	}

	if len(model.Conditions) > 0 {
		conditions := make([]lssconfigcontroller.PolicyRuleResourceConditions, 0, len(model.Conditions))
		for _, conditionModel := range model.Conditions {
			condition := lssconfigcontroller.PolicyRuleResourceConditions{
				Operator: helpers.StringValue(conditionModel.Operator),
			}

			if len(conditionModel.Operands) > 0 {
				operands := make([]lssconfigcontroller.PolicyRuleResourceOperands, 0, len(conditionModel.Operands))
				for _, operandModel := range conditionModel.Operands {
					values, valuesDiags := helpers.SetValueToStringSlice(ctx, operandModel.Values)
					diags.Append(valuesDiags...)

					entryValues := expandPolicyRuleEntryValues(operandModel.EntryValues)

					var entryPointer *[]lssconfigcontroller.OperandsResourceLHSRHSValue
					if len(entryValues) > 0 {
						entryPointer = &entryValues
					}

					operands = append(operands, lssconfigcontroller.PolicyRuleResourceOperands{
						ObjectType:                  helpers.StringValue(operandModel.ObjectType),
						Values:                      values,
						OperandsResourceLHSRHSValue: entryPointer,
					})
				}
				condition.Operands = &operands
			}

			conditions = append(conditions, condition)
		}
		resource.Conditions = conditions
	}

	return resource, diags
}

func expandPolicyRuleEntryValues(models []PolicyRuleEntryValueModel) []lssconfigcontroller.OperandsResourceLHSRHSValue {
	if len(models) == 0 {
		return nil
	}
	result := make([]lssconfigcontroller.OperandsResourceLHSRHSValue, 0, len(models))
	for _, model := range models {
		result = append(result, lssconfigcontroller.OperandsResourceLHSRHSValue{
			LHS: helpers.StringValue(model.LHS),
			RHS: helpers.StringValue(model.RHS),
		})
	}
	return result
}

func flattenLSSConfig(ctx context.Context, cfg *lssconfigcontroller.LSSConfig) (*LSSConfigBlockModel, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	filter, filterDiags := types.SetValueFrom(ctx, types.StringType, cfg.Filter)
	diags.Append(filterDiags...)

	model := &LSSConfigBlockModel{
		AuditMessage:  helpers.StringValueOrNull(html.UnescapeString(cfg.AuditMessage)),
		Description:   helpers.StringValueOrNull(cfg.Description),
		Enabled:       types.BoolValue(cfg.Enabled),
		Filter:        filter,
		Format:        helpers.StringValueOrNull(html.UnescapeString(cfg.Format)),
		ID:            helpers.StringValueOrNull(cfg.ID),
		Name:          helpers.StringValueOrNull(cfg.Name),
		LSSHost:       helpers.StringValueOrNull(cfg.LSSHost),
		LSSPort:       helpers.StringValueOrNull(cfg.LSSPort),
		SourceLogType: helpers.StringValueOrNull(cfg.SourceLogType),
		UseTLS:        types.BoolValue(cfg.UseTLS),
	}

	return model, diags
}

func flattenConnectorGroups(ctx context.Context, groups []lssconfigcontroller.ConnectorGroups) ([]ConnectorGroupModel, diag.Diagnostics) {
	if len(groups) == 0 {
		return nil, nil
	}

	var diags diag.Diagnostics
	ids := make([]string, 0, len(groups))
	for _, group := range groups {
		if strings.TrimSpace(group.ID) == "" {
			continue
		}
		ids = append(ids, group.ID)
	}

	if len(ids) == 0 {
		return nil, nil
	}

	listValue, listDiags := types.ListValueFrom(ctx, types.StringType, ids)
	diags.Append(listDiags...)

	if diags.HasError() {
		return nil, diags
	}

	return []ConnectorGroupModel{{IDs: listValue}}, diags
}

func flattenPolicyRuleResource(ctx context.Context, resource *lssconfigcontroller.PolicyRuleResource) (*PolicyRuleResourceModel, diag.Diagnostics) {
	if resource == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	model := &PolicyRuleResourceModel{
		Action:                 helpers.StringValueOrNull(resource.Action),
		ActionID:               helpers.StringValueOrNull(resource.ActionID),
		BypassDefaultRule:      types.BoolNull(),
		CustomMsg:              helpers.StringValueOrNull(resource.CustomMsg),
		DefaultRule:            types.BoolNull(),
		Description:            helpers.StringValueOrNull(resource.Description),
		ID:                     helpers.StringValueOrNull(resource.ID),
		Name:                   helpers.StringValueOrNull(resource.Name),
		Operator:               helpers.StringValueOrNull(resource.Operator),
		PolicySetID:            helpers.StringValueOrNull(resource.PolicySetID),
		PolicyType:             helpers.StringValueOrNull(resource.PolicyType),
		Priority:               helpers.StringValueOrNull(resource.Priority),
		ReauthDefaultRule:      types.BoolNull(),
		ReauthIdleTimeout:      helpers.StringValueOrNull(resource.ReauthIdleTimeout),
		ReauthTimeout:          helpers.StringValueOrNull(resource.ReauthTimeout),
		ZPNIsolationProfileID:  types.StringNull(),
		ZPNCBIProfileID:        helpers.StringValueOrNull(resource.ZpnCbiProfileID),
		ZPNInspectionProfileID: helpers.StringValueOrNull(resource.ZpnInspectionProfileID),
		RuleOrder:              helpers.StringValueOrNull(resource.RuleOrder),
		MicrotenantID:          helpers.StringValueOrNull(resource.MicroTenantID),
		LSSDefaultRule:         types.BoolNull(),
	}

	if len(resource.Conditions) > 0 {
		conditions := make([]PolicyRuleConditionModel, 0, len(resource.Conditions))
		for _, condition := range resource.Conditions {
			conditionModel := PolicyRuleConditionModel{
				Operator: helpers.StringValueOrNull(condition.Operator),
			}

			if condition.Operands != nil {
				operands := *condition.Operands
				operandModels := make([]PolicyRuleOperandModel, 0, len(operands))
				for _, operand := range operands {
					values, valuesDiags := types.SetValueFrom(ctx, types.StringType, operand.Values)
					diags.Append(valuesDiags...)

					entryValues := flattenPolicyRuleEntryValues(operand.OperandsResourceLHSRHSValue)

					operandModels = append(operandModels, PolicyRuleOperandModel{
						Values:      values,
						ObjectType:  helpers.StringValueOrNull(operand.ObjectType),
						EntryValues: entryValues,
					})
				}
				conditionModel.Operands = operandModels
			}

			conditions = append(conditions, conditionModel)
		}
		model.Conditions = conditions
	}

	return model, diags
}

func flattenPolicyRuleEntryValues(entryValues *[]lssconfigcontroller.OperandsResourceLHSRHSValue) []PolicyRuleEntryValueModel {
	if entryValues == nil || len(*entryValues) == 0 {
		return nil
	}

	result := make([]PolicyRuleEntryValueModel, 0, len(*entryValues))
	for _, value := range *entryValues {
		result = append(result, PolicyRuleEntryValueModel{
			LHS: helpers.StringValueOrNull(value.LHS),
			RHS: helpers.StringValueOrNull(value.RHS),
		})
	}
	return result
}

// mergePolicyRuleResource merges preserved values from plan with computed values from API response
// Preserves user-provided values, populates computed fields from API if null/unknown in plan
func mergePolicyRuleResource(preserved *PolicyRuleResourceModel, apiResponse *PolicyRuleResourceModel) *PolicyRuleResourceModel {
	if preserved == nil {
		return apiResponse
	}
	if apiResponse == nil {
		return preserved
	}

	merged := &PolicyRuleResourceModel{
		// Preserve user-provided values
		Action:      preserved.Action,
		Name:        preserved.Name,
		Description: preserved.Description,
		CustomMsg:   preserved.CustomMsg,
		PolicySetID: preserved.PolicySetID,
		Conditions:  preserved.Conditions,
	}

	// Populate computed fields from API if null/unknown in plan
	// If API response is null/unknown, set to null (not unknown) to satisfy Terraform's requirement
	if preserved.ActionID.IsNull() || preserved.ActionID.IsUnknown() {
		if apiResponse.ActionID.IsNull() || apiResponse.ActionID.IsUnknown() {
			merged.ActionID = types.StringNull()
		} else {
			merged.ActionID = apiResponse.ActionID
		}
	} else {
		merged.ActionID = preserved.ActionID
	}

	if preserved.ID.IsNull() || preserved.ID.IsUnknown() {
		if apiResponse.ID.IsNull() || apiResponse.ID.IsUnknown() {
			merged.ID = types.StringNull()
		} else {
			merged.ID = apiResponse.ID
		}
	} else {
		merged.ID = preserved.ID
	}

	if preserved.Operator.IsNull() || preserved.Operator.IsUnknown() {
		if apiResponse.Operator.IsNull() || apiResponse.Operator.IsUnknown() {
			merged.Operator = types.StringNull()
		} else {
			merged.Operator = apiResponse.Operator
		}
	} else {
		merged.Operator = preserved.Operator
	}

	if preserved.PolicyType.IsNull() || preserved.PolicyType.IsUnknown() {
		if apiResponse.PolicyType.IsNull() || apiResponse.PolicyType.IsUnknown() {
			merged.PolicyType = types.StringNull()
		} else {
			merged.PolicyType = apiResponse.PolicyType
		}
	} else {
		merged.PolicyType = preserved.PolicyType
	}

	if preserved.Priority.IsNull() || preserved.Priority.IsUnknown() {
		if apiResponse.Priority.IsNull() || apiResponse.Priority.IsUnknown() {
			merged.Priority = types.StringNull()
		} else {
			merged.Priority = apiResponse.Priority
		}
	} else {
		merged.Priority = preserved.Priority
	}

	if preserved.RuleOrder.IsNull() || preserved.RuleOrder.IsUnknown() {
		if apiResponse.RuleOrder.IsNull() || apiResponse.RuleOrder.IsUnknown() {
			merged.RuleOrder = types.StringNull()
		} else {
			merged.RuleOrder = apiResponse.RuleOrder
		}
	} else {
		merged.RuleOrder = preserved.RuleOrder
	}

	if preserved.ZPNCBIProfileID.IsNull() || preserved.ZPNCBIProfileID.IsUnknown() {
		if apiResponse.ZPNCBIProfileID.IsNull() || apiResponse.ZPNCBIProfileID.IsUnknown() {
			merged.ZPNCBIProfileID = types.StringNull()
		} else {
			merged.ZPNCBIProfileID = apiResponse.ZPNCBIProfileID
		}
	} else {
		merged.ZPNCBIProfileID = preserved.ZPNCBIProfileID
	}

	if preserved.ZPNInspectionProfileID.IsNull() || preserved.ZPNInspectionProfileID.IsUnknown() {
		if apiResponse.ZPNInspectionProfileID.IsNull() || apiResponse.ZPNInspectionProfileID.IsUnknown() {
			merged.ZPNInspectionProfileID = types.StringNull()
		} else {
			merged.ZPNInspectionProfileID = apiResponse.ZPNInspectionProfileID
		}
	} else {
		merged.ZPNInspectionProfileID = preserved.ZPNInspectionProfileID
	}

	if preserved.ZPNIsolationProfileID.IsNull() || preserved.ZPNIsolationProfileID.IsUnknown() {
		if apiResponse.ZPNIsolationProfileID.IsNull() || apiResponse.ZPNIsolationProfileID.IsUnknown() {
			merged.ZPNIsolationProfileID = types.StringNull()
		} else {
			merged.ZPNIsolationProfileID = apiResponse.ZPNIsolationProfileID
		}
	} else {
		merged.ZPNIsolationProfileID = preserved.ZPNIsolationProfileID
	}

	if preserved.MicrotenantID.IsNull() || preserved.MicrotenantID.IsUnknown() {
		if apiResponse.MicrotenantID.IsNull() || apiResponse.MicrotenantID.IsUnknown() {
			merged.MicrotenantID = types.StringNull()
		} else {
			merged.MicrotenantID = apiResponse.MicrotenantID
		}
	} else {
		merged.MicrotenantID = preserved.MicrotenantID
	}

	if preserved.ReauthIdleTimeout.IsNull() || preserved.ReauthIdleTimeout.IsUnknown() {
		if apiResponse.ReauthIdleTimeout.IsNull() || apiResponse.ReauthIdleTimeout.IsUnknown() {
			merged.ReauthIdleTimeout = types.StringNull()
		} else {
			merged.ReauthIdleTimeout = apiResponse.ReauthIdleTimeout
		}
	} else {
		merged.ReauthIdleTimeout = preserved.ReauthIdleTimeout
	}

	if preserved.ReauthTimeout.IsNull() || preserved.ReauthTimeout.IsUnknown() {
		if apiResponse.ReauthTimeout.IsNull() || apiResponse.ReauthTimeout.IsUnknown() {
			merged.ReauthTimeout = types.StringNull()
		} else {
			merged.ReauthTimeout = apiResponse.ReauthTimeout
		}
	} else {
		merged.ReauthTimeout = preserved.ReauthTimeout
	}

	// Boolean fields: populate from API if null/unknown in plan
	// If API response is null/unknown, set to null (not unknown)
	if preserved.DefaultRule.IsNull() || preserved.DefaultRule.IsUnknown() {
		if apiResponse.DefaultRule.IsNull() || apiResponse.DefaultRule.IsUnknown() {
			merged.DefaultRule = types.BoolNull()
		} else {
			merged.DefaultRule = apiResponse.DefaultRule
		}
	} else {
		merged.DefaultRule = preserved.DefaultRule
	}

	if preserved.BypassDefaultRule.IsNull() || preserved.BypassDefaultRule.IsUnknown() {
		if apiResponse.BypassDefaultRule.IsNull() || apiResponse.BypassDefaultRule.IsUnknown() {
			merged.BypassDefaultRule = types.BoolNull()
		} else {
			merged.BypassDefaultRule = apiResponse.BypassDefaultRule
		}
	} else {
		merged.BypassDefaultRule = preserved.BypassDefaultRule
	}

	if preserved.ReauthDefaultRule.IsNull() || preserved.ReauthDefaultRule.IsUnknown() {
		if apiResponse.ReauthDefaultRule.IsNull() || apiResponse.ReauthDefaultRule.IsUnknown() {
			merged.ReauthDefaultRule = types.BoolNull()
		} else {
			merged.ReauthDefaultRule = apiResponse.ReauthDefaultRule
		}
	} else {
		merged.ReauthDefaultRule = preserved.ReauthDefaultRule
	}

	if preserved.LSSDefaultRule.IsNull() || preserved.LSSDefaultRule.IsUnknown() {
		if apiResponse.LSSDefaultRule.IsNull() || apiResponse.LSSDefaultRule.IsUnknown() {
			merged.LSSDefaultRule = types.BoolNull()
		} else {
			merged.LSSDefaultRule = apiResponse.LSSDefaultRule
		}
	} else {
		merged.LSSDefaultRule = preserved.LSSDefaultRule
	}

	return merged
}

// convertUnknownToNullPolicyRuleResource converts all unknown values to null in the preserved model
// This ensures Terraform's requirement that all values must be known after apply
func convertUnknownToNullPolicyRuleResource(preserved *PolicyRuleResourceModel) *PolicyRuleResourceModel {
	if preserved == nil {
		return nil
	}

	converted := &PolicyRuleResourceModel{
		// Preserve user-provided values, converting unknown to null
		Action:      convertUnknownStringToNull(preserved.Action),
		Name:        convertUnknownStringToNull(preserved.Name),
		Description: convertUnknownStringToNull(preserved.Description),
		CustomMsg:   convertUnknownStringToNull(preserved.CustomMsg),
		PolicySetID: convertUnknownStringToNull(preserved.PolicySetID),
		Conditions:  preserved.Conditions, // Conditions are preserved as-is (they're user-provided)
	}

	// Convert all computed/optional fields from unknown to null
	converted.ActionID = convertUnknownStringToNull(preserved.ActionID)
	converted.ID = convertUnknownStringToNull(preserved.ID)
	converted.Operator = convertUnknownStringToNull(preserved.Operator)
	converted.PolicyType = convertUnknownStringToNull(preserved.PolicyType)
	converted.Priority = convertUnknownStringToNull(preserved.Priority)
	converted.RuleOrder = convertUnknownStringToNull(preserved.RuleOrder)
	converted.ZPNCBIProfileID = convertUnknownStringToNull(preserved.ZPNCBIProfileID)
	converted.ZPNInspectionProfileID = convertUnknownStringToNull(preserved.ZPNInspectionProfileID)
	converted.ZPNIsolationProfileID = convertUnknownStringToNull(preserved.ZPNIsolationProfileID)
	converted.MicrotenantID = convertUnknownStringToNull(preserved.MicrotenantID)
	converted.ReauthIdleTimeout = convertUnknownStringToNull(preserved.ReauthIdleTimeout)
	converted.ReauthTimeout = convertUnknownStringToNull(preserved.ReauthTimeout)
	converted.DefaultRule = convertUnknownBoolToNull(preserved.DefaultRule)
	converted.BypassDefaultRule = convertUnknownBoolToNull(preserved.BypassDefaultRule)
	converted.ReauthDefaultRule = convertUnknownBoolToNull(preserved.ReauthDefaultRule)
	converted.LSSDefaultRule = convertUnknownBoolToNull(preserved.LSSDefaultRule)

	return converted
}

// convertUnknownStringToNull converts unknown string values to null
func convertUnknownStringToNull(value types.String) types.String {
	if value.IsUnknown() {
		return types.StringNull()
	}
	return value
}

// convertUnknownBoolToNull converts unknown bool values to null
func convertUnknownBoolToNull(value types.Bool) types.Bool {
	if value.IsUnknown() {
		return types.BoolNull()
	}
	return value
}
