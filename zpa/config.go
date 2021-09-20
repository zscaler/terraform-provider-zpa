package zpa

import (
	"log"

	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/appconnectorgroup"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/applicationsegment"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/appservercontroller"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/bacertificate"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/browseraccess"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/cloudconnectorgroup"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/idpcontroller"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/machinegroup"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/policysetglobal"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/policysetrule"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/postureprofile"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/samlattribute"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/scimattributeheader"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/scimgroup"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/segmentgroup"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/servergroup"
	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/trustednetwork"
)

func init() {
	// remove timestamp from Zscaler provider logger, use the timestamp from the default terraform logger
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

type Client struct {
	appconnectorgroup   appconnectorgroup.Service
	applicationsegment  applicationsegment.Service
	appservercontroller appservercontroller.Service
	bacertificate       bacertificate.Service
	cloudconnectorgroup cloudconnectorgroup.Service
	idpcontroller       idpcontroller.Service
	machinegroup        machinegroup.Service
	postureprofile      postureprofile.Service
	policysetglobal     policysetglobal.Service
	policysetrule       policysetrule.Service
	samlattribute       samlattribute.Service
	scimgroup           scimgroup.Service
	scimattributeheader scimattributeheader.Service
	segmentgroup        segmentgroup.Service
	servergroup         servergroup.Service
	trustednetwork      trustednetwork.Service
	browseraccess       browseraccess.Service
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
		appconnectorgroup:   *appconnectorgroup.New(config),
		applicationsegment:  *applicationsegment.New(config),
		appservercontroller: *appservercontroller.New(config),
		bacertificate:       *bacertificate.New(config),
		cloudconnectorgroup: *cloudconnectorgroup.New(config),
		idpcontroller:       *idpcontroller.New(config),
		machinegroup:        *machinegroup.New(config),
		postureprofile:      *postureprofile.New(config),
		policysetglobal:     *policysetglobal.New(config),
		policysetrule:       *policysetrule.New(config),
		samlattribute:       *samlattribute.New(config),
		scimgroup:           *scimgroup.New(config),
		scimattributeheader: *scimattributeheader.New(config),
		segmentgroup:        *segmentgroup.New(config),
		servergroup:         *servergroup.New(config),
		trustednetwork:      *trustednetwork.New(config),
		browseraccess:       *browseraccess.New(config),
	}

	log.Println("[INFO] initialized ZPA client")
	return client, nil
}
