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
	_ resource.Resource                = &PolicyAccessForwardingRuleV2Resource{}
	_ resource.ResourceWithConfigure   = &PolicyAccessForwardingRuleV2Resource{}
	_ resource.ResourceWithImportState = &PolicyAccessForwardingRuleV2Resource{}
)

func NewPolicyAccessForwardingRuleV2Resource() resource.Resource {
	return &PolicyAccessForwardingRuleV2Resource{}
}

type PolicyAccessForwardingRuleV2Resource struct {
	client *client.Client
}

type PolicyAccessForwardingRuleV2Model struct {
	ID            types.String                 `tfsdk:"id"`
	Name          types.String                 `tfsdk:"name"`
	Description   types.String                 `tfsdk:"description"`
	Action        types.String                 `tfsdk:"action"`
	PolicySetID   types.String                 `tfsdk:"policy_set_id"`
	Conditions    []PolicyAccessConditionModel `tfsdk:"conditions"`
	MicrotenantID types.String                 `tfsdk:"microtenant_id"`
}

func (r *PolicyAccessForwardingRuleV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_access_forwarding_rule_v2"
}

func (r *PolicyAccessForwardingRuleV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				stringvalidator.OneOf("BYPASS", "INTERCEPT", "INTERCEPT_ACCESSIBLE"),
			},
		},
		"policy_set_id": schema.StringAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				fwstringplanmodifier.UseStateForUnknown(),
			},
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
		"BRANCH_CONNECTOR_GROUP",
		"EDGE_CONNECTOR_GROUP",
		"POSTURE",
		"MACHINE_GRP",
		"TRUSTED_NETWORK",
		"PLATFORM",
		"IDP",
		"SAML",
		"SCIM",
		"SCIM_GROUP",
	}

	resp.Schema = schema.Schema{
		Description: "Manages ZPA Client Forwarding policy rules (v2).",
		Attributes:  attrs,
		Blocks: map[string]schema.Block{
			"conditions": helpers.PolicyConditionsV2Block(objectTypes),
		},
	}
}

func (r *PolicyAccessForwardingRuleV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PolicyAccessForwardingRuleV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy forwarding rules.")
		return
	}

	var plan PolicyAccessForwardingRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeClientForwarding, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	reqPayload, diags := expandPolicyAccessForwardingRuleV2(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := policysetcontrollerv2.CreateRule(ctx, service, reqPayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create policy forwarding rule: %v", err))
		return
	}

	state, readDiags := r.readPolicyAccessForwardingRuleV2(ctx, service, policySetID, created.ID, plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PolicyAccessForwardingRuleV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy forwarding rules.")
		return
	}

	var state PolicyAccessForwardingRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeClientForwarding, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	newState, diags := r.readPolicyAccessForwardingRuleV2(ctx, service, policySetID, state.ID.ValueString(), state.MicrotenantID)
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

func (r *PolicyAccessForwardingRuleV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy forwarding rules.")
		return
	}

	var plan PolicyAccessForwardingRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeClientForwarding, microTenantID)
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

	reqPayload, diags := expandPolicyAccessForwardingRuleV2(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := policysetcontrollerv2.UpdateRule(ctx, service, policySetID, plan.ID.ValueString(), reqPayload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update policy forwarding rule: %v", err))
		return
	}

	newState, readDiags := r.readPolicyAccessForwardingRuleV2(ctx, service, policySetID, plan.ID.ValueString(), plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *PolicyAccessForwardingRuleV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy forwarding rules.")
		return
	}

	var state PolicyAccessForwardingRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeClientForwarding, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	if _, err := policysetcontrollerv2.Delete(ctx, service, policySetID, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete policy forwarding rule: %v", err))
	}
}

func (r *PolicyAccessForwardingRuleV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before importing policy forwarding rules.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the policy forwarding rule ID or name.")
		return
	}

	if _, err := fmt.Sscan(id, new(int64)); err != nil {
		rule, _, lookupErr := policysetcontrollerv2.GetByNameAndTypes(ctx, r.client.Service, []string{helpers.PolicyTypeClientForwarding, helpers.PolicyTypeBypass}, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate policy forwarding rule %q: %v", id, lookupErr))
			return
		}
		id = rule.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *PolicyAccessForwardingRuleV2Resource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && microtenantID.ValueString() != "" {
		service = service.WithMicroTenant(microtenantID.ValueString())
	}
	return service
}

func (r *PolicyAccessForwardingRuleV2Resource) readPolicyAccessForwardingRuleV2(ctx context.Context, service *zscaler.Service, policySetID, ruleID string, microTenantID types.String) (PolicyAccessForwardingRuleV2Model, diag.Diagnostics) {
	ruleResource, _, err := policysetcontrollerv2.GetPolicyRule(ctx, service, policySetID, ruleID)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return PolicyAccessForwardingRuleV2Model{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Policy forwarding rule %s not found", ruleID))}
		}
		return PolicyAccessForwardingRuleV2Model{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read policy forwarding rule: %v", err))}
	}

	rule := helpers.ConvertV1ResponseToV2Request(*ruleResource)

	conditions, condDiags := flattenPolicyRuleConditionsV2(ctx, rule.Conditions)
	if condDiags.HasError() {
		return PolicyAccessForwardingRuleV2Model{}, condDiags
	}

	// Ensure microtenant_id is always known (not unknown)
	var microtenantID types.String
	if microTenantID.IsUnknown() {
		microtenantID = types.StringNull()
	} else {
		microtenantID = microTenantID
	}

	model := PolicyAccessForwardingRuleV2Model{
		ID:            types.StringValue(rule.ID),
		Name:          types.StringValue(rule.Name),
		Description:   types.StringValue(rule.Description),
		Action:        types.StringValue(rule.Action),
		PolicySetID:   types.StringValue(policySetID),
		Conditions:    conditions,
		MicrotenantID: microtenantID,
	}

	return model, condDiags
}

func expandPolicyAccessForwardingRuleV2(ctx context.Context, model *PolicyAccessForwardingRuleV2Model, policySetID string) (*policysetcontrollerv2.PolicyRule, diag.Diagnostics) {
	conditions, diags := expandPolicyRuleConditionsV2(ctx, model.Conditions)
	if diags.HasError() {
		return nil, diags
	}

	rule := &policysetcontrollerv2.PolicyRule{
		ID:            helpers.StringValue(model.ID),
		Name:          helpers.StringValue(model.Name),
		Description:   helpers.StringValue(model.Description),
		Action:        helpers.StringValue(model.Action),
		PolicySetID:   policySetID,
		MicroTenantID: helpers.StringValue(model.MicrotenantID),
		Conditions:    conditions,
	}

	return rule, diags
}
