package resources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	aupservice "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/aup"
)

var (
	_ resource.Resource                = &UserPortalAUPResource{}
	_ resource.ResourceWithConfigure   = &UserPortalAUPResource{}
	_ resource.ResourceWithImportState = &UserPortalAUPResource{}
)

func NewUserPortalAUPResource() resource.Resource {
	return &UserPortalAUPResource{}
}

type UserPortalAUPResource struct {
	client *client.Client
}

type UserPortalAUPModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	AUP           types.String `tfsdk:"aup"`
	Email         types.String `tfsdk:"email"`
	PhoneNum      types.String `tfsdk:"phone_num"`
	MicrotenantID types.String `tfsdk:"microtenant_id"`
}

func (r *UserPortalAUPResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_portal_aup"
}

func (r *UserPortalAUPResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages ZPA User Portal Acceptable Use Policy (AUP).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"enabled": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"aup": schema.StringAttribute{
				Optional: true,
			},
			"email": schema.StringAttribute{
				Optional: true,
			},
			"phone_num": schema.StringAttribute{
				Optional: true,
			},
			"microtenant_id": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *UserPortalAUPResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserPortalAUPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan UserPortalAUPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload := expandUserPortalAUP(&plan)
	created, _, err := aupservice.Create(ctx, service, &payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create user portal AUP: %v", err))
		return
	}

	state, diags := r.readIntoState(ctx, created.ID, plan.MicrotenantID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserPortalAUPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state UserPortalAUPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readIntoState(ctx, state.ID.ValueString(), state.MicrotenantID)
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

func (r *UserPortalAUPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan UserPortalAUPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload := expandUserPortalAUP(&plan)
	if _, err := aupservice.Update(ctx, service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update user portal AUP: %v", err))
		return
	}

	state, diags := r.readIntoState(ctx, plan.ID.ValueString(), plan.MicrotenantID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserPortalAUPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state UserPortalAUPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	if _, err := aupservice.Delete(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete user portal AUP: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *UserPortalAUPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before importing resources.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the user portal AUP ID or name.")
		return
	}

	service := r.client.Service
	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		resource, _, lookupErr := aupservice.GetByName(ctx, service, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate user portal AUP %q: %v", id, lookupErr))
			return
		}
		id = resource.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *UserPortalAUPResource) readIntoState(ctx context.Context, id string, microtenantID types.String) (UserPortalAUPModel, diag.Diagnostics) {
	service := r.serviceForMicrotenant(microtenantID)

	resource, _, err := aupservice.Get(ctx, service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return UserPortalAUPModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("User portal AUP %s not found", id))}
		}
		return UserPortalAUPModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read user portal AUP: %v", err))}
	}

	state := UserPortalAUPModel{
		ID:            helpers.StringValueOrNull(resource.ID),
		Name:          helpers.StringValueOrNull(resource.Name),
		Description:   helpers.StringValueOrNull(resource.Description),
		Enabled:       types.BoolValue(resource.Enabled),
		AUP:           helpers.StringValueOrNull(resource.Aup),
		Email:         helpers.StringValueOrNull(resource.Email),
		PhoneNum:      helpers.StringValueOrNull(resource.PhoneNum),
		MicrotenantID: helpers.StringValueOrNull(resource.MicrotenantID),
	}

	return state, nil
}

func expandUserPortalAUP(model *UserPortalAUPModel) aupservice.UserPortalAup {
	return aupservice.UserPortalAup{
		ID:            helpers.StringValue(model.ID),
		Name:          helpers.StringValue(model.Name),
		Description:   helpers.StringValue(model.Description),
		Enabled:       helpers.BoolValue(model.Enabled, false),
		Aup:           helpers.StringValue(model.AUP),
		Email:         helpers.StringValue(model.Email),
		PhoneNum:      helpers.StringValue(model.PhoneNum),
		MicrotenantID: helpers.StringValue(model.MicrotenantID),
	}
}

func (r *UserPortalAUPResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	if r.client == nil {
		return nil
	}
	service := r.client.Service
	id := helpers.StringValue(microtenantID)
	if id != "" {
		service = service.WithMicroTenant(id)
	}
	return service
}
