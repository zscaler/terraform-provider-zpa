package lssconfigcontroller

import (
	"fmt"
	"net/http"
)

type LSSFormats struct {
	Csv  string `json:"csv"`
	Tsv  string `json:"tsv"`
	Json string `json:"json"`
}

func (service *Service) GetFormats(logType string) (*LSSFormats, *http.Response, error) {
	v := new(LSSFormats)
	relativeURL := fmt.Sprintf("%slssConfig/logType/formats", mgmtConfigTypesAndFormats)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, struct {
		LogType string `url:"logType"`
	}{
		LogType: logType,
	}, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
