package applicationsegmentpra

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig            = "/mgmtconfig/v1/admin/customers/"
	browserAccessEndpoint = "/application"
)

type BrowserAccess struct {
	ID                   string                `json:"id,omitempty"`
	SegmentGroupID       string                `json:"segmentGroupId,omitempty"`
	SegmentGroupName     string                `json:"segmentGroupName,omitempty"`
	BypassType           string                `json:"bypassType,omitempty"`
	ConfigSpace          string                `json:"configSpace,omitempty"`
	DomainNames          []string              `json:"domainNames,omitempty"`
	Name                 string                `json:"name,omitempty"`
	Description          string                `json:"description,omitempty"`
	Enabled              bool                  `json:"enabled"`
	PassiveHealthEnabled bool                  `json:"passiveHealthEnabled"`
	DoubleEncrypt        bool                  `json:"doubleEncrypt"`
	HealthCheckType      string                `json:"healthCheckType,omitempty"`
	IsCnameEnabled       bool                  `json:"isCnameEnabled"`
	IPAnchored           bool                  `json:"ipAnchored"`
	HealthReporting      string                `json:"healthReporting,omitempty"`
	CreationTime         string                `json:"creationTime,omitempty"`
	ModifiedBy           string                `json:"modifiedBy,omitempty"`
	ModifiedTime         string                `json:"modifiedTime,omitempty"`
	TCPPortRanges        []string              `json:"tcpPortRanges,omitempty"`
	UDPPortRanges        []string              `json:"udpPortRanges,omitempty"`
	TCPAppPortRange      []common.NetworkPorts `json:"tcpPortRange,omitempty"`
	UDPAppPortRange      []common.NetworkPorts `json:"udpPortRange,omitempty"`
	SRAAppsDto           []SRAAppsDto          `json:"sraApps,omitempty"`
	CommonAppsDto        CommonAppsDto         `json:"clientlessApps,omitempty"`
	AppServerGroups      []AppServerGroups     `json:"serverGroups,omitempty"`
}

type CommonAppsDto struct {
	AppConfig []AppConfig `json:"appsConfig,omitempty"`
}

type AppConfig struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	AllowOptions        bool   `json:"allowOptions"`
	AppID               string `json:"appId,omitempty"`
	AppTypes            string `json:"appTypes,omitempty"`
	ApplicationPort     string `json:"applicationPort,omitempty"`
	ApplicationProtocol string `json:"applicationProtocol,omitempty"`
	BaAppID             string `json:"baAppId,omitempty"`
	CertificateID       string `json:"certificateId,omitempty"`
	CertificateName     string `json:"certificateName,omitempty"`
	Cname               string `json:"cname,omitempty"`
	Description         string `json:"description,omitempty"`
	Domain              string `json:"domain,omitempty"`
	Enabled             bool   `json:"enabled"`
	Hidden              bool   `json:"hidden"`
	InspectAppId        string `json:"inspectAppId,omitempty"`
	LocalDomain         string `json:"localDomain,omitempty"`
	Path                string `json:"path,omitempty"`
	Portal              string `json:"portal,omitempty"`
	TrustUntrustedCert  bool   `json:"trustUntrustedCert"`
}

type SRAAppsDto struct {
	AppID               string `json:"appId,omitempty"`
	ApplicationPort     string `json:"applicationPort,omitempty"`
	ApplicationProtocol string `json:"applicationProtocol,omitempty"`
	CertificateID       string `json:"certificateId,omitempty"`
	CertificateName     string `json:"certificateName,omitempty"`
	Description         string `json:"description,omitempty"`
	Domain              string `json:"domain,omitempty"`
	Enabled             bool   `json:"enabled"`
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
}

func (service *Service) Get(id string) (*BrowserAccess, *http.Response, error) {
	v := new(BrowserAccess)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+browserAccessEndpoint, id)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(BaName string) (*BrowserAccess, *http.Response, error) {
	var v struct {
		List []BrowserAccess `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + browserAccessEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize, Search: BaName}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, app := range v.List {
		if strings.EqualFold(app.Name, BaName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no browser access application named '%s' was found", BaName)
}

func (service *Service) Create(browserAccess BrowserAccess) (*BrowserAccess, *http.Response, error) {
	v := new(BrowserAccess)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+browserAccessEndpoint, nil, browserAccess, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(id string, browserAccess *BrowserAccess) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+browserAccessEndpoint, id)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, browserAccess, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(id string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+browserAccessEndpoint, id)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
