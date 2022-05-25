package applicationsegmentinspection

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig                   = "/mgmtconfig/v1/admin/customers/"
	appSegmentInspectionEndpoint = "/application"
)

type AppSegmentInspection struct {
	ID                        string                `json:"id,omitempty"`
	SegmentGroupID            string                `json:"segmentGroupId,omitempty"`
	SegmentGroupName          string                `json:"segmentGroupName,omitempty"`
	BypassType                string                `json:"bypassType,omitempty"`
	ConfigSpace               string                `json:"configSpace,omitempty"`
	DomainNames               []string              `json:"domainNames,omitempty"`
	Name                      string                `json:"name,omitempty"`
	Description               string                `json:"description,omitempty"`
	Enabled                   bool                  `json:"enabled"`
	ICMPAccessType            string                `json:"icmpAccessType,omitempty"`
	PassiveHealthEnabled      bool                  `json:"passiveHealthEnabled"`
	SelectConnectorCloseToApp bool                  `json:"selectConnectorCloseToApp"`
	DoubleEncrypt             bool                  `json:"doubleEncrypt"`
	HealthCheckType           string                `json:"healthCheckType,omitempty"`
	IsCnameEnabled            bool                  `json:"isCnameEnabled"`
	IPAnchored                bool                  `json:"ipAnchored"`
	HealthReporting           string                `json:"healthReporting,omitempty"`
	CreationTime              string                `json:"creationTime,omitempty"`
	ModifiedBy                string                `json:"modifiedBy,omitempty"`
	ModifiedTime              string                `json:"modifiedTime,omitempty"`
	TCPPortRanges             []string              `json:"tcpPortRanges,omitempty"`
	UDPPortRanges             []string              `json:"udpPortRanges,omitempty"`
	TCPAppPortRange           []common.NetworkPorts `json:"tcpPortRange,omitempty"`
	UDPAppPortRange           []common.NetworkPorts `json:"udpPortRange,omitempty"`
	InspectionAppDto          []InspectionAppDto    `json:"inspectionApps,omitempty"`
	CommonAppsDto             CommonAppsDto         `json:"commonAppsDto,omitempty"`
	AppServerGroups           []AppServerGroups     `json:"serverGroups,omitempty"`
}

type CommonAppsDto struct {
	AppsConfig []AppsConfig `json:"appsConfig,omitempty"`
}

type AppsConfig struct {
	Name                string `json:"name,omitempty"`
	AllowOptions        bool   `json:"allowOptions"`
	AppID               string `json:"appId,omitempty"`
	AppTypes            string `json:"appTypes,omitempty"`
	ApplicationPort     string `json:"applicationPort,omitempty"`
	ApplicationProtocol string `json:"applicationProtocol,omitempty"`
	InspectAppID        string `json:"inspectAppId,omitempty"`
	CertificateID       string `json:"certificateId,omitempty"`
	CertificateName     string `json:"certificateName,omitempty"`
	Cname               string `json:"cname,omitempty"`
	Description         string `json:"description,omitempty"`
	Domain              string `json:"domain,omitempty"`
	Enabled             bool   `json:"enabled"`
	Hidden              bool   `json:"hidden"`
	LocalDomain         string `json:"localDomain,omitempty"`
	Portal              bool   `json:"portal"`
}

type InspectionAppDto struct {
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

type AppServerGroups struct {
	ID string `json:"id"`
}

func (service *Service) Get(id string) (*AppSegmentInspection, *http.Response, error) {
	v := new(AppSegmentInspection)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentInspectionEndpoint, id)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(BaName string) (*AppSegmentInspection, *http.Response, error) {
	var v struct {
		List []AppSegmentInspection `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appSegmentInspectionEndpoint
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

func (service *Service) Create(appSegmentPra AppSegmentInspection) (*AppSegmentInspection, *http.Response, error) {
	v := new(AppSegmentInspection)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+appSegmentInspectionEndpoint, nil, appSegmentPra, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(id string, appSegmentPra *AppSegmentInspection) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentInspectionEndpoint, id)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, appSegmentPra, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(id string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentInspectionEndpoint, id)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
