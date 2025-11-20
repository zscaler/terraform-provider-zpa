package resources

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgeschedule"
)

var (
	_ resource.Resource                = &ServiceEdgeAssistantScheduleResource{}
	_ resource.ResourceWithConfigure   = &ServiceEdgeAssistantScheduleResource{}
	_ resource.ResourceWithImportState = &ServiceEdgeAssistantScheduleResource{}
)

func NewServiceEdgeAssistantScheduleResource() resource.Resource {
	return &ServiceEdgeAssistantScheduleResource{}
}

type ServiceEdgeAssistantScheduleResource struct {
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

func (r *ServiceEdgeAssistantScheduleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_edge_assistant_schedule"
}

func (r *ServiceEdgeAssistantScheduleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the Service Edge assistant schedule used to delete inactive service edges.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"customer_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ZPA customer identifier. Defaults to the customer configured in the provider or ZPA_CUSTOMER_ID environment variable.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the schedule is enabled.",
			},
			"delete_disabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether disabled Service Edges are also deleted when the schedule runs.",
			},
			"frequency": schema.StringAttribute{
				Optional:    true,
				Description: "Frequency of the schedule. Currently only \"days\" is supported.",
				Validators: []validator.String{
					stringvalidator.OneOf("days"),
				},
			},
			"frequency_interval": schema.StringAttribute{
				Optional:    true,
				Description: "Frequency interval in days. Supported values: 5, 7, 14, 30, 60, 90.",
				Validators: []validator.String{
					stringvalidator.OneOf("5", "7", "14", "30", "60", "90"),
				},
			},
		},
	}
}

func (r *ServiceEdgeAssistantScheduleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cl, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = cl
}

func (r *ServiceEdgeAssistantScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider was not configured before use.")
		return
	}

	var plan ServiceEdgeAssistantScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	schedule, diags := r.buildSchedule(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating service edge assistant schedule")
	created, _, err := serviceedgeschedule.CreateSchedule(ctx, r.client.Service, schedule)
	if err != nil {
		if strings.Contains(err.Error(), "resource.already.exist") {
			tflog.Warn(ctx, "Schedule already exists, converting create to update")
			current, _, getErr := serviceedgeschedule.GetSchedule(ctx, r.client.Service)
			if getErr != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve existing schedule: %v", getErr))
				return
			}
			plan.ID = types.StringValue(current.ID)
			schedule.ID = current.ID
			if _, updErr := serviceedgeschedule.UpdateSchedule(ctx, r.client.Service, current.ID, &schedule); updErr != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update existing schedule: %v", updErr))
				return
			}
			created = current
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create schedule: %v", err))
			return
		}
	}

	plan.ID = types.StringValue(created.ID)

	state, readDiags := r.readSchedule(ctx)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = plan.ID
	state.CustomerID = plan.CustomerID
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceEdgeAssistantScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider was not configured before use.")
		return
	}

	var state ServiceEdgeAssistantScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readSchedule(ctx)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	newState.ID = state.ID
	newState.CustomerID = state.CustomerID
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ServiceEdgeAssistantScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider was not configured before use.")
		return
	}

	var plan ServiceEdgeAssistantScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	schedule, diags := r.buildSchedule(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := serviceedgeschedule.UpdateSchedule(ctx, r.client.Service, plan.ID.ValueString(), &schedule); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update schedule: %v", err))
		return
	}

	state, readDiags := r.readSchedule(ctx)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = plan.ID
	state.CustomerID = plan.CustomerID
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceEdgeAssistantScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Delete is a no-op as per SDKv2 implementation
}

func (r *ServiceEdgeAssistantScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before importing service edge assistant schedule.")
		return
	}

	state, diags := r.readSchedule(ctx)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceEdgeAssistantScheduleResource) buildSchedule(ctx context.Context, plan ServiceEdgeAssistantScheduleModel) (serviceedgeschedule.AssistantSchedule, diag.Diagnostics) {
	var diags diag.Diagnostics

	customerID := helpers.StringValue(plan.CustomerID)
	if customerID == "" {
		customerID = os.Getenv("ZPA_CUSTOMER_ID")
	}
	if customerID == "" {
		diags.AddError(
			"Missing customer_id",
			"customer_id must be provided either in the HCL or as an environment variable ZPA_CUSTOMER_ID",
		)
		return serviceedgeschedule.AssistantSchedule{}, diags
	}

	return serviceedgeschedule.AssistantSchedule{
		ID:                plan.ID.ValueString(),
		CustomerID:        customerID,
		Enabled:           helpers.BoolValue(plan.Enabled, false),
		DeleteDisabled:    helpers.BoolValue(plan.DeleteDisabled, false),
		FrequencyInterval: helpers.StringValue(plan.FrequencyInterval),
		Frequency:         helpers.StringValue(plan.Frequency),
	}, diags
}

func (r *ServiceEdgeAssistantScheduleResource) readSchedule(ctx context.Context) (ServiceEdgeAssistantScheduleModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	schedule, _, err := serviceedgeschedule.GetSchedule(ctx, r.client.Service)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return ServiceEdgeAssistantScheduleModel{}, diag.Diagnostics{
				diag.NewErrorDiagnostic("Not Found", "Service edge assistant schedule not found"),
			}
		}
		return ServiceEdgeAssistantScheduleModel{}, diag.Diagnostics{
			diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read service edge assistant schedule: %v", err)),
		}
	}

	return ServiceEdgeAssistantScheduleModel{
		ID:                types.StringValue(schedule.ID),
		CustomerID:        helpers.StringValueOrNull(schedule.CustomerID),
		Enabled:           types.BoolValue(schedule.Enabled),
		DeleteDisabled:    types.BoolValue(schedule.DeleteDisabled),
		Frequency:         helpers.StringValueOrNull(schedule.Frequency),
		FrequencyInterval: helpers.StringValueOrNull(schedule.FrequencyInterval),
	}, diags
}
