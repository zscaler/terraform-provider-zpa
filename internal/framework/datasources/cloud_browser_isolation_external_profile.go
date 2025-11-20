package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	cbiprofilecontroller "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbiprofilecontroller"
)

var (
	_ datasource.DataSource              = &CBIExternalProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &CBIExternalProfileDataSource{}
)

func NewCBIExternalProfileDataSource() datasource.DataSource {
	return &CBIExternalProfileDataSource{}
}

type CBIExternalProfileDataSource struct {
	client *client.Client
}

type CBIExternalProfileModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	BannerID       types.String `tfsdk:"banner_id"`
	IsDefault      types.Bool   `tfsdk:"is_default"`
	Href           types.String `tfsdk:"href"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	CreationTime   types.String `tfsdk:"creation_time"`
	ModifiedBy     types.String `tfsdk:"modified_by"`
	ModifiedTime   types.String `tfsdk:"modified_time"`
	CBITenantID    types.String `tfsdk:"cbi_tenant_id"`
	CBIProfileID   types.String `tfsdk:"cbi_profile_id"`
	CBIURL         types.String `tfsdk:"cbi_url"`
	Regions        types.List   `tfsdk:"regions"`
	CertificateIDs types.List   `tfsdk:"certificate_ids"`
	UserExperience types.Set    `tfsdk:"user_experience"`
	SecurityCtrl   types.Set    `tfsdk:"security_controls"`
	DebugMode      types.Set    `tfsdk:"debug_mode"`
}

func (d *CBIExternalProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_browser_isolation_external_profile"
}

func (d *CBIExternalProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a Cloud Browser Isolation external profile by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the CBI external profile.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the CBI external profile.",
			},
			"description":    schema.StringAttribute{Computed: true},
			"banner_id":      schema.StringAttribute{Computed: true},
			"is_default":     schema.BoolAttribute{Computed: true},
			"href":           schema.StringAttribute{Computed: true},
			"enabled":        schema.BoolAttribute{Computed: true},
			"creation_time":  schema.StringAttribute{Computed: true},
			"modified_by":    schema.StringAttribute{Computed: true},
			"modified_time":  schema.StringAttribute{Computed: true},
			"cbi_tenant_id":  schema.StringAttribute{Computed: true},
			"cbi_profile_id": schema.StringAttribute{Computed: true},
			"cbi_url":        schema.StringAttribute{Computed: true},
			"certificate_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"regions": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
					},
				},
			},
			"user_experience": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"zgpu":                  schema.BoolAttribute{Computed: true},
						"browser_in_browser":    schema.BoolAttribute{Computed: true},
						"persist_isolation_bar": schema.BoolAttribute{Computed: true},
						"translate":             schema.BoolAttribute{Computed: true},
						"session_persistence":   schema.BoolAttribute{Computed: true},
					},
					Blocks: map[string]schema.Block{
						"forward_to_zia": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"enabled":         schema.BoolAttribute{Computed: true},
									"organization_id": schema.StringAttribute{Computed: true},
									"cloud_name":      schema.StringAttribute{Computed: true},
									"pac_file_url":    schema.StringAttribute{Computed: true},
								},
							},
						},
					},
				},
			},
			"security_controls": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"copy_paste":          schema.StringAttribute{Computed: true},
						"upload_download":     schema.StringAttribute{Computed: true},
						"document_viewer":     schema.BoolAttribute{Computed: true},
						"local_render":        schema.BoolAttribute{Computed: true},
						"allow_printing":      schema.BoolAttribute{Computed: true},
						"restrict_keystrokes": schema.BoolAttribute{Computed: true},
						"flattened_pdf":       schema.BoolAttribute{Computed: true},
					},
					Blocks: map[string]schema.Block{
						"deep_link": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{Computed: true},
									"applications": schema.ListAttribute{
										ElementType: types.StringType,
										Computed:    true,
									},
								},
							},
						},
						"watermark": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"enabled":        schema.BoolAttribute{Computed: true},
									"show_user_id":   schema.BoolAttribute{Computed: true},
									"show_timestamp": schema.BoolAttribute{Computed: true},
									"show_message":   schema.BoolAttribute{Computed: true},
									"message":        schema.StringAttribute{Computed: true},
								},
							},
						},
					},
				},
			},
			"debug_mode": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"allowed":       schema.BoolAttribute{Computed: true},
						"file_password": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *CBIExternalProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CBIExternalProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data CBIExternalProfileModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())
	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided to read a CBI external profile.")
		return
	}

	identifier := id
	if identifier == "" {
		identifier = name
	}

	profile, _, err := cbiprofilecontroller.GetByNameOrID(ctx, d.client.Service, identifier)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read CBI external profile: %v", err))
		return
	}

	if profile == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("CBI external profile with identifier %q was not found.", identifier))
		return
	}

	regions, regionsDiags := flattenCBIRegions(ctx, profile.Regions)
	resp.Diagnostics.Append(regionsDiags...)

	certificateIDs, certDiags := types.ListValueFrom(ctx, types.StringType, profile.CertificateIDs)
	resp.Diagnostics.Append(certDiags...)

	userExperience, uxDiags := flattenCBIUserExperience(ctx, profile.UserExperience)
	resp.Diagnostics.Append(uxDiags...)

	securityControls, scDiags := flattenCBISecurityControls(ctx, profile.SecurityControls)
	resp.Diagnostics.Append(scDiags...)

	debugMode, debugDiags := flattenCBIDebugMode(ctx, profile.DebugMode)
	resp.Diagnostics.Append(debugDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(profile.ID)
	data.Name = stringOrNull(profile.Name)
	data.Description = stringOrNull(profile.Description)
	data.BannerID = stringOrNull(profile.BannerID)
	data.IsDefault = types.BoolValue(profile.IsDefault)
	data.Href = stringOrNull(profile.Href)
	data.Enabled = types.BoolValue(profile.Enabled)
	data.CreationTime = stringOrNull(profile.CreationTime)
	data.ModifiedBy = stringOrNull(profile.ModifiedBy)
	data.ModifiedTime = stringOrNull(profile.ModifiedTime)
	data.CBITenantID = stringOrNull(profile.CBITenantID)
	data.CBIProfileID = stringOrNull(profile.CBIProfileID)
	data.CBIURL = stringOrNull(profile.CBIURL)
	data.Regions = regions
	data.CertificateIDs = certificateIDs
	data.UserExperience = userExperience
	data.SecurityCtrl = securityControls
	data.DebugMode = debugMode

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenCBIRegions(ctx context.Context, regions []cbiprofilecontroller.Regions) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}

	if len(regions) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(regions))
	var diags diag.Diagnostics
	for _, region := range regions {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":   stringOrNull(region.ID),
			"name": stringOrNull(region.Name),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}

func flattenCBIUserExperience(ctx context.Context, ux *cbiprofilecontroller.UserExperience) (types.Set, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"zgpu":                  types.BoolType,
		"browser_in_browser":    types.BoolType,
		"persist_isolation_bar": types.BoolType,
		"translate":             types.BoolType,
		"session_persistence":   types.BoolType,
		"forward_to_zia":        types.ListType{ElemType: types.ObjectType{AttrTypes: cbiForwardToZiaAttrTypes()}},
	}

	if ux == nil {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	forwardToZia, diags := flattenCBIForwardToZia(ctx, ux.ForwardToZia)

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"zgpu":                  types.BoolValue(ux.ZGPU),
		"browser_in_browser":    types.BoolValue(ux.BrowserInBrowser),
		"persist_isolation_bar": types.BoolValue(ux.PersistIsolationBar),
		"translate":             types.BoolValue(ux.Translate),
		"session_persistence":   types.BoolValue(ux.SessionPersistence),
		"forward_to_zia":        forwardToZia,
	})
	diags.Append(objDiags...)

	set, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	diags.Append(setDiags...)
	return set, diags
}

func cbiForwardToZiaAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":         types.BoolType,
		"organization_id": types.StringType,
		"cloud_name":      types.StringType,
		"pac_file_url":    types.StringType,
	}
}

func flattenCBIForwardToZia(ctx context.Context, fz *cbiprofilecontroller.ForwardToZia) (types.List, diag.Diagnostics) {
	attrTypes := cbiForwardToZiaAttrTypes()
	if fz == nil {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"enabled":         types.BoolValue(fz.Enabled),
		"organization_id": stringOrNull(fz.OrganizationID),
		"cloud_name":      stringOrNull(fz.CloudName),
		"pac_file_url":    stringOrNull(fz.PacFileUrl),
	})

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	var diags diag.Diagnostics
	diags.Append(objDiags...)
	diags.Append(listDiags...)
	return list, diags
}

func flattenCBISecurityControls(ctx context.Context, sc *cbiprofilecontroller.SecurityControls) (types.Set, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"copy_paste":          types.StringType,
		"upload_download":     types.StringType,
		"document_viewer":     types.BoolType,
		"local_render":        types.BoolType,
		"allow_printing":      types.BoolType,
		"restrict_keystrokes": types.BoolType,
		"flattened_pdf":       types.BoolType,
		"deep_link":           types.SetType{ElemType: types.ObjectType{AttrTypes: cbiDeepLinkAttrTypes()}},
		"watermark":           types.SetType{ElemType: types.ObjectType{AttrTypes: cbiWatermarkAttrTypes()}},
	}

	if sc == nil {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	deepLink, diags := flattenCBIDeepLink(ctx, sc.DeepLink)
	watermark, wmDiags := flattenCBIWatermark(ctx, sc.Watermark)
	diags.Append(wmDiags...)

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"copy_paste":          stringOrNull(sc.CopyPaste),
		"upload_download":     stringOrNull(sc.UploadDownload),
		"document_viewer":     types.BoolValue(sc.DocumentViewer),
		"local_render":        types.BoolValue(sc.LocalRender),
		"allow_printing":      types.BoolValue(sc.AllowPrinting),
		"restrict_keystrokes": types.BoolValue(sc.RestrictKeystrokes),
		"flattened_pdf":       types.BoolValue(sc.FlattenedPdf),
		"deep_link":           deepLink,
		"watermark":           watermark,
	})
	diags.Append(objDiags...)

	set, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	diags.Append(setDiags...)
	return set, diags
}

func cbiDeepLinkAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":      types.BoolType,
		"applications": types.ListType{ElemType: types.StringType},
	}
}

func flattenCBIDeepLink(ctx context.Context, deepLink *cbiprofilecontroller.DeepLink) (types.Set, diag.Diagnostics) {
	attrTypes := cbiDeepLinkAttrTypes()
	if deepLink == nil {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	applications := types.ListNull(types.StringType)
	var diags diag.Diagnostics
	if len(deepLink.Applications) > 0 {
		list, listDiags := types.ListValueFrom(ctx, types.StringType, deepLink.Applications)
		diags.Append(listDiags...)
		applications = list
	}

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"enabled":      types.BoolValue(deepLink.Enabled),
		"applications": applications,
	})
	diags.Append(objDiags...)

	set, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	diags.Append(setDiags...)
	return set, diags
}

func cbiWatermarkAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":        types.BoolType,
		"show_user_id":   types.BoolType,
		"show_timestamp": types.BoolType,
		"show_message":   types.BoolType,
		"message":        types.StringType,
	}
}

func flattenCBIWatermark(ctx context.Context, watermark *cbiprofilecontroller.Watermark) (types.Set, diag.Diagnostics) {
	attrTypes := cbiWatermarkAttrTypes()
	if watermark == nil {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"enabled":        types.BoolValue(watermark.Enabled),
		"show_user_id":   types.BoolValue(watermark.ShowUserID),
		"show_timestamp": types.BoolValue(watermark.ShowTimestamp),
		"show_message":   types.BoolValue(watermark.ShowMessage),
		"message":        stringOrNull(watermark.Message),
	})

	set, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	var diags diag.Diagnostics
	diags.Append(objDiags...)
	diags.Append(setDiags...)
	return set, diags
}

func flattenCBIDebugMode(ctx context.Context, debug *cbiprofilecontroller.DebugMode) (types.Set, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"allowed":       types.BoolType,
		"file_password": types.StringType,
	}

	if debug == nil {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"allowed":       types.BoolValue(debug.Allowed),
		"file_password": stringOrNull(debug.FilePassword),
	})

	set, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	var diags diag.Diagnostics
	diags.Append(objDiags...)
	diags.Append(setDiags...)
	return set, diags
}
