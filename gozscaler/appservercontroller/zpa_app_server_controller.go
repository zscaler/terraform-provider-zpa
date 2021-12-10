package appservercontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/common"
)

const (
	mgmtConfig                  = "/mgmtconfig/v1/admin/customers/"
	appServerControllerEndpoint = "/server"
)

type ApplicationServer struct {
	Address           string   `json:"address"`
	AppServerGroupIds []string `json:"appServerGroupIds"`
	ConfigSpace       string   `json:"configSpace,omitempty"`
	CreationTime      string   `json:"creationTime,"`
	Description       string   `json:"description"`
	Enabled           bool     `json:"enabled"`
	ID                string   `json:"id,"`
	ModifiedBy        string   `json:"modifiedBy,"`
	ModifiedTime      string   `json:"modifiedTime,"`
	Name              string   `json:"name"`
}

func (service *Service) Get(id string) (*ApplicationServer, *http.Response, error) {
	v := new(ApplicationServer)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appServerControllerEndpoint, id)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(appServerName string) (*ApplicationServer, *http.Response, error) {
	var v struct {
		List []ApplicationServer `json:"list"`
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appServerControllerEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Pagination{PageSize: common.DefaultPageSize}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	for _, app := range v.List {
		if strings.EqualFold(app.Name, appServerName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no application named '%s' was found", appServerName)
}

func (service *Service) Create(server ApplicationServer) (*ApplicationServer, *http.Response, error) {
	v := new(ApplicationServer)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+appServerControllerEndpoint, nil, server, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(id string, appServer ApplicationServer) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appServerControllerEndpoint, id)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, appServer, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(id string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appServerControllerEndpoint, id)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
