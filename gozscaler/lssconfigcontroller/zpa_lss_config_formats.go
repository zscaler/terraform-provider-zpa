package lssconfigcontroller

import (
	"fmt"
	"net/http"
)

const (
	lssFormatsEndpoint = "/lssConfig/statusCodes"
)

type LSSFormats struct {
	ZPNAuthLog      map[string]interface{} `json:"zpn_auth_log"`
	ZPNAstAuthLog   map[string]interface{} `json:"zpn_ast_auth_log"`
	ZPNAuditLog     map[string]interface{} `json:"zpn_audit_log"`
	ZPNTransLog     map[string]interface{} `json:"zpn_trans_log"`
	ZPNHTTPTransLog map[string]interface{} `json:"zpn_http_trans_log"`
}

func (service *Service) GetFormats() (*LSSFormats, *http.Response, error) {
	v := new(LSSFormats)
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + lssFormatsEndpoint)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
