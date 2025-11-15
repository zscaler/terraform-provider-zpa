package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praconsole"
)

var (
	_ resource.Resource                = &PRAConsoleControllerResource{}
	_ resource.ResourceWithConfigure   = &PRAConsoleControllerResource{}
	_ resource.ResourceWithImportState = &PRAConsoleControllerResource{}
)

var praConsolePolicyLock sync.Mutex

func NewPRAConsoleControllerResource() resource.Resource {
	return &PRAConsoleControllerResource{}
}

type PRAConsoleControllerResource struct {
	client *client.Client
}

type PRAConsoleControllerModel struct {
	ID             types.String                 `tfsdk:"id"`
	Name           types.String                 `tfsdk:"name"`
	Description    types.String                 `tfsdk:"description"`
	Enabled        types.Bool                   `tfsdk:"enabled"`
	IconText       types.String                 `tfsdk:"icon_text"`
	PRAPortals     []PRAConsolePortalModel      `tfsdk:"pra_portals"`
	PRAApplication []PRAConsoleApplicationModel `tfsdk:"pra_application"`
	MicrotenantID  types.String                 `tfsdk:"microtenant_id"`
}

type PRAConsolePortalModel struct {
	IDs types.Set `tfsdk:"id"`
}

type PRAConsoleApplicationModel struct {
	ID types.String `tfsdk:"id"`
}

func (r *PRAConsoleControllerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pra_console_controller"
}

func (r *PRAConsoleControllerResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZPA PRA console controller.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the PRA console controller.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the PRA console controller.",
			},
			"enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Description: "Indicates whether the console controller is enabled.",
			},
			"icon_text": schema.StringAttribute{
				Optional:    true,
				Description: "Base64 encoded icon text.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Description: "Micro-tenant ID for scoping.",
			},
		},
		Blocks: map[string]schema.Block{
			// pra_portals: TypeSet in SDKv2
			// Using SetNestedBlock for block syntax support
			"pra_portals": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
			// pra_application: TypeList in SDKv2
			// Using ListNestedBlock for block syntax support
			"pra_application": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}

func (r *PRAConsoleControllerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *PRAConsoleControllerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing PRA console controllers.")
		return
	}

	var plan PRAConsoleControllerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload, diags := expandPRAConsole(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := praconsole.Create(ctx, service, &payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create PRA console controller: %v", err))
		return
	}

	tflog.Info(ctx, "Created PRA console controller", map[string]any{"id": created.ID})

	// Pass plan to preserve null values for optional fields (matching SDKv2 behavior)
	state, readDiags := r.readConsole(ctx, service, created.ID, &plan)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PRAConsoleControllerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing PRA console controllers.")
		return
	}

	var state PRAConsoleControllerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	// Pass existing state to preserve null values for optional fields (matching SDKv2 behavior)
	newState, diags := r.readConsole(ctx, service, state.ID.ValueString(), &state)
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

func (r *PRAConsoleControllerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing PRA console controllers.")
		return
	}

	var plan PRAConsoleControllerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing ID", "id must be known during update.")
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload, diags := expandPRAConsole(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := praconsole.Update(ctx, service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update PRA console controller: %v", err))
		return
	}

	// Pass plan to preserve null values for optional fields (matching SDKv2 behavior)
	state, readDiags := r.readConsole(ctx, service, plan.ID.ValueString(), &plan)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PRAConsoleControllerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing PRA console controllers.")
		return
	}

	var state PRAConsoleControllerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	if err := detachPRAConsoleFromPolicies(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Detach Error", fmt.Sprintf("Failed to detach PRA console controller from policies: %v", err))
		return
	}

	if _, err := praconsole.Delete(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete PRA console controller: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *PRAConsoleControllerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *PRAConsoleControllerResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && microtenantID.ValueString() != "" {
		service = service.WithMicroTenant(strings.TrimSpace(microtenantID.ValueString()))
	}
	return service
}

func (r *PRAConsoleControllerResource) readConsole(ctx context.Context, service *zscaler.Service, id string, existingState *PRAConsoleControllerModel) (PRAConsoleControllerModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	console, _, err := praconsole.Get(ctx, service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return PRAConsoleControllerModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("PRA console controller %s not found", id))}
		}
		return PRAConsoleControllerModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read PRA console controller: %v", err))}
	}

	model, flattenDiags := flattenPRAConsole(ctx, console, existingState)
	diags.Append(flattenDiags...)
	return model, diags
}

func expandPRAConsole(ctx context.Context, plan PRAConsoleControllerModel) (praconsole.PRAConsole, diag.Diagnostics) {
	var diags diag.Diagnostics

	portals, portalDiags := expandPRAPortals(ctx, plan.PRAPortals)
	diags.Append(portalDiags...)

	application, appDiags := expandPRAConsoleApplication(plan.PRAApplication)
	diags.Append(appDiags...)

	result := praconsole.PRAConsole{
		ID:            plan.ID.ValueString(),
		Name:          plan.Name.ValueString(),
		Description:   plan.Description.ValueString(),
		Enabled:       plan.Enabled.ValueBool(),
		IconText:      plan.IconText.ValueString(),
		PRAPortals:    portals,
		MicroTenantID: plan.MicrotenantID.ValueString(),
	}

	if application != nil {
		result.PRAApplication = *application
	}

	return result, diags
}

func expandPRAPortals(ctx context.Context, models []PRAConsolePortalModel) ([]praconsole.PRAPortals, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(models) == 0 {
		diags.AddError("Validation Error", "pra_portals must be provided")
		return nil, diags
	}

	var result []praconsole.PRAPortals
	for _, model := range models {
		if model.IDs.IsNull() || model.IDs.IsUnknown() {
			continue
		}
		var ids []string
		diags.Append(model.IDs.ElementsAs(ctx, &ids, false)...)
		for _, id := range ids {
			id = strings.TrimSpace(id)
			if id == "" {
				continue
			}
			result = append(result, praconsole.PRAPortals{ID: id})
		}
	}

	if len(result) == 0 {
		diags.AddError("Validation Error", "pra_portals.id must contain at least one value")
	}

	return result, diags
}

func expandPRAConsoleApplication(models []PRAConsoleApplicationModel) (*praconsole.PRAApplication, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(models) == 0 {
		diags.AddError("Validation Error", "pra_application must be provided")
		return nil, diags
	}

	id := strings.TrimSpace(models[0].ID.ValueString())
	if id == "" {
		diags.AddError("Validation Error", "pra_application.id must be provided")
		return nil, diags
	}

	return &praconsole.PRAApplication{ID: id}, diags
}

func flattenPRAConsole(ctx context.Context, console *praconsole.PRAConsole, existingState *PRAConsoleControllerModel) (PRAConsoleControllerModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	portals, portalDiags := flattenPRAPortals(ctx, console.PRAPortals)
	diags.Append(portalDiags...)

	applications, appDiags := flattenPRAConsoleApplication(ctx, console.PRAApplication)
	diags.Append(appDiags...)

	// Preserve null values for optional fields if they were null in the plan/state
	// SDKv2 sets these directly from API, but we need to preserve null to avoid inconsistent results
	// If API returns empty string and plan/state had null, preserve null
	iconText := helpers.StringValueOrNull(console.IconText)
	if existingState != nil && (existingState.IconText.IsNull() || existingState.IconText.IsUnknown()) && console.IconText == "" {
		iconText = existingState.IconText
	}

	microtenantID := helpers.StringValueOrNull(console.MicroTenantID)
	if existingState != nil && (existingState.MicrotenantID.IsNull() || existingState.MicrotenantID.IsUnknown()) && console.MicroTenantID == "" {
		microtenantID = existingState.MicrotenantID
	}

	return PRAConsoleControllerModel{
		ID:             types.StringValue(console.ID),
		Name:           helpers.StringValueOrNull(console.Name),
		Description:    helpers.StringValueOrNull(console.Description),
		Enabled:        types.BoolValue(console.Enabled),
		IconText:       iconText,
		PRAPortals:     portals,
		PRAApplication: applications,
		MicrotenantID:  microtenantID,
	}, diags
}

func flattenPRAPortals(ctx context.Context, portals []praconsole.PRAPortals) ([]PRAConsolePortalModel, diag.Diagnostics) {
	if len(portals) == 0 {
		return nil, nil
	}

	ids := make([]string, 0, len(portals))
	for _, portal := range portals {
		if strings.TrimSpace(portal.ID) != "" {
			ids = append(ids, portal.ID)
		}
	}

	setValue, diags := types.SetValueFrom(ctx, types.StringType, ids)
	if diags.HasError() {
		return nil, diags
	}

	return []PRAConsolePortalModel{{IDs: setValue}}, nil
}

func flattenPRAConsoleApplication(ctx context.Context, application praconsole.PRAApplication) ([]PRAConsoleApplicationModel, diag.Diagnostics) {
	if strings.TrimSpace(application.ID) == "" {
		return nil, nil
	}
	return []PRAConsoleApplicationModel{{ID: types.StringValue(application.ID)}}, nil
}

func detachPRAConsoleFromPolicies(ctx context.Context, service *zscaler.Service, consoleID string) error {
	praConsolePolicyLock.Lock()
	defer praConsolePolicyLock.Unlock()

	types := []string{"CREDENTIAL_POLICY"}
	var allRules []policysetcontrollerv2.PolicyRuleResource

	for _, policyType := range types {
		policySet, _, err := policysetcontrollerv2.GetByPolicyType(ctx, service, policyType)
		if err != nil {
			return fmt.Errorf("failed to get policy set for type %s: %w", policyType, err)
		}

		rules, _, err := policysetcontrollerv2.GetAllByType(ctx, service, policyType)
		if err != nil {
			return fmt.Errorf("failed to get rules for policy type %s: %w", policyType, err)
		}

		for _, rule := range rules {
			rule.PolicySetID = policySet.ID
			allRules = append(allRules, rule)
		}
	}

	for _, rule := range allRules {
		modified := false
		newConditions := make([]policysetcontrollerv2.PolicyRuleResourceConditions, 0, len(rule.Conditions))

		for _, condition := range rule.Conditions {
			newOperands := make([]policysetcontrollerv2.PolicyRuleResourceOperands, 0, len(condition.Operands))
			for _, operand := range condition.Operands {
				if strings.EqualFold(operand.ObjectType, "CONSOLE") && strings.EqualFold(operand.LHS, "id") {
					filtered := make([]string, 0, len(operand.Values))
					for _, value := range operand.Values {
						if value == consoleID {
							modified = true
							continue
						}
						filtered = append(filtered, value)
					}

					if operand.RHS == consoleID {
						modified = true
						continue
					}

					if len(filtered) > 0 {
						operand.Values = filtered
						newOperands = append(newOperands, operand)
					}
				} else {
					newOperands = append(newOperands, operand)
				}
			}

			if len(newOperands) > 0 {
				condition.Operands = newOperands
				newConditions = append(newConditions, condition)
			}
		}

		if modified {
			converted := helpers.ConvertV1ResponseToV2Request(rule)

			payload, _ := json.MarshalIndent(rule.Conditions, "", "  ")
			tflog.Debug(ctx, "Detaching PRA console from policy rule", map[string]any{
				"rule_id":   rule.ID,
				"payload":   string(payload),
				"consoleID": consoleID,
			})

			const maxRetries = 3
			for attempt := 0; attempt < maxRetries; attempt++ {
				if _, err := policysetcontrollerv2.UpdateRule(ctx, service, rule.PolicySetID, rule.ID, &converted); err != nil {
					if attempt == maxRetries-1 {
						return fmt.Errorf("failed to update rule %s after %d attempts: %w", rule.ID, maxRetries, err)
					}
					time.Sleep(1 * time.Second)
					continue
				}
				break
			}
		}
	}

	return nil
}
