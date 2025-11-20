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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

var policyIsolationImportTypes = []string{helpers.PolicyTypeIsolation}

var (
	_ resource.Resource                = &PolicyAccessIsolationRuleResource{}
	_ resource.ResourceWithConfigure   = &PolicyAccessIsolationRuleResource{}
	_ resource.ResourceWithImportState = &PolicyAccessIsolationRuleResource{}
)

func NewPolicyAccessIsolationRuleResource() resource.Resource {
	return &PolicyAccessIsolationRuleResource{}
}

type PolicyAccessIsolationRuleResource struct {
	client *client.Client
}

type PolicyAccessIsolationRuleModel struct {
	ID                     types.String                   `tfsdk:"id"`
	Name                   types.String                   `tfsdk:"name"`
	Description            types.String                   `tfsdk:"description"`
	Action                 types.String                   `tfsdk:"action"`
	ActionID               types.String                   `tfsdk:"action_id"`
	CustomMsg              types.String                   `tfsdk:"custom_msg"`
	BypassDefaultRule      types.Bool                     `tfsdk:"bypass_default_rule"`
	DefaultRule            types.Bool                     `tfsdk:"default_rule"`
	Operator               types.String                   `tfsdk:"operator"`
	PolicySetID            types.String                   `tfsdk:"policy_set_id"`
	PolicyType             types.String                   `tfsdk:"policy_type"`
	Priority               types.String                   `tfsdk:"priority"`
	ReauthDefaultRule      types.Bool                     `tfsdk:"reauth_default_rule"`
	ReauthIdleTimeout      types.String                   `tfsdk:"reauth_idle_timeout"`
	ReauthTimeout          types.String                   `tfsdk:"reauth_timeout"`
	ZPNIsolationProfileID  types.String                   `tfsdk:"zpn_isolation_profile_id"`
	ZPNInspectionProfileID types.String                   `tfsdk:"zpn_inspection_profile_id"`
	ZPNCBIProfileID        types.String                   `tfsdk:"zpn_cbi_profile_id"`
	RuleOrder              types.String                   `tfsdk:"rule_order"`
	MicrotenantID          types.String                   `tfsdk:"microtenant_id"`
	LSSDefaultRule         types.Bool                     `tfsdk:"lss_default_rule"`
	Conditions             []helpers.PolicyConditionModel `tfsdk:"conditions"`
}

func (r *PolicyAccessIsolationRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_access_isolation_rule"
}

func (r *PolicyAccessIsolationRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := helpers.PolicyCommonSchemaAttributes()

	attributes["action"] = schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf("ISOLATE", "BYPASS_ISOLATE"),
		},
	}

	objectTypes := []string{
		"APP",
		"APP_GROUP",
		"CLIENT_TYPE",
		"EDGE_CONNECTOR_GROUP",
		"PLATFORM",
		"IDP",
		"SAML",
		"SCIM",
		"SCIM_GROUP",
	}

	resp.Schema = schema.Schema{
		Description: "Manages ZPA isolation policy rules.",
		Attributes:  attributes,
		Blocks: map[string]schema.Block{
			"conditions": helpers.PolicyConditionsBlock(objectTypes),
		},
	}
}

func (r *PolicyAccessIsolationRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PolicyAccessIsolationRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing isolation policy rules.")
		return
	}

	var plan PolicyAccessIsolationRuleModel
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeIsolation, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	payload, diags := expandPolicyAccessIsolationRule(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := helpers.ValidateConditions(ctx, payload.Conditions, helperClient, microTenantID); err != nil {
		resp.Diagnostics.AddError("Validation Error", err.Error())
		return
	}

	created, _, err := policysetcontroller.CreateRule(ctx, service, payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create isolation policy rule: %v", err))
		return
	}

	state, readDiags := r.readPolicyAccessIsolationRule(ctx, service, policySetID, created.ID, plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PolicyAccessIsolationRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing isolation policy rules.")
		return
	}

	var state PolicyAccessIsolationRuleModel
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeIsolation, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	newState, diags := r.readPolicyAccessIsolationRule(ctx, service, policySetID, helpers.StringValue(state.ID), state.MicrotenantID)
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

func (r *PolicyAccessIsolationRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing isolation policy rules.")
		return
	}

	var plan PolicyAccessIsolationRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.IsUnknown() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)
	helperClient := helpers.NewHelperClient(r.client)

	policySetID := helpers.StringValue(plan.PolicySetID)
	microTenantID := helpers.StringValue(plan.MicrotenantID)
	if policySetID == "" {
		var err error
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeIsolation, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	payload, diags := expandPolicyAccessIsolationRule(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := helpers.ValidateConditions(ctx, payload.Conditions, helperClient, microTenantID); err != nil {
		resp.Diagnostics.AddError("Validation Error", err.Error())
		return
	}

	if _, err := policysetcontroller.UpdateRule(ctx, service, policySetID, plan.ID.ValueString(), payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update isolation policy rule: %v", err))
		return
	}

	newState, readDiags := r.readPolicyAccessIsolationRule(ctx, service, policySetID, plan.ID.ValueString(), plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *PolicyAccessIsolationRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing isolation policy rules.")
		return
	}

	var state PolicyAccessIsolationRuleModel
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeIsolation, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	if _, err := policysetcontroller.Delete(ctx, service, policySetID, helpers.StringValue(state.ID)); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete isolation policy rule: %v", err))
		return
	}
}

func (r *PolicyAccessIsolationRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before importing isolation policy rules.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the isolation policy rule ID or name.")
		return
	}

	if _, err := strconv.ParseInt(id, 10, 64); err == nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		return
	}

	rule, _, err := policysetcontroller.GetByNameAndTypes(ctx, r.client.Service, policyIsolationImportTypes, id)
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate isolation policy rule %q: %v", id, err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(rule.ID))...)
}

func (r *PolicyAccessIsolationRuleResource) readPolicyAccessIsolationRule(ctx context.Context, service *zscaler.Service, policySetID, ruleID string, microtenantID types.String) (PolicyAccessIsolationRuleModel, diag.Diagnostics) {
	rule, _, err := policysetcontroller.GetPolicyRule(ctx, service, policySetID, ruleID)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return PolicyAccessIsolationRuleModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Isolation policy rule %s not found", ruleID))}
		}
		return PolicyAccessIsolationRuleModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read isolation policy rule: %v", err))}
	}

	return flattenPolicyAccessIsolationRule(ctx, rule, policySetID, microtenantID)
}

func (r *PolicyAccessIsolationRuleResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	id := helpers.StringValue(microtenantID)
	if id != "" {
		service = service.WithMicroTenant(id)
	}
	return service
}

func expandPolicyAccessIsolationRule(ctx context.Context, model *PolicyAccessIsolationRuleModel, policySetID string) (*policysetcontroller.PolicyRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	conditions, condDiags := helpers.PolicyConditionModelsToSDK(ctx, model.Conditions)
	diags.Append(condDiags...)

	payload := &policysetcontroller.PolicyRule{
		ID:                     helpers.StringValue(model.ID),
		Name:                   helpers.StringValue(model.Name),
		Description:            helpers.StringValue(model.Description),
		Action:                 helpers.StringValue(model.Action),
		ActionID:               helpers.StringValue(model.ActionID),
		CustomMsg:              helpers.StringValue(model.CustomMsg),
		BypassDefaultRule:      helpers.BoolValueDefaultFalse(model.BypassDefaultRule),
		DefaultRule:            helpers.BoolValueDefaultFalse(model.DefaultRule),
		Operator:               helpers.StringValue(model.Operator),
		PolicySetID:            policySetID,
		PolicyType:             helpers.StringValue(model.PolicyType),
		Priority:               helpers.StringValue(model.Priority),
		ReauthDefaultRule:      helpers.BoolValueDefaultFalse(model.ReauthDefaultRule),
		ReauthIdleTimeout:      helpers.StringValue(model.ReauthIdleTimeout),
		ReauthTimeout:          helpers.StringValue(model.ReauthTimeout),
		ZpnIsolationProfileID:  helpers.StringValue(model.ZPNIsolationProfileID),
		ZpnInspectionProfileID: helpers.StringValue(model.ZPNInspectionProfileID),
		ZpnCbiProfileID:        helpers.StringValue(model.ZPNCBIProfileID),
		MicroTenantID:          helpers.StringValue(model.MicrotenantID),
		LSSDefaultRule:         helpers.BoolValueDefaultFalse(model.LSSDefaultRule),
		Conditions:             conditions,
	}

	return payload, diags
}

func flattenPolicyAccessIsolationRule(ctx context.Context, rule *policysetcontroller.PolicyRule, policySetID string, microtenantID types.String) (PolicyAccessIsolationRuleModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	conditions, condDiags := helpers.PolicyConditionsToModels(ctx, rule.Conditions)
	diags.Append(condDiags...)

	state := PolicyAccessIsolationRuleModel{
		ID:                     helpers.StringValueOrNull(rule.ID),
		Name:                   helpers.StringValueOrNull(rule.Name),
		Description:            helpers.StringValueOrNull(rule.Description),
		Action:                 helpers.StringValueOrNull(rule.Action),
		ActionID:               helpers.StringValueOrNull(rule.ActionID),
		CustomMsg:              helpers.StringValueOrNull(rule.CustomMsg),
		BypassDefaultRule:      types.BoolValue(rule.BypassDefaultRule),
		DefaultRule:            types.BoolValue(rule.DefaultRule),
		Operator:               helpers.StringValueOrNull(rule.Operator),
		PolicySetID:            types.StringValue(policySetID),
		PolicyType:             helpers.StringValueOrNull(rule.PolicyType),
		Priority:               helpers.StringValueOrNull(rule.Priority),
		ReauthDefaultRule:      types.BoolValue(rule.ReauthDefaultRule),
		ReauthIdleTimeout:      helpers.StringValueOrNull(rule.ReauthIdleTimeout),
		ReauthTimeout:          helpers.StringValueOrNull(rule.ReauthTimeout),
		ZPNIsolationProfileID:  helpers.StringValueOrNull(rule.ZpnIsolationProfileID),
		ZPNInspectionProfileID: helpers.StringValueOrNull(rule.ZpnInspectionProfileID),
		ZPNCBIProfileID:        helpers.StringValueOrNull(rule.ZpnCbiProfileID),
		RuleOrder:              helpers.StringValueOrNull(rule.RuleOrder),
		MicrotenantID:          helpers.StringValueOrNull(rule.MicroTenantID),
		LSSDefaultRule:         types.BoolValue(rule.LSSDefaultRule),
		Conditions:             conditions,
	}

	if microtenantID != types.StringNull() && !microtenantID.IsNull() && !microtenantID.IsUnknown() {
		state.MicrotenantID = microtenantID
	}

	return state, diags
}
