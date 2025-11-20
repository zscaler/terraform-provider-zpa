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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
)

var (
	_ datasource.DataSource              = &ServerGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &ServerGroupDataSource{}
)

func NewServerGroupDataSource() datasource.DataSource {
	return &ServerGroupDataSource{}
}

type ServerGroupDataSource struct {
	client *client.Client
}

type ServerGroupDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	MicroTenantID    types.String `tfsdk:"microtenant_id"`
	MicroTenantName  types.String `tfsdk:"microtenant_name"`
	ConfigSpace      types.String `tfsdk:"config_space"`
	CreationTime     types.String `tfsdk:"creation_time"`
	IPAanchored      types.Bool   `tfsdk:"ip_anchored"`
	DynamicDiscovery types.Bool   `tfsdk:"dynamic_discovery"`
	ModifiedBy       types.String `tfsdk:"modifiedby"`
	ModifiedTime     types.String `tfsdk:"modified_time"`
	Applications     types.Set    `tfsdk:"applications"`
	AppConnector     types.Set    `tfsdk:"app_connector_groups"`
	Servers          types.Set    `tfsdk:"servers"`
}

func (d *ServerGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_group"
}

func (d *ServerGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	nestedIDNameBlock := schema.SetNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"id":   schema.StringAttribute{Computed: true},
				"name": schema.StringAttribute{Computed: true},
			},
		},
	}

	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA Server Group by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id":                schema.StringAttribute{Optional: true},
			"name":              schema.StringAttribute{Optional: true},
			"description":       schema.StringAttribute{Computed: true},
			"enabled":           schema.BoolAttribute{Computed: true},
			"config_space":      schema.StringAttribute{Computed: true},
			"creation_time":     schema.StringAttribute{Computed: true},
			"ip_anchored":       schema.BoolAttribute{Computed: true},
			"dynamic_discovery": schema.BoolAttribute{Computed: true},
			"modifiedby":        schema.StringAttribute{Computed: true},
			"modified_time":     schema.StringAttribute{Computed: true},
			"microtenant_id":    schema.StringAttribute{Optional: true},
			"microtenant_name":  schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			"applications":         nestedIDNameBlock,
			"app_connector_groups": nestedIDNameBlock,
			"servers":              nestedIDNameBlock,
		},
	}
}

func (d *ServerGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServerGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ServerGroupDataSourceModel
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
		group *servergroup.ServerGroup
		err   error
	)

	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		id := data.ID.ValueString()
		tflog.Info(ctx, "Fetching server group", map[string]any{"id": id})
		group, _, err = servergroup.Get(ctx, service, id)
	} else {
		name := data.Name.ValueString()
		tflog.Info(ctx, "Fetching server group", map[string]any{"name": name})
		group, _, err = servergroup.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read server group: %v", err))
		return
	}

	flattened, diags := flattenServerGroup(ctx, group)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &flattened)...)
}

func flattenServerGroup(ctx context.Context, group *servergroup.ServerGroup) (ServerGroupDataSourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	applications, appDiags := flattenIDNameList(ctx, group.Applications)
	diags.Append(appDiags...)

	connectors, connectorDiags := flattenConnectorGroups(ctx, group.AppConnectorGroups)
	diags.Append(connectorDiags...)

	servers, serverDiags := flattenServersList(ctx, group.Servers)
	diags.Append(serverDiags...)

	model := ServerGroupDataSourceModel{
		ID:               types.StringValue(group.ID),
		Name:             types.StringValue(group.Name),
		Description:      types.StringValue(group.Description),
		Enabled:          types.BoolValue(group.Enabled),
		MicroTenantID:    types.StringValue(group.MicroTenantID),
		MicroTenantName:  types.StringValue(group.MicroTenantName),
		ConfigSpace:      types.StringValue(group.ConfigSpace),
		CreationTime:     types.StringValue(group.CreationTime),
		IPAanchored:      types.BoolValue(group.IpAnchored),
		DynamicDiscovery: types.BoolValue(group.DynamicDiscovery),
		ModifiedBy:       types.StringValue(group.ModifiedBy),
		ModifiedTime:     types.StringValue(group.ModifiedTime),
		Applications:     applications,
		AppConnector:     connectors,
		Servers:          servers,
	}

	return model, diags
}

func flattenIDNameList(ctx context.Context, items []servergroup.Applications) (types.Set, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
	values := make([]attr.Value, 0, len(items))

	for _, item := range items {
		obj, diags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":   types.StringValue(item.ID),
			"name": types.StringValue(item.Name),
		})
		if diags.HasError() {
			return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diags
		}
		values = append(values, obj)
	}

	set, diags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, values)
	return set, diags
}

func flattenConnectorGroups(ctx context.Context, groups []appconnectorgroup.AppConnectorGroup) (types.Set, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
	values := make([]attr.Value, 0, len(groups))

	for _, group := range groups {
		obj, diags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":   types.StringValue(group.ID),
			"name": types.StringValue(group.Name),
		})
		if diags.HasError() {
			return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diags
		}
		values = append(values, obj)
	}

	set, diags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, values)
	return set, diags
}

func flattenServersList(ctx context.Context, servers []appservercontroller.ApplicationServer) (types.Set, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
	values := make([]attr.Value, 0, len(servers))

	for _, server := range servers {
		obj, diags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":   types.StringValue(server.ID),
			"name": types.StringValue(server.Name),
		})
		if diags.HasError() {
			return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diags
		}
		values = append(values, obj)
	}

	set, diags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, values)
	return set, diags
}
