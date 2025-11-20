package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/c2c_ip_ranges"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &C2CIPRangesResource{}
	_ resource.ResourceWithConfigure   = &C2CIPRangesResource{}
	_ resource.ResourceWithImportState = &C2CIPRangesResource{}
)

func NewC2CIPRangesResource() resource.Resource {
	return &C2CIPRangesResource{}
}

type C2CIPRangesResource struct {
	client *client.Client
}

type C2CIPRangesModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	IPRangeBegin  types.String `tfsdk:"ip_range_begin"`
	IPRangeEnd    types.String `tfsdk:"ip_range_end"`
	Location      types.String `tfsdk:"location"`
	LocationHint  types.String `tfsdk:"location_hint"`
	SccmFlag      types.Bool   `tfsdk:"sccm_flag"`
	SubnetCIDR    types.String `tfsdk:"subnet_cidr"`
	CountryCode   types.String `tfsdk:"country_code"`
	LatitudeInDB  types.String `tfsdk:"latitude_in_db"`
	LongitudeInDB types.String `tfsdk:"longitude_in_db"`
}

func (r *C2CIPRangesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_c2c_ip_ranges"
}

func (r *C2CIPRangesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZPA Cloud Connector (C2C) IP range.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "Identifier of the C2C IP range.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the C2C IP range.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the C2C IP range.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates whether the IP range is enabled.",
			},
			"ip_range_begin": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Beginning IP address of the range.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"ip_range_end": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Ending IP address of the range.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"location": schema.StringAttribute{
				Optional:    true,
				Description: "Location name for the IP range.",
			},
			"location_hint": schema.StringAttribute{
				Optional:    true,
				Description: "Location hint for the IP range.",
			},
			"sccm_flag": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates whether the SCCM flag is enabled.",
			},
			"subnet_cidr": schema.StringAttribute{
				Optional:    true,
				Description: "Subnet CIDR associated with the IP range.",
			},
			"country_code": schema.StringAttribute{
				Optional:    true,
				Description: "Country code of the IP range.",
			},
			"latitude_in_db": schema.StringAttribute{
				Optional:    true,
				Description: "Latitude value stored in the database.",
			},
			"longitude_in_db": schema.StringAttribute{
				Optional:    true,
				Description: "Longitude value stored in the database.",
			},
		},
	}
}

func (r *C2CIPRangesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cl, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = cl
}

func (r *C2CIPRangesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan C2CIPRangesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.client.Service

	payload := expandC2CIPRanges(plan)

	created, _, err := c2c_ip_ranges.Create(ctx, service, &payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create C2C IP range: %v", err))
		return
	}

	state, diags := r.readIntoState(ctx, service, created.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *C2CIPRangesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state C2CIPRangesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.client.Service

	newState, diags := r.readIntoState(ctx, service, state.ID.ValueString())
	if diags.HasError() {
		for _, d := range diags {
			if d.Severity() == diag.SeverityError && strings.Contains(strings.ToLower(d.Detail()), "not found") {
				resp.State.RemoveResource(ctx)
				return
			}
		}
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *C2CIPRangesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan C2CIPRangesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	service := r.client.Service

	// Check if resource still exists before updating
	if _, _, err := c2c_ip_ranges.Get(ctx, service, plan.ID.ValueString()); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
	}

	payload := expandC2CIPRanges(plan)

	if _, err := c2c_ip_ranges.Update(ctx, service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update C2C IP range: %v", err))
		return
	}

	state, diags := r.readIntoState(ctx, service, plan.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *C2CIPRangesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state C2CIPRangesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.client.Service

	if _, err := c2c_ip_ranges.Delete(ctx, service, state.ID.ValueString()); err != nil {
		if helpers.IsObjectNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete C2C IP range: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *C2CIPRangesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
}

func (r *C2CIPRangesResource) readIntoState(ctx context.Context, service *zscaler.Service, id string) (C2CIPRangesModel, diag.Diagnostics) {
	rangeResp, _, err := c2c_ip_ranges.Get(ctx, service, id)
	if err != nil {
		if helpers.IsObjectNotFoundError(err) {
			return C2CIPRangesModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("C2C IP range %s not found", id))}
		}
		return C2CIPRangesModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read C2C IP range: %v", err))}
	}

	state := C2CIPRangesModel{
		ID:            helpers.StringValueOrNull(rangeResp.ID),
		Name:          helpers.StringValueOrNull(rangeResp.Name),
		Description:   helpers.StringValueOrNull(rangeResp.Description),
		Enabled:       types.BoolValue(rangeResp.Enabled),
		IPRangeBegin:  helpers.StringValueOrNull(rangeResp.IpRangeBegin),
		IPRangeEnd:    helpers.StringValueOrNull(rangeResp.IpRangeEnd),
		Location:      helpers.StringValueOrNull(rangeResp.Location),
		LocationHint:  helpers.StringValueOrNull(rangeResp.LocationHint),
		SccmFlag:      types.BoolValue(rangeResp.SccmFlag),
		SubnetCIDR:    helpers.StringValueOrNull(rangeResp.SubnetCidr),
		CountryCode:   helpers.StringValueOrNull(rangeResp.CountryCode),
		LatitudeInDB:  helpers.StringValueOrNull(rangeResp.LatitudeInDb),
		LongitudeInDB: helpers.StringValueOrNull(rangeResp.LongitudeInDb),
	}

	return state, nil
}

func expandC2CIPRanges(plan C2CIPRangesModel) c2c_ip_ranges.IPRanges {
	payload := c2c_ip_ranges.IPRanges{
		ID:            strings.TrimSpace(plan.ID.ValueString()),
		Name:          strings.TrimSpace(plan.Name.ValueString()),
		Description:   strings.TrimSpace(plan.Description.ValueString()),
		Enabled:       helpers.BoolValue(plan.Enabled, false),
		IpRangeBegin:  strings.TrimSpace(plan.IPRangeBegin.ValueString()),
		IpRangeEnd:    strings.TrimSpace(plan.IPRangeEnd.ValueString()),
		Location:      strings.TrimSpace(plan.Location.ValueString()),
		LocationHint:  strings.TrimSpace(plan.LocationHint.ValueString()),
		SccmFlag:      helpers.BoolValue(plan.SccmFlag, false),
		SubnetCidr:    strings.TrimSpace(plan.SubnetCIDR.ValueString()),
		CountryCode:   strings.TrimSpace(plan.CountryCode.ValueString()),
		LatitudeInDb:  strings.TrimSpace(plan.LatitudeInDB.ValueString()),
		LongitudeInDb: strings.TrimSpace(plan.LongitudeInDB.ValueString()),
	}

	return payload
}
