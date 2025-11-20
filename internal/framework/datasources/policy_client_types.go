package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/clienttypes"
)

var (
	_ datasource.DataSource              = &AccessPolicyClientTypesDataSource{}
	_ datasource.DataSourceWithConfigure = &AccessPolicyClientTypesDataSource{}
)

func NewAccessPolicyClientTypesDataSource() datasource.DataSource {
	return &AccessPolicyClientTypesDataSource{}
}

type AccessPolicyClientTypesDataSource struct {
	client *client.Client
}

type AccessPolicyClientTypesModel struct {
	ID                            types.String `tfsdk:"id"`
	ZPNClientTypeExporter         types.String `tfsdk:"zpn_client_type_exporter"`
	ZPNClientTypeExporterNoAuth   types.String `tfsdk:"zpn_client_type_exporter_noauth"`
	ZPNClientTypeBrowserIsolation types.String `tfsdk:"zpn_client_type_browser_isolation"`
	ZPNClientTypeMachineTunnel    types.String `tfsdk:"zpn_client_type_machine_tunnel"`
	ZPNClientTypeIPAnchoring      types.String `tfsdk:"zpn_client_type_ip_anchoring"`
	ZPNClientTypeEdgeConnector    types.String `tfsdk:"zpn_client_type_edge_connector"`
	ZPNClientTypeZAPP             types.String `tfsdk:"zpn_client_type_zapp"`
	ZPNClientTypeSLogger          types.String `tfsdk:"zpn_client_type_slogger"`
	ZPNClientTypeBranchConnector  types.String `tfsdk:"zpn_client_type_branch_connector"`
	ZPNClientTypeZAPPPartner      types.String `tfsdk:"zpn_client_type_zapp_partner"`
	ZPNClientTypeVDI              types.String `tfsdk:"zpn_client_type_vdi"`
	ZPNClientTypeZIAInspection    types.String `tfsdk:"zpn_client_type_zia_inspection"`
}

func (d *AccessPolicyClientTypesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_policy_client_types"
}

func (d *AccessPolicyClientTypesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the available ZPA client type identifiers for use in access policies.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Synthetic identifier for this data source.",
			},
			"zpn_client_type_exporter": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for the ZPN client type exporter.",
			},
			"zpn_client_type_exporter_noauth": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for the ZPN client type exporter without authentication.",
			},
			"zpn_client_type_browser_isolation": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for the ZPN browser isolation client type.",
			},
			"zpn_client_type_machine_tunnel": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for the ZPN machine tunnel client type.",
			},
			"zpn_client_type_ip_anchoring": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for the ZPN IP anchoring client type.",
			},
			"zpn_client_type_edge_connector": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for the ZPN edge connector client type.",
			},
			"zpn_client_type_zapp": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for the Zscaler Client Connector (Z App) client type.",
			},
			"zpn_client_type_slogger": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for the ZPN slogger client type.",
			},
			"zpn_client_type_branch_connector": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for the ZPN branch connector client type.",
			},
			"zpn_client_type_zapp_partner": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for the ZPN partner client type.",
			},
			"zpn_client_type_vdi": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for the ZPN VDI client type.",
			},
			"zpn_client_type_zia_inspection": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for the ZPA ZIA inspection client type.",
			},
		},
	}
}

func (d *AccessPolicyClientTypesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AccessPolicyClientTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data AccessPolicyClientTypesModel

	tflog.Debug(ctx, "Retrieving access policy client types")
	clientTypesResp, _, err := clienttypes.GetAllClientTypes(ctx, d.client.Service)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read access policy client types: %v", err))
		return
	}

	data.ID = types.StringValue("access_policy_client_types")
	data.ZPNClientTypeExporter = types.StringValue(clientTypesResp.ZPNClientTypeExplorer)
	data.ZPNClientTypeExporterNoAuth = types.StringValue(clientTypesResp.ZPNClientTypeNoAuth)
	data.ZPNClientTypeBrowserIsolation = types.StringValue(clientTypesResp.ZPNClientTypeBrowserIsolation)
	data.ZPNClientTypeMachineTunnel = types.StringValue(clientTypesResp.ZPNClientTypeMachineTunnel)
	data.ZPNClientTypeIPAnchoring = types.StringValue(clientTypesResp.ZPNClientTypeIPAnchoring)
	data.ZPNClientTypeEdgeConnector = types.StringValue(clientTypesResp.ZPNClientTypeEdgeConnector)
	data.ZPNClientTypeZAPP = types.StringValue(clientTypesResp.ZPNClientTypeZAPP)
	data.ZPNClientTypeSLogger = types.StringValue(clientTypesResp.ZPNClientTypeSlogger)
	data.ZPNClientTypeBranchConnector = types.StringValue(clientTypesResp.ZPNClientTypeBranchConnector)
	data.ZPNClientTypeZAPPPartner = types.StringValue(clientTypesResp.ZPNClientTypePartner)
	data.ZPNClientTypeVDI = types.StringValue(clientTypesResp.ZPNClientTypeVDI)
	data.ZPNClientTypeZIAInspection = types.StringValue(clientTypesResp.ZPNClientTypeZIAInspection)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
