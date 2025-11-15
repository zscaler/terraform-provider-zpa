package datasources

import (
	"context"
	"fmt"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/lssconfigcontroller"
)

var (
	_ datasource.DataSource              = &LSSLogTypeFormatsDataSource{}
	_ datasource.DataSourceWithConfigure = &LSSLogTypeFormatsDataSource{}
)

func NewLSSLogTypeFormatsDataSource() datasource.DataSource {
	return &LSSLogTypeFormatsDataSource{}
}

type LSSLogTypeFormatsDataSource struct {
	client *client.Client
}

type LSSLogTypeFormatsModel struct {
	ID      types.String `tfsdk:"id"`
	LogType types.String `tfsdk:"log_type"`
	TSV     types.String `tfsdk:"tsv"`
	CSV     types.String `tfsdk:"csv"`
	JSON    types.String `tfsdk:"json"`
}

func (d *LSSLogTypeFormatsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lss_config_log_type_formats"
}

func (d *LSSLogTypeFormatsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves available log formats for a specified LSS log type.",
		Attributes: map[string]schema.Attribute{
			"log_type": schema.StringAttribute{
				Required:    true,
				Description: "The LSS log type to query.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"zpn_ast_comprehensive_stats",
						"zpn_auth_log",
						"zpn_pbroker_comprehensive_stats",
						"zpn_ast_auth_log",
						"zpn_audit_log",
						"zpn_trans_log",
						"zpn_http_trans_log",
						"zpn_waf_http_exchanges_log",
						"zpn_sys_auth_log",
						"zpn_smb_inspection_log",
						"zpn_auth_log_1id",
						"zpn_sitec_auth_log",
						"zpn_sitec_comprehensive_stats",
						"zpn_ldap_inspection_log",
						"zms_flow_log",
						"zpn_krb_inspection_log",
					),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"tsv": schema.StringAttribute{
				Computed:    true,
				Description: "The TSV format for the provided log type.",
			},
			"csv": schema.StringAttribute{
				Computed:    true,
				Description: "The CSV format for the provided log type.",
			},
			"json": schema.StringAttribute{
				Computed:    true,
				Description: "The JSON format for the provided log type.",
			},
		},
	}
}

func (d *LSSLogTypeFormatsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LSSLogTypeFormatsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data LSSLogTypeFormatsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	logType := data.LogType.ValueString()
	tflog.Debug(ctx, "Retrieving LSS log type formats", map[string]any{"log_type": logType})

	formats, _, err := lssconfigcontroller.GetFormats(ctx, d.client.Service, logType)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve LSS log type formats: %v", err))
		return
	}

	data.ID = types.StringValue("lss_log_types_" + logType)
	data.TSV = types.StringValue(formats.Tsv)
	data.CSV = types.StringValue(formats.Csv)
	data.JSON = types.StringValue(formats.Json)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
