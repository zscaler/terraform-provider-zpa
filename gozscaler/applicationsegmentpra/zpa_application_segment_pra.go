package applicationsegmentpra

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig            = "/mgmtconfig/v1/admin/customers/"
	appSegmentPraEndpoint = "/application"
)

type AppSegmentPRA struct {
	ID                   string                `json:"id,omitempty"`
	DomainNames          []string              `json:"domainNames,omitempty"`
	Name                 string                `json:"name,omitempty"`
	Description          string                `json:"description,omitempty"`
	Enabled              bool                  `json:"enabled"`
	PassiveHealthEnabled bool                  `json:"passiveHealthEnabled"`
	DoubleEncrypt        bool                  `json:"doubleEncrypt"`
	ConfigSpace          string                `json:"configSpace,omitempty"`
	Applications         string                `json:"applications,omitempty"`
	BypassType           string                `json:"bypassType,omitempty"`
	HealthCheckType      string                `json:"healthCheckType,omitempty"`
	IsCnameEnabled       bool                  `json:"isCnameEnabled"`
	IpAnchored           bool                  `json:"ipAnchored"`
	HealthReporting      string                `json:"healthReporting,omitempty"`
	IcmpAccessType       string                `json:"icmpAccessType,omitempty"`
	SegmentGroupID       string                `json:"segmentGroupId"`
	SegmentGroupName     string                `json:"segmentGroupName,omitempty"`
	CreationTime         string                `json:"creationTime,omitempty"`
	ModifiedBy           string                `json:"modifiedBy,omitempty"`
	ModifiedTime         string                `json:"modifiedTime,omitempty"`
	TCPAppPortRange      []common.NetworkPorts `json:"tcpPortRange,omitempty"`
	UDPAppPortRange      []common.NetworkPorts `json:"udpPortRange,omitempty"`
	ServerGroups         []AppServerGroups     `json:"serverGroups,omitempty"`
	DefaultIdleTimeout   string                `json:"defaultIdleTimeout,omitempty"`
	DefaultMaxAge        string                `json:"defaultMaxAge,omitempty"`
	TCPPortRanges        []string              `json:"tcpPortRanges,omitempty"`
	UDPPortRanges        []string              `json:"udpPortRanges,omitempty"`
	SRAAppsDto           []SRAAppsDto          `json:"sraApps,omitempty"`
	CommonAppsDto        CommonAppsDto         `json:"commonAppsDto,omitempty"`
}

type CommonAppsDto struct {
	AppsConfig []AppsConfig `json:"appsConfig,omitempty"`
}

type AppsConfig struct {
	Name                string   `json:"name,omitempty"`
	AllowOptions        bool     `json:"allowOptions"`
	ID                  string   `json:"id,omitempty"`
	AppID               string   `json:"appId,omitempty"`
	AppTypes            []string `json:"appTypes,omitempty"`
	ApplicationPort     string   `json:"applicationPort,omitempty"`
	ApplicationProtocol string   `json:"applicationProtocol,omitempty"`
	// CertificateID       string `json:"certificateId,omitempty"`
	// CertificateName     string `json:"certificateName,omitempty"`
	Cname              string `json:"cname,omitempty"`
	ConnectionSecurity string `json:"connectionSecurity,omitempty"`
	Description        string `json:"description,omitempty"`
	Domain             string `json:"domain,omitempty"`
	Enabled            bool   `json:"enabled"`
	Hidden             bool   `json:"hidden"`
	LocalDomain        string `json:"localDomain,omitempty"`
	Portal             bool   `json:"portal"`
}

type SRAAppsDto struct {
	AppID               string `json:"appId,omitempty"`
	ApplicationPort     string `json:"applicationPort,omitempty"`
	ApplicationProtocol string `json:"applicationProtocol,omitempty"`
	CertificateID       string `json:"certificateId,omitempty"`
	CertificateName     string `json:"certificateName,omitempty"`
	ConnectionSecurity  string `json:"connectionSecurity,omitempty"`
	Hidden              bool   `json:"hidden"`
	Portal              bool   `json:"portal"`
	Description         string `json:"description,omitempty"`
	Domain              string `json:"domain,omitempty"`
	Enabled             bool   `json:"enabled"`
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
}

type AppServerGroups struct {
	ID string `json:"id"`
}

func (service *Service) Get(id string) (*AppSegmentPRA, *http.Response, error) {
	v := new(AppSegmentPRA)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentPraEndpoint, id)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(BaName string) (*AppSegmentPRA, *http.Response, error) {
	var v struct {
		List []AppSegmentPRA `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appSegmentPraEndpoint
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

func (service *Service) Create(appSegmentPra AppSegmentPRA) (*AppSegmentPRA, *http.Response, error) {
	v := new(AppSegmentPRA)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+appSegmentPraEndpoint, nil, appSegmentPra, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(id string, appSegmentPra *AppSegmentPRA) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentPraEndpoint, id)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, appSegmentPra, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(id string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentPraEndpoint, id)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
