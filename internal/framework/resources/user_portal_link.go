package resources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	controller "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
	linksvc "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_link"
)

var (
	_ resource.Resource                = &UserPortalLinkResource{}
	_ resource.ResourceWithConfigure   = &UserPortalLinkResource{}
	_ resource.ResourceWithImportState = &UserPortalLinkResource{}
)

func NewUserPortalLinkResource() resource.Resource {
	return &UserPortalLinkResource{}
}

type UserPortalLinkResource struct {
	client *client.Client
}

type UserPortalLinkModel struct {
	ID            types.String               `tfsdk:"id"`
	Name          types.String               `tfsdk:"name"`
	Description   types.String               `tfsdk:"description"`
	Enabled       types.Bool                 `tfsdk:"enabled"`
	IconText      types.String               `tfsdk:"icon_text"`
	Link          types.String               `tfsdk:"link"`
	LinkPath      types.String               `tfsdk:"link_path"`
	Protocol      types.String               `tfsdk:"protocol"`
	MicrotenantID types.String               `tfsdk:"microtenant_id"`
	UserPortals   []UserPortalReferenceModel `tfsdk:"user_portals"`
}

type UserPortalReferenceModel struct {
	IDs types.Set `tfsdk:"id"`
}

func (r *UserPortalLinkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_portal_link"
}

func (r *UserPortalLinkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages ZPA User Portal Links.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the User Portal Link.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the User Portal Link.",
			},
			"enabled": schema.BoolAttribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Description:   "Whether the link is enabled.",
			},
			"icon_text": schema.StringAttribute{
				Optional:    true,
				Description: "Icon text displayed for the link.",
			},
			"link": schema.StringAttribute{
				Optional:    true,
				Description: "Link URL.",
			},
			"link_path": schema.StringAttribute{
				Optional:    true,
				Description: "Link path.",
			},
			"protocol": schema.StringAttribute{
				Optional:    true,
				Description: "Protocol used for the link.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Description: "Microtenant ID for scoping operations.",
			},
		},
		Blocks: map[string]schema.Block{
			"user_portals": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *UserPortalLinkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserPortalLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan UserPortalLinkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload, diags := expandUserPortalLink(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := linksvc.Create(ctx, service, payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create user portal link: %v", err))
		return
	}

	state, readDiags := r.readIntoState(ctx, created.ID, plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserPortalLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state UserPortalLinkModel
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

func (r *UserPortalLinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan UserPortalLinkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload, diags := expandUserPortalLink(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := linksvc.Update(ctx, service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update user portal link: %v", err))
		return
	}

	state, readDiags := r.readIntoState(ctx, plan.ID.ValueString(), plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserPortalLinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state UserPortalLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	if _, err := linksvc.Delete(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete user portal link: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *UserPortalLinkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before importing resources.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the user portal link ID or name.")
		return
	}

	service := r.client.Service
	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		resource, _, lookupErr := linksvc.GetByName(ctx, service, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate user portal link %q: %v", id, lookupErr))
			return
		}
		id = resource.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *UserPortalLinkResource) readIntoState(ctx context.Context, id string, microtenantID types.String) (UserPortalLinkModel, diag.Diagnostics) {
	service := r.serviceForMicrotenant(microtenantID)

	resource, _, err := linksvc.Get(ctx, service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return UserPortalLinkModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("User portal link %s not found", id))}
		}
		return UserPortalLinkModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read user portal link: %v", err))}
	}

	userPortalModels, diags := flattenUserPortalReferences(ctx, resource.UserPortals)

	state := UserPortalLinkModel{
		ID:            helpers.StringValueOrNull(resource.ID),
		Name:          helpers.StringValueOrNull(resource.Name),
		Description:   helpers.StringValueOrNull(resource.Description),
		Enabled:       types.BoolValue(resource.Enabled),
		IconText:      helpers.StringValueOrNull(resource.IconText),
		Link:          helpers.StringValueOrNull(resource.Link),
		LinkPath:      helpers.StringValueOrNull(resource.LinkPath),
		Protocol:      helpers.StringValueOrNull(resource.Protocol),
		MicrotenantID: helpers.StringValueOrNull(resource.MicrotenantID),
		UserPortals:   userPortalModels,
	}

	return state, diags
}

func expandUserPortalLink(ctx context.Context, model *UserPortalLinkModel) (linksvc.UserPortalLink, diag.Diagnostics) {
	var diags diag.Diagnostics

	userPortals, userDiags := expandUserPortalReferences(ctx, model.UserPortals)
	diags.Append(userDiags...)

	result := linksvc.UserPortalLink{
		ID:            helpers.StringValue(model.ID),
		Name:          helpers.StringValue(model.Name),
		Description:   helpers.StringValue(model.Description),
		Enabled:       helpers.BoolValue(model.Enabled, false),
		IconText:      helpers.StringValue(model.IconText),
		Link:          helpers.StringValue(model.Link),
		LinkPath:      helpers.StringValue(model.LinkPath),
		Protocol:      helpers.StringValue(model.Protocol),
		MicrotenantID: helpers.StringValue(model.MicrotenantID),
		UserPortals:   userPortals,
	}

	return result, diags
}

func expandUserPortalReferences(ctx context.Context, models []UserPortalReferenceModel) ([]controller.UserPortalController, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(models) == 0 {
		return nil, nil
	}

	refSet := make(map[string]struct{})

	for _, model := range models {
		ids, idsDiags := helpers.SetValueToStringSlice(ctx, model.IDs)
		diags.Append(idsDiags...)
		for _, id := range ids {
			id = strings.TrimSpace(id)
			if id == "" {
				continue
			}
			refSet[id] = struct{}{}
		}
	}

	if len(refSet) == 0 {
		return nil, diags
	}

	result := make([]controller.UserPortalController, 0, len(refSet))
	for id := range refSet {
		result = append(result, controller.UserPortalController{ID: id})
	}

	return result, diags
}

func flattenUserPortalReferences(ctx context.Context, portals []controller.UserPortalController) ([]UserPortalReferenceModel, diag.Diagnostics) {
	if len(portals) == 0 {
		return nil, nil
	}

	idValues := make([]string, 0, len(portals))
	for _, portal := range portals {
		if strings.TrimSpace(portal.ID) != "" {
			idValues = append(idValues, portal.ID)
		}
	}

	setValue, diags := types.SetValueFrom(ctx, types.StringType, idValues)
	if diags.HasError() {
		return nil, diags
	}

	return []UserPortalReferenceModel{{IDs: setValue}}, nil
}

func (r *UserPortalLinkResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
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
