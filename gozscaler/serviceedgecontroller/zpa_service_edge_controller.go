package serviceedgecontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig                    = "/mgmtconfig/v1/admin/customers/"
	serviceEdgeControllerEndpoint = "/serviceEdge"
)

type ServiceEdgeController struct {
	ApplicationStartTime             string                 `json:"applicationStartTime,omitempty"`
	ServiceEdgeGroupID               string                 `json:"serviceEdgeGroupId,omitempty"`
	ServiceEdgeGroupName             string                 `json:"serviceEdgeGroupName,omitempty"`
	ControlChannelStatus             string                 `json:"controlChannelStatus,omitempty"`
	CreationTime                     string                 `json:"creationTime,omitempty"`
	CtrlBrokerName                   string                 `json:"ctrlBrokerName,omitempty"`
	CurrentVersion                   string                 `json:"currentVersion,omitempty"`
	Description                      string                 `json:"description,omitempty"`
	Enabled                          bool                   `json:"enabled,omitempty"`
	ExpectedUpgradeTime              string                 `json:"expectedUpgradeTime,omitempty"`
	ExpectedVersion                  string                 `json:"expectedVersion,omitempty"`
	Fingerprint                      string                 `json:"fingerprint,omitempty"`
	ID                               string                 `json:"id,omitempty"`
	IPACL                            string                 `json:"ipAcl,omitempty"`
	IssuedCertID                     string                 `json:"issuedCertId,omitempty"`
	LastBrokerConnectTime            string                 `json:"lastBrokerConnectTime,omitempty"`
	LastBrokerConnectTimeDuration    string                 `json:"lastBrokerConnectTimeDuration,omitempty"`
	LastBrokerDisconnectTime         string                 `json:"lastBrokerDisconnectTime,omitempty"`
	LastBrokerDisconnectTimeDuration string                 `json:"lastBrokerDisconnectTimeDuration,omitempty"`
	LastUpgradeTime                  string                 `json:"lastUpgradeTime,omitempty"`
	Latitude                         string                 `json:"latitude,omitempty"`
	Location                         string                 `json:"location,omitempty"`
	Longitude                        string                 `json:"longitude,omitempty"`
	ListenIPs                        string                 `json:"listenIps,omitempty"`
	ModifiedBy                       string                 `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                 `json:"modifiedTime,omitempty"`
	Name                             string                 `json:"name,omitempty"`
	ProvisioningKeyID                string                 `json:"provisioningKeyId"`
	ProvisioningKeyName              string                 `json:"provisioningKeyName"`
	Platform                         string                 `json:"platform,omitempty"`
	PreviousVersion                  string                 `json:"previousVersion,omitempty"`
	PrivateIP                        string                 `json:"privateIp,omitempty"`
	PublicIP                         string                 `json:"publicIp,omitempty"`
	PublishIPs                       string                 `json:"publishIps,omitempty"`
	SargeVersion                     string                 `json:"sargeVersion,omitempty"`
	EnrollmentCert                   map[string]interface{} `json:"enrollmentCert,omitempty"`
	UpgradeAttempt                   string                 `json:"upgradeAttempt,omitempty"`
	UpgradeStatus                    string                 `json:"upgradeStatus,omitempty"`
}

func (service *Service) Get(serviceEdgeID string) (*ServiceEdgeController, *http.Response, error) {
	v := new(ServiceEdgeController)
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+serviceEdgeControllerEndpoint, serviceEdgeID)
	resp, err := service.Client.NewRequestDo("GET", path, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(serviceEdgeName string) (*ServiceEdgeController, *http.Response, error) {
	var v struct {
		List []ServiceEdgeController `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + serviceEdgeControllerEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize, Search: serviceEdgeName}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, service := range v.List {
		if strings.EqualFold(service.Name, serviceEdgeName) {
			return &service, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no service edge named '%s' was found", serviceEdgeName)
}
