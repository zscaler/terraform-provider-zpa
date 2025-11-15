package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/extranet_resource"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/location_controller"
)

var (
	_ datasource.DataSource              = &LocationControllerDataSource{}
	_ datasource.DataSourceWithConfigure = &LocationControllerDataSource{}
)

func NewLocationControllerDataSource() datasource.DataSource {
	return &LocationControllerDataSource{}
}

type LocationControllerDataSource struct {
	client *client.Client
}

type LocationControllerModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	ZIAErName types.String `tfsdk:"zia_er_name"`
	Enabled   types.Bool   `tfsdk:"enabled"`
}

func (d *LocationControllerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_location_controller"
}

func (d *LocationControllerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a location associated with an extranet resource partner.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the location to retrieve.",
			},
			"zia_er_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the extranet resource partner.",
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (d *LocationControllerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LocationControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data LocationControllerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	locationName := strings.TrimSpace(data.Name.ValueString())
	extranetName := strings.TrimSpace(data.ZIAErName.ValueString())

	if locationName == "" || extranetName == "" {
		resp.Diagnostics.AddError("Missing Required Attributes", "'name' and 'zia_er_name' must be provided.")
		return
	}

	tflog.Debug(ctx, "Resolving extranet resource partner", map[string]any{"name": extranetName})
	extranet, _, err := extranet_resource.GetExtranetResourcePartnerByName(ctx, d.client.Service, extranetName)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read extranet resource partner %q: %v", extranetName, err))
		return
	}
	if extranet == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Extranet resource partner %q not found.", extranetName))
		return
	}

	tflog.Debug(ctx, "Retrieving locations for extranet resource", map[string]any{"er_id": extranet.ID})
	locations, _, err := location_controller.GetLocationExtranetResource(ctx, d.client.Service, extranet.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read locations for extranet resource %q: %v", extranetName, err))
		return
	}

	var location *common.CommonSummary
	for i := range locations {
		if strings.EqualFold(locations[i].Name, locationName) {
			location = &locations[i]
			break
		}
	}

	if location == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Location %q not found for extranet resource %q.", locationName, extranetName))
		return
	}

	data.ID = types.StringValue(location.ID)
	data.Enabled = types.BoolValue(location.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
