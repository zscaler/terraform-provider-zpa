package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/bacertificate"
)

var (
	_ datasource.DataSource              = &BaCertificateDataSource{}
	_ datasource.DataSourceWithConfigure = &BaCertificateDataSource{}
)

func NewBaCertificateDataSource() datasource.DataSource {
	return &BaCertificateDataSource{}
}

type BaCertificateDataSource struct {
	client *client.Client
}

type BaCertificateDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	CName            types.String `tfsdk:"cname"`
	CertChain        types.String `tfsdk:"cert_chain"`
	Certificate      types.String `tfsdk:"certificate"`
	CreationTime     types.String `tfsdk:"creation_time"`
	IssuedBy         types.String `tfsdk:"issued_by"`
	IssuedTo         types.String `tfsdk:"issued_to"`
	ModifiedBy       types.String `tfsdk:"modifiedby"`
	ModifiedTime     types.String `tfsdk:"modified_time"`
	PublicKey        types.String `tfsdk:"public_key"`
	SAN              types.List   `tfsdk:"san"`
	SerialNo         types.String `tfsdk:"serial_no"`
	Status           types.String `tfsdk:"status"`
	ValidFromInEpoch types.String `tfsdk:"valid_from_in_epochsec"`
	ValidToInEpoch   types.String `tfsdk:"valid_to_in_epochsec"`
	MicroTenantID    types.String `tfsdk:"microtenant_id"`
	MicroTenantName  types.String `tfsdk:"microtenant_name"`
}

func (d *BaCertificateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ba_certificate"
}

func (d *BaCertificateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a Browser Access (BA) certificate by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the BA certificate.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the BA certificate.",
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"cname": schema.StringAttribute{
				Computed: true,
			},
			"cert_chain": schema.StringAttribute{
				Computed: true,
			},
			"certificate": schema.StringAttribute{
				Computed: true,
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"issued_by": schema.StringAttribute{
				Computed: true,
			},
			"issued_to": schema.StringAttribute{
				Computed: true,
			},
			"modifiedby": schema.StringAttribute{
				Computed: true,
			},
			"modified_time": schema.StringAttribute{
				Computed: true,
			},
			"public_key": schema.StringAttribute{
				Computed: true,
			},
			"san": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"serial_no": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"valid_from_in_epochsec": schema.StringAttribute{
				Computed: true,
			},
			"valid_to_in_epochsec": schema.StringAttribute{
				Computed: true,
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID used to scope the lookup.",
			},
			"microtenant_name": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *BaCertificateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BaCertificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data BaCertificateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())
	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided to read a BA certificate.")
		return
	}

	service := d.client.Service
	if !data.MicroTenantID.IsNull() && !data.MicroTenantID.IsUnknown() {
		if microID := strings.TrimSpace(data.MicroTenantID.ValueString()); microID != "" {
			service = service.WithMicroTenant(microID)
		}
	}

	tflog.Debug(ctx, "Retrieving BA certificate", map[string]any{
		"id":   id,
		"name": name,
	})

	var (
		cert *bacertificate.BaCertificate
		err  error
	)

	if id != "" {
		cert, _, err = bacertificate.Get(ctx, service, id)
	} else {
		cert, _, err = bacertificate.GetIssuedByName(ctx, service, name)
	}

	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			resp.Diagnostics.AddError("Not Found", fmt.Sprintf("BA certificate with id %q or name %q was not found.", id, name))
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read BA certificate: %v", err))
		return
	}

	if cert == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("BA certificate with id %q or name %q was not found.", id, name))
		return
	}

	state, diags := flattenBaCertificate(ctx, cert)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.MicroTenantID.IsNull() && strings.TrimSpace(data.MicroTenantID.ValueString()) != "" {
		state.MicroTenantID = types.StringValue(strings.TrimSpace(data.MicroTenantID.ValueString()))
	} else if cert.MicrotenantID != "" {
		state.MicroTenantID = types.StringValue(cert.MicrotenantID)
	} else {
		state.MicroTenantID = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenBaCertificate(ctx context.Context, cert *bacertificate.BaCertificate) (BaCertificateDataSourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	state := BaCertificateDataSourceModel{
		ID:               types.StringValue(cert.ID),
		Name:             types.StringValue(cert.Name),
		Description:      types.StringValue(cert.Description),
		CName:            types.StringValue(cert.CName),
		CertChain:        types.StringValue(cert.CertChain),
		Certificate:      types.StringValue(cert.Certificate),
		CreationTime:     types.StringValue(cert.CreationTime),
		IssuedBy:         types.StringValue(cert.IssuedBy),
		IssuedTo:         types.StringValue(cert.IssuedTo),
		ModifiedBy:       types.StringValue(cert.ModifiedBy),
		ModifiedTime:     types.StringValue(cert.ModifiedTime),
		PublicKey:        types.StringValue(cert.PublicKey),
		SerialNo:         types.StringValue(cert.SerialNo),
		Status:           types.StringValue(cert.Status),
		ValidFromInEpoch: types.StringValue(cert.ValidFromInEpochSec),
		ValidToInEpoch:   types.StringValue(cert.ValidToInEpochSec),
		MicroTenantName:  types.StringValue(cert.MicrotenantName),
	}

	if cert.San != nil {
		sanList, sanDiags := types.ListValueFrom(ctx, types.StringType, cert.San)
		diags.Append(sanDiags...)
		state.SAN = sanList
	} else {
		state.SAN = types.ListNull(types.StringType)
	}

	return state, diags
}
