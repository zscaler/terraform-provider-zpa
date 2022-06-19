package variable

// App Connector Group
const (
	AppConnectorResourceName    = "testAcc_app_connector_group"
	AppConnectorDescription     = "testAcc_app_connector_group"
	AppConnectorEnabled         = true
	AppConnectorOverrideProfile = true
)

// Service Edge Group
const (
	ServiceEdgeResourceName = "testAcc_service_edge_group"
	ServiceEdgeDescription  = "testAcc_service_edge_group"
	ServiceEdgeEnabled      = true
)

// Provisioning Key
const (
	ConnectorGroupType     = "CONNECTOR_GRP"
	ServiceEdgeGroupType   = "SERVICE_EDGE_GRP"
	ProvisioningKeyUsage   = "2"
	ProvisioningKeyEnabled = true
)

// Customer Version Profile
const (
	VersionProfileDefault         = "Default"
	VersionProfilePreviousDefault = "Previous Default"
	VersionProfileNewRelease      = "New Release"
)

// Application Server
const (
	AppServerResourceName = "testAcc_application_server"
	AppServerDescription  = "testAcc_application_server"
	AppServerAddress      = "192.168.1.1"
	AppServerEnabled      = true
)

// Server Server
const (
	ServerGroupResourceName     = "testAcc_server_group"
	ServerGroupDescription      = "testAcc_server_group"
	ServerGroupEnabled          = true
	ServerGroupDynamicDiscovery = true
)

// Segment Group
const (
	SegmentGroupDescription = "testAcc_segment_group"
	SegmentGroupEnabled     = true
)

// Application Segment
const (
	AppSegmentResourceName = "testAcc_app_segment"
	AppSegmentDescription  = "testAcc_app_segment"
	AppSegmentEnabled      = true
	AppSegmentCnameEnabled = true
	AppSegmentGroupID      = ""
)

// Browser Access Segment
const (
	BrowserAccessEnabled      = true
	BrowserAccessCnameEnabled = true
)

// LSS Controller
const (
	LSSControllerEnabled    = true
	LSSControllerTLSEnabled = true
)

// Policy Access Rule
const (
	AccessRuleDescription = "testAcc_access_rule"
	AccessRuleAction      = "ALLOW"
	AccessRuleOrder       = 1
)

// Inspection Custom Control
const (
	InspectionCustomControlAction        = "PASS"
	InspectionCustomControlDefaultAction = "PASS"
	InspectionCustomControlSeverity      = "CRITICAL"
	InspectionCustomControlType          = "RESPONSE"
)
