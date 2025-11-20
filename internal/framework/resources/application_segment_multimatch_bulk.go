package resources

import (
	"context"
	"fmt"
	"sort"
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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &ApplicationSegmentMultimatchBulkResource{}
	_ resource.ResourceWithConfigure   = &ApplicationSegmentMultimatchBulkResource{}
	_ resource.ResourceWithImportState = &ApplicationSegmentMultimatchBulkResource{}
)

func NewApplicationSegmentMultimatchBulkResource() resource.Resource {
	return &ApplicationSegmentMultimatchBulkResource{}
}

type ApplicationSegmentMultimatchBulkResource struct {
	client *client.Client
}

type ApplicationSegmentMultimatchBulkModel struct {
	ID             types.String `tfsdk:"id"`
	ApplicationIDs types.Set    `tfsdk:"application_ids"`
	MatchStyle     types.String `tfsdk:"match_style"`
	MicroTenantID  types.String `tfsdk:"microtenant_id"`
}

func (r *ApplicationSegmentMultimatchBulkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_segment_multimatch_bulk"
}

func (r *ApplicationSegmentMultimatchBulkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages bulk updates of application segment match styles.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"application_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "Set of application segment IDs to update.",
			},
			"match_style": schema.StringAttribute{
				Required:    true,
				Description: "Match style applied to the specified application segments.",
				Validators: []validator.String{
					stringvalidator.OneOf("EXCLUSIVE", "INCLUSIVE"),
				},
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Description: "Micro-tenant identifier used to scope API calls.",
			},
		},
	}
}

func (r *ApplicationSegmentMultimatchBulkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationSegmentMultimatchBulkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before creating resources.")
		return
	}

	var plan ApplicationSegmentMultimatchBulkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicroTenantID)

	applicationIDs, diags := helpers.SetValueToStringSlice(ctx, plan.ApplicationIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(applicationIDs) == 0 {
		resp.Diagnostics.AddError("Validation Error", "At least one application_id must be provided.")
		return
	}

	intIDs, intDiags := helpers.StringSliceToIntSlice(applicationIDs)
	resp.Diagnostics.Append(intDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	matchStyle := strings.TrimSpace(plan.MatchStyle.ValueString())
	if matchStyle == "" {
		resp.Diagnostics.AddError("Validation Error", "match_style must be provided.")
		return
	}

	payload := applicationsegment.BulkUpdateMultiMatchPayload{
		ApplicationIDs: intIDs,
		MatchStyle:     matchStyle,
	}

	tflog.Info(ctx, "Updating application segments match style", map[string]any{"count": len(intIDs), "match_style": matchStyle})
	if _, err := applicationsegment.UpdatebulkUpdateMultiMatch(ctx, service, payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update application segments: %v", err))
		return
	}

	syntheticID := r.generateID(applicationIDs, matchStyle)
	plan.ID = types.StringValue(syntheticID)

	state, diags := r.refreshState(ctx, service, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationSegmentMultimatchBulkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before reading resources.")
		return
	}

	var state ApplicationSegmentMultimatchBulkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicroTenantID)

	newState, diags := r.refreshState(ctx, service, state)
	if diags.HasError() {
		for _, diagnostic := range diags {
			if diagnostic.Severity() == diag.SeverityError && strings.Contains(strings.ToLower(diagnostic.Detail()), "not found") {
				resp.State.RemoveResource(ctx)
				return
			}
		}
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ApplicationSegmentMultimatchBulkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before updating resources.")
		return
	}

	var plan ApplicationSegmentMultimatchBulkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicroTenantID)

	applicationIDs, diags := helpers.SetValueToStringSlice(ctx, plan.ApplicationIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(applicationIDs) == 0 {
		resp.Diagnostics.AddError("Validation Error", "At least one application_id must be provided.")
		return
	}

	intIDs, intDiags := helpers.StringSliceToIntSlice(applicationIDs)
	resp.Diagnostics.Append(intDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	matchStyle := strings.TrimSpace(plan.MatchStyle.ValueString())
	if matchStyle == "" {
		resp.Diagnostics.AddError("Validation Error", "match_style must be provided.")
		return
	}

	payload := applicationsegment.BulkUpdateMultiMatchPayload{
		ApplicationIDs: intIDs,
		MatchStyle:     matchStyle,
	}

	tflog.Info(ctx, "Updating application segments match style", map[string]any{"count": len(intIDs), "match_style": matchStyle})
	if _, err := applicationsegment.UpdatebulkUpdateMultiMatch(ctx, service, payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update application segments: %v", err))
		return
	}

	plan.ID = types.StringValue(r.generateID(applicationIDs, matchStyle))

	state, refreshDiags := r.refreshState(ctx, service, plan)
	resp.Diagnostics.Append(refreshDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationSegmentMultimatchBulkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}

func (r *ApplicationSegmentMultimatchBulkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
}

func (r *ApplicationSegmentMultimatchBulkResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && !microtenantID.IsUnknown() && microtenantID.ValueString() != "" {
		service = service.WithMicroTenant(microtenantID.ValueString())
	}
	return service
}

func (r *ApplicationSegmentMultimatchBulkResource) refreshState(ctx context.Context, service *zscaler.Service, model ApplicationSegmentMultimatchBulkModel) (ApplicationSegmentMultimatchBulkModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	applicationIDs, idsDiags := helpers.SetValueToStringSlice(ctx, model.ApplicationIDs)
	diags.Append(idsDiags...)
	if diags.HasError() {
		return model, diags
	}
	if len(applicationIDs) == 0 {
		diags.AddError("Validation Error", "Resource has no application_ids to manage.")
		return model, diags
	}

	actualMatchStyle, fetchDiags := r.determineMatchStyle(ctx, service, applicationIDs)
	diags.Append(fetchDiags...)
	if diags.HasError() {
		return model, diags
	}

	sortedIDs := make([]string, len(applicationIDs))
	copy(sortedIDs, applicationIDs)
	sort.Strings(sortedIDs)

	idsSet, setDiags := types.SetValueFrom(ctx, types.StringType, sortedIDs)
	diags.Append(setDiags...)
	if diags.HasError() {
		return model, diags
	}

	model.ApplicationIDs = idsSet
	model.MatchStyle = helpers.StringValueOrNull(actualMatchStyle)
	model.ID = types.StringValue(r.generateID(sortedIDs, actualMatchStyle))

	return model, diags
}

func (r *ApplicationSegmentMultimatchBulkResource) determineMatchStyle(ctx context.Context, service *zscaler.Service, applicationIDs []string) (string, diag.Diagnostics) {
	var diags diag.Diagnostics
	matchStyle := ""

	for _, id := range applicationIDs {
		segment, _, err := applicationsegment.Get(ctx, service, id)
		if err != nil {
			if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
				diags.AddError("Not Found", fmt.Sprintf("Application segment %s was not found", id))
				return "", diags
			}
			diags.AddError("Client Error", fmt.Sprintf("Failed to read application segment %s: %v", id, err))
			return "", diags
		}

		tflog.Debug(ctx, "Fetched application segment", map[string]any{"id": id, "match_style": segment.MatchStyle})
		if matchStyle == "" {
			matchStyle = segment.MatchStyle
		} else if segment.MatchStyle != matchStyle {
			tflog.Warn(ctx, "Application segments have differing match styles", map[string]any{"segment_id": id, "match_style": segment.MatchStyle, "expected": matchStyle})
		}
	}

	return matchStyle, diags
}

func (r *ApplicationSegmentMultimatchBulkResource) generateID(applicationIDs []string, matchStyle string) string {
	sorted := make([]string, len(applicationIDs))
	copy(sorted, applicationIDs)
	sort.Strings(sorted)
	key := fmt.Sprintf("%s|%s", strings.Join(sorted, ","), strings.TrimSpace(matchStyle))
	return helpers.GenerateShortID(key)
}
