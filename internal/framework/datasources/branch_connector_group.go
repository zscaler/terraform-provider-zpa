package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/branch_connector_group"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

var (
	_ datasource.DataSource              = &BranchConnectorGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &BranchConnectorGroupDataSource{}
)

func NewBranchConnectorGroupDataSource() datasource.DataSource {
	return &BranchConnectorGroupDataSource{}
}

type BranchConnectorGroupDataSource struct {
	client *client.Client
}

type BranchConnectorGroupModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func (d *BranchConnectorGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch_connector_group"
}

func (d *BranchConnectorGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA branch connector group summary by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the branch connector group.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the branch connector group.",
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (d *BranchConnectorGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BranchConnectorGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data BranchConnectorGroupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())

	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided.")
		return
	}

	service := d.client.Service

	var summary *common.CommonSummary
	var err error

	groups, _, err := branch_connector_group.GetBranchConnectorGroupSummary(ctx, service)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list branch connector groups: %v", err))
		return
	}

	for i := range groups {
		if id != "" && groups[i].ID == id {
			summary = &groups[i]
			break
		}
		if name != "" && strings.EqualFold(groups[i].Name, name) {
			summary = &groups[i]
			break
		}
	}

	if summary == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Branch connector group with id %q or name %q not found.", id, name))
		return
	}

	tflog.Debug(ctx, "Found branch connector group", map[string]any{"id": summary.ID, "name": summary.Name})

	data.ID = types.StringValue(summary.ID)
	data.Name = types.StringValue(summary.Name)
	data.Enabled = types.BoolValue(summary.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
