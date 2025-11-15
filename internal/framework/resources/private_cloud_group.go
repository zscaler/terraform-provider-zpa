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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	privateCloud "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_group"
)

var (
	_ resource.Resource                = &PrivateCloudGroupResource{}
	_ resource.ResourceWithConfigure   = &PrivateCloudGroupResource{}
	_ resource.ResourceWithImportState = &PrivateCloudGroupResource{}
)

func NewPrivateCloudGroupResource() resource.Resource {
	return &PrivateCloudGroupResource{}
}

type PrivateCloudGroupResource struct {
	client *client.Client
}

type PrivateCloudGroupModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	CityCountry            types.String `tfsdk:"city_country"`
	CountryCode            types.String `tfsdk:"country_code"`
	Description            types.String `tfsdk:"description"`
	Enabled                types.Bool   `tfsdk:"enabled"`
	IsPublic               types.String `tfsdk:"is_public"`
	Latitude               types.String `tfsdk:"latitude"`
	Location               types.String `tfsdk:"location"`
	Longitude              types.String `tfsdk:"longitude"`
	OverrideVersionProfile types.Bool   `tfsdk:"override_version_profile"`
	MicrotenantID          types.String `tfsdk:"microtenant_id"`
	SiteID                 types.String `tfsdk:"site_id"`
	UpgradeDay             types.String `tfsdk:"upgrade_day"`
	UpgradeTimeInSecs      types.String `tfsdk:"upgrade_time_in_secs"`
	VersionProfileID       types.String `tfsdk:"version_profile_id"`
}

func (r *PrivateCloudGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_cloud_group"
}

func (r *PrivateCloudGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages ZPA Private Cloud Controller groups.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the Private Cloud Group.",
			},
			"city_country": schema.StringAttribute{
				Optional:    true,
				Description: "City and country of the Private Cloud Group.",
			},
			"country_code": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Country code of the Private Cloud Group.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the Private Cloud Group.",
			},
			"enabled": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Description:   "Whether this Private Cloud Group is enabled.",
			},
			"is_public": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether the Private Cloud Group is public.",
			},
			"latitude": schema.StringAttribute{
				Optional:    true,
				Description: "Latitude of the Private Cloud Group.",
			},
			"location": schema.StringAttribute{
				Optional:    true,
				Description: "Location description of the Private Cloud Group.",
			},
			"longitude": schema.StringAttribute{
				Optional:    true,
				Description: "Longitude of the Private Cloud Group.",
			},
			"override_version_profile": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the default version profile is overridden.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Description: "Microtenant ID to scope the request.",
			},
			"site_id": schema.StringAttribute{
				Optional:    true,
				Description: "Site ID associated with the Private Cloud Group.",
			},
			"upgrade_day": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("SUNDAY", "MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY"),
				},
				Description: "Day when controllers attempt to upgrade.",
			},
			"upgrade_time_in_secs": schema.StringAttribute{
				Optional:    true,
				Description: "Time in seconds when controllers attempt to upgrade.",
			},
			"version_profile_id": schema.StringAttribute{
				Optional:    true,
				Description: "Version profile ID assigned to the group.",
			},
		},
	}
}

func (r *PrivateCloudGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PrivateCloudGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan PrivateCloudGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload, diags := expandPrivateCloudGroup(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := privateCloud.Create(ctx, service, payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create private cloud group: %v", err))
		return
	}

	state, readDiags := r.readIntoState(ctx, created.ID, plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PrivateCloudGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state PrivateCloudGroupModel
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

func (r *PrivateCloudGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan PrivateCloudGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload, diags := expandPrivateCloudGroup(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := privateCloud.Update(ctx, service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update private cloud group: %v", err))
		return
	}

	state, readDiags := r.readIntoState(ctx, plan.ID.ValueString(), plan.MicrotenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PrivateCloudGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state PrivateCloudGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	if _, err := privateCloud.Delete(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete private cloud group: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *PrivateCloudGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before importing resources.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the private cloud group ID or name.")
		return
	}

	service := r.client.Service

	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		resource, _, lookupErr := privateCloud.GetByName(ctx, service, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate private cloud group %q: %v", id, lookupErr))
			return
		}
		id = resource.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *PrivateCloudGroupResource) readIntoState(ctx context.Context, id string, microtenantID types.String) (PrivateCloudGroupModel, diag.Diagnostics) {
	service := r.serviceForMicrotenant(microtenantID)

	resource, _, err := privateCloud.Get(ctx, service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return PrivateCloudGroupModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Private cloud group %s not found", id))}
		}
		return PrivateCloudGroupModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read private cloud group: %v", err))}
	}

	var diags diag.Diagnostics

	model := PrivateCloudGroupModel{
		ID:                     helpers.StringValueOrNull(resource.ID),
		Name:                   helpers.StringValueOrNull(resource.Name),
		CityCountry:            helpers.StringValueOrNull(resource.CityCountry),
		CountryCode:            helpers.StringValueOrNull(resource.CountryCode),
		Description:            helpers.StringValueOrNull(resource.Description),
		Enabled:                types.BoolValue(resource.Enabled),
		IsPublic:               helpers.StringValueOrNull(resource.IsPublic),
		Latitude:               helpers.StringValueOrNull(resource.Latitude),
		Location:               helpers.StringValueOrNull(resource.Location),
		Longitude:              helpers.StringValueOrNull(resource.Longitude),
		OverrideVersionProfile: types.BoolValue(resource.OverrideVersionProfile),
		MicrotenantID:          helpers.StringValueOrNull(resource.MicrotenantID),
		SiteID:                 helpers.StringValueOrNull(resource.SiteID),
		UpgradeDay:             helpers.StringValueOrNull(resource.UpgradeDay),
		UpgradeTimeInSecs:      helpers.StringValueOrNull(resource.UpgradeTimeInSecs),
		VersionProfileID:       helpers.StringValueOrNull(resource.VersionProfileID),
	}

	return model, diags
}

func expandPrivateCloudGroup(ctx context.Context, model *PrivateCloudGroupModel) (privateCloud.PrivateCloudGroup, diag.Diagnostics) {
	var diags diag.Diagnostics

	latitude := helpers.StringValue(model.Latitude)
	if latitude != "" {
		if _, errs := helpers.ValidateLatitude(latitude, "latitude"); len(errs) > 0 {
			diags.AddError("Validation Error", errs[0].Error())
		}
	}

	longitude := helpers.StringValue(model.Longitude)
	if longitude != "" {
		if _, errs := helpers.ValidateLongitude(longitude, "longitude"); len(errs) > 0 {
			diags.AddError("Validation Error", errs[0].Error())
		}
	}

	result := privateCloud.PrivateCloudGroup{
		ID:                     helpers.StringValue(model.ID),
		Name:                   helpers.StringValue(model.Name),
		CityCountry:            helpers.StringValue(model.CityCountry),
		CountryCode:            helpers.StringValue(model.CountryCode),
		Description:            helpers.StringValue(model.Description),
		Enabled:                helpers.BoolValue(model.Enabled, false),
		IsPublic:               helpers.StringValue(model.IsPublic),
		Latitude:               latitude,
		Location:               helpers.StringValue(model.Location),
		Longitude:              longitude,
		OverrideVersionProfile: helpers.BoolValue(model.OverrideVersionProfile, false),
		MicrotenantID:          helpers.StringValue(model.MicrotenantID),
		SiteID:                 helpers.StringValue(model.SiteID),
		UpgradeDay:             helpers.StringValue(model.UpgradeDay),
		UpgradeTimeInSecs:      helpers.StringValue(model.UpgradeTimeInSecs),
		VersionProfileID:       helpers.StringValue(model.VersionProfileID),
	}

	return result, diags
}

func (r *PrivateCloudGroupResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
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
