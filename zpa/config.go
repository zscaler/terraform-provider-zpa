package zpa

import (
	"log"

	gozscaler "github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/appconnectorcontroller"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/applicationsegmentinspection"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/applicationsegmentpra"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/appservercontroller"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/bacertificate"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/browseraccess"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/clienttypes"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/cloudbrowserisolation/isolationprofile"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/cloudconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/customerversionprofile"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/enrollmentcert"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/inspectioncontrol/inspection_custom_controls"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/inspectioncontrol/inspection_predefined_controls"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/inspectioncontrol/inspection_profile"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/lssconfigcontroller"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/machinegroup"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/microtenants"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/platforms"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/postureprofile"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/provisioningkey"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/samlattribute"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/scimattributeheader"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/scimgroup"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/segmentgroup"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/servergroup"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/serviceedgecontroller"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/serviceedgegroup"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/trustednetwork"
)

func init() {
	// remove timestamp from Zscaler provider logger, use the timestamp from the default terraform logger
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

type Client struct {
	appconnectorgroup              appconnectorgroup.Service
	appconnectorcontroller         appconnectorcontroller.Service
	applicationsegment             applicationsegment.Service
	applicationsegmentpra          applicationsegmentpra.Service
	applicationsegmentinspection   applicationsegmentinspection.Service
	appservercontroller            appservercontroller.Service
	bacertificate                  bacertificate.Service
	cloudconnectorgroup            cloudconnectorgroup.Service
	customerversionprofile         customerversionprofile.Service
	enrollmentcert                 enrollmentcert.Service
	idpcontroller                  idpcontroller.Service
	lssconfigcontroller            lssconfigcontroller.Service
	machinegroup                   machinegroup.Service
	microtenants                   microtenants.Service
	postureprofile                 postureprofile.Service
	isolationprofile               isolationprofile.Service
	policysetcontroller            policysetcontroller.Service
	provisioningkey                provisioningkey.Service
	samlattribute                  samlattribute.Service
	scimgroup                      scimgroup.Service
	scimattributeheader            scimattributeheader.Service
	segmentgroup                   segmentgroup.Service
	servergroup                    servergroup.Service
	serviceedgegroup               serviceedgegroup.Service
	serviceedgecontroller          serviceedgecontroller.Service
	trustednetwork                 trustednetwork.Service
	platforms                      platforms.Service
	clienttypes                    clienttypes.Service
	browseraccess                  browseraccess.Service
	inspection_custom_controls     inspection_custom_controls.Service
	inspection_predefined_controls inspection_predefined_controls.Service
	inspection_profile             inspection_profile.Service
}

type Config struct {

	// ZPA Client ID for API Client
	ClientID string

	// ZPA Client Secret for API Client
	ClientSecret string

	// ZPA Customer ID for API Client
	CustomerID string

	// ZPA Base URL for API Client
	BaseURL string

	// UserAgent for API Client
	UserAgent string
}

func (c *Config) Client() (*Client, error) {
	config, err := gozscaler.NewConfig(c.ClientID, c.ClientSecret, c.CustomerID, c.BaseURL, c.UserAgent)
	if err != nil {
		return nil, err
	}
	zpaClient := gozscaler.NewClient(config)
	client := &Client{
		appconnectorgroup:              *appconnectorgroup.New(zpaClient),
		appconnectorcontroller:         *appconnectorcontroller.New(zpaClient),
		applicationsegment:             *applicationsegment.New(zpaClient),
		applicationsegmentpra:          *applicationsegmentpra.New(zpaClient),
		applicationsegmentinspection:   *applicationsegmentinspection.New(zpaClient),
		appservercontroller:            *appservercontroller.New(zpaClient),
		bacertificate:                  *bacertificate.New(zpaClient),
		cloudconnectorgroup:            *cloudconnectorgroup.New(zpaClient),
		customerversionprofile:         *customerversionprofile.New(zpaClient),
		enrollmentcert:                 *enrollmentcert.New(zpaClient),
		idpcontroller:                  *idpcontroller.New(zpaClient),
		lssconfigcontroller:            *lssconfigcontroller.New(zpaClient),
		machinegroup:                   *machinegroup.New(zpaClient),
		microtenants:                   *microtenants.New(zpaClient),
		postureprofile:                 *postureprofile.New(zpaClient),
		isolationprofile:               *isolationprofile.New(zpaClient),
		policysetcontroller:            *policysetcontroller.New(zpaClient),
		provisioningkey:                *provisioningkey.New(zpaClient),
		samlattribute:                  *samlattribute.New(zpaClient),
		scimgroup:                      *scimgroup.New(zpaClient),
		scimattributeheader:            *scimattributeheader.New(zpaClient),
		segmentgroup:                   *segmentgroup.New(zpaClient),
		servergroup:                    *servergroup.New(zpaClient),
		serviceedgegroup:               *serviceedgegroup.New(zpaClient),
		serviceedgecontroller:          *serviceedgecontroller.New(zpaClient),
		trustednetwork:                 *trustednetwork.New(zpaClient),
		platforms:                      *platforms.New(zpaClient),
		clienttypes:                    *clienttypes.New(zpaClient),
		browseraccess:                  *browseraccess.New(zpaClient),
		inspection_custom_controls:     *inspection_custom_controls.New(zpaClient),
		inspection_predefined_controls: *inspection_predefined_controls.New(zpaClient),
		inspection_profile:             *inspection_profile.New(zpaClient),
	}

	log.Println("[INFO] initialized ZPA client")
	return client, nil
}
