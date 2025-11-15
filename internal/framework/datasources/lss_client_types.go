package datasources

import (
	"context"
	"fmt"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/lssconfigcontroller"
)

var (
	_ datasource.DataSource              = &LSSClientTypesDataSource{}
	_ datasource.DataSourceWithConfigure = &LSSClientTypesDataSource{}
)

func NewLSSClientTypesDataSource() datasource.DataSource {
	return &LSSClientTypesDataSource{}
}

type LSSClientTypesDataSource struct {
	client *client.Client
}

type LSSClientTypesModel struct {
	ID                         types.String `tfsdk:"id"`
	ZPNClientTypeExporter      types.String `tfsdk:"zpn_client_type_exporter"`
	ZPNClientTypeMachineTunnel types.String `tfsdk:"zpn_client_type_machine_tunnel"`
	ZPNClientTypeIPAnchoring   types.String `tfsdk:"zpn_client_type_ip_anchoring"`
	ZPNClientTypeEdgeConnector types.String `tfsdk:"zpn_client_type_edge_connector"`
	ZPNClientTypeZAPP          types.String `tfsdk:"zpn_client_type_zapp"`
	ZPNClientTypeSlogger       types.String `tfsdk:"zpn_client_type_slogger"`
}

func (d *LSSClientTypesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lss_client_types"
}

func (d *LSSClientTypesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the available LSS client types.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"zpn_client_type_exporter": schema.StringAttribute{
				Computed: true,
			},
			"zpn_client_type_machine_tunnel": schema.StringAttribute{
				Computed: true,
			},
			"zpn_client_type_ip_anchoring": schema.StringAttribute{
				Computed: true,
			},
			"zpn_client_type_edge_connector": schema.StringAttribute{
				Computed: true,
			},
			"zpn_client_type_zapp": schema.StringAttribute{
				Computed: true,
			},
			"zpn_client_type_slogger": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *LSSClientTypesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LSSClientTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data LSSClientTypesModel
	tflog.Debug(ctx, "Retrieving LSS client types")

	typesResp, _, err := lssconfigcontroller.GetClientTypes(ctx, d.client.Service)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve LSS client types: %v", err))
		return
	}

	data.ID = types.StringValue("lss_client_types")
	data.ZPNClientTypeExporter = types.StringValue(typesResp.ZPNClientTypeExporter)
	data.ZPNClientTypeMachineTunnel = types.StringValue(typesResp.ZPNClientTypeMachineTunnel)
	data.ZPNClientTypeIPAnchoring = types.StringValue(typesResp.ZPNClientTypeIPAnchoring)
	data.ZPNClientTypeEdgeConnector = types.StringValue(typesResp.ZPNClientTypeEdgeConnector)
	data.ZPNClientTypeZAPP = types.StringValue(typesResp.ZPNClientTypeZAPP)
	data.ZPNClientTypeSlogger = types.StringValue(typesResp.ZPNClientTypeSlogger)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
