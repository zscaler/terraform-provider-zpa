package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredential"
)

var (
	_ datasource.DataSource              = &PRACredentialControllerDataSource{}
	_ datasource.DataSourceWithConfigure = &PRACredentialControllerDataSource{}
)

func NewPRACredentialControllerDataSource() datasource.DataSource {
	return &PRACredentialControllerDataSource{}
}

type PRACredentialControllerDataSource struct {
	client *client.Client
}

type PRACredentialControllerDataSourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	LastCredentialResetTime types.String `tfsdk:"last_credential_reset_time"`
	CredentialType          types.String `tfsdk:"credential_type"`
	UserDomain              types.String `tfsdk:"user_domain"`
	Username                types.String `tfsdk:"username"`
	Password                types.String `tfsdk:"password"`
	CreationTime            types.String `tfsdk:"creation_time"`
	ModifiedBy              types.String `tfsdk:"modified_by"`
	ModifiedTime            types.String `tfsdk:"modified_time"`
	MicroTenantID           types.String `tfsdk:"microtenant_id"`
	MicroTenantName         types.String `tfsdk:"microtenant_name"`
}

func (d *PRACredentialControllerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pra_credential_controller"
}

func (d *PRACredentialControllerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a Privileged Remote Access (PRA) credential controller by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "The unique identifier of the privileged credential",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the privileged credential",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description of the privileged credential",
			},
			"last_credential_reset_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time the privileged credential was last reset",
			},
			"credential_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of supported credential",
			},
			"user_domain": schema.StringAttribute{
				Computed:    true,
				Description: "The domain name associated with the username. You can also include the domain name as part of the username. The domain name only needs to be specified with logging in to an RDP console that is connected to an Active Directory Domain.",
			},
			"username": schema.StringAttribute{
				Computed:    true,
				Description: "The username for the login you want to use for the privileged credential",
			},
			"password": schema.StringAttribute{
				Computed:    true,
				Description: "The password associated with the username for the login you want to use for the privileged credential",
			},
			"creation_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time the privileged credential is created",
			},
			"modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the tenant who modified the privileged credential",
			},
			"modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time the privileged credential is modified",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Description: "The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.",
			},
			"microtenant_name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the Microtenant",
			},
		},
	}
}

func (d *PRACredentialControllerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PRACredentialControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before reading data sources.")
		return
	}

	var model PRACredentialControllerDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := d.serviceForMicrotenant(model.MicroTenantID)

	id := strings.TrimSpace(model.ID.ValueString())
	name := strings.TrimSpace(model.Name.ValueString())

	if id == "" && name == "" {
		resp.Diagnostics.AddError(
			"Missing Identifier",
			"Either 'id' or 'name' must be provided.",
		)
		return
	}

	var credential *pracredential.Credential
	var err error

	if id != "" {
		credential, _, err = pracredential.Get(ctx, service, id)
	} else {
		credential, _, err = pracredential.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to find credential controller with name '%s' or id '%s': %v", name, id, err),
		)
		return
	}

	if credential == nil {
		resp.Diagnostics.AddError(
			"Not Found",
			fmt.Sprintf("Couldn't find any credential controller with name '%s' or id '%s'", name, id),
		)
		return
	}

	model.ID = types.StringValue(credential.ID)
	model.Name = helpers.StringValueOrNull(credential.Name)
	model.Description = helpers.StringValueOrNull(credential.Description)
	model.LastCredentialResetTime = helpers.StringValueOrNull(credential.LastCredentialResetTime)
	model.CredentialType = helpers.StringValueOrNull(credential.CredentialType)
	model.UserDomain = helpers.StringValueOrNull(credential.UserDomain)
	model.Username = helpers.StringValueOrNull(credential.UserName)
	model.Password = helpers.StringValueOrNull(credential.Password)
	model.CreationTime = helpers.StringValueOrNull(credential.CreationTime)
	model.ModifiedBy = helpers.StringValueOrNull(credential.ModifiedBy)
	model.ModifiedTime = helpers.StringValueOrNull(credential.ModifiedTime)
	model.MicroTenantID = helpers.StringValueOrNull(credential.MicroTenantID)
	model.MicroTenantName = helpers.StringValueOrNull(credential.MicroTenantName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (d *PRACredentialControllerDataSource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := d.client.Service
	id := helpers.StringValue(microtenantID)
	if id != "" {
		service = service.WithMicroTenant(id)
	}
	return service
}
