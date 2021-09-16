package segmentgroup

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	mgmtConfig           = "/mgmtconfig/v1/admin/customers/"
	segmentGroupEndpoint = "/segmentGroup"
)

type SegmentGroup struct {
	Applications        []Application `json:"applications"`
	ConfigSpace         string        `json:"configSpace,omitempty"`
	CreationTime        string        `json:"creationTime,omitempty"`
	Description         string        `json:"description,omitempty"`
	Enabled             bool          `json:"enabled"`
	ID                  string        `json:"id,omitempty"`
	ModifiedBy          string        `json:"modifiedBy,omitempty"`
	ModifiedTime        string        `json:"modifiedTime,omitempty"`
	Name                string        `json:"name"`
	PolicyMigrated      bool          `json:"policyMigrated"`
	TcpKeepAliveEnabled string        `json:"tcpKeepAliveEnabled,omitempty"`
}

type Application struct {
	BypassType           string           `json:"bypassType,omitempty"`
	ConfigSpace          string           `json:"configSpace,omitempty"`
	CreationTime         string           `json:"creationTime,omitempty"`
	DefaultIdleTimeout   string           `json:"defaultIdleTimeout,omitempty"`
	DefaultMaxAge        string           `json:"defaultMaxAge,omitempty"`
	Description          string           `json:"description,omitempty"`
	DomainName           string           `json:"domainName,omitempty"`
	DomainNames          []string         `json:"domainNames,omitempty"`
	DoubleEncrypt        bool             `json:"doubleEncrypt"`
	Enabled              bool             `json:"enabled"`
	HealthCheckType      string           `json:"healthCheckType,omitempty"`
	ID                   string           `json:"id,omitempty"`
	IPAnchored           bool             `json:"ipAnchored"`
	LogFeatures          []string         `json:"logFeatures,omitempty"`
	ModifiedBy           string           `json:"modifiedBy,omitempty"`
	ModifiedTime         string           `json:"modifiedTime,omitempty"`
	Name                 string           `json:"name"`
	PassiveHealthEnabled bool             `json:"passiveHealthEnabled"`
	ServerGroup          []AppServerGroup `json:"serverGroups,omitempty"`
	TCPPortRanges        interface{}      `json:"tcpPortRanges,omitempty"`
	TCPPortsIn           interface{}      `json:"tcpPortsIn,omitempty"`
	TCPPortsOut          interface{}      `json:"tcpPortsOut,omitempty"`
	UDPPortRanges        interface{}      `json:"udpPortRangesg,omitempty"`
}
type AppServerGroup struct {
	ConfigSpace      string `json:"configSpace,omitempty"`
	CreationTime     string `json:"creationTime,omitempty"`
	Description      string `json:"description,omitempty"`
	Enabled          bool   `json:"enabled"`
	ID               string `json:"id,omitempty"`
	DynamicDiscovery bool   `json:"dynamicDiscovery"`
	ModifiedBy       string `json:"modifiedBy,omitempty"`
	ModifiedTime     string `json:"modifiedTime,omitempty"`
	Name             string `json:"name"`
}

func (service *Service) Get(segmentGroupId string) (*SegmentGroup, *http.Response, error) {
	v := new(SegmentGroup)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+segmentGroupEndpoint, segmentGroupId)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(segmentName string) (*SegmentGroup, *http.Response, error) {
	var v struct {
		List []SegmentGroup `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + segmentGroupEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, struct{ pagesize int }{
		pagesize: 500,
	}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, app := range v.List {
		if strings.EqualFold(app.Name, segmentName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no application named '%s' was found", segmentName)
}

func (service *Service) Create(segmentGroup *SegmentGroup) (*SegmentGroup, *http.Response, error) {
	v := new(SegmentGroup)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+segmentGroupEndpoint, nil, segmentGroup, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(segmentGroupId string, segmentGroupRequest *SegmentGroup) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+segmentGroupEndpoint, segmentGroupId)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, segmentGroupRequest, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(segmentGroupId string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+segmentGroupEndpoint, segmentGroupId)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
