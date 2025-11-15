package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
)

var (
	_ datasource.DataSource              = &AppServerControllerDataSource{}
	_ datasource.DataSourceWithConfigure = &AppServerControllerDataSource{}
)

func NewAppServerControllerDataSource() datasource.DataSource {
	return &AppServerControllerDataSource{}
}

type AppServerControllerDataSource struct {
	client *client.Client
}

type AppServerControllerDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	Address           types.String `tfsdk:"address"`
	AppServerGroupIDs types.Set    `tfsdk:"app_server_group_ids"`
	ConfigSpace       types.String `tfsdk:"config_space"`
	CreationTime      types.String `tfsdk:"creation_time"`
	ModifiedBy        types.String `tfsdk:"modifiedby"`
	ModifiedTime      types.String `tfsdk:"modified_time"`
	MicroTenantID     types.String `tfsdk:"microtenant_id"`
	MicroTenantName   types.String `tfsdk:"microtenant_name"`
}

func (d *AppServerControllerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_server"
}

func (d *AppServerControllerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a ZPA Application Server (App Connector Server)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"description": schema.StringAttribute{Computed: true},
			"enabled":     schema.BoolAttribute{Computed: true},
			"address":     schema.StringAttribute{Computed: true},
			"app_server_group_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"config_space":  schema.StringAttribute{Computed: true},
			"creation_time": schema.StringAttribute{Computed: true},
			"modifiedby":    schema.StringAttribute{Computed: true},
			"modified_time": schema.StringAttribute{Computed: true},
			"microtenant_id": schema.StringAttribute{
				Optional: true,
			},
			"microtenant_name": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *AppServerControllerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AppServerControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AppServerControllerDataSourceModel
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
		server *appservercontroller.ApplicationServer
		err    error
	)

	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		id := data.ID.ValueString()
		tflog.Info(ctx, "Fetching application server", map[string]any{"id": id})
		server, _, err = appservercontroller.Get(ctx, service, id)
	} else {
		name := data.Name.ValueString()
		tflog.Info(ctx, "Fetching application server", map[string]any{"name": name})
		server, _, err = appservercontroller.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read application server: %v", err))
		return
	}

	fetched, diags := flattenApplicationServer(ctx, server)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data = fetched

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenApplicationServer(ctx context.Context, server *appservercontroller.ApplicationServer) (AppServerControllerDataSourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	appServerGroupIDs, setDiags := types.SetValueFrom(ctx, types.StringType, server.AppServerGroupIds)
	diags.Append(setDiags...)

	model := AppServerControllerDataSourceModel{
		ID:                types.StringValue(server.ID),
		Name:              types.StringValue(server.Name),
		Description:       types.StringValue(server.Description),
		Enabled:           types.BoolValue(server.Enabled),
		Address:           types.StringValue(server.Address),
		AppServerGroupIDs: appServerGroupIDs,
		ConfigSpace:       types.StringValue(server.ConfigSpace),
		CreationTime:      types.StringValue(server.CreationTime),
		ModifiedBy:        types.StringValue(server.ModifiedBy),
		ModifiedTime:      types.StringValue(server.ModifiedTime),
		MicroTenantID:     types.StringValue(server.MicroTenantID),
		MicroTenantName:   types.StringValue(server.MicroTenantName),
	}

	return model, diags
}
