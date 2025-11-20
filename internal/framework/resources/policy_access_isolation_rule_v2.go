package resources

import (
	"context"
	"fmt"
	"strings"

	stringvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	fwstringplanmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

var (
	_ resource.Resource                = &PolicyAccessIsolationRuleV2Resource{}
	_ resource.ResourceWithConfigure   = &PolicyAccessIsolationRuleV2Resource{}
	_ resource.ResourceWithImportState = &PolicyAccessIsolationRuleV2Resource{}
)

func NewPolicyAccessIsolationRuleV2Resource() resource.Resource {
	return &PolicyAccessIsolationRuleV2Resource{}
}

type PolicyAccessIsolationRuleV2Resource struct {
	client *client.Client
}

type PolicyAccessIsolationRuleV2Model struct {
	ID                    types.String                 `tfsdk:"id"`
	Name                  types.String                 `tfsdk:"name"`
	Description           types.String                 `tfsdk:"description"`
	Action                types.String                 `tfsdk:"action"`
	PolicySetID           types.String                 `tfsdk:"policy_set_id"`
	ZpnIsolationProfileID types.String                 `tfsdk:"zpn_isolation_profile_id"`
	Conditions            []PolicyAccessConditionModel `tfsdk:"conditions"`
	MicrotenantID         types.String                 `tfsdk:"microtenant_id"`
}

func (r *PolicyAccessIsolationRuleV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_access_isolation_rule_v2"
}

func (r *PolicyAccessIsolationRuleV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	attrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				fwstringplanmodifier.UseStateForUnknown(),
			},
		},
		"name":        schema.StringAttribute{Required: true},
		"description": schema.StringAttribute{Optional: true},
		"action": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("ISOLATE", "BYPASS_ISOLATE"),
			},
		},
		"policy_set_id": schema.StringAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				fwstringplanmodifier.UseStateForUnknown(),
			},
		},
		"zpn_isolation_profile_id": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"microtenant_id": schema.StringAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				fwstringplanmodifier.UseStateForUnknown(),
			},
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
		"CHROME_ENTERPRISE",
		"CHROME_POSTURE_PROFILE",
	}

	resp.Schema = schema.Schema{
		Description: "Manages ZPA Isolation policy rules (v2).",
		Attributes:  attrs,
		Blocks: map[string]schema.Block{
			"conditions": helpers.PolicyConditionsV2Block(objectTypes),
		},
	}
}

func (r *PolicyAccessIsolationRuleV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PolicyAccessIsolationRuleV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy isolation rules.")
		return
	}

	var plan PolicyAccessIsolationRuleV2Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateObjectTypeUniquenessV2(ctx, plan.Conditions)...)
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeIsolation, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	reqPayload, diags := expandPolicyAccessIsolationRuleV2(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := policysetcontrollerv2.CreateRule(ctx, service, reqPayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create policy isolation rule: %v", err))
		return
	}

	state, readDiags := r.readPolicyAccessIsolationRuleV2(ctx, service, policySetID, created.ID, plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PolicyAccessIsolationRuleV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy isolation rules.")
		return
	}

	var state PolicyAccessIsolationRuleV2Model
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

	newState, diags := r.readPolicyAccessIsolationRuleV2(ctx, service, policySetID, state.ID.ValueString(), state.MicrotenantID)
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

func (r *PolicyAccessIsolationRuleV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy isolation rules.")
		return
	}

	var plan PolicyAccessIsolationRuleV2Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	resp.Diagnostics.Append(validateObjectTypeUniquenessV2(ctx, plan.Conditions)...)
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeIsolation, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	// Check if resource still exists before updating
	ruleResource, _, err := policysetcontrollerv2.GetPolicyRule(ctx, service, policySetID, plan.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
	}
	_ = ruleResource // Use the retrieved rule if needed

	reqPayload, diags := expandPolicyAccessIsolationRuleV2(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := policysetcontrollerv2.UpdateRule(ctx, service, policySetID, plan.ID.ValueString(), reqPayload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update policy isolation rule: %v", err))
		return
	}

	newState, readDiags := r.readPolicyAccessIsolationRuleV2(ctx, service, policySetID, plan.ID.ValueString(), plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *PolicyAccessIsolationRuleV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy isolation rules.")
		return
	}

	var state PolicyAccessIsolationRuleV2Model
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

	if _, err := policysetcontrollerv2.Delete(ctx, service, policySetID, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete policy isolation rule: %v", err))
	}
}

func (r *PolicyAccessIsolationRuleV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before importing policy isolation rules.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the policy isolation rule ID or name.")
		return
	}

	if _, err := fmt.Sscan(id, new(int64)); err != nil {
		rule, _, lookupErr := policysetcontrollerv2.GetByNameAndTypes(ctx, r.client.Service, []string{helpers.PolicyTypeIsolation}, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate policy isolation rule %q: %v", id, lookupErr))
			return
		}
		id = rule.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *PolicyAccessIsolationRuleV2Resource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && microtenantID.ValueString() != "" {
		service = service.WithMicroTenant(microtenantID.ValueString())
	}
	return service
}

func (r *PolicyAccessIsolationRuleV2Resource) readPolicyAccessIsolationRuleV2(ctx context.Context, service *zscaler.Service, policySetID, ruleID string, microTenantID types.String) (PolicyAccessIsolationRuleV2Model, diag.Diagnostics) {
	ruleResource, _, err := policysetcontrollerv2.GetPolicyRule(ctx, service, policySetID, ruleID)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return PolicyAccessIsolationRuleV2Model{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Policy isolation rule %s not found", ruleID))}
		}
		return PolicyAccessIsolationRuleV2Model{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read policy isolation rule: %v", err))}
	}

	rule := helpers.ConvertV1ResponseToV2Request(*ruleResource)

	conditions, condDiags := flattenPolicyRuleConditionsV2(ctx, rule.Conditions)
	if condDiags.HasError() {
		return PolicyAccessIsolationRuleV2Model{}, condDiags
	}

	model := PolicyAccessIsolationRuleV2Model{
		ID:                    types.StringValue(rule.ID),
		Name:                  types.StringValue(rule.Name),
		Description:           types.StringValue(rule.Description),
		Action:                types.StringValue(rule.Action),
		PolicySetID:           types.StringValue(policySetID),
		ZpnIsolationProfileID: types.StringValue(rule.ZpnIsolationProfileID),
		Conditions:            conditions,
		MicrotenantID:         microTenantID,
	}

	return model, condDiags
}

func expandPolicyAccessIsolationRuleV2(ctx context.Context, model *PolicyAccessIsolationRuleV2Model, policySetID string) (*policysetcontrollerv2.PolicyRule, diag.Diagnostics) {
	conditions, diags := expandPolicyRuleConditionsV2(ctx, model.Conditions)
	if diags.HasError() {
		return nil, diags
	}

	rule := &policysetcontrollerv2.PolicyRule{
		ID:                    helpers.StringValue(model.ID),
		Name:                  helpers.StringValue(model.Name),
		Description:           helpers.StringValue(model.Description),
		Action:                helpers.StringValue(model.Action),
		PolicySetID:           policySetID,
		ZpnIsolationProfileID: helpers.StringValue(model.ZpnIsolationProfileID),
		MicroTenantID:         helpers.StringValue(model.MicrotenantID),
		Conditions:            conditions,
	}

	return rule, diags
}
