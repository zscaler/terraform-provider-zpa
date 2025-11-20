package resources

import (
	"context"
	"fmt"
	"strconv"
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
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praapproval"
)

var (
	_ resource.Resource                = &PRAApprovalResource{}
	_ resource.ResourceWithConfigure   = &PRAApprovalResource{}
	_ resource.ResourceWithImportState = &PRAApprovalResource{}
)

func NewPRAApprovalResource() resource.Resource {
	return &PRAApprovalResource{}
}

type PRAApprovalResource struct {
	client *client.Client
}

type PRAApprovalModel struct {
	ID            types.String             `tfsdk:"id"`
	EmailIDs      types.Set                `tfsdk:"email_ids"`
	StartTime     types.String             `tfsdk:"start_time"`
	EndTime       types.String             `tfsdk:"end_time"`
	Status        types.String             `tfsdk:"status"`
	WorkingHours  []PRAWorkingHoursModel   `tfsdk:"working_hours"`
	Applications  []PRAApprovalApplication `tfsdk:"applications"`
	MicrotenantID types.String             `tfsdk:"microtenant_id"`
}

type PRAWorkingHoursModel struct {
	Days          types.Set    `tfsdk:"days"`
	StartTime     types.String `tfsdk:"start_time"`
	StartTimeCron types.String `tfsdk:"start_time_cron"`
	EndTime       types.String `tfsdk:"end_time"`
	EndTimeCron   types.String `tfsdk:"end_time_cron"`
	Timezone      types.String `tfsdk:"timezone"`
}

type PRAApprovalApplication struct {
	IDs types.List `tfsdk:"id"`
}

func (r *PRAApprovalResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pra_approval_controller"
}

func (r *PRAApprovalResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZPA PRA privileged approval.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "Email IDs associated with the approval.",
			},
			"start_time": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Approval start time in RFC1123 format.",
			},
			"end_time": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Approval end time in RFC1123 format.",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Approval status. Supported values: INVALID, ACTIVE, FUTURE, EXPIRED.",
				Validators: []validator.String{
					stringvalidator.OneOf("INVALID", "ACTIVE", "FUTURE", "EXPIRED"),
				},
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID for scoping.",
			},
		},
		Blocks: map[string]schema.Block{
			// working_hours: TypeSet in SDKv2
			// Using SetNestedBlock for block syntax support
			"working_hours": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"days": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,
						},
						"start_time": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"start_time_cron": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"end_time": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"end_time_cron": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"timezone": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			// applications: TypeSet in SDKv2, id is TypeList (Optional)
			// Using SetNestedBlock for block syntax support
			"applications": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *PRAApprovalResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *PRAApprovalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing PRA approvals.")
		return
	}

	var plan PRAApprovalModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload, diags := expandPRAApproval(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := praapproval.Create(ctx, service, &payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create PRA approval: %v", err))
		return
	}

	state, readDiags := r.readApproval(ctx, service, created.ID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PRAApprovalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing PRA approvals.")
		return
	}

	var state PRAApprovalModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	newState, diags := r.readApproval(ctx, service, state.ID.ValueString())
	if diags.HasError() {
		for _, d := range diags {
			if d.Severity() == diag.SeverityError && strings.Contains(strings.ToLower(d.Detail()), "not found") {
				resp.State.RemoveResource(ctx)
				return
			}
		}
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *PRAApprovalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing PRA approvals.")
		return
	}

	var plan PRAApprovalModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing ID", "id must be known during update.")
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload, diags := expandPRAApproval(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := praapproval.Update(ctx, service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update PRA approval: %v", err))
		return
	}

	state, readDiags := r.readApproval(ctx, service, plan.ID.ValueString())
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PRAApprovalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing PRA approvals.")
		return
	}

	var state PRAApprovalModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	if _, err := praapproval.Delete(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete PRA approval: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *PRAApprovalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before importing PRA approvals.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the PRA approval ID or email address.")
		return
	}

	service := r.client.Service
	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		resource, _, lookupErr := praapproval.GetByEmailID(ctx, service, id)
		if lookupErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate PRA approval %q: %v", id, lookupErr))
			return
		}
		id = resource.ID
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
}

func (r *PRAApprovalResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && microtenantID.ValueString() != "" {
		service = service.WithMicroTenant(strings.TrimSpace(microtenantID.ValueString()))
	}
	return service
}

func (r *PRAApprovalResource) readApproval(ctx context.Context, service *zscaler.Service, id string) (PRAApprovalModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	result, _, err := praapproval.Get(ctx, service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return PRAApprovalModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("PRA approval %s not found", id))}
		}
		return PRAApprovalModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read PRA approval: %v", err))}
	}

	model, flattenDiags := flattenPRAApproval(ctx, result)
	diags.Append(flattenDiags...)
	return model, diags
}

func expandPRAApproval(ctx context.Context, model *PRAApprovalModel) (praapproval.PrivilegedApproval, diag.Diagnostics) {
	var diags diag.Diagnostics

	emails, emailDiags := helpers.SetValueToStringSlice(ctx, model.EmailIDs)
	diags.Append(emailDiags...)

	applications, appsDiags := expandPRAApprovalApplications(ctx, model.Applications)
	diags.Append(appsDiags...)

	var workingHours *praapproval.WorkingHours
	if len(model.WorkingHours) > 0 {
		var workingDiag diag.Diagnostics
		workingHours, workingDiag = expandPRAWorkingHours(ctx, model.WorkingHours[0])
		diags.Append(workingDiag...)
	}

	start := strings.TrimSpace(model.StartTime.ValueString())
	end := strings.TrimSpace(model.EndTime.ValueString())

	if start != "" && end != "" {
		if err := helpers.ValidatePRATimeRange(start, end); err != nil {
			diags.AddError("Validation Error", err.Error())
		}
	}

	var startEpoch, endEpoch string
	if start != "" {
		if epoch, err := helpers.ConvertRFC1123ToEpoch(start); err != nil {
			diags.AddError("Validation Error", fmt.Sprintf("start_time conversion error: %v", err))
		} else {
			startEpoch = fmt.Sprintf("%d", epoch)
		}
	}
	if end != "" {
		if epoch, err := helpers.ConvertRFC1123ToEpoch(end); err != nil {
			diags.AddError("Validation Error", fmt.Sprintf("end_time conversion error: %v", err))
		} else {
			endEpoch = fmt.Sprintf("%d", epoch)
		}
	}

	result := praapproval.PrivilegedApproval{
		ID:            model.ID.ValueString(),
		EmailIDs:      emails,
		StartTime:     startEpoch,
		EndTime:       endEpoch,
		Status:        model.Status.ValueString(),
		WorkingHours:  workingHours,
		Applications:  applications,
		MicroTenantID: model.MicrotenantID.ValueString(),
	}

	return result, diags
}

func expandPRAApprovalApplications(ctx context.Context, models []PRAApprovalApplication) ([]praapproval.Applications, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(models) == 0 {
		diags.AddError("Validation Error", "applications must be provided")
		return nil, diags
	}

	applications := make([]praapproval.Applications, 0)
	for _, model := range models {
		ids, idsDiags := helpers.ListValueToStringSlice(ctx, model.IDs)
		diags.Append(idsDiags...)
		for _, id := range ids {
			if strings.TrimSpace(id) == "" {
				continue
			}
			applications = append(applications, praapproval.Applications{ID: id})
		}
	}

	if len(applications) == 0 {
		diags.AddError("Validation Error", "applications.id must contain at least one value")
	}

	return applications, diags
}

func expandPRAWorkingHours(ctx context.Context, model PRAWorkingHoursModel) (*praapproval.WorkingHours, diag.Diagnostics) {
	var diags diag.Diagnostics

	if model.Days.IsNull() && model.StartTime.IsNull() && model.EndTime.IsNull() &&
		model.StartTimeCron.IsNull() && model.EndTimeCron.IsNull() && model.Timezone.IsNull() {
		return nil, diags
	}

	days, dayDiags := helpers.SetValueToStringSlice(ctx, model.Days)
	diags.Append(dayDiags...)

	start := strings.TrimSpace(model.StartTime.ValueString())
	end := strings.TrimSpace(model.EndTime.ValueString())
	startCron := strings.TrimSpace(model.StartTimeCron.ValueString())
	endCron := strings.TrimSpace(model.EndTimeCron.ValueString())
	timezone := strings.TrimSpace(model.Timezone.ValueString())

	if start != "" {
		if err := helpers.Validate24HourTimeFormat(start); err != nil {
			diags.AddError("Validation Error", fmt.Sprintf("working_hours.start_time: %v", err))
		}
	}
	if end != "" {
		if err := helpers.Validate24HourTimeFormat(end); err != nil {
			diags.AddError("Validation Error", fmt.Sprintf("working_hours.end_time: %v", err))
		}
	}
	if timezone != "" {
		if err := helpers.ValidateTimeZone(timezone); err != nil {
			diags.AddError("Validation Error", fmt.Sprintf("working_hours.timezone: %v", err))
		}
	}

	return &praapproval.WorkingHours{
		Days:          days,
		StartTime:     start,
		EndTime:       end,
		StartTimeCron: startCron,
		EndTimeCron:   endCron,
		TimeZone:      timezone,
	}, diags
}

func flattenPRAApproval(ctx context.Context, approval *praapproval.PrivilegedApproval) (PRAApprovalModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	emailSet, emailDiags := types.SetValueFrom(ctx, types.StringType, approval.EmailIDs)
	diags.Append(emailDiags...)

	startRFC, endRFC := types.StringNull(), types.StringNull()
	if strings.TrimSpace(approval.StartTime) != "" {
		if value, err := helpers.RFC1123FromEpoch(approval.StartTime); err == nil {
			startRFC = types.StringValue(value)
		} else {
			diags.AddError("Conversion Error", fmt.Sprintf("Unable to convert start_time: %v", err))
		}
	}
	if strings.TrimSpace(approval.EndTime) != "" {
		if value, err := helpers.RFC1123FromEpoch(approval.EndTime); err == nil {
			endRFC = types.StringValue(value)
		} else {
			diags.AddError("Conversion Error", fmt.Sprintf("Unable to convert end_time: %v", err))
		}
	}

	apps, appsDiags := flattenPRAApprovalApplications(ctx, approval.Applications)
	diags.Append(appsDiags...)

	var workingHours []PRAWorkingHoursModel
	if approval.WorkingHours != nil {
		wh, whDiags := flattenPRAWorkingHours(ctx, approval.WorkingHours)
		diags.Append(whDiags...)
		workingHours = wh
	}

	return PRAApprovalModel{
		ID:            types.StringValue(approval.ID),
		EmailIDs:      emailSet,
		StartTime:     startRFC,
		EndTime:       endRFC,
		Status:        types.StringValue(approval.Status),
		WorkingHours:  workingHours,
		Applications:  apps,
		MicrotenantID: types.StringValue(approval.MicroTenantID),
	}, diags
}

func flattenPRAApprovalApplications(ctx context.Context, applications []praapproval.Applications) ([]PRAApprovalApplication, diag.Diagnostics) {
	if len(applications) == 0 {
		return nil, nil
	}

	ids := make([]string, 0, len(applications))
	for _, app := range applications {
		if strings.TrimSpace(app.ID) != "" {
			ids = append(ids, app.ID)
		}
	}

	listValue, diags := types.ListValueFrom(ctx, types.StringType, ids)
	if diags.HasError() {
		return nil, diags
	}

	return []PRAApprovalApplication{{IDs: listValue}}, nil
}

func flattenPRAWorkingHours(ctx context.Context, hours *praapproval.WorkingHours) ([]PRAWorkingHoursModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	if hours == nil {
		return nil, diags
	}

	days, dayDiags := types.SetValueFrom(ctx, types.StringType, hours.Days)
	diags.Append(dayDiags...)

	model := PRAWorkingHoursModel{
		Days:          days,
		StartTime:     types.StringValue(hours.StartTime),
		EndTime:       types.StringValue(hours.EndTime),
		StartTimeCron: types.StringValue(hours.StartTimeCron),
		EndTimeCron:   types.StringValue(hours.EndTimeCron),
		Timezone:      types.StringValue(hours.TimeZone),
	}

	return []PRAWorkingHoursModel{model}, diags
}
