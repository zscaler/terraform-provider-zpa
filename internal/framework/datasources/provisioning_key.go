package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/provisioningkey"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
)

var (
	_ datasource.DataSource              = &ProvisioningKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &ProvisioningKeyDataSource{}
)

func NewProvisioningKeyDataSource() datasource.DataSource {
	return &ProvisioningKeyDataSource{}
}

type ProvisioningKeyDataSource struct {
	client *client.Client
}

type ProvisioningKeyDataSourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	AssociationType       types.String `tfsdk:"association_type"`
	AppConnectorGroupID   types.String `tfsdk:"app_connector_group_id"`
	AppConnectorGroupName types.String `tfsdk:"app_connector_group_name"`
	CreationTime          types.String `tfsdk:"creation_time"`
	Enabled               types.Bool   `tfsdk:"enabled"`
	ExpirationInEpochSec  types.String `tfsdk:"expiration_in_epoch_sec"`
	IPAcl                 types.Set    `tfsdk:"ip_acl"`
	MaxUsage              types.String `tfsdk:"max_usage"`
	ModifiedBy            types.String `tfsdk:"modifiedby"`
	ModifiedTime          types.String `tfsdk:"modified_time"`
	ProvisioningKey       types.String `tfsdk:"provisioning_key"`
	EnrollmentCertID      types.String `tfsdk:"enrollment_cert_id"`
	EnrollmentCertName    types.String `tfsdk:"enrollment_cert_name"`
	UIConfig              types.String `tfsdk:"ui_config"`
	UsageCount            types.String `tfsdk:"usage_count"`
	ZComponentID          types.String `tfsdk:"zcomponent_id"`
	ZComponentName        types.String `tfsdk:"zcomponent_name"`
	MicroTenantID         types.String `tfsdk:"microtenant_id"`
	MicroTenantName       types.String `tfsdk:"microtenant_name"`
}

func (d *ProvisioningKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provisioning_key"
}

func (d *ProvisioningKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA provisioning key by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"association_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(provisioningkey.ProvisioningKeyAssociationTypes...),
				},
			},
			"app_connector_group_id":   schema.StringAttribute{Computed: true},
			"app_connector_group_name": schema.StringAttribute{Computed: true},
			"creation_time":            schema.StringAttribute{Computed: true},
			"enabled":                  schema.BoolAttribute{Computed: true},
			"expiration_in_epoch_sec":  schema.StringAttribute{Computed: true},
			"ip_acl": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"max_usage":            schema.StringAttribute{Computed: true},
			"modifiedby":           schema.StringAttribute{Computed: true},
			"modified_time":        schema.StringAttribute{Computed: true},
			"provisioning_key":     schema.StringAttribute{Computed: true},
			"enrollment_cert_id":   schema.StringAttribute{Computed: true},
			"enrollment_cert_name": schema.StringAttribute{Computed: true},
			"ui_config":            schema.StringAttribute{Computed: true},
			"usage_count":          schema.StringAttribute{Computed: true},
			"zcomponent_id":        schema.StringAttribute{Computed: true},
			"zcomponent_name":      schema.StringAttribute{Computed: true},
			"microtenant_id": schema.StringAttribute{
				Optional: true,
			},
			"microtenant_name": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *ProvisioningKeyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *ProvisioningKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProvisioningKeyDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (data.ID.IsNull() || data.ID.ValueString() == "") && (data.Name.IsNull() || data.Name.ValueString() == "") {
		resp.Diagnostics.AddError("Missing Identifier", "Either id or name must be provided")
		return
	}

	associationType := data.AssociationType.ValueString()
	if associationType == "" {
		resp.Diagnostics.AddError("Missing Association Type", "association_type must be provided")
		return
	}

	service := d.client.Service
	if !data.MicroTenantID.IsNull() && data.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(data.MicroTenantID.ValueString())
	}

	var (
		key *provisioningkey.ProvisioningKey
		err error
	)

	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		id := data.ID.ValueString()
		tflog.Info(ctx, "Fetching provisioning key", map[string]any{"id": id, "association_type": associationType})
		key, _, err = provisioningkey.Get(ctx, service, associationType, id)
	} else {
		name := data.Name.ValueString()
		tflog.Info(ctx, "Fetching provisioning key", map[string]any{"name": name, "association_type": associationType})
		key, _, err = provisioningkey.GetByName(ctx, service, associationType, data.Name.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read provisioning key: %v", err))
		return
	}

	flattened, diags := flattenProvisioningKey(ctx, key, associationType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &flattened)...)
}

func flattenProvisioningKey(ctx context.Context, key *provisioningkey.ProvisioningKey, associationType string) (ProvisioningKeyDataSourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	ipAcl, setDiags := types.SetValueFrom(ctx, types.StringType, key.IPACL)
	diags.Append(setDiags...)

	model := ProvisioningKeyDataSourceModel{
		ID:                    types.StringValue(key.ID),
		Name:                  types.StringValue(key.Name),
		AssociationType:       types.StringValue(associationType),
		AppConnectorGroupID:   types.StringValue(key.AppConnectorGroupID),
		AppConnectorGroupName: types.StringValue(key.AppConnectorGroupName),
		CreationTime:          types.StringValue(key.CreationTime),
		Enabled:               types.BoolValue(key.Enabled),
		ExpirationInEpochSec:  types.StringValue(key.ExpirationInEpochSec),
		IPAcl:                 ipAcl,
		MaxUsage:              types.StringValue(key.MaxUsage),
		ModifiedBy:            types.StringValue(key.ModifiedBy),
		ModifiedTime:          types.StringValue(key.ModifiedTime),
		ProvisioningKey:       types.StringValue(key.ProvisioningKey),
		EnrollmentCertID:      types.StringValue(key.EnrollmentCertID),
		EnrollmentCertName:    types.StringValue(key.EnrollmentCertName),
		UIConfig:              types.StringValue(key.UIConfig),
		UsageCount:            types.StringValue(key.UsageCount),
		ZComponentID:          types.StringValue(key.ZcomponentID),
		ZComponentName:        types.StringValue(key.ZcomponentName),
		MicroTenantID:         types.StringValue(key.MicroTenantID),
		MicroTenantName:       types.StringValue(key.MicroTenantName),
	}

	return model, diags
}
