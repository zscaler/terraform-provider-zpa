package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentbrowseraccess"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &ApplicationSegmentBrowserAccessResource{}
	_ resource.ResourceWithConfigure   = &ApplicationSegmentBrowserAccessResource{}
	_ resource.ResourceWithImportState = &ApplicationSegmentBrowserAccessResource{}
)

func NewApplicationSegmentBrowserAccessResource() resource.Resource {
	return &ApplicationSegmentBrowserAccessResource{}
}

type ApplicationSegmentBrowserAccessResource struct {
	client *client.Client
}

type ApplicationSegmentBrowserAccessModel struct {
	ID                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	SegmentGroupID            types.String `tfsdk:"segment_group_id"`
	SegmentGroupName          types.String `tfsdk:"segment_group_name"`
	BypassType                types.String `tfsdk:"bypass_type"`
	ConfigSpace               types.String `tfsdk:"config_space"`
	Description               types.String `tfsdk:"description"`
	DomainNames               types.List   `tfsdk:"domain_names"`
	DoubleEncrypt             types.Bool   `tfsdk:"double_encrypt"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
	PassiveHealthEnabled      types.Bool   `tfsdk:"passive_health_enabled"`
	HealthCheckType           types.String `tfsdk:"health_check_type"`
	HealthReporting           types.String `tfsdk:"health_reporting"`
	IcmpAccessType            types.String `tfsdk:"icmp_access_type"`
	IPAnchored                types.Bool   `tfsdk:"ip_anchored"`
	IsCnameEnabled            types.Bool   `tfsdk:"is_cname_enabled"`
	SelectConnectorCloseToApp types.Bool   `tfsdk:"select_connector_close_to_app"`
	UseInDrMode               types.Bool   `tfsdk:"use_in_dr_mode"`
	IsIncompleteDRConfig      types.Bool   `tfsdk:"is_incomplete_dr_config"`
	APIProtectionEnabled      types.Bool   `tfsdk:"api_protection_enabled"`
	FqdndnsCheck              types.Bool   `tfsdk:"fqdn_dns_check"`
	TCPKeepAlive              types.String `tfsdk:"tcp_keep_alive"`
	MatchStyle                types.String `tfsdk:"match_style"`
	MicroTenantID             types.String `tfsdk:"microtenant_id"`
	MicroTenantName           types.String `tfsdk:"microtenant_name"`
	TCPPortRanges             types.List   `tfsdk:"tcp_port_ranges"`
	UDPPortRanges             types.List   `tfsdk:"udp_port_ranges"`
	TCPPortRange              types.List   `tfsdk:"tcp_port_range"`
	UDPPortRange              types.List   `tfsdk:"udp_port_range"`
	ClientlessApps            types.List   `tfsdk:"clientless_apps"`
	ServerGroups              types.List   `tfsdk:"server_groups"`
	ZpnERID                   types.List   `tfsdk:"zpn_er_id"`
}

type clientlessAppModel struct {
	ID                  types.String `tfsdk:"id"`
	AppID               types.String `tfsdk:"app_id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	MicroTenantID       types.String `tfsdk:"microtenant_id"`
	MicroTenantName     types.String `tfsdk:"microtenant_name"`
	AllowOptions        types.Bool   `tfsdk:"allow_options"`
	ApplicationPort     types.String `tfsdk:"application_port"`
	ApplicationProtocol types.String `tfsdk:"application_protocol"`
	CertificateID       types.String `tfsdk:"certificate_id"`
	CertificateName     types.String `tfsdk:"certificate_name"`
	Cname               types.String `tfsdk:"cname"`
	Domain              types.String `tfsdk:"domain"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	Hidden              types.Bool   `tfsdk:"hidden"`
	LocalDomain         types.String `tfsdk:"local_domain"`
	Path                types.String `tfsdk:"path"`
	TrustUntrustedCert  types.Bool   `tfsdk:"trust_untrusted_cert"`
	ExtLabel            types.String `tfsdk:"ext_label"`
	ExtDomain           types.String `tfsdk:"ext_domain"`
}

type zpnERModel struct {
	ID types.List `tfsdk:"id"`
}

func (r *ApplicationSegmentBrowserAccessResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_segment_browser_access"
}

func (r *ApplicationSegmentBrowserAccessResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Browser Access application segment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the browser access application.",
			},
			"segment_group_id": schema.StringAttribute{
				Required:    true,
				Description: "Segment group identifier associated with the application.",
			},
			"segment_group_name": schema.StringAttribute{
				Computed: true,
			},
			"bypass_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("NEVER"),
			},
			"config_space": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("DEFAULT"),
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"domain_names": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "List of domains and IPs for the application segment.",
			},
			"double_encrypt":                schema.BoolAttribute{Optional: true, Computed: true},
			"enabled":                       schema.BoolAttribute{Optional: true, Computed: true},
			"passive_health_enabled":        schema.BoolAttribute{Optional: true, Computed: true},
			"health_check_type":             schema.StringAttribute{Optional: true, Computed: true},
			"health_reporting":              schema.StringAttribute{Optional: true, Computed: true},
			"icmp_access_type":              schema.StringAttribute{Optional: true, Computed: true},
			"ip_anchored":                   schema.BoolAttribute{Optional: true, Computed: true},
			"is_cname_enabled":              schema.BoolAttribute{Optional: true, Computed: true},
			"select_connector_close_to_app": schema.BoolAttribute{Optional: true, Computed: true},
			"use_in_dr_mode":                schema.BoolAttribute{Optional: true, Computed: true},
			"is_incomplete_dr_config":       schema.BoolAttribute{Optional: true, Computed: true},
			"api_protection_enabled":        schema.BoolAttribute{Optional: true, Computed: true},
			"fqdn_dns_check":                schema.BoolAttribute{Optional: true, Computed: true},
			"tcp_keep_alive":                schema.StringAttribute{Optional: true, Computed: true},
			"match_style":                   schema.StringAttribute{Optional: true, Computed: true},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant identifier used to scope API requests.",
			},
			"microtenant_name": schema.StringAttribute{Computed: true},
			"tcp_port_ranges": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"udp_port_ranges": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"zpn_er_id": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.SetAttribute{ElementType: types.StringType, Optional: true, Computed: true},
					},
				},
			},
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
			// clientless_apps: TypeList in SDKv2, Required
			// Using ListNestedBlock for block syntax support
			"clientless_apps": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":                   schema.StringAttribute{Computed: true},
						"app_id":               schema.StringAttribute{Computed: true},
						"name":                 schema.StringAttribute{Required: true},
						"description":          schema.StringAttribute{Optional: true},
						"microtenant_id":       schema.StringAttribute{Optional: true},
						"microtenant_name":     schema.StringAttribute{Computed: true},
						"allow_options":        schema.BoolAttribute{Optional: true, Computed: true},
						"application_port":     schema.StringAttribute{Required: true},
						"application_protocol": schema.StringAttribute{Required: true},
						"certificate_id":       schema.StringAttribute{Optional: true},
						"certificate_name":     schema.StringAttribute{Computed: true},
						"cname":                schema.StringAttribute{Computed: true},
						"domain":               schema.StringAttribute{Optional: true},
						"enabled":              schema.BoolAttribute{Optional: true, Computed: true},
						"hidden":               schema.BoolAttribute{Optional: true, Computed: true},
						"local_domain":         schema.StringAttribute{Optional: true},
						"path":                 schema.StringAttribute{Optional: true},
						"trust_untrusted_cert": schema.BoolAttribute{Optional: true, Computed: true},
						"ext_label":            schema.StringAttribute{Optional: true},
						"ext_domain":           schema.StringAttribute{Optional: true},
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

func (r *ApplicationSegmentBrowserAccessResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationSegmentBrowserAccessResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider was not configured before use.")
		return
	}

	var plan ApplicationSegmentBrowserAccessModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if diag := validateClientlessApps(ctx, plan); diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	service := r.serviceForMicrotenant(plan.MicroTenantID)

	payload, diags := r.expandBrowserAccess(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if diag := helpers.ValidateAppPorts(payload.SelectConnectorCloseToApp, payload.UDPAppPortRange, payload.UDPPortRanges); diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	created, _, err := applicationsegmentbrowseraccess.Create(ctx, service, payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create browser access application: %v", err))
		return
	}

	tflog.Info(ctx, "Created browser access application", map[string]any{"id": created.ID})

	state, diags := r.readBrowserAccess(ctx, service, created.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationSegmentBrowserAccessResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider was not configured before use.")
		return
	}

	var state ApplicationSegmentBrowserAccessModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicroTenantID)
	newState, diags := r.readBrowserAccess(ctx, service, state.ID.ValueString())
	if diags.HasError() {
		notFound := false
		for _, d := range diags {
			if d.Severity() == diag.SeverityError && strings.Contains(strings.ToLower(d.Detail()), "not found") {
				notFound = true
				break
			}
		}
		if notFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ApplicationSegmentBrowserAccessResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider was not configured before use.")
		return
	}

	var plan ApplicationSegmentBrowserAccessModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if diag := validateClientlessApps(ctx, plan); diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	service := r.serviceForMicrotenant(plan.MicroTenantID)

	payload, diags := r.expandBrowserAccess(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	existing, _, err := applicationsegmentbrowseraccess.Get(ctx, service, plan.ID.ValueString())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve existing browser access application: %v", err))
		return
	}

	r.reconcileClientlessAppIDs(existing, &payload)

	if diag := helpers.ValidateAppPorts(payload.SelectConnectorCloseToApp, payload.UDPAppPortRange, payload.UDPPortRanges); diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	if _, err := applicationsegmentbrowseraccess.Update(ctx, service, existing.ID, &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update browser access application: %v", err))
		return
	}

	state, diags := r.readBrowserAccess(ctx, service, existing.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationSegmentBrowserAccessResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider was not configured before use.")
		return
	}

	var state ApplicationSegmentBrowserAccessModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicroTenantID)

	if _, err := applicationsegmentbrowseraccess.Delete(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete browser access application: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *ApplicationSegmentBrowserAccessResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
}

func (r *ApplicationSegmentBrowserAccessResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && microtenantID.ValueString() != "" {
		service = service.WithMicroTenant(microtenantID.ValueString())
	}
	return service
}

func (r *ApplicationSegmentBrowserAccessResource) expandBrowserAccess(ctx context.Context, plan ApplicationSegmentBrowserAccessModel) (applicationsegmentbrowseraccess.BrowserAccess, diag.Diagnostics) {
	var diags diag.Diagnostics

	domainNames, domainDiags := helpers.ListValueToStringSlice(ctx, plan.DomainNames)
	diags.Append(domainDiags...)

	tcpRanges, tcpRangesDiags := helpers.ListValueToStringSlice(ctx, plan.TCPPortRanges)
	diags.Append(tcpRangesDiags...)

	udpRanges, udpRangesDiags := helpers.ListValueToStringSlice(ctx, plan.UDPPortRanges)
	diags.Append(udpRangesDiags...)

	tcpPorts, tcpDiags := helpers.ExpandNetworkPorts(ctx, plan.TCPPortRange)
	diags.Append(tcpDiags...)

	udpPorts, udpDiags := helpers.ExpandNetworkPorts(ctx, plan.UDPPortRange)
	diags.Append(udpDiags...)

	clientlessApps, clientlessDiags := r.expandClientlessApps(ctx, plan.ClientlessApps)
	diags.Append(clientlessDiags...)

	serverGroups, sgDiags := helpers.ExpandServerGroups(ctx, plan.ServerGroups)
	diags.Append(sgDiags...)

	zpnER, zpnDiags := expandZpnERID(ctx, plan.ZpnERID)
	diags.Append(zpnDiags...)

	payload := applicationsegmentbrowseraccess.BrowserAccess{
		ID:                        strings.TrimSpace(plan.ID.ValueString()),
		Name:                      plan.Name.ValueString(),
		SegmentGroupID:            plan.SegmentGroupID.ValueString(),
		SegmentGroupName:          plan.SegmentGroupName.ValueString(),
		BypassType:                plan.BypassType.ValueString(),
		ConfigSpace:               plan.ConfigSpace.ValueString(),
		Description:               plan.Description.ValueString(),
		DomainNames:               domainNames,
		DoubleEncrypt:             helpers.BoolValue(plan.DoubleEncrypt, false),
		Enabled:                   helpers.BoolValue(plan.Enabled, false),
		PassiveHealthEnabled:      helpers.BoolValue(plan.PassiveHealthEnabled, false),
		HealthCheckType:           plan.HealthCheckType.ValueString(),
		HealthReporting:           plan.HealthReporting.ValueString(),
		ICMPAccessType:            plan.IcmpAccessType.ValueString(),
		IPAnchored:                helpers.BoolValue(plan.IPAnchored, false),
		IsCnameEnabled:            helpers.BoolValue(plan.IsCnameEnabled, false),
		SelectConnectorCloseToApp: helpers.BoolValue(plan.SelectConnectorCloseToApp, false),
		UseInDrMode:               helpers.BoolValue(plan.UseInDrMode, false),
		IsIncompleteDRConfig:      helpers.BoolValue(plan.IsIncompleteDRConfig, false),
		APIProtectionEnabled:      helpers.BoolValue(plan.APIProtectionEnabled, false),
		FQDNDnsCheck:              helpers.BoolValue(plan.FqdndnsCheck, false),
		TCPKeepAlive:              plan.TCPKeepAlive.ValueString(),
		MatchStyle:                plan.MatchStyle.ValueString(),
		MicroTenantID:             plan.MicroTenantID.ValueString(),
		TCPPortRanges:             tcpRanges,
		UDPPortRanges:             udpRanges,
		TCPAppPortRange:           tcpPorts,
		UDPAppPortRange:           udpPorts,
		ClientlessApps:            clientlessApps,
		AppServerGroups:           serverGroups,
		ZPNERID:                   zpnER,
	}

	return payload, diags
}

func (r *ApplicationSegmentBrowserAccessResource) expandClientlessApps(ctx context.Context, list types.List) ([]applicationsegmentbrowseraccess.ClientlessApps, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Missing clientless_apps", "At least one clientless application must be specified.")}
	}

	var models []clientlessAppModel
	var diags diag.Diagnostics
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil, diags
	}

	if len(models) == 0 {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Missing clientless_apps", "At least one clientless application must be specified.")}
	}

	apps := make([]applicationsegmentbrowseraccess.ClientlessApps, 0, len(models))
	for _, app := range models {
		apps = append(apps, applicationsegmentbrowseraccess.ClientlessApps{
			ID:                  app.ID.ValueString(),
			AppID:               app.AppID.ValueString(),
			Name:                app.Name.ValueString(),
			Description:         app.Description.ValueString(),
			AllowOptions:        helpers.BoolValue(app.AllowOptions, false),
			ApplicationPort:     app.ApplicationPort.ValueString(),
			ApplicationProtocol: app.ApplicationProtocol.ValueString(),
			CertificateID:       app.CertificateID.ValueString(),
			CertificateName:     app.CertificateName.ValueString(),
			Cname:               app.Cname.ValueString(),
			Domain:              app.Domain.ValueString(),
			Enabled:             helpers.BoolValue(app.Enabled, false),
			Hidden:              helpers.BoolValue(app.Hidden, false),
			LocalDomain:         app.LocalDomain.ValueString(),
			Path:                app.Path.ValueString(),
			TrustUntrustedCert:  helpers.BoolValue(app.TrustUntrustedCert, false),
			MicroTenantID:       app.MicroTenantID.ValueString(),
			MicroTenantName:     app.MicroTenantName.ValueString(),
			ExtLabel:            app.ExtLabel.ValueString(),
			ExtDomain:           app.ExtDomain.ValueString(),
		})
	}

	return apps, diags
}

func expandZpnERID(ctx context.Context, list types.List) (*common.ZPNERID, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, diag.Diagnostics{}
	}

	var models []zpnERModel
	var diags diag.Diagnostics
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil, diags
	}

	for _, model := range models {
		if model.ID.IsNull() || model.ID.IsUnknown() {
			continue
		}
		var ids []string
		diags.Append(model.ID.ElementsAs(ctx, &ids, false)...)
		if diags.HasError() {
			return nil, diags
		}
		if len(ids) == 0 {
			continue
		}
		for _, id := range ids {
			id = strings.TrimSpace(id)
			if id != "" {
				return &common.ZPNERID{ID: id}, diags
			}
		}
	}

	return nil, diags
}

func (r *ApplicationSegmentBrowserAccessResource) readBrowserAccess(ctx context.Context, service *zscaler.Service, id string) (ApplicationSegmentBrowserAccessModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var state ApplicationSegmentBrowserAccessModel

	app, _, err := applicationsegmentbrowseraccess.Get(ctx, service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			diags.AddError("Not Found", fmt.Sprintf("Browser access application %s was not found", id))
			return state, diags
		}
		diags.AddError("Client Error", fmt.Sprintf("Unable to read browser access application: %v", err))
		return state, diags
	}

	domainNames, domainsDiags := helpers.StringSliceToList(ctx, app.DomainNames)
	diags.Append(domainsDiags...)

	tcpRanges, tcpRangesDiags := helpers.StringSliceToList(ctx, app.TCPPortRanges)
	diags.Append(tcpRangesDiags...)

	udpRanges, udpRangesDiags := helpers.StringSliceToList(ctx, app.UDPPortRanges)
	diags.Append(udpRangesDiags...)

	tcpPorts, tcpDiags := helpers.FlattenNetworkPorts(ctx, app.TCPAppPortRange)
	diags.Append(tcpDiags...)

	udpPorts, udpDiags := helpers.FlattenNetworkPorts(ctx, app.UDPAppPortRange)
	diags.Append(udpDiags...)

	serverGroups, sgDiags := helpers.FlattenServerGroups(ctx, app.AppServerGroups)
	diags.Append(sgDiags...)

	clientlessApps, clientlessDiags := flattenClientlessApps(ctx, app.ClientlessApps)
	diags.Append(clientlessDiags...)

	zpnER, zpnDiags := flattenZpnERID(ctx, app.ZPNERID)
	diags.Append(zpnDiags...)

	state = ApplicationSegmentBrowserAccessModel{
		ID:                        helpers.StringValueOrNull(app.ID),
		Name:                      helpers.StringValueOrNull(app.Name),
		SegmentGroupID:            helpers.StringValueOrNull(app.SegmentGroupID),
		SegmentGroupName:          helpers.StringValueOrNull(app.SegmentGroupName),
		BypassType:                helpers.StringValueOrNull(app.BypassType),
		ConfigSpace:               helpers.StringValueOrNull(app.ConfigSpace),
		Description:               helpers.StringValueOrNull(app.Description),
		DomainNames:               domainNames,
		DoubleEncrypt:             types.BoolValue(app.DoubleEncrypt),
		Enabled:                   types.BoolValue(app.Enabled),
		PassiveHealthEnabled:      types.BoolValue(app.PassiveHealthEnabled),
		HealthCheckType:           helpers.StringValueOrNull(app.HealthCheckType),
		HealthReporting:           helpers.StringValueOrNull(app.HealthReporting),
		IcmpAccessType:            helpers.StringValueOrNull(app.ICMPAccessType),
		IPAnchored:                types.BoolValue(app.IPAnchored),
		IsCnameEnabled:            types.BoolValue(app.IsCnameEnabled),
		SelectConnectorCloseToApp: types.BoolValue(app.SelectConnectorCloseToApp),
		UseInDrMode:               types.BoolValue(app.UseInDrMode),
		IsIncompleteDRConfig:      types.BoolValue(app.IsIncompleteDRConfig),
		APIProtectionEnabled:      types.BoolValue(app.APIProtectionEnabled),
		FqdndnsCheck:              types.BoolValue(app.FQDNDnsCheck),
		TCPKeepAlive:              helpers.StringValueOrNull(app.TCPKeepAlive),
		MatchStyle:                helpers.StringValueOrNull(app.MatchStyle),
		MicroTenantID:             helpers.StringValueOrNull(app.MicroTenantID),
		MicroTenantName:           helpers.StringValueOrNull(app.MicroTenantName),
		TCPPortRanges:             tcpRanges,
		UDPPortRanges:             udpRanges,
		TCPPortRange:              tcpPorts,
		UDPPortRange:              udpPorts,
		ClientlessApps:            clientlessApps,
		ServerGroups:              serverGroups,
		ZpnERID:                   zpnER,
	}

	return state, diags
}

func flattenClientlessApps(ctx context.Context, apps []applicationsegmentbrowseraccess.ClientlessApps) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id":                   types.StringType,
		"app_id":               types.StringType,
		"name":                 types.StringType,
		"description":          types.StringType,
		"microtenant_id":       types.StringType,
		"microtenant_name":     types.StringType,
		"allow_options":        types.BoolType,
		"application_port":     types.StringType,
		"application_protocol": types.StringType,
		"certificate_id":       types.StringType,
		"certificate_name":     types.StringType,
		"cname":                types.StringType,
		"domain":               types.StringType,
		"enabled":              types.BoolType,
		"hidden":               types.BoolType,
		"local_domain":         types.StringType,
		"path":                 types.StringType,
		"trust_untrusted_cert": types.BoolType,
		"ext_label":            types.StringType,
		"ext_domain":           types.StringType,
	}

	if len(apps) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(apps))
	var diags diag.Diagnostics
	for _, app := range apps {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":                   helpers.StringValueOrNull(app.ID),
			"app_id":               helpers.StringValueOrNull(app.AppID),
			"name":                 helpers.StringValueOrNull(app.Name),
			"description":          helpers.StringValueOrNull(app.Description),
			"microtenant_id":       helpers.StringValueOrNull(app.MicroTenantID),
			"microtenant_name":     helpers.StringValueOrNull(app.MicroTenantName),
			"allow_options":        types.BoolValue(app.AllowOptions),
			"application_port":     helpers.StringValueOrNull(app.ApplicationPort),
			"application_protocol": helpers.StringValueOrNull(app.ApplicationProtocol),
			"certificate_id":       helpers.StringValueOrNull(app.CertificateID),
			"certificate_name":     helpers.StringValueOrNull(app.CertificateName),
			"cname":                helpers.StringValueOrNull(app.Cname),
			"domain":               helpers.StringValueOrNull(app.Domain),
			"enabled":              types.BoolValue(app.Enabled),
			"hidden":               types.BoolValue(app.Hidden),
			"local_domain":         helpers.StringValueOrNull(app.LocalDomain),
			"path":                 helpers.StringValueOrNull(app.Path),
			"trust_untrusted_cert": types.BoolValue(app.TrustUntrustedCert),
			"ext_label":            helpers.StringValueOrNull(app.ExtLabel),
			"ext_domain":           helpers.StringValueOrNull(app.ExtDomain),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}

func flattenZpnERID(ctx context.Context, value *common.ZPNERID) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id": types.SetType{ElemType: types.StringType},
	}
	if value == nil || strings.TrimSpace(value.ID) == "" {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	idSet, diags := types.SetValueFrom(ctx, types.StringType, []string{value.ID})
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"id": idSet,
	})
	diags.Append(objDiags...)

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	diags.Append(listDiags...)
	return list, diags
}

func (r *ApplicationSegmentBrowserAccessResource) reconcileClientlessAppIDs(existing *applicationsegmentbrowseraccess.BrowserAccess, payload *applicationsegmentbrowseraccess.BrowserAccess) {
	if existing == nil || payload == nil {
		return
	}

	existingByName := make(map[string]applicationsegmentbrowseraccess.ClientlessApps)
	for _, app := range existing.ClientlessApps {
		existingByName[strings.ToLower(app.Name)] = app
	}

	for i := range payload.ClientlessApps {
		app := &payload.ClientlessApps[i]
		if app.ID != "" {
			continue
		}
		if match, ok := existingByName[strings.ToLower(app.Name)]; ok {
			app.ID = match.ID
			app.AppID = match.AppID
		} else if i < len(existing.ClientlessApps) {
			app.ID = existing.ClientlessApps[i].ID
			app.AppID = existing.ClientlessApps[i].AppID
		}
	}
}

func validateClientlessApps(ctx context.Context, plan ApplicationSegmentBrowserAccessModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if plan.ClientlessApps.IsNull() || plan.ClientlessApps.IsUnknown() {
		diags.AddError("Missing clientless_apps", "At least one clientless application must be specified.")
		return diags
	}

	var apps []clientlessAppModel
	di := plan.ClientlessApps.ElementsAs(ctx, &apps, false)
	if di.HasError() {
		diags.Append(di...)
		return diags
	}

	for i, app := range apps {
		hasExt := (!app.ExtLabel.IsNull() && app.ExtLabel.ValueString() != "") || (!app.ExtDomain.IsNull() && app.ExtDomain.ValueString() != "")
		hasCert := !app.CertificateID.IsNull() && app.CertificateID.ValueString() != ""
		if hasExt && hasCert {
			diags.AddError("Invalid clientless app configuration", fmt.Sprintf("clientless_apps[%d]: certificate_id cannot be set when ext_label or ext_domain is specified", i))
		}
	}

	return diags
}
