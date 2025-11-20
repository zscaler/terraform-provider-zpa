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
	_ resource.Resource                = &PolicyAccessTimeoutRuleV2Resource{}
	_ resource.ResourceWithConfigure   = &PolicyAccessTimeoutRuleV2Resource{}
	_ resource.ResourceWithImportState = &PolicyAccessTimeoutRuleV2Resource{}
)

func NewPolicyAccessTimeoutRuleV2Resource() resource.Resource {
	return &PolicyAccessTimeoutRuleV2Resource{}
}

type PolicyAccessTimeoutRuleV2Resource struct {
	client *client.Client
}

type PolicyAccessTimeoutRuleV2Model struct {
	ID                types.String                 `tfsdk:"id"`
	Name              types.String                 `tfsdk:"name"`
	Description       types.String                 `tfsdk:"description"`
	Action            types.String                 `tfsdk:"action"`
	CustomMsg         types.String                 `tfsdk:"custom_msg"`
	PolicySetID       types.String                 `tfsdk:"policy_set_id"`
	ReauthIdleTimeout types.String                 `tfsdk:"reauth_idle_timeout"`
	ReauthTimeout     types.String                 `tfsdk:"reauth_timeout"`
	Conditions        []PolicyAccessConditionModel `tfsdk:"conditions"`
	MicrotenantID     types.String                 `tfsdk:"microtenant_id"`
}

func (r *PolicyAccessTimeoutRuleV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_access_timeout_rule_v2"
}

func (r *PolicyAccessTimeoutRuleV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				stringvalidator.OneOf("RE_AUTH"),
			},
		},
		"custom_msg": schema.StringAttribute{
			Optional: true,
		},
		"policy_set_id": schema.StringAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				fwstringplanmodifier.UseStateForUnknown(),
			},
		},
		"reauth_idle_timeout": schema.StringAttribute{Optional: true},
		"reauth_timeout":      schema.StringAttribute{Optional: true},
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
		"IDP",
		"POSTURE",
		"PLATFORM",
		"SAML",
		"SCIM",
		"SCIM_GROUP",
	}

	resp.Schema = schema.Schema{
		Description: "Manages ZPA Timeout/Re-auth policy rules (v2).",
		Attributes:  attrs,
		Blocks: map[string]schema.Block{
			"conditions": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Computed: true},
						"operator": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf("AND", "OR"),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"operands": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"values": schema.SetAttribute{
										ElementType: types.StringType,
										Optional:    true,
									},
									"object_type": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											stringvalidator.OneOf(objectTypes...),
										},
									},
								},
								Blocks: map[string]schema.Block{
									"entry_values": schema.SetNestedBlock{
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"rhs": schema.StringAttribute{Optional: true},
												"lhs": schema.StringAttribute{Optional: true},
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

func (r *PolicyAccessTimeoutRuleV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PolicyAccessTimeoutRuleV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy timeout rules.")
		return
	}

	var plan PolicyAccessTimeoutRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeTimeout, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	if plan.ReauthIdleTimeout.ValueString() != "" {
		if err := validateTimeoutIntervals(plan.ReauthIdleTimeout.ValueString()); err != nil {
			resp.Diagnostics.AddError("Validation Error", err.Error())
			return
		}
	}
	if plan.ReauthTimeout.ValueString() != "" {
		if err := validateTimeoutIntervals(plan.ReauthTimeout.ValueString()); err != nil {
			resp.Diagnostics.AddError("Validation Error", err.Error())
			return
		}
	}

	reqPayload, diags := expandPolicyAccessTimeoutRuleV2(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := policysetcontrollerv2.CreateRule(ctx, service, reqPayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create policy timeout rule: %v", err))
		return
	}

	state, readDiags := r.readPolicyAccessTimeoutRuleV2(ctx, service, policySetID, created.ID, plan.MicrotenantID, &plan)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PolicyAccessTimeoutRuleV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy timeout rules.")
		return
	}

	var state PolicyAccessTimeoutRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeTimeout, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	newState, diags := r.readPolicyAccessTimeoutRuleV2(ctx, service, policySetID, state.ID.ValueString(), state.MicrotenantID, &state)
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

func (r *PolicyAccessTimeoutRuleV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy timeout rules.")
		return
	}

	var plan PolicyAccessTimeoutRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeTimeout, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	if plan.ReauthIdleTimeout.ValueString() != "" {
		if err := validateTimeoutIntervals(plan.ReauthIdleTimeout.ValueString()); err != nil {
			resp.Diagnostics.AddError("Validation Error", err.Error())
			return
		}
	}
	if plan.ReauthTimeout.ValueString() != "" {
		if err := validateTimeoutIntervals(plan.ReauthTimeout.ValueString()); err != nil {
			resp.Diagnostics.AddError("Validation Error", err.Error())
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

	reqPayload, diags := expandPolicyAccessTimeoutRuleV2(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := policysetcontrollerv2.UpdateRule(ctx, service, policySetID, plan.ID.ValueString(), reqPayload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update policy timeout rule: %v", err))
		return
	}

	newState, readDiags := r.readPolicyAccessTimeoutRuleV2(ctx, service, policySetID, plan.ID.ValueString(), plan.MicrotenantID, &plan)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *PolicyAccessTimeoutRuleV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy timeout rules.")
		return
	}

	var state PolicyAccessTimeoutRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeTimeout, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	if _, err := policysetcontrollerv2.Delete(ctx, service, policySetID, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete policy timeout rule: %v", err))
	}
}

func (r *PolicyAccessTimeoutRuleV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before importing policy timeout rules.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the policy timeout rule ID or name.")
		return
	}

	if _, err := fmt.Sscan(id, new(int64)); err != nil {
		rule, _, lookupErr := policysetcontrollerv2.GetByNameAndTypes(ctx, r.client.Service, []string{helpers.PolicyTypeTimeout, helpers.PolicyTypeReauth}, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate policy timeout rule %q: %v", id, lookupErr))
			return
		}
		id = rule.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *PolicyAccessTimeoutRuleV2Resource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && microtenantID.ValueString() != "" {
		service = service.WithMicroTenant(microtenantID.ValueString())
	}
	return service
}

func (r *PolicyAccessTimeoutRuleV2Resource) readPolicyAccessTimeoutRuleV2(ctx context.Context, service *zscaler.Service, policySetID, ruleID string, microTenantID types.String, existingState *PolicyAccessTimeoutRuleV2Model) (PolicyAccessTimeoutRuleV2Model, diag.Diagnostics) {
	ruleResource, _, err := policysetcontrollerv2.GetPolicyRule(ctx, service, policySetID, ruleID)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return PolicyAccessTimeoutRuleV2Model{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Policy timeout rule %s not found", ruleID))}
		}
		return PolicyAccessTimeoutRuleV2Model{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read policy timeout rule: %v", err))}
	}

	rule := helpers.ConvertV1ResponseToV2Request(*ruleResource)

	conditions, condDiags := flattenPolicyRuleConditionsV2(ctx, rule.Conditions)
	if condDiags.HasError() {
		return PolicyAccessTimeoutRuleV2Model{}, condDiags
	}

	idleTimeout := rule.ReauthIdleTimeout
	if idleTimeout == "-1" {
		idleTimeout = "never"
	} else {
		idleTimeout = helpers.SecondsToHumanReadable(idleTimeout)
	}

	timeout := rule.ReauthTimeout
	if timeout == "-1" {
		timeout = "never"
	} else {
		timeout = helpers.SecondsToHumanReadable(timeout)
	}

	// Ensure microtenant_id is always known (not unknown)
	var microtenantID types.String
	if microTenantID.IsUnknown() {
		microtenantID = types.StringNull()
	} else {
		microtenantID = microTenantID
	}

	// Preserve null values for CustomMsg if it wasn't set in the plan/state
	var customMsg types.String
	if existingState != nil && (existingState.CustomMsg.IsNull() || existingState.CustomMsg.IsUnknown()) {
		// If it was null in the plan/state, keep it null unless the API returned a non-empty value
		if rule.CustomMsg != "" {
			customMsg = types.StringValue(rule.CustomMsg)
		} else {
			customMsg = types.StringNull()
		}
	} else {
		customMsg = types.StringValue(rule.CustomMsg)
	}

	model := PolicyAccessTimeoutRuleV2Model{
		ID:                types.StringValue(rule.ID),
		Name:              types.StringValue(rule.Name),
		Description:       types.StringValue(rule.Description),
		Action:            types.StringValue(rule.Action),
		CustomMsg:         customMsg,
		PolicySetID:       types.StringValue(policySetID),
		ReauthIdleTimeout: types.StringValue(idleTimeout),
		ReauthTimeout:     types.StringValue(timeout),
		Conditions:        conditions,
		MicrotenantID:     microtenantID,
	}

	return model, condDiags
}

func expandPolicyAccessTimeoutRuleV2(ctx context.Context, model *PolicyAccessTimeoutRuleV2Model, policySetID string) (*policysetcontrollerv2.PolicyRule, diag.Diagnostics) {
	conditions, diags := expandPolicyRuleConditionsV2(ctx, model.Conditions)
	if diags.HasError() {
		return nil, diags
	}

	reauthIdleSeconds, err := helpers.ParseHumanReadableTimeout(model.ReauthIdleTimeout.ValueString())
	if err != nil {
		diags.AddError("Validation Error", err.Error())
		return nil, diags
	}

	reauthSeconds, err := helpers.ParseHumanReadableTimeout(model.ReauthTimeout.ValueString())
	if err != nil {
		diags.AddError("Validation Error", err.Error())
		return nil, diags
	}

	rule := &policysetcontrollerv2.PolicyRule{
		ID:                helpers.StringValue(model.ID),
		Name:              helpers.StringValue(model.Name),
		Description:       helpers.StringValue(model.Description),
		CustomMsg:         helpers.StringValue(model.CustomMsg),
		Action:            helpers.StringValue(model.Action),
		PolicySetID:       policySetID,
		ReauthIdleTimeout: strconv.Itoa(reauthIdleSeconds),
		ReauthTimeout:     strconv.Itoa(reauthSeconds),
		Conditions:        conditions,
		MicroTenantID:     helpers.StringValue(model.MicrotenantID),
	}

	return rule, diags
}

func validateTimeoutIntervals(input string) error {
	if strings.TrimSpace(strings.ToLower(input)) == "never" || input == "" {
		return nil
	}

	seconds, err := helpers.ParseHumanReadableTimeout(input)
	if err != nil {
		return err
	}

	if seconds >= 0 && seconds < 600 {
		return fmt.Errorf("timeout interval must be at least 10 minutes or 'never'")
	}

	return nil
}
