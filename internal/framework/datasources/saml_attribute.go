package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/samlattribute"
)

var (
	_ datasource.DataSource              = &SAMLAttributeDataSource{}
	_ datasource.DataSourceWithConfigure = &SAMLAttributeDataSource{}
)

func NewSAMLAttributeDataSource() datasource.DataSource {
	return &SAMLAttributeDataSource{}
}

type SAMLAttributeDataSource struct {
	client *client.Client
}

type SAMLAttributeModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	IdpID         types.String `tfsdk:"idp_id"`
	IdpName       types.String `tfsdk:"idp_name"`
	CreationTime  types.String `tfsdk:"creation_time"`
	ModifiedBy    types.String `tfsdk:"modifiedby"`
	ModifiedTime  types.String `tfsdk:"modified_time"`
	SAMLName      types.String `tfsdk:"saml_name"`
	UserAttribute types.Bool   `tfsdk:"user_attribute"`
}

func (d *SAMLAttributeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_saml_attribute"
}

func (d *SAMLAttributeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a SAML attribute by ID or name within a specific IDP.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the SAML attribute.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the SAML attribute.",
			},
			"idp_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the associated IDP.",
			},
			"idp_name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the associated IDP.",
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"modifiedby": schema.StringAttribute{
				Computed: true,
			},
			"modified_time": schema.StringAttribute{
				Computed: true,
			},
			"saml_name": schema.StringAttribute{
				Computed: true,
			},
			"user_attribute": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (d *SAMLAttributeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SAMLAttributeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data SAMLAttributeModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())
	idpID := strings.TrimSpace(data.IdpID.ValueString())
	idpName := strings.TrimSpace(data.IdpName.ValueString())

	if idpID == "" && idpName == "" {
		resp.Diagnostics.AddError("Missing Required Attribute", "Either 'idp_id' or 'idp_name' must be provided.")
		return
	}

	service := d.client.Service

	tflog.Debug(ctx, "Resolving IDP for SAML attribute lookup", map[string]any{
		"idp_id":   idpID,
		"idp_name": idpName,
	})

	var (
		idp *idpcontroller.IdpController
		err error
	)

	if idpID != "" {
		idp, _, err = idpcontroller.Get(ctx, service, idpID)
	} else {
		idp, _, err = idpcontroller.GetByName(ctx, service, idpName)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read IDP: %v", err))
		return
	}

	if idp == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Unable to locate IDP with id %q or name %q.", idpID, idpName))
		return
	}

	tflog.Debug(ctx, "Retrieving SAML attribute", map[string]any{
		"id":       id,
		"name":     name,
		"idp_id":   idp.ID,
		"idp_name": idp.Name,
	})

	var attr *samlattribute.SamlAttribute
	if id != "" {
		attr, _, err = samlattribute.GetByIdpAndAttributeID(ctx, service, idp.ID, id)
	} else if name != "" {
		attrs, _, e := samlattribute.GetAllByIdp(ctx, service, idp.ID)
		if e != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list SAML attributes: %v", e))
			return
		}

		for _, candidate := range attrs {
			if strings.EqualFold(candidate.Name, name) {
				attr = &candidate
				break
			}
		}

		if attr == nil {
			resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No SAML attribute named %q found in IDP %q.", name, idp.Name))
			return
		}
	} else {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided.")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SAML attribute: %v", err))
		return
	}

	if attr == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("SAML attribute not found with id %q or name %q.", id, name))
		return
	}

	state, diags := flattenSAMLAttribute(attr)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.IdpID = types.StringValue(idp.ID)
	state.IdpName = types.StringValue(idp.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenSAMLAttribute(attr *samlattribute.SamlAttribute) (SAMLAttributeModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := SAMLAttributeModel{
		ID:            types.StringValue(attr.ID),
		Name:          types.StringValue(attr.Name),
		CreationTime:  types.StringValue(attr.CreationTime),
		ModifiedBy:    types.StringValue(attr.ModifiedBy),
		ModifiedTime:  types.StringValue(attr.ModifiedTime),
		SAMLName:      types.StringValue(attr.SamlName),
		UserAttribute: types.BoolValue(attr.UserAttribute),
	}

	return model, diags
}
