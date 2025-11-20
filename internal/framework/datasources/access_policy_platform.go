package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/platforms"
)

var (
	_ datasource.DataSource              = &AccessPolicyPlatformDataSource{}
	_ datasource.DataSourceWithConfigure = &AccessPolicyPlatformDataSource{}
)

func NewAccessPolicyPlatformDataSource() datasource.DataSource {
	return &AccessPolicyPlatformDataSource{}
}

type AccessPolicyPlatformDataSource struct {
	client *client.Client
}

type AccessPolicyPlatformModel struct {
	ID      types.String `tfsdk:"id"`
	Linux   types.String `tfsdk:"linux"`
	Android types.String `tfsdk:"android"`
	Windows types.String `tfsdk:"windows"`
	IOS     types.String `tfsdk:"ios"`
	Mac     types.String `tfsdk:"mac"`
}

func (d *AccessPolicyPlatformDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_policy_platform"
}

func (d *AccessPolicyPlatformDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the available platform identifiers for ZPA access policy conditions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Synthetic identifier for this data source.",
			},
			"linux": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier representing Linux clients.",
			},
			"android": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier representing Android clients.",
			},
			"windows": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier representing Windows clients.",
			},
			"ios": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier representing iOS clients.",
			},
			"mac": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier representing macOS clients.",
			},
		},
	}
}

func (d *AccessPolicyPlatformDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AccessPolicyPlatformDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data AccessPolicyPlatformModel

	tflog.Debug(ctx, "Retrieving access policy platforms")
	platformsResp, _, err := platforms.GetAllPlatforms(ctx, d.client.Service)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read access policy platforms: %v", err))
		return
	}

	data.Linux = types.StringValue(platformsResp.Linux)
	data.Android = types.StringValue(platformsResp.Android)
	data.Windows = types.StringValue(platformsResp.Windows)
	data.IOS = types.StringValue(platformsResp.IOS)
	data.Mac = types.StringValue(platformsResp.MacOS)
	data.ID = types.StringValue("platforms")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
