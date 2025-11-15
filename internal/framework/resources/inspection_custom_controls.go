package resources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_custom_controls"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_profile"
)

var (
	_ resource.Resource                = &InspectionCustomControlsResource{}
	_ resource.ResourceWithConfigure   = &InspectionCustomControlsResource{}
	_ resource.ResourceWithImportState = &InspectionCustomControlsResource{}
)

func NewInspectionCustomControlsResource() resource.Resource {
	return &InspectionCustomControlsResource{}
}

type InspectionCustomControlsResource struct {
	client *client.Client
}

type InspectionCustomControlModel struct {
	ID                           types.String                       `tfsdk:"id"`
	Name                         types.String                       `tfsdk:"name"`
	Description                  types.String                       `tfsdk:"description"`
	DefaultAction                types.String                       `tfsdk:"default_action"`
	DefaultActionValue           types.String                       `tfsdk:"default_action_value"`
	ParanoiaLevel                types.String                       `tfsdk:"paranoia_level"`
	ProtocolType                 types.String                       `tfsdk:"protocol_type"`
	Severity                     types.String                       `tfsdk:"severity"`
	Type                         types.String                       `tfsdk:"type"`
	Version                      types.String                       `tfsdk:"version"`
	ControlType                  types.String                       `tfsdk:"control_type"`
	AssociatedInspectionProfiles []AssociatedInspectionProfileModel `tfsdk:"associated_inspection_profiles"`
	Rules                        []InspectionCustomControlRuleModel `tfsdk:"rules"`
}

type InspectionCustomControlRuleModel struct {
	Names      types.Set                               `tfsdk:"names"`
	Type       types.String                            `tfsdk:"type"`
	Conditions []InspectionCustomControlConditionModel `tfsdk:"conditions"`
}

type InspectionCustomControlConditionModel struct {
	LHS types.String `tfsdk:"lhs"`
	OP  types.String `tfsdk:"op"`
	RHS types.String `tfsdk:"rhs"`
}

type AssociatedInspectionProfileModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (r *InspectionCustomControlsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inspection_custom_controls"
}

func (r *InspectionCustomControlsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages ZPA inspection custom controls.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the custom control.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Description of the custom control.",
			},
			"default_action": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{stringvalidator.OneOf(
					"PASS", "BLOCK", "REDIRECT",
				)},
			},
			"default_action_value": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Redirect URL when default action is REDIRECT.",
			},
			"paranoia_level": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"protocol_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{stringvalidator.OneOf(
					"HTTP", "HTTPS", "FTP", "RDP", "SSH", "WEBSOCKET", "VNC", "NONE", "AUTO", "DYNAMIC",
				)},
			},
			"severity": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{stringvalidator.OneOf(
					"CRITICAL", "ERROR", "WARNING", "INFO",
				)},
			},
			"type": schema.StringAttribute{
				Required:   true,
				Validators: []validator.String{stringvalidator.OneOf("REQUEST", "RESPONSE")},
			},
			"version": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"control_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{stringvalidator.OneOf(
					"WEBSOCKET_PREDEFINED", "WEBSOCKET_CUSTOM", "THREATLABZ", "CUSTOM", "PREDEFINED", "API_PREDEFINED",
				)},
			},
		},
		Blocks: map[string]schema.Block{
			"associated_inspection_profiles": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"rules": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"names": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,
							Description: "Required when rules.type is REQUEST_HEADERS, REQUEST_COOKIES, or RESPONSE_HEADERS.",
						},
						"type": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Validators: []validator.String{stringvalidator.OneOf(
								"REQUEST_HEADERS", "REQUEST_URI", "QUERY_STRING", "REQUEST_COOKIES", "REQUEST_METHOD",
								"REQUEST_BODY", "RESPONSE_HEADERS", "RESPONSE_BODY", "WS_MAX_PAYLOAD_SIZE", "WS_MAX_FRAGMENT_PER_MESSAGE",
							)},
						},
					},
					Blocks: map[string]schema.Block{
						"conditions": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"lhs": schema.StringAttribute{
										Optional:   true,
										Computed:   true,
										Validators: []validator.String{stringvalidator.OneOf("SIZE", "VALUE")},
									},
									"op": schema.StringAttribute{
										Optional: true,
										Computed: true,
										Validators: []validator.String{stringvalidator.OneOf(
											"RX", "EQ", "LE", "GE", "CONTAINS", "STARTS_WITH", "ENDS_WITH",
										)},
									},
									"rhs": schema.StringAttribute{
										Optional: true,
										Computed: true,
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

func (r *InspectionCustomControlsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cl, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = cl
}

func (r *InspectionCustomControlsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan InspectionCustomControlModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diags := r.expand(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := inspection_custom_controls.Create(ctx, r.client.Service, payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create inspection custom control: %v", err))
		return
	}

	if err := r.updateInspectionProfiles(ctx, created.ID, payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update associated inspection profiles: %v", err))
		return
	}

	state, readDiags := r.readIntoState(ctx, created.ID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *InspectionCustomControlsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state InspectionCustomControlModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readIntoState(ctx, state.ID.ValueString())
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

func (r *InspectionCustomControlsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan InspectionCustomControlModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	payload, diags := r.expand(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := inspection_custom_controls.Update(ctx, r.client.Service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update inspection custom control: %v", err))
		return
	}

	if err := r.updateInspectionProfiles(ctx, plan.ID.ValueString(), payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update associated inspection profiles: %v", err))
		return
	}

	state, readDiags := r.readIntoState(ctx, plan.ID.ValueString())
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *InspectionCustomControlsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state InspectionCustomControlModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.detachFromProfiles(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to detach inspection custom control from inspection profiles: %v", err))
		return
	}

	if _, err := inspection_custom_controls.Delete(ctx, r.client.Service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete inspection custom control: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *InspectionCustomControlsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before importing resources.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires either the inspection custom control ID or name.")
		return
	}

	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		control, _, lookupErr := inspection_custom_controls.GetByName(ctx, r.client.Service, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate inspection custom control %q: %v", id, lookupErr))
			return
		}
		id = control.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *InspectionCustomControlsResource) readIntoState(ctx context.Context, id string) (InspectionCustomControlModel, diag.Diagnostics) {
	resource, _, err := inspection_custom_controls.Get(ctx, r.client.Service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return InspectionCustomControlModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Inspection custom control %s not found", id))}
		}
		return InspectionCustomControlModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read inspection custom control: %v", err))}
	}

	rules, ruleDiags := r.flattenRules(ctx, resource.Rules)
	if ruleDiags.HasError() {
		return InspectionCustomControlModel{}, ruleDiags
	}

	state := InspectionCustomControlModel{
		ID:                           helpers.StringValueOrNull(resource.ID),
		Name:                         helpers.StringValueOrNull(resource.Name),
		Description:                  helpers.StringValueOrNull(resource.Description),
		DefaultAction:                helpers.StringValueOrNull(resource.DefaultAction),
		DefaultActionValue:           helpers.StringValueOrNull(resource.DefaultActionValue),
		ParanoiaLevel:                helpers.StringValueOrNull(resource.ParanoiaLevel),
		ProtocolType:                 helpers.StringValueOrNull(resource.ProtocolType),
		Severity:                     helpers.StringValueOrNull(resource.Severity),
		Type:                         helpers.StringValueOrNull(resource.Type),
		Version:                      helpers.StringValueOrNull(resource.Version),
		ControlType:                  helpers.StringValueOrNull(resource.ControlType),
		AssociatedInspectionProfiles: flattenAssociatedProfiles(resource.AssociatedInspectionProfileNames),
		Rules:                        rules,
	}

	return state, nil
}

func (r *InspectionCustomControlsResource) expand(ctx context.Context, plan InspectionCustomControlModel) (inspection_custom_controls.InspectionCustomControl, diag.Diagnostics) {
	var diags diag.Diagnostics

	rules, validations, ruleDiags := r.expandRules(ctx, plan)
	diags.Append(ruleDiags...)

	associatedProfiles, assocDiags := expandAssociatedProfiles(plan.AssociatedInspectionProfiles)
	diags.Append(assocDiags...)

	payload := inspection_custom_controls.InspectionCustomControl{
		ID:                               helpers.StringValue(plan.ID),
		Name:                             helpers.StringValue(plan.Name),
		Description:                      helpers.StringValue(plan.Description),
		DefaultAction:                    helpers.StringValue(plan.DefaultAction),
		DefaultActionValue:               helpers.StringValue(plan.DefaultActionValue),
		ParanoiaLevel:                    helpers.StringValue(plan.ParanoiaLevel),
		ProtocolType:                     helpers.StringValue(plan.ProtocolType),
		Severity:                         helpers.StringValue(plan.Severity),
		Type:                             helpers.StringValue(plan.Type),
		Version:                          helpers.StringValue(plan.Version),
		ControlType:                      helpers.StringValue(plan.ControlType),
		Rules:                            rules,
		AssociatedInspectionProfileNames: associatedProfiles,
	}

	diags.Append(validateInspectionCustomControl(payload, validations)...) // augment diagnostics with validation errors

	return payload, diags
}

func (r *InspectionCustomControlsResource) expandRules(ctx context.Context, plan InspectionCustomControlModel) ([]inspection_custom_controls.Rules, [][]inspection_custom_controls.Conditions, diag.Diagnostics) {
	var diags diag.Diagnostics
	rules := make([]inspection_custom_controls.Rules, 0, len(plan.Rules))
	validations := make([][]inspection_custom_controls.Conditions, 0, len(plan.Rules))

	for _, rule := range plan.Rules {
		names, namesDiags := helpers.SetValueToStringSlice(ctx, rule.Names)
		diags.Append(namesDiags...)

		conditions := make([]inspection_custom_controls.Conditions, 0, len(rule.Conditions))
		for _, condition := range rule.Conditions {
			conditions = append(conditions, inspection_custom_controls.Conditions{
				LHS: helpers.StringValue(condition.LHS),
				OP:  helpers.StringValue(condition.OP),
				RHS: helpers.StringValue(condition.RHS),
			})
		}

		rules = append(rules, inspection_custom_controls.Rules{
			Names:      names,
			Type:       helpers.StringValue(rule.Type),
			Conditions: conditions,
		})
		validations = append(validations, conditions)
	}

	return rules, validations, diags
}

func expandAssociatedProfiles(profiles []AssociatedInspectionProfileModel) ([]common.AssociatedProfileNames, diag.Diagnostics) {
	if len(profiles) == 0 {
		return nil, diag.Diagnostics{}
	}

	result := make([]common.AssociatedProfileNames, 0, len(profiles))
	for _, p := range profiles {
		result = append(result, common.AssociatedProfileNames{
			ID:   helpers.StringValue(p.ID),
			Name: helpers.StringValue(p.Name),
		})
	}

	return result, diag.Diagnostics{}
}

func (r *InspectionCustomControlsResource) flattenRules(ctx context.Context, rules []inspection_custom_controls.Rules) ([]InspectionCustomControlRuleModel, diag.Diagnostics) {
	result := make([]InspectionCustomControlRuleModel, 0, len(rules))
	var diags diag.Diagnostics

	for _, rule := range rules {
		names, nameDiags := types.SetValueFrom(ctx, types.StringType, rule.Names)
		diags.Append(nameDiags...)

		conditions := make([]InspectionCustomControlConditionModel, 0, len(rule.Conditions))
		for _, condition := range rule.Conditions {
			conditions = append(conditions, InspectionCustomControlConditionModel{
				LHS: helpers.StringValueOrNull(condition.LHS),
				OP:  helpers.StringValueOrNull(condition.OP),
				RHS: helpers.StringValueOrNull(condition.RHS),
			})
		}

		result = append(result, InspectionCustomControlRuleModel{
			Names:      names,
			Type:       helpers.StringValueOrNull(rule.Type),
			Conditions: conditions,
		})
	}

	return result, diags
}

func flattenAssociatedProfiles(profiles []common.AssociatedProfileNames) []AssociatedInspectionProfileModel {
	if len(profiles) == 0 {
		return nil
	}

	result := make([]AssociatedInspectionProfileModel, 0, len(profiles))
	for _, profile := range profiles {
		result = append(result, AssociatedInspectionProfileModel{
			ID:   helpers.StringValueOrNull(profile.ID),
			Name: helpers.StringValueOrNull(profile.Name),
		})
	}

	return result
}

func (r *InspectionCustomControlsResource) updateInspectionProfiles(ctx context.Context, controlID string, payload inspection_custom_controls.InspectionCustomControl) error {
	if len(payload.AssociatedInspectionProfileNames) == 0 {
		return nil
	}

	control, _, err := inspection_custom_controls.Get(ctx, r.client.Service, controlID)
	if err != nil {
		return err
	}

	for _, profile := range payload.AssociatedInspectionProfileNames {
		profileObj, _, err := inspection_profile.Get(ctx, r.client.Service, profile.ID)
		if err != nil {
			return err
		}

		customControls := make([]inspection_profile.InspectionCustomControl, 0, len(profileObj.CustomControls))
		for _, existing := range profileObj.CustomControls {
			if existing.ID == control.ID {
				continue
			}
			customControls = append(customControls, existing)
		}

		customControls = append(customControls, inspection_profile.InspectionCustomControl{
			ID:                 control.ID,
			DefaultAction:      payload.DefaultAction,
			DefaultActionValue: payload.DefaultActionValue,
		})

		update := &inspection_profile.InspectionProfile{
			CustomControls:     customControls,
			PredefinedControls: profileObj.PredefinedControls,
		}

		if _, err := inspection_profile.Patch(ctx, r.client.Service, profileObj.ID, update); err != nil {
			return err
		}
	}

	return nil
}

func (r *InspectionCustomControlsResource) detachFromProfiles(ctx context.Context, controlID string) error {
	control, _, err := inspection_custom_controls.Get(ctx, r.client.Service, controlID)
	if err != nil {
		return err
	}

	for _, profile := range control.AssociatedInspectionProfileNames {
		profileObj, _, err := inspection_profile.Get(ctx, r.client.Service, profile.ID)
		if err != nil {
			continue
		}

		customControls := make([]inspection_profile.InspectionCustomControl, 0, len(profileObj.CustomControls))
		for _, existing := range profileObj.CustomControls {
			if existing.ID == control.ID {
				continue
			}
			customControls = append(customControls, existing)
		}

		profileObj.CustomControls = customControls
		if _, err := inspection_profile.Update(ctx, r.client.Service, profile.ID, profileObj); err != nil {
			return err
		}
	}

	return nil
}

func validateInspectionCustomControl(control inspection_custom_controls.InspectionCustomControl, rulesValues [][]inspection_custom_controls.Conditions) diag.Diagnostics {
	var diags diag.Diagnostics

	if strings.EqualFold(control.DefaultAction, "REDIRECT") && strings.TrimSpace(control.DefaultActionValue) == "" {
		diags.Append(diag.NewErrorDiagnostic("Validation Error", "when default_action is REDIRECT, default_action_value must be set"))
	}

	for idx, rule := range control.Rules {
		var conditions []inspection_custom_controls.Conditions
		if idx < len(rulesValues) {
			conditions = rulesValues[idx]
		}
		if strings.EqualFold(control.Type, "RESPONSE") {
			if rule.Type != "RESPONSE_HEADERS" && rule.Type != "RESPONSE_BODY" {
				diags.Append(diag.NewErrorDiagnostic("Validation Error", "when type is RESPONSE, rules.type must be RESPONSE_HEADERS or RESPONSE_BODY"))
			}
		} else if strings.EqualFold(control.Type, "REQUEST") {
			if (rule.Type == "REQUEST_HEADERS" || rule.Type == "REQUEST_COOKIES") && len(rule.Names) == 0 {
				diags.Append(diag.NewErrorDiagnostic("Validation Error", "when type is REQUEST and rules.type is REQUEST_HEADERS or REQUEST_COOKIES, rules.names must be set"))
			}
			if (rule.Type == "REQUEST_URI" || rule.Type == "QUERY_STRING" || rule.Type == "REQUEST_BODY" || rule.Type == "REQUEST_METHOD") && len(rule.Names) > 0 {
				diags.Append(diag.NewErrorDiagnostic("Validation Error", "when rules.type is REQUEST_URI, QUERY_STRING, REQUEST_BODY, or REQUEST_METHOD, rules.names cannot be set"))
			}
		}

		for _, condition := range conditions {
			if in(rule.Type, []string{"REQUEST_HEADERS", "REQUEST_COOKIES", "REQUEST_URI", "QUERY_STRING", "REQUEST_BODY"}) {
				if strings.EqualFold(condition.LHS, "SIZE") {
					if !in(condition.OP, []string{"EQ", "LE", "GE"}) || !isNumber(condition.RHS) {
						diags.Append(diag.NewErrorDiagnostic("Validation Error", fmt.Sprintf("when rules.type is %s and conditions.lhs == SIZE, conditions.op must be EQ, LE, or GE and rhs must be a number", rule.Type)))
					}
				}
				if strings.EqualFold(condition.LHS, "VALUE") {
					if !in(condition.OP, []string{"CONTAINS", "STARTS_WITH", "ENDS_WITH", "RX"}) {
						diags.Append(diag.NewErrorDiagnostic("Validation Error", fmt.Sprintf("when rules.type is %s and conditions.lhs == VALUE, conditions.op must be CONTAINS, STARTS_WITH, ENDS_WITH, or RX", rule.Type)))
					}
				}
			}
			if strings.EqualFold(rule.Type, "REQUEST_METHOD") {
				if strings.EqualFold(condition.LHS, "SIZE") {
					if !in(condition.OP, []string{"EQ", "LE", "GE"}) || !isNumber(condition.RHS) {
						diags.Append(diag.NewErrorDiagnostic("Validation Error", "when rules.type is REQUEST_METHOD and conditions.lhs == SIZE, conditions.op must be EQ, LE, or GE and rhs must be a number"))
					}
				}
				if strings.EqualFold(condition.LHS, "VALUE") {
					if !in(condition.OP, []string{"CONTAINS", "STARTS_WITH", "ENDS_WITH", "RX"}) || !in(strings.ToUpper(condition.RHS), []string{"GET", "POST", "PUT", "PATCH", "CONNECT", "HEAD", "OPTIONS", "DELETE", "TRACE"}) {
						diags.Append(diag.NewErrorDiagnostic("Validation Error", "when rules.type is REQUEST_METHOD and conditions.lhs == VALUE, conditions.op must be CONTAINS, STARTS_WITH, ENDS_WITH, or RX and rhs must be a valid HTTP method"))
					}
				}
			}
		}
	}

	return diags
}

func (r *InspectionCustomControlsResource) updateInspectionProfilesFromState(ctx context.Context, id string, plan InspectionCustomControlModel) diag.Diagnostics {
	payload, diags := r.expand(ctx, plan)
	if diags.HasError() {
		return diags
	}
	if err := r.updateInspectionProfiles(ctx, id, payload); err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", err.Error())}
	}
	return diag.Diagnostics{}
}

func in(val string, list []string) bool {
	for _, v := range list {
		if strings.EqualFold(v, val) {
			return true
		}
	}
	return false
}

func isNumber(str string) bool {
	if str == "" {
		return false
	}
	for _, ch := range str {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
