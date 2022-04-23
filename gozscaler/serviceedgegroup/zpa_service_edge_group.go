package serviceedgegroup

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig               = "/mgmtconfig/v1/admin/customers/"
	serviceEdgeGroupEndpoint = "/serviceEdgeGroup"
)

type ServiceEdgeGroup struct {
	CityCountry                   string            `json:"cityCountry,omitempty"`
	CountryCode                   string            `json:"countryCode,omitempty"`
	CreationTime                  string            `json:"creationTime"`
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
	ServiceEdges                  []ServiceEdges    `json:"serviceEdges,omitempty"`
	TrustedNetworks               []TrustedNetworks `json:"trustedNetworks,omitempty"`
	UpgradeDay                    string            `json:"upgradeDay,omitempty"`
	UpgradeTimeInSecs             string            `json:"upgradeTimeInSecs,omitempty"`
	VersionProfileID              string            `json:"versionProfileId,omitempty"`
	VersionProfileName            string            `json:"versionProfileName,omitempty"`
	VersionProfileVisibilityScope string            `json:"versionProfileVisibilityScope,omitempty"`
}

type ServiceEdges struct {
	ApplicationStartTime             string                 `json:"applicationStartTime,omitempty"`
	ControlChannelStatus             string                 `json:"controlChannelStatus,omitempty"`
	CreationTime                     string                 `json:"creationTime"`
	CtrlBrokerName                   string                 `json:"ctrlBrokerName,omitempty"`
	CurrentVersion                   string                 `json:"currentVersion,omitempty"`
	Description                      string                 `json:"description,omitempty"`
	Enabled                          bool                   `json:"enabled"`
	ExpectedUpgradeTime              string                 `json:"expectedUpgradeTime,omitempty"`
	ExpectedVersion                  string                 `json:"expectedVersion,omitempty"`
	Fingerprint                      string                 `json:"fingerprint,omitempty"`
	ID                               string                 `json:"id,omitempty"`
	IPACL                            []string               `json:"ipAcl,omitempty"`
	IssuedCertID                     string                 `json:"issuedCertId,omitempty"`
	LastBrokerConnectTime            string                 `json:"lastBrokerConnectTime,omitempty"`
	LastBrokerConnectTimeDuration    string                 `json:"lastBrokerConnectTimeDuration,omitempty"`
	LastBrokerDisconnectTime         string                 `json:"lastBrokerDisconnectTime,omitempty"`
	LastBrokerDisconnectTimeDuration string                 `json:"lastBrokerDisconnectTimeDuration,omitempty"`
	LastUpgradeTime                  string                 `json:"lastUpgradeTime,omitempty"`
	Latitude                         string                 `json:"latitude,omitempty"`
	ListenIPs                        []string               `json:"listenIps,omitempty"`
	Location                         string                 `json:"location,omitempty"`
	Longitude                        string                 `json:"longitude,omitempty"`
	ModifiedBy                       string                 `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                 `json:"modifiedTime,omitempty"`
	Name                             string                 `json:"name,omitempty"`
	ProvisioningKeyID                string                 `json:"provisioningKeyId,omitempty"`
	ProvisioningKeyName              string                 `json:"provisioningKeyName,omitempty"`
	Platform                         string                 `json:"platform,omitempty"`
	PreviousVersion                  string                 `json:"previousVersion,omitempty"`
	ServiceEdgeGroupID               string                 `json:"serviceEdgeGroupId,omitempty"`
	ServiceEdgeGroupName             string                 `json:"serviceEdgeGroupName,omitempty"`
	PrivateIP                        string                 `json:"privateIp,omitempty"`
	PublicIP                         string                 `json:"publicIp,omitempty"`
	PublishIPs                       []string               `json:"publishIps,omitempty"`
	SargeVersion                     string                 `json:"sargeVersion,omitempty"`
	EnrollmentCert                   map[string]interface{} `json:"enrollmentCert"`
	UpgradeAttempt                   string                 `json:"upgradeAttempt,omitempty"`
	UpgradeStatus                    string                 `json:"upgradeStatus,omitempty"`
}
type TrustedNetworks struct {
	CreationTime     string `json:"creationTime,omitempty"`
	Domain           string `json:"domain,omitempty"`
	ID               string `json:"id"`
	MasterCustomerID string `json:"masterCustomerId"`
	ModifiedBy       string `json:"modifiedBy,omitempty"`
	ModifiedTime     string `json:"modifiedTime,omitempty"`
	Name             string `json:"name,omitempty"`
	NetworkID        string `json:"networkId,omitempty"`
	ZscalerCloud     string `json:"zscalerCloud,omitempty"`
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
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize, Search: serviceEdgeGroupName}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, app := range v.List {
		if strings.EqualFold(app.Name, serviceEdgeGroupName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no service edge group named '%s' was found", serviceEdgeGroupName)
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
