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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentbrowseraccess"
)

var (
	_ datasource.DataSource              = &ApplicationSegmentBrowserAccessDataSource{}
	_ datasource.DataSourceWithConfigure = &ApplicationSegmentBrowserAccessDataSource{}
)

func NewApplicationSegmentBrowserAccessDataSource() datasource.DataSource {
	return &ApplicationSegmentBrowserAccessDataSource{}
}

type ApplicationSegmentBrowserAccessDataSource struct {
	client *client.Client
}

type ApplicationSegmentBrowserAccessModel struct {
	ID                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	MicroTenantID             types.String `tfsdk:"microtenant_id"`
	MicroTenantName           types.String `tfsdk:"microtenant_name"`
	SegmentGroupID            types.String `tfsdk:"segment_group_id"`
	SegmentGroupName          types.String `tfsdk:"segment_group_name"`
	BypassType                types.String `tfsdk:"bypass_type"`
	ConfigSpace               types.String `tfsdk:"config_space"`
	Description               types.String `tfsdk:"description"`
	DomainNames               types.List   `tfsdk:"domain_names"`
	DoubleEncrypt             types.Bool   `tfsdk:"double_encrypt"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
	PassiveHealthEnabled      types.Bool   `tfsdk:"passive_health_enabled"`
	HealthCheckType           types.String `tfsdk:"health_check_type"`
	HealthReporting           types.String `tfsdk:"health_reporting"`
	MatchStyle                types.String `tfsdk:"match_style"`
	IPAnchored                types.Bool   `tfsdk:"ip_anchored"`
	IsCnameEnabled            types.Bool   `tfsdk:"is_cname_enabled"`
	SelectConnectorCloseToApp types.Bool   `tfsdk:"select_connector_close_to_app"`
	UseInDrMode               types.Bool   `tfsdk:"use_in_dr_mode"`
	IsIncompleteDRConfig      types.Bool   `tfsdk:"is_incomplete_dr_config"`
	APIProtectionEnabled      types.Bool   `tfsdk:"api_protection_enabled"`
	TCPPortRanges             types.List   `tfsdk:"tcp_port_ranges"`
	UDPPortRanges             types.List   `tfsdk:"udp_port_ranges"`
	TCPPortRange              types.Set    `tfsdk:"tcp_port_range"`
	UDPPortRange              types.Set    `tfsdk:"udp_port_range"`
	ServerGroups              types.Set    `tfsdk:"server_groups"`
	ClientlessApps            types.List   `tfsdk:"clientless_apps"`
}

func (d *ApplicationSegmentBrowserAccessDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_segment_browser_access"
}

func (d *ApplicationSegmentBrowserAccessDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a Browser Access application segment by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the browser access application segment.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the browser access application segment.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID used to scope the lookup.",
			},
			"microtenant_name":   schema.StringAttribute{Computed: true},
			"segment_group_id":   schema.StringAttribute{Computed: true},
			"segment_group_name": schema.StringAttribute{Computed: true},
			"bypass_type":        schema.StringAttribute{Computed: true},
			"config_space":       schema.StringAttribute{Computed: true},
			"description":        schema.StringAttribute{Computed: true},
			"domain_names": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"double_encrypt":                schema.BoolAttribute{Computed: true},
			"enabled":                       schema.BoolAttribute{Computed: true},
			"passive_health_enabled":        schema.BoolAttribute{Computed: true},
			"health_check_type":             schema.StringAttribute{Computed: true},
			"health_reporting":              schema.StringAttribute{Computed: true},
			"match_style":                   schema.StringAttribute{Computed: true},
			"ip_anchored":                   schema.BoolAttribute{Computed: true},
			"is_cname_enabled":              schema.BoolAttribute{Computed: true},
			"select_connector_close_to_app": schema.BoolAttribute{Computed: true},
			"use_in_dr_mode":                schema.BoolAttribute{Computed: true},
			"is_incomplete_dr_config":       schema.BoolAttribute{Computed: true},
			"api_protection_enabled":        schema.BoolAttribute{Computed: true},
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
			"clientless_apps": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":                   schema.StringAttribute{Computed: true},
						"name":                 schema.StringAttribute{Computed: true},
						"app_id":               schema.StringAttribute{Computed: true},
						"microtenant_id":       schema.StringAttribute{Computed: true},
						"microtenant_name":     schema.StringAttribute{Computed: true},
						"allow_options":        schema.BoolAttribute{Computed: true},
						"application_port":     schema.StringAttribute{Computed: true},
						"application_protocol": schema.StringAttribute{Computed: true},
						"certificate_id":       schema.StringAttribute{Computed: true},
						"certificate_name":     schema.StringAttribute{Computed: true},
						"cname":                schema.StringAttribute{Computed: true},
						"description":          schema.StringAttribute{Computed: true},
						"domain":               schema.StringAttribute{Computed: true},
						"enabled":              schema.BoolAttribute{Computed: true},
						"hidden":               schema.BoolAttribute{Computed: true},
						"local_domain":         schema.StringAttribute{Computed: true},
						"path":                 schema.StringAttribute{Computed: true},
						"trust_untrusted_cert": schema.BoolAttribute{Computed: true},
						"ext_label":            schema.StringAttribute{Computed: true},
						"ext_domain":           schema.StringAttribute{Computed: true},
					},
				},
			},
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
		},
	}
}

func (d *ApplicationSegmentBrowserAccessDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ApplicationSegmentBrowserAccessDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data ApplicationSegmentBrowserAccessModel
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
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided to read a browser access application segment.")
		return
	}

	var (
		segment *applicationsegmentbrowseraccess.BrowserAccess
		err     error
	)

	if id != "" {
		tflog.Debug(ctx, "Retrieving browser access application segment by ID", map[string]any{"id": id})
		segment, _, err = applicationsegmentbrowseraccess.Get(ctx, service, id)
	} else {
		tflog.Debug(ctx, "Retrieving browser access application segment by name", map[string]any{"name": name})
		segment, _, err = applicationsegmentbrowseraccess.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read browser access application segment: %v", err))
		return
	}

	if segment == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Browser access application segment with id %q or name %q was not found.", id, name))
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

	serverGroupsList, sgDiags := helpers.FlattenServerGroups(ctx, segment.AppServerGroups)
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

	clientlessApps, caDiags := flattenClientlessApps(ctx, segment.ClientlessApps)
	resp.Diagnostics.Append(caDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(segment.ID)
	data.Name = stringOrNull(segment.Name)
	data.MicroTenantID = stringOrNull(segment.MicroTenantID)
	data.MicroTenantName = stringOrNull(segment.MicroTenantName)
	data.SegmentGroupID = stringOrNull(segment.SegmentGroupID)
	data.SegmentGroupName = stringOrNull(segment.SegmentGroupName)
	data.BypassType = stringOrNull(segment.BypassType)
	data.ConfigSpace = stringOrNull(segment.ConfigSpace)
	data.Description = stringOrNull(segment.Description)
	data.DomainNames = domainNames
	data.DoubleEncrypt = types.BoolValue(segment.DoubleEncrypt)
	data.Enabled = types.BoolValue(segment.Enabled)
	data.PassiveHealthEnabled = types.BoolValue(segment.PassiveHealthEnabled)
	data.HealthCheckType = stringOrNull(segment.HealthCheckType)
	data.HealthReporting = stringOrNull(segment.HealthReporting)
	data.MatchStyle = stringOrNull(segment.MatchStyle)
	data.IPAnchored = types.BoolValue(segment.IPAnchored)
	data.IsCnameEnabled = types.BoolValue(segment.IsCnameEnabled)
	data.SelectConnectorCloseToApp = types.BoolValue(segment.SelectConnectorCloseToApp)
	data.UseInDrMode = types.BoolValue(segment.UseInDrMode)
	data.IsIncompleteDRConfig = types.BoolValue(segment.IsIncompleteDRConfig)
	data.APIProtectionEnabled = types.BoolValue(segment.APIProtectionEnabled)
	data.TCPPortRanges = tcpRanges
	data.UDPPortRanges = udpRanges
	data.TCPPortRange = tcpPorts
	data.UDPPortRange = udpPorts
	data.ServerGroups = serverGroups
	data.ClientlessApps = clientlessApps

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenClientlessApps(ctx context.Context, apps []applicationsegmentbrowseraccess.ClientlessApps) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id":                   types.StringType,
		"name":                 types.StringType,
		"app_id":               types.StringType,
		"microtenant_id":       types.StringType,
		"microtenant_name":     types.StringType,
		"allow_options":        types.BoolType,
		"application_port":     types.StringType,
		"application_protocol": types.StringType,
		"certificate_id":       types.StringType,
		"certificate_name":     types.StringType,
		"cname":                types.StringType,
		"description":          types.StringType,
		"domain":               types.StringType,
		"enabled":              types.BoolType,
		"hidden":               types.BoolType,
		"local_domain":         types.StringType,
		"path":                 types.StringType,
		"trust_untrusted_cert": types.BoolType,
		"ext_label":            types.StringType,
		"ext_domain":           types.StringType,
	}

	if len(apps) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(apps))
	var diags diag.Diagnostics
	for _, app := range apps {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":                   stringOrNull(app.ID),
			"name":                 stringOrNull(app.Name),
			"app_id":               stringOrNull(app.AppID),
			"microtenant_id":       stringOrNull(app.MicroTenantID),
			"microtenant_name":     stringOrNull(app.MicroTenantName),
			"allow_options":        types.BoolValue(app.AllowOptions),
			"application_port":     stringOrNull(app.ApplicationPort),
			"application_protocol": stringOrNull(app.ApplicationProtocol),
			"certificate_id":       stringOrNull(app.CertificateID),
			"certificate_name":     stringOrNull(app.CertificateName),
			"cname":                stringOrNull(app.Cname),
			"description":          stringOrNull(app.Description),
			"domain":               stringOrNull(app.Domain),
			"enabled":              types.BoolValue(app.Enabled),
			"hidden":               types.BoolValue(app.Hidden),
			"local_domain":         stringOrNull(app.LocalDomain),
			"path":                 stringOrNull(app.Path),
			"trust_untrusted_cert": types.BoolValue(app.TrustUntrustedCert),
			"ext_label":            stringOrNull(app.ExtLabel),
			"ext_domain":           stringOrNull(app.ExtDomain),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}
