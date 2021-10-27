package lssconfigcontroller

import (
	"fmt"
	"net/http"
)

const (
	lssClientTypesEndpoint = "/lssConfig/clientTypes"
)

type LSSClientTypes struct {
	//ID                         string `json:"id"`
	ZPNClientTypeExporter      string `json:"zpn_client_type_exporter"`
	ZPNClientTypeMachineTunnel string `json:"zpn_client_type_machine_tunnel"`
	ZPNClientTypeIPAnchoring   string `json:"zpn_client_type_ip_anchoring"`
	ZPNClientTypeEdgeConnector string `json:"zpn_client_type_edge_connector"`
	ZPNClientTypeZAPP          string `json:"zpn_client_type_zapp"`
	ZPNClientTypeSlogger       string `json:"zpn_client_type_slogger,omitempty"`
}

func (service *Service) GetClientTypes() (*LSSClientTypes, *http.Response, error) {
	v := new(LSSClientTypes)
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + lssClientTypesEndpoint)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
