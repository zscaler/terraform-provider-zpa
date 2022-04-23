package client

/*
import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zscaler/terraform-provider-zpa/gozscaler"
)

type dummyStruct struct {
	ID int `json:"id"`
}

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

const getResponse = `{"id": 1234}`

func TestClient_NewRequestDo(t *testing.T) {

	type args struct {
		method string
		url    string
		body   interface{}
		v      interface{}
	}
	tests := []struct {
		name       string
		args       args
		muxHandler func(w http.ResponseWriter, r *http.Request)
		wantResp   *http.Response
		wantErr    bool
		wantVal    *dummyStruct
	}{
		// NewRequestDo test cases
		{
			name: "GET happy path",
			args: struct {
				method string
				url    string
				body   interface{}
				v      interface{}
			}{
				method: "GET",
				url:    "/test",
				body:   nil,
				v:      new(dummyStruct),
			},
			muxHandler: func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte(getResponse))
				w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
				if err != nil {
					t.Fatal(err)
				}
			},
			wantResp: &http.Response{
				StatusCode: 200,
			},
			wantVal: &dummyStruct{
				ID: 1234,
			},
		},
	}

	for _, tt := range tests {
		client = NewClient(setupMuxConfig())
		client.WriteLog("Server URL: %v", client.Config.BaseURL)
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc(tt.args.url, tt.muxHandler)
			res, err := client.NewRequestDo(tt.args.method, tt.args.url, nil, tt.args.body, tt.args.v)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.NewRequestDo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantResp.StatusCode != res.StatusCode {
				t.Errorf("Client.NewRequestDo() = %v, want %v", res, tt.wantResp)
			}

			if !reflect.DeepEqual(tt.args.v, tt.wantVal) {
				t.Errorf("returned %#v; want %#v", tt.args.v, tt.wantVal)
			}
		})
	}
	teardown()
}

func TestNewClient(t *testing.T) {
	type args struct {
		config *gozscaler.Config
	}
	tests := []struct {
		name  string
		args  args
		wantC *Client
	}{
		// NewClient test cases
		{
			name: "Successful Client creation with default config values",
			args: struct{ config *gozscaler.Config }{config: nil},
			wantC: &Client{&gozscaler.Config{
				BaseURL: &url.URL{
					Scheme: "https",
					Host:   "config.private.zscaler.com",
					Path:   "/signin",
				},
			}},
		},
		{
			name: "Successful Client creation with custom config values",
			args: struct{ config *gozscaler.Config }{config: &gozscaler.Config{
				BaseURL: &url.URL{Host: "https://otherhost.com"},
			}},
			wantC: &Client{&gozscaler.Config{
				BaseURL: &url.URL{Host: "https://otherhost.com"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC := NewClient(tt.args.config)
			assert.Equal(t, gotC.Config.BaseURL.Host, tt.wantC.Config.BaseURL.Host)
			assert.Equal(t, gotC.Config.BaseURL.Scheme, tt.wantC.Config.BaseURL.Scheme)
			assert.Equal(t, gotC.Config.BaseURL.Path, tt.wantC.Config.BaseURL.Path)
		})
	}
}

func setupMuxConfig() *gozscaler.Config {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	config, err := gozscaler.NewConfig("", "", "", server.URL)
	if err != nil {
		panic(err)
	}
	return config
}

func teardown() {
	server.Close()
}
*/
