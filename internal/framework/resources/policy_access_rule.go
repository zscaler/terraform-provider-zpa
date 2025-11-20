package resources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	stringvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	fwrschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	fwvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

var (
	_ resource.Resource                = &PolicyAccessRuleResource{}
	_ resource.ResourceWithConfigure   = &PolicyAccessRuleResource{}
	_ resource.ResourceWithImportState = &PolicyAccessRuleResource{}
)

const policyAccessType = "ACCESS_POLICY"

func NewPolicyAccessRuleResource() resource.Resource {
	return &PolicyAccessRuleResource{}
}

type PolicyAccessRuleResource struct {
	client *client.Client
}

type PolicyAccessRuleModel struct {
	ID                     types.String                   `tfsdk:"id"`
	Name                   types.String                   `tfsdk:"name"`
	Description            types.String                   `tfsdk:"description"`
	Action                 types.String                   `tfsdk:"action"`
	ActionID               types.String                   `tfsdk:"action_id"`
	BypassDefaultRule      types.Bool                     `tfsdk:"bypass_default_rule"`
	CustomMsg              types.String                   `tfsdk:"custom_msg"`
	DefaultRule            types.Bool                     `tfsdk:"default_rule"`
	Operator               types.String                   `tfsdk:"operator"`
	PolicySetID            types.String                   `tfsdk:"policy_set_id"`
	PolicyType             types.String                   `tfsdk:"policy_type"`
	Priority               types.String                   `tfsdk:"priority"`
	ReauthDefaultRule      types.Bool                     `tfsdk:"reauth_default_rule"`
	ReauthIdleTimeout      types.String                   `tfsdk:"reauth_idle_timeout"`
	ReauthTimeout          types.String                   `tfsdk:"reauth_timeout"`
	ZPNIsolationProfileID  types.String                   `tfsdk:"zpn_isolation_profile_id"`
	ZPNCBIProfileID        types.String                   `tfsdk:"zpn_cbi_profile_id"`
	ZPNInspectionProfileID types.String                   `tfsdk:"zpn_inspection_profile_id"`
	RuleOrder              types.String                   `tfsdk:"rule_order"`
	MicrotenantID          types.String                   `tfsdk:"microtenant_id"`
	LSSDefaultRule         types.Bool                     `tfsdk:"lss_default_rule"`
	AppServerGroups        types.List                     `tfsdk:"app_server_groups"`
	AppConnectorGroups     types.List                     `tfsdk:"app_connector_groups"`
	Conditions             []helpers.PolicyConditionModel `tfsdk:"conditions"`
}

func (r *PolicyAccessRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_access_rule"
}

func (r *PolicyAccessRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	base := helpers.PolicyCommonSchemaAttributes()
	attributes := make(map[string]fwrschema.Attribute, len(base)+4)

	for k, v := range base {
		attributes[k] = v
	}

	attributes["action"] = fwrschema.StringAttribute{
		Optional: true,
		Validators: []fwvalidator.String{
			stringvalidator.OneOf("ALLOW", "DENY", "REQUIRE_APPROVAL"),
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
	}
	resp.Schema = fwrschema.Schema{
		Description: "Manages ZPA Access Policy rules.",
		Attributes:  attributes,
		Blocks: map[string]fwrschema.Block{
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
			"conditions": helpers.PolicyConditionsBlock(objectTypes),
		},
	}
}

func (r *PolicyAccessRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PolicyAccessRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy access rules.")
		return
	}

	var plan PolicyAccessRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)
	helperClient := helpers.NewHelperClient(r.client)

	policySetID := helpers.StringValue(plan.PolicySetID)
	microTenantID := helpers.StringValue(plan.MicrotenantID)
	if policySetID == "" {
		var err error
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, policyAccessType, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	request, diags := expandPolicyAccessRule(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := helpers.ValidateConditions(ctx, request.Conditions, helperClient, microTenantID); err != nil {
		resp.Diagnostics.AddError("Validation Error", err.Error())
		return
	}

	created, _, err := policysetcontroller.CreateRule(ctx, service, request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create policy access rule: %v", err))
		return
	}

	state, readDiags := r.readPolicyAccessRule(ctx, service, policySetID, created.ID, plan.MicrotenantID, &plan)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PolicyAccessRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy access rules.")
		return
	}

	var state PolicyAccessRuleModel
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, policyAccessType, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	newState, diags := r.readPolicyAccessRule(ctx, service, policySetID, state.ID.ValueString(), state.MicrotenantID, &state)
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

func (r *PolicyAccessRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy access rules.")
		return
	}

	var plan PolicyAccessRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)
	helperClient := helpers.NewHelperClient(r.client)

	policySetID := helpers.StringValue(plan.PolicySetID)
	microTenantID := helpers.StringValue(plan.MicrotenantID)
	if policySetID == "" {
		var err error
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, policyAccessType, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	// Check if resource still exists before updating
	rule, _, err := policysetcontroller.GetPolicyRule(ctx, service, policySetID, plan.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
	}
	_ = rule // Use the retrieved rule if needed

	request, diags := expandPolicyAccessRule(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := helpers.ValidateConditions(ctx, request.Conditions, helperClient, microTenantID); err != nil {
		resp.Diagnostics.AddError("Validation Error", err.Error())
		return
	}

	if _, err := policysetcontroller.UpdateRule(ctx, service, policySetID, plan.ID.ValueString(), request); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update policy access rule: %v", err))
		return
	}

	newState, readDiags := r.readPolicyAccessRule(ctx, service, policySetID, plan.ID.ValueString(), plan.MicrotenantID, &plan)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *PolicyAccessRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy access rules.")
		return
	}

	var state PolicyAccessRuleModel
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, policyAccessType, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	if _, err := policysetcontroller.Delete(ctx, service, policySetID, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete policy access rule: %v", err))
		return
	}
}

func (r *PolicyAccessRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before importing policy access rules.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the policy access rule ID or name.")
		return
	}

	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		rule, _, lookupErr := policysetcontroller.GetByNameAndTypes(ctx, r.client.Service, []string{policyAccessType, "GLOBAL_POLICY"}, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate policy access rule %q: %v", id, lookupErr))
			return
		}
		id = rule.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *PolicyAccessRuleResource) readPolicyAccessRule(ctx context.Context, service *zscaler.Service, policySetID, ruleID string, microTenantID types.String, existingState *PolicyAccessRuleModel) (PolicyAccessRuleModel, diag.Diagnostics) {
	rule, _, err := policysetcontroller.GetPolicyRule(ctx, service, policySetID, ruleID)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return PolicyAccessRuleModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Policy access rule %s not found", ruleID))}
		}
		return PolicyAccessRuleModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read policy access rule: %v", err))}
	}

	return flattenPolicyAccessRule(ctx, rule, policySetID, microTenantID, existingState)
}

func expandPolicyAccessRule(ctx context.Context, model *PolicyAccessRuleModel, policySetID string) (*policysetcontroller.PolicyRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	serverGroups, sgDiags := helpers.ExpandServerGroups(ctx, model.AppServerGroups)
	diags.Append(sgDiags...)

	connectorGroups, cgDiags := helpers.ExpandAppConnectorGroups(ctx, model.AppConnectorGroups)
	diags.Append(cgDiags...)

	conditions, condDiags := helpers.PolicyConditionModelsToSDK(ctx, model.Conditions)
	diags.Append(condDiags...)

	rule := &policysetcontroller.PolicyRule{
		ID:                     helpers.StringValue(model.ID),
		Name:                   helpers.StringValue(model.Name),
		Description:            helpers.StringValue(model.Description),
		Action:                 helpers.StringValue(model.Action),
		ActionID:               helpers.StringValue(model.ActionID),
		BypassDefaultRule:      helpers.BoolValue(model.BypassDefaultRule, false),
		CustomMsg:              helpers.StringValue(model.CustomMsg),
		DefaultRule:            helpers.BoolValue(model.DefaultRule, false),
		Operator:               helpers.StringValue(model.Operator),
		PolicySetID:            policySetID,
		PolicyType:             helpers.StringValue(model.PolicyType),
		Priority:               helpers.StringValue(model.Priority),
		ReauthDefaultRule:      helpers.BoolValue(model.ReauthDefaultRule, false),
		ReauthIdleTimeout:      helpers.StringValue(model.ReauthIdleTimeout),
		ReauthTimeout:          helpers.StringValue(model.ReauthTimeout),
		ZpnIsolationProfileID:  helpers.StringValue(model.ZPNIsolationProfileID),
		ZpnCbiProfileID:        helpers.StringValue(model.ZPNCBIProfileID),
		ZpnInspectionProfileID: helpers.StringValue(model.ZPNInspectionProfileID),
		RuleOrder:              helpers.StringValue(model.RuleOrder),
		MicroTenantID:          helpers.StringValue(model.MicrotenantID),
		LSSDefaultRule:         helpers.BoolValue(model.LSSDefaultRule, false),
		AppServerGroups:        serverGroups,
		AppConnectorGroups:     connectorGroups,
		Conditions:             conditions,
	}

	return rule, diags
}

func flattenPolicyAccessRule(ctx context.Context, rule *policysetcontroller.PolicyRule, policySetID string, microTenantID types.String, existingState *PolicyAccessRuleModel) (PolicyAccessRuleModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	serverGroups, sgDiags := helpers.FlattenServerGroups(ctx, rule.AppServerGroups)
	diags.Append(sgDiags...)

	connectorGroups, cgDiags := helpers.FlattenAppConnectorGroups(ctx, rule.AppConnectorGroups)
	diags.Append(cgDiags...)

	conditions, condDiags := helpers.PolicyConditionsToModels(ctx, rule.Conditions)
	diags.Append(condDiags...)

	// Preserve null values for LSSDefaultRule if it wasn't set in the plan/state
	var lssDefaultRule types.Bool
	if existingState != nil && (existingState.LSSDefaultRule.IsNull() || existingState.LSSDefaultRule.IsUnknown()) {
		// If it was null in the plan/state, keep it null unless the API returned true
		if rule.LSSDefaultRule {
			lssDefaultRule = types.BoolValue(true)
		} else {
			lssDefaultRule = types.BoolNull()
		}
	} else {
		lssDefaultRule = types.BoolValue(rule.LSSDefaultRule)
	}

	state := PolicyAccessRuleModel{
		ID:                     helpers.StringValueOrNull(rule.ID),
		Name:                   helpers.StringValueOrNull(rule.Name),
		Description:            helpers.StringValueOrNull(rule.Description),
		Action:                 helpers.StringValueOrNull(rule.Action),
		ActionID:               helpers.StringValueOrNull(rule.ActionID),
		BypassDefaultRule:      types.BoolValue(rule.BypassDefaultRule),
		CustomMsg:              helpers.StringValueOrNull(rule.CustomMsg),
		DefaultRule:            types.BoolValue(rule.DefaultRule),
		Operator:               helpers.StringValueOrNull(rule.Operator),
		PolicySetID:            types.StringValue(policySetID),
		PolicyType:             helpers.StringValueOrNull(rule.PolicyType),
		Priority:               helpers.StringValueOrNull(rule.Priority),
		ReauthDefaultRule:      types.BoolValue(rule.ReauthDefaultRule),
		ReauthIdleTimeout:      helpers.StringValueOrNull(rule.ReauthIdleTimeout),
		ReauthTimeout:          helpers.StringValueOrNull(rule.ReauthTimeout),
		ZPNIsolationProfileID:  helpers.StringValueOrNull(rule.ZpnIsolationProfileID),
		ZPNCBIProfileID:        helpers.StringValueOrNull(rule.ZpnCbiProfileID),
		ZPNInspectionProfileID: helpers.StringValueOrNull(rule.ZpnInspectionProfileID),
		RuleOrder:              helpers.StringValueOrNull(rule.RuleOrder),
		MicrotenantID:          helpers.StringValueOrNull(rule.MicroTenantID),
		LSSDefaultRule:         lssDefaultRule,
		AppServerGroups:        serverGroups,
		AppConnectorGroups:     connectorGroups,
		Conditions:             conditions,
	}

	if microTenantID != types.StringNull() && !microTenantID.IsUnknown() && !microTenantID.IsNull() {
		state.MicrotenantID = microTenantID
	}

	return state, diags
}

func (r *PolicyAccessRuleResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	id := helpers.StringValue(microtenantID)
	if id != "" {
		service = service.WithMicroTenant(id)
	}
	return service
}
