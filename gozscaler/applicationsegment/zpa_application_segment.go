package applicationsegment

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig         = "/mgmtconfig/v1/admin/customers/"
	appSegmentEndpoint = "/application"
)

type ApplicationSegmentResource struct {
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
	TCPPortRanges        []string              `json:"tcpPortRanges,omitempty"`
	UDPPortRanges        []string              `json:"udpPortRanges,omitempty"`
	TCPAppPortRange      []AppSegmentPortRange `json:"tcpPortRange,omitempty"`
	UDPAppPortRange      []AppSegmentPortRange `json:"udpPortRange,omitempty"`
	ClientlessApps       []ClientlessApps      `json:"clientlessApps,omitempty"`
	ServerGroups         []AppServerGroups     `json:"serverGroups,omitempty"`
	SRAApps              []SRAApps             `json:"sraApps,omitempty"`
	DefaultIdleTimeout   string                `json:"defaultIdleTimeout,omitempty"`
	DefaultMaxAge        string                `json:"defaultMaxAge,omitempty"`
}

type AppSegmentPortRange struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

type ClientlessApps struct {
	AllowOptions        bool   `json:"allowOptions"`
	AppID               string `json:"appId,omitempty"`
	ApplicationPort     string `json:"applicationPort,omitempty"`
	ApplicationProtocol string `json:"applicationProtocol,omitempty"`
	CertificateID       string `json:"certificateId,omitempty"`
	CertificateName     string `json:"certificateName,omitempty"`
	Cname               string `json:"cname,omitempty"`
	CreationTime        string `json:"creationTime,omitempty"`
	Description         string `json:"description,omitempty"`
	Domain              string `json:"domain,omitempty"`
	Enabled             bool   `json:"enabled"`
	Hidden              bool   `json:"hidden"`
	ID                  string `json:"id,omitempty"`
	LocalDomain         string `json:"localDomain,omitempty"`
	ModifiedBy          string `json:"modifiedBy,omitempty"`
	ModifiedTime        string `json:"modifiedTime,omitempty"`
	Name                string `json:"name"`
	Path                string `json:"path,omitempty"`
	Portal              bool   `json:"portal"`
	TrustUntrustedCert  bool   `json:"trustUntrustedCert"`
}

type AppServerGroups struct {
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

func (service *Service) Get(applicationID string) (*ApplicationSegmentResource, *http.Response, error) {
	v := new(ApplicationSegmentResource)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentEndpoint, applicationID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(appName string) (*ApplicationSegmentResource, *http.Response, error) {
	var v struct {
		List []ApplicationSegmentResource `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appSegmentEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize, Search: appName}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, app := range v.List {
		if strings.EqualFold(app.Name, appName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no application segment named '%s' was found", appName)
}

func (service *Service) Create(appSegment ApplicationSegmentResource) (*ApplicationSegmentResource, *http.Response, error) {
	v := new(ApplicationSegmentResource)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+appSegmentEndpoint, nil, appSegment, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) Update(applicationId string, appSegmentRequest ApplicationSegmentResource) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentEndpoint, applicationId)
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, nil, appSegmentRequest, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (service *Service) Delete(applicationId string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentEndpoint, applicationId)
	resp, err := service.Client.NewRequestDo("DELETE", relativeURL, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
