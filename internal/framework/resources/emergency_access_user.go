package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/emergencyaccess"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &EmergencyAccessResource{}
	_ resource.ResourceWithConfigure   = &EmergencyAccessResource{}
	_ resource.ResourceWithImportState = &EmergencyAccessResource{}
)

func NewEmergencyAccessResource() resource.Resource {
	return &EmergencyAccessResource{}
}

type EmergencyAccessResource struct {
	client *client.Client
}

type EmergencyAccessModel struct {
	ID            types.String `tfsdk:"id"`
	EmailID       types.String `tfsdk:"email_id"`
	FirstName     types.String `tfsdk:"first_name"`
	LastName      types.String `tfsdk:"last_name"`
	MicrotenantID types.String `tfsdk:"microtenant_id"`
}

func (r *EmergencyAccessResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_emergency_access_user"
}

func (r *EmergencyAccessResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages ZPA emergency access users.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "Unique identifier of the emergency access user.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"email_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Email address of the emergency access user.",
			},
			"first_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "First name of the user.",
			},
			"last_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Last name of the user.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the micro-tenant scoped to this user.",
			},
		},
	}
}

func (r *EmergencyAccessResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EmergencyAccessResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan EmergencyAccessModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload := expandEmergencyAccess(plan)

	created, _, err := emergencyaccess.Create(ctx, service, &payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create emergency access user: %v", err))
		return
	}

	state, diags := r.readIntoState(ctx, service, created.UserID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EmergencyAccessResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state EmergencyAccessModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	newState, diags := r.readIntoState(ctx, service, state.ID.ValueString())
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

func (r *EmergencyAccessResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan EmergencyAccessModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload := expandEmergencyAccess(plan)

	if _, err := emergencyaccess.Update(ctx, service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update emergency access user: %v", err))
		return
	}

	state, diags := r.readIntoState(ctx, service, plan.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EmergencyAccessResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state EmergencyAccessModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	if _, err := emergencyaccess.Deactivate(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to deactivate emergency access user: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *EmergencyAccessResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
}

func (r *EmergencyAccessResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	trimmed := helpers.StringValue(microtenantID)
	if trimmed != "" {
		service = service.WithMicroTenant(trimmed)
	}
	return service
}

func (r *EmergencyAccessResource) readIntoState(ctx context.Context, service *zscaler.Service, id string) (EmergencyAccessModel, diag.Diagnostics) {
	user, _, err := emergencyaccess.Get(ctx, service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return EmergencyAccessModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Emergency access user %s not found", id))}
		}
		return EmergencyAccessModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read emergency access user: %v", err))}
	}

	microTenantID := ""
	if service != nil {
		if mt := service.MicroTenantID(); mt != nil {
			microTenantID = strings.TrimSpace(*mt)
		}
	}

	return EmergencyAccessModel{
		ID:            helpers.StringValueOrNull(user.UserID),
		EmailID:       helpers.StringValueOrNull(user.EmailID),
		FirstName:     helpers.StringValueOrNull(user.FirstName),
		LastName:      helpers.StringValueOrNull(user.LastName),
		MicrotenantID: helpers.StringValueOrNull(microTenantID),
	}, nil
}

func expandEmergencyAccess(plan EmergencyAccessModel) emergencyaccess.EmergencyAccess {
	return emergencyaccess.EmergencyAccess{
		UserID:    helpers.StringValue(plan.ID),
		EmailID:   helpers.StringValue(plan.EmailID),
		FirstName: helpers.StringValue(plan.FirstName),
		LastName:  helpers.StringValue(plan.LastName),
	}
}
