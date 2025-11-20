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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimattributeheader"
)

var (
	_ datasource.DataSource              = &SCIMAttributeHeaderDataSource{}
	_ datasource.DataSourceWithConfigure = &SCIMAttributeHeaderDataSource{}
)

func NewSCIMAttributeHeaderDataSource() datasource.DataSource {
	return &SCIMAttributeHeaderDataSource{}
}

type SCIMAttributeHeaderDataSource struct {
	client *client.Client
}

type SCIMAttributeHeaderModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	IdpID           types.String `tfsdk:"idp_id"`
	IdpName         types.String `tfsdk:"idp_name"`
	CanonicalValues types.List   `tfsdk:"canonical_values"`
	CaseSensitive   types.Bool   `tfsdk:"case_sensitive"`
	CreationTime    types.String `tfsdk:"creation_time"`
	DataType        types.String `tfsdk:"data_type"`
	Description     types.String `tfsdk:"description"`
	ModifiedBy      types.String `tfsdk:"modifiedby"`
	ModifiedTime    types.String `tfsdk:"modified_time"`
	Multivalued     types.Bool   `tfsdk:"multivalued"`
	Mutability      types.String `tfsdk:"mutability"`
	Required        types.Bool   `tfsdk:"required"`
	Returned        types.String `tfsdk:"returned"`
	SchemaURI       types.String `tfsdk:"schema_uri"`
	Uniqueness      types.Bool   `tfsdk:"uniqueness"`
	Values          types.Set    `tfsdk:"values"`
}

func (d *SCIMAttributeHeaderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scim_attribute_header"
}

func (d *SCIMAttributeHeaderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a SCIM attribute header definition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the SCIM attribute header.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the SCIM attribute header.",
			},
			"idp_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the IDP to scope the lookup.",
			},
			"idp_name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the IDP to scope the lookup.",
			},
			"canonical_values": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"case_sensitive": schema.BoolAttribute{
				Computed: true,
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"data_type": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"modifiedby": schema.StringAttribute{
				Computed: true,
			},
			"modified_time": schema.StringAttribute{
				Computed: true,
			},
			"multivalued": schema.BoolAttribute{
				Computed: true,
			},
			"mutability": schema.StringAttribute{
				Computed: true,
			},
			"required": schema.BoolAttribute{
				Computed: true,
			},
			"returned": schema.StringAttribute{
				Computed: true,
			},
			"schema_uri": schema.StringAttribute{
				Computed: true,
			},
			"uniqueness": schema.BoolAttribute{
				Computed: true,
			},
			"values": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (d *SCIMAttributeHeaderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SCIMAttributeHeaderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data SCIMAttributeHeaderModel
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

	tflog.Debug(ctx, "Resolving IDP for SCIM attribute header", map[string]any{
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

	var header *scimattributeheader.ScimAttributeHeader

	if id != "" {
		header, _, err = scimattributeheader.Get(ctx, service, idp.ID, id)
	} else if name != "" {
		header, _, err = scimattributeheader.GetByName(ctx, service, name, idp.ID)
	} else {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided.")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SCIM attribute header: %v", err))
		return
	}

	if header == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("SCIM attribute header not found with id %q or name %q.", id, name))
		return
	}

	values, err := scimattributeheader.GetValues(ctx, service, header.IdpID, header.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SCIM attribute header values: %v", err))
		return
	}

	state, diags := flattenSCIMAttributeHeader(ctx, header, values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.IdpID = types.StringValue(idp.ID)
	state.IdpName = types.StringValue(idp.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenSCIMAttributeHeader(ctx context.Context, header *scimattributeheader.ScimAttributeHeader, values []string) (SCIMAttributeHeaderModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	canonicalList, canonicalDiags := types.ListValueFrom(ctx, types.StringType, header.CanonicalValues)
	diags.Append(canonicalDiags...)

	var valuesSet types.Set
	var setDiags diag.Diagnostics
	if len(values) > 0 {
		elements := make([]attr.Value, 0, len(values))
		for _, v := range values {
			elements = append(elements, types.StringValue(v))
		}
		valuesSet, setDiags = types.SetValue(types.StringType, elements)
	} else {
		valuesSet = types.SetNull(types.StringType)
	}
	diags.Append(setDiags...)

	model := SCIMAttributeHeaderModel{
		ID:              types.StringValue(header.ID),
		Name:            types.StringValue(header.Name),
		CanonicalValues: canonicalList,
		CaseSensitive:   types.BoolValue(header.CaseSensitive),
		CreationTime:    types.StringValue(header.CreationTime),
		DataType:        types.StringValue(header.DataType),
		Description:     types.StringValue(header.Description),
		ModifiedBy:      types.StringValue(header.ModifiedBy),
		ModifiedTime:    types.StringValue(header.ModifiedTime),
		Multivalued:     types.BoolValue(header.MultiValued),
		Mutability:      types.StringValue(header.Mutability),
		Required:        types.BoolValue(header.Required),
		Returned:        types.StringValue(header.Returned),
		SchemaURI:       types.StringValue(header.SchemaURI),
		Uniqueness:      types.BoolValue(header.Uniqueness),
		Values:          valuesSet,
	}

	return model, diags
}
