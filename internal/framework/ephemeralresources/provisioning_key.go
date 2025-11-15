package ephemeralresources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	ephschema "github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/provisioningkey"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
)

var (
	_ ephemeral.EphemeralResource              = &ProvisioningKeyEphemeralResource{}
	_ ephemeral.EphemeralResourceWithConfigure = &ProvisioningKeyEphemeralResource{}
)

func NewProvisioningKeyEphemeralResource() ephemeral.EphemeralResource {
	return &ProvisioningKeyEphemeralResource{}
}

type ProvisioningKeyEphemeralResource struct {
	client *client.Client
}

type ProvisioningKeyEphemeralModel struct {
	ID              types.String `tfsdk:"id"`
	AssociationType types.String `tfsdk:"association_type"`
	MicroTenantID   types.String `tfsdk:"microtenant_id"`
	ProvisioningKey types.String `tfsdk:"provisioning_key"`
}

func (r *ProvisioningKeyEphemeralResource) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provisioning_key"
}

func (r *ProvisioningKeyEphemeralResource) Schema(ctx context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = ephschema.Schema{
		Description: "Retrieves a ZPA provisioning key value without persisting it in Terraform state.",
		Attributes: map[string]ephschema.Attribute{
			"id": ephschema.StringAttribute{
				Required:    true,
				Description: "The identifier of the provisioning key.",
			},
			"association_type": ephschema.StringAttribute{
				Required:    true,
				Description: "Provisioning key association type. Valid values match the provisioning key resource.",
				Validators: []validator.String{
					stringvalidator.OneOf(provisioningkey.ProvisioningKeyAssociationTypes...),
				},
			},
			"microtenant_id": ephschema.StringAttribute{
				Optional: true,
				// Computed:    true,
				Description: "Optional micro-tenant ID used to scope the lookup.",
			},
			"provisioning_key": ephschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "The provisioning key value returned by the API.",
			},
		},
	}
}

func (r *ProvisioningKeyEphemeralResource) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cl, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = cl
}

func (r *ProvisioningKeyEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Ensure the provider block is configured before using the provisioning key ephemeral resource.",
		)
		return
	}

	var config ProvisioningKeyEphemeralModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(config.ID.ValueString())
	associationType := strings.TrimSpace(config.AssociationType.ValueString())
	if id == "" || associationType == "" {
		resp.Diagnostics.AddError("Missing Required Attributes", "Both id and association_type must be provided.")
		return
	}

	service := r.client.Service
	if !config.MicroTenantID.IsNull() && !config.MicroTenantID.IsUnknown() {
		microtenantID := strings.TrimSpace(config.MicroTenantID.ValueString())
		if microtenantID != "" {
			service = service.WithMicroTenant(microtenantID)
		}
	}

	tflog.Debug(ctx, "Fetching provisioning key via ephemeral resource", map[string]any{
		"id":               id,
		"association_type": associationType,
	})

	key, _, err := provisioningkey.Get(ctx, service, associationType, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			resp.Diagnostics.AddError(
				"Provisioning Key Not Found",
				fmt.Sprintf("No provisioning key found with id %q and association type %q.", id, associationType),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read provisioning key: %v", err),
		)
		return
	}

	if key == nil {
		resp.Diagnostics.AddError(
			"Empty Response",
			fmt.Sprintf("Received empty provisioning key for id %q.", id),
		)
		return
	}

	result := ProvisioningKeyEphemeralModel{
		ID:              types.StringValue(id),
		AssociationType: types.StringValue(associationType),
		ProvisioningKey: types.StringValue(key.ProvisioningKey),
	}

	if !config.MicroTenantID.IsNull() && !config.MicroTenantID.IsUnknown() {
		result.MicroTenantID = config.MicroTenantID
	} else if key.MicroTenantID != "" {
		result.MicroTenantID = types.StringValue(key.MicroTenantID)
	} else {
		result.MicroTenantID = types.StringNull()
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &result)...)
}
