package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
)

var (
	_ datasource.DataSource              = &IdpControllerDataSource{}
	_ datasource.DataSourceWithConfigure = &IdpControllerDataSource{}
)

func NewIdpControllerDataSource() datasource.DataSource {
	return &IdpControllerDataSource{}
}

type IdpControllerDataSource struct {
	client *client.Client
}

type IdpControllerDataSourceModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	Description                 types.String `tfsdk:"description"`
	Enabled                     types.Bool   `tfsdk:"enabled"`
	AdminMetadata               types.Set    `tfsdk:"admin_metadata"`
	AdminSpSigningCertID        types.String `tfsdk:"admin_sp_signing_cert_id"`
	AutoProvision               types.String `tfsdk:"auto_provision"`
	CreationTime                types.String `tfsdk:"creation_time"`
	DisableSamlBasedPolicy      types.Bool   `tfsdk:"disable_saml_based_policy"`
	DomainList                  types.List   `tfsdk:"domain_list"`
	EnableScimBasedPolicy       types.Bool   `tfsdk:"enable_scim_based_policy"`
	IdpEntityID                 types.String `tfsdk:"idp_entity_id"`
	LoginNameAttribute          types.String `tfsdk:"login_name_attribute"`
	LoginURL                    types.String `tfsdk:"login_url"`
	LoginHint                   types.Bool   `tfsdk:"login_hint"`
	ForceAuth                   types.Bool   `tfsdk:"force_auth"`
	EnableArbitraryAuthDomains  types.String `tfsdk:"enable_arbitrary_auth_domains"`
	ModifiedBy                  types.String `tfsdk:"modifiedby"`
	ModifiedTime                types.String `tfsdk:"modified_time"`
	ReauthOnUserUpdate          types.Bool   `tfsdk:"reauth_on_user_update"`
	RedirectBinding             types.Bool   `tfsdk:"redirect_binding"`
	ScimEnabled                 types.Bool   `tfsdk:"scim_enabled"`
	ScimServiceProviderEndpoint types.String `tfsdk:"scim_service_provider_endpoint"`
	ScimSharedSecretExists      types.Bool   `tfsdk:"scim_shared_secret_exists"`
	SignSamlRequest             types.String `tfsdk:"sign_saml_request"`
	SsoType                     types.List   `tfsdk:"sso_type"`
	UseCustomSpMetadata         types.Bool   `tfsdk:"use_custom_sp_metadata"`
	UserMetadata                types.Set    `tfsdk:"user_metadata"`
	UserSpSigningCertID         types.String `tfsdk:"user_sp_signing_cert_id"`
}

func (d *IdpControllerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_idp_controller"
}

func (d *IdpControllerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	metadataAttributes := map[string]schema.Attribute{
		"certificate_url": schema.StringAttribute{Computed: true},
		"sp_base_url":     schema.StringAttribute{Computed: true},
		"sp_entity_id":    schema.StringAttribute{Computed: true},
		"sp_metadata_url": schema.StringAttribute{Computed: true},
		"sp_post_url":     schema.StringAttribute{Computed: true},
	}

	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA IdP controller by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id":                        schema.StringAttribute{Optional: true},
			"name":                      schema.StringAttribute{Optional: true},
			"description":               schema.StringAttribute{Computed: true},
			"enabled":                   schema.BoolAttribute{Computed: true},
			"admin_sp_signing_cert_id":  schema.StringAttribute{Computed: true},
			"auto_provision":            schema.StringAttribute{Computed: true},
			"creation_time":             schema.StringAttribute{Computed: true},
			"disable_saml_based_policy": schema.BoolAttribute{Computed: true},
			"domain_list": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"enable_scim_based_policy":       schema.BoolAttribute{Computed: true},
			"idp_entity_id":                  schema.StringAttribute{Computed: true},
			"login_name_attribute":           schema.StringAttribute{Computed: true},
			"login_url":                      schema.StringAttribute{Computed: true},
			"login_hint":                     schema.BoolAttribute{Computed: true},
			"force_auth":                     schema.BoolAttribute{Computed: true},
			"enable_arbitrary_auth_domains":  schema.StringAttribute{Computed: true},
			"modifiedby":                     schema.StringAttribute{Computed: true},
			"modified_time":                  schema.StringAttribute{Computed: true},
			"reauth_on_user_update":          schema.BoolAttribute{Computed: true},
			"redirect_binding":               schema.BoolAttribute{Computed: true},
			"scim_enabled":                   schema.BoolAttribute{Computed: true},
			"scim_service_provider_endpoint": schema.StringAttribute{Computed: true},
			"scim_shared_secret_exists":      schema.BoolAttribute{Computed: true},
			"sign_saml_request":              schema.StringAttribute{Computed: true},
			"sso_type": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"use_custom_sp_metadata":  schema.BoolAttribute{Computed: true},
			"user_sp_signing_cert_id": schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			"admin_metadata": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: metadataAttributes,
				},
			},
			"user_metadata": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: metadataAttributes,
				},
			},
		},
	}
}

func (d *IdpControllerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *IdpControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IdpControllerDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (data.ID.IsNull() || data.ID.ValueString() == "") && (data.Name.IsNull() || data.Name.ValueString() == "") {
		resp.Diagnostics.AddError("Missing Identifier", "Either id or name must be specified")
		return
	}

	service := d.client.Service

	var (
		controller *idpcontroller.IdpController
		err        error
	)

	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		id := data.ID.ValueString()
		tflog.Info(ctx, "Fetching IdP controller", map[string]any{"id": id})
		controller, _, err = idpcontroller.Get(ctx, service, id)
	} else {
		name := data.Name.ValueString()
		tflog.Info(ctx, "Fetching IdP controller", map[string]any{"name": name})
		controller, _, err = idpcontroller.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read IdP controller: %v", err))
		return
	}

	adminMeta, adminDiags := flattenAdminMetadata(ctx, controller.AdminMetadata)
	resp.Diagnostics.Append(adminDiags...)
	userMeta, userDiags := flattenUserMetadata(ctx, controller.UserMetadata)
	resp.Diagnostics.Append(userDiags...)

	domainList, domainDiags := types.ListValueFrom(ctx, types.StringType, controller.Domainlist)
	resp.Diagnostics.Append(domainDiags...)

	ssoList, ssoDiags := types.ListValueFrom(ctx, types.StringType, controller.SsoType)
	resp.Diagnostics.Append(ssoDiags...)

	model := IdpControllerDataSourceModel{
		ID:                          types.StringValue(controller.ID),
		Name:                        types.StringValue(controller.Name),
		Description:                 types.StringValue(controller.Description),
		Enabled:                     types.BoolValue(controller.Enabled),
		AdminMetadata:               adminMeta,
		AdminSpSigningCertID:        types.StringValue(controller.AdminSpSigningCertID),
		AutoProvision:               types.StringValue(controller.AutoProvision),
		CreationTime:                types.StringValue(controller.CreationTime),
		DisableSamlBasedPolicy:      types.BoolValue(controller.DisableSamlBasedPolicy),
		DomainList:                  domainList,
		EnableScimBasedPolicy:       types.BoolValue(controller.EnableScimBasedPolicy),
		IdpEntityID:                 types.StringValue(controller.IdpEntityID),
		LoginNameAttribute:          types.StringValue(controller.LoginNameAttribute),
		LoginURL:                    types.StringValue(controller.LoginURL),
		LoginHint:                   types.BoolValue(controller.LoginHint),
		ForceAuth:                   types.BoolValue(controller.ForceAuth),
		EnableArbitraryAuthDomains:  types.StringValue(controller.EnableArbitraryAuthDomains),
		ModifiedBy:                  types.StringValue(controller.ModifiedBy),
		ModifiedTime:                types.StringValue(controller.ModifiedTime),
		ReauthOnUserUpdate:          types.BoolValue(controller.ReauthOnUserUpdate),
		RedirectBinding:             types.BoolValue(controller.RedirectBinding),
		ScimEnabled:                 types.BoolValue(controller.ScimEnabled),
		ScimServiceProviderEndpoint: types.StringValue(controller.ScimServiceProviderEndpoint),
		ScimSharedSecretExists:      types.BoolValue(controller.ScimSharedSecretExists),
		SignSamlRequest:             types.StringValue(controller.SignSamlRequest),
		SsoType:                     ssoList,
		UseCustomSpMetadata:         types.BoolValue(controller.UseCustomSpMetadata),
		UserMetadata:                userMeta,
		UserSpSigningCertID:         types.StringValue(controller.UserSpSigningCertID),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func flattenAdminMetadata(ctx context.Context, meta *idpcontroller.AdminMetadata) (types.Set, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"certificate_url": types.StringType,
		"sp_base_url":     types.StringType,
		"sp_entity_id":    types.StringType,
		"sp_metadata_url": types.StringType,
		"sp_post_url":     types.StringType,
	}

	if meta == nil {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), nil
	}

	obj, diags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"certificate_url": types.StringValue(meta.CertificateURL),
		"sp_base_url":     types.StringValue(meta.SpBaseURL),
		"sp_entity_id":    types.StringValue(meta.SpEntityID),
		"sp_metadata_url": types.StringValue(meta.SpMetadataURL),
		"sp_post_url":     types.StringValue(meta.SpPostURL),
	})
	if diags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	set, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	if setDiags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), setDiags
	}

	return set, setDiags
}

func flattenUserMetadata(ctx context.Context, meta *idpcontroller.UserMetadata) (types.Set, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"certificate_url": types.StringType,
		"sp_base_url":     types.StringType,
		"sp_entity_id":    types.StringType,
		"sp_metadata_url": types.StringType,
		"sp_post_url":     types.StringType,
	}

	if meta == nil {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), nil
	}

	obj, diags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"certificate_url": types.StringValue(meta.CertificateURL),
		"sp_base_url":     types.StringValue(meta.SpBaseURL),
		"sp_entity_id":    types.StringValue(meta.SpEntityID),
		"sp_metadata_url": types.StringValue(meta.SpMetadataURL),
		"sp_post_url":     types.StringValue(meta.SpPostURL),
	})
	if diags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	set, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	if setDiags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), setDiags
	}

	return set, setDiags
}
