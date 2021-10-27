package lssconfigcontroller

import (
	"fmt"
	"net/http"
)

const (
	lssStatusCodesEndpoint = "/lssConfig/statusCodes"
)

type LSSStatusCodes struct {
	ZPNAuthLog    map[string]interface{} `json:"zpn_auth_log"`
	ZPNAstAuthLog map[string]interface{} `json:"zpn_ast_auth_log"`
	ZPNTransLog   map[string]interface{} `json:"zpn_trans_log"`
}

func (service *Service) GetStatusCodes() (*LSSStatusCodes, *http.Response, error) {
	v := new(LSSStatusCodes)
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + lssStatusCodesEndpoint)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
