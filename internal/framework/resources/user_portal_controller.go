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
	controller "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
)

var (
	_ resource.Resource                = &UserPortalControllerResource{}
	_ resource.ResourceWithConfigure   = &UserPortalControllerResource{}
	_ resource.ResourceWithImportState = &UserPortalControllerResource{}
)

func NewUserPortalControllerResource() resource.Resource {
	return &UserPortalControllerResource{}
}

type UserPortalControllerResource struct {
	client *client.Client
}

type UserPortalControllerModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	CertificateID           types.String `tfsdk:"certificate_id"`
	Description             types.String `tfsdk:"description"`
	Domain                  types.String `tfsdk:"domain"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
	ExtDomain               types.String `tfsdk:"ext_domain"`
	ExtDomainName           types.String `tfsdk:"ext_domain_name"`
	ExtDomainTranslation    types.String `tfsdk:"ext_domain_translation"`
	ExtLabel                types.String `tfsdk:"ext_label"`
	MicrotenantID           types.String `tfsdk:"microtenant_id"`
	UserNotification        types.String `tfsdk:"user_notification"`
	UserNotificationEnabled types.Bool   `tfsdk:"user_notification_enabled"`
}

func (r *UserPortalControllerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_portal_controller"
}

func (r *UserPortalControllerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages ZPA User Portal Controllers.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the User Portal Controller.",
			},
			"certificate_id": schema.StringAttribute{
				Optional:    true,
				Description: "Certificate ID associated with the controller.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the User Portal Controller.",
			},
			"domain": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Domain of the User Portal Controller.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether this User Portal Controller is enabled.",
			},
			"ext_domain": schema.StringAttribute{
				Optional:    true,
				Description: "External domain for the controller.",
			},
			"ext_domain_name": schema.StringAttribute{
				Optional:    true,
				Description: "External domain name for the controller.",
			},
			"ext_domain_translation": schema.StringAttribute{
				Optional:    true,
				Description: "External domain translation.",
			},
			"ext_label": schema.StringAttribute{
				Optional:    true,
				Description: "External label.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Description: "Microtenant ID to scope the controller.",
			},
			"user_notification": schema.StringAttribute{
				Optional:    true,
				Description: "User notification message.",
			},
			"user_notification_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether user notifications are enabled.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *UserPortalControllerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserPortalControllerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan UserPortalControllerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload := expandUserPortalController(&plan)
	created, _, err := controller.Create(ctx, service, payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create user portal controller: %v", err))
		return
	}

	state, diags := r.readIntoState(ctx, created.ID, plan.MicrotenantID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserPortalControllerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state UserPortalControllerModel
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

func (r *UserPortalControllerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan UserPortalControllerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	// Check if resource still exists before updating
	if _, _, err := controller.Get(ctx, service, plan.ID.ValueString()); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
	}

	payload := expandUserPortalController(&plan)
	if _, err := controller.Update(ctx, service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update user portal controller: %v", err))
		return
	}

	state, diags := r.readIntoState(ctx, plan.ID.ValueString(), plan.MicrotenantID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserPortalControllerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state UserPortalControllerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	if _, err := controller.Delete(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete user portal controller: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *UserPortalControllerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before importing resources.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the user portal controller ID or name.")
		return
	}

	service := r.client.Service
	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		resource, _, lookupErr := controller.GetByName(ctx, service, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate user portal controller %q: %v", id, lookupErr))
			return
		}
		id = resource.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *UserPortalControllerResource) readIntoState(ctx context.Context, id string, microtenantID types.String) (UserPortalControllerModel, diag.Diagnostics) {
	service := r.serviceForMicrotenant(microtenantID)

	resource, _, err := controller.Get(ctx, service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return UserPortalControllerModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("User portal controller %s not found", id))}
		}
		return UserPortalControllerModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read user portal controller: %v", err))}
	}

	model := UserPortalControllerModel{
		ID:                      helpers.StringValueOrNull(resource.ID),
		Name:                    helpers.StringValueOrNull(resource.Name),
		CertificateID:           helpers.StringValueOrNull(resource.CertificateId),
		Description:             helpers.StringValueOrNull(resource.Description),
		Domain:                  helpers.StringValueOrNull(resource.Domain),
		Enabled:                 types.BoolValue(resource.Enabled),
		ExtDomain:               helpers.StringValueOrNull(resource.ExtDomain),
		ExtDomainName:           helpers.StringValueOrNull(resource.ExtDomainName),
		ExtDomainTranslation:    helpers.StringValueOrNull(resource.ExtDomainTranslation),
		ExtLabel:                helpers.StringValueOrNull(resource.ExtLabel),
		MicrotenantID:           helpers.StringValueOrNull(resource.MicrotenantId),
		UserNotification:        helpers.StringValueOrNull(resource.UserNotification),
		UserNotificationEnabled: types.BoolValue(resource.UserNotificationEnabled),
	}

	return model, nil
}

func expandUserPortalController(model *UserPortalControllerModel) controller.UserPortalController {
	return controller.UserPortalController{
		ID:                      helpers.StringValue(model.ID),
		Name:                    helpers.StringValue(model.Name),
		CertificateId:           helpers.StringValue(model.CertificateID),
		Description:             helpers.StringValue(model.Description),
		Domain:                  helpers.StringValue(model.Domain),
		Enabled:                 helpers.BoolValue(model.Enabled, false),
		ExtDomain:               helpers.StringValue(model.ExtDomain),
		ExtDomainName:           helpers.StringValue(model.ExtDomainName),
		ExtDomainTranslation:    helpers.StringValue(model.ExtDomainTranslation),
		ExtLabel:                helpers.StringValue(model.ExtLabel),
		MicrotenantId:           helpers.StringValue(model.MicrotenantID),
		UserNotification:        helpers.StringValue(model.UserNotification),
		UserNotificationEnabled: helpers.BoolValue(model.UserNotificationEnabled, false),
	}
}

func (r *UserPortalControllerResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
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
