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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/aup"
)

var (
	_ datasource.DataSource              = &UserPortalAUPDataSource{}
	_ datasource.DataSourceWithConfigure = &UserPortalAUPDataSource{}
)

func NewUserPortalAUPDataSource() datasource.DataSource {
	return &UserPortalAUPDataSource{}
}

type UserPortalAUPDataSource struct {
	client *client.Client
}

type UserPortalAUPModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	AUP             types.String `tfsdk:"aup"`
	Email           types.String `tfsdk:"email"`
	PhoneNum        types.String `tfsdk:"phone_num"`
	MicroTenantID   types.String `tfsdk:"microtenant_id"`
	MicroTenantName types.String `tfsdk:"microtenant_name"`
	CreationTime    types.String `tfsdk:"creation_time"`
	ModifiedTime    types.String `tfsdk:"modified_time"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
}

func (d *UserPortalAUPDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_portal_aup"
}

func (d *UserPortalAUPDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a user portal Acceptable Use Policy (AUP) by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the user portal AUP.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the user portal AUP.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID used to scope the lookup.",
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
			},
			"aup": schema.StringAttribute{
				Computed: true,
			},
			"email": schema.StringAttribute{
				Computed: true,
			},
			"phone_num": schema.StringAttribute{
				Computed: true,
			},
			"microtenant_name": schema.StringAttribute{
				Computed: true,
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"modified_time": schema.StringAttribute{
				Computed: true,
			},
			"modified_by": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *UserPortalAUPDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserPortalAUPDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data UserPortalAUPModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := d.client.Service
	if !data.MicroTenantID.IsNull() && !data.MicroTenantID.IsUnknown() {
		microID := strings.TrimSpace(data.MicroTenantID.ValueString())
		if microID != "" {
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

	var aupResp *aup.UserPortalAup
	var err error

	if id != "" {
		tflog.Debug(ctx, "Retrieving user portal AUP by ID", map[string]any{"id": id})
		aupResp, _, err = aup.Get(ctx, service, id)
	} else {
		tflog.Debug(ctx, "Retrieving user portal AUP by name", map[string]any{"name": name})
		aupResp, _, err = aup.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user portal AUP: %v", err))
		return
	}

	if aupResp == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("User portal AUP with id %q or name %q not found.", id, name))
		return
	}

	state := UserPortalAUPModel{
		ID:              types.StringValue(aupResp.ID),
		Name:            types.StringValue(aupResp.Name),
		Description:     types.StringValue(aupResp.Description),
		Enabled:         types.BoolValue(aupResp.Enabled),
		AUP:             types.StringValue(aupResp.Aup),
		Email:           types.StringValue(aupResp.Email),
		PhoneNum:        types.StringValue(aupResp.PhoneNum),
		MicroTenantID:   types.StringValue(aupResp.MicrotenantID),
		MicroTenantName: types.StringValue(aupResp.MicrotenantName),
		CreationTime:    types.StringValue(aupResp.CreationTime),
		ModifiedTime:    types.StringValue(aupResp.ModifiedTime),
		ModifiedBy:      types.StringValue(aupResp.ModifiedBy),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
