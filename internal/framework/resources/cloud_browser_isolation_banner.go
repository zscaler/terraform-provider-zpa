package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	cbibannercontroller "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbibannercontroller"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &CBIBannerResource{}
	_ resource.ResourceWithConfigure   = &CBIBannerResource{}
	_ resource.ResourceWithImportState = &CBIBannerResource{}
)

func NewCBIBannerResource() resource.Resource {
	return &CBIBannerResource{}
}

type CBIBannerResource struct {
	client *client.Client
}

type CBIBannerModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	PrimaryColor      types.String `tfsdk:"primary_color"`
	TextColor         types.String `tfsdk:"text_color"`
	NotificationTitle types.String `tfsdk:"notification_title"`
	NotificationText  types.String `tfsdk:"notification_text"`
	Logo              types.String `tfsdk:"logo"`
	Banner            types.Bool   `tfsdk:"banner"`
	Persist           types.Bool   `tfsdk:"persist"`
}

func (r *CBIBannerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_browser_isolation_banner"
}

func (r *CBIBannerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Cloud Browser Isolation banner.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the banner.",
			},
			"primary_color": schema.StringAttribute{
				Optional:    true,
				Description: "Primary color of the banner.",
			},
			"text_color": schema.StringAttribute{
				Optional:    true,
				Description: "Text color of the banner.",
			},
			"notification_title": schema.StringAttribute{
				Optional:    true,
				Description: "Title displayed on the banner notification.",
			},
			"notification_text": schema.StringAttribute{
				Optional:    true,
				Description: "Body text displayed on the banner notification.",
			},
			"logo": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Logo displayed on the banner (base64 encoded).",
			},
			"banner": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Indicates if the banner is enabled.",
			},
			"persist": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Indicates if the banner should persist for the user.",
			},
		},
	}
}

func (r *CBIBannerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CBIBannerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan CBIBannerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := expandCBIBanner(plan)

	created, _, err := cbibannercontroller.Create(ctx, r.client.Service, &payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create CBI banner: %v", err))
		return
	}

	state, diags := r.readIntoState(ctx, created.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CBIBannerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state CBIBannerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readIntoState(ctx, state.ID.ValueString())
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

func (r *CBIBannerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan CBIBannerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.IsUnknown() {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	payload := expandCBIBanner(plan)

	if _, err := cbibannercontroller.Update(ctx, r.client.Service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update CBI banner: %v", err))
		return
	}

	state, diags := r.readIntoState(ctx, plan.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CBIBannerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state CBIBannerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := cbibannercontroller.Delete(ctx, r.client.Service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete CBI banner: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *CBIBannerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
}

func (r *CBIBannerResource) readIntoState(ctx context.Context, id string) (CBIBannerModel, diag.Diagnostics) {
	banner, _, err := cbibannercontroller.Get(ctx, r.client.Service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return CBIBannerModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("CBI banner %s not found", id))}
		}
		return CBIBannerModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read CBI banner: %v", err))}
	}

	state := CBIBannerModel{
		ID:                helpers.StringValueOrNull(banner.ID),
		Name:              helpers.StringValueOrNull(banner.Name),
		PrimaryColor:      helpers.StringValueOrNull(banner.PrimaryColor),
		TextColor:         helpers.StringValueOrNull(banner.TextColor),
		NotificationTitle: helpers.StringValueOrNull(banner.NotificationTitle),
		NotificationText:  helpers.StringValueOrNull(banner.NotificationText),
		Logo:              helpers.StringValueOrNull(banner.Logo),
		Banner:            types.BoolValue(banner.Banner),
		Persist:           types.BoolValue(banner.Persist),
	}

	return state, nil
}

func expandCBIBanner(plan CBIBannerModel) cbibannercontroller.CBIBannerController {
	return cbibannercontroller.CBIBannerController{
		ID:                helpers.StringValue(plan.ID),
		Name:              helpers.StringValue(plan.Name),
		PrimaryColor:      helpers.StringValue(plan.PrimaryColor),
		TextColor:         helpers.StringValue(plan.TextColor),
		NotificationTitle: helpers.StringValue(plan.NotificationTitle),
		NotificationText:  helpers.StringValue(plan.NotificationText),
		Logo:              helpers.StringValue(plan.Logo),
		Banner:            helpers.BoolValue(plan.Banner, false),
		Persist:           helpers.BoolValue(plan.Persist, false),
	}
}
