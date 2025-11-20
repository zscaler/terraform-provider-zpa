package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/enrollmentcert"
)

var (
	_ datasource.DataSource              = &EnrollmentCertDataSource{}
	_ datasource.DataSourceWithConfigure = &EnrollmentCertDataSource{}
)

func NewEnrollmentCertDataSource() datasource.DataSource {
	return &EnrollmentCertDataSource{}
}

type EnrollmentCertDataSource struct {
	client *client.Client
}

type EnrollmentCertDataSourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	MicroTenantID           types.String `tfsdk:"microtenant_id"`
	AllowSigning            types.Bool   `tfsdk:"allow_signing"`
	Cname                   types.String `tfsdk:"cname"`
	Certificate             types.String `tfsdk:"certificate"`
	ClientCertType          types.String `tfsdk:"client_cert_type"`
	CreationTime            types.String `tfsdk:"creation_time"`
	CSR                     types.String `tfsdk:"csr"`
	Description             types.String `tfsdk:"description"`
	IssuedBy                types.String `tfsdk:"issued_by"`
	IssuedTo                types.String `tfsdk:"issued_to"`
	ModifiedBy              types.String `tfsdk:"modified_by"`
	ModifiedTime            types.String `tfsdk:"modified_time"`
	ParentCertID            types.String `tfsdk:"parent_cert_id"`
	ParentCertName          types.String `tfsdk:"parent_cert_name"`
	PrivateKey              types.String `tfsdk:"private_key"`
	PrivateKeyPresent       types.Bool   `tfsdk:"private_key_present"`
	SerialNo                types.String `tfsdk:"serial_no"`
	ValidFromInEpochSec     types.String `tfsdk:"valid_from_in_epoch_sec"`
	ValidToInEpochSec       types.String `tfsdk:"valid_to_in_epoch_sec"`
	ZRSAEncryptedPrivateKey types.String `tfsdk:"zrsa_encrypted_private_key"`
	ZRSAEncryptedSessionKey types.String `tfsdk:"zrsa_encrypted_session_key"`
}

func (d *EnrollmentCertDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enrollment_cert"
}

func (d *EnrollmentCertDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA enrollment certificate by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "Identifier of the enrollment certificate.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the enrollment certificate.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID to scope the lookup.",
			},
			"allow_signing": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates if signing is allowed.",
			},
			"cname": schema.StringAttribute{
				Computed: true,
			},
			"certificate": schema.StringAttribute{
				Computed:    true,
				Description: "Certificate text in PEM format.",
			},
			"client_cert_type": schema.StringAttribute{
				Computed: true,
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"csr": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"issued_by": schema.StringAttribute{
				Computed: true,
			},
			"issued_to": schema.StringAttribute{
				Computed: true,
			},
			"modified_by": schema.StringAttribute{
				Computed: true,
			},
			"modified_time": schema.StringAttribute{
				Computed: true,
			},
			"parent_cert_id": schema.StringAttribute{
				Computed: true,
			},
			"parent_cert_name": schema.StringAttribute{
				Computed: true,
			},
			"private_key": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"private_key_present": schema.BoolAttribute{
				Computed: true,
			},
			"serial_no": schema.StringAttribute{
				Computed: true,
			},
			"valid_from_in_epoch_sec": schema.StringAttribute{
				Computed: true,
			},
			"valid_to_in_epoch_sec": schema.StringAttribute{
				Computed: true,
			},
			"zrsa_encrypted_private_key": schema.StringAttribute{
				Computed: true,
			},
			"zrsa_encrypted_session_key": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *EnrollmentCertDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EnrollmentCertDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EnrollmentCertDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())

	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either id or name must be specified.")
		return
	}

	service := d.client.Service
	var configuredMicroTenant string
	if !data.MicroTenantID.IsNull() && !data.MicroTenantID.IsUnknown() {
		configuredMicroTenant = strings.TrimSpace(data.MicroTenantID.ValueString())
		if configuredMicroTenant != "" {
			service = service.WithMicroTenant(configuredMicroTenant)
		} else {
			configuredMicroTenant = ""
		}
	}

	tflog.Debug(ctx, "Reading enrollment certificate", map[string]any{"id": id, "name": name})

	var cert *enrollmentcert.EnrollmentCert
	var err error

	if id != "" {
		cert, _, err = enrollmentcert.Get(ctx, service, id)
	} else {
		cert, _, err = enrollmentcert.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read enrollment certificate: %v", err))
		return
	}

	if cert == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Could not find enrollment certificate with id %q or name %q.", id, name))
		return
	}

	state, diags := flattenEnrollmentCert(cert)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if configuredMicroTenant != "" {
		state.MicroTenantID = types.StringValue(configuredMicroTenant)
	} else if cert.MicrotenantID != "" {
		state.MicroTenantID = types.StringValue(cert.MicrotenantID)
	} else {
		state.MicroTenantID = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenEnrollmentCert(cert *enrollmentcert.EnrollmentCert) (EnrollmentCertDataSourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := EnrollmentCertDataSourceModel{
		MicroTenantID:           types.StringNull(),
		ID:                      types.StringValue(cert.ID),
		Name:                    types.StringValue(cert.Name),
		AllowSigning:            types.BoolValue(cert.AllowSigning),
		Cname:                   types.StringValue(cert.Cname),
		Certificate:             types.StringValue(cert.Certificate),
		ClientCertType:          types.StringValue(cert.ClientCertType),
		CreationTime:            types.StringValue(cert.CreationTime),
		CSR:                     types.StringValue(cert.CSR),
		Description:             types.StringValue(cert.Description),
		IssuedBy:                types.StringValue(cert.IssuedBy),
		IssuedTo:                types.StringValue(cert.IssuedTo),
		ModifiedBy:              types.StringValue(cert.ModifiedBy),
		ModifiedTime:            types.StringValue(cert.ModifiedTime),
		ParentCertID:            types.StringValue(cert.ParentCertID),
		ParentCertName:          types.StringValue(cert.ParentCertName),
		PrivateKey:              types.StringValue(cert.PrivateKey),
		PrivateKeyPresent:       types.BoolValue(cert.PrivateKeyPresent),
		SerialNo:                types.StringValue(cert.SerialNo),
		ValidFromInEpochSec:     types.StringValue(cert.ValidFromInEpochSec),
		ValidToInEpochSec:       types.StringValue(cert.ValidToInEpochSec),
		ZRSAEncryptedPrivateKey: types.StringValue(cert.ZrsaEncryptedPrivateKey),
		ZRSAEncryptedSessionKey: types.StringValue(cert.ZrsaEncryptedSessionKey),
	}

	return model, diags
}
