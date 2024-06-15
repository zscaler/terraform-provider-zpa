package zpa

import (
	"log"

	gozscaler "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/scimgroup"
)

func init() {
	// remove timestamp from Zscaler provider logger, use the timestamp from the default terraform logger
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

type Client struct {
	AppConnectorGroup            *services.Service
	AppConnectorController       *services.Service
	AppConnectorSchedule         *services.Service
	ApplicationSegment           *services.Service
	ApplicationSegmentPRA        *services.Service
	ApplicationSegmentInspection *services.Service
	ApplicationSegmentByType     *services.Service
	AppServerController          *services.Service
	BACertificate                *services.Service
	BrowserAccess                *services.Service
	CBIRegions                   *services.Service
	CBIProfileController         *services.Service
	CBIZpaProfile                *services.Service
	CBICertificateController     *services.Service
	CBIBannerController          *services.Service
	CloudConnectorGroup          *services.Service
	CustomerVersionProfile       *services.Service
	ClientTypes                  *services.Service
	EmergencyAccess              *services.Service
	EnrollmentCert               *services.Service
	IDPController                *services.Service
	InspectionCustomControls     *services.Service
	InspectionPredefinedControls *services.Service
	InspectionProfile            *services.Service
	IsolationProfile             *services.Service
	LSSConfigController          *services.Service
	MachineGroup                 *services.Service
	MicroTenants                 *services.Service
	Platforms                    *services.Service
	PolicySetController          *services.Service
	PolicySetControllerV2        *services.Service
	PostureProfile               *services.Service
	PRAApproval                  *services.Service
	PRAConsole                   *services.Service
	PRACredential                *services.Service
	PRAPortal                    *services.Service
	ProvisioningKey              *services.Service
	SAMLAttribute                *services.Service
	ScimGroup                    *scimgroup.Service
	ScimAttributeHeader          *services.Service
	SegmentGroup                 *services.Service
	ServerGroup                  *services.Service
	ServiceEdgeGroup             *services.Service
	ServiceEdgeSchedule          *services.Service
	ServiceEdgeController        *services.Service
	TrustedNetwork               *services.Service
}

type Config struct {
	ClientID     string
	ClientSecret string
	CustomerID   string
	BaseURL      string
	UserAgent    string
}

func (c *Config) Client() (*Client, error) {
	config, err := gozscaler.NewConfig(c.ClientID, c.ClientSecret, c.CustomerID, c.BaseURL, c.UserAgent)
	if err != nil {
		return nil, err
	}
	zpaClient := gozscaler.NewClient(config)

	client := &Client{
		AppConnectorGroup:            services.New(zpaClient),
		AppConnectorController:       services.New(zpaClient),
		AppConnectorSchedule:         services.New(zpaClient),
		ApplicationSegment:           services.New(zpaClient),
		ApplicationSegmentPRA:        services.New(zpaClient),
		ApplicationSegmentInspection: services.New(zpaClient),
		ApplicationSegmentByType:     services.New(zpaClient),
		AppServerController:          services.New(zpaClient),
		BACertificate:                services.New(zpaClient),
		BrowserAccess:                services.New(zpaClient),
		CBIRegions:                   services.New(zpaClient),
		CBIProfileController:         services.New(zpaClient),
		CBIZpaProfile:                services.New(zpaClient),
		CBICertificateController:     services.New(zpaClient),
		CBIBannerController:          services.New(zpaClient),
		CloudConnectorGroup:          services.New(zpaClient),
		CustomerVersionProfile:       services.New(zpaClient),
		ClientTypes:                  services.New(zpaClient),
		EmergencyAccess:              services.New(zpaClient),
		EnrollmentCert:               services.New(zpaClient),
		IDPController:                services.New(zpaClient),
		InspectionCustomControls:     services.New(zpaClient),
		InspectionPredefinedControls: services.New(zpaClient),
		InspectionProfile:            services.New(zpaClient),
		IsolationProfile:             services.New(zpaClient),
		LSSConfigController:          services.New(zpaClient),
		MachineGroup:                 services.New(zpaClient),
		MicroTenants:                 services.New(zpaClient),
		Platforms:                    services.New(zpaClient),
		PolicySetController:          services.New(zpaClient), // Correct initialization
		PolicySetControllerV2:        services.New(zpaClient), // Correct initialization
		PostureProfile:               services.New(zpaClient),
		PRAApproval:                  services.New(zpaClient),
		PRAConsole:                   services.New(zpaClient),
		PRACredential:                services.New(zpaClient),
		PRAPortal:                    services.New(zpaClient),
		ProvisioningKey:              services.New(zpaClient),
		SAMLAttribute:                services.New(zpaClient),
		ScimGroup:                    scimgroup.New(zpaClient),
		ScimAttributeHeader:          services.New(zpaClient),
		SegmentGroup:                 services.New(zpaClient),
		ServerGroup:                  services.New(zpaClient),
		ServiceEdgeGroup:             services.New(zpaClient),
		ServiceEdgeSchedule:          services.New(zpaClient),
		ServiceEdgeController:        services.New(zpaClient),
		TrustedNetwork:               services.New(zpaClient),
	}

	log.Println("[INFO] initialized ZPA client")
	return client, nil
}
