---
subcategory: "Application Segment"
layout: "zscaler"
page_title: "ZPA: application_segment"
description: |-
  Creates and manages ZPA Application Segments.
---

# zpa_application_segment

The **zpa_application_segment** resource creates an application segment in the Zscaler Private Access cloud. This resource can then be referenced in an access policy rule, access policy timeout rule or access policy client forwarding rule.

## Example 1 Usage

```hcl
# ZPA Application Segment resource
resource "zpa_application_segment" "test_app_segment" {
    name              = "test1-app-segment"
    description       = "test1-app-segment"
    enabled           = true
    health_reporting  = "ON_ACCESS"
    bypass_type       = "NEVER"
    is_cname_enabled  = true
    tcp_port_ranges   = ["8080", "8080"]
    domain_names      = ["server.acme.com"]
    segment_group_id  = zpa_segment_group.test_segment_group.id
    server_groups {
        id = [ zpa_server_group.test_server_group.id]
    }
    depends_on = [ zpa_server_group.test_server_group, zpa_server_group.test_segment_group]
}

# ZPA Segment Group resource
resource "zpa_segment_group" "test_segment_group" {
  name            = "test1-segment-group"
  description     = "test1-segment-group"
  enabled         = true
}

# ZPA Server Group resource
resource "zpa_server_group" "test_server_group" {
  name              = "test1-server-group"
  description       = "test1-server-group"
  enabled           = true
  dynamic_discovery = false
  app_connector_groups {
    id = [ zpa_app_connector_group.example.id ]
  }
  depends_on = [ zpa_app_connector_group.test_app_connector_group ]
}

# ZPA App Connector Group resource
resource "zpa_app_connector_group" "test_app_connector_group" {
  name                          = "test1-appconnector-group"
  description                   = "test1-appconnector-group"
  enabled                       = true
  city_country                  = "San Jose, CA"
  country_code                  = "US"
  latitude                      = "37.338"
  longitude                     = "-121.8863"
  location                      = "San Jose, CA, US"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  override_version_profile      = true
  version_profile_id            = 0
  dns_query_type                = "IPV4"
}
```

## Example 2 Usage

```hcl
# ZPA Application Segment resource
resource "zpa_application_segment" "test_app_segment" {
    name              = "test1-app-segment"
    description       = "test1-app-segment"
    enabled           = true
    health_reporting  = "ON_ACCESS"
    bypass_type       = "NEVER"
    is_cname_enabled  = true
  tcp_port_range = [
    {
    from = "8080"
    to = "8080"
    }
  ]
  udp_port_range = [
    {
    from = "8080"
    to = "8080"
    }
  ]
    domain_names      = ["server.acme.com"]
    segment_group_id  = zpa_segment_group.test_segment_group.id
    server_groups {
        id = [ zpa_server_group.test_server_group.id]
    }
    depends_on = [ zpa_server_group.test_server_group, zpa_server_group.test_segment_group]
}

# ZPA Segment Group resource
resource "zpa_segment_group" "test_segment_group" {
  name            = "test1-segment-group"
  description     = "test1-segment-group"
  enabled         = true
}

# ZPA Server Group resource
resource "zpa_server_group" "test_server_group" {
  name              = "test1-server-group"
  description       = "test1-server-group"
  enabled           = true
  dynamic_discovery = false
  app_connector_groups {
    id = [ zpa_app_connector_group.example.id ]
  }
  depends_on = [ zpa_app_connector_group.test_app_connector_group ]
}

# ZPA App Connector Group resource
resource "zpa_app_connector_group" "test_app_connector_group" {
  name                          = "test1-appconnector-group"
  description                   = "test1-appconnector-group"
  enabled                       = true
  city_country                  = "San Jose, CA"
  country_code                  = "US"
  latitude                      = "37.338"
  longitude                     = "-121.8863"
  location                      = "San Jose, CA, US"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  override_version_profile      = true
  version_profile_id            = 0
  dns_query_type                = "IPV4"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name. The name of the App Connector Group to be exported.
* `domain_names` - (Required) List of domains and IPs.
* `server_groups` - (Required) List of Server Group IDs
* `segment_group_id` - (Required) List of Segment Group IDs
* `tcp_port_ranges` - (Required) TCP port ranges used to access the app.
* `udp_port_ranges` - (Required) UDP port ranges used to access the app.

-> **NOTE:**  TCP and UDP ports can also be defined using the following model:

* `tcp_port_range` - (Required) TCP port ranges used to access the app.

  * `from:`
  * `to:`

* `udp_port_range` - (Required) UDP port ranges used to access the app.
  * `from:`
  * `to:`

## Attributes Reference

* `description` - (Optional) Description of the application.
* `bypass_type` - (Optional) Indicates whether users can bypass ZPA to access applications.
* `config_space` - (Optional)
* `double_encrypt` - (Optional) Whether Double Encryption is enabled or disabled for the app.
* `enabled` - (Optional) Whether this application is enabled or not.
* `health_reporting` - (Optional) Whether health reporting for the app is Continuous or On Access. Supported values: NONE, ON_ACCESS, CONTINUOUS.
* `health_check_type` (Optional)
* `icmp_access_type` - (Optional)
* `p_anchored` - (Optional)
* `is_cname_enabled` - (Optional) Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors.
* `log_features` - (Optional)
* `passive_health_enabled` - (Optional)

## Import

Application Segment can be imported by using `<APPLICATION SEGMENT ID>` or `<APPLICATION SEGMENT NAME>` as the import ID.

```shell
terraform import zpa_application_segment.example <application_segment_id>
```

or

```shell
terraform import zpa_application_segment.example <application_segment_name>
```
