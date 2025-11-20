// Copyright (c) SecurityGeekIO, Inc.
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/datasources"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/ephemeralresources"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/resources"
)

// Ensure ZPAProvider satisfies various provider interfaces.
var (
	_ provider.Provider                       = &ZPAProvider{}
	_ provider.ProviderWithEphemeralResources = &ZPAProvider{}
)

// ZPAProvider defines the provider implementation.
type ZPAProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ZPAProviderModel describes the provider data model.
type ZPAProviderModel struct {
	ClientID        types.String `tfsdk:"client_id"`
	ClientSecret    types.String `tfsdk:"client_secret"`
	PrivateKey      types.String `tfsdk:"private_key"`
	VanityDomain    types.String `tfsdk:"vanity_domain"`
	ZscalerCloud    types.String `tfsdk:"zscaler_cloud"`
	CustomerID      types.String `tfsdk:"customer_id"`
	MicrotenantID   types.String `tfsdk:"microtenant_id"`
	ZPAClientID     types.String `tfsdk:"zpa_client_id"`
	ZPAClientSecret types.String `tfsdk:"zpa_client_secret"`
	ZPACustomerID   types.String `tfsdk:"zpa_customer_id"`
	ZPACloud        types.String `tfsdk:"zpa_cloud"`
	UseLegacyClient types.Bool   `tfsdk:"use_legacy_client"`
	HTTPProxy       types.String `tfsdk:"http_proxy"`
	MaxRetries      types.Int64  `tfsdk:"max_retries"`
	Parallelism     types.Int64  `tfsdk:"parallelism"`
	RequestTimeout  types.Int64  `tfsdk:"request_timeout"`
	MinWaitSeconds  types.Int64  `tfsdk:"min_wait_seconds"`
	MaxWaitSeconds  types.Int64  `tfsdk:"max_wait_seconds"`
}

func (p *ZPAProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "zpa"
	resp.Version = p.version
}

func (p *ZPAProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "ZPA client ID",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "ZPA client secret",
				Optional:    true,
				Sensitive:   true,
			},
			"private_key": schema.StringAttribute{
				Description: "ZPA private key",
				Optional:    true,
				Sensitive:   true,
			},
			"vanity_domain": schema.StringAttribute{
				Description: "Zscaler Vanity Domain",
				Optional:    true,
				Sensitive:   true,
			},
			"zscaler_cloud": schema.StringAttribute{
				Description: "Zscaler Cloud Name",
				Optional:    true,
				Sensitive:   true,
			},
			"customer_id": schema.StringAttribute{
				Description: "ZPA customer ID",
				Optional:    true,
			},
			"microtenant_id": schema.StringAttribute{
				Description: "ZPA microtenant ID",
				Optional:    true,
			},
			"zpa_client_id": schema.StringAttribute{
				Description: "ZPA legacy API client ID",
				Optional:    true,
			},
			"zpa_client_secret": schema.StringAttribute{
				Description: "ZPA legacy API client secret",
				Optional:    true,
				Sensitive:   true,
			},
			"zpa_customer_id": schema.StringAttribute{
				Description: "ZPA legacy API customer ID",
				Optional:    true,
			},
			"zpa_cloud": schema.StringAttribute{
				Description: "ZPA cloud",
				Optional:    true,
			},
			"use_legacy_client": schema.BoolAttribute{
				Description: "Enable ZPA V2 (legacy) client",
				Optional:    true,
			},
			"http_proxy": schema.StringAttribute{
				Description: "Alternate HTTP proxy of scheme://hostname or scheme://hostname:port format",
				Optional:    true,
			},
			"max_retries": schema.Int64Attribute{
				Description: "Maximum number of retries to attempt before erroring out",
				Optional:    true,
			},
			"parallelism": schema.Int64Attribute{
				Description: "Number of concurrent requests to make within a resource where bulk operations are not possible",
				Optional:    true,
			},
			"request_timeout": schema.Int64Attribute{
				Description: "Timeout for single request (in seconds) which is made to Zscaler",
				Optional:    true,
			},
			"min_wait_seconds": schema.Int64Attribute{
				Description: "Minimum wait in seconds between retry attempts when rate limited",
				Optional:    true,
			},
			"max_wait_seconds": schema.Int64Attribute{
				Description: "Maximum wait in seconds between retry attempts when rate limited",
				Optional:    true,
			},
		},
	}
}

func (p *ZPAProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewAppConnectorAssistantScheduleResource,
		resources.NewAppServerControllerResource,
		resources.NewAppConnectorGroupResource,
		resources.NewPRACredentialControllerResource,
		resources.NewProvisioningKeyResource,
		resources.NewApplicationSegmentResource,
		resources.NewApplicationSegmentBrowserAccessResource,
		resources.NewBaCertificateResource,
		resources.NewApplicationSegmentInspectionResource,
		resources.NewApplicationSegmentMultimatchBulkResource,
		resources.NewApplicationSegmentPRAResource,
		resources.NewCBIBannerResource,
		resources.NewCBICertificateResource,
		resources.NewCBIExternalProfileResource,
		resources.NewC2CIPRangesResource,
		resources.NewLSSConfigControllerResource,
		resources.NewMicrotenantControllerResource,
		resources.NewPrivateCloudGroupResource,
		resources.NewPRAApprovalResource,
		resources.NewPRAConsoleControllerResource,
		resources.NewPRACredentialPoolResource,
		resources.NewUserPortalAUPResource,
		resources.NewUserPortalControllerResource,
		resources.NewUserPortalLinkResource,
		resources.NewInspectionCustomControlsResource,
		resources.NewInspectionProfileResource,
		resources.NewServiceEdgeGroupResource,
		resources.NewSegmentGroupsResource,
		resources.NewServerGroupResource,
		resources.NewPolicyAccessRuleResource,
		resources.NewPolicyAccessRuleReorderResource,
		resources.NewPolicyAccessForwardingRuleResource,
		resources.NewPolicyAccessInspectionRuleResource,
		resources.NewPolicyAccessIsolationRuleResource,
		resources.NewPolicyAccessRedirectionRuleResource,
		resources.NewPolicyAccessTimeoutRuleResource,
		resources.NewPolicyAccessRuleV2Resource,
		resources.NewPolicyAccessForwardingRuleV2Resource,
		resources.NewPolicyAccessBrowserProtectionRuleV2Resource,
		resources.NewPolicyAccessInspectionRuleV2Resource,
		resources.NewPolicyAccessIsolationRuleV2Resource,
		resources.NewPolicyAccessTimeoutRuleV2Resource,
		resources.NewPolicyAccessPortalRuleV2Resource,
		resources.NewPolicyAccessCapabilitiesRuleV2Resource,
		resources.NewPolicyAccessCredentialRuleV2Resource,
		resources.NewApplicationSegmentWeightedLBConfigResource,
		resources.NewEmergencyAccessResource,
		resources.NewPRAPortalControllerResource,
		resources.NewServiceEdgeAssistantScheduleResource,
		resources.NewZIACloudConfigResource,
		resources.NewBrowserProtectionResource,
	}
}

func (p *ZPAProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewAppServerControllerDataSource,
		datasources.NewAccessPolicyClientTypesDataSource,
		datasources.NewAccessPolicyPlatformDataSource,
		datasources.NewAppConnectorGroupDataSource,
		datasources.NewAppConnectorAssistantScheduleDataSource,
		datasources.NewIdpControllerDataSource,
		datasources.NewBaCertificateDataSource,
		datasources.NewBranchConnectorGroupDataSource,
		datasources.NewPostureProfileDataSource,
		datasources.NewEnrollmentCertDataSource,
		datasources.NewRiskScoreValuesDataSource,
		datasources.NewSAMLAttributeDataSource,
		datasources.NewSCIMAttributeHeaderDataSource,
		datasources.NewSCIMGroupDataSource,
		datasources.NewMachineGroupDataSource,
		datasources.NewIsolationProfileDataSource,
		datasources.NewApplicationSegmentByTypeDataSource,
		datasources.NewApplicationSegmentWeightedLBConfigDataSource,
		datasources.NewBrowserProtectionDataSource,
		datasources.NewCloudConfigDataSource,
		datasources.NewCustomerVersionProfileDataSource,
		datasources.NewExtranetResourcePartnerDataSource,
		datasources.NewLocationControllerSummaryDataSource,
		datasources.NewLocationControllerDataSource,
		datasources.NewLocationGroupControllerDataSource,
		datasources.NewUserPortalLinkDataSource,
		datasources.NewUserPortalControllerDataSource,
		datasources.NewUserPortalAUPDataSource,
		datasources.NewPRACredentialPoolDataSource,
		datasources.NewPRAPortalControllerDataSource,
		datasources.NewPrivateCloudControllerDataSource,
		datasources.NewPrivateCloudGroupDataSource,
		datasources.NewPRAConsoleControllerDataSource,
		datasources.NewPRAPrivilegedApprovalDataSource,
		datasources.NewPolicyTypeDataSource,
		datasources.NewMicrotenantControllerDataSource,
		datasources.NewManagedBrowserProfileDataSource,
		datasources.NewServiceEdgeControllerDataSource,
		datasources.NewProvisioningKeyDataSource,
		datasources.NewServiceEdgeAssistantScheduleDataSource,
		datasources.NewServiceEdgeGroupDataSource,
		datasources.NewSegmentGroupsDataSource,
		datasources.NewTrustedNetworkDataSource,
		datasources.NewWorkloadTagGroupDataSource,
		datasources.NewServerGroupDataSource,
		datasources.NewLSSStatusCodesDataSource,
		datasources.NewLSSLogTypeFormatsDataSource,
		datasources.NewLSSConfigControllerDataSource,
		datasources.NewLSSClientTypesDataSource,
		datasources.NewInspectionProfileDataSource,
		datasources.NewInspectionPredefinedControlsDataSource,
		datasources.NewInspectionCustomControlsDataSource,
		datasources.NewInspectionAllPredefinedControlsDataSource,
		datasources.NewCloudConnectorGroupDataSource,
		datasources.NewCBIZPAProfilesDataSource,
		datasources.NewCBIRegionsDataSource,
		datasources.NewCBIExternalProfileDataSource,
		datasources.NewCBICertificatesDataSource,
		datasources.NewCBIBannersDataSource,
		datasources.NewC2CIPRangesDataSource,
		datasources.NewApplicationSegmentDataSource,
		datasources.NewApplicationSegmentPRADataSource,
		datasources.NewApplicationSegmentMultimatchBulkDataSource,
		datasources.NewApplicationSegmentInspectionDataSource,
		datasources.NewApplicationSegmentBrowserAccessDataSource,
		datasources.NewAppConnectorControllerDataSource,
		datasources.NewPRACredentialControllerDataSource,
	}
}

func (p *ZPAProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		ephemeralresources.NewProvisioningKeyEphemeralResource,
		ephemeralresources.NewPRACredentialControllerEphemeralResource,
	}
}

func New(version string) provider.Provider {
	return &ZPAProvider{
		version: version,
	}
}
