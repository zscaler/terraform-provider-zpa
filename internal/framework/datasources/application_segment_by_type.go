package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentbytype"
)

var (
	_ datasource.DataSource              = &ApplicationSegmentByTypeDataSource{}
	_ datasource.DataSourceWithConfigure = &ApplicationSegmentByTypeDataSource{}
)

func NewApplicationSegmentByTypeDataSource() datasource.DataSource {
	return &ApplicationSegmentByTypeDataSource{}
}

type ApplicationSegmentByTypeDataSource struct {
	client *client.Client
}

type ApplicationSegmentByTypeModel struct {
	ID                  types.String `tfsdk:"id"`
	AppID               types.String `tfsdk:"app_id"`
	Name                types.String `tfsdk:"name"`
	ApplicationType     types.String `tfsdk:"application_type"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	Domain              types.String `tfsdk:"domain"`
	ApplicationPort     types.String `tfsdk:"application_port"`
	ApplicationProtocol types.String `tfsdk:"application_protocol"`
	CertificateID       types.String `tfsdk:"certificate_id"`
	CertificateName     types.String `tfsdk:"certificate_name"`
	MicroTenantID       types.String `tfsdk:"microtenant_id"`
	MicroTenantName     types.String `tfsdk:"microtenant_name"`
}

func (d *ApplicationSegmentByTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_segment_by_type"
}

func (d *ApplicationSegmentByTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves application segments filtered by type (Browser Access, Inspect, or Secure Remote Access).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the application segment result.",
			},
			"app_id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Optional name filter for the application segment.",
			},
			"application_type": schema.StringAttribute{
				Required:    true,
				Description: "Application type to filter by. Valid values: BROWSER_ACCESS, INSPECT, SECURE_REMOTE_ACCESS.",
				Validators: []validator.String{
					stringvalidator.OneOf("BROWSER_ACCESS", "INSPECT", "SECURE_REMOTE_ACCESS"),
				},
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
			},
			"domain": schema.StringAttribute{
				Computed: true,
			},
			"application_port": schema.StringAttribute{
				Computed: true,
			},
			"application_protocol": schema.StringAttribute{
				Computed: true,
			},
			"certificate_id": schema.StringAttribute{
				Computed: true,
			},
			"certificate_name": schema.StringAttribute{
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

func (d *ApplicationSegmentByTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ApplicationSegmentByTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data ApplicationSegmentByTypeModel
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

	appType := data.ApplicationType.ValueString()
	name := strings.TrimSpace(data.Name.ValueString())

	tflog.Debug(ctx, "Retrieving application segments by type", map[string]any{
		"type": appType,
		"name": name,
	})

	segments, _, err := applicationsegmentbytype.GetByApplicationType(ctx, service, name, appType, true)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read application segments by type: %v", err))
		return
	}

	if len(segments) == 0 {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No application segment found with name %q and type %q.", name, appType))
		return
	}

	segment := segments[0]

	data.ID = types.StringValue(segment.ID)
	data.AppID = types.StringValue(segment.AppID)
	data.Name = types.StringValue(segment.Name)
	data.Enabled = types.BoolValue(segment.Enabled)
	data.Domain = types.StringValue(segment.Domain)
	data.ApplicationPort = types.StringValue(segment.ApplicationPort)
	data.ApplicationProtocol = types.StringValue(segment.ApplicationProtocol)
	data.CertificateID = types.StringValue(segment.CertificateID)
	data.CertificateName = types.StringValue(segment.CertificateName)
	data.MicroTenantID = types.StringValue(segment.MicroTenantID)
	data.MicroTenantName = types.StringValue(segment.MicroTenantName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
