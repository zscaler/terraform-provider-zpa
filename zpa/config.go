package zpa

import (
	"log"

	"github.com/willguibr/terraform-provider-zpa/gozscaler"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/appconnectorcontroller"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/appconnectorgroup"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/applicationsegment"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/appservercontroller"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/bacertificate"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/browseraccess"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/cloudconnectorgroup"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/customerversionprofile"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/enrollmentcert"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/idpcontroller"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/lssconfigcontroller"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/machinegroup"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/policysetcontroller"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/postureprofile"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/provisioningkey"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/samlattribute"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/scimattributeheader"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/scimgroup"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/segmentgroup"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/servergroup"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/serviceedgecontroller"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/serviceedgegroup"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/trustednetwork"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/inspectioncontrol/inspection_custom_controls"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/inspectioncontrol/inspection_predefined_controls"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/inspectioncontrol/inspection_profile"
)

func init() {
	// remove timestamp from Zscaler provider logger, use the timestamp from the default terraform logger
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

type Client struct {
	appconnectorgroup              appconnectorgroup.Service
	appconnectorcontroller         appconnectorcontroller.Service
	applicationsegment             applicationsegment.Service
	appservercontroller            appservercontroller.Service
	bacertificate                  bacertificate.Service
	cloudconnectorgroup            cloudconnectorgroup.Service
	customerversionprofile         customerversionprofile.Service
	enrollmentcert                 enrollmentcert.Service
	idpcontroller                  idpcontroller.Service
	lssconfigcontroller            lssconfigcontroller.Service
	machinegroup                   machinegroup.Service
	postureprofile                 postureprofile.Service
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
	browseraccess                  browseraccess.Service
	inspection_profile             inspection_profile.Service
	inspection_custom_controls     inspection_custom_controls.Service
	inspection_predefined_controls inspection_predefined_controls.Service
}

type Config struct {
	ClientID     string
	ClientSecret string
	CustomerID   string
	BaseURL      string
}

func (c *Config) Client() (*Client, error) {
	config, err := gozscaler.NewConfig(c.ClientID, c.ClientSecret, c.CustomerID, c.BaseURL)
	if err != nil {
		return nil, err
	}

	client := &Client{
		appconnectorgroup:              *appconnectorgroup.New(config),
		appconnectorcontroller:         *appconnectorcontroller.New(config),
		applicationsegment:             *applicationsegment.New(config),
		appservercontroller:            *appservercontroller.New(config),
		bacertificate:                  *bacertificate.New(config),
		cloudconnectorgroup:            *cloudconnectorgroup.New(config),
		customerversionprofile:         *customerversionprofile.New(config),
		enrollmentcert:                 *enrollmentcert.New(config),
		idpcontroller:                  *idpcontroller.New(config),
		lssconfigcontroller:            *lssconfigcontroller.New(config),
		machinegroup:                   *machinegroup.New(config),
		postureprofile:                 *postureprofile.New(config),
		policysetcontroller:            *policysetcontroller.New(config),
		provisioningkey:                *provisioningkey.New(config),
		samlattribute:                  *samlattribute.New(config),
		scimgroup:                      *scimgroup.New(config),
		scimattributeheader:            *scimattributeheader.New(config),
		segmentgroup:                   *segmentgroup.New(config),
		servergroup:                    *servergroup.New(config),
		serviceedgegroup:               *serviceedgegroup.New(config),
		serviceedgecontroller:          *serviceedgecontroller.New(config),
		trustednetwork:                 *trustednetwork.New(config),
		browseraccess:                  *browseraccess.New(config),
		inspection_profile:             *inspection_profile.New(config),
		inspection_custom_controls:     *inspection_custom_controls.New(config),
		inspection_predefined_controls: *inspection_predefined_controls.New(config),
	}

	log.Println("[INFO] initialized ZPA client")
	return client, nil
}
