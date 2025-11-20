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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/location_controller"
)

var (
	_ datasource.DataSource              = &LocationControllerSummaryDataSource{}
	_ datasource.DataSourceWithConfigure = &LocationControllerSummaryDataSource{}
)

func NewLocationControllerSummaryDataSource() datasource.DataSource {
	return &LocationControllerSummaryDataSource{}
}

type LocationControllerSummaryDataSource struct {
	client *client.Client
}

type LocationControllerSummaryModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func (d *LocationControllerSummaryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_location_controller_summary"
}

func (d *LocationControllerSummaryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a location controller summary entry by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the location.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the location.",
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (d *LocationControllerSummaryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LocationControllerSummaryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data LocationControllerSummaryModel
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

	var summary *common.CommonSummary
	var err error

	if id != "" {
		tflog.Debug(ctx, "Retrieving location controller summary by ID", map[string]any{"id": id})
		summaries, _, e := location_controller.GetLocationSummary(ctx, d.client.Service)
		err = e
		if err == nil {
			for i := range summaries {
				if summaries[i].ID == id {
					summary = &summaries[i]
					break
				}
			}
			if summary == nil {
				err = fmt.Errorf("location with id %q not found", id)
			}
		}
	} else {
		tflog.Debug(ctx, "Retrieving location controller summary by name", map[string]any{"name": name})
		summary, _, err = location_controller.GetLocationSummaryByName(ctx, d.client.Service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read location controller summary: %v", err))
		return
	}

	if summary == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Location summary with id %q or name %q not found.", id, name))
		return
	}

	data.ID = types.StringValue(summary.ID)
	data.Name = types.StringValue(summary.Name)
	data.Enabled = types.BoolValue(summary.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
