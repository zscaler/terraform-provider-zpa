package datasources

import (
	"context"
	"fmt"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

var (
	_ datasource.DataSource              = &RiskScoreValuesDataSource{}
	_ datasource.DataSourceWithConfigure = &RiskScoreValuesDataSource{}
)

func NewRiskScoreValuesDataSource() datasource.DataSource {
	return &RiskScoreValuesDataSource{}
}

type RiskScoreValuesDataSource struct {
	client *client.Client
}

type RiskScoreValuesModel struct {
	ID             types.String `tfsdk:"id"`
	ExcludeUnknown types.Bool   `tfsdk:"exclude_unknown"`
	Values         types.List   `tfsdk:"values"`
}

func (d *RiskScoreValuesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_risk_score_values"
}

func (d *RiskScoreValuesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of supported risk score values for policy conditions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Synthetic identifier for this data source.",
			},
			"exclude_unknown": schema.BoolAttribute{
				Optional:    true,
				Description: "Exclude unknown risk score values from the results.",
			},
			"values": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of risk score values.",
			},
		},
	}
}

func (d *RiskScoreValuesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RiskScoreValuesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data RiskScoreValuesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	excludeUnknown := false
	if !data.ExcludeUnknown.IsNull() && !data.ExcludeUnknown.IsUnknown() {
		excludeUnknown = data.ExcludeUnknown.ValueBool()
	}

	tflog.Debug(ctx, "Retrieving risk score values", map[string]any{
		"exclude_unknown": excludeUnknown,
	})

	values, _, err := policysetcontrollerv2.GetRiskScoreValues(ctx, d.client.Service, &excludeUnknown)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read risk score values: %v", err))
		return
	}

	list, diags := types.ListValueFrom(ctx, types.StringType, values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue("risk_score_values")
	data.Values = list
	data.ExcludeUnknown = types.BoolValue(excludeUnknown)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
