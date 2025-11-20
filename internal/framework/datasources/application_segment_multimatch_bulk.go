package datasources

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
)

var (
	_ datasource.DataSource              = &ApplicationSegmentMultimatchBulkDataSource{}
	_ datasource.DataSourceWithConfigure = &ApplicationSegmentMultimatchBulkDataSource{}
)

func NewApplicationSegmentMultimatchBulkDataSource() datasource.DataSource {
	return &ApplicationSegmentMultimatchBulkDataSource{}
}

type ApplicationSegmentMultimatchBulkDataSource struct {
	client *client.Client
}

type ApplicationSegmentMultimatchBulkModel struct {
	ID              types.String `tfsdk:"id"`
	DomainNames     types.List   `tfsdk:"domain_names"`
	UnsupportedRefs types.List   `tfsdk:"unsupported_references"`
	MicroTenantID   types.String `tfsdk:"microtenant_id"`
}

func (d *ApplicationSegmentMultimatchBulkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_segment_multimatch_bulk"
}

func (d *ApplicationSegmentMultimatchBulkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves application segments that do not support multimatch for the provided domain list.",
		Attributes: map[string]schema.Attribute{
			"domain_names": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "List of domain names to check for unsupported multimatch references.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Synthetic identifier derived from the requested domain names.",
			},
			"microtenant_id": schema.StringAttribute{
				Computed:    true,
				Description: "Micro-tenant ID used to scope the lookup.",
			},
		},
		Blocks: map[string]schema.Block{
			"unsupported_references": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":               schema.StringAttribute{Computed: true},
						"app_segment_name": schema.StringAttribute{Computed: true},
						"domains": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"tcp_ports": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"match_style":      schema.StringAttribute{Computed: true},
						"microtenant_name": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *ApplicationSegmentMultimatchBulkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ApplicationSegmentMultimatchBulkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data ApplicationSegmentMultimatchBulkModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainNames := make([]string, 0)
	if !data.DomainNames.IsNull() && !data.DomainNames.IsUnknown() {
		var tmp []string
		diag := data.DomainNames.ElementsAs(ctx, &tmp, false)
		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, domain := range tmp {
			trimmed := strings.TrimSpace(domain)
			if trimmed != "" {
				domainNames = append(domainNames, trimmed)
			}
		}
	}

	if len(domainNames) == 0 {
		resp.Diagnostics.AddError("Validation Error", "At least one domain name must be provided.")
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

	tflog.Debug(ctx, "Retrieving multimatch unsupported references", map[string]any{"domains": domainNames})

	payload := applicationsegment.MultiMatchUnsupportedReferencesPayload(domainNames)
	refs, _, err := applicationsegment.GetMultiMatchUnsupportedReferences(ctx, service, payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve multimatch unsupported references: %v", err))
		return
	}

	unsupported, unsupportedDiags := flattenUnsupportedReferences(ctx, refs)
	resp.Diagnostics.Append(unsupportedDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sort.Strings(domainNames)
	syntheticID := helpers.GenerateShortID(strings.Join(domainNames, ","))

	domainsList, domainsDiags := helpers.StringSliceToList(ctx, domainNames)
	resp.Diagnostics.Append(domainsDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.DomainNames = domainsList
	data.UnsupportedRefs = unsupported
	data.ID = types.StringValue(syntheticID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenUnsupportedReferences(ctx context.Context, refs []applicationsegment.MultiMatchUnsupportedReferencesResponse) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id":               types.StringType,
		"app_segment_name": types.StringType,
		"domains":          types.ListType{ElemType: types.StringType},
		"tcp_ports":        types.ListType{ElemType: types.StringType},
		"match_style":      types.StringType,
		"microtenant_name": types.StringType,
	}

	if len(refs) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(refs))
	var diags diag.Diagnostics
	for _, ref := range refs {
		domains, domainDiags := types.ListValueFrom(ctx, types.StringType, ref.Domains)
		diags.Append(domainDiags...)
		tcpPorts, portDiags := types.ListValueFrom(ctx, types.StringType, ref.TCPPorts)
		diags.Append(portDiags...)

		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":               helpers.StringValueOrNull(ref.ID),
			"app_segment_name": helpers.StringValueOrNull(ref.AppSegmentName),
			"domains":          domains,
			"tcp_ports":        tcpPorts,
			"match_style":      helpers.StringValueOrNull(ref.MatchStyle),
			"microtenant_name": helpers.StringValueOrNull(ref.MicrotenantName),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}
