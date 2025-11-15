package resources

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &ServerGroupResource{}
	_ resource.ResourceWithConfigure   = &ServerGroupResource{}
	_ resource.ResourceWithImportState = &ServerGroupResource{}
)

var serverGroupDetachLock sync.Mutex

func NewServerGroupResource() resource.Resource {
	return &ServerGroupResource{}
}

type ServerGroupResource struct {
	client *client.Client
}

type serverGroupServiceEdges struct {
	ID types.Set `tfsdk:"id"`
}

type serverGroupAppConnectorGroups struct {
	ID types.Set `tfsdk:"id"`
}

type serverGroupServers struct {
	ID types.Set `tfsdk:"id"`
}

type serverGroupExtranetLocation struct {
	ID types.String `tfsdk:"id"`
}

type serverGroupExtranetDTO struct {
	ZPNERID          types.String `tfsdk:"zpn_er_id"`
	LocationDTO      types.Set    `tfsdk:"location_dto"`
	LocationGroupDTO types.Set    `tfsdk:"location_group_dto"`
}

type ServerGroupResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	ConfigSpace       types.String `tfsdk:"config_space"`
	IPAanchored       types.Bool   `tfsdk:"ip_anchored"`
	DynamicDiscovery  types.Bool   `tfsdk:"dynamic_discovery"`
	MicroTenantID     types.String `tfsdk:"microtenant_id"`
	AppConnectorGroup types.List   `tfsdk:"app_connector_groups"` // TypeList in SDKv2
	Servers           types.List   `tfsdk:"servers"`              // TypeList in SDKv2
	ExtranetEnabled   types.Bool   `tfsdk:"extranet_enabled"`
	ExtranetDTO       types.Set    `tfsdk:"extranet_dto"`
}

func (r *ServerGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_group"
}

func (r *ServerGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// app_connector_groups: TypeList in SDKv2, id is TypeSet (Required)
	// Using ListNestedBlock for block syntax support
	appConnectorGroupsBlock := schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.SetAttribute{
					ElementType: types.StringType,
					Required:    true,
				},
			},
		},
	}

	// servers: TypeList in SDKv2, id is TypeSet (Required)
	// Using ListNestedBlock for block syntax support
	serversBlock := schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.SetAttribute{
					ElementType: types.StringType,
					Required:    true,
				},
			},
		},
	}

	extranetLocationBlock := schema.SetNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{Required: true},
			},
		},
	}

	extranetBlock := schema.SetNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"zpn_er_id": schema.StringAttribute{Optional: true},
			},
			Blocks: map[string]schema.Block{
				"location_dto":       extranetLocationBlock,
				"location_group_dto": extranetLocationBlock,
			},
		},
	}

	resp.Schema = schema.Schema{
		Description: "Manages a ZPA Server Group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the server group.",
			},
			"description": schema.StringAttribute{Optional: true},
			"enabled":     schema.BoolAttribute{Optional: true},
			"config_space": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("DEFAULT"),
				Validators: []validator.String{
					stringvalidator.OneOf("DEFAULT", "SIEM"),
				},
			},
			"ip_anchored":       schema.BoolAttribute{Optional: true},
			"dynamic_discovery": schema.BoolAttribute{Optional: true},
			"microtenant_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"extranet_enabled": schema.BoolAttribute{Optional: true},
		},
		Blocks: map[string]schema.Block{
			"app_connector_groups": appConnectorGroupsBlock,
			"servers":              serversBlock,
			"extranet_dto":         extranetBlock,
		},
	}
}

func (r *ServerGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ServerGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServerGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.client.Service
	if !plan.MicroTenantID.IsNull() && plan.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(plan.MicroTenantID.ValueString())
	}

	payload, diags := expandServerGroup(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate dynamic_discovery and servers relationship (SDKv2 behavior)
	dynamicDiscovery := helpers.BoolValue(plan.DynamicDiscovery, false)
	if dynamicDiscovery && len(payload.Servers) > 0 {
		resp.Diagnostics.AddError("Invalid Configuration", "an application server can only be attached to a server when DynamicDiscovery is disabled")
		return
	}
	if !dynamicDiscovery && len(payload.Servers) == 0 {
		resp.Diagnostics.AddError("Invalid Configuration", "servers must not be empty when DynamicDiscovery is disabled")
		return
	}

	created, _, err := servergroup.Create(ctx, service, &payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create server group: %v", err))
		return
	}

	state, readDiags := r.readServerGroup(ctx, service, created.ID, plan)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServerGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServerGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing ID", "Server group ID is required")
		return
	}

	service := r.client.Service
	if !state.MicroTenantID.IsNull() && state.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(state.MicroTenantID.ValueString())
	}

	newState, readDiags := r.readServerGroup(ctx, service, state.ID.ValueString(), state)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState.ID.IsNull() || newState.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ServerGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServerGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ServerGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.client.Service
	if !plan.MicroTenantID.IsNull() && plan.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(plan.MicroTenantID.ValueString())
	}

	// Check if servers or dynamic_discovery changed (SDKv2 behavior)
	serversChanged := !plan.Servers.Equal(state.Servers)
	dynamicDiscoveryChanged := !plan.DynamicDiscovery.Equal(state.DynamicDiscovery)

	payload, diags := expandServerGroup(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Only validate if servers or dynamic_discovery changed (SDKv2 behavior)
	if serversChanged || dynamicDiscoveryChanged {
		dynamicDiscovery := helpers.BoolValue(plan.DynamicDiscovery, false)
		if dynamicDiscovery && len(payload.Servers) > 0 {
			resp.Diagnostics.AddError("Invalid Configuration", "an application server can only be attached to a server when DynamicDiscovery is disabled")
			return
		}
		if !dynamicDiscovery && len(payload.Servers) == 0 {
			resp.Diagnostics.AddError("Invalid Configuration", "servers must not be empty when DynamicDiscovery is disabled")
			return
		}
	}

	if _, err := servergroup.Update(ctx, service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update server group: %v", err))
		return
	}

	updatedState, readDiags := r.readServerGroup(ctx, service, plan.ID.ValueString(), plan)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *ServerGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServerGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.client.Service
	if !state.MicroTenantID.IsNull() && state.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(state.MicroTenantID.ValueString())
	}

	if err := detachServerGroupFromAppConnectorGroups(ctx, state.ID.ValueString(), service); err != nil {
		tflog.Warn(ctx, "Failed to detach server group from app connector groups", map[string]any{"error": err.Error()})
	}

	detachServerGroupFromPolicies(ctx, state.ID.ValueString(), service)
	detachServerGroupFromAppSegments(ctx, state.ID.ValueString(), service)

	if _, err := servergroup.Delete(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete server group: %v", err))
		return
	}
}

func (r *ServerGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := req.ID
	service := r.client.Service

	// Check if the import ID is numeric (assume it's an ID)
	_, parseErr := strconv.ParseInt(importID, 10, 64)
	if parseErr == nil {
		// It's numeric, use it as ID
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(importID))...)
		return
	}

	// Not numeric, try to lookup by name
	group, _, err := servergroup.GetByName(ctx, service, importID)
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to find server group by name '%s': %v", importID, err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(group.ID))...)
}

func (r *ServerGroupResource) readServerGroup(ctx context.Context, service *zscaler.Service, id string, existingState ServerGroupResourceModel) (ServerGroupResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var state ServerGroupResourceModel

	group, _, err := servergroup.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			tflog.Warn(ctx, "Server group not found; removing from state", map[string]any{"id": id})
			state.ID = types.StringValue("")
			return state, diags
		}
		diags.AddError("Client Error", fmt.Sprintf("Unable to read server group: %v", err))
		return state, diags
	}

	state = flattenServerGroupResource(ctx, group, existingState)
	return state, diags
}

func expandServerGroup(ctx context.Context, plan ServerGroupResourceModel) (servergroup.ServerGroup, diag.Diagnostics) {
	var diags diag.Diagnostics

	connectorGroups, connectorDiags := expandServerGroupMembership(ctx, plan.AppConnectorGroup)
	diags.Append(connectorDiags...)

	servers, serverDiags := expandServerGroupServers(ctx, plan.Servers)
	diags.Append(serverDiags...)

	extranetDTO, extranetDiags := expandExtranetDTO(ctx, plan.ExtranetDTO)
	diags.Append(extranetDiags...)

	// Note: dynamic_discovery and servers validation is done in Create/Update functions, not here
	// to match SDKv2 behavior where Update only validates if fields changed
	dynamicDiscovery := helpers.BoolValue(plan.DynamicDiscovery, false)

	payload := servergroup.ServerGroup{
		ID:                 plan.ID.ValueString(),
		Name:               plan.Name.ValueString(),
		Description:        plan.Description.ValueString(),
		Enabled:            helpers.BoolValue(plan.Enabled, false),
		ConfigSpace:        plan.ConfigSpace.ValueString(),
		IpAnchored:         helpers.BoolValue(plan.IPAanchored, false),
		DynamicDiscovery:   dynamicDiscovery,
		MicroTenantID:      plan.MicroTenantID.ValueString(),
		AppConnectorGroups: connectorGroups,
		Servers:            servers,
		ExtranetEnabled:    helpers.BoolValue(plan.ExtranetEnabled, false),
		ExtranetDTO:        extranetDTO,
	}

	return payload, diags
}

func flattenServerGroupResource(ctx context.Context, group *servergroup.ServerGroup, existingState ServerGroupResourceModel) ServerGroupResourceModel {
	connectors := flattenServerGroupConnectorGroups(ctx, group.AppConnectorGroups)
	servers := flattenServerGroupServers(ctx, group.Servers)
	extranet := flattenExtranetDTO(ctx, group.ExtranetDTO)

	// Preserve null values for optional boolean fields if they were null in the plan/state
	ipAnchored := types.BoolValue(group.IpAnchored)
	if existingState.IPAanchored.IsNull() || existingState.IPAanchored.IsUnknown() {
		ipAnchored = existingState.IPAanchored
	}

	extranetEnabled := types.BoolValue(group.ExtranetEnabled)
	if existingState.ExtranetEnabled.IsNull() || existingState.ExtranetEnabled.IsUnknown() {
		extranetEnabled = existingState.ExtranetEnabled
	}

	return ServerGroupResourceModel{
		ID:                types.StringValue(group.ID),
		Name:              types.StringValue(group.Name),
		Description:       types.StringValue(group.Description),
		Enabled:           types.BoolValue(group.Enabled),
		ConfigSpace:       types.StringValue(group.ConfigSpace),
		IPAanchored:       ipAnchored,
		DynamicDiscovery:  types.BoolValue(group.DynamicDiscovery),
		MicroTenantID:     types.StringValue(group.MicroTenantID),
		AppConnectorGroup: connectors,
		Servers:           servers,
		ExtranetEnabled:   extranetEnabled,
		ExtranetDTO:       extranet,
	}
}

func expandServerGroupMembership(ctx context.Context, list types.List) ([]appconnectorgroup.AppConnectorGroup, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var memberships []serverGroupAppConnectorGroups
	diags.Append(list.ElementsAs(ctx, &memberships, false)...)
	if diags.HasError() {
		return nil, diags
	}

	var result []appconnectorgroup.AppConnectorGroup
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
			result = append(result, appconnectorgroup.AppConnectorGroup{ID: id})
		}
	}

	return result, diags
}

func expandServerGroupServers(ctx context.Context, list types.List) ([]appservercontroller.ApplicationServer, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var memberships []serverGroupServers
	diags.Append(list.ElementsAs(ctx, &memberships, false)...)
	if diags.HasError() {
		return nil, diags
	}

	var result []appservercontroller.ApplicationServer
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
			result = append(result, appservercontroller.ApplicationServer{ID: id})
		}
	}

	return result, diags
}

func expandExtranetDTO(ctx context.Context, set types.Set) (common.ExtranetDTO, diag.Diagnostics) {
	var diags diag.Diagnostics
	dto := common.ExtranetDTO{}
	if set.IsNull() || set.IsUnknown() {
		return dto, diags
	}

	var dtos []serverGroupExtranetDTO
	diags.Append(set.ElementsAs(ctx, &dtos, false)...)
	if diags.HasError() || len(dtos) == 0 {
		return dto, diags
	}

	first := dtos[0]
	dto.ZpnErID = first.ZPNERID.ValueString()

	locations, locDiags := expandExtranetLocations(ctx, first.LocationDTO)
	diags.Append(locDiags...)
	dto.LocationDTO = locations

	locationGroups, groupDiags := expandExtranetLocationGroups(ctx, first.LocationGroupDTO)
	diags.Append(groupDiags...)
	dto.LocationGroupDTO = locationGroups

	return dto, diags
}

func expandExtranetLocations(ctx context.Context, set types.Set) ([]common.LocationDTO, diag.Diagnostics) {
	var diags diag.Diagnostics
	if set.IsNull() || set.IsUnknown() {
		return nil, diags
	}

	var locations []serverGroupExtranetLocation
	diags.Append(set.ElementsAs(ctx, &locations, false)...)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]common.LocationDTO, 0, len(locations))
	for _, loc := range locations {
		result = append(result, common.LocationDTO{ID: loc.ID.ValueString()})
	}

	return result, diags
}

func expandExtranetLocationGroups(ctx context.Context, set types.Set) ([]common.LocationGroupDTO, diag.Diagnostics) {
	var diags diag.Diagnostics
	if set.IsNull() || set.IsUnknown() {
		return nil, diags
	}

	var locations []serverGroupExtranetLocation
	diags.Append(set.ElementsAs(ctx, &locations, false)...)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]common.LocationGroupDTO, 0, len(locations))
	for _, loc := range locations {
		result = append(result, common.LocationGroupDTO{ID: loc.ID.ValueString()})
	}

	return result, diags
}

func flattenServerGroupConnectorGroups(ctx context.Context, groups []appconnectorgroup.AppConnectorGroup) types.List {
	attrTypes := map[string]attr.Type{
		"id": types.SetType{ElemType: types.StringType},
	}

	if len(groups) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes})
	}

	ids := make([]string, len(groups))
	for i, group := range groups {
		ids[i] = group.ID
	}

	setValue, diags := types.SetValueFrom(ctx, types.StringType, ids)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes})
	}

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{"id": setValue})
	if objDiags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes})
	}

	result, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	if listDiags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes})
	}

	return result
}

func flattenServerGroupServers(ctx context.Context, servers []appservercontroller.ApplicationServer) types.List {
	attrTypes := map[string]attr.Type{
		"id": types.SetType{ElemType: types.StringType},
	}

	if len(servers) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes})
	}

	ids := make([]string, len(servers))
	for i, server := range servers {
		ids[i] = server.ID
	}

	setValue, diags := types.SetValueFrom(ctx, types.StringType, ids)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes})
	}

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{"id": setValue})
	if objDiags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes})
	}

	result, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	if listDiags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes})
	}

	return result
}

func flattenExtranetDTO(ctx context.Context, dto common.ExtranetDTO) types.Set {
	attrTypes := map[string]attr.Type{
		"zpn_er_id":          types.StringType,
		"location_dto":       types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"id": types.StringType}}},
		"location_group_dto": types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"id": types.StringType}}},
	}

	if dto.ZpnErID == "" && len(dto.LocationDTO) == 0 && len(dto.LocationGroupDTO) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes})
	}

	locationValues := make([]attr.Value, 0, len(dto.LocationDTO))
	for _, loc := range dto.LocationDTO {
		obj, diags := types.ObjectValue(map[string]attr.Type{"id": types.StringType}, map[string]attr.Value{"id": types.StringValue(loc.ID)})
		if diags.HasError() {
			return types.SetNull(types.ObjectType{AttrTypes: attrTypes})
		}
		locationValues = append(locationValues, obj)
	}
	locationSet, diags := types.SetValue(types.ObjectType{AttrTypes: map[string]attr.Type{"id": types.StringType}}, locationValues)
	if diags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes})
	}

	locationGroupValues := make([]attr.Value, 0, len(dto.LocationGroupDTO))
	for _, loc := range dto.LocationGroupDTO {
		obj, diags := types.ObjectValue(map[string]attr.Type{"id": types.StringType}, map[string]attr.Value{"id": types.StringValue(loc.ID)})
		if diags.HasError() {
			return types.SetNull(types.ObjectType{AttrTypes: attrTypes})
		}
		locationGroupValues = append(locationGroupValues, obj)
	}
	locationGroupSet, groupDiags := types.SetValue(types.ObjectType{AttrTypes: map[string]attr.Type{"id": types.StringType}}, locationGroupValues)
	if groupDiags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes})
	}

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"zpn_er_id":          types.StringValue(dto.ZpnErID),
		"location_dto":       locationSet,
		"location_group_dto": locationGroupSet,
	})
	if objDiags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes})
	}

	result, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	if setDiags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes})
	}

	return result
}

func detachServerGroupFromPolicies(ctx context.Context, id string, service *zscaler.Service) {
	helpers.PolicyRulesDetachLock.Lock()
	defer helpers.PolicyRulesDetachLock.Unlock()

	policySet, _, err := policysetcontroller.GetByPolicyType(ctx, service, "ACCESS_POLICY")
	if err != nil {
		tflog.Warn(ctx, "Failed to fetch access policy set", map[string]any{"error": err.Error()})
		return
	}

	rules, _, err := policysetcontroller.GetAllByType(ctx, service, "ACCESS_POLICY")
	if err != nil {
		tflog.Warn(ctx, "Failed to fetch access policy rules", map[string]any{"error": err.Error()})
		return
	}

	for _, rule := range rules {
		changed := false
		updated := make([]servergroup.ServerGroup, 0, len(rule.AppServerGroups))
		for _, sg := range rule.AppServerGroups {
			if sg.ID == id {
				changed = true
				continue
			}
			updated = append(updated, servergroup.ServerGroup{ID: sg.ID})
		}
		if !changed {
			continue
		}
		rule.AppServerGroups = updated
		if _, err := policysetcontroller.UpdateRule(ctx, service, policySet.ID, rule.ID, &rule); err != nil {
			tflog.Warn(ctx, "Failed to update access policy rule", map[string]any{"id": rule.ID, "error": err.Error()})
		}
	}
}

func detachServerGroupFromAppSegments(ctx context.Context, id string, service *zscaler.Service) {
	apps, _, err := applicationsegment.GetAll(ctx, service)
	if err != nil {
		tflog.Warn(ctx, "Failed to fetch application segments", map[string]any{"error": err.Error()})
		return
	}

	for _, app := range apps {
		changed := false
		updated := make([]servergroup.ServerGroup, 0, len(app.ServerGroups))
		for _, sg := range app.ServerGroups {
			if sg.ID == id {
				changed = true
				continue
			}
			updated = append(updated, servergroup.ServerGroup{ID: sg.ID})
		}
		if !changed {
			continue
		}
		app.ServerGroups = updated
		if _, err := applicationsegment.Update(ctx, service, app.ID, app); err != nil {
			tflog.Warn(ctx, "Failed to update application segment", map[string]any{"id": app.ID, "error": err.Error()})
		}
	}
}

func detachServerGroupFromAppConnectorGroups(ctx context.Context, id string, service *zscaler.Service) error {
	serverGroupDetachLock.Lock()
	defer serverGroupDetachLock.Unlock()

	group, _, err := servergroup.Get(ctx, service, id)
	if err != nil {
		return fmt.Errorf("failed to fetch server group: %w", err)
	}

	for _, connector := range group.AppConnectorGroups {
		appConnector, _, err := appconnectorgroup.Get(ctx, service, connector.ID)
		if err != nil {
			tflog.Warn(ctx, "Failed to fetch app connector group", map[string]any{"id": connector.ID, "error": err.Error()})
			continue
		}
		updated := make([]appconnectorgroup.AppServerGroup, 0, len(appConnector.AppServerGroup))
		for _, sg := range appConnector.AppServerGroup {
			if sg.ID == id {
				continue
			}
			updated = append(updated, sg)
		}
		appConnector.AppServerGroup = updated
		if _, err := appconnectorgroup.Update(ctx, service, appConnector.ID, appConnector); err != nil {
			tflog.Warn(ctx, "Failed to update app connector group", map[string]any{"id": appConnector.ID, "error": err.Error()})
		}
	}

	return nil
}
