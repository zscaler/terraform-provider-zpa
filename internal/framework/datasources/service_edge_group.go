package datasources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgecontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgegroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/trustednetwork"
)

var (
	_ datasource.DataSource              = &ServiceEdgeGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &ServiceEdgeGroupDataSource{}
)

func NewServiceEdgeGroupDataSource() datasource.DataSource {
	return &ServiceEdgeGroupDataSource{}
}

type ServiceEdgeGroupDataSource struct {
	client *client.Client
}

type ServiceEdgeGroupModel struct {
	ID                            types.String `tfsdk:"id"`
	Name                          types.String `tfsdk:"name"`
	AltCloud                      types.String `tfsdk:"alt_cloud"`
	CityCountry                   types.String `tfsdk:"city_country"`
	CountryCode                   types.String `tfsdk:"country_code"`
	CreationTime                  types.String `tfsdk:"creation_time"`
	Description                   types.String `tfsdk:"description"`
	Enabled                       types.Bool   `tfsdk:"enabled"`
	GeoLocationID                 types.String `tfsdk:"geo_location_id"`
	IsPublic                      types.String `tfsdk:"is_public"`
	Latitude                      types.String `tfsdk:"latitude"`
	Location                      types.String `tfsdk:"location"`
	Longitude                     types.String `tfsdk:"longitude"`
	OverrideVersionProfile        types.Bool   `tfsdk:"override_version_profile"`
	ModifiedBy                    types.String `tfsdk:"modified_by"`
	ModifiedTime                  types.String `tfsdk:"modified_time"`
	UseInDrMode                   types.Bool   `tfsdk:"use_in_dr_mode"`
	SiteID                        types.String `tfsdk:"site_id"`
	SiteName                      types.String `tfsdk:"site_name"`
	UpgradeDay                    types.String `tfsdk:"upgrade_day"`
	UpgradeTimeInSecs             types.String `tfsdk:"upgrade_time_in_secs"`
	VersionProfileID              types.String `tfsdk:"version_profile_id"`
	VersionProfileName            types.String `tfsdk:"version_profile_name"`
	VersionProfileVisibilityScope types.String `tfsdk:"version_profile_visibility_scope"`
	GraceDistanceEnabled          types.Bool   `tfsdk:"grace_distance_enabled"`
	GraceDistanceValue            types.String `tfsdk:"grace_distance_value"`
	GraceDistanceValueUnit        types.String `tfsdk:"grace_distance_value_unit"`
	ServiceEdges                  types.Set    `tfsdk:"service_edges"`
	TrustedNetworks               types.Set    `tfsdk:"trusted_networks"`
	MicroTenantID                 types.String `tfsdk:"microtenant_id"`
	MicroTenantName               types.String `tfsdk:"microtenant_name"`
}

func (d *ServiceEdgeGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_edge_group"
}

func (d *ServiceEdgeGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	trustedNetworkAttributes := map[string]schema.Attribute{
		"id":   schema.StringAttribute{Computed: true},
		"name": schema.StringAttribute{Computed: true},
	}

	serviceEdgeAttributes := map[string]schema.Attribute{
		"id":                          schema.StringAttribute{Computed: true},
		"name":                        schema.StringAttribute{Computed: true},
		"description":                 schema.StringAttribute{Computed: true},
		"enabled":                     schema.BoolAttribute{Computed: true},
		"current_version":             schema.StringAttribute{Computed: true},
		"previous_version":            schema.StringAttribute{Computed: true},
		"platform":                    schema.StringAttribute{Computed: true},
		"private_ip":                  schema.StringAttribute{Computed: true},
		"public_ip":                   schema.StringAttribute{Computed: true},
		"provisioning_key_id":         schema.StringAttribute{Computed: true},
		"provisioning_key_name":       schema.StringAttribute{Computed: true},
		"control_channel_status":      schema.StringAttribute{Computed: true},
		"last_broker_connect_time":    schema.StringAttribute{Computed: true},
		"last_broker_disconnect_time": schema.StringAttribute{Computed: true},
		"listen_ips": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
		},
		"publish_ips": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
		},
		"publish_ipv6": schema.BoolAttribute{Computed: true},
		"enrollment_cert": schema.MapAttribute{
			ElementType: types.StringType,
			Computed:    true,
		},
	}

	privateBrokerAttributes := map[string]schema.Attribute{
		"id":                   schema.StringAttribute{Computed: true},
		"current_version":      schema.StringAttribute{Computed: true},
		"upgrade_status":       schema.StringAttribute{Computed: true},
		"upgrade_attempt":      schema.StringAttribute{Computed: true},
		"last_connect_time":    schema.StringAttribute{Computed: true},
		"last_disconnect_time": schema.StringAttribute{Computed: true},
	}

	subModuleAttributes := map[string]schema.Attribute{
		"id":               schema.StringAttribute{Computed: true},
		"current_version":  schema.StringAttribute{Computed: true},
		"expected_version": schema.StringAttribute{Computed: true},
		"role":             schema.StringAttribute{Computed: true},
		"upgrade_status":   schema.StringAttribute{Computed: true},
	}

	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA Service Edge Group by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id":                               schema.StringAttribute{Optional: true},
			"name":                             schema.StringAttribute{Optional: true},
			"alt_cloud":                        schema.StringAttribute{Computed: true},
			"city_country":                     schema.StringAttribute{Computed: true},
			"country_code":                     schema.StringAttribute{Computed: true},
			"creation_time":                    schema.StringAttribute{Computed: true},
			"description":                      schema.StringAttribute{Computed: true},
			"enabled":                          schema.BoolAttribute{Computed: true},
			"geo_location_id":                  schema.StringAttribute{Computed: true},
			"is_public":                        schema.StringAttribute{Computed: true},
			"latitude":                         schema.StringAttribute{Computed: true},
			"location":                         schema.StringAttribute{Computed: true},
			"longitude":                        schema.StringAttribute{Computed: true},
			"override_version_profile":         schema.BoolAttribute{Computed: true},
			"modified_by":                      schema.StringAttribute{Computed: true},
			"modified_time":                    schema.StringAttribute{Computed: true},
			"use_in_dr_mode":                   schema.BoolAttribute{Computed: true},
			"site_id":                          schema.StringAttribute{Computed: true},
			"site_name":                        schema.StringAttribute{Computed: true},
			"upgrade_day":                      schema.StringAttribute{Computed: true},
			"upgrade_time_in_secs":             schema.StringAttribute{Computed: true},
			"version_profile_id":               schema.StringAttribute{Computed: true},
			"version_profile_name":             schema.StringAttribute{Computed: true},
			"version_profile_visibility_scope": schema.StringAttribute{Computed: true},
			"grace_distance_enabled":           schema.BoolAttribute{Computed: true},
			"grace_distance_value":             schema.StringAttribute{Computed: true},
			"grace_distance_value_unit":        schema.StringAttribute{Computed: true},
			"microtenant_id":                   schema.StringAttribute{Optional: true},
			"microtenant_name":                 schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			"service_edges": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: serviceEdgeAttributes,
					Blocks: map[string]schema.Block{
						"private_broker_version": schema.SingleNestedBlock{
							Attributes: privateBrokerAttributes,
							Blocks: map[string]schema.Block{
								"sub_module_upgrades": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: subModuleAttributes,
									},
								},
							},
						},
					},
				},
			},
			"trusted_networks": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: trustedNetworkAttributes,
				},
			},
		},
	}
}

func (d *ServiceEdgeGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceEdgeGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ServiceEdgeGroupModel
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
		group *serviceedgegroup.ServiceEdgeGroup
		err   error
	)

	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		id := data.ID.ValueString()
		tflog.Info(ctx, "Fetching service edge group", map[string]any{"id": id})
		group, _, err = serviceedgegroup.Get(ctx, service, id)
	} else {
		name := data.Name.ValueString()
		tflog.Info(ctx, "Fetching service edge group", map[string]any{"name": name})
		group, _, err = serviceedgegroup.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service edge group: %v", err))
		return
	}

	flattened, diags := flattenServiceEdgeGroup(ctx, group)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &flattened)...)
}

func flattenServiceEdgeGroup(ctx context.Context, group *serviceedgegroup.ServiceEdgeGroup) (ServiceEdgeGroupModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	serviceEdges, seDiags := flattenServiceEdges(ctx, group.ServiceEdges)
	diags.Append(seDiags...)

	trustedNetworks, tnDiags := flattenTrustedNetworks(ctx, group.TrustedNetworks)
	diags.Append(tnDiags...)

	graceValue := group.GraceDistanceValue
	if graceValue != "" {
		if parsed, err := strconv.ParseFloat(graceValue, 64); err == nil {
			graceValue = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", parsed), "0"), ".")
		}
	}

	model := ServiceEdgeGroupModel{
		ID:                            types.StringValue(group.ID),
		Name:                          types.StringValue(group.Name),
		AltCloud:                      types.StringValue(group.AltCloud),
		CityCountry:                   types.StringValue(group.CityCountry),
		CountryCode:                   types.StringValue(group.CountryCode),
		CreationTime:                  types.StringValue(group.CreationTime),
		Description:                   types.StringValue(group.Description),
		Enabled:                       types.BoolValue(group.Enabled),
		GeoLocationID:                 types.StringValue(group.GeoLocationID),
		IsPublic:                      types.StringValue(group.IsPublic),
		Latitude:                      types.StringValue(group.Latitude),
		Location:                      types.StringValue(group.Location),
		Longitude:                     types.StringValue(group.Longitude),
		OverrideVersionProfile:        types.BoolValue(group.OverrideVersionProfile),
		ModifiedBy:                    types.StringValue(group.ModifiedBy),
		ModifiedTime:                  types.StringValue(group.ModifiedTime),
		UseInDrMode:                   types.BoolValue(group.UseInDrMode),
		SiteID:                        types.StringValue(group.SiteID),
		SiteName:                      types.StringValue(group.SiteName),
		UpgradeDay:                    types.StringValue(group.UpgradeDay),
		UpgradeTimeInSecs:             types.StringValue(group.UpgradeTimeInSecs),
		VersionProfileID:              types.StringValue(group.VersionProfileID),
		VersionProfileName:            types.StringValue(group.VersionProfileName),
		VersionProfileVisibilityScope: types.StringValue(group.VersionProfileVisibilityScope),
		GraceDistanceEnabled:          types.BoolValue(group.GraceDistanceEnabled),
		GraceDistanceValue:            types.StringValue(graceValue),
		GraceDistanceValueUnit:        types.StringValue(group.GraceDistanceValueUnit),
		ServiceEdges:                  serviceEdges,
		TrustedNetworks:               trustedNetworks,
		MicroTenantID:                 types.StringValue(group.MicroTenantID),
		MicroTenantName:               types.StringValue(group.MicroTenantName),
	}

	return model, diags
}

func flattenServiceEdges(ctx context.Context, edges []serviceedgecontroller.ServiceEdgeController) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := map[string]attr.Type{
		"id":                          types.StringType,
		"name":                        types.StringType,
		"description":                 types.StringType,
		"enabled":                     types.BoolType,
		"current_version":             types.StringType,
		"previous_version":            types.StringType,
		"platform":                    types.StringType,
		"private_ip":                  types.StringType,
		"public_ip":                   types.StringType,
		"provisioning_key_id":         types.StringType,
		"provisioning_key_name":       types.StringType,
		"control_channel_status":      types.StringType,
		"last_broker_connect_time":    types.StringType,
		"last_broker_disconnect_time": types.StringType,
		"listen_ips":                  types.SetType{ElemType: types.StringType},
		"publish_ips":                 types.SetType{ElemType: types.StringType},
		"publish_ipv6":                types.BoolType,
		"enrollment_cert":             types.MapType{ElemType: types.StringType},
		"private_broker_version": types.ObjectType{AttrTypes: map[string]attr.Type{
			"id":                   types.StringType,
			"current_version":      types.StringType,
			"upgrade_status":       types.StringType,
			"upgrade_attempt":      types.StringType,
			"last_connect_time":    types.StringType,
			"last_disconnect_time": types.StringType,
			"sub_module_upgrades": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"id":               types.StringType,
				"current_version":  types.StringType,
				"expected_version": types.StringType,
				"role":             types.StringType,
				"upgrade_status":   types.StringType,
			}}},
		}},
	}

	if len(edges) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	values := make([]attr.Value, 0, len(edges))

	for _, edge := range edges {
		listenIPs, listenDiags := types.SetValueFrom(ctx, types.StringType, edge.ListenIPs)
		diags.Append(listenDiags...)

		publishIPs, publishDiags := types.SetValueFrom(ctx, types.StringType, edge.PublishIPs)
		diags.Append(publishDiags...)

		enrollment, enrollmentDiags := mapStringInterfaceToMap(edge.EnrollmentCert)
		diags.Append(enrollmentDiags...)

		brokerValue, brokerDiags := helpers.FlattenPrivateBrokerVersionToObject(ctx, edge.PrivateBrokerVersion)
		diags.Append(brokerDiags...)

		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":                          types.StringValue(edge.ID),
			"name":                        types.StringValue(edge.Name),
			"description":                 types.StringValue(edge.Description),
			"enabled":                     types.BoolValue(edge.Enabled),
			"current_version":             types.StringValue(edge.CurrentVersion),
			"previous_version":            types.StringValue(edge.PreviousVersion),
			"platform":                    types.StringValue(edge.Platform),
			"private_ip":                  types.StringValue(edge.PrivateIP),
			"public_ip":                   types.StringValue(edge.PublicIP),
			"provisioning_key_id":         types.StringValue(edge.ProvisioningKeyID),
			"provisioning_key_name":       types.StringValue(edge.ProvisioningKeyName),
			"control_channel_status":      types.StringValue(edge.ControlChannelStatus),
			"last_broker_connect_time":    types.StringValue(edge.LastBrokerConnectTime),
			"last_broker_disconnect_time": types.StringValue(edge.LastBrokerDisconnectTime),
			"listen_ips":                  listenIPs,
			"publish_ips":                 publishIPs,
			"publish_ipv6":                types.BoolValue(edge.PublishIPv6),
			"enrollment_cert":             enrollment,
			"private_broker_version":      brokerValue,
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	set, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(setDiags...)
	return set, diags
}

func flattenTrustedNetworks(ctx context.Context, networks []trustednetwork.TrustedNetwork) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}

	if len(networks) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	values := make([]attr.Value, 0, len(networks))
	for _, network := range networks {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":   types.StringValue(network.ID),
			"name": types.StringValue(network.Name),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	set, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(setDiags...)
	return set, diags
}
