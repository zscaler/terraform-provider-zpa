package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
	stringvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	fwstringplanmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

var (
	_ resource.Resource                = &PolicyAccessCredentialRuleV2Resource{}
	_ resource.ResourceWithConfigure   = &PolicyAccessCredentialRuleV2Resource{}
	_ resource.ResourceWithImportState = &PolicyAccessCredentialRuleV2Resource{}
)

func NewPolicyAccessCredentialRuleV2Resource() resource.Resource {
	return &PolicyAccessCredentialRuleV2Resource{}
}

type PolicyAccessCredentialRuleV2Resource struct {
	client *client.Client
}

type PolicyAccessCredentialRuleV2Model struct {
	ID             types.String                 `tfsdk:"id"`
	Name           types.String                 `tfsdk:"name"`
	Description    types.String                 `tfsdk:"description"`
	Action         types.String                 `tfsdk:"action"`
	PolicySetID    types.String                 `tfsdk:"policy_set_id"`
	Conditions     []PolicyAccessConditionModel `tfsdk:"conditions"`
	Credential     types.List                   `tfsdk:"credential"`
	CredentialPool types.List                   `tfsdk:"credential_pool"`
	MicrotenantID  types.String                 `tfsdk:"microtenant_id"`
}

var credentialAttrTypes = map[string]attr.Type{
	"id": types.StringType,
}

func (r *PolicyAccessCredentialRuleV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_access_credential_rule_v2"
}

func (r *PolicyAccessCredentialRuleV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				stringvalidator.OneOf("INJECT_CREDENTIALS"),
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
		"CONSOLE",
		"SAML",
		"SCIM",
		"SCIM_GROUP",
	}

	resp.Schema = schema.Schema{
		Description: "Manages ZPA Credential Injection policy rules (v2).",
		Attributes:  attrs,
		Blocks: map[string]schema.Block{
			"conditions": helpers.PolicyConditionsV2Block(objectTypes),
			"credential": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Optional: true},
					},
				},
			},
			"credential_pool": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Optional: true},
					},
				},
			},
		},
	}
}

func (r *PolicyAccessCredentialRuleV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PolicyAccessCredentialRuleV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy credential rules.")
		return
	}

	var plan PolicyAccessCredentialRuleV2Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateObjectTypeUniquenessV2(ctx, plan.Conditions)...)
	resp.Diagnostics.Append(validatePolicyRuleConditionsV2(ctx, plan.Conditions)...)
	resp.Diagnostics.Append(validateCredentialSelection(plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)
	helperClient := helpers.NewHelperClient(r.client)

	policySetID := helpers.StringValue(plan.PolicySetID)
	microTenantID := helpers.StringValue(plan.MicrotenantID)
	if policySetID == "" {
		var err error
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeCredential, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	request, diags := expandPolicyAccessCredentialRuleV2(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := policysetcontrollerv2.CreateRule(ctx, service, request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create policy credential rule: %v", err))
		return
	}

	state, readDiags := r.readPolicyAccessCredentialRuleV2(ctx, service, policySetID, created.ID, plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PolicyAccessCredentialRuleV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy credential rules.")
		return
	}

	var state PolicyAccessCredentialRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeCredential, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	newState, diags := r.readPolicyAccessCredentialRuleV2(ctx, service, policySetID, state.ID.ValueString(), state.MicrotenantID)
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

func (r *PolicyAccessCredentialRuleV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy credential rules.")
		return
	}

	var plan PolicyAccessCredentialRuleV2Model
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
	resp.Diagnostics.Append(validateCredentialSelection(plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)
	helperClient := helpers.NewHelperClient(r.client)

	policySetID := helpers.StringValue(plan.PolicySetID)
	microTenantID := helpers.StringValue(plan.MicrotenantID)
	if policySetID == "" {
		var err error
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeCredential, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	request, diags := expandPolicyAccessCredentialRuleV2(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := policysetcontrollerv2.UpdateRule(ctx, service, policySetID, plan.ID.ValueString(), request); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update policy credential rule: %v", err))
		return
	}

	newState, readDiags := r.readPolicyAccessCredentialRuleV2(ctx, service, policySetID, plan.ID.ValueString(), plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *PolicyAccessCredentialRuleV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy credential rules.")
		return
	}

	var state PolicyAccessCredentialRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeCredential, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	if _, err := policysetcontrollerv2.Delete(ctx, service, policySetID, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete policy credential rule: %v", err))
	}
}

func (r *PolicyAccessCredentialRuleV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before importing policy credential rules.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the policy credential rule ID or name.")
		return
	}

	if _, err := fmt.Sscan(id, new(int64)); err != nil {
		rule, _, lookupErr := policysetcontrollerv2.GetByNameAndTypes(ctx, r.client.Service, []string{helpers.PolicyTypeCredential}, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate policy credential rule %q: %v", id, lookupErr))
			return
		}
		id = rule.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *PolicyAccessCredentialRuleV2Resource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && microtenantID.ValueString() != "" {
		service = service.WithMicroTenant(microtenantID.ValueString())
	}
	return service
}

func (r *PolicyAccessCredentialRuleV2Resource) readPolicyAccessCredentialRuleV2(ctx context.Context, service *zscaler.Service, policySetID, ruleID string, microTenantID types.String) (PolicyAccessCredentialRuleV2Model, diag.Diagnostics) {
	ruleResource, _, err := policysetcontrollerv2.GetPolicyRule(ctx, service, policySetID, ruleID)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return PolicyAccessCredentialRuleV2Model{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Policy credential rule %s not found", ruleID))}
		}
		return PolicyAccessCredentialRuleV2Model{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read policy credential rule: %v", err))}
	}

	rule := helpers.ConvertV1ResponseToV2Request(*ruleResource)

	conditions, condDiags := flattenPolicyRuleConditionsV2(ctx, rule.Conditions)
	if condDiags.HasError() {
		return PolicyAccessCredentialRuleV2Model{}, condDiags
	}

	credential := flattenCredentialToList(ctx, rule.Credential)
	credentialPool := flattenCredentialToList(ctx, rule.CredentialPool)

	model := PolicyAccessCredentialRuleV2Model{
		ID:             types.StringValue(rule.ID),
		Name:           types.StringValue(rule.Name),
		Description:    types.StringValue(rule.Description),
		Action:         types.StringValue(rule.Action),
		PolicySetID:    types.StringValue(policySetID),
		Conditions:     conditions,
		Credential:     credential,
		CredentialPool: credentialPool,
		MicrotenantID:  microTenantID,
	}

	return model, condDiags
}

func expandPolicyAccessCredentialRuleV2(ctx context.Context, model *PolicyAccessCredentialRuleV2Model, policySetID string) (*policysetcontrollerv2.PolicyRule, diag.Diagnostics) {
	conditions, diags := expandPolicyRuleConditionsV2(ctx, model.Conditions)
	if diags.HasError() {
		return nil, diags
	}

	credential := expandCredentialFromList(ctx, model.Credential)
	credentialPool := expandCredentialFromList(ctx, model.CredentialPool)

	rule := &policysetcontrollerv2.PolicyRule{
		ID:             helpers.StringValue(model.ID),
		Name:           helpers.StringValue(model.Name),
		Description:    helpers.StringValue(model.Description),
		Action:         helpers.StringValue(model.Action),
		PolicySetID:    policySetID,
		Conditions:     conditions,
		Credential:     credential,
		CredentialPool: credentialPool,
		MicroTenantID:  helpers.StringValue(model.MicrotenantID),
	}

	return rule, diags
}

func flattenCredentialToList(ctx context.Context, credential *policysetcontrollerv2.Credential) types.List {
	if credential == nil || credential.ID == "" {
		return types.ListNull(types.ObjectType{AttrTypes: credentialAttrTypes})
	}

	obj, diags := types.ObjectValue(credentialAttrTypes, map[string]attr.Value{
		"id": types.StringValue(credential.ID),
	})
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: credentialAttrTypes})
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: credentialAttrTypes}, []attr.Value{obj})
	if listDiags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: credentialAttrTypes})
	}
	return list
}

func expandCredentialFromList(ctx context.Context, list types.List) *policysetcontrollerv2.Credential {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var items []struct {
		ID types.String `tfsdk:"id"`
	}
	if diags := list.ElementsAs(ctx, &items, false); diags.HasError() {
		return nil
	}
	if len(items) == 0 {
		return nil
	}

	id := helpers.StringValue(items[0].ID)
	if id == "" {
		return nil
	}

	return &policysetcontrollerv2.Credential{ID: id}
}

func validateCredentialSelection(model PolicyAccessCredentialRuleV2Model) diag.Diagnostics {
	var diags diag.Diagnostics

	hasCredential := listHasElements(model.Credential)
	hasCredentialPool := listHasElements(model.CredentialPool)

	if hasCredential && hasCredentialPool {
		diags.AddError("Invalid configuration", "Only one of `credential` or `credential_pool` can be specified.")
		return diags
	}

	if !hasCredential && !hasCredentialPool {
		diags.AddError("Invalid configuration", "One of `credential` or `credential_pool` must be specified.")
	}

	return diags
}

func listHasElements(list types.List) bool {
	return !list.IsNull() && !list.IsUnknown() && len(list.Elements()) > 0
}
