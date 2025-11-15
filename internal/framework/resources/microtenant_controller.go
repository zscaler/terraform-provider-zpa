package resources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/microtenants"
)

var (
	_ resource.Resource                = &MicrotenantControllerResource{}
	_ resource.ResourceWithConfigure   = &MicrotenantControllerResource{}
	_ resource.ResourceWithImportState = &MicrotenantControllerResource{}
)

func NewMicrotenantControllerResource() resource.Resource {
	return &MicrotenantControllerResource{}
}

type MicrotenantControllerResource struct {
	client *client.Client
}

type MicrotenantModel struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Description                types.String `tfsdk:"description"`
	Enabled                    types.Bool   `tfsdk:"enabled"`
	CriteriaAttribute          types.String `tfsdk:"criteria_attribute"`
	CriteriaAttributeValues    types.Set    `tfsdk:"criteria_attribute_values"`
	PrivilegedApprovalsEnabled types.Bool   `tfsdk:"privileged_approvals_enabled"`
	User                       types.Set    `tfsdk:"user"`
}

type MicrotenantUserModel struct {
	DisplayName   types.String `tfsdk:"display_name"`
	Username      types.String `tfsdk:"username"`
	Password      types.String `tfsdk:"password"`
	MicrotenantID types.String `tfsdk:"microtenant_id"`
}

func (r *MicrotenantControllerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_microtenant_controller"
}

func (r *MicrotenantControllerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages ZPA microtenant configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Description:   "The unique identifier of the microtenant.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the microtenant.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the microtenant.",
			},
			"enabled": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Description:   "Whether the microtenant is enabled.",
			},
			"criteria_attribute": schema.StringAttribute{
				Optional:    true,
				Description: "The criteria attribute for the microtenant. Supported value is `AuthDomain`.",
			},
			"criteria_attribute_values": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "Criteria attribute values such as authentication domains.",
			},
			"privileged_approvals_enabled": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Description:   "Indicates if privileged approvals are enabled for the microtenant.",
			},
		},
		Blocks: map[string]schema.Block{
			"user": schema.SetNestedBlock{
				Description: "Microtenant user information.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"display_name": schema.StringAttribute{
							Computed: true,
						},
						"username": schema.StringAttribute{
							Computed: true,
						},
						"password": schema.StringAttribute{
							Computed: true,
						},
						"microtenant_id": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (r *MicrotenantControllerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *MicrotenantControllerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan MicrotenantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diags := expandMicrotenant(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := microtenants.Create(ctx, r.client.Service, payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create microtenant: %v", err))
		return
	}

	state, readDiags := r.readIntoState(ctx, created.ID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MicrotenantControllerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state MicrotenantModel
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

func (r *MicrotenantControllerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan MicrotenantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	payload, diags := expandMicrotenant(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := microtenants.Update(ctx, r.client.Service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update microtenant: %v", err))
		return
	}

	state, readDiags := r.readIntoState(ctx, plan.ID.ValueString())
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MicrotenantControllerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state MicrotenantModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := microtenants.Delete(ctx, r.client.Service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete microtenant: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *MicrotenantControllerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before importing resources.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the microtenant ID or name.")
		return
	}

	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		resource, _, lookupErr := microtenants.GetByName(ctx, r.client.Service, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate microtenant %q: %v", id, lookupErr))
			return
		}
		id = resource.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *MicrotenantControllerResource) readIntoState(ctx context.Context, id string) (MicrotenantModel, diag.Diagnostics) {
	resource, _, err := microtenants.Get(ctx, r.client.Service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return MicrotenantModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Microtenant %s not found", id))}
		}
		return MicrotenantModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read microtenant: %v", err))}
	}

	var diags diag.Diagnostics

	// Handle criteria_attribute_values: return empty set if empty (matching SDKv2 behavior)
	// SDKv2 sets it directly from API response, which results in empty set for empty slice
	var criteriaValues types.Set
	if len(resource.CriteriaAttributeValues) == 0 {
		criteriaValues = types.SetValueMust(types.StringType, []attr.Value{})
	} else {
		var valuesDiags diag.Diagnostics
		criteriaValues, valuesDiags = types.SetValueFrom(ctx, types.StringType, resource.CriteriaAttributeValues)
		diags.Append(valuesDiags...)
	}

	userSet, userDiags := flattenMicrotenantUser(ctx, resource.UserResource)
	diags.Append(userDiags...)

	model := MicrotenantModel{
		ID:                         helpers.StringValueOrNull(resource.ID),
		Name:                       helpers.StringValueOrNull(resource.Name),
		Description:                helpers.StringValueOrNull(resource.Description),
		Enabled:                    types.BoolValue(resource.Enabled),
		CriteriaAttribute:          helpers.StringValueOrNull(resource.CriteriaAttribute),
		CriteriaAttributeValues:    criteriaValues,
		PrivilegedApprovalsEnabled: types.BoolValue(resource.PrivilegedApprovalsEnabled),
		User:                       userSet,
	}

	return model, diags
}

func expandMicrotenant(ctx context.Context, model *MicrotenantModel) (microtenants.MicroTenant, diag.Diagnostics) {
	var diags diag.Diagnostics

	values, valuesDiags := helpers.SetValueToStringSlice(ctx, model.CriteriaAttributeValues)
	diags.Append(valuesDiags...)

	payload := microtenants.MicroTenant{
		ID:                         helpers.StringValue(model.ID),
		Name:                       helpers.StringValue(model.Name),
		Description:                helpers.StringValue(model.Description),
		Enabled:                    helpers.BoolValue(model.Enabled, false),
		CriteriaAttribute:          helpers.StringValue(model.CriteriaAttribute),
		CriteriaAttributeValues:    values,
		PrivilegedApprovalsEnabled: helpers.BoolValue(model.PrivilegedApprovalsEnabled, false),
	}

	return payload, diags
}

func flattenMicrotenantUser(ctx context.Context, user *microtenants.UserResource) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	if user == nil {
		return types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"display_name":   types.StringType,
				"username":       types.StringType,
				"password":       types.StringType,
				"microtenant_id": types.StringType,
			},
		}), diags
	}

	attrTypes := map[string]attr.Type{
		"display_name":   types.StringType,
		"username":       types.StringType,
		"password":       types.StringType,
		"microtenant_id": types.StringType,
	}

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"display_name":   helpers.StringValueOrNull(user.DisplayName),
		"username":       helpers.StringValueOrNull(user.Username),
		"password":       helpers.StringValueOrNull(user.Password),
		"microtenant_id": helpers.StringValueOrNull(user.MicrotenantID),
	})
	diags.Append(objDiags...)
	if diags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	set, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	diags.Append(setDiags...)
	return set, diags
}
