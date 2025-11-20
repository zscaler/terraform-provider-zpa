package resources

import (
	"context"
	"fmt"
	"strings"

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
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

var (
	_ resource.Resource                = &PolicyAccessCapabilitiesRuleV2Resource{}
	_ resource.ResourceWithConfigure   = &PolicyAccessCapabilitiesRuleV2Resource{}
	_ resource.ResourceWithImportState = &PolicyAccessCapabilitiesRuleV2Resource{}
)

func NewPolicyAccessCapabilitiesRuleV2Resource() resource.Resource {
	return &PolicyAccessCapabilitiesRuleV2Resource{}
}

type PolicyAccessCapabilitiesRuleV2Resource struct {
	client *client.Client
}

type PolicyAccessCapabilitiesRuleV2Model struct {
	ID                     types.String                 `tfsdk:"id"`
	Name                   types.String                 `tfsdk:"name"`
	Description            types.String                 `tfsdk:"description"`
	Action                 types.String                 `tfsdk:"action"`
	PolicySetID            types.String                 `tfsdk:"policy_set_id"`
	Conditions             []PolicyAccessConditionModel `tfsdk:"conditions"`
	PrivilegedCapabilities types.List                   `tfsdk:"privileged_capabilities"`
	MicrotenantID          types.String                 `tfsdk:"microtenant_id"`
}

var privilegedCapabilitiesAttrTypes = map[string]attr.Type{
	"clipboard_copy":        types.BoolType,
	"clipboard_paste":       types.BoolType,
	"file_upload":           types.BoolType,
	"file_download":         types.BoolType,
	"inspect_file_download": types.BoolType,
	"inspect_file_upload":   types.BoolType,
	"monitor_session":       types.BoolType,
	"record_session":        types.BoolType,
	"share_session":         types.BoolType,
}

func (r *PolicyAccessCapabilitiesRuleV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_access_capabilities_rule_v2"
}

func (r *PolicyAccessCapabilitiesRuleV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				stringvalidator.OneOf("CHECK_CAPABILITIES"),
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
		"SAML",
		"SCIM",
		"SCIM_GROUP",
	}

	resp.Schema = schema.Schema{
		Description: "Manages ZPA Privileged Capabilities policy rules (v2).",
		Attributes:  attrs,
		Blocks: map[string]schema.Block{
			"conditions": helpers.PolicyConditionsV2Block(objectTypes),
			"privileged_capabilities": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"clipboard_copy":        schema.BoolAttribute{Optional: true},
						"clipboard_paste":       schema.BoolAttribute{Optional: true},
						"file_upload":           schema.BoolAttribute{Optional: true},
						"file_download":         schema.BoolAttribute{Optional: true},
						"inspect_file_download": schema.BoolAttribute{Optional: true},
						"inspect_file_upload":   schema.BoolAttribute{Optional: true},
						"monitor_session":       schema.BoolAttribute{Optional: true},
						"record_session":        schema.BoolAttribute{Optional: true},
						"share_session":         schema.BoolAttribute{Optional: true},
					},
				},
			},
		},
	}
}

func (r *PolicyAccessCapabilitiesRuleV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PolicyAccessCapabilitiesRuleV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy capabilities rules.")
		return
	}

	var plan PolicyAccessCapabilitiesRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeCapabilities, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	request, diags := expandPolicyAccessCapabilitiesRuleV2(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := policysetcontrollerv2.CreateRule(ctx, service, request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create policy capabilities rule: %v", err))
		return
	}

	state, readDiags := r.readPolicyAccessCapabilitiesRuleV2(ctx, service, policySetID, created.ID, plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PolicyAccessCapabilitiesRuleV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy capabilities rules.")
		return
	}

	var state PolicyAccessCapabilitiesRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeCapabilities, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	newState, diags := r.readPolicyAccessCapabilitiesRuleV2(ctx, service, policySetID, state.ID.ValueString(), state.MicrotenantID)
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

func (r *PolicyAccessCapabilitiesRuleV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy capabilities rules.")
		return
	}

	var plan PolicyAccessCapabilitiesRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeCapabilities, microTenantID)
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

	request, diags := expandPolicyAccessCapabilitiesRuleV2(ctx, &plan, policySetID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := policysetcontrollerv2.UpdateRule(ctx, service, policySetID, plan.ID.ValueString(), request); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update policy capabilities rule: %v", err))
		return
	}

	newState, readDiags := r.readPolicyAccessCapabilitiesRuleV2(ctx, service, policySetID, plan.ID.ValueString(), plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *PolicyAccessCapabilitiesRuleV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing policy capabilities rules.")
		return
	}

	var state PolicyAccessCapabilitiesRuleV2Model
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
		policySetID, err = helpers.FetchPolicySetIDByType(ctx, helperClient, helpers.PolicyTypeCapabilities, microTenantID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to determine policy set ID: %v", err))
			return
		}
	}

	if _, err := policysetcontrollerv2.Delete(ctx, service, policySetID, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete policy capabilities rule: %v", err))
	}
}

func (r *PolicyAccessCapabilitiesRuleV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before importing policy capabilities rules.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the policy capabilities rule ID or name.")
		return
	}

	if _, err := fmt.Sscan(id, new(int64)); err != nil {
		rule, _, lookupErr := policysetcontrollerv2.GetByNameAndTypes(ctx, r.client.Service, []string{helpers.PolicyTypeCapabilities}, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate policy capabilities rule %q: %v", id, lookupErr))
			return
		}
		id = rule.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *PolicyAccessCapabilitiesRuleV2Resource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && microtenantID.ValueString() != "" {
		service = service.WithMicroTenant(microtenantID.ValueString())
	}
	return service
}

func (r *PolicyAccessCapabilitiesRuleV2Resource) readPolicyAccessCapabilitiesRuleV2(ctx context.Context, service *zscaler.Service, policySetID, ruleID string, microTenantID types.String) (PolicyAccessCapabilitiesRuleV2Model, diag.Diagnostics) {
	ruleResource, _, err := policysetcontrollerv2.GetPolicyRule(ctx, service, policySetID, ruleID)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return PolicyAccessCapabilitiesRuleV2Model{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Policy capabilities rule %s not found", ruleID))}
		}
		return PolicyAccessCapabilitiesRuleV2Model{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read policy capabilities rule: %v", err))}
	}

	rule := helpers.ConvertV1ResponseToV2Request(*ruleResource)

	conditions, condDiags := flattenPolicyRuleConditionsV2(ctx, rule.Conditions)
	if condDiags.HasError() {
		return PolicyAccessCapabilitiesRuleV2Model{}, condDiags
	}

	capabilities, capDiags := flattenPrivilegedCapabilitiesToList(ctx, rule.PrivilegedCapabilities)
	if capDiags.HasError() {
		return PolicyAccessCapabilitiesRuleV2Model{}, capDiags
	}

	model := PolicyAccessCapabilitiesRuleV2Model{
		ID:                     types.StringValue(rule.ID),
		Name:                   types.StringValue(rule.Name),
		Description:            types.StringValue(rule.Description),
		Action:                 types.StringValue(rule.Action),
		PolicySetID:            types.StringValue(policySetID),
		Conditions:             conditions,
		PrivilegedCapabilities: capabilities,
		MicrotenantID:          microTenantID,
	}

	return model, condDiags
}

func expandPolicyAccessCapabilitiesRuleV2(ctx context.Context, model *PolicyAccessCapabilitiesRuleV2Model, policySetID string) (*policysetcontrollerv2.PolicyRule, diag.Diagnostics) {
	conditions, diags := expandPolicyRuleConditionsV2(ctx, model.Conditions)
	if diags.HasError() {
		return nil, diags
	}

	capabilities := expandPrivilegedCapabilitiesFromList(ctx, model.PrivilegedCapabilities)

	rule := &policysetcontrollerv2.PolicyRule{
		ID:                     helpers.StringValue(model.ID),
		Name:                   helpers.StringValue(model.Name),
		Description:            helpers.StringValue(model.Description),
		Action:                 helpers.StringValue(model.Action),
		PolicySetID:            policySetID,
		Conditions:             conditions,
		PrivilegedCapabilities: capabilities,
		MicroTenantID:          helpers.StringValue(model.MicrotenantID),
	}

	return rule, diags
}

func flattenPrivilegedCapabilitiesToList(ctx context.Context, capabilities policysetcontrollerv2.PrivilegedCapabilities) (types.List, diag.Diagnostics) {
	if len(capabilities.Capabilities) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: privilegedCapabilitiesAttrTypes}), diag.Diagnostics{}
	}

	flags := map[string]bool{}
	for name := range privilegedCapabilitiesAttrTypes {
		flags[name] = false
	}

	for _, cap := range capabilities.Capabilities {
		switch cap {
		case "CLIPBOARD_COPY":
			flags["clipboard_copy"] = true
		case "CLIPBOARD_PASTE":
			flags["clipboard_paste"] = true
		case "FILE_UPLOAD":
			flags["file_upload"] = true
		case "FILE_DOWNLOAD":
			flags["file_download"] = true
		case "INSPECT_FILE_DOWNLOAD":
			flags["inspect_file_download"] = true
		case "INSPECT_FILE_UPLOAD":
			flags["inspect_file_upload"] = true
		case "MONITOR_SESSION":
			flags["monitor_session"] = true
		case "RECORD_SESSION":
			flags["record_session"] = true
		case "SHARE_SESSION":
			flags["share_session"] = true
		}
	}

	obj, diags := types.ObjectValue(privilegedCapabilitiesAttrTypes, map[string]attr.Value{
		"clipboard_copy":        types.BoolValue(flags["clipboard_copy"]),
		"clipboard_paste":       types.BoolValue(flags["clipboard_paste"]),
		"file_upload":           types.BoolValue(flags["file_upload"]),
		"file_download":         types.BoolValue(flags["file_download"]),
		"inspect_file_download": types.BoolValue(flags["inspect_file_download"]),
		"inspect_file_upload":   types.BoolValue(flags["inspect_file_upload"]),
		"monitor_session":       types.BoolValue(flags["monitor_session"]),
		"record_session":        types.BoolValue(flags["record_session"]),
		"share_session":         types.BoolValue(flags["share_session"]),
	})
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: privilegedCapabilitiesAttrTypes}), diags
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: privilegedCapabilitiesAttrTypes}, []attr.Value{obj})
	diags.Append(listDiags...)
	return list, diags
}

func expandPrivilegedCapabilitiesFromList(ctx context.Context, list types.List) policysetcontrollerv2.PrivilegedCapabilities {
	capabilities := policysetcontrollerv2.PrivilegedCapabilities{}
	if list.IsNull() || list.IsUnknown() {
		return capabilities
	}

	var items []struct {
		ClipboardCopy       types.Bool `tfsdk:"clipboard_copy"`
		ClipboardPaste      types.Bool `tfsdk:"clipboard_paste"`
		FileUpload          types.Bool `tfsdk:"file_upload"`
		FileDownload        types.Bool `tfsdk:"file_download"`
		InspectFileDownload types.Bool `tfsdk:"inspect_file_download"`
		InspectFileUpload   types.Bool `tfsdk:"inspect_file_upload"`
		MonitorSession      types.Bool `tfsdk:"monitor_session"`
		RecordSession       types.Bool `tfsdk:"record_session"`
		ShareSession        types.Bool `tfsdk:"share_session"`
	}
	if diags := list.ElementsAs(ctx, &items, false); diags.HasError() {
		return capabilities
	}
	if len(items) == 0 {
		return capabilities
	}

	item := items[0]
	if helpers.BoolValue(item.ClipboardCopy, false) {
		capabilities.Capabilities = append(capabilities.Capabilities, "CLIPBOARD_COPY")
	}
	if helpers.BoolValue(item.ClipboardPaste, false) {
		capabilities.Capabilities = append(capabilities.Capabilities, "CLIPBOARD_PASTE")
	}
	if helpers.BoolValue(item.FileUpload, false) {
		capabilities.Capabilities = append(capabilities.Capabilities, "FILE_UPLOAD")
	}
	if helpers.BoolValue(item.FileDownload, false) {
		capabilities.Capabilities = append(capabilities.Capabilities, "FILE_DOWNLOAD")
	}
	if helpers.BoolValue(item.InspectFileDownload, false) {
		capabilities.Capabilities = append(capabilities.Capabilities, "INSPECT_FILE_DOWNLOAD")
	}
	if helpers.BoolValue(item.InspectFileUpload, false) {
		capabilities.Capabilities = append(capabilities.Capabilities, "INSPECT_FILE_UPLOAD")
	}
	if helpers.BoolValue(item.MonitorSession, false) {
		capabilities.Capabilities = append(capabilities.Capabilities, "MONITOR_SESSION")
	}
	if helpers.BoolValue(item.RecordSession, false) {
		capabilities.Capabilities = append(capabilities.Capabilities, "RECORD_SESSION")
	}
	if helpers.BoolValue(item.ShareSession, false) {
		capabilities.Capabilities = append(capabilities.Capabilities, "SHARE_SESSION")
	}

	return capabilities
}
