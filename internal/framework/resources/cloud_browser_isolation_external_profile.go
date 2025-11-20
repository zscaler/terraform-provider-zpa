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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	cbiprofilecontroller "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbiprofilecontroller"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &CBIExternalProfileResource{}
	_ resource.ResourceWithConfigure   = &CBIExternalProfileResource{}
	_ resource.ResourceWithImportState = &CBIExternalProfileResource{}
)

func NewCBIExternalProfileResource() resource.Resource {
	return &CBIExternalProfileResource{}
}

type CBIExternalProfileResource struct {
	client *client.Client
}

type CBIExternalProfileModel struct {
	ID               types.String               `tfsdk:"id"`
	Name             types.String               `tfsdk:"name"`
	Description      types.String               `tfsdk:"description"`
	BannerID         types.String               `tfsdk:"banner_id"`
	RegionIDs        types.Set                  `tfsdk:"region_ids"`
	CertificateIDs   types.Set                  `tfsdk:"certificate_ids"`
	UserExperience   []CBIUserExperienceModel   `tfsdk:"user_experience"`
	SecurityControls []CBISecurityControlsModel `tfsdk:"security_controls"`
	DebugMode        []CBIDebugModeModel        `tfsdk:"debug_mode"`
}

type CBIUserExperienceModel struct {
	ZGPU                types.Bool             `tfsdk:"zgpu"`
	BrowserInBrowser    types.Bool             `tfsdk:"browser_in_browser"`
	PersistIsolationBar types.Bool             `tfsdk:"persist_isolation_bar"`
	Translate           types.Bool             `tfsdk:"translate"`
	SessionPersistence  types.Bool             `tfsdk:"session_persistence"`
	ForwardToZia        []CBIForwardToZiaModel `tfsdk:"forward_to_zia"`
}

type CBIForwardToZiaModel struct {
	Enabled        types.Bool   `tfsdk:"enabled"`
	OrganizationID types.String `tfsdk:"organization_id"`
	CloudName      types.String `tfsdk:"cloud_name"`
	PacFileURL     types.String `tfsdk:"pac_file_url"`
}

type CBISecurityControlsModel struct {
	CopyPaste          types.String        `tfsdk:"copy_paste"`
	UploadDownload     types.String        `tfsdk:"upload_download"`
	DocumentViewer     types.Bool          `tfsdk:"document_viewer"`
	LocalRender        types.Bool          `tfsdk:"local_render"`
	AllowPrinting      types.Bool          `tfsdk:"allow_printing"`
	RestrictKeystrokes types.Bool          `tfsdk:"restrict_keystrokes"`
	FlattenedPDF       types.Bool          `tfsdk:"flattened_pdf"`
	DeepLink           []CBIDeepLinkModel  `tfsdk:"deep_link"`
	Watermark          []CBIWatermarkModel `tfsdk:"watermark"`
}

type CBIDeepLinkModel struct {
	Enabled      types.Bool `tfsdk:"enabled"`
	Applications types.Set  `tfsdk:"applications"`
}

type CBIWatermarkModel struct {
	Enabled       types.Bool   `tfsdk:"enabled"`
	ShowUserID    types.Bool   `tfsdk:"show_user_id"`
	ShowTimestamp types.Bool   `tfsdk:"show_timestamp"`
	ShowMessage   types.Bool   `tfsdk:"show_message"`
	Message       types.String `tfsdk:"message"`
}

type CBIDebugModeModel struct {
	Allowed      types.Bool   `tfsdk:"allowed"`
	FilePassword types.String `tfsdk:"file_password"`
}

func (r *CBIExternalProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_browser_isolation_external_profile"
}

func (r *CBIExternalProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Cloud Browser Isolation external profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the external profile.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the profile.",
			},
			"banner_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifier of the banner associated with this profile.",
			},
			"region_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "List of region IDs associated with the profile. Must include at least two IDs.",
			},
			"certificate_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of certificate IDs associated with the profile.",
			},
		},
		Blocks: map[string]schema.Block{
			"user_experience": schema.ListNestedBlock{
				Description: "User experience configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"zgpu":                  schema.BoolAttribute{Optional: true, Computed: true},
						"browser_in_browser":    schema.BoolAttribute{Optional: true, Computed: true},
						"persist_isolation_bar": schema.BoolAttribute{Optional: true, Computed: true},
						"translate":             schema.BoolAttribute{Optional: true, Computed: true},
						"session_persistence":   schema.BoolAttribute{Optional: true, Computed: true},
					},
					Blocks: map[string]schema.Block{
						"forward_to_zia": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"enabled":         schema.BoolAttribute{Optional: true, Computed: true},
									"organization_id": schema.StringAttribute{Optional: true, Computed: true},
									"cloud_name":      schema.StringAttribute{Optional: true, Computed: true},
									"pac_file_url":    schema.StringAttribute{Optional: true, Computed: true},
								},
							},
						},
					},
				},
			},
			"security_controls": schema.ListNestedBlock{
				Description: "Security control configuration for isolated sessions.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"copy_paste":          schema.StringAttribute{Optional: true, Computed: true},
						"upload_download":     schema.StringAttribute{Optional: true, Computed: true},
						"document_viewer":     schema.BoolAttribute{Optional: true, Computed: true},
						"local_render":        schema.BoolAttribute{Optional: true, Computed: true},
						"allow_printing":      schema.BoolAttribute{Optional: true, Computed: true},
						"restrict_keystrokes": schema.BoolAttribute{Optional: true, Computed: true},
						"flattened_pdf":       schema.BoolAttribute{Optional: true, Computed: true},
					},
					Blocks: map[string]schema.Block{
						"deep_link": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{Optional: true, Computed: true},
									"applications": schema.SetAttribute{
										ElementType: types.StringType,
										Optional:    true,
										Computed:    true,
									},
								},
							},
						},
						"watermark": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"enabled":        schema.BoolAttribute{Optional: true, Computed: true},
									"show_user_id":   schema.BoolAttribute{Optional: true, Computed: true},
									"show_timestamp": schema.BoolAttribute{Optional: true, Computed: true},
									"show_message":   schema.BoolAttribute{Optional: true, Computed: true},
									"message":        schema.StringAttribute{Optional: true, Computed: true},
								},
							},
						},
					},
				},
			},
			"debug_mode": schema.ListNestedBlock{
				Description: "Debug mode configuration for isolated sessions.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"allowed":       schema.BoolAttribute{Optional: true, Computed: true},
						"file_password": schema.StringAttribute{Optional: true, Computed: true},
					},
				},
			},
		},
	}
}

func (r *CBIExternalProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CBIExternalProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan CBIExternalProfileModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diags := expandCBIExternalProfile(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(payload.RegionIDs) < 2 {
		resp.Diagnostics.AddError("Validation Error", "expected region_ids to contain at least 2 items")
		return
	}

	// Clear Regions, Certificates, and Banner before sending - API only accepts IDs
	payload.Regions = nil
	payload.Certificates = nil
	payload.Banner = nil

	created, _, err := cbiprofilecontroller.Create(ctx, r.client.Service, &payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create CBI external profile: %v", err))
		return
	}

	state, readDiags := r.readIntoState(ctx, created.ID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CBIExternalProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state CBIExternalProfileModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readIntoState(ctx, state.ID.ValueString())
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

func (r *CBIExternalProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan CBIExternalProfileModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.IsUnknown() {
		resp.Diagnostics.AddError("Validation Error", "id must be known during update")
		return
	}

	payload, diags := expandCBIExternalProfile(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(payload.RegionIDs) < 2 {
		resp.Diagnostics.AddError("Validation Error", "expected region_ids to contain at least 2 items")
		return
	}

	// Clear Regions, Certificates, and Banner before sending - API only accepts IDs
	payload.Regions = nil
	payload.Certificates = nil
	payload.Banner = nil

	if _, err := cbiprofilecontroller.Update(ctx, r.client.Service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update CBI external profile: %v", err))
		return
	}

	state, readDiags := r.readIntoState(ctx, plan.ID.ValueString())
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CBIExternalProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state CBIExternalProfileModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := cbiprofilecontroller.Delete(ctx, r.client.Service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete CBI external profile: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *CBIExternalProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
}

func (r *CBIExternalProfileResource) readIntoState(ctx context.Context, id string) (CBIExternalProfileModel, diag.Diagnostics) {
	profile, _, err := cbiprofilecontroller.Get(ctx, r.client.Service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return CBIExternalProfileModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("CBI external profile %s not found", id))}
		}
		return CBIExternalProfileModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read CBI external profile: %v", err))}
	}

	// Handle banner_id - API may return Banner object instead of BannerID
	if profile.BannerID == "" && profile.Banner != nil && profile.Banner.ID != "" {
		profile.BannerID = profile.Banner.ID
	}

	// Handle certificate_ids - API may return Certificates array instead of CertificateIDs
	if len(profile.CertificateIDs) == 0 && len(profile.Certificates) > 0 {
		for _, cert := range profile.Certificates {
			profile.CertificateIDs = append(profile.CertificateIDs, cert.ID)
		}
	}

	// Handle region_ids - API may return Regions array instead of RegionIDs
	if len(profile.RegionIDs) == 0 && len(profile.Regions) > 0 {
		for _, region := range profile.Regions {
			profile.RegionIDs = append(profile.RegionIDs, region.ID)
		}
	}

	state := CBIExternalProfileModel{
		ID:          helpers.StringValueOrNull(profile.ID),
		Name:        helpers.StringValueOrNull(profile.Name),
		Description: helpers.StringValueOrNull(profile.Description),
		BannerID:    helpers.StringValueOrNull(profile.BannerID),
	}

	regionIDs, regDiags := types.SetValueFrom(ctx, types.StringType, profile.RegionIDs)
	if regDiags.HasError() {
		return CBIExternalProfileModel{}, regDiags
	}
	state.RegionIDs = regionIDs

	certIDs, certDiags := types.SetValueFrom(ctx, types.StringType, profile.CertificateIDs)
	if certDiags.HasError() {
		return CBIExternalProfileModel{}, certDiags
	}
	state.CertificateIDs = certIDs

	ux, uxDiags := flattenCBIUserExperience(ctx, profile.UserExperience)
	if uxDiags.HasError() {
		return CBIExternalProfileModel{}, uxDiags
	}
	state.UserExperience = ux

	sc, scDiags := flattenCBISecurityControls(ctx, profile.SecurityControls)
	if scDiags.HasError() {
		return CBIExternalProfileModel{}, scDiags
	}
	state.SecurityControls = sc

	state.DebugMode = flattenCBIDebugMode(profile.DebugMode)

	return state, nil
}

func expandCBIExternalProfile(ctx context.Context, plan CBIExternalProfileModel) (cbiprofilecontroller.IsolationProfile, diag.Diagnostics) {
	var diags diag.Diagnostics

	regionIDs, regionDiags := helpers.SetValueToStringSlice(ctx, plan.RegionIDs)
	diags.Append(regionDiags...)

	certificateIDs, certDiags := helpers.SetValueToStringSlice(ctx, plan.CertificateIDs)
	diags.Append(certDiags...)

	if diags.HasError() {
		return cbiprofilecontroller.IsolationProfile{}, diags
	}

	payload := cbiprofilecontroller.IsolationProfile{
		ID:             helpers.StringValue(plan.ID),
		Name:           helpers.StringValue(plan.Name),
		Description:    helpers.StringValue(plan.Description),
		BannerID:       helpers.StringValue(plan.BannerID),
		RegionIDs:      regionIDs,
		CertificateIDs: certificateIDs,
	}

	if payload.BannerID != "" {
		payload.Banner = &cbiprofilecontroller.Banner{ID: payload.BannerID}
	}

	payload.Regions = make([]cbiprofilecontroller.Regions, 0, len(regionIDs))
	for _, id := range regionIDs {
		payload.Regions = append(payload.Regions, cbiprofilecontroller.Regions{ID: id})
	}

	payload.Certificates = make([]cbiprofilecontroller.Certificates, 0, len(certificateIDs))
	for _, id := range certificateIDs {
		payload.Certificates = append(payload.Certificates, cbiprofilecontroller.Certificates{ID: id})
	}

	if ux := expandCBIUserExperience(plan.UserExperience); ux != nil {
		payload.UserExperience = ux
	}

	if sc, scDiags := expandCBISecurityControls(ctx, plan.SecurityControls); scDiags.HasError() {
		diags.Append(scDiags...)
	} else {
		payload.SecurityControls = sc
	}

	if dm := expandCBIDebugMode(plan.DebugMode); dm != nil {
		payload.DebugMode = dm
	}

	return payload, diags
}

func expandCBIUserExperience(items []CBIUserExperienceModel) *cbiprofilecontroller.UserExperience {
	if len(items) == 0 {
		return nil
	}

	v := items[0]
	ux := &cbiprofilecontroller.UserExperience{
		ZGPU:                helpers.BoolValue(v.ZGPU, false),
		BrowserInBrowser:    helpers.BoolValue(v.BrowserInBrowser, false),
		PersistIsolationBar: helpers.BoolValue(v.PersistIsolationBar, false),
		Translate:           helpers.BoolValue(v.Translate, false),
		SessionPersistence:  helpers.BoolValue(v.SessionPersistence, false),
	}

	if len(v.ForwardToZia) > 0 {
		fz := v.ForwardToZia[0]
		ux.ForwardToZia = &cbiprofilecontroller.ForwardToZia{
			Enabled:        helpers.BoolValue(fz.Enabled, false),
			OrganizationID: helpers.StringValue(fz.OrganizationID),
			CloudName:      helpers.StringValue(fz.CloudName),
			PacFileUrl:     helpers.StringValue(fz.PacFileURL),
		}
	}

	return ux
}

func expandCBISecurityControls(ctx context.Context, items []CBISecurityControlsModel) (*cbiprofilecontroller.SecurityControls, diag.Diagnostics) {
	if len(items) == 0 {
		return nil, diag.Diagnostics{}
	}

	sc := items[0]
	result := &cbiprofilecontroller.SecurityControls{
		CopyPaste:          helpers.StringValue(sc.CopyPaste),
		UploadDownload:     helpers.StringValue(sc.UploadDownload),
		DocumentViewer:     helpers.BoolValue(sc.DocumentViewer, false),
		LocalRender:        helpers.BoolValue(sc.LocalRender, false),
		AllowPrinting:      helpers.BoolValue(sc.AllowPrinting, false),
		RestrictKeystrokes: helpers.BoolValue(sc.RestrictKeystrokes, false),
		FlattenedPdf:       helpers.BoolValue(sc.FlattenedPDF, false),
	}

	var diags diag.Diagnostics

	if len(sc.DeepLink) > 0 {
		dl := sc.DeepLink[0]
		applications, appDiags := helpers.SetValueToStringSlice(ctx, dl.Applications)
		diags.Append(appDiags...)
		result.DeepLink = &cbiprofilecontroller.DeepLink{
			Enabled:      helpers.BoolValue(dl.Enabled, false),
			Applications: applications,
		}
	}

	if len(sc.Watermark) > 0 {
		wm := sc.Watermark[0]
		result.Watermark = &cbiprofilecontroller.Watermark{
			Enabled:       helpers.BoolValue(wm.Enabled, false),
			ShowUserID:    helpers.BoolValue(wm.ShowUserID, false),
			ShowTimestamp: helpers.BoolValue(wm.ShowTimestamp, false),
			ShowMessage:   helpers.BoolValue(wm.ShowMessage, false),
			Message:       helpers.StringValue(wm.Message),
		}
	}

	return result, diags
}

func expandCBIDebugMode(items []CBIDebugModeModel) *cbiprofilecontroller.DebugMode {
	if len(items) == 0 {
		return nil
	}

	dm := items[0]
	return &cbiprofilecontroller.DebugMode{
		Allowed:      helpers.BoolValue(dm.Allowed, false),
		FilePassword: helpers.StringValue(dm.FilePassword),
	}
}

func flattenCBIUserExperience(ctx context.Context, ux *cbiprofilecontroller.UserExperience) ([]CBIUserExperienceModel, diag.Diagnostics) {
	if ux == nil {
		return nil, diag.Diagnostics{}
	}

	model := CBIUserExperienceModel{
		ZGPU:                types.BoolValue(ux.ZGPU),
		BrowserInBrowser:    types.BoolValue(ux.BrowserInBrowser),
		PersistIsolationBar: types.BoolValue(ux.PersistIsolationBar),
		Translate:           types.BoolValue(ux.Translate),
		SessionPersistence:  types.BoolValue(ux.SessionPersistence),
	}

	if ux.ForwardToZia != nil {
		model.ForwardToZia = []CBIForwardToZiaModel{
			{
				Enabled:        types.BoolValue(ux.ForwardToZia.Enabled),
				OrganizationID: helpers.StringValueOrNull(ux.ForwardToZia.OrganizationID),
				CloudName:      helpers.StringValueOrNull(ux.ForwardToZia.CloudName),
				PacFileURL:     helpers.StringValueOrNull(ux.ForwardToZia.PacFileUrl),
			},
		}
	}

	return []CBIUserExperienceModel{model}, diag.Diagnostics{}
}

func flattenCBISecurityControls(ctx context.Context, sc *cbiprofilecontroller.SecurityControls) ([]CBISecurityControlsModel, diag.Diagnostics) {
	if sc == nil {
		return nil, diag.Diagnostics{}
	}

	model := CBISecurityControlsModel{
		CopyPaste:          helpers.StringValueOrNull(sc.CopyPaste),
		UploadDownload:     helpers.StringValueOrNull(sc.UploadDownload),
		DocumentViewer:     types.BoolValue(sc.DocumentViewer),
		LocalRender:        types.BoolValue(sc.LocalRender),
		AllowPrinting:      types.BoolValue(sc.AllowPrinting),
		RestrictKeystrokes: types.BoolValue(sc.RestrictKeystrokes),
		FlattenedPDF:       types.BoolValue(sc.FlattenedPdf),
	}

	var diags diag.Diagnostics

	if sc.DeepLink != nil {
		applications, appDiags := types.SetValueFrom(ctx, types.StringType, sc.DeepLink.Applications)
		diags.Append(appDiags...)
		model.DeepLink = []CBIDeepLinkModel{
			{
				Enabled:      types.BoolValue(sc.DeepLink.Enabled),
				Applications: applications,
			},
		}
	}

	if sc.Watermark != nil {
		model.Watermark = []CBIWatermarkModel{
			{
				Enabled:       types.BoolValue(sc.Watermark.Enabled),
				ShowUserID:    types.BoolValue(sc.Watermark.ShowUserID),
				ShowTimestamp: types.BoolValue(sc.Watermark.ShowTimestamp),
				ShowMessage:   types.BoolValue(sc.Watermark.ShowMessage),
				Message:       helpers.StringValueOrNull(sc.Watermark.Message),
			},
		}
	}

	return []CBISecurityControlsModel{model}, diags
}

func flattenCBIDebugMode(debug *cbiprofilecontroller.DebugMode) []CBIDebugModeModel {
	if debug == nil {
		return nil
	}

	if !debug.Allowed && strings.TrimSpace(debug.FilePassword) == "" {
		return nil
	}

	return []CBIDebugModeModel{
		{
			Allowed:      types.BoolValue(debug.Allowed),
			FilePassword: helpers.StringValueOrNull(debug.FilePassword),
		},
	}
}
