package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredentialpool"
)

var (
	_ datasource.DataSource              = &PRACredentialPoolDataSource{}
	_ datasource.DataSourceWithConfigure = &PRACredentialPoolDataSource{}
)

func NewPRACredentialPoolDataSource() datasource.DataSource {
	return &PRACredentialPoolDataSource{}
}

type PRACredentialPoolDataSource struct {
	client *client.Client
}

type PRACredentialPoolModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	CredentialMappingCount types.String `tfsdk:"credential_mapping_count"`
	CredentialType         types.String `tfsdk:"credential_type"`
	Credentials            types.List   `tfsdk:"credentials"`
	CreationTime           types.String `tfsdk:"creation_time"`
	ModifiedBy             types.String `tfsdk:"modified_by"`
	ModifiedTime           types.String `tfsdk:"modified_time"`
	MicroTenantID          types.String `tfsdk:"microtenant_id"`
	MicroTenantName        types.String `tfsdk:"microtenant_name"`
}

func (d *PRACredentialPoolDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pra_credential_pool"
}

func (d *PRACredentialPoolDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA Privileged Remote Access (PRA) credential pool by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the PRA credential pool.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the PRA credential pool.",
			},
			"credential_mapping_count": schema.StringAttribute{
				Computed:    true,
				Description: "Number of credential mappings associated with the pool.",
			},
			"credential_type": schema.StringAttribute{
				Computed:    true,
				Description: "Protocol type for the credentials (SSH, RDP, or VNC).",
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"modified_by": schema.StringAttribute{
				Computed: true,
			},
			"modified_time": schema.StringAttribute{
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
		Blocks: map[string]schema.Block{
			"credentials": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *PRACredentialPoolDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PRACredentialPoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data PRACredentialPoolModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := d.client.Service

	if !data.MicroTenantID.IsNull() && !data.MicroTenantID.IsUnknown() {
		if microID := strings.TrimSpace(data.MicroTenantID.ValueString()); microID != "" {
			service = service.WithMicroTenant(microID)
			data.MicroTenantID = types.StringValue(microID)
		}
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())

	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided.")
		return
	}

	var pool *pracredentialpool.CredentialPool
	var err error

	if id != "" {
		tflog.Debug(ctx, "Retrieving PRA credential pool by ID", map[string]any{"id": id})
		pool, _, err = pracredentialpool.Get(ctx, service, id)
	} else {
		tflog.Debug(ctx, "Retrieving PRA credential pool by name", map[string]any{"name": name})
		pool, _, err = pracredentialpool.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read PRA credential pool: %v", err))
		return
	}

	if pool == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("PRA credential pool with id %q or name %q was not found.", id, name))
		return
	}

	credentialsList, diags := flattenPRACredentials(ctx, pool.PRACredentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(pool.ID)
	data.Name = types.StringValue(pool.Name)
	data.CredentialType = types.StringValue(pool.CredentialType)
	data.CredentialMappingCount = types.StringValue(pool.CredentialMappingCount)
	data.CreationTime = types.StringValue(pool.CreationTime)
	data.ModifiedBy = types.StringValue(pool.ModifiedBy)
	data.ModifiedTime = types.StringValue(pool.ModifiedTime)
	data.Credentials = credentialsList

	if !data.MicroTenantID.IsNull() && strings.TrimSpace(data.MicroTenantID.ValueString()) != "" {
		// keep user-provided
	} else if pool.MicroTenantID != "" {
		data.MicroTenantID = types.StringValue(pool.MicroTenantID)
	} else {
		data.MicroTenantID = types.StringNull()
	}

	if pool.MicroTenantName != "" {
		data.MicroTenantName = types.StringValue(pool.MicroTenantName)
	} else {
		data.MicroTenantName = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenPRACredentials(ctx context.Context, creds []common.CommonIDName) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}

	if len(creds) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	values := make([]attr.Value, 0, len(creds))
	for _, cred := range creds {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":   types.StringValue(cred.ID),
			"name": types.StringValue(cred.Name),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}
