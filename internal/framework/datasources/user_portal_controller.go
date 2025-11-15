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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
)

var (
	_ datasource.DataSource              = &UserPortalControllerDataSource{}
	_ datasource.DataSourceWithConfigure = &UserPortalControllerDataSource{}
)

func NewUserPortalControllerDataSource() datasource.DataSource {
	return &UserPortalControllerDataSource{}
}

type UserPortalControllerDataSource struct {
	client *client.Client
}

type UserPortalControllerModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	CertificateID           types.String `tfsdk:"certificate_id"`
	CertificateName         types.String `tfsdk:"certificate_name"`
	CreationTime            types.String `tfsdk:"creation_time"`
	Description             types.String `tfsdk:"description"`
	Domain                  types.String `tfsdk:"domain"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
	ExtDomain               types.String `tfsdk:"ext_domain"`
	ExtDomainName           types.String `tfsdk:"ext_domain_name"`
	ExtDomainTranslation    types.String `tfsdk:"ext_domain_translation"`
	ExtLabel                types.String `tfsdk:"ext_label"`
	GetcName                types.String `tfsdk:"getc_name"`
	ImageData               types.String `tfsdk:"image_data"`
	ModifiedBy              types.String `tfsdk:"modified_by"`
	ModifiedTime            types.String `tfsdk:"modified_time"`
	MicroTenantID           types.String `tfsdk:"microtenant_id"`
	MicroTenantName         types.String `tfsdk:"microtenant_name"`
	UserNotification        types.String `tfsdk:"user_notification"`
	UserNotificationEnabled types.Bool   `tfsdk:"user_notification_enabled"`
	ManagedByZS             types.Bool   `tfsdk:"managed_by_zs"`
}

func (d *UserPortalControllerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_portal_controller"
}

func (d *UserPortalControllerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a user portal controller by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the user portal controller.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the user portal controller.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID used to scope the lookup.",
			},
			"certificate_id": schema.StringAttribute{
				Computed: true,
			},
			"certificate_name": schema.StringAttribute{
				Computed: true,
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"domain": schema.StringAttribute{
				Computed: true,
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
			},
			"ext_domain": schema.StringAttribute{
				Computed: true,
			},
			"ext_domain_name": schema.StringAttribute{
				Computed: true,
			},
			"ext_domain_translation": schema.StringAttribute{
				Computed: true,
			},
			"ext_label": schema.StringAttribute{
				Computed: true,
			},
			"getc_name": schema.StringAttribute{
				Computed: true,
			},
			"image_data": schema.StringAttribute{
				Computed: true,
			},
			"modified_by": schema.StringAttribute{
				Computed: true,
			},
			"modified_time": schema.StringAttribute{
				Computed: true,
			},
			"microtenant_name": schema.StringAttribute{
				Computed: true,
			},
			"user_notification": schema.StringAttribute{
				Computed: true,
			},
			"user_notification_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"managed_by_zs": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (d *UserPortalControllerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserPortalControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data UserPortalControllerModel
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

	var controller *portal_controller.UserPortalController
	var err error

	if id != "" {
		tflog.Debug(ctx, "Retrieving user portal controller by ID", map[string]any{"id": id})
		controller, _, err = portal_controller.Get(ctx, service, id)
	} else {
		tflog.Debug(ctx, "Retrieving user portal controller by name", map[string]any{"name": name})
		controller, _, err = portal_controller.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user portal controller: %v", err))
		return
	}

	if controller == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("User portal controller with id %q or name %q not found.", id, name))
		return
	}

	state := UserPortalControllerModel{
		ID:                      types.StringValue(controller.ID),
		Name:                    types.StringValue(controller.Name),
		CertificateID:           types.StringValue(controller.CertificateId),
		CertificateName:         types.StringValue(controller.CertificateName),
		CreationTime:            types.StringValue(controller.CreationTime),
		Description:             types.StringValue(controller.Description),
		Domain:                  types.StringValue(controller.Domain),
		Enabled:                 types.BoolValue(controller.Enabled),
		ExtDomain:               types.StringValue(controller.ExtDomain),
		ExtDomainName:           types.StringValue(controller.ExtDomainName),
		ExtDomainTranslation:    types.StringValue(controller.ExtDomainTranslation),
		ExtLabel:                types.StringValue(controller.ExtLabel),
		GetcName:                types.StringValue(controller.GetcName),
		ImageData:               types.StringValue(controller.ImageData),
		ModifiedBy:              types.StringValue(controller.ModifiedBy),
		ModifiedTime:            types.StringValue(controller.ModifiedTime),
		MicroTenantID:           types.StringValue(controller.MicrotenantId),
		MicroTenantName:         types.StringValue(controller.MicrotenantName),
		UserNotification:        types.StringValue(controller.UserNotification),
		UserNotificationEnabled: types.BoolValue(controller.UserNotificationEnabled),
		ManagedByZS:             types.BoolValue(controller.ManagedByZS),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
