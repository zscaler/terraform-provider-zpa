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
	c2cipranges "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/c2c_ip_ranges"
)

var (
	_ datasource.DataSource              = &C2CIPRangesDataSource{}
	_ datasource.DataSourceWithConfigure = &C2CIPRangesDataSource{}
)

func NewC2CIPRangesDataSource() datasource.DataSource {
	return &C2CIPRangesDataSource{}
}

type C2CIPRangesDataSource struct {
	client *client.Client
}

type C2CIPRangesModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	AvailableIPs  types.String `tfsdk:"available_ips"`
	CountryCode   types.String `tfsdk:"country_code"`
	CreationTime  types.String `tfsdk:"creation_time"`
	CustomerID    types.String `tfsdk:"customer_id"`
	Description   types.String `tfsdk:"description"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	IPRangeBegin  types.String `tfsdk:"ip_range_begin"`
	IPRangeEnd    types.String `tfsdk:"ip_range_end"`
	IsDeleted     types.String `tfsdk:"is_deleted"`
	LatitudeInDB  types.String `tfsdk:"latitude_in_db"`
	Location      types.String `tfsdk:"location"`
	LocationHint  types.String `tfsdk:"location_hint"`
	LongitudeInDB types.String `tfsdk:"longitude_in_db"`
	ModifiedBy    types.String `tfsdk:"modified_by"`
	ModifiedTime  types.String `tfsdk:"modified_time"`
	SCCMFlag      types.Bool   `tfsdk:"sccm_flag"`
	SubnetCIDR    types.String `tfsdk:"subnet_cidr"`
	TotalIPs      types.String `tfsdk:"total_ips"`
	UsedIPs       types.String `tfsdk:"used_ips"`
}

func (d *C2CIPRangesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_c2c_ip_ranges"
}

func (d *C2CIPRangesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a Cloud to Cloud IP range by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the C2C IP range.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the C2C IP range.",
			},
			"available_ips":   schema.StringAttribute{Computed: true},
			"country_code":    schema.StringAttribute{Computed: true},
			"creation_time":   schema.StringAttribute{Computed: true},
			"customer_id":     schema.StringAttribute{Computed: true},
			"description":     schema.StringAttribute{Computed: true},
			"enabled":         schema.BoolAttribute{Computed: true},
			"ip_range_begin":  schema.StringAttribute{Computed: true},
			"ip_range_end":    schema.StringAttribute{Computed: true},
			"is_deleted":      schema.StringAttribute{Computed: true},
			"latitude_in_db":  schema.StringAttribute{Computed: true},
			"location":        schema.StringAttribute{Computed: true},
			"location_hint":   schema.StringAttribute{Computed: true},
			"longitude_in_db": schema.StringAttribute{Computed: true},
			"modified_by":     schema.StringAttribute{Computed: true},
			"modified_time":   schema.StringAttribute{Computed: true},
			"sccm_flag":       schema.BoolAttribute{Computed: true},
			"subnet_cidr":     schema.StringAttribute{Computed: true},
			"total_ips":       schema.StringAttribute{Computed: true},
			"used_ips":        schema.StringAttribute{Computed: true},
		},
	}
}

func (d *C2CIPRangesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *C2CIPRangesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data C2CIPRangesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())
	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided to read a C2C IP range.")
		return
	}

	var (
		ipRange *c2cipranges.IPRanges
		err     error
	)

	if id != "" {
		tflog.Debug(ctx, "Retrieving C2C IP range by ID", map[string]any{"id": id})
		ipRange, _, err = c2cipranges.Get(ctx, d.client.Service, id)
	} else {
		tflog.Debug(ctx, "Retrieving C2C IP range by name", map[string]any{"name": name})
		ipRange, _, err = c2cipranges.GetByName(ctx, d.client.Service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read C2C IP range: %v", err))
		return
	}

	if ipRange == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("C2C IP range with id %q or name %q was not found.", id, name))
		return
	}

	data.ID = types.StringValue(ipRange.ID)
	data.Name = stringOrNull(ipRange.Name)
	data.AvailableIPs = stringOrNull(ipRange.AvailableIps)
	data.CountryCode = stringOrNull(ipRange.CountryCode)
	data.CreationTime = stringOrNull(ipRange.CreationTime)
	data.CustomerID = stringOrNull(ipRange.CustomerId)
	data.Description = stringOrNull(ipRange.Description)
	data.Enabled = types.BoolValue(ipRange.Enabled)
	data.IPRangeBegin = stringOrNull(ipRange.IpRangeBegin)
	data.IPRangeEnd = stringOrNull(ipRange.IpRangeEnd)
	data.IsDeleted = stringOrNull(ipRange.IsDeleted)
	data.LatitudeInDB = stringOrNull(ipRange.LatitudeInDb)
	data.Location = stringOrNull(ipRange.Location)
	data.LocationHint = stringOrNull(ipRange.LocationHint)
	data.LongitudeInDB = stringOrNull(ipRange.LongitudeInDb)
	data.ModifiedBy = stringOrNull(ipRange.ModifiedBy)
	data.ModifiedTime = stringOrNull(ipRange.ModifiedTime)
	data.SCCMFlag = types.BoolValue(ipRange.SccmFlag)
	data.SubnetCIDR = stringOrNull(ipRange.SubnetCidr)
	data.TotalIPs = stringOrNull(ipRange.TotalIps)
	data.UsedIPs = stringOrNull(ipRange.UsedIps)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
