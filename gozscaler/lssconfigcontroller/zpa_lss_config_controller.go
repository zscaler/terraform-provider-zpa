package lssconfigcontroller

type LSSResource struct {
	Config             Config             `json:"config"`
	ID                 int                `json:"id"`
	ConnectorGroups    []ConnectorGroups  `json:"connectorGroups"`
	PolicyRule         PolicyRule         `json:"policyRule"`
	PolicyRuleResource PolicyRuleResource `json:"policyRuleResource"`
}
type Config struct {
	AuditMessage  string   `json:"auditMessage"`
	CreationTime  int      `json:"creationTime"`
	Description   string   `json:"description"`
	Enabled       bool     `json:"enabled"`
	Filter        []string `json:"filter"`
	Format        string   `json:"format"`
	ID            int      `json:"id"`
	ModifiedBy    int      `json:"modifiedBy"`
	ModifiedTime  int      `json:"modifiedTime"`
	Name          string   `json:"name"`
	LssHost       string   `json:"lssHost"`
	LssPort       int      `json:"lssPort"`
	SourceLogType string   `json:"sourceLogType"`
	UseTLS        bool     `json:"useTls"`
}

type ConnectorGroups struct {
	CityCountry                   string         `json:"cityCountry"`
	CountryCode                   string         `json:"countryCode"`
	CreationTime                  int            `json:"creationTime"`
	Description                   string         `json:"description"`
	DNSQueryType                  string         `json:"dnsQueryType"`
	Enabled                       bool           `json:"enabled"`
	GeoLocationID                 int            `json:"geoLocationId"`
	ID                            int            `json:"id"`
	Latitude                      string         `json:"latitude"`
	Location                      string         `json:"location"`
	Longitude                     string         `json:"longitude"`
	ModifiedBy                    int            `json:"modifiedBy"`
	ModifiedTime                  int            `json:"modifiedTime"`
	Name                          string         `json:"name"`
	OverrideVersionProfile        bool           `json:"overrideVersionProfile"`
	LssAppConnectorGroup          bool           `json:"lssAppConnectorGroup"`
	UpgradeDay                    string         `json:"upgradeDay"`
	UpgradeTimeInSecs             string         `json:"upgradeTimeInSecs"`
	VersionProfileID              int            `json:"versionProfileId"`
	VersionProfileName            string         `json:"versionProfileName"`
	VersionProfileVisibilityScope string         `json:"versionProfileVisibilityScope"`
	ServerGroups                  []ServerGroups `json:"serverGroups"`
	// Connectors                    []Connectors   `json:"connectors"`
}

/*
type Connectors struct {
	ApplicationStartTime             string                 `json:"applicationStartTime"`
	AppConnectorGroupID              string                 `json:"appConnectorGroupId"`
	AppConnectorGroupName            string                 `json:"appConnectorGroupName"`
	ControlChannelStatus             string                 `json:"controlChannelStatus"`
	CreationTime                     string                 `json:"creationTime"`
	CtrlBrokerName                   string                 `json:"ctrlBrokerName"`
	CurrentVersion                   string                 `json:"currentVersion"`
	Description                      string                 `json:"description"`
	Enabled                          bool                   `json:"enabled"`
	ExpectedUpgradeTime              string                 `json:"expectedUpgradeTime"`
	ExpectedVersion                  string                 `json:"expectedVersion"`
	Fingerprint                      string                 `json:"fingerprint"`
	ID                               string                 `json:"id"`
	IPACL                            []string               `json:"ipAcl"`
	IssuedCertID                     string                 `json:"issuedCertId"`
	LastBrokerConnectTime            string                 `json:"lastBrokerConnectTime"`
	LastBrokerConnectTimeDuration    string                 `json:"lastBrokerConnectTimeDuration"`
	LastBrokerDisconnectTime         string                 `json:"lastBrokerDisconnectTime"`
	LastBrokerDisconnectTimeDuration string                 `json:"lastBrokerDisconnectTimeDuration"`
	LastUpgradeTime                  string                 `json:"lastUpgradeTime"`
	Latitude                         string                 `json:"latitude"`
	Location                         string                 `json:"location"`
	Longitude                        string                 `json:"longitude"`
	ModifiedBy                       string                 `json:"modifiedBy"`
	ModifiedTime                     string                 `json:"modifiedTime"`
	Name                             string                 `json:"name"`
	ProvisioningKeyID                string                 `json:"provisioningKeyId"`
	ProvisioningKeyName              string                 `json:"provisioningKeyName"`
	Platform                         string                 `json:"platform"`
	PreviousVersion                  string                 `json:"previousVersion"`
	PrivateIP                        string                 `json:"privateIp"`
	PublicIP                         string                 `json:"publicIp"`
	SargeVersion                     string                 `json:"sargeVersion"`
	EnrollmentCert                   map[string]interface{} `json:"enrollmentCert"`
	UpgradeAttempt                   string                 `json:"upgradeAttempt"`
	UpgradeStatus                    string                 `json:"upgradeStatus"`
}

type ServerGroups struct {
	ConfigSpace      string `json:"configSpace"`
	CreationTime     string `json:"creationTime"`
	Description      string `json:"description"`
	Enabled          bool   `json:"enabled"`
	ID               string `json:"id"`
	DynamicDiscovery bool   `json:"dynamicDiscovery"`
	ModifiedBy       string `json:"modifiedBy"`
	ModifiedTime     string `json:"modifiedTime"`
	Name             string `json:"name"`
}
*/

type PolicyRule struct {
	Action                       string                         `json:"action"`
	ActionID                     string                         `json:"actionId"`
	BypassDefaultRule            bool                           `json:"bypassDefaultRule"`
	CreationTime                 string                         `json:"creationTime"`
	CustomMsg                    string                         `json:"customMsg"`
	DefaultRule                  bool                           `json:"defaultRule"`
	Description                  string                         `json:"description"`
	ID                           string                         `json:"id"`
	IsolationDefaultRule         bool                           `json:"isolationDefaultRule"`
	ModifiedBy                   string                         `json:"modifiedBy"`
	ModifiedTime                 string                         `json:"modifiedTime"`
	Name                         string                         `json:"name"`
	Operator                     string                         `json:"operator"`
	PolicySetID                  string                         `json:"policySetId"`
	PolicyType                   string                         `json:"policyType"`
	Priority                     string                         `json:"priority"`
	ReauthDefaultRule            bool                           `json:"reauthDefaultRule"`
	ReauthIdleTimeout            string                         `json:"reauthIdleTimeout"`
	ReauthTimeout                string                         `json:"reauthTimeout"`
	RuleOrder                    string                         `json:"ruleOrder"`
	LssDefaultRule               bool                           `json:"lssDefaultRule"`
	ZpnCbiProfileID              string                         `json:"zpnCbiProfileId"`
	ZpnInspectionProfileID       string                         `json:"zpnInspectionProfileId"`
	ZpnInspectionProfileName     string                         `json:"zpnInspectionProfileName"`
	PolicyRuleAppServerGroups    []PolicyRuleAppServerGroups    `json:"appServerGroups"`
	PolicyRuleAppConnectorGroups []PolicyRuleAppConnectorGroups `json:"appConnectorGroups"`
	Conditions                   []Conditions                   `json:"conditions"`
}

type PolicyRuleAppServerGroups struct {
	ConfigSpace      string `json:"configSpace"`
	CreationTime     int    `json:"creationTime"`
	Description      string `json:"description"`
	Enabled          bool   `json:"enabled"`
	ID               int    `json:"id"`
	DynamicDiscovery bool   `json:"dynamicDiscovery"`
	ModifiedBy       int    `json:"modifiedBy"`
	ModifiedTime     int    `json:"modifiedTime"`
	Name             string `json:"name"`
}
type PolicyRuleAppConnectorGroups struct {
	CityCountry                   string         `json:"cityCountry"`
	CountryCode                   string         `json:"countryCode"`
	CreationTime                  int            `json:"creationTime"`
	Description                   string         `json:"description"`
	DNSQueryType                  string         `json:"dnsQueryType"`
	Enabled                       bool           `json:"enabled"`
	GeoLocationID                 int            `json:"geoLocationId"`
	ID                            int            `json:"id"`
	Latitude                      string         `json:"latitude"`
	Location                      string         `json:"location"`
	Longitude                     string         `json:"longitude"`
	ModifiedBy                    int            `json:"modifiedBy"`
	ModifiedTime                  int            `json:"modifiedTime"`
	Name                          string         `json:"name"`
	OverrideVersionProfile        bool           `json:"overrideVersionProfile"`
	LssAppConnectorGroup          bool           `json:"lssAppConnectorGroup"`
	UpgradeDay                    string         `json:"upgradeDay"`
	UpgradeTimeInSecs             string         `json:"upgradeTimeInSecs"`
	VersionProfileID              int            `json:"versionProfileId"`
	VersionProfileName            string         `json:"versionProfileName"`
	VersionProfileVisibilityScope string         `json:"versionProfileVisibilityScope"`
	ServerGroups                  []ServerGroups `json:"serverGroups"`
	Connectors                    []Connectors   `json:"connectors"`
}
type PolicyRuleOperands struct {
	CreationTime string `json:"creationTime"`
	ID           string `json:"id"`
	IdpID        string `json:"idpId"`
	LHS          string `json:"lhs"`
	ModifiedBy   string `json:"modifiedBy"`
	ModifiedTime string `json:"modifiedTime"`
	Name         string `json:"name"`
	ObjectType   string `json:"objectType"`
	RHS          string `json:"rhs"`
}
type PolicyRuleConditions struct {
	CreationTime string     `json:"creationTime"`
	ID           string     `json:"id"`
	ModifiedBy   string     `json:"modifiedBy"`
	ModifiedTime string     `json:"modifiedTime"`
	Negated      bool       `json:"negated"`
	Operands     []Operands `json:"operands"`
	Operator     string     `json:"operator"`
}

type AppServerGroups struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type AppConnectorGroups struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type CommonProperties struct {
	CreationTime string `json:"creationTime"`
	ID           string `json:"id"`
	IdpID        string `json:"idpId"`
	LHS          string `json:"lhs"`
	ModifiedBy   string `json:"modifiedBy"`
	ModifiedTime string `json:"modifiedTime"`
	Name         string `json:"name"`
	ObjectType   string `json:"objectType"`
	RHS          string `json:"rhs"`
}
type EntryValues struct {
	LHS string `json:"lhs"`
	RHS string `json:"rhs"`
}
type Operands struct {
	CommonProperties CommonProperties `json:"commonProperties"`
	CreationTime     int              `json:"creationTime"`
	EntryValues      []EntryValues    `json:"entryValues"`
	ID               string           `json:"id"`
	IdpID            string           `json:"idpId"`
	ModifiedBy       string           `json:"modifiedBy"`
	ModifiedTime     string           `json:"modifiedTime"`
	ObjectType       string           `json:"objectType"`
	Values           []string         `json:"values"`
}
type Conditions struct {
	CreationTime string     `json:"creationTime"`
	ID           string     `json:"id"`
	ModifiedBy   string     `json:"modifiedBy"`
	ModifiedTime string     `json:"modifiedTime"`
	Negated      bool       `json:"negated"`
	Operands     []Operands `json:"operands"`
	Operator     string     `json:"operator"`
	SetIds       []int      `json:"setIds"`
}
type PolicyRuleResource struct {
	Action                   string               `json:"action"`
	ActionID                 string               `json:"actionId"`
	AppServerGroups          []AppServerGroups    `json:"appServerGroups"`
	AppConnectorGroups       []AppConnectorGroups `json:"appConnectorGroups"`
	AuditMessage             string               `json:"auditMessage"`
	Conditions               []Conditions         `json:"conditions"`
	CreationTime             string               `json:"creationTime"`
	CustomMsg                string               `json:"customMsg"`
	Description              string               `json:"description"`
	ID                       string               `json:"id"`
	ModifiedBy               string               `json:"modifiedBy"`
	ModifiedTime             string               `json:"modifiedTime"`
	Name                     string               `json:"name"`
	Operator                 string               `json:"operator"`
	PolicySetID              string               `json:"policySetId"`
	PolicyType               string               `json:"policyType"`
	Priority                 string               `json:"priority"`
	ReauthIdleTimeout        string               `json:"reauthIdleTimeout"`
	ReauthTimeout            string               `json:"reauthTimeout"`
	RuleOrder                string               `json:"ruleOrder"`
	Version                  string               `json:"version"`
	ZpnCbiProfileID          string               `json:"zpnCbiProfileId"`
	ZpnInspectionProfileID   string               `json:"zpnInspectionProfileId"`
	ZpnInspectionProfileName string               `json:"zpnInspectionProfileName"`
}
