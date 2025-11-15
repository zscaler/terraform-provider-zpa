package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cbicertificatecontroller "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbicertificatecontroller"
)

var (
	_ datasource.DataSource              = &CBICertificatesDataSource{}
	_ datasource.DataSourceWithConfigure = &CBICertificatesDataSource{}
)

func NewCBICertificatesDataSource() datasource.DataSource {
	return &CBICertificatesDataSource{}
}

type CBICertificatesDataSource struct {
	client *client.Client
}

type CBICertificatesModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	PEM       types.String `tfsdk:"pem"`
	IsDefault types.Bool   `tfsdk:"is_default"`
}

func (d *CBICertificatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_browser_isolation_certificate"
}

func (d *CBICertificatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a Cloud Browser Isolation certificate by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the CBI certificate.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the CBI certificate.",
			},
			"pem": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"is_default": schema.BoolAttribute{Computed: true},
		},
	}
}

func (d *CBICertificatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cl, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = cl
}

func (d *CBICertificatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data CBICertificatesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())
	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided to read a CBI certificate.")
		return
	}

	var (
		certificate *cbicertificatecontroller.CBICertificate
		err         error
	)

	if id != "" {
		tflog.Debug(ctx, "Retrieving CBI certificate by ID", map[string]any{"id": id})
		certificate, _, err = cbicertificatecontroller.Get(ctx, d.client.Service, id)
	} else {
		tflog.Debug(ctx, "Retrieving CBI certificate by name", map[string]any{"name": name})
		certificate, _, err = cbicertificatecontroller.GetByName(ctx, d.client.Service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read CBI certificate: %v", err))
		return
	}

	if certificate == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("CBI certificate with id %q or name %q was not found.", id, name))
		return
	}

	data.ID = types.StringValue(certificate.ID)
	data.Name = stringOrNull(certificate.Name)
	data.PEM = stringOrNull(certificate.PEM)
	data.IsDefault = types.BoolValue(certificate.IsDefault)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
