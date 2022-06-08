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
	AppsConfig         []AppsConfig `json:"appsConfig,omitempty"`
	DeletedInspectApps []string     `json:"deletedInspectApps,omitempty"`
}

type AppsConfig struct {
	Name                string   `json:"name,omitempty"`
	AllowOptions        bool     `json:"allowOptions"`
	ID                  string   `json:"id,omitempty"`
	AppID               string   `json:"appId,omitempty"`
	AppTypes            []string `json:"appTypes,omitempty"`
	ApplicationPort     string   `json:"applicationPort,omitempty"`
	ApplicationProtocol string   `json:"applicationProtocol,omitempty"`
	InspectAppID        string   `json:"inspectAppId,omitempty"`
	CertificateID       string   `json:"certificateId,omitempty"`
	CertificateName     string   `json:"certificateName,omitempty"`
	Cname               string   `json:"cname,omitempty"`
	Description         string   `json:"description,omitempty"`
	Domain              string   `json:"domain,omitempty"`
	Enabled             bool     `json:"enabled"`
	Hidden              bool     `json:"hidden"`
	LocalDomain         string   `json:"localDomain,omitempty"`
	Portal              bool     `json:"portal"`
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

func (service *Service) GetByName(appSegmentName string) (*AppSegmentInspection, *http.Response, error) {
	var v struct {
		List []AppSegmentInspection `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appSegmentInspectionEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize, Search: appSegmentName}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, app := range v.List {
		if strings.EqualFold(app.Name, appSegmentName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no inspection application segment named '%s' was found", appSegmentName)
}

func (service *Service) Create(appSegmentInspection AppSegmentInspection) (*AppSegmentInspection, *http.Response, error) {
	v := new(AppSegmentInspection)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+appSegmentInspectionEndpoint, nil, appSegmentInspection, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// return the new items that were added to slice1
func difference(slice1 []AppsConfig, slice2 []AppsConfig) []AppsConfig {
	var diff []AppsConfig
	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1.Domain == s2.Domain || s1.Name == s2.Name {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, s1)
		}
	}
	return diff
}

func mapInspectionApp(InspectionAppDto []InspectionAppDto) []AppsConfig {
	result := []AppsConfig{}
	for _, app := range InspectionAppDto {
		result = append(result, AppsConfig{
			Name:   app.Name,
			Domain: app.Domain,
			ID:     app.ID,
			AppID:  app.AppID,
		})
	}
	return result
}

func appToListStringIDs(apps []AppsConfig) []string {
	result := []string{}
	for _, app := range apps {
		result = append(result, app.ID)
	}
	return result
}

func (service *Service) Update(id string, appSegmentInspection *AppSegmentInspection) (*http.Response, error) {
	existingResource, _, err := service.Get(id)
	if err != nil {
		return nil, err
	}
	existingApps := mapInspectionApp(existingResource.InspectionAppDto)
	newApps := difference(appSegmentInspection.CommonAppsDto.AppsConfig, existingApps)
	removedApps := difference(existingApps, appSegmentInspection.CommonAppsDto.AppsConfig)
	appSegmentInspection.CommonAppsDto.AppsConfig = newApps
	appSegmentInspection.CommonAppsDto.DeletedInspectApps = appToListStringIDs(removedApps)
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentInspectionEndpoint, id)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, appSegmentInspection, nil)
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
