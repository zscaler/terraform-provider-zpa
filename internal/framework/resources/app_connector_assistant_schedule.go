package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorschedule"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &AppConnectorAssistantScheduleResource{}
	_ resource.ResourceWithConfigure   = &AppConnectorAssistantScheduleResource{}
	_ resource.ResourceWithImportState = &AppConnectorAssistantScheduleResource{}
)

func NewAppConnectorAssistantScheduleResource() resource.Resource {
	return &AppConnectorAssistantScheduleResource{}
}

type AppConnectorAssistantScheduleResource struct {
	client *client.Client
}

type AppConnectorAssistantScheduleModel struct {
	ID                types.String `tfsdk:"id"`
	CustomerID        types.String `tfsdk:"customer_id"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	DeleteDisabled    types.Bool   `tfsdk:"delete_disabled"`
	Frequency         types.String `tfsdk:"frequency"`
	FrequencyInterval types.String `tfsdk:"frequency_interval"`
}

func (r *AppConnectorAssistantScheduleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_connector_assistant_schedule"
}

func (r *AppConnectorAssistantScheduleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the App Connector assistant schedule used to delete inactive connectors.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"customer_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ZPA customer identifier. Defaults to the customer configured in the provider.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the schedule is enabled.",
			},
			"delete_disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether disabled App Connectors are also deleted when the schedule runs.",
			},
			"frequency": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Frequency of the schedule. Currently only \"days\" is supported.",
				Validators: []validator.String{
					stringvalidator.OneOf("days"),
				},
			},
			"frequency_interval": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Frequency interval in days. Supported values: 5, 7, 14, 30, 60, 90.",
				Validators: []validator.String{
					stringvalidator.OneOf("5", "7", "14", "30", "60", "90"),
				},
			},
		},
	}
}

func (r *AppConnectorAssistantScheduleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AppConnectorAssistantScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider was not configured before use.")
		return
	}

	var plan AppConnectorAssistantScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	schedule, diags := r.buildSchedule(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating app connector assistant schedule")
	created, _, err := appconnectorschedule.CreateSchedule(ctx, r.client.Service, schedule)
	if err != nil {
		if strings.Contains(err.Error(), "resource.already.exist") {
			tflog.Warn(ctx, "Schedule already exists, converting create to update")
			current, _, getErr := appconnectorschedule.GetSchedule(ctx, r.client.Service)
			if getErr != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve existing schedule: %v", getErr))
				return
			}
			plan.ID = types.StringValue(current.ID)
			schedule.ID = current.ID
			if _, updErr := appconnectorschedule.UpdateSchedule(ctx, r.client.Service, current.ID, &schedule); updErr != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update existing schedule: %v", updErr))
				return
			}
			created = current
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create schedule: %v", err))
			return
		}
	}

	state, diags := r.readSchedule(ctx, created)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AppConnectorAssistantScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider was not configured before use.")
		return
	}

	var state AppConnectorAssistantScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	schedule, _, err := appconnectorschedule.GetSchedule(ctx, r.client.Service)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read schedule: %v", err))
		return
	}

	newState, diags := r.readSchedule(ctx, schedule)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *AppConnectorAssistantScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider was not configured before use.")
		return
	}

	var plan AppConnectorAssistantScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	schedule, diags := r.buildSchedule(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(plan.ID.ValueString())
	if id == "" {
		current, _, err := appconnectorschedule.GetSchedule(ctx, r.client.Service)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve existing schedule: %v", err))
			return
		}
		id = current.ID
		schedule.ID = id
	}

	if _, err := appconnectorschedule.UpdateSchedule(ctx, r.client.Service, id, &schedule); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update schedule: %v", err))
		return
	}

	latest, _, err := appconnectorschedule.GetSchedule(ctx, r.client.Service)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to refresh schedule: %v", err))
		return
	}

	state, diags := r.readSchedule(ctx, latest)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AppConnectorAssistantScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// The API does not support deletion of the schedule; simply forget the resource.
	resp.State.RemoveResource(ctx)
}

func (r *AppConnectorAssistantScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, pathRootID, types.StringValue(req.ID))...)
}

func (r *AppConnectorAssistantScheduleResource) buildSchedule(ctx context.Context, plan AppConnectorAssistantScheduleModel) (appconnectorschedule.AssistantSchedule, diag.Diagnostics) {
	var diags diag.Diagnostics

	customerID := strings.TrimSpace(plan.CustomerID.ValueString())
	if customerID == "" {
		customerID = r.client.Service.Client.GetCustomerID()
	}

	schedule := appconnectorschedule.AssistantSchedule{
		ID:                strings.TrimSpace(plan.ID.ValueString()),
		CustomerID:        customerID,
		Enabled:           helpers.BoolValue(plan.Enabled, false),
		DeleteDisabled:    helpers.BoolValue(plan.DeleteDisabled, false),
		Frequency:         strings.TrimSpace(plan.Frequency.ValueString()),
		FrequencyInterval: strings.TrimSpace(plan.FrequencyInterval.ValueString()),
	}

	if schedule.CustomerID == "" {
		diags.AddError("Missing customer ID", "customer_id must be specified either in the resource or provider configuration.")
	}

	return schedule, diags
}

func (r *AppConnectorAssistantScheduleResource) readSchedule(ctx context.Context, schedule *appconnectorschedule.AssistantSchedule) (AppConnectorAssistantScheduleModel, diag.Diagnostics) {
	if schedule == nil {
		return AppConnectorAssistantScheduleModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Empty response", "The schedule response was empty")}
	}

	model := AppConnectorAssistantScheduleModel{
		ID:                helpers.StringValueOrNull(schedule.ID),
		CustomerID:        helpers.StringValueOrNull(schedule.CustomerID),
		Enabled:           types.BoolValue(schedule.Enabled),
		DeleteDisabled:    types.BoolValue(schedule.DeleteDisabled),
		Frequency:         helpers.StringValueOrNull(schedule.Frequency),
		FrequencyInterval: helpers.StringValueOrNull(schedule.FrequencyInterval),
	}

	return model, diag.Diagnostics{}
}

var pathRootID = path.Root("id")
