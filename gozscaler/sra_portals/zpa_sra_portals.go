package sra_portals

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig            = "/mgmtconfig/v1/admin/customers/"
	sraAppSegmentEndpoint = "/application"
)

type ApplicationSegmentSRA struct {
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
	SRAApps              []SRAApps             `json:"sraApps,omitempty"`
	AppServerGroups      []AppServerGroups     `json:"serverGroups,omitempty"`
}

type SRAApps struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	Enabled             bool   `json:"enabled"`
	ApplicationPort     string `json:"applicationPort,omitempty"`
	ApplicationProtocol string `json:"applicationProtocol,omitempty"`
	Domain              string `json:"domain,omitempty"`
	AppID               string `json:"appId,omitempty"`
	Hidden              bool   `json:"hidden,omitempty"`
	Portal              bool   `json:"portal,omitempty"`
	ConnectionSecurity  string `json:"connectionSecurity,omitempty"`
}

type NetworkPorts struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

type AppServerGroups struct {
	ID string `json:"id"`
}

func (service *Service) Get(sraApplicationId string) (*ApplicationSegmentSRA, *http.Response, error) {
	v := new(ApplicationSegmentSRA)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+sraAppSegmentEndpoint, sraApplicationId)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(sraName string) (*ApplicationSegmentSRA, *http.Response, error) {
	var v struct {
		List []ApplicationSegmentSRA `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + sraAppSegmentEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize, Search: sraName}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, sra := range v.List {
		if strings.EqualFold(sra.Name, sraName) {
			return &sra, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no sra application segment named '%s' was found", sraName)
}

func (service *Service) Create(appSegment ApplicationSegmentSRA) (*ApplicationSegmentSRA, *http.Response, error) {
	v := new(ApplicationSegmentSRA)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+sraAppSegmentEndpoint, nil, appSegment, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) Update(id string, browserAccess *ApplicationSegmentSRA) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+sraAppSegmentEndpoint, id)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, browserAccess, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(sraApplicationId string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+sraAppSegmentEndpoint, sraApplicationId)
	resp, err := service.Client.NewRequestDo("DELETE", relativeURL, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
