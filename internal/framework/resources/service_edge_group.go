package resources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/customerversionprofile"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgecontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgegroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/trustednetwork"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &ServiceEdgeGroupResource{}
	_ resource.ResourceWithConfigure   = &ServiceEdgeGroupResource{}
	_ resource.ResourceWithImportState = &ServiceEdgeGroupResource{}
)

func NewServiceEdgeGroupResource() resource.Resource {
	return &ServiceEdgeGroupResource{}
}

type ServiceEdgeGroupResource struct {
	client *client.Client
}

type serviceEdgeMembership struct {
	ID types.Set `tfsdk:"id"`
}

type trustedNetworkMembership struct {
	ID types.Set `tfsdk:"id"`
}

type ServiceEdgeGroupResourceModel struct {
	ID                            types.String `tfsdk:"id"`
	Name                          types.String `tfsdk:"name"`
	CityCountry                   types.String `tfsdk:"city_country"`
	CountryCode                   types.String `tfsdk:"country_code"`
	Description                   types.String `tfsdk:"description"`
	Enabled                       types.Bool   `tfsdk:"enabled"`
	IsPublic                      types.Bool   `tfsdk:"is_public"`
	Latitude                      types.String `tfsdk:"latitude"`
	Location                      types.String `tfsdk:"location"`
	Longitude                     types.String `tfsdk:"longitude"`
	OverrideVersionProfile        types.Bool   `tfsdk:"override_version_profile"`
	UseInDrMode                   types.Bool   `tfsdk:"use_in_dr_mode"`
	ServiceEdges                  types.List   `tfsdk:"service_edges"`
	TrustedNetworks               types.List   `tfsdk:"trusted_networks"`
	UpgradeDay                    types.String `tfsdk:"upgrade_day"`
	UpgradeTimeInSecs             types.String `tfsdk:"upgrade_time_in_secs"`
	VersionProfileID              types.String `tfsdk:"version_profile_id"`
	VersionProfileName            types.String `tfsdk:"version_profile_name"`
	VersionProfileVisibilityScope types.String `tfsdk:"version_profile_visibility_scope"`
	MicroTenantID                 types.String `tfsdk:"microtenant_id"`
	GraceDistanceEnabled          types.Bool   `tfsdk:"grace_distance_enabled"`
	GraceDistanceValue            types.String `tfsdk:"grace_distance_value"`
	GraceDistanceValueUnit        types.String `tfsdk:"grace_distance_value_unit"`
	AltCloud                      types.String `tfsdk:"alt_cloud"`
	SiteID                        types.String `tfsdk:"site_id"`
	SiteName                      types.String `tfsdk:"site_name"`
	GeoLocationID                 types.String `tfsdk:"geo_location_id"`
	CreationTime                  types.String `tfsdk:"creation_time"`
	ModifiedBy                    types.String `tfsdk:"modified_by"`
	ModifiedTime                  types.String `tfsdk:"modified_time"`
}

func (r *ServiceEdgeGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_edge_group"
}

func (r *ServiceEdgeGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZPA Service Edge Group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the Service Edge Group.",
			},
			"city_country": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"country_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 2),
				},
			},
			"description": schema.StringAttribute{Optional: true},
			"enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"is_public": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"latitude": schema.StringAttribute{
				Required: true,
			},
			"location": schema.StringAttribute{
				Required: true,
			},
			"longitude": schema.StringAttribute{
				Required: true,
			},
			"override_version_profile": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"use_in_dr_mode": schema.BoolAttribute{Computed: true},
			"upgrade_day": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"upgrade_time_in_secs": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"version_profile_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"version_profile_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"version_profile_visibility_scope": schema.StringAttribute{Computed: true},
			"microtenant_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"grace_distance_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"grace_distance_value": schema.StringAttribute{
				Optional: true,
			},
			"grace_distance_value_unit": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("MILES", "KMS"),
				},
			},
			"alt_cloud":       schema.StringAttribute{Computed: true},
			"site_id":         schema.StringAttribute{Computed: true},
			"site_name":       schema.StringAttribute{Computed: true},
			"geo_location_id": schema.StringAttribute{Computed: true},
			"creation_time":   schema.StringAttribute{Computed: true},
			"modified_by":     schema.StringAttribute{Computed: true},
			"modified_time":   schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			// service_edges: TypeList in SDKv2, id is TypeSet (Optional)
			// Using ListNestedBlock for block syntax support
			"service_edges": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			// trusted_networks: TypeList in SDKv2, id is TypeSet (Optional)
			// Using ListNestedBlock for block syntax support
			"trusted_networks": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (r *ServiceEdgeGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ServiceEdgeGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServiceEdgeGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.client.Service
	if !plan.MicroTenantID.IsNull() && plan.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(plan.MicroTenantID.ValueString())
	}

	if diags := r.ensureVersionProfileID(ctx, service, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	payload, diags := expandServiceEdgeGroup(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := serviceedgegroup.Create(ctx, service, payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create service edge group: %v", err))
		return
	}

	state, stateDiags := r.readServiceEdgeGroup(ctx, service, created.ID)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceEdgeGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServiceEdgeGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing ID", "Service edge group ID is required")
		return
	}

	service := r.client.Service
	if !state.MicroTenantID.IsNull() && state.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(state.MicroTenantID.ValueString())
	}

	newState, diags := r.readServiceEdgeGroup(ctx, service, state.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState.ID.IsNull() || newState.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ServiceEdgeGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServiceEdgeGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.client.Service
	if !plan.MicroTenantID.IsNull() && plan.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(plan.MicroTenantID.ValueString())
	}

	if diags := r.ensureVersionProfileID(ctx, service, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	payload, diags := expandServiceEdgeGroup(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := serviceedgegroup.Update(ctx, service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update service edge group: %v", err))
		return
	}

	state, stateDiags := r.readServiceEdgeGroup(ctx, service, plan.ID.ValueString())
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceEdgeGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServiceEdgeGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.client.Service
	if !state.MicroTenantID.IsNull() && state.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(state.MicroTenantID.ValueString())
	}

	if _, err := serviceedgegroup.Delete(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service edge group: %v", err))
		return
	}
}

func (r *ServiceEdgeGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
}

func (r *ServiceEdgeGroupResource) ensureVersionProfileID(ctx context.Context, service *zscaler.Service, plan *ServiceEdgeGroupResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if !helpers.BoolValue(plan.OverrideVersionProfile, false) {
		return diags
	}

	if !plan.VersionProfileID.IsNull() {
		idValue := plan.VersionProfileID.ValueString()
		if idValue != "" && idValue != "0" {
			return diags
		}
	}

	if plan.VersionProfileName.IsNull() || plan.VersionProfileName.ValueString() == "" {
		diags.AddError("Missing Version Profile", "version_profile_name must be provided when override_version_profile is true")
		return diags
	}

	profile, _, err := customerversionprofile.GetByName(ctx, service, plan.VersionProfileName.ValueString())
	if err != nil {
		diags.AddError("Lookup Error", fmt.Sprintf("Unable to find version profile '%s': %v", plan.VersionProfileName.ValueString(), err))
		return diags
	}

	plan.VersionProfileID = types.StringValue(profile.ID)
	plan.VersionProfileName = types.StringValue(profile.Name)
	plan.VersionProfileVisibilityScope = types.StringValue(profile.VisibilityScope)
	return diags
}

func (r *ServiceEdgeGroupResource) readServiceEdgeGroup(ctx context.Context, service *zscaler.Service, id string) (ServiceEdgeGroupResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var state ServiceEdgeGroupResourceModel

	group, _, err := serviceedgegroup.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			tflog.Warn(ctx, "Service edge group not found; removing from state", map[string]any{"id": id})
			state.ID = types.StringValue("")
			return state, diags
		}
		diags.AddError("Client Error", fmt.Sprintf("Unable to read service edge group: %v", err))
		return state, diags
	}

	state = flattenServiceEdgeGroupResource(ctx, group)
	return state, diags
}

func expandServiceEdgeGroup(ctx context.Context, plan ServiceEdgeGroupResourceModel) (serviceedgegroup.ServiceEdgeGroup, diag.Diagnostics) {
	var diags diag.Diagnostics

	serviceEdges, seDiags := expandServiceEdges(ctx, plan.ServiceEdges)
	diags.Append(seDiags...)
	if diags.HasError() {
		return serviceedgegroup.ServiceEdgeGroup{}, diags
	}

	trustedNetworks, tnDiags := expandTrustedNetworks(ctx, plan.TrustedNetworks)
	diags.Append(tnDiags...)
	if diags.HasError() {
		return serviceedgegroup.ServiceEdgeGroup{}, diags
	}

	payload := serviceedgegroup.ServiceEdgeGroup{
		ID:                            plan.ID.ValueString(),
		Name:                          plan.Name.ValueString(),
		CityCountry:                   plan.CityCountry.ValueString(),
		CountryCode:                   plan.CountryCode.ValueString(),
		Description:                   plan.Description.ValueString(),
		Enabled:                       helpers.BoolValue(plan.Enabled, true),
		IsPublic:                      strings.ToUpper(strconv.FormatBool(helpers.BoolValue(plan.IsPublic, false))),
		Latitude:                      plan.Latitude.ValueString(),
		Location:                      plan.Location.ValueString(),
		Longitude:                     plan.Longitude.ValueString(),
		OverrideVersionProfile:        helpers.BoolValue(plan.OverrideVersionProfile, false),
		UseInDrMode:                   helpers.BoolValue(plan.UseInDrMode, false),
		UpgradeDay:                    plan.UpgradeDay.ValueString(),
		UpgradeTimeInSecs:             plan.UpgradeTimeInSecs.ValueString(),
		VersionProfileID:              plan.VersionProfileID.ValueString(),
		VersionProfileName:            plan.VersionProfileName.ValueString(),
		VersionProfileVisibilityScope: plan.VersionProfileVisibilityScope.ValueString(),
		MicroTenantID:                 plan.MicroTenantID.ValueString(),
		GraceDistanceEnabled:          helpers.BoolValue(plan.GraceDistanceEnabled, false),
		GraceDistanceValue:            plan.GraceDistanceValue.ValueString(),
		GraceDistanceValueUnit:        plan.GraceDistanceValueUnit.ValueString(),
		ServiceEdges:                  serviceEdges,
		TrustedNetworks:               trustedNetworks,
	}

	return payload, diags
}

func flattenServiceEdgeGroupResource(ctx context.Context, group *serviceedgegroup.ServiceEdgeGroup) ServiceEdgeGroupResourceModel {
	seList, _ := flattenServiceEdgeMembership(ctx, group.ServiceEdges)
	tnList, _ := flattenTrustedNetworkMembership(ctx, group.TrustedNetworks)

	isPublic, _ := strconv.ParseBool(group.IsPublic)
	graceValue := group.GraceDistanceValue
	if graceValue != "" {
		if parsed, err := strconv.ParseFloat(graceValue, 64); err == nil {
			graceValue = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", parsed), "0"), ".")
		}
	}

	return ServiceEdgeGroupResourceModel{
		ID:                            types.StringValue(group.ID),
		Name:                          types.StringValue(group.Name),
		CityCountry:                   types.StringValue(group.CityCountry),
		CountryCode:                   types.StringValue(group.CountryCode),
		Description:                   types.StringValue(group.Description),
		Enabled:                       types.BoolValue(group.Enabled),
		IsPublic:                      types.BoolValue(isPublic),
		Latitude:                      types.StringValue(group.Latitude),
		Location:                      types.StringValue(group.Location),
		Longitude:                     types.StringValue(group.Longitude),
		OverrideVersionProfile:        types.BoolValue(group.OverrideVersionProfile),
		UseInDrMode:                   types.BoolValue(group.UseInDrMode),
		ServiceEdges:                  seList,
		TrustedNetworks:               tnList,
		UpgradeDay:                    types.StringValue(group.UpgradeDay),
		UpgradeTimeInSecs:             types.StringValue(group.UpgradeTimeInSecs),
		VersionProfileID:              types.StringValue(group.VersionProfileID),
		VersionProfileName:            types.StringValue(group.VersionProfileName),
		VersionProfileVisibilityScope: types.StringValue(group.VersionProfileVisibilityScope),
		MicroTenantID:                 types.StringValue(group.MicroTenantID),
		GraceDistanceEnabled:          types.BoolValue(group.GraceDistanceEnabled),
		GraceDistanceValue:            types.StringValue(graceValue),
		GraceDistanceValueUnit:        types.StringValue(group.GraceDistanceValueUnit),
		AltCloud:                      types.StringValue(group.AltCloud),
		SiteID:                        types.StringValue(group.SiteID),
		SiteName:                      types.StringValue(group.SiteName),
		GeoLocationID:                 types.StringValue(group.GeoLocationID),
		CreationTime:                  types.StringValue(group.CreationTime),
		ModifiedBy:                    types.StringValue(group.ModifiedBy),
		ModifiedTime:                  types.StringValue(group.ModifiedTime),
	}
}

func expandServiceEdges(ctx context.Context, list types.List) ([]serviceedgecontroller.ServiceEdgeController, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var memberships []serviceEdgeMembership
	diags.Append(list.ElementsAs(ctx, &memberships, false)...)
	if diags.HasError() {
		return nil, diags
	}

	var result []serviceedgecontroller.ServiceEdgeController
	for _, membership := range memberships {
		if membership.ID.IsNull() || membership.ID.IsUnknown() {
			continue
		}
		var ids []string
		diags.Append(membership.ID.ElementsAs(ctx, &ids, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for _, id := range ids {
			result = append(result, serviceedgecontroller.ServiceEdgeController{ID: id})
		}
	}

	return result, diags
}

func expandTrustedNetworks(ctx context.Context, list types.List) ([]trustednetwork.TrustedNetwork, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var memberships []trustedNetworkMembership
	diags.Append(list.ElementsAs(ctx, &memberships, false)...)
	if diags.HasError() {
		return nil, diags
	}

	var result []trustednetwork.TrustedNetwork
	for _, membership := range memberships {
		if membership.ID.IsNull() || membership.ID.IsUnknown() {
			continue
		}
		var ids []string
		diags.Append(membership.ID.ElementsAs(ctx, &ids, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for _, id := range ids {
			result = append(result, trustednetwork.TrustedNetwork{ID: id})
		}
	}

	return result, diags
}

func flattenServiceEdgeMembership(ctx context.Context, edges []serviceedgecontroller.ServiceEdgeController) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id": types.SetType{ElemType: types.StringType},
	}

	if len(edges) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), nil
	}

	ids := make([]string, len(edges))
	for i, edge := range edges {
		ids[i] = edge.ID
	}

	setValue, diags := types.SetValueFrom(ctx, types.StringType, ids)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{"id": setValue})
	if objDiags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), objDiags
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	if listDiags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), listDiags
	}

	return list, nil
}

func flattenTrustedNetworkMembership(ctx context.Context, networks []trustednetwork.TrustedNetwork) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id": types.SetType{ElemType: types.StringType},
	}

	if len(networks) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), nil
	}

	ids := make([]string, len(networks))
	for i, network := range networks {
		ids[i] = network.ID
	}

	setValue, diags := types.SetValueFrom(ctx, types.StringType, ids)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{"id": setValue})
	if objDiags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), objDiags
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	if listDiags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), listDiags
	}

	return list, nil
}
