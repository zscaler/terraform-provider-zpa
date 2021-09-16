package common

type ServerGroup struct {
	ID                 string              `json:"id"`
	Enabled            bool                `json:"enabled"`
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	IpAnchored         bool                `json:"ipAnchored"`
	ConfigSpace        string              `json:"configSpace"`
	DynamicDiscovery   bool                `json:"dynamicDiscovery"`
	CreationTime       int32               `json:"creationTime,string"`
	ModifiedBy         string              `json:"modifiedBy"`
	ModifiedTime       int32               `json:"modifiedTime,string"`
	AppConnectorGroups []AppConnectorGroup `json:"appConnectorGroups,omitempty"`
	ApplicationServers []ApplicationServer `json:"servers,omitempty"`
	Applications       []Application       `json:"applications,omitempty"`
}
type Application struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
}

type AppConnectorGroup struct {
	CityCountry           string        `json:"cityCountry"`
	CountryCode           string        `json:"countryCode"`
	CreationTime          int           `json:"creationTime,string"`
	Description           string        `json:"description"`
	DnsQueryType          string        `json:"dnsQueryType"`
	Enabled               bool          `json:"enabled"`
	GeolocationId         int64         `json:"geoLocationId,string"`
	ID                    int64         `json:"id,string"`
	Latitude              string        `json:"latitude"`
	Location              string        `json:"location"`
	Longitude             string        `json:"longitude"`
	ModifiedBy            int64         `json:"modifiedBy,string"`
	ModifiedTime          int32         `json:"modifiedTime,string"`
	Name                  string        `json:"name"`
	SiemAppConnectorGroup bool          `json:"siemAppConnectorGroup"`
	UpgradeDay            string        `json:"upgradeDay"`
	UpgradeTimeInSecs     string        `json:"upgradeTimeInSecs"`
	VersionProfileId      int64         `json:"versionProfileId,string"`
	ServerGroups          []ServerGroup `json:"serverGroups"`
	Connectors            []Connector   `json:"connectors"`
}

type Connector struct {
	ApplicationStartTime     int64    `json:"applicationStartTime,string"`
	AppConnectorGroupID      string   `json:"appConnectorGroupId,omitempty"`
	AppConnectorGroupName    string   `json:"appConnectorGroupName,omitempty"`
	ControlChannelStatus     string   `json:"controlChannelStatus,omitempty"`
	CreationTime             int32    `json:"creationTime,string"`
	CtrlBrokerName           string   `json:"ctrlBrokerName,omitempty"`
	CurrentVersion           string   `json:"currentVersion,omitempty"`
	Description              string   `json:"description"`
	Enabled                  bool     `json:"enabled,omitempty"`
	ExpectedUpgradeTime      int64    `json:"expectedUpgradeTime,string"`
	ExpectedVersion          string   `json:"expectedVersion,omitempty"`
	Fingerprint              string   `json:"fingerprint,omitempty"`
	ID                       int      `json:"id,string"`
	IpAcl                    []string `json:"ipAcl,omitempty"`
	IssuedCertID             int64    `json:"issuedCertId,string,omitempty"`
	LastBrokerConnectTime    int64    `json:"lastBrokerConnectTime,string,omitempty"`
	LastBrokerDisconnectTime int64    `json:"lastBrokerDisconnectTime,string,omitempty"`
	LastUpgradeTime          int64    `json:"lastUpgradeTime,string,omitempty"`
	Latitude                 int      `json:"latitude,string,omitempty"`
	Location                 string   `json:"location,omitempty"`
	Longitude                int      `json:"longitude,string,omitempty"`
	ModifiedBy               int64    `json:"modifiedBy,string,omitempty"`
	ModifiedTime             int32    `json:"modifiedTime,string,omitempty"`
	Name                     string   `json:"name,omitempty"`
	Platform                 string   `json:"platform"`
	PreviousVersion          string   `json:"previousVersion"`
	PrivateIp                string   `json:"privateIp"`
	PublicIp                 string   `json:"publicIp"`
	UpgradeAttempt           int32    `json:"upgradeAttempt,string"`
	UpgradeStatus            string   `json:"upgradeStatus"`
}

type ApplicationServer struct {
	Address           string   `json:"address"`
	AppServerGroupIds []string `json:"appServerGroupIds"` // Don't omitempty. We need empty slice in JSON for update.
	ConfigSpace       string   `json:"configSpace,omitempty"`
	CreationTime      int32    `json:"creationTime,string"`
	Description       string   `json:"description"`
	Enabled           bool     `json:"enabled"`
	ID                int64    `json:"id,string"`
	ModifiedBy        int64    `json:"modifiedBy,string"`
	ModifiedTime      int32    `json:"modifiedTime,string"`
	Name              string   `json:"name"`
}
