package datasources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/lssconfigcontroller"
)

var (
	_ datasource.DataSource              = &LSSStatusCodesDataSource{}
	_ datasource.DataSourceWithConfigure = &LSSStatusCodesDataSource{}
)

func NewLSSStatusCodesDataSource() datasource.DataSource {
	return &LSSStatusCodesDataSource{}
}

type LSSStatusCodesDataSource struct {
	client *client.Client
}

type LSSStatusCodesModel struct {
	ID            types.String `tfsdk:"id"`
	ZPNAuthLog    types.Map    `tfsdk:"zpn_auth_log"`
	ZPNAstAuthLog types.Map    `tfsdk:"zpn_ast_auth_log"`
	ZPNTransLog   types.Map    `tfsdk:"zpn_trans_log"`
	ZPNSysAuthLog types.Map    `tfsdk:"zpn_sys_auth_log"`
}

func (d *LSSStatusCodesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lss_config_status_codes"
}

func (d *LSSStatusCodesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the LSS status codes maps for each log type.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"zpn_auth_log": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Status codes for zpn_auth_log.",
			},
			"zpn_ast_auth_log": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Status codes for zpn_ast_auth_log.",
			},
			"zpn_trans_log": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Status codes for zpn_trans_log.",
			},
			"zpn_sys_auth_log": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Status codes for zpn_sys_auth_log.",
			},
		},
	}
}

func (d *LSSStatusCodesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LSSStatusCodesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data LSSStatusCodesModel
	tflog.Debug(ctx, "Retrieving LSS status codes")

	statusCodes, _, err := lssconfigcontroller.GetStatusCodes(ctx, d.client.Service)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve LSS status codes: %v", err))
		return
	}

	data.ID = types.StringValue("lss_status_codes")

	// Match SDKv2 exactly: zpn_auth_log uses ZPNAstAuthLog (line 76 in SDKv2)
	authLog, diags := mapInterfaceToStringMap(ctx, statusCodes.ZPNAstAuthLog)
	resp.Diagnostics.Append(diags...)
	astAuthLog, diags := mapInterfaceToStringMap(ctx, statusCodes.ZPNAstAuthLog)
	resp.Diagnostics.Append(diags...)
	transLog, diags := mapInterfaceToStringMap(ctx, statusCodes.ZPNTransLog)
	resp.Diagnostics.Append(diags...)
	sysAuthLog, diags := mapInterfaceToStringMap(ctx, statusCodes.ZPNSysAuthLog)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.ZPNAuthLog = authLog
	data.ZPNAstAuthLog = astAuthLog
	data.ZPNTransLog = transLog
	data.ZPNSysAuthLog = sysAuthLog

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapInterfaceToStringMap(ctx context.Context, input map[string]interface{}) (types.Map, diag.Diagnostics) {
	if len(input) == 0 {
		return types.MapNull(types.StringType), diag.Diagnostics{}
	}

	stringMap := make(map[string]string, len(input))
	for key, value := range input {
		bytes, err := json.MarshalIndent(&value, "", " ")
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("MarshalIndent failed for key %q: %v", key, err))
			continue
		}
		stringMap[key] = string(bytes)
	}

	mapValue, diags := types.MapValueFrom(ctx, types.StringType, stringMap)
	return mapValue, diags
}
