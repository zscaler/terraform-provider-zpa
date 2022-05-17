package postureprofile

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig             = "/mgmtconfig/v2/admin/customers/"
	postureProfileEndpoint = "/posture"
)

// PostureProfile ...
type PostureProfile struct {
	CreationTime      string `json:"creationTime,omitempty"`
	Domain            string `json:"domain,omitempty"`
	ID                string `json:"id,omitempty"`
	MasterCustomerID  string `json:"masterCustomerId,omitempty"`
	ModifiedBy        string `json:"modifiedBy,omitempty"`
	ModifiedTime      string `json:"modifiedTime,omitempty"`
	Name              string `json:"name,omitempty"`
	PostureudID       string `json:"postureUdid,omitempty"`
	ZscalerCloud      string `json:"zscalerCloud,omitempty"`
	ZscalerCustomerID string `json:"zscalerCustomerId,omitempty"`
}

func (service *Service) Get(id string) (*PostureProfile, *http.Response, error) {
	v := new(PostureProfile)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+postureProfileEndpoint, id)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByPostureUDID(postureUDID string) (*PostureProfile, *http.Response, error) {
	var v struct {
		List []PostureProfile `json:"list"`
	}
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + postureProfileEndpoint)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, postureProfile := range v.List {
		if postureProfile.PostureudID == postureUDID {
			return &postureProfile, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no posture profile with postureUDID '%s' was found", postureUDID)
}

func (service *Service) GetByName(postureName string) (*PostureProfile, *http.Response, error) {
	var v struct {
		List []PostureProfile `json:"list"`
	}
	// remove appended cloud name (zscalerthree.net) i.e "CrowdStrike_ZPA_Pre-ZTA (zscalerthree.net)"
	adaptedPostureName := postureName[:strings.LastIndex(postureName, " ")]
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + postureProfileEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize, Search2: adaptedPostureName}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, postureProfile := range v.List {
		if strings.EqualFold(postureProfile.Name, postureName) {
			return &postureProfile, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no posture profile named '%s' was found", postureName)
}
