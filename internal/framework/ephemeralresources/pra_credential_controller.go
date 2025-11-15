package ephemeralresources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	ephschema "github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ ephemeral.EphemeralResource = &PRACredentialControllerEphemeralResource{}

type PRACredentialControllerEphemeralResource struct{}

type PRACredentialControllerEphemeralModel struct {
	Passphrase types.String `tfsdk:"passphrase"`
	Password   types.String `tfsdk:"password"`
	PrivateKey types.String `tfsdk:"private_key"`
}

func NewPRACredentialControllerEphemeralResource() ephemeral.EphemeralResource {
	return &PRACredentialControllerEphemeralResource{}
}

func (r *PRACredentialControllerEphemeralResource) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pra_credential_controller"
}

func (r *PRACredentialControllerEphemeralResource) Schema(ctx context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = ephschema.Schema{
		Description: "Holds PRA credential secrets without persisting them in Terraform state.",
		Attributes: map[string]ephschema.Attribute{
			"passphrase": ephschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "Passphrase protecting the SSH private key.",
			},
			"password": ephschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "Password associated with the credential.",
			},
			"private_key": ephschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "SSH private key associated with the credential.",
			},
		},
	}
}

func (r *PRACredentialControllerEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data PRACredentialControllerEphemeralModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
