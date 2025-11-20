package resources

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/customerversionprofile"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &AppConnectorGroupResource{}
	_ resource.ResourceWithConfigure   = &AppConnectorGroupResource{}
	_ resource.ResourceWithImportState = &AppConnectorGroupResource{}
)

var appConnectorPolicyLock sync.Mutex

func NewAppConnectorGroupResource() resource.Resource {
	return &AppConnectorGroupResource{}
}

type AppConnectorGroupResource struct {
	client *client.Client
}

type AppConnectorGroupResourceModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	CityCountry              types.String `tfsdk:"city_country"`
	CountryCode              types.String `tfsdk:"country_code"`
	Description              types.String `tfsdk:"description"`
	DNSQueryType             types.String `tfsdk:"dns_query_type"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	Latitude                 types.String `tfsdk:"latitude"`
	Location                 types.String `tfsdk:"location"`
	Longitude                types.String `tfsdk:"longitude"`
	LSSAppConnectorGroup     types.Bool   `tfsdk:"lss_app_connector_group"`
	TCPQuickAckApp           types.Bool   `tfsdk:"tcp_quick_ack_app"`
	TCPQuickAckAssistant     types.Bool   `tfsdk:"tcp_quick_ack_assistant"`
	TCPQuickAckReadAssistant types.Bool   `tfsdk:"tcp_quick_ack_read_assistant"`
	UseInDrMode              types.Bool   `tfsdk:"use_in_dr_mode"`
	PRAEnabled               types.Bool   `tfsdk:"pra_enabled"`
	WAFDisabled              types.Bool   `tfsdk:"waf_disabled"`
	OverrideVersionProfile   types.Bool   `tfsdk:"override_version_profile"`
	UpgradeDay               types.String `tfsdk:"upgrade_day"`
	UpgradeTimeInSecs        types.String `tfsdk:"upgrade_time_in_secs"`
	VersionProfileID         types.String `tfsdk:"version_profile_id"`
	VersionProfileName       types.String `tfsdk:"version_profile_name"`
	MicroTenantID            types.String `tfsdk:"microtenant_id"`
	MicroTenantName          types.String `tfsdk:"microtenant_name"`
	CreationTime             types.String `tfsdk:"creation_time"`
	ModifiedBy               types.String `tfsdk:"modifiedby"`
	ModifiedTime             types.String `tfsdk:"modified_time"`
	GeoLocationID            types.String `tfsdk:"geo_location_id"`
	VersionProfileVisibility types.String `tfsdk:"version_profile_visibility_scope"`
}

func (r *AppConnectorGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_connector_group"
}

func (r *AppConnectorGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZPA App Connector Group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the App Connector Group.",
				Required:    true,
			},
			"city_country": schema.StringAttribute{Optional: true},
			"country_code": schema.StringAttribute{Optional: true},
			"description":  schema.StringAttribute{Optional: true},
			"dns_query_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("IPV4_IPV6", "IPV4", "IPV6"),
				},
				Default: stringdefault.StaticString("IPV4_IPV6"),
			},
			"enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"latitude": schema.StringAttribute{
				Description: "Latitude of the App Connector Group.",
				Required:    true,
			},
			"location": schema.StringAttribute{
				Description: "Location description of the App Connector Group.",
				Required:    true,
			},
			"longitude": schema.StringAttribute{
				Description: "Longitude of the App Connector Group.",
				Required:    true,
			},
			"lss_app_connector_group":          schema.BoolAttribute{Optional: true, Computed: true},
			"tcp_quick_ack_app":                schema.BoolAttribute{Optional: true, Computed: true},
			"tcp_quick_ack_assistant":          schema.BoolAttribute{Optional: true, Computed: true},
			"tcp_quick_ack_read_assistant":     schema.BoolAttribute{Optional: true, Computed: true},
			"use_in_dr_mode":                   schema.BoolAttribute{Optional: true, Computed: true},
			"pra_enabled":                      schema.BoolAttribute{Optional: true, Computed: true},
			"waf_disabled":                     schema.BoolAttribute{Optional: true, Computed: true},
			"override_version_profile":         schema.BoolAttribute{Optional: true, Computed: true},
			"upgrade_day":                      schema.StringAttribute{Optional: true},
			"upgrade_time_in_secs":             schema.StringAttribute{Optional: true},
			"version_profile_id":               schema.StringAttribute{Optional: true, Computed: true},
			"version_profile_name":             schema.StringAttribute{Optional: true, Computed: true},
			"microtenant_id":                   schema.StringAttribute{Optional: true, Computed: true},
			"microtenant_name":                 schema.StringAttribute{Computed: true},
			"creation_time":                    schema.StringAttribute{Computed: true},
			"modifiedby":                       schema.StringAttribute{Computed: true},
			"modified_time":                    schema.StringAttribute{Computed: true},
			"geo_location_id":                  schema.StringAttribute{Computed: true},
			"version_profile_visibility_scope": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *AppConnectorGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AppConnectorGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AppConnectorGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicroTenantID)

	if diags := r.resolveVersionProfile(ctx, service, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	group := expandAppConnectorGroup(plan)

	if diags := validateTCPQuickAckSettings(plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	created, _, err := appconnectorgroup.Create(ctx, service, group)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create app connector group: %v", err))
		return
	}

	tflog.Info(ctx, "Created app connector group", map[string]any{"id": created.ID})

	state, diags := r.readIntoState(ctx, service, created.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AppConnectorGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AppConnectorGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing ID", "App Connector Group ID is required to read the resource")
		return
	}

	service := r.serviceForMicrotenant(state.MicroTenantID)

	updatedState, diags := r.readIntoState(ctx, service, state.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if updatedState.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *AppConnectorGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AppConnectorGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing ID", "App Connector Group ID is required to update the resource")
		return
	}

	service := r.serviceForMicrotenant(plan.MicroTenantID)

	// Check if resource still exists before updating
	if _, _, err := appconnectorgroup.Get(ctx, service, plan.ID.ValueString()); err != nil {
		if helpers.IsObjectNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
	}

	if diags := r.resolveVersionProfile(ctx, service, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if diags := validateTCPQuickAckSettings(plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	payload := expandAppConnectorGroup(plan)

	if _, err := appconnectorgroup.Update(ctx, service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update app connector group: %v", err))
		return
	}

	updatedState, diags := r.readIntoState(ctx, service, plan.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if updatedState.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *AppConnectorGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AppConnectorGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing ID", "App Connector Group ID is required to delete the resource")
		return
	}

	service := r.serviceForMicrotenant(state.MicroTenantID)

	if err := r.detachFromPolicies(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to detach app connector group from policies: %v", err))
		return
	}

	if _, err := appconnectorgroup.Delete(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete app connector group: %v", err))
		return
	}
}

func (r *AppConnectorGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	service := r.client.Service

	if _, err := strconv.Atoi(id); err != nil {
		group, _, err := appconnectorgroup.GetByName(ctx, service, id)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import app connector group by name: %v", err))
			return
		}
		id = group.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *AppConnectorGroupResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && microtenantID.ValueString() != "" {
		service = service.WithMicroTenant(microtenantID.ValueString())
	}
	return service
}

func (r *AppConnectorGroupResource) resolveVersionProfile(ctx context.Context, service *zscaler.Service, model *AppConnectorGroupResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if model.OverrideVersionProfile.IsNull() || !model.OverrideVersionProfile.ValueBool() {
		return diags
	}

	if !model.VersionProfileID.IsNull() {
		idValue := model.VersionProfileID.ValueString()
		if idValue != "" && idValue != "0" {
			return diags
		}
	}

	if model.VersionProfileName.IsNull() || model.VersionProfileName.ValueString() == "" {
		diags.AddError("Missing version profile", "version_profile_name must be provided when override_version_profile is true")
		return diags
	}

	profile, _, err := customerversionprofile.GetByName(ctx, service, model.VersionProfileName.ValueString())
	if err != nil {
		diags.AddError("Version profile lookup failed", fmt.Sprintf("Unable to find version profile %s: %v", model.VersionProfileName.ValueString(), err))
		return diags
	}

	model.VersionProfileID = types.StringValue(profile.ID)
	model.VersionProfileName = types.StringValue(profile.Name)
	return diags
}

func (r *AppConnectorGroupResource) readIntoState(ctx context.Context, service *zscaler.Service, id string) (AppConnectorGroupResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var state AppConnectorGroupResourceModel

	group, _, err := appconnectorgroup.Get(ctx, service, id)
	if err != nil {
		if helpers.IsObjectNotFoundError(err) {
			tflog.Warn(ctx, "App connector group not found; removing from state", map[string]any{"id": id})
			state.ID = types.StringValue("")
			return state, diags
		}
		diags.AddError("Client Error", fmt.Sprintf("Unable to read app connector group: %v", err))
		return state, diags
	}

	state = flattenAppConnectorGroup(group)
	return state, diags
}

func flattenAppConnectorGroup(group *appconnectorgroup.AppConnectorGroup) AppConnectorGroupResourceModel {
	return AppConnectorGroupResourceModel{
		ID:                       types.StringValue(group.ID),
		Name:                     types.StringValue(group.Name),
		CityCountry:              types.StringValue(group.CityCountry),
		CountryCode:              types.StringValue(group.CountryCode),
		Description:              types.StringValue(group.Description),
		DNSQueryType:             types.StringValue(group.DNSQueryType),
		Enabled:                  types.BoolValue(group.Enabled),
		Latitude:                 types.StringValue(group.Latitude),
		Location:                 types.StringValue(group.Location),
		Longitude:                types.StringValue(group.Longitude),
		LSSAppConnectorGroup:     types.BoolValue(group.LSSAppConnectorGroup),
		TCPQuickAckApp:           types.BoolValue(group.TCPQuickAckApp),
		TCPQuickAckAssistant:     types.BoolValue(group.TCPQuickAckAssistant),
		TCPQuickAckReadAssistant: types.BoolValue(group.TCPQuickAckReadAssistant),
		UseInDrMode:              types.BoolValue(group.UseInDrMode),
		PRAEnabled:               types.BoolValue(group.PRAEnabled),
		WAFDisabled:              types.BoolValue(group.WAFDisabled),
		OverrideVersionProfile:   types.BoolValue(group.OverrideVersionProfile),
		UpgradeDay:               types.StringValue(group.UpgradeDay),
		UpgradeTimeInSecs:        types.StringValue(group.UpgradeTimeInSecs),
		VersionProfileID:         types.StringValue(group.VersionProfileID),
		VersionProfileName:       types.StringValue(group.VersionProfileName),
		MicroTenantID:            types.StringValue(group.MicroTenantID),
		MicroTenantName:          types.StringValue(group.MicroTenantName),
		CreationTime:             types.StringValue(group.CreationTime),
		ModifiedBy:               types.StringValue(group.ModifiedBy),
		ModifiedTime:             types.StringValue(group.ModifiedTime),
		GeoLocationID:            types.StringValue(group.GeoLocationID),
		VersionProfileVisibility: types.StringValue(group.VersionProfileVisibilityScope),
	}
}

func expandAppConnectorGroup(model AppConnectorGroupResourceModel) appconnectorgroup.AppConnectorGroup {
	return appconnectorgroup.AppConnectorGroup{
		ID:                       model.ID.ValueString(),
		Name:                     model.Name.ValueString(),
		CityCountry:              model.CityCountry.ValueString(),
		CountryCode:              model.CountryCode.ValueString(),
		Description:              model.Description.ValueString(),
		DNSQueryType:             model.DNSQueryType.ValueString(),
		Enabled:                  helpers.BoolValue(model.Enabled, true),
		Latitude:                 model.Latitude.ValueString(),
		Longitude:                model.Longitude.ValueString(),
		Location:                 model.Location.ValueString(),
		LSSAppConnectorGroup:     helpers.BoolValue(model.LSSAppConnectorGroup, false),
		TCPQuickAckApp:           helpers.BoolValue(model.TCPQuickAckApp, false),
		TCPQuickAckAssistant:     helpers.BoolValue(model.TCPQuickAckAssistant, false),
		TCPQuickAckReadAssistant: helpers.BoolValue(model.TCPQuickAckReadAssistant, false),
		UseInDrMode:              helpers.BoolValue(model.UseInDrMode, false),
		PRAEnabled:               helpers.BoolValue(model.PRAEnabled, false),
		WAFDisabled:              helpers.BoolValue(model.WAFDisabled, false),
		OverrideVersionProfile:   helpers.BoolValue(model.OverrideVersionProfile, false),
		UpgradeDay:               model.UpgradeDay.ValueString(),
		UpgradeTimeInSecs:        model.UpgradeTimeInSecs.ValueString(),
		VersionProfileID:         model.VersionProfileID.ValueString(),
		VersionProfileName:       model.VersionProfileName.ValueString(),
		MicroTenantID:            model.MicroTenantID.ValueString(),
	}
}

func validateTCPQuickAckSettings(plan AppConnectorGroupResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	app := helpers.BoolValue(plan.TCPQuickAckApp, false)
	assistant := helpers.BoolValue(plan.TCPQuickAckAssistant, false)
	readAssistant := helpers.BoolValue(plan.TCPQuickAckReadAssistant, false)

	if app != assistant || app != readAssistant {
		diags.AddError("Invalid TCP Quick Ack configuration", "tcp_quick_ack_app, tcp_quick_ack_assistant, and tcp_quick_ack_read_assistant must all have the same value")
	}

	return diags
}

func (r *AppConnectorGroupResource) detachFromPolicies(ctx context.Context, service *zscaler.Service, id string) error {
	appConnectorPolicyLock.Lock()
	defer appConnectorPolicyLock.Unlock()

	policySet, _, err := policysetcontroller.GetByPolicyType(ctx, service, "ACCESS_POLICY")
	if err != nil {
		return fmt.Errorf("failed to fetch access policy set: %w", err)
	}

	rules, _, err := policysetcontroller.GetAllByType(ctx, service, "ACCESS_POLICY")
	if err != nil {
		return fmt.Errorf("failed to fetch access policy rules: %w", err)
	}

	for _, rule := range rules {
		changed := false
		updated := make([]appconnectorgroup.AppConnectorGroup, 0, len(rule.AppConnectorGroups))

		for _, group := range rule.AppConnectorGroups {
			if group.ID == id {
				changed = true
				continue
			}
			updated = append(updated, appconnectorgroup.AppConnectorGroup{ID: group.ID})
		}

		if !changed {
			continue
		}

		rule.AppConnectorGroups = updated
		if _, err := policysetcontroller.UpdateRule(ctx, service, policySet.ID, rule.ID, &rule); err != nil {
			log.Printf("[WARN] Failed to update policy rule %s: %v", rule.ID, err)
		}
	}

	return nil
}
