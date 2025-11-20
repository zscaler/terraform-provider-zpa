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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgeschedule"
)

var (
	_ datasource.DataSource              = &ServiceEdgeAssistantScheduleDataSource{}
	_ datasource.DataSourceWithConfigure = &ServiceEdgeAssistantScheduleDataSource{}
)

func NewServiceEdgeAssistantScheduleDataSource() datasource.DataSource {
	return &ServiceEdgeAssistantScheduleDataSource{}
}

type ServiceEdgeAssistantScheduleDataSource struct {
	client *client.Client
}

type ServiceEdgeAssistantScheduleModel struct {
	ID                types.String `tfsdk:"id"`
	CustomerID        types.String `tfsdk:"customer_id"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	DeleteDisabled    types.Bool   `tfsdk:"delete_disabled"`
	Frequency         types.String `tfsdk:"frequency"`
	FrequencyInterval types.String `tfsdk:"frequency_interval"`
}

func (d *ServiceEdgeAssistantScheduleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_edge_assistant_schedule"
}

func (d *ServiceEdgeAssistantScheduleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the Service Edge Assistant schedule configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the assistant schedule.",
			},
			"customer_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Customer identifier associated with the schedule.",
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether the assistant schedule is enabled.",
			},
			"delete_disabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether deleting the schedule is disabled.",
			},
			"frequency": schema.StringAttribute{
				Computed:    true,
				Description: "Frequency of the assistant schedule.",
			},
			"frequency_interval": schema.StringAttribute{
				Computed:    true,
				Description: "Frequency interval of the assistant schedule.",
			},
		},
	}
}

func (d *ServiceEdgeAssistantScheduleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceEdgeAssistantScheduleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data ServiceEdgeAssistantScheduleModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestedID := strings.TrimSpace(data.ID.ValueString())
	requestedCustomerID := strings.TrimSpace(data.CustomerID.ValueString())

	tflog.Debug(ctx, "Retrieving service edge assistant schedule", map[string]any{
		"id":          requestedID,
		"customer_id": requestedCustomerID,
	})

	schedule, _, err := serviceedgeschedule.GetSchedule(ctx, d.client.Service)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service edge assistant schedule: %v", err))
		return
	}

	if schedule == nil {
		resp.Diagnostics.AddError("Not Found", "Service edge assistant schedule was not found.")
		return
	}

	data.ID = types.StringValue(schedule.ID)
	data.CustomerID = types.StringValue(schedule.CustomerID)
	data.Enabled = types.BoolValue(schedule.Enabled)
	data.DeleteDisabled = types.BoolValue(schedule.DeleteDisabled)
	data.Frequency = types.StringValue(schedule.Frequency)
	data.FrequencyInterval = types.StringValue(schedule.FrequencyInterval)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
