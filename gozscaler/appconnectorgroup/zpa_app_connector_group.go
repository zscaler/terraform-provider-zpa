package appconnectorgroup

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	mgmtConfig                = "/mgmtconfig/v1/admin/customers/"
	appConnectorGroupEndpoint = "/appConnectorGroup"
)

type AppConnectorGroup struct {
	CityCountry           string           `json:"cityCountry,omitempty"`
	CountryCode           string           `json:"countryCode,omitempty"`
	CreationTime          string           `json:"creationTime,omitempty"`
	Description           string           `json:"description,omitempty"`
	DNSQueryType          string           `json:"dnsQueryType,omitempty"`
	Enabled               bool             `json:"enabled,omitempty"`
	GeoLocationID         string           `json:"geoLocationId,omitempty"`
	ID                    string           `json:"id,omitempty"`
	Latitude              string           `json:"latitude,omitempty"`
	Location              string           `json:"location,omitempty"`
	Longitude             string           `json:"longitude,omitempty"`
	ModifiedBy            string           `json:"modifiedBy,omitempty"`
	ModifiedTime          string           `json:"modifiedTime,omitempty"`
	Name                  string           `json:"name,omitempty"`
	SiemAppConnectorGroup bool             `json:"siemAppConnectorGroup,omitempty"`
	UpgradeDay            string           `json:"upgradeDay,omitempty"`
	UpgradeTimeInSecs     string           `json:"upgradeTimeInSecs,omitempty"`
	VersionProfileID      string           `json:"versionProfileId,omitempty"`
	AppServerGroup        []AppServerGroup `json:"serverGroups,omitempty"`
	Connectors            []*Connector     `json:"connectors,omitempty"`
}
type Connector struct {
	ApplicationStartTime     string                 `json:"applicationStartTime,omitempty"`
	AppConnectorGroupID      string                 `json:"appConnectorGroupId,omitempty"`
	AppConnectorGroupName    string                 `json:"appConnectorGroupName,omitempty"`
	ControlChannelStatus     string                 `json:"controlChannelStatus,omitempty"`
	CreationTime             string                 `json:"creationTime,omitempty"`
	CtrlBrokerName           string                 `json:"ctrlBrokerName,omitempty"`
	CurrentVersion           string                 `json:"currentVersion,omitempty"`
	Description              string                 `json:"description,omitempty"`
	Enabled                  bool                   `json:"enabled,omitempty"`
	ExpectedUpgradeTime      string                 `json:"expectedUpgradeTime,omitempty"`
	ExpectedVersion          string                 `json:"expectedVersion,omitempty"`
	Fingerprint              string                 `json:"fingerprint,omitempty"`
	ID                       string                 `json:"id,omitempty"`
	IPACL                    string                 `json:"ipAcl,omitempty"`
	IssuedCertID             string                 `json:"issuedCertId,omitempty"`
	LastBrokerConnectTime    string                 `json:"lastBrokerConnectTime,omitempty"`
	LastBrokerDisconnectTime string                 `json:"lastBrokerDisconnectTime,omitempty"`
	LastUpgradeTime          string                 `json:"lastUpgradeTime,omitempty"`
	Latitude                 float64                `json:"latitude,omitempty"`
	Location                 string                 `json:"location,omitempty"`
	Longitude                float64                `json:"longitude,omitempty"`
	ModifiedBy               string                 `json:"modifiedBy,omitempty"`
	ModifiedTime             string                 `json:"modifiedTime,omitempty"`
	Name                     string                 `json:"name,omitempty"`
	Platform                 string                 `json:"platform,omitempty"`
	PreviousVersion          string                 `json:"previousVersion,omitempty"`
	PrivateIP                string                 `json:"privateIp,omitempty"`
	SigningCert              map[string]interface{} `json:"signingCert,omitempty"`
	PublicIP                 string                 `json:"publicIp,omitempty"`
	UpgradeAttempt           string                 `json:"upgradeAttempt,omitempty"`
	UpgradeStatus            string                 `json:"upgradeStatus,omitempty"`
}
type AppServerGroup struct {
	ConfigSpace      string `json:"configSpace,omitempty"`
	CreationTime     string `json:"creationTime,omitempty"`
	Description      string `json:"description,omitempty"`
	Enabled          bool   `json:"enabled,omitempty"`
	ID               string `json:"id,omitempty"`
	DynamicDiscovery bool   `json:"dynamicDiscovery,omitempty"`
	ModifiedBy       string `json:"modifiedBy,omitempty"`
	ModifiedTime     string `json:"modifiedTime,omitempty"`
	Name             string `json:"name,omitempty"`
}

func (service *Service) Get(appConnectorGroupID string) (*AppConnectorGroup, *http.Response, error) {
	v := new(AppConnectorGroup)
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+appConnectorGroupEndpoint, appConnectorGroupID)
	resp, err := service.Client.NewRequestDo("GET", path, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(appConnectorGroupName string) (*AppConnectorGroup, *http.Response, error) {
	var v struct {
		List []AppConnectorGroup `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appConnectorGroupEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, struct{ pagesize int }{
		pagesize: 500,
	}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, app := range v.List {
		if strings.EqualFold(app.Name, appConnectorGroupName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no application named '%s' was found", appConnectorGroupName)
}
