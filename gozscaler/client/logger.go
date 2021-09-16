package client

import (
	"net/http"
	"net/http/httputil"
)

const (
	logReqMsg = `Request "%s %s" details:
---[ ZSCALER SDK REQUEST ]-------------------------------
%s
---------------------------------------------------------`

	logRespMsg = `Response "%s %s" details:
---[ ZSCALER SDK RESPONSE ]--------------------------------
%s
-------------------------------------------------------`
)

func (client *Client) WriteLog(format string, args ...interface{}) {
	if client.Config.Logger != nil {
		client.Config.Logger.Printf(format, args...)
	}
}

func (client *Client) logRequest(req *http.Request) {
	if client.Config.Logger != nil && req != nil {
		out, err := httputil.DumpRequestOut(req, true)
		if err == nil {
			client.WriteLog(logReqMsg, req.Method, req.URL, string(out))
		}
	}
}

func (client *Client) logResponse(resp *http.Response) {
	if client.Config.Logger != nil && resp != nil {
		out, err := httputil.DumpResponse(resp, true)
		if err == nil {
			client.WriteLog(logRespMsg, resp.Request.Method, resp.Request.URL, string(out))
		}
	}
}
