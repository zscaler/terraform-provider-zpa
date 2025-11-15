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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/bacertificate"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &BaCertificateResource{}
	_ resource.ResourceWithConfigure   = &BaCertificateResource{}
	_ resource.ResourceWithImportState = &BaCertificateResource{}
)

func NewBaCertificateResource() resource.Resource {
	return &BaCertificateResource{}
}

type BaCertificateResource struct {
	client *client.Client
}

type BaCertificateModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CertBlob      types.String `tfsdk:"cert_blob"`
	Certificate   types.String `tfsdk:"certificate"`
	MicrotenantID types.String `tfsdk:"microtenant_id"`
}

func (r *BaCertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ba_certificate"
}

func (r *BaCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZPA Browser Access certificate.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "Unique identifier of the certificate.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the certificate.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the certificate.",
			},
			"cert_blob": schema.StringAttribute{
				Optional:    true,
				Description: "Certificate blob content in PEM format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate": schema.StringAttribute{
				Computed:    true,
				Description: "Certificate text returned by the API in PEM format.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Microtenant identifier used to scope API calls.",
			},
		},
	}
}

func (r *BaCertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BaCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan BaCertificateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload := expandBaCertificate(plan)

	created, _, err := bacertificate.Create(ctx, service, payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create BA certificate: %v", err))
		return
	}

	state, diags := r.readIntoState(ctx, service, created.ID, plan.CertBlob)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BaCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state BaCertificateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	newState, diags := r.readIntoState(ctx, service, state.ID.ValueString(), state.CertBlob)
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

func (r *BaCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	// No update operation is supported by the API; retain existing state.
	var state BaCertificateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BaCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state BaCertificateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	if _, err := bacertificate.Delete(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete BA certificate: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *BaCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *BaCertificateResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && !microtenantID.IsUnknown() {
		trimmed := strings.TrimSpace(microtenantID.ValueString())
		if trimmed != "" {
			service = service.WithMicroTenant(trimmed)
		}
	}
	return service
}

func (r *BaCertificateResource) readIntoState(ctx context.Context, service *zscaler.Service, id string, certBlob types.String) (BaCertificateModel, diag.Diagnostics) {
	cert, _, err := bacertificate.Get(ctx, service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return BaCertificateModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("BA certificate %s not found", id))}
		}
		return BaCertificateModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read BA certificate: %v", err))}
	}

	state := BaCertificateModel{
		ID:            helpers.StringValueOrNull(cert.ID),
		Name:          helpers.StringValueOrNull(cert.Name),
		Description:   helpers.StringValueOrNull(cert.Description),
		Certificate:   helpers.StringValueOrNull(cert.Certificate),
		MicrotenantID: helpers.StringValueOrNull(cert.MicrotenantID),
		CertBlob:      certBlob,
	}

	return state, nil
}

func expandBaCertificate(plan BaCertificateModel) bacertificate.BaCertificate {
	name := strings.TrimSpace(plan.Name.ValueString())
	description := ""
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		description = strings.TrimSpace(plan.Description.ValueString())
	}
	certBlob := ""
	if !plan.CertBlob.IsNull() && !plan.CertBlob.IsUnknown() {
		certBlob = strings.TrimSpace(plan.CertBlob.ValueString())
	}
	microtenantID := ""
	if !plan.MicrotenantID.IsNull() && !plan.MicrotenantID.IsUnknown() {
		microtenantID = strings.TrimSpace(plan.MicrotenantID.ValueString())
	}

	return bacertificate.BaCertificate{
		Name:          name,
		Description:   description,
		CertBlob:      certBlob,
		MicrotenantID: microtenantID,
	}
}
