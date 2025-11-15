package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
)

var (
	_ datasource.DataSource              = &AppConnectorGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &AppConnectorGroupDataSource{}
)

func NewAppConnectorGroupDataSource() datasource.DataSource {
	return &AppConnectorGroupDataSource{}
}

type AppConnectorGroupDataSource struct {
	client *client.Client
}

type AppConnectorGroupDataSourceModel struct {
	ID                            types.String `tfsdk:"id"`
	Name                          types.String `tfsdk:"name"`
	CityCountry                   types.String `tfsdk:"city_country"`
	CountryCode                   types.String `tfsdk:"country_code"`
	CreationTime                  types.String `tfsdk:"creation_time"`
	Description                   types.String `tfsdk:"description"`
	DNSQueryType                  types.String `tfsdk:"dns_query_type"`
	Enabled                       types.Bool   `tfsdk:"enabled"`
	GeoLocationID                 types.String `tfsdk:"geo_location_id"`
	Latitude                      types.String `tfsdk:"latitude"`
	Location                      types.String `tfsdk:"location"`
	Longitude                     types.String `tfsdk:"longitude"`
	ModifiedBy                    types.String `tfsdk:"modifiedby"`
	ModifiedTime                  types.String `tfsdk:"modified_time"`
	OverrideVersionProfile        types.Bool   `tfsdk:"override_version_profile"`
	PRAEnabled                    types.Bool   `tfsdk:"pra_enabled"`
	WAFDisabled                   types.Bool   `tfsdk:"waf_disabled"`
	UpgradeDay                    types.String `tfsdk:"upgrade_day"`
	UpgradeTimeInSecs             types.String `tfsdk:"upgrade_time_in_secs"`
	VersionProfileID              types.String `tfsdk:"version_profile_id"`
	VersionProfileName            types.String `tfsdk:"version_profile_name"`
	VersionProfileVisibilityScope types.String `tfsdk:"version_profile_visibility_scope"`
	LSSAppConnectorGroup          types.Bool   `tfsdk:"lss_app_connector_group"`
	TCPQuickAckApp                types.Bool   `tfsdk:"tcp_quick_ack_app"`
	TCPQuickAckAssistant          types.Bool   `tfsdk:"tcp_quick_ack_assistant"`
	TCPQuickAckReadAssistant      types.Bool   `tfsdk:"tcp_quick_ack_read_assistant"`
	UseInDrMode                   types.Bool   `tfsdk:"use_in_dr_mode"`
	MicroTenantID                 types.String `tfsdk:"microtenant_id"`
	MicroTenantName               types.String `tfsdk:"microtenant_name"`
	ServerGroups                  types.List   `tfsdk:"server_groups"`
}

func (d *AppConnectorGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_connector_group"
}

func (d *AppConnectorGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	serverGroupAttrTypes := map[string]schema.Attribute{
		"config_space":      schema.StringAttribute{Computed: true},
		"creation_time":     schema.StringAttribute{Computed: true},
		"description":       schema.StringAttribute{Computed: true},
		"enabled":           schema.BoolAttribute{Computed: true},
		"id":                schema.StringAttribute{Computed: true},
		"dynamic_discovery": schema.BoolAttribute{Computed: true},
		"modifiedby":        schema.StringAttribute{Computed: true},
		"modified_time":     schema.StringAttribute{Computed: true},
		"name":              schema.StringAttribute{Computed: true},
	}

	resp.Schema = schema.Schema{
		Description: "Data source for retrieving a ZPA App Connector Group",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "App Connector Group ID",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "App Connector Group name",
				Optional:    true,
			},
			"city_country":                     schema.StringAttribute{Computed: true},
			"country_code":                     schema.StringAttribute{Computed: true},
			"creation_time":                    schema.StringAttribute{Computed: true},
			"description":                      schema.StringAttribute{Computed: true},
			"dns_query_type":                   schema.StringAttribute{Computed: true},
			"enabled":                          schema.BoolAttribute{Computed: true},
			"geo_location_id":                  schema.StringAttribute{Computed: true},
			"latitude":                         schema.StringAttribute{Computed: true},
			"location":                         schema.StringAttribute{Computed: true},
			"longitude":                        schema.StringAttribute{Computed: true},
			"modifiedby":                       schema.StringAttribute{Computed: true},
			"modified_time":                    schema.StringAttribute{Computed: true},
			"override_version_profile":         schema.BoolAttribute{Computed: true},
			"pra_enabled":                      schema.BoolAttribute{Computed: true},
			"waf_disabled":                     schema.BoolAttribute{Computed: true},
			"upgrade_day":                      schema.StringAttribute{Computed: true},
			"upgrade_time_in_secs":             schema.StringAttribute{Computed: true},
			"version_profile_id":               schema.StringAttribute{Computed: true},
			"version_profile_name":             schema.StringAttribute{Computed: true},
			"version_profile_visibility_scope": schema.StringAttribute{Computed: true},
			"lss_app_connector_group":          schema.BoolAttribute{Computed: true},
			"tcp_quick_ack_app":                schema.BoolAttribute{Computed: true},
			"tcp_quick_ack_assistant":          schema.BoolAttribute{Computed: true},
			"tcp_quick_ack_read_assistant":     schema.BoolAttribute{Computed: true},
			"use_in_dr_mode":                   schema.BoolAttribute{Computed: true},
			"microtenant_id": schema.StringAttribute{
				Description: "Microtenant ID to scope the lookup",
				Optional:    true,
				Computed:    true,
			},
			"microtenant_name": schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			"server_groups": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: serverGroupAttrTypes,
				},
			},
		},
	}
}

func (d *AppConnectorGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *AppConnectorGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AppConnectorGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (data.ID.IsNull() || data.ID.ValueString() == "") && (data.Name.IsNull() || data.Name.ValueString() == "") {
		resp.Diagnostics.AddError("Missing Identifier", "Either id or name must be specified")
		return
	}

	service := d.client.Service
	if !data.MicroTenantID.IsNull() && data.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(data.MicroTenantID.ValueString())
	}

	var (
		group *appconnectorgroup.AppConnectorGroup
		err   error
	)

	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		id := data.ID.ValueString()
		tflog.Info(ctx, "Fetching app connector group", map[string]any{"id": id})
		group, _, err = appconnectorgroup.Get(ctx, service, id)
	} else {
		name := data.Name.ValueString()
		tflog.Info(ctx, "Fetching app connector group", map[string]any{"name": name})
		group, _, err = appconnectorgroup.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read app connector group: %v", err))
		return
	}

	if diags := populateAppConnectorGroupData(ctx, service, group, &data); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func populateAppConnectorGroupData(ctx context.Context, service *zscaler.Service, group *appconnectorgroup.AppConnectorGroup, data *AppConnectorGroupDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(group.ID)
	data.Name = types.StringValue(group.Name)
	data.CityCountry = types.StringValue(group.CityCountry)
	data.CountryCode = types.StringValue(group.CountryCode)
	data.CreationTime = types.StringValue(group.CreationTime)
	data.Description = types.StringValue(group.Description)
	data.DNSQueryType = types.StringValue(group.DNSQueryType)
	data.Enabled = types.BoolValue(group.Enabled)
	data.GeoLocationID = types.StringValue(group.GeoLocationID)
	data.Latitude = types.StringValue(group.Latitude)
	data.Location = types.StringValue(group.Location)
	data.Longitude = types.StringValue(group.Longitude)
	data.ModifiedBy = types.StringValue(group.ModifiedBy)
	data.ModifiedTime = types.StringValue(group.ModifiedTime)
	data.OverrideVersionProfile = types.BoolValue(group.OverrideVersionProfile)
	data.PRAEnabled = types.BoolValue(group.PRAEnabled)
	data.WAFDisabled = types.BoolValue(group.WAFDisabled)
	data.UpgradeDay = types.StringValue(group.UpgradeDay)
	data.UpgradeTimeInSecs = types.StringValue(group.UpgradeTimeInSecs)
	data.VersionProfileID = types.StringValue(group.VersionProfileID)
	data.VersionProfileName = types.StringValue(group.VersionProfileName)
	data.VersionProfileVisibilityScope = types.StringValue(group.VersionProfileVisibilityScope)
	data.LSSAppConnectorGroup = types.BoolValue(group.LSSAppConnectorGroup)
	data.TCPQuickAckApp = types.BoolValue(group.TCPQuickAckApp)
	data.TCPQuickAckAssistant = types.BoolValue(group.TCPQuickAckAssistant)
	data.TCPQuickAckReadAssistant = types.BoolValue(group.TCPQuickAckReadAssistant)
	data.UseInDrMode = types.BoolValue(group.UseInDrMode)
	data.MicroTenantID = types.StringValue(group.MicroTenantID)
	data.MicroTenantName = types.StringValue(group.MicroTenantName)

	serverGroups, serverGroupDiags := flattenServerGroupList(group.AppServerGroup)
	diags.Append(serverGroupDiags...)
	data.ServerGroups = serverGroups

	return diags
}

func flattenServerGroupList(groups []appconnectorgroup.AppServerGroup) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := map[string]attr.Type{
		"config_space":      types.StringType,
		"creation_time":     types.StringType,
		"description":       types.StringType,
		"enabled":           types.BoolType,
		"id":                types.StringType,
		"dynamic_discovery": types.BoolType,
		"modifiedby":        types.StringType,
		"modified_time":     types.StringType,
		"name":              types.StringType,
	}

	elements := make([]attr.Value, 0, len(groups))

	for _, group := range groups {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"config_space":      types.StringValue(group.ConfigSpace),
			"creation_time":     types.StringValue(group.CreationTime),
			"description":       types.StringValue(group.Description),
			"enabled":           types.BoolValue(group.Enabled),
			"id":                types.StringValue(group.ID),
			"dynamic_discovery": types.BoolValue(group.DynamicDiscovery),
			"modifiedby":        types.StringValue(group.ModifiedBy),
			"modified_time":     types.StringValue(group.ModifiedTime),
			"name":              types.StringValue(group.Name),
		})
		diags.Append(objDiags...)
		elements = append(elements, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, elements)
	diags.Append(listDiags...)
	return list, diags
}

func mapStringInterfaceToMap(value map[string]interface{}) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(value) == 0 {
		return types.MapNull(types.StringType), diags
	}

	elements := make(map[string]attr.Value, len(value))
	for k, v := range value {
		elements[k] = types.StringValue(fmt.Sprint(v))
	}

	mapValue, mapDiags := types.MapValue(types.StringType, elements)
	diags.Append(mapDiags...)
	return mapValue, diags
}
