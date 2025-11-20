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
	cbicertificatecontroller "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbicertificatecontroller"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &CBICertificateResource{}
	_ resource.ResourceWithConfigure   = &CBICertificateResource{}
	_ resource.ResourceWithImportState = &CBICertificateResource{}
)

func NewCBICertificateResource() resource.Resource {
	return &CBICertificateResource{}
}

type CBICertificateResource struct {
	client *client.Client
}

type CBICertificateModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	PEM  types.String `tfsdk:"pem"`
}

func (r *CBICertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_browser_isolation_certificate"
}

func (r *CBICertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Cloud Browser Isolation certificate.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the certificate.",
			},
			"pem": schema.StringAttribute{
				Optional:    true,
				Description: "Certificate body in PEM format.",
			},
		},
	}
}

func (r *CBICertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CBICertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan CBICertificateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := expandCBICertificate(plan)

	created, _, err := cbicertificatecontroller.Create(ctx, r.client.Service, &payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create CBI certificate: %v", err))
		return
	}

	state, diags := r.readIntoState(ctx, created.ID, plan.PEM)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CBICertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state CBICertificateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readIntoState(ctx, state.ID.ValueString(), state.PEM)
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

func (r *CBICertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan CBICertificateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.IsUnknown() {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	payload := expandCBICertificate(plan)

	if _, err := cbicertificatecontroller.Update(ctx, r.client.Service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update CBI certificate: %v", err))
		return
	}

	state, diags := r.readIntoState(ctx, plan.ID.ValueString(), plan.PEM)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CBICertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state CBICertificateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := cbicertificatecontroller.Delete(ctx, r.client.Service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete CBI certificate: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *CBICertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
}

func (r *CBICertificateResource) readIntoState(ctx context.Context, id string, pem types.String) (CBICertificateModel, diag.Diagnostics) {
	cert, _, err := cbicertificatecontroller.Get(ctx, r.client.Service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return CBICertificateModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("CBI certificate %s not found", id))}
		}
		return CBICertificateModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read CBI certificate: %v", err))}
	}

	state := CBICertificateModel{
		ID:   helpers.StringValueOrNull(cert.ID),
		Name: helpers.StringValueOrNull(cert.Name),
		PEM:  pem, // Always preserve the provided PEM value (sensitive field, API may return different formatting)
	}

	return state, nil
}

func expandCBICertificate(plan CBICertificateModel) cbicertificatecontroller.CBICertificate {
	return cbicertificatecontroller.CBICertificate{
		ID:   helpers.StringValue(plan.ID),
		Name: helpers.StringValue(plan.Name),
		PEM:  helpers.StringValue(plan.PEM),
	}
}
