package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/trustednetwork"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
)

var (
	_ datasource.DataSource              = &TrustedNetworkDataSource{}
	_ datasource.DataSourceWithConfigure = &TrustedNetworkDataSource{}
)

func NewTrustedNetworkDataSource() datasource.DataSource {
	return &TrustedNetworkDataSource{}
}

type TrustedNetworkDataSource struct {
	client *client.Client
}

type TrustedNetworkDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	CreationTime types.String `tfsdk:"creation_time"`
	Domain       types.String `tfsdk:"domain"`
	ModifiedBy   types.String `tfsdk:"modified_by"`
	ModifiedTime types.String `tfsdk:"modified_time"`
	NetworkID    types.String `tfsdk:"network_id"`
	ZscalerCloud types.String `tfsdk:"zscaler_cloud"`
}

func (d *TrustedNetworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trusted_network"
}

func (d *TrustedNetworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA trusted network by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Optional: true},
			"name":          schema.StringAttribute{Optional: true},
			"creation_time": schema.StringAttribute{Computed: true},
			"domain":        schema.StringAttribute{Computed: true},
			"modified_by":   schema.StringAttribute{Computed: true},
			"modified_time": schema.StringAttribute{Computed: true},
			"network_id":    schema.StringAttribute{Computed: true},
			"zscaler_cloud": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *TrustedNetworkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TrustedNetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TrustedNetworkDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (data.ID.IsNull() || data.ID.ValueString() == "") && (data.Name.IsNull() || data.Name.ValueString() == "") {
		resp.Diagnostics.AddError("Missing Identifier", "Either id or name must be specified")
		return
	}

	service := d.client.Service

	var (
		net *trustednetwork.TrustedNetwork
		err error
	)

	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		id := data.ID.ValueString()
		tflog.Info(ctx, "Fetching trusted network", map[string]any{"id": id})
		net, _, err = trustednetwork.Get(ctx, service, id)
	} else {
		name := data.Name.ValueString()
		tflog.Info(ctx, "Fetching trusted network", map[string]any{"name": name})
		net, _, err = trustednetwork.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read trusted network: %v", err))
		return
	}

	model := TrustedNetworkDataSourceModel{
		ID:           types.StringValue(net.ID),
		Name:         types.StringValue(net.Name),
		CreationTime: types.StringValue(net.CreationTime),
		Domain:       types.StringValue(net.Domain),
		ModifiedBy:   types.StringValue(net.ModifiedBy),
		ModifiedTime: types.StringValue(net.ModifiedTime),
		NetworkID:    types.StringValue(net.NetworkID),
		ZscalerCloud: types.StringValue(net.ZscalerCloud),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
