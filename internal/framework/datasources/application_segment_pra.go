package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentpra"
)

var (
	_ datasource.DataSource              = &ApplicationSegmentPRADataSource{}
	_ datasource.DataSourceWithConfigure = &ApplicationSegmentPRADataSource{}
)

func NewApplicationSegmentPRADataSource() datasource.DataSource {
	return &ApplicationSegmentPRADataSource{}
}

type ApplicationSegmentPRADataSource struct {
	client *client.Client
}

type ApplicationSegmentPRAModel struct {
	ID                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	SegmentGroupID            types.String `tfsdk:"segment_group_id"`
	SegmentGroupName          types.String `tfsdk:"segment_group_name"`
	BypassType                types.String `tfsdk:"bypass_type"`
	ConfigSpace               types.String `tfsdk:"config_space"`
	Description               types.String `tfsdk:"description"`
	DomainNames               types.List   `tfsdk:"domain_names"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
	PassiveHealthEnabled      types.Bool   `tfsdk:"passive_health_enabled"`
	SelectConnectorCloseToApp types.Bool   `tfsdk:"select_connector_close_to_app"`
	DoubleEncrypt             types.Bool   `tfsdk:"double_encrypt"`
	MatchStyle                types.String `tfsdk:"match_style"`
	HealthCheckType           types.String `tfsdk:"health_check_type"`
	IsCnameEnabled            types.Bool   `tfsdk:"is_cname_enabled"`
	IPAnchored                types.Bool   `tfsdk:"ip_anchored"`
	HealthReporting           types.String `tfsdk:"health_reporting"`
	CreationTime              types.String `tfsdk:"creation_time"`
	ModifiedBy                types.String `tfsdk:"modified_by"`
	ModifiedTime              types.String `tfsdk:"modified_time"`
	UseInDrMode               types.Bool   `tfsdk:"use_in_dr_mode"`
	IsIncompleteDRConfig      types.Bool   `tfsdk:"is_incomplete_dr_config"`
	MicroTenantID             types.String `tfsdk:"microtenant_id"`
	MicroTenantName           types.String `tfsdk:"microtenant_name"`
	TCPPortRanges             types.List   `tfsdk:"tcp_port_ranges"`
	UDPPortRanges             types.List   `tfsdk:"udp_port_ranges"`
	TCPPortRange              types.List   `tfsdk:"tcp_port_range"`
	UDPPortRange              types.List   `tfsdk:"udp_port_range"`
	ServerGroups              types.List   `tfsdk:"server_groups"`
	SRAApps                   types.List   `tfsdk:"sra_apps"`
}

func (d *ApplicationSegmentPRADataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_segment_pra"
}

func (d *ApplicationSegmentPRADataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA Privileged Remote Access (PRA) application segment by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the PRA application segment.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the PRA application segment.",
			},
			"segment_group_id":   schema.StringAttribute{Computed: true},
			"segment_group_name": schema.StringAttribute{Computed: true},
			"bypass_type":        schema.StringAttribute{Computed: true},
			"config_space":       schema.StringAttribute{Computed: true},
			"description":        schema.StringAttribute{Computed: true},
			"domain_names": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"enabled":                       schema.BoolAttribute{Computed: true},
			"passive_health_enabled":        schema.BoolAttribute{Computed: true},
			"select_connector_close_to_app": schema.BoolAttribute{Computed: true},
			"double_encrypt":                schema.BoolAttribute{Computed: true},
			"match_style":                   schema.StringAttribute{Computed: true},
			"health_check_type":             schema.StringAttribute{Computed: true},
			"is_cname_enabled":              schema.BoolAttribute{Computed: true},
			"ip_anchored":                   schema.BoolAttribute{Computed: true},
			"health_reporting":              schema.StringAttribute{Computed: true},
			"creation_time":                 schema.StringAttribute{Computed: true},
			"modified_by":                   schema.StringAttribute{Computed: true},
			"modified_time":                 schema.StringAttribute{Computed: true},
			"use_in_dr_mode":                schema.BoolAttribute{Computed: true},
			"is_incomplete_dr_config":       schema.BoolAttribute{Computed: true},
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
			"sra_apps": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":                   schema.StringAttribute{Computed: true},
						"app_id":               schema.StringAttribute{Computed: true},
						"application_port":     schema.StringAttribute{Computed: true},
						"application_protocol": schema.StringAttribute{Computed: true},
						"certificate_id":       schema.StringAttribute{Computed: true},
						"certificate_name":     schema.StringAttribute{Computed: true},
						"connection_security":  schema.StringAttribute{Computed: true},
						"description":          schema.StringAttribute{Computed: true},
						"domain":               schema.StringAttribute{Computed: true},
						"enabled":              schema.BoolAttribute{Computed: true},
						"hidden":               schema.BoolAttribute{Computed: true},
						"name":                 schema.StringAttribute{Computed: true},
						"microtenant_id":       schema.StringAttribute{Computed: true},
						"microtenant_name":     schema.StringAttribute{Computed: true},
						"portal":               schema.BoolAttribute{Computed: true},
					},
				},
			},
			"tcp_port_range": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"from": schema.StringAttribute{Computed: true},
						"to":   schema.StringAttribute{Computed: true},
					},
				},
			},
			"udp_port_range": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"from": schema.StringAttribute{Computed: true},
						"to":   schema.StringAttribute{Computed: true},
					},
				},
			},
			"server_groups": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *ApplicationSegmentPRADataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ApplicationSegmentPRADataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data ApplicationSegmentPRAModel
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
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided to read a PRA application segment.")
		return
	}

	var (
		segment *applicationsegmentpra.AppSegmentPRA
		err     error
	)

	if id != "" {
		tflog.Debug(ctx, "Retrieving PRA application segment by ID", map[string]any{"id": id})
		segment, _, err = applicationsegmentpra.Get(ctx, service, id)
	} else {
		tflog.Debug(ctx, "Retrieving PRA application segment by name", map[string]any{"name": name})
		segment, _, err = applicationsegmentpra.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read PRA application segment: %v", err))
		return
	}

	if segment == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("PRA application segment with id %q or name %q was not found.", id, name))
		return
	}

	domainNames, domainDiags := helpers.StringSliceToList(ctx, segment.DomainNames)
	resp.Diagnostics.Append(domainDiags...)

	tcpRanges, tcpRangesDiags := helpers.StringSliceToList(ctx, segment.TCPPortRanges)
	resp.Diagnostics.Append(tcpRangesDiags...)

	udpRanges, udpRangesDiags := helpers.StringSliceToList(ctx, segment.UDPPortRanges)
	resp.Diagnostics.Append(udpRangesDiags...)

	tcpPorts, tcpDiags := helpers.FlattenNetworkPorts(ctx, segment.TCPAppPortRange)
	resp.Diagnostics.Append(tcpDiags...)

	udpPorts, udpDiags := helpers.FlattenNetworkPorts(ctx, segment.UDPAppPortRange)
	resp.Diagnostics.Append(udpDiags...)

	serverGroups, sgDiags := helpers.FlattenServerGroups(ctx, segment.ServerGroups)
	resp.Diagnostics.Append(sgDiags...)

	sraApps, sraDiags := flattenPRAApps(ctx, segment.PRAApps)
	resp.Diagnostics.Append(sraDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(segment.ID)
	data.Name = stringOrNull(segment.Name)
	data.SegmentGroupID = stringOrNull(segment.SegmentGroupID)
	data.SegmentGroupName = stringOrNull(segment.SegmentGroupName)
	data.BypassType = stringOrNull(segment.BypassType)
	data.ConfigSpace = stringOrNull(segment.ConfigSpace)
	data.Description = stringOrNull(segment.Description)
	data.DomainNames = domainNames
	data.Enabled = types.BoolValue(segment.Enabled)
	data.PassiveHealthEnabled = types.BoolValue(segment.PassiveHealthEnabled)
	data.SelectConnectorCloseToApp = types.BoolValue(segment.SelectConnectorCloseToApp)
	data.DoubleEncrypt = types.BoolValue(segment.DoubleEncrypt)
	data.MatchStyle = stringOrNull(segment.MatchStyle)
	data.HealthCheckType = stringOrNull(segment.HealthCheckType)
	data.IsCnameEnabled = types.BoolValue(segment.IsCnameEnabled)
	data.IPAnchored = types.BoolValue(segment.IpAnchored)
	data.HealthReporting = stringOrNull(segment.HealthReporting)
	data.CreationTime = stringOrNull(segment.CreationTime)
	data.ModifiedBy = stringOrNull(segment.ModifiedBy)
	data.ModifiedTime = stringOrNull(segment.ModifiedTime)
	data.UseInDrMode = types.BoolValue(segment.UseInDrMode)
	data.IsIncompleteDRConfig = types.BoolValue(segment.IsIncompleteDRConfig)
	data.MicroTenantID = stringOrNull(segment.MicroTenantID)
	data.MicroTenantName = stringOrNull(segment.MicroTenantName)
	data.TCPPortRanges = tcpRanges
	data.UDPPortRanges = udpRanges
	data.TCPPortRange = tcpPorts
	data.UDPPortRange = udpPorts
	data.ServerGroups = serverGroups
	data.SRAApps = sraApps

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenPRAApps(ctx context.Context, apps []applicationsegmentpra.PRAApps) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id":                   types.StringType,
		"app_id":               types.StringType,
		"application_port":     types.StringType,
		"application_protocol": types.StringType,
		"certificate_id":       types.StringType,
		"certificate_name":     types.StringType,
		"connection_security":  types.StringType,
		"description":          types.StringType,
		"domain":               types.StringType,
		"enabled":              types.BoolType,
		"hidden":               types.BoolType,
		"name":                 types.StringType,
		"microtenant_id":       types.StringType,
		"microtenant_name":     types.StringType,
		"portal":               types.BoolType,
	}

	if len(apps) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(apps))
	var diags diag.Diagnostics
	for _, app := range apps {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":                   stringOrNull(app.ID),
			"app_id":               stringOrNull(app.AppID),
			"application_port":     stringOrNull(app.ApplicationPort),
			"application_protocol": stringOrNull(app.ApplicationProtocol),
			"certificate_id":       stringOrNull(app.CertificateID),
			"certificate_name":     stringOrNull(app.CertificateName),
			"connection_security":  stringOrNull(app.ConnectionSecurity),
			"description":          stringOrNull(app.Description),
			"domain":               stringOrNull(app.Domain),
			"enabled":              types.BoolValue(app.Enabled),
			"hidden":               types.BoolValue(app.Hidden),
			"name":                 stringOrNull(app.Name),
			"microtenant_id":       stringOrNull(app.MicroTenantID),
			"microtenant_name":     stringOrNull(app.MicroTenantName),
			"portal":               types.BoolValue(app.Portal),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}
