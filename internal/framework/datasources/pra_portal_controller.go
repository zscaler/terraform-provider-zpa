package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praportal"
)

var (
	_ datasource.DataSource              = &PRAPortalControllerDataSource{}
	_ datasource.DataSourceWithConfigure = &PRAPortalControllerDataSource{}
)

func NewPRAPortalControllerDataSource() datasource.DataSource {
	return &PRAPortalControllerDataSource{}
}

type PRAPortalControllerDataSource struct {
	client *client.Client
}

type PRAPortalControllerModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
	CName                   types.String `tfsdk:"cname"`
	Domain                  types.String `tfsdk:"domain"`
	CertificateID           types.String `tfsdk:"certificate_id"`
	CertificateName         types.String `tfsdk:"certificate_name"`
	UserNotification        types.String `tfsdk:"user_notification"`
	UserNotificationEnabled types.Bool   `tfsdk:"user_notification_enabled"`
	CreationTime            types.String `tfsdk:"creation_time"`
	ModifiedBy              types.String `tfsdk:"modified_by"`
	ModifiedTime            types.String `tfsdk:"modified_time"`
	MicroTenantID           types.String `tfsdk:"microtenant_id"`
	MicroTenantName         types.String `tfsdk:"microtenant_name"`
	ExtLabel                types.String `tfsdk:"ext_label"`
	ExtDomain               types.String `tfsdk:"ext_domain"`
	ExtDomainName           types.String `tfsdk:"ext_domain_name"`
	ExtDomainTranslation    types.String `tfsdk:"ext_domain_translation"`
	UserPortalGID           types.String `tfsdk:"user_portal_gid"`
	UserPortalName          types.String `tfsdk:"user_portal_name"`
	GetcName                types.String `tfsdk:"getc_name"`
	ApprovalReviewers       types.Set    `tfsdk:"approval_reviewers"`
}

func (d *PRAPortalControllerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pra_portal_controller"
}

func (d *PRAPortalControllerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA PRA portal controller by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the PRA portal controller.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the PRA portal controller.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID used to scope the lookup.",
			},
			"description": schema.StringAttribute{Computed: true},
			"enabled":     schema.BoolAttribute{Computed: true},
			"cname":       schema.StringAttribute{Computed: true},
			"domain":      schema.StringAttribute{Computed: true},
			"certificate_id": schema.StringAttribute{
				Computed: true,
			},
			"certificate_name": schema.StringAttribute{
				Computed: true,
			},
			"user_notification": schema.StringAttribute{
				Computed: true,
			},
			"user_notification_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"creation_time": schema.StringAttribute{Computed: true},
			"modified_by":   schema.StringAttribute{Computed: true},
			"modified_time": schema.StringAttribute{Computed: true},
			"microtenant_name": schema.StringAttribute{
				Computed: true,
			},
			"ext_label": schema.StringAttribute{
				Optional: true,
			},
			"ext_domain": schema.StringAttribute{
				Optional: true,
			},
			"ext_domain_name": schema.StringAttribute{
				Optional: true,
			},
			"ext_domain_translation": schema.StringAttribute{
				Optional: true,
			},
			"user_portal_gid": schema.StringAttribute{
				Computed: true,
			},
			"user_portal_name": schema.StringAttribute{
				Computed: true,
			},
			"getc_name": schema.StringAttribute{
				Computed: true,
			},
			"approval_reviewers": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (d *PRAPortalControllerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PRAPortalControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data PRAPortalControllerModel
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

	var portal *praportal.PRAPortal
	var err error

	if id != "" {
		tflog.Debug(ctx, "Retrieving PRA portal controller by ID", map[string]any{"id": id})
		portal, _, err = praportal.Get(ctx, service, id)
	} else {
		tflog.Debug(ctx, "Retrieving PRA portal controller by name", map[string]any{"name": name})
		portal, _, err = praportal.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read PRA portal controller: %v", err))
		return
	}

	if portal == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("PRA portal controller with id %q or name %q was not found.", id, name))
		return
	}

	state := PRAPortalControllerModel{
		ID:                      types.StringValue(portal.ID),
		Name:                    types.StringValue(portal.Name),
		Description:             types.StringValue(portal.Description),
		Enabled:                 types.BoolValue(portal.Enabled),
		CName:                   types.StringValue(portal.CName),
		Domain:                  types.StringValue(portal.Domain),
		CertificateID:           types.StringValue(portal.CertificateID),
		CertificateName:         types.StringValue(portal.CertificateName),
		UserNotification:        types.StringValue(portal.UserNotification),
		UserNotificationEnabled: types.BoolValue(portal.UserNotificationEnabled),
		CreationTime:            types.StringValue(portal.CreationTime),
		ModifiedBy:              types.StringValue(portal.ModifiedBy),
		ModifiedTime:            types.StringValue(portal.ModifiedTime),
		MicroTenantName:         types.StringValue(portal.MicroTenantName),
		ExtLabel:                types.StringValue(portal.ExtLabel),
		ExtDomain:               types.StringValue(portal.ExtDomain),
		ExtDomainName:           types.StringValue(portal.ExtDomainName),
		ExtDomainTranslation:    types.StringValue(portal.ExtDomainTranslation),
		UserPortalGID:           types.StringValue(portal.UserPortalGid),
		UserPortalName:          types.StringValue(portal.UserPortalName),
		GetcName:                types.StringValue(portal.GetcName),
	}

	if !data.MicroTenantID.IsNull() && strings.TrimSpace(data.MicroTenantID.ValueString()) != "" {
		state.MicroTenantID = data.MicroTenantID
	} else if portal.MicroTenantID != "" {
		state.MicroTenantID = types.StringValue(portal.MicroTenantID)
	} else {
		state.MicroTenantID = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
