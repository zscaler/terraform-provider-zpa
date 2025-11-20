package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment_share"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &ApplicationSegmentResource{}
	_ resource.ResourceWithConfigure   = &ApplicationSegmentResource{}
	_ resource.ResourceWithImportState = &ApplicationSegmentResource{}
)

func NewApplicationSegmentResource() resource.Resource {
	return &ApplicationSegmentResource{}
}

type ApplicationSegmentResource struct {
	client *client.Client
}

type ApplicationSegmentModel struct {
	ID                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	SegmentGroupID            types.String `tfsdk:"segment_group_id"`
	SegmentGroupName          types.String `tfsdk:"segment_group_name"`
	BypassType                types.String `tfsdk:"bypass_type"`
	BypassOnReauth            types.Bool   `tfsdk:"bypass_on_reauth"`
	ConfigSpace               types.String `tfsdk:"config_space"`
	DomainNames               types.List   `tfsdk:"domain_names"`
	DoubleEncrypt             types.Bool   `tfsdk:"double_encrypt"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
	InspectTrafficWithZia     types.Bool   `tfsdk:"inspect_traffic_with_zia"`
	PassiveHealthEnabled      types.Bool   `tfsdk:"passive_health_enabled"`
	HealthCheckType           types.String `tfsdk:"health_check_type"`
	HealthReporting           types.String `tfsdk:"health_reporting"`
	IcmpAccessType            types.String `tfsdk:"icmp_access_type"`
	IPAnchored                types.Bool   `tfsdk:"ip_anchored"`
	FqdnDnsCheck              types.Bool   `tfsdk:"fqdn_dns_check"`
	SelectConnectorCloseToApp types.Bool   `tfsdk:"select_connector_close_to_app"`
	UseInDrMode               types.Bool   `tfsdk:"use_in_dr_mode"`
	IsIncompleteDRConfig      types.Bool   `tfsdk:"is_incomplete_dr_config"`
	IsCnameEnabled            types.Bool   `tfsdk:"is_cname_enabled"`
	TCPKeepAlive              types.String `tfsdk:"tcp_keep_alive"`
	MicroTenantID             types.String `tfsdk:"microtenant_id"`
	MicroTenantName           types.String `tfsdk:"microtenant_name"`
	ShareToMicrotenants       types.Set    `tfsdk:"share_to_microtenants"`
	TCPPortRanges             types.List   `tfsdk:"tcp_port_ranges"`
	UDPPortRanges             types.List   `tfsdk:"udp_port_ranges"`
	TCPPortRange              types.List   `tfsdk:"tcp_port_range"`
	UDPPortRange              types.List   `tfsdk:"udp_port_range"`
	ServerGroups              types.List   `tfsdk:"server_groups"`
	ZpnERID                   types.List   `tfsdk:"zpn_er_id"`
	MatchStyle                types.String `tfsdk:"match_style"`
	APIProtectionEnabled      types.Bool   `tfsdk:"api_protection_enabled"`
}

func (r *ApplicationSegmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_segment"
}

func (r *ApplicationSegmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZPA application segment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the application segment.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"segment_group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the associated segment group.",
			},
			"segment_group_name": schema.StringAttribute{Computed: true},
			"bypass_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Controls whether users can bypass ZPA to access applications.",
				Validators: []validator.String{
					stringvalidator.OneOf("ALWAYS", "NEVER", "ON_NET"),
				},
			},
			"bypass_on_reauth": schema.BoolAttribute{Optional: true, Computed: true},
			"config_space": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("DEFAULT"),
				Validators: []validator.String{
					stringvalidator.OneOf("DEFAULT", "SIEM"),
				},
			},
			"domain_names": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "List of domains and IPs for the application segment.",
			},
			"double_encrypt":           schema.BoolAttribute{Optional: true, Computed: true},
			"enabled":                  schema.BoolAttribute{Optional: true, Computed: true},
			"inspect_traffic_with_zia": schema.BoolAttribute{Optional: true, Computed: true},
			"passive_health_enabled":   schema.BoolAttribute{Optional: true, Computed: true},
			"health_check_type":        schema.StringAttribute{Optional: true, Computed: true},
			"health_reporting":         schema.StringAttribute{Optional: true, Computed: true},
			"icmp_access_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("NONE"),
				Validators: []validator.String{
					stringvalidator.OneOf("PING_TRACEROUTING", "PING", "NONE"),
				},
			},
			"ip_anchored":                   schema.BoolAttribute{Optional: true, Computed: true},
			"fqdn_dns_check":                schema.BoolAttribute{Optional: true, Computed: true},
			"select_connector_close_to_app": schema.BoolAttribute{Optional: true, Computed: true},
			"use_in_dr_mode":                schema.BoolAttribute{Optional: true, Computed: true},
			"is_incomplete_dr_config":       schema.BoolAttribute{Optional: true, Computed: true},
			"is_cname_enabled":              schema.BoolAttribute{Optional: true, Computed: true},
			"tcp_keep_alive":                schema.StringAttribute{Optional: true, Computed: true},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant identifier used to scope API calls.",
			},
			"microtenant_name": schema.StringAttribute{Computed: true},
			"share_to_microtenants": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of microtenant IDs to share this application segment with.",
			},
			"tcp_port_ranges":        schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"udp_port_ranges":        schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"match_style":            schema.StringAttribute{Computed: true},
			"api_protection_enabled": schema.BoolAttribute{Optional: true, Computed: true},
		},
		Blocks: map[string]schema.Block{
			"tcp_port_range": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"from": schema.StringAttribute{Optional: true, Computed: true},
						"to":   schema.StringAttribute{Optional: true, Computed: true},
					},
				},
			},
			"udp_port_range": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"from": schema.StringAttribute{Optional: true, Computed: true},
						"to":   schema.StringAttribute{Optional: true, Computed: true},
					},
				},
			},
			"zpn_er_id": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.SetAttribute{ElementType: types.StringType, Optional: true, Computed: true},
					},
				},
			},
			// server_groups: TypeList in SDKv2, id is TypeSet (Required)
			// Using ListNestedBlock for block syntax support
			"server_groups": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.SetAttribute{ElementType: types.StringType, Required: true},
					},
				},
			},
		},
	}
}

func (r *ApplicationSegmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationSegmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan ApplicationSegmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := validateSegmentGroup(plan.SegmentGroupID); err != nil {
		resp.Diagnostics.AddError("Validation Error", err.Error())
		return
	}

	service := r.serviceForMicrotenant(plan.MicroTenantID)

	payload, diags := r.expandApplicationSegment(ctx, plan, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if portDiags := helpers.ValidateAppPorts(payload.SelectConnectorCloseToApp, payload.UDPAppPortRange, payload.UDPPortRanges); portDiags.HasError() {
		resp.Diagnostics.Append(portDiags...)
		return
	}

	created, _, err := applicationsegment.Create(ctx, service, payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create application segment: %v", err))
		return
	}

	plan.ID = types.StringValue(created.ID)

	if shareDiags := r.shareApplicationSegment(ctx, service, created.ID, plan.ShareToMicrotenants, plan.MicroTenantID); shareDiags.HasError() {
		resp.Diagnostics.Append(shareDiags...)
		return
	}

	state, readDiags := r.readIntoState(ctx, service, created.ID, plan.MicroTenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationSegmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state ApplicationSegmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicroTenantID)

	newState, diags := r.readIntoState(ctx, service, state.ID.ValueString(), state.MicroTenantID)
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

func (r *ApplicationSegmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan ApplicationSegmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	if err := validateSegmentGroupOnUpdate(plan.SegmentGroupID); err != nil {
		resp.Diagnostics.AddError("Validation Error", err.Error())
		return
	}

	service := r.serviceForMicrotenant(plan.MicroTenantID)

	existing, _, err := applicationsegment.Get(ctx, service, plan.ID.ValueString())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve application segment: %v", err))
		return
	}

	payload, diags := r.expandApplicationSegment(ctx, plan, existing)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if portDiags := helpers.ValidateAppPorts(payload.SelectConnectorCloseToApp, payload.UDPAppPortRange, payload.UDPPortRanges); portDiags.HasError() {
		resp.Diagnostics.Append(portDiags...)
		return
	}

	if _, err := applicationsegment.Update(ctx, service, plan.ID.ValueString(), payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update application segment: %v", err))
		return
	}

	if shareDiags := r.shareApplicationSegment(ctx, service, plan.ID.ValueString(), plan.ShareToMicrotenants, plan.MicroTenantID); shareDiags.HasError() {
		resp.Diagnostics.Append(shareDiags...)
		return
	}

	state, readDiags := r.readIntoState(ctx, service, plan.ID.ValueString(), plan.MicroTenantID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationSegmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state ApplicationSegmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicroTenantID)

	helpers.DetachAppFromAllPolicyRules(ctx, service, state.ID.ValueString())

	if _, err := applicationsegment.Delete(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete application segment: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *ApplicationSegmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
}

func (r *ApplicationSegmentResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && !microtenantID.IsUnknown() {
		trimmed := strings.TrimSpace(microtenantID.ValueString())
		if trimmed != "" {
			service = service.WithMicroTenant(trimmed)
		}
	}
	return service
}

func (r *ApplicationSegmentResource) expandApplicationSegment(ctx context.Context, plan ApplicationSegmentModel, existing *applicationsegment.ApplicationSegmentResource) (applicationsegment.ApplicationSegmentResource, diag.Diagnostics) {
	var diags diag.Diagnostics

	domainNames, domainDiags := helpers.ListValueToStringSlice(ctx, plan.DomainNames)
	diags.Append(domainDiags...)

	tcpRanges, tcpRangeDiags := helpers.ListValueToStringSlice(ctx, plan.TCPPortRanges)
	diags.Append(tcpRangeDiags...)

	udpRanges, udpRangeDiags := helpers.ListValueToStringSlice(ctx, plan.UDPPortRanges)
	diags.Append(udpRangeDiags...)

	tcpPorts, tcpDiags := helpers.ExpandNetworkPorts(ctx, plan.TCPPortRange)
	diags.Append(tcpDiags...)

	udpPorts, udpDiags := helpers.ExpandNetworkPorts(ctx, plan.UDPPortRange)
	diags.Append(udpDiags...)

	serverGroups, sgDiags := helpers.ExpandServerGroups(ctx, plan.ServerGroups)
	diags.Append(sgDiags...)

	zpnER, zpnDiags := helpers.ExpandZPNERID(ctx, plan.ZpnERID)
	diags.Append(zpnDiags...)

	shareTo, shareDiags := helpers.SetValueToStringSlice(ctx, plan.ShareToMicrotenants)
	diags.Append(shareDiags...)

	if diags.HasError() {
		return applicationsegment.ApplicationSegmentResource{}, diags
	}

	payload := applicationsegment.ApplicationSegmentResource{
		ID:                        strings.TrimSpace(plan.ID.ValueString()),
		Name:                      strings.TrimSpace(plan.Name.ValueString()),
		SegmentGroupID:            strings.TrimSpace(plan.SegmentGroupID.ValueString()),
		SegmentGroupName:          strings.TrimSpace(plan.SegmentGroupName.ValueString()),
		BypassType:                strings.TrimSpace(plan.BypassType.ValueString()),
		BypassOnReauth:            helpers.BoolValue(plan.BypassOnReauth, false),
		ConfigSpace:               strings.TrimSpace(plan.ConfigSpace.ValueString()),
		IcmpAccessType:            strings.TrimSpace(plan.IcmpAccessType.ValueString()),
		Description:               strings.TrimSpace(plan.Description.ValueString()),
		DomainNames:               domainNames,
		HealthCheckType:           strings.TrimSpace(plan.HealthCheckType.ValueString()),
		MatchStyle:                strings.TrimSpace(plan.MatchStyle.ValueString()),
		HealthReporting:           strings.TrimSpace(plan.HealthReporting.ValueString()),
		TCPKeepAlive:              strings.TrimSpace(plan.TCPKeepAlive.ValueString()),
		MicroTenantID:             strings.TrimSpace(plan.MicroTenantID.ValueString()),
		PassiveHealthEnabled:      helpers.BoolValue(plan.PassiveHealthEnabled, false),
		InspectTrafficWithZia:     helpers.BoolValue(plan.InspectTrafficWithZia, false),
		DoubleEncrypt:             helpers.BoolValue(plan.DoubleEncrypt, false),
		Enabled:                   helpers.BoolValue(plan.Enabled, false),
		IpAnchored:                helpers.BoolValue(plan.IPAnchored, false),
		IsCnameEnabled:            helpers.BoolValue(plan.IsCnameEnabled, false),
		SelectConnectorCloseToApp: helpers.BoolValue(plan.SelectConnectorCloseToApp, false),
		UseInDrMode:               helpers.BoolValue(plan.UseInDrMode, false),
		IsIncompleteDRConfig:      helpers.BoolValue(plan.IsIncompleteDRConfig, false),
		FQDNDnsCheck:              helpers.BoolValue(plan.FqdnDnsCheck, false),
		APIProtectionEnabled:      helpers.BoolValue(plan.APIProtectionEnabled, false),
		ShareToMicrotenants:       shareTo,
		ServerGroups:              serverGroups,
		ZPNERID:                   zpnER,
		TCPPortRanges:             tcpRanges,
		UDPPortRanges:             udpRanges,
		TCPAppPortRange:           tcpPorts,
		UDPAppPortRange:           udpPorts,
	}

	if existing != nil {
		payload.ID = existing.ID

		existingTCP := existing.TCPPortRanges
		existingUDP := existing.UDPPortRanges

		if len(payload.TCPPortRanges) == 0 && len(existingTCP) > 0 {
			payload.TCPPortRanges = existingTCP
		}
		if len(payload.UDPPortRanges) == 0 && len(existingUDP) > 0 {
			payload.UDPPortRanges = existingUDP
		}

		if len(payload.TCPAppPortRange) == 0 && len(existing.TCPAppPortRange) > 0 {
			payload.TCPAppPortRange = existing.TCPAppPortRange
		}
		if len(payload.UDPAppPortRange) == 0 && len(existing.UDPAppPortRange) > 0 {
			payload.UDPAppPortRange = existing.UDPAppPortRange
		}
	}

	return payload, diags
}

func (r *ApplicationSegmentResource) readIntoState(ctx context.Context, service *zscaler.Service, id string, microTenantID types.String) (ApplicationSegmentModel, diag.Diagnostics) {
	segment, _, err := applicationsegment.Get(ctx, service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return ApplicationSegmentModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Application segment %s not found", id))}
		}
		return ApplicationSegmentModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read application segment: %v", err))}
	}

	domainNames, domainDiags := helpers.StringSliceToList(ctx, segment.DomainNames)
	tcpRanges, tcpRangeDiags := helpers.StringSliceToList(ctx, segment.TCPPortRanges)
	udpRanges, udpRangeDiags := helpers.StringSliceToList(ctx, segment.UDPPortRanges)
	tcpPorts, tcpPortsDiags := helpers.FlattenNetworkPorts(ctx, segment.TCPAppPortRange)
	udpPorts, udpPortsDiags := helpers.FlattenNetworkPorts(ctx, segment.UDPAppPortRange)
	serverGroups, sgDiags := helpers.FlattenServerGroups(ctx, segment.ServerGroups)
	zpnER, zpnDiags := helpers.FlattenZPNERID(ctx, segment.ZPNERID)

	shareTo := make([]string, 0)
	if details := segment.SharedMicrotenantDetails; details.SharedToMicrotenants != nil {
		for _, mt := range details.SharedToMicrotenants {
			if strings.TrimSpace(mt.ID) != "" {
				shareTo = append(shareTo, mt.ID)
			}
		}
	}

	var shareSet types.Set
	if len(shareTo) == 0 {
		shareSet = types.SetNull(types.StringType)
	} else {
		shareSetValue, shareDiags := types.SetValueFrom(ctx, types.StringType, shareTo)
		if shareDiags.HasError() {
			return ApplicationSegmentModel{}, shareDiags
		}
		shareSet = shareSetValue
	}

	diags := diag.Diagnostics{}
	diags.Append(domainDiags...)
	diags.Append(tcpRangeDiags...)
	diags.Append(udpRangeDiags...)
	diags.Append(tcpPortsDiags...)
	diags.Append(udpPortsDiags...)
	diags.Append(sgDiags...)
	diags.Append(zpnDiags...)
	if diags.HasError() {
		return ApplicationSegmentModel{}, diags
	}

	state := ApplicationSegmentModel{
		ID:                        helpers.StringValueOrNull(segment.ID),
		Name:                      helpers.StringValueOrNull(segment.Name),
		Description:               helpers.StringValueOrNull(segment.Description),
		SegmentGroupID:            helpers.StringValueOrNull(segment.SegmentGroupID),
		SegmentGroupName:          helpers.StringValueOrNull(segment.SegmentGroupName),
		BypassType:                helpers.StringValueOrNull(segment.BypassType),
		BypassOnReauth:            types.BoolValue(segment.BypassOnReauth),
		ConfigSpace:               helpers.StringValueOrNull(segment.ConfigSpace),
		DomainNames:               domainNames,
		DoubleEncrypt:             types.BoolValue(segment.DoubleEncrypt),
		Enabled:                   types.BoolValue(segment.Enabled),
		InspectTrafficWithZia:     types.BoolValue(segment.InspectTrafficWithZia),
		PassiveHealthEnabled:      types.BoolValue(segment.PassiveHealthEnabled),
		HealthCheckType:           helpers.StringValueOrNull(segment.HealthCheckType),
		HealthReporting:           helpers.StringValueOrNull(segment.HealthReporting),
		IcmpAccessType:            helpers.StringValueOrNull(segment.IcmpAccessType),
		IPAnchored:                types.BoolValue(segment.IpAnchored),
		FqdnDnsCheck:              types.BoolValue(segment.FQDNDnsCheck),
		SelectConnectorCloseToApp: types.BoolValue(segment.SelectConnectorCloseToApp),
		UseInDrMode:               types.BoolValue(segment.UseInDrMode),
		IsIncompleteDRConfig:      types.BoolValue(segment.IsIncompleteDRConfig),
		IsCnameEnabled:            types.BoolValue(segment.IsCnameEnabled),
		TCPKeepAlive:              helpers.StringValueOrNull(segment.TCPKeepAlive),
		MicroTenantID:             helpers.StringValueOrNull(segment.MicroTenantID),
		MicroTenantName:           helpers.StringValueOrNull(segment.MicroTenantName),
		ShareToMicrotenants:       shareSet,
		TCPPortRanges:             tcpRanges,
		UDPPortRanges:             udpRanges,
		TCPPortRange:              tcpPorts,
		UDPPortRange:              udpPorts,
		ServerGroups:              serverGroups,
		ZpnERID:                   zpnER,
		MatchStyle:                helpers.StringValueOrNull(segment.MatchStyle),
		APIProtectionEnabled:      types.BoolValue(segment.APIProtectionEnabled),
	}

	return state, diags
}

func (r *ApplicationSegmentResource) shareApplicationSegment(ctx context.Context, service *zscaler.Service, appID string, share types.Set, microTenantID types.String) diag.Diagnostics {
	var diags diag.Diagnostics

	shareValues, shareDiags := helpers.SetValueToStringSlice(ctx, share)
	diags.Append(shareDiags...)
	if diags.HasError() {
		return diags
	}

	if len(shareValues) == 0 {
		return diags
	}

	mtID := strings.TrimSpace(microTenantID.ValueString())

	shareReq := applicationsegment_share.AppSegmentSharedToMicrotenant{
		ApplicationID:       appID,
		ShareToMicrotenants: shareValues,
		MicroTenantID:       mtID,
	}

	if _, err := applicationsegment_share.AppSegmentMicrotenantShare(ctx, service, appID, shareReq); err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Failed to share application segment with microtenants: %v", err))
	}

	return diags
}

func validateSegmentGroup(segmentGroupID types.String) error {
	if segmentGroupID.IsUnknown() || segmentGroupID.IsNull() {
		return fmt.Errorf("segment_group_id must be specified")
	}
	if strings.TrimSpace(segmentGroupID.ValueString()) == "" {
		return fmt.Errorf("segment_group_id must be specified")
	}
	return nil
}

func validateSegmentGroupOnUpdate(segmentGroupID types.String) error {
	if segmentGroupID.IsNull() || segmentGroupID.IsUnknown() {
		return nil
	}
	if strings.TrimSpace(segmentGroupID.ValueString()) == "" {
		return fmt.Errorf("segment_group_id must not be empty when updating")
	}
	return nil
}
