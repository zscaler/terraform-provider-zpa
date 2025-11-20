package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/custom_config_controller"
)

var (
	_ datasource.DataSource              = &CloudConfigDataSource{}
	_ datasource.DataSourceWithConfigure = &CloudConfigDataSource{}
)

func NewCloudConfigDataSource() datasource.DataSource {
	return &CloudConfigDataSource{}
}

type CloudConfigDataSource struct {
	client *client.Client
}

type CloudConfigModel struct {
	ID             types.String `tfsdk:"id"`
	ZIACloudDomain types.String `tfsdk:"zia_cloud_domain"`
	ZIAUsername    types.String `tfsdk:"zia_username"`
}

func (d *CloudConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_config"
}

func (d *CloudConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the associated ZIA cloud configuration for the tenant.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"zia_cloud_domain": schema.StringAttribute{
				Computed: true,
			},
			"zia_username": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *CloudConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	tflog.Debug(ctx, "Retrieving ZIA cloud configuration")
	config, _, err := custom_config_controller.GetZIACloudConfig(ctx, d.client.Service)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ZIA cloud config: %v", err))
		return
	}

	if config == nil {
		resp.Diagnostics.AddError("Not Found", "ZIA cloud configuration not found.")
		return
	}

	state := CloudConfigModel{
		ID:             types.StringValue("zia_cloud_config"),
		ZIACloudDomain: types.StringValue(config.ZIACloudDomain),
		ZIAUsername:    types.StringValue(config.ZIAUsername),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
