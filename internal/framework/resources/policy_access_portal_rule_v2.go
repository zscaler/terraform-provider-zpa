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
	_ resource.Resource                = &PolicyAccessPortalRuleV2Resource{}
	_ resource.ResourceWithConfigure   = &PolicyAccessPortalRuleV2Resource{}
	_ resource.ResourceWithImportState = &PolicyAccessPortalRuleV2Resource{}
)

func NewPolicyAccessPortalRuleV2Resource() resource.Resource {
	return &PolicyAccessPortalRuleV2Resource{}
}

type PolicyAccessPortalRuleV2Resource struct {
	client *client.Client
}

type PolicyAccessPortalRuleV2Model struct {
	ID                           types.String                 `tfsdk:"id"`
	Name                         types.String                 `tfsdk:"name"`
	Description                  types.String                 `tfsdk:"description"`
	Action                       types.String                 `tfsdk:"action"`
	PolicySetID                  types.String                 `tfsdk:"policy_set_id"`
	Conditions                   []PolicyAccessConditionModel `tfsdk:"conditions"`
	PrivilegedPortalCapabilities types.List                   `tfsdk:"privileged_portal_capabilities"`
	MicrotenantID                types.String                 `tfsdk:"microtenant_id"`
}

func (r *PolicyAccessPortalRuleV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_access_portal_rule_v2"
}

var portalCapabilitiesAttrTypes = map[string]attr.Type{
	"delete_file":             types.BoolType,
	"access_uninspected_file": types.BoolType,
	"request_approvals":       types.BoolType,
	"review_approvals":        types.BoolType,
}

func (r *PolicyAccessPortalRuleV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				stringvalidator.OneOf("CHECK_PRIVILEGED_PORTAL_CAPABILITIES"),
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
		"COUNTRY_CODE",
		"PRIVILEGE_PORTAL",
		"SAML",
		"SCIM",
		"SCIM_GROUP",
	}

	resp.Schema = schema.Schema{
		Description: "Manages ZPA Privileged Portal policy rules (v2).",
		Attributes:  attrs,
		Blocks: map[string]schema.Block{
			"conditions": helpers.PolicyConditionsV2Block(objectTypes),
			"privileged_portal_capabilities": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"delete_file":             schema.BoolAttribute{Optional: true},
						"access_uninspected_file": schema.BoolAttribute{Optional: true},
						"request_approvals":       schema.BoolAttribute{Optional: true},
						"review_approvals":        schema.BoolAttribute{Optional: true},
					},
				},
			},
		},
	}
}

func (r *PolicyAccessPortalRuleV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PolicyAccessPortalRuleV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy portal rules.")
		return
	}

	var plan PolicyAccessPortalRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypePrivilegedPortal, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	request, diags := expandPolicyAccessPortalRuleV2(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := policysetcontrollerv2.CreateRule(ctx, service, request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create policy portal rule: %v", err))
		return
	}

	state, readDiags := r.readPolicyAccessPortalRuleV2(ctx, service, policySetID, created.ID, plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PolicyAccessPortalRuleV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy portal rules.")
		return
	}

	var state PolicyAccessPortalRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypePrivilegedPortal, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	newState, diags := r.readPolicyAccessPortalRuleV2(ctx, service, policySetID, state.ID.ValueString(), state.MicrotenantID)
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

func (r *PolicyAccessPortalRuleV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy portal rules.")
		return
	}

	var plan PolicyAccessPortalRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypePrivilegedPortal, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	request, diags := expandPolicyAccessPortalRuleV2(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := policysetcontrollerv2.UpdateRule(ctx, service, policySetID, plan.ID.ValueString(), request); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update policy portal rule: %v", err))
		return
	}

	newState, readDiags := r.readPolicyAccessPortalRuleV2(ctx, service, policySetID, plan.ID.ValueString(), plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *PolicyAccessPortalRuleV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy portal rules.")
		return
	}

	var state PolicyAccessPortalRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypePrivilegedPortal, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	if _, err := policysetcontrollerv2.Delete(ctx, service, policySetID, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete policy portal rule: %v", err))
	}
}

func (r *PolicyAccessPortalRuleV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before importing policy portal rules.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the policy portal rule ID or name.")
		return
	}

	if _, err := fmt.Sscan(id, new(int64)); err != nil {
		rule, _, lookupErr := policysetcontrollerv2.GetByNameAndTypes(ctx, r.client.Service, []string{helpers.PolicyTypePrivilegedPortal}, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate policy portal rule %q: %v", id, lookupErr))
			return
		}
		id = rule.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *PolicyAccessPortalRuleV2Resource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && microtenantID.ValueString() != "" {
		service = service.WithMicroTenant(microtenantID.ValueString())
	}
	return service
}

func (r *PolicyAccessPortalRuleV2Resource) readPolicyAccessPortalRuleV2(ctx context.Context, service *zscaler.Service, policySetID, ruleID string, microTenantID types.String) (PolicyAccessPortalRuleV2Model, diag.Diagnostics) {
	ruleResource, _, err := policysetcontrollerv2.GetPolicyRule(ctx, service, policySetID, ruleID)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return PolicyAccessPortalRuleV2Model{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Policy portal rule %s not found", ruleID))}
		}
		return PolicyAccessPortalRuleV2Model{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read policy portal rule: %v", err))}
	}

	rule := helpers.ConvertV1ResponseToV2Request(*ruleResource)

	conditions, condDiags := flattenPolicyRuleConditionsV2(ctx, rule.Conditions)
	if condDiags.HasError() {
		return PolicyAccessPortalRuleV2Model{}, condDiags
	}

	capabilities, capDiags := flattenPrivilegedPortalCapabilitiesToList(ctx, rule.PrivilegedPortalCapabilities)
	if capDiags.HasError() {
		return PolicyAccessPortalRuleV2Model{}, capDiags
	}

	model := PolicyAccessPortalRuleV2Model{
		ID:                           types.StringValue(rule.ID),
		Name:                         types.StringValue(rule.Name),
		Description:                  types.StringValue(rule.Description),
		Action:                       types.StringValue(rule.Action),
		PolicySetID:                  types.StringValue(policySetID),
		Conditions:                   conditions,
		PrivilegedPortalCapabilities: capabilities,
		MicrotenantID:                microTenantID,
	}

	return model, condDiags
}

func expandPolicyAccessPortalRuleV2(ctx context.Context, model *PolicyAccessPortalRuleV2Model, policySetID string) (*policysetcontrollerv2.PolicyRule, diag.Diagnostics) {
	conditions, diags := expandPolicyRuleConditionsV2(ctx, model.Conditions)
	if diags.HasError() {
		return nil, diags
	}

	capabilities := expandPrivilegedPortalCapabilitiesFromList(ctx, model.PrivilegedPortalCapabilities)

	rule := &policysetcontrollerv2.PolicyRule{
		ID:                           helpers.StringValue(model.ID),
		Name:                         helpers.StringValue(model.Name),
		Description:                  helpers.StringValue(model.Description),
		Action:                       helpers.StringValue(model.Action),
		PolicySetID:                  policySetID,
		Conditions:                   conditions,
		PrivilegedPortalCapabilities: capabilities,
		MicroTenantID:                helpers.StringValue(model.MicrotenantID),
	}

	return rule, diags
}

func flattenPrivilegedPortalCapabilitiesToList(ctx context.Context, capabilities policysetcontrollerv2.PrivilegedPortalCapabilities) (types.List, diag.Diagnostics) {
	if len(capabilities.Capabilities) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: portalCapabilitiesAttrTypes}), diag.Diagnostics{}
	}

	flags := map[string]bool{
		"delete_file":             false,
		"access_uninspected_file": false,
		"request_approvals":       false,
		"review_approvals":        false,
	}

	for _, cap := range capabilities.Capabilities {
		switch cap {
		case "DELETE_FILE":
			flags["delete_file"] = true
		case "ACCESS_UNINSPECTED_FILE":
			flags["access_uninspected_file"] = true
		case "REQUEST_APPROVALS":
			flags["request_approvals"] = true
		case "REVIEW_APPROVALS":
			flags["review_approvals"] = true
		}
	}

	obj, diags := types.ObjectValue(portalCapabilitiesAttrTypes, map[string]attr.Value{
		"delete_file":             types.BoolValue(flags["delete_file"]),
		"access_uninspected_file": types.BoolValue(flags["access_uninspected_file"]),
		"request_approvals":       types.BoolValue(flags["request_approvals"]),
		"review_approvals":        types.BoolValue(flags["review_approvals"]),
	})
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: portalCapabilitiesAttrTypes}), diags
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: portalCapabilitiesAttrTypes}, []attr.Value{obj})
	diags.Append(listDiags...)
	return list, diags
}

func expandPrivilegedPortalCapabilitiesFromList(ctx context.Context, list types.List) policysetcontrollerv2.PrivilegedPortalCapabilities {
	capabilities := policysetcontrollerv2.PrivilegedPortalCapabilities{}
	if list.IsNull() || list.IsUnknown() {
		return capabilities
	}

	var items []struct {
		DeleteFile            types.Bool `tfsdk:"delete_file"`
		AccessUninspectedFile types.Bool `tfsdk:"access_uninspected_file"`
		RequestApprovals      types.Bool `tfsdk:"request_approvals"`
		ReviewApprovals       types.Bool `tfsdk:"review_approvals"`
	}
	if diags := list.ElementsAs(ctx, &items, false); diags.HasError() {
		return capabilities
	}
	if len(items) == 0 {
		return capabilities
	}

	item := items[0]
	if helpers.BoolValue(item.DeleteFile, false) {
		capabilities.Capabilities = append(capabilities.Capabilities, "DELETE_FILE")
	}
	if helpers.BoolValue(item.AccessUninspectedFile, false) {
		capabilities.Capabilities = append(capabilities.Capabilities, "ACCESS_UNINSPECTED_FILE")
	}
	if helpers.BoolValue(item.RequestApprovals, false) {
		capabilities.Capabilities = append(capabilities.Capabilities, "REQUEST_APPROVALS")
	}
	if helpers.BoolValue(item.ReviewApprovals, false) {
		capabilities.Capabilities = append(capabilities.Capabilities, "REVIEW_APPROVALS")
	}

	return capabilities
}
