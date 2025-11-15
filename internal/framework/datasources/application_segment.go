package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
)

var (
	_ datasource.DataSource              = &ApplicationSegmentDataSource{}
	_ datasource.DataSourceWithConfigure = &ApplicationSegmentDataSource{}
)

func NewApplicationSegmentDataSource() datasource.DataSource {
	return &ApplicationSegmentDataSource{}
}

type ApplicationSegmentDataSource struct {
	client *client.Client
}

type ApplicationSegmentModel struct {
	ID                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	SegmentGroupID            types.String `tfsdk:"segment_group_id"`
	SegmentGroupName          types.String `tfsdk:"segment_group_name"`
	BypassType                types.String `tfsdk:"bypass_type"`
	ConfigSpace               types.String `tfsdk:"config_space"`
	CreationTime              types.String `tfsdk:"creation_time"`
	DefaultIdleTimeout        types.String `tfsdk:"default_idle_timeout"`
	DefaultMaxAge             types.String `tfsdk:"default_max_age"`
	Description               types.String `tfsdk:"description"`
	DomainNames               types.List   `tfsdk:"domain_names"`
	DoubleEncrypt             types.Bool   `tfsdk:"double_encrypt"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
	HealthCheckType           types.String `tfsdk:"health_check_type"`
	HealthReporting           types.String `tfsdk:"health_reporting"`
	SelectConnectorCloseToApp types.Bool   `tfsdk:"select_connector_close_to_app"`
	UseInDrMode               types.Bool   `tfsdk:"use_in_dr_mode"`
	IsIncompleteDRConfig      types.Bool   `tfsdk:"is_incomplete_dr_config"`
	IPAnchored                types.Bool   `tfsdk:"ip_anchored"`
	IsCnameEnabled            types.Bool   `tfsdk:"is_cname_enabled"`
	ModifiedBy                types.String `tfsdk:"modified_by"`
	ModifiedTime              types.String `tfsdk:"modified_time"`
	PassiveHealthEnabled      types.Bool   `tfsdk:"passive_health_enabled"`
	APIProtectionEnabled      types.Bool   `tfsdk:"api_protection_enabled"`
	ServerGroups              types.Set    `tfsdk:"server_groups"`
	MicroTenantID             types.String `tfsdk:"microtenant_id"`
	MicroTenantName           types.String `tfsdk:"microtenant_name"`
	TCPPortRanges             types.List   `tfsdk:"tcp_port_ranges"`
	UDPPortRanges             types.List   `tfsdk:"udp_port_ranges"`
	TCPPortRange              types.Set    `tfsdk:"tcp_port_range"`
	UDPPortRange              types.Set    `tfsdk:"udp_port_range"`
	MatchStyle                types.String `tfsdk:"match_style"`
}

func (d *ApplicationSegmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_segment"
}

func (d *ApplicationSegmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA application segment by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the application segment.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the application segment.",
			},
			"segment_group_id":     schema.StringAttribute{Computed: true},
			"segment_group_name":   schema.StringAttribute{Computed: true},
			"bypass_type":          schema.StringAttribute{Computed: true},
			"config_space":         schema.StringAttribute{Computed: true},
			"creation_time":        schema.StringAttribute{Computed: true},
			"default_idle_timeout": schema.StringAttribute{Computed: true},
			"default_max_age":      schema.StringAttribute{Computed: true},
			"description":          schema.StringAttribute{Computed: true},
			"domain_names": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"double_encrypt":                schema.BoolAttribute{Computed: true},
			"enabled":                       schema.BoolAttribute{Computed: true},
			"health_check_type":             schema.StringAttribute{Computed: true},
			"health_reporting":              schema.StringAttribute{Computed: true},
			"select_connector_close_to_app": schema.BoolAttribute{Computed: true},
			"use_in_dr_mode":                schema.BoolAttribute{Computed: true},
			"is_incomplete_dr_config":       schema.BoolAttribute{Computed: true},
			"ip_anchored":                   schema.BoolAttribute{Computed: true},
			"is_cname_enabled":              schema.BoolAttribute{Computed: true},
			"modified_by":                   schema.StringAttribute{Computed: true},
			"modified_time":                 schema.StringAttribute{Computed: true},
			"passive_health_enabled":        schema.BoolAttribute{Computed: true},
			"api_protection_enabled":        schema.BoolAttribute{Computed: true},
			"match_style":                   schema.StringAttribute{Computed: true},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID used to scope the lookup.",
			},
			"microtenant_name": schema.StringAttribute{Computed: true},
			"tcp_port_ranges": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"udp_port_ranges": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"tcp_port_range": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"from": schema.StringAttribute{Computed: true},
						"to":   schema.StringAttribute{Computed: true},
					},
				},
			},
			"udp_port_range": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"from": schema.StringAttribute{Computed: true},
						"to":   schema.StringAttribute{Computed: true},
					},
				},
			},
			"server_groups": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *ApplicationSegmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cl, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = cl
}

func (d *ApplicationSegmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data ApplicationSegmentModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := d.client.Service
	if !data.MicroTenantID.IsNull() && !data.MicroTenantID.IsUnknown() {
		microID := strings.TrimSpace(data.MicroTenantID.ValueString())
		if microID != "" {
			service = service.WithMicroTenant(microID)
			data.MicroTenantID = types.StringValue(microID)
		}
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())
	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided to read an application segment.")
		return
	}

	var (
		segment *applicationsegment.ApplicationSegmentResource
		err     error
	)

	if id != "" {
		tflog.Debug(ctx, "Retrieving application segment by ID", map[string]any{"id": id})
		segment, _, err = applicationsegment.Get(ctx, service, id)
	} else {
		tflog.Debug(ctx, "Retrieving application segment by name", map[string]any{"name": name})
		segment, _, err = applicationsegment.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read application segment: %v", err))
		return
	}

	if segment == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Application segment with id %q or name %q was not found.", id, name))
		return
	}

	domainNames, domainDiags := helpers.StringSliceToList(ctx, segment.DomainNames)
	resp.Diagnostics.Append(domainDiags...)

	tcpRanges, tcpRangesDiags := helpers.StringSliceToList(ctx, segment.TCPPortRanges)
	resp.Diagnostics.Append(tcpRangesDiags...)

	udpRanges, udpRangesDiags := helpers.StringSliceToList(ctx, segment.UDPPortRanges)
	resp.Diagnostics.Append(udpRangesDiags...)

	tcpPortsList, tcpDiags := helpers.FlattenNetworkPorts(ctx, segment.TCPAppPortRange)
	resp.Diagnostics.Append(tcpDiags...)

	udpPortsList, udpDiags := helpers.FlattenNetworkPorts(ctx, segment.UDPAppPortRange)
	resp.Diagnostics.Append(udpDiags...)

	// Convert List to Set for tcp_port_range and udp_port_range to match SDKv2 (TypeSet)
	attrTypes := map[string]attr.Type{
		"from": types.StringType,
		"to":   types.StringType,
	}
	portRangeObjectType := types.ObjectType{AttrTypes: attrTypes}

	var tcpPorts types.Set
	if !tcpPortsList.IsNull() && !tcpPortsList.IsUnknown() {
		var portValues []attr.Value
		elemDiags := tcpPortsList.ElementsAs(ctx, &portValues, false)
		resp.Diagnostics.Append(elemDiags...)
		if !resp.Diagnostics.HasError() {
			set, setDiags := types.SetValue(portRangeObjectType, portValues)
			resp.Diagnostics.Append(setDiags...)
			tcpPorts = set
		} else {
			tcpPorts = types.SetNull(portRangeObjectType)
		}
	} else {
		tcpPorts = types.SetNull(portRangeObjectType)
	}

	var udpPorts types.Set
	if !udpPortsList.IsNull() && !udpPortsList.IsUnknown() {
		var portValues []attr.Value
		elemDiags := udpPortsList.ElementsAs(ctx, &portValues, false)
		resp.Diagnostics.Append(elemDiags...)
		if !resp.Diagnostics.HasError() {
			set, setDiags := types.SetValue(portRangeObjectType, portValues)
			resp.Diagnostics.Append(setDiags...)
			udpPorts = set
		} else {
			udpPorts = types.SetNull(portRangeObjectType)
		}
	} else {
		udpPorts = types.SetNull(portRangeObjectType)
	}

	serverGroupsList, sgDiags := helpers.FlattenServerGroups(ctx, segment.ServerGroups)
	resp.Diagnostics.Append(sgDiags...)

	// Convert List to Set for Protocol 5 compliance
	// Also convert id from Set to List to match SDKv2 (TypeList)
	serverGroupAttrTypes := map[string]attr.Type{
		"id": types.ListType{ElemType: types.StringType},
	}
	serverGroupObjectType := types.ObjectType{AttrTypes: serverGroupAttrTypes}

	var serverGroups types.Set
	if !serverGroupsList.IsNull() && !serverGroupsList.IsUnknown() {
		type serverGroupModel struct {
			IDs types.Set `tfsdk:"id"` // FlattenServerGroups returns id as Set
		}
		var models []serverGroupModel
		elemDiags := serverGroupsList.ElementsAs(ctx, &models, false)
		resp.Diagnostics.Append(elemDiags...)
		if !resp.Diagnostics.HasError() && len(models) > 0 {
			objectValues := make([]types.Object, 0, len(models))
			for _, model := range models {
				// Convert Set to List for id
				var idList types.List
				if !model.IDs.IsNull() && !model.IDs.IsUnknown() {
					// Extract string values from Set
					var idStrings []string
					elemDiags := model.IDs.ElementsAs(ctx, &idStrings, false)
					resp.Diagnostics.Append(elemDiags...)
					if !resp.Diagnostics.HasError() {
						// Create List from string slice
						listElements := make([]attr.Value, len(idStrings))
						for i, idStr := range idStrings {
							listElements[i] = types.StringValue(idStr)
						}
						idListValue, listDiags := types.ListValue(types.StringType, listElements)
						resp.Diagnostics.Append(listDiags...)
						if !resp.Diagnostics.HasError() {
							idList = idListValue
						} else {
							idList = types.ListNull(types.StringType)
						}
					} else {
						idList = types.ListNull(types.StringType)
					}
				} else {
					idList = types.ListNull(types.StringType)
				}

				obj, objDiags := types.ObjectValue(serverGroupAttrTypes, map[string]attr.Value{
					"id": idList,
				})
				resp.Diagnostics.Append(objDiags...)
				if !resp.Diagnostics.HasError() {
					objectValues = append(objectValues, obj)
				}
			}
			if !resp.Diagnostics.HasError() {
				// Convert []types.Object to []attr.Value
				elements := make([]attr.Value, len(objectValues))
				for i, obj := range objectValues {
					elements[i] = obj
				}
				serverGroupsSet, setDiags := types.SetValue(serverGroupObjectType, elements)
				resp.Diagnostics.Append(setDiags...)
				serverGroups = serverGroupsSet
			} else {
				serverGroups = types.SetNull(serverGroupObjectType)
			}
		} else {
			serverGroups = types.SetNull(serverGroupObjectType)
		}
	} else {
		serverGroups = types.SetNull(serverGroupObjectType)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(segment.ID)
	data.Name = helpers.StringValueOrNull(segment.Name)
	data.SegmentGroupID = helpers.StringValueOrNull(segment.SegmentGroupID)
	data.SegmentGroupName = helpers.StringValueOrNull(segment.SegmentGroupName)
	data.BypassType = helpers.StringValueOrNull(segment.BypassType)
	data.ConfigSpace = helpers.StringValueOrNull(segment.ConfigSpace)
	data.CreationTime = helpers.StringValueOrNull(segment.CreationTime)
	data.DefaultIdleTimeout = helpers.StringValueOrNull(segment.DefaultIdleTimeout)
	data.DefaultMaxAge = helpers.StringValueOrNull(segment.DefaultMaxAge)
	data.Description = helpers.StringValueOrNull(segment.Description)
	data.DomainNames = domainNames
	data.DoubleEncrypt = types.BoolValue(segment.DoubleEncrypt)
	data.Enabled = types.BoolValue(segment.Enabled)
	data.HealthCheckType = helpers.StringValueOrNull(segment.HealthCheckType)
	data.HealthReporting = helpers.StringValueOrNull(segment.HealthReporting)
	data.SelectConnectorCloseToApp = types.BoolValue(segment.SelectConnectorCloseToApp)
	data.UseInDrMode = types.BoolValue(segment.UseInDrMode)
	data.IsIncompleteDRConfig = types.BoolValue(segment.IsIncompleteDRConfig)
	data.IPAnchored = types.BoolValue(segment.IpAnchored)
	data.IsCnameEnabled = types.BoolValue(segment.IsCnameEnabled)
	data.ModifiedBy = helpers.StringValueOrNull(segment.ModifiedBy)
	data.ModifiedTime = helpers.StringValueOrNull(segment.ModifiedTime)
	data.PassiveHealthEnabled = types.BoolValue(segment.PassiveHealthEnabled)
	data.APIProtectionEnabled = types.BoolValue(segment.APIProtectionEnabled)
	data.ServerGroups = serverGroups
	data.MicroTenantID = helpers.StringValueOrNull(segment.MicroTenantID)
	data.MicroTenantName = helpers.StringValueOrNull(segment.MicroTenantName)
	data.TCPPortRanges = tcpRanges
	data.UDPPortRanges = udpRanges
	data.TCPPortRange = tcpPorts
	data.UDPPortRange = udpPorts
	data.MatchStyle = helpers.StringValueOrNull(segment.MatchStyle)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
