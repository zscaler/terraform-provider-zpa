package serviceedgegroup

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	mgmtConfig               = "/mgmtconfig/v1/admin/customers/"
	serviceEdgeGroupEndpoint = "/serviceEdgeGroup"
)

type ServiceEdgeGroup struct {
	CityCountry                   string            `json:"cityCountry,omitempty"`
	CountryCode                   string            `json:"countryCode,omitempty"`
	CreationTime                  string            `json:"creationTime,omitempty"`
	Description                   string            `json:"description,omitempty"`
	Enabled                       bool              `json:"enabled"`
	GeoLocationID                 string            `json:"geoLocationId,omitempty"`
	ID                            string            `json:"id"`
	IsPublic                      string            `json:"isPublic,omitempty"`
	Latitude                      string            `json:"latitude,omitempty"`
	Location                      string            `json:"location,omitempty"`
	Longitude                     string            `json:"longitude,omitempty"`
	ModifiedBy                    string            `json:"modifiedBy,omitempty"`
	ModifiedTime                  string            `json:"modifiedTime,omitempty"`
	Name                          string            `json:"name,omitempty"`
	OverrideVersionProfile        bool              `json:"overrideVersionProfile"`
	ServiceEdges                  []ServiceEdges    `json:"serviceEdges"`
	TrustedNetworks               []TrustedNetworks `json:"trustedNetworks"`
	UpgradeDay                    string            `json:"upgradeDay,omitempty"`
	UpgradeTimeInSecs             string            `json:"upgradeTimeInSecs,omitempty"`
	VersionProfileID              string            `json:"versionProfileId,omitempty"`
	VersionProfileName            string            `json:"versionProfileName,omitempty"`
	VersionProfileVisibilityScope string            `json:"versionProfileVisibilityScope,omitempty"`
}

type ServiceEdges struct {
	ApplicationStartTime             string                 `json:"applicationStartTime"`
	ControlChannelStatus             string                 `json:"controlChannelStatus"`
	CreationTime                     string                 `json:"creationTime"`
	CtrlBrokerName                   string                 `json:"ctrlBrokerName"`
	CurrentVersion                   string                 `json:"currentVersion"`
	Description                      string                 `json:"description"`
	Enabled                          bool                   `json:"enabled"`
	ExpectedUpgradeTime              string                 `json:"expectedUpgradeTime"`
	ExpectedVersion                  string                 `json:"expectedVersion"`
	Fingerprint                      string                 `json:"fingerprint"`
	ID                               string                 `json:"id"`
	IPACL                            []string               `json:"ipAcl"`
	IssuedCertID                     string                 `json:"issuedCertId"`
	LastBrokerConnectTime            string                 `json:"lastBrokerConnectTime"`
	LastBrokerConnectTimeDuration    string                 `json:"lastBrokerConnectTimeDuration"`
	LastBrokerDisconnectTime         string                 `json:"lastBrokerDisconnectTime"`
	LastBrokerDisconnectTimeDuration string                 `json:"lastBrokerDisconnectTimeDuration"`
	LastUpgradeTime                  string                 `json:"lastUpgradeTime"`
	Latitude                         string                 `json:"latitude"`
	ListenIps                        []string               `json:"listenIps"`
	Location                         string                 `json:"location"`
	Longitude                        string                 `json:"longitude"`
	ModifiedBy                       string                 `json:"modifiedBy"`
	ModifiedTime                     string                 `json:"modifiedTime"`
	Name                             string                 `json:"name"`
	ProvisioningKeyID                string                 `json:"provisioningKeyId"`
	ProvisioningKeyName              string                 `json:"provisioningKeyName"`
	Platform                         string                 `json:"platform"`
	PreviousVersion                  string                 `json:"previousVersion"`
	ServiceEdgeGroupID               string                 `json:"serviceEdgeGroupId"`
	ServiceEdgeGroupName             string                 `json:"serviceEdgeGroupName"`
	PrivateIP                        string                 `json:"privateIp"`
	PublicIP                         string                 `json:"publicIp"`
	PublishIps                       []string               `json:"publishIps"`
	SargeVersion                     string                 `json:"sargeVersion"`
	EnrollmentCert                   map[string]interface{} `json:"enrollmentCert"`
	UpgradeAttempt                   string                 `json:"upgradeAttempt"`
	UpgradeStatus                    string                 `json:"upgradeStatus"`
}
type TrustedNetworks struct {
	CreationTime     string `json:"creationTime"`
	Domain           string `json:"domain"`
	ID               string `json:"id"`
	MasterCustomerID string `json:"masterCustomerId"`
	ModifiedBy       string `json:"modifiedBy"`
	ModifiedTime     string `json:"modifiedTime"`
	Name             string `json:"name"`
	NetworkID        string `json:"networkId"`
	ZscalerCloud     string `json:"zscalerCloud"`
}

func (service *Service) Get(serviceEdgeGroupID string) (*ServiceEdgeGroup, *http.Response, error) {
	v := new(ServiceEdgeGroup)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+serviceEdgeGroupEndpoint, serviceEdgeGroupID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(serviceEdgeGroupName string) (*ServiceEdgeGroup, *http.Response, error) {
	var v struct {
		List []ServiceEdgeGroup `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + serviceEdgeGroupEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, struct{ pagesize int }{
		pagesize: 500,
	}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, app := range v.List {
		if strings.EqualFold(app.Name, serviceEdgeGroupName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no server group named '%s' was found", serviceEdgeGroupName)
}

func (service *Service) Create(serviceEdge ServiceEdgeGroup) (*ServiceEdgeGroup, *http.Response, error) {
	v := new(ServiceEdgeGroup)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+serviceEdgeGroupEndpoint, nil, serviceEdge, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(serviceEdgeGroupID string, serviceEdge *ServiceEdgeGroup) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+serviceEdgeGroupEndpoint, serviceEdgeGroupID)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, serviceEdge, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(serviceEdgeGroupID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+serviceEdgeGroupEndpoint, serviceEdgeGroupID)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
