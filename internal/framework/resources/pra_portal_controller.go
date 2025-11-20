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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praconsole"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praportal"
)

var (
	_ resource.Resource                = &PRAPortalControllerResource{}
	_ resource.ResourceWithConfigure   = &PRAPortalControllerResource{}
	_ resource.ResourceWithImportState = &PRAPortalControllerResource{}
)

func NewPRAPortalControllerResource() resource.Resource {
	return &PRAPortalControllerResource{}
}

type PRAPortalControllerResource struct {
	client *client.Client
}

type PRAPortalControllerModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
	Domain                  types.String `tfsdk:"domain"`
	CertificateID           types.String `tfsdk:"certificate_id"`
	UserNotification        types.String `tfsdk:"user_notification"`
	UserNotificationEnabled types.Bool   `tfsdk:"user_notification_enabled"`
	ExtLabel                types.String `tfsdk:"ext_label"`
	ExtDomain               types.String `tfsdk:"ext_domain"`
	ExtDomainName           types.String `tfsdk:"ext_domain_name"`
	ExtDomainTranslation    types.String `tfsdk:"ext_domain_translation"`
	MicrotenantID           types.String `tfsdk:"microtenant_id"`
	UserPortalGid           types.String `tfsdk:"user_portal_gid"`
}

func (r *PRAPortalControllerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pra_portal_controller"
}

func (r *PRAPortalControllerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Privileged Remote Access (PRA) Portal Controller.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the privileged portal",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the privileged portal",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether or not the privileged portal is enabled",
			},
			"domain": schema.StringAttribute{
				Optional:    true,
				Description: "The domain of the privileged portal",
			},
			"certificate_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the certificate",
			},
			"user_notification": schema.StringAttribute{
				Optional:    true,
				Description: "The notification message displayed in the banner of the privileged portallink, if enabled",
			},
			"user_notification_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates if the Notification Banner is enabled (true) or disabled (false)",
			},
			"ext_label": schema.StringAttribute{
				Optional:    true,
				Description: "The domain prefix for the privileged portal URL. The supported string can include numbers, lower case characters, and only supports a hyphen (-).",
			},
			"ext_domain": schema.StringAttribute{
				Optional:    true,
				Description: "The external domain name prefix of the Browser Access application that is used for Zscaler-managed certificates when creating a privileged portal.",
			},
			"ext_domain_name": schema.StringAttribute{
				Optional:    true,
				Description: "The domain suffix for the privileged portal URL. This field must be one of the customer's authentication domains.",
			},
			"ext_domain_translation": schema.StringAttribute{
				Optional:    true,
				Description: "The translation of the external domain name prefix of the Browser Access application that is used for Zscaler-managed certificates when creating a privileged portal.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Description: "The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.",
			},
			"user_portal_gid": schema.StringAttribute{
				Optional:    true,
				Description: "The unique identifier of the user portal.",
			},
		},
	}
}

func (r *PRAPortalControllerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PRAPortalControllerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan PRAPortalControllerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := r.validatePRAPortalController(&plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload := r.expandPRAPortalController(&plan)

	praPortal, _, err := praportal.Create(ctx, service, &payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create PRA portal controller: %v", err))
		return
	}

	plan.ID = types.StringValue(praPortal.ID)

	state, readDiags := r.readPRAPortalController(ctx, service, praPortal.ID, plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = plan.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PRAPortalControllerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state PRAPortalControllerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	newState, diags := r.readPRAPortalController(ctx, service, state.ID.ValueString(), state.MicrotenantID)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	newState.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *PRAPortalControllerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan PRAPortalControllerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := r.validatePRAPortalController(&plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload := r.expandPRAPortalController(&plan)

	if _, err := praportal.Update(ctx, service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update PRA portal controller: %v", err))
		return
	}

	state, readDiags := r.readPRAPortalController(ctx, service, plan.ID.ValueString(), plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = plan.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PRAPortalControllerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state PRAPortalControllerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	portalID := state.ID.ValueString()
	service := r.serviceForMicrotenant(state.MicrotenantID)

	if err := r.detachAndCleanUpPRAPortals(ctx, portalID, service); err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Error detaching PRAPortal with ID %s from PRAConsoleControllers: %v", portalID, err),
		)
		return
	}

	if _, err := praportal.Delete(ctx, service, portalID); err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Error deleting PRA Portal Controller with ID %s: %v", portalID, err),
		)
		return
	}
}

func (r *PRAPortalControllerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before importing PRA portal controller.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the PRA portal controller ID or name.")
		return
	}

	service := r.client.Service

	if _, err := strconv.ParseInt(id, 10, 64); err == nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		return
	}

	portal, _, err := praportal.GetByName(ctx, service, id)
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate PRA portal controller %q: %v", id, err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(portal.ID))...)
}

func (r *PRAPortalControllerResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	id := helpers.StringValue(microtenantID)
	if id != "" {
		service = service.WithMicroTenant(id)
	}
	return service
}

func (r *PRAPortalControllerResource) validatePRAPortalController(plan *PRAPortalControllerModel) diag.Diagnostics {
	var diags diag.Diagnostics

	externalFields := []string{
		"ext_label",
		"ext_domain",
		"ext_domain_name",
		"ext_domain_translation",
		"user_portal_gid",
	}

	externalSet := false
	for _, fieldName := range externalFields {
		var fieldValue types.String
		switch fieldName {
		case "ext_label":
			fieldValue = plan.ExtLabel
		case "ext_domain":
			fieldValue = plan.ExtDomain
		case "ext_domain_name":
			fieldValue = plan.ExtDomainName
		case "ext_domain_translation":
			fieldValue = plan.ExtDomainTranslation
		case "user_portal_gid":
			fieldValue = plan.UserPortalGid
		}

		if !fieldValue.IsNull() && !fieldValue.IsUnknown() && fieldValue.ValueString() != "" {
			externalSet = true
			break
		}
	}

	if externalSet {
		if !plan.CertificateID.IsNull() && !plan.CertificateID.IsUnknown() && plan.CertificateID.ValueString() != "" {
			diags.AddError(
				"Invalid Configuration",
				"'certificate_id' cannot be set when any of the following are configured: ext_label, ext_domain, ext_domain_name, ext_domain_translation, user_portal_gid",
			)
		}
	}

	return diags
}

func (r *PRAPortalControllerResource) expandPRAPortalController(plan *PRAPortalControllerModel) praportal.PRAPortal {
	return praportal.PRAPortal{
		ID:                      plan.ID.ValueString(),
		Name:                    helpers.StringValue(plan.Name),
		Description:             helpers.StringValue(plan.Description),
		Enabled:                 helpers.BoolValue(plan.Enabled, false),
		Domain:                  helpers.StringValue(plan.Domain),
		CertificateID:           helpers.StringValue(plan.CertificateID),
		MicroTenantID:           helpers.StringValue(plan.MicrotenantID),
		UserNotification:        helpers.StringValue(plan.UserNotification),
		UserNotificationEnabled: helpers.BoolValue(plan.UserNotificationEnabled, false),
		ExtLabel:                helpers.StringValue(plan.ExtLabel),
		ExtDomain:               helpers.StringValue(plan.ExtDomain),
		ExtDomainName:           helpers.StringValue(plan.ExtDomainName),
		ExtDomainTranslation:    helpers.StringValue(plan.ExtDomainTranslation),
		UserPortalGid:           helpers.StringValue(plan.UserPortalGid),
	}
}

func (r *PRAPortalControllerResource) readPRAPortalController(ctx context.Context, service *zscaler.Service, id string, microtenantID types.String) (PRAPortalControllerModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	portal, _, err := praportal.Get(ctx, service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return PRAPortalControllerModel{}, diag.Diagnostics{
				diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("PRA portal controller %s not found", id)),
			}
		}
		return PRAPortalControllerModel{}, diag.Diagnostics{
			diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read PRA portal controller: %v", err)),
		}
	}

	return PRAPortalControllerModel{
		ID:                      types.StringValue(portal.ID),
		Name:                    helpers.StringValueOrNull(portal.Name),
		Description:             helpers.StringValueOrNull(portal.Description),
		Enabled:                 types.BoolValue(portal.Enabled),
		Domain:                  helpers.StringValueOrNull(portal.Domain),
		CertificateID:           helpers.StringValueOrNull(portal.CertificateID),
		MicrotenantID:           helpers.StringValueOrNull(portal.MicroTenantID),
		UserNotification:        helpers.StringValueOrNull(portal.UserNotification),
		UserNotificationEnabled: types.BoolValue(portal.UserNotificationEnabled),
		ExtLabel:                helpers.StringValueOrNull(portal.ExtLabel),
		ExtDomain:               helpers.StringValueOrNull(portal.ExtDomain),
		ExtDomainName:           helpers.StringValueOrNull(portal.ExtDomainName),
		ExtDomainTranslation:    helpers.StringValueOrNull(portal.ExtDomainTranslation),
		UserPortalGid:           helpers.StringValueOrNull(portal.UserPortalGid),
	}, diags
}

func (r *PRAPortalControllerResource) detachAndCleanUpPRAPortals(ctx context.Context, portalID string, service *zscaler.Service) error {
	consoles, _, err := praconsole.GetAll(ctx, service)
	if err != nil {
		return fmt.Errorf("failed to list all PRAConsoleControllers: %w", err)
	}

	for _, console := range consoles {
		var portalFound bool
		for _, portal := range console.PRAPortals {
			if portal.ID == portalID {
				portalFound = true
				break
			}
		}

		if portalFound {
			updatedPortals := []praconsole.PRAPortals{}
			for _, portal := range console.PRAPortals {
				if portal.ID != portalID {
					updatedPortals = append(updatedPortals, portal)
				}
			}

			if len(updatedPortals) == 0 {
				if _, err := praconsole.Delete(ctx, service, console.ID); err != nil {
					return fmt.Errorf("failed to delete PRAConsoleController with ID %s: %w", console.ID, err)
				}
			} else {
				console.PRAPortals = updatedPortals
				if _, err := praconsole.Update(ctx, service, console.ID, &console); err != nil {
					return fmt.Errorf("failed to update PRAConsoleController with ID %s: %w", console.ID, err)
				}
			}
		}
	}

	return nil
}
