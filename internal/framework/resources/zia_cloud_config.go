package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/custom_config_controller"
)

var (
	_ resource.Resource                = &ZIACloudConfigResource{}
	_ resource.ResourceWithConfigure   = &ZIACloudConfigResource{}
	_ resource.ResourceWithImportState = &ZIACloudConfigResource{}
)

func NewZIACloudConfigResource() resource.Resource {
	return &ZIACloudConfigResource{}
}

type ZIACloudConfigResource struct {
	client *client.Client
}

type ZIACloudConfigModel struct {
	ID                    types.String `tfsdk:"id"`
	ZIACloudDomain        types.String `tfsdk:"zia_cloud_domain"`
	ZIAUsername           types.String `tfsdk:"zia_username"`
	ZIAPassword           types.String `tfsdk:"zia_password"`
	ZIASandboxApiToken    types.String `tfsdk:"zia_sandbox_api_token"`
	ZIACloudServiceApiKey types.String `tfsdk:"zia_cloud_service_api_key"`
}

func (r *ZIACloudConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zia_cloud_config"
}

func (r *ZIACloudConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages ZIA (Zscaler Internet Access) cloud configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"zia_cloud_domain": schema.StringAttribute{
				Required:    true,
				Description: "ZIA cloud domain (without .net suffix). Valid values: zscaler, zscloud, zscalerone, zscalertwo, zscalerthree, zscalerbeta, zscalergov, zscalerten, zspreview",
				Validators: []validator.String{
					stringvalidator.OneOf("zscaler", "zscloud", "zscalerone", "zscalertwo", "zscalerthree", "zscalerbeta", "zscalergov", "zscalerten", "zspreview"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					&normalizeDomainPlanModifier{},
				},
			},
			"zia_username": schema.StringAttribute{
				Required:    true,
				Description: "ZIA username",
			},
			"zia_password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "ZIA password (write-only, not returned by API)",
			},
			"zia_sandbox_api_token": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "ZIA sandbox API token (write-only, not returned by API)",
			},
			"zia_cloud_service_api_key": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "ZIA cloud service API key (write-only, not returned by API)",
			},
		},
	}
}

type normalizeDomainPlanModifier struct{}

func (m *normalizeDomainPlanModifier) Description(ctx context.Context) string {
	return "Normalizes the domain by ensuring it has .net suffix"
}

func (m *normalizeDomainPlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Normalizes the domain by ensuring it has .net suffix"
}

func (m *normalizeDomainPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	domain := req.ConfigValue.ValueString()
	if domain != "" && !strings.HasSuffix(domain, ".net") {
		resp.PlanValue = types.StringValue(domain + ".net")
	} else {
		resp.PlanValue = req.ConfigValue
	}
}

func (r *ZIACloudConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ZIACloudConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan ZIACloudConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := plan.ZIACloudDomain.ValueString()
	if !strings.HasSuffix(domain, ".net") {
		domain = domain + ".net"
	}

	cloudConfig := custom_config_controller.ZIACloudConfig{
		ZIACloudDomain:        domain,
		ZIAUsername:           plan.ZIAUsername.ValueString(),
		ZIAPassword:           plan.ZIAPassword.ValueString(),
		ZIASandboxApiToken:    plan.ZIASandboxApiToken.ValueString(),
		ZIACloudServiceApiKey: plan.ZIACloudServiceApiKey.ValueString(),
	}

	if _, _, err := custom_config_controller.AddZIACloudConfig(ctx, r.client.Service, &cloudConfig); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create ZIA cloud config: %v", err))
		return
	}

	plan.ID = types.StringValue("zia_cloud_config")

	state, readDiags := r.readZIACloudConfig(ctx)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = plan.ID
	state.ZIAPassword = plan.ZIAPassword
	state.ZIASandboxApiToken = plan.ZIASandboxApiToken
	state.ZIACloudServiceApiKey = plan.ZIACloudServiceApiKey
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ZIACloudConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state ZIACloudConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readZIACloudConfig(ctx)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	newState.ID = state.ID
	newState.ZIAPassword = state.ZIAPassword
	newState.ZIASandboxApiToken = state.ZIASandboxApiToken
	newState.ZIACloudServiceApiKey = state.ZIACloudServiceApiKey
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ZIACloudConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan ZIACloudConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := plan.ZIACloudDomain.ValueString()
	if !strings.HasSuffix(domain, ".net") {
		domain = domain + ".net"
	}

	cloudConfig := custom_config_controller.ZIACloudConfig{
		ZIACloudDomain:        domain,
		ZIAUsername:           plan.ZIAUsername.ValueString(),
		ZIAPassword:           plan.ZIAPassword.ValueString(),
		ZIASandboxApiToken:    plan.ZIASandboxApiToken.ValueString(),
		ZIACloudServiceApiKey: plan.ZIACloudServiceApiKey.ValueString(),
	}

	if _, _, err := custom_config_controller.AddZIACloudConfig(ctx, r.client.Service, &cloudConfig); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update ZIA cloud config: %v", err))
		return
	}

	plan.ID = types.StringValue("zia_cloud_config")

	state, readDiags := r.readZIACloudConfig(ctx)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = plan.ID
	state.ZIAPassword = plan.ZIAPassword
	state.ZIASandboxApiToken = plan.ZIASandboxApiToken
	state.ZIACloudServiceApiKey = plan.ZIACloudServiceApiKey
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ZIACloudConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Delete is a no-op as per SDKv2 implementation
}

func (r *ZIACloudConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before importing ZIA cloud config.")
		return
	}

	state, diags := r.readZIACloudConfig(ctx)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	state.ID = types.StringValue("zia_cloud_config")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ZIACloudConfigResource) readZIACloudConfig(ctx context.Context) (ZIACloudConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	config, _, err := custom_config_controller.GetZIACloudConfig(ctx, r.client.Service)
	if err != nil {
		return ZIACloudConfigModel{}, diag.Diagnostics{
			diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read ZIA cloud config: %v", err)),
		}
	}

	if config == nil {
		return ZIACloudConfigModel{}, diag.Diagnostics{
			diag.NewErrorDiagnostic("Not Found", "Couldn't read ZIA cloud config"),
		}
	}

	return ZIACloudConfigModel{
		ID:             types.StringValue("zia_cloud_config"),
		ZIACloudDomain: helpers.StringValueOrNull(config.ZIACloudDomain),
		ZIAUsername:    helpers.StringValueOrNull(config.ZIAUsername),
		// Password, sandbox token, and API key are write-only and not returned by API
		ZIAPassword:           types.StringNull(),
		ZIASandboxApiToken:    types.StringNull(),
		ZIACloudServiceApiKey: types.StringNull(),
	}, diags
}
