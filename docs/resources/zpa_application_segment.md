---
subcategory: "Application Segment"
layout: "zscaler"
page_title: "ZPA: application_segment"
description: |-
  Creates and manages ZPA Application Segments.
---

# Resource: zpa_application_segment

The **zpa_application_segment** resource creates an application segment in the Zscaler Private Access cloud. This resource can then be referenced in an access policy rule, access policy timeout rule or access policy client forwarding rule.

## Zenith Community - ZPA Application Segment

[![ZPA Terraform provider Video Series Ep7 - Application Segment](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_application_segments.svg)](https://community.zscaler.com/t/video-zpa-terraform-provider-video-series-ep-7-zpa-application-segment/18946)

## Example 1 Usage

```hcl
# ZPA Application Segment resource
resource "zpa_application_segment" "this" {
    name              = "Example"
    description       = "Example"
    enabled           = true
    health_reporting  = "ON_ACCESS"
    bypass_type       = "NEVER"
    is_cname_enabled  = true
    tcp_port_ranges   = ["8080", "8080"]
    domain_names      = ["server.acme.com"]
    segment_group_id  = zpa_segment_group.this.id
    server_groups {
        id = [ zpa_server_group.this.id]
    }
    depends_on = [ zpa_server_group.this, zpa_segment_group.this]
}

# ZPA Segment Group resource
resource "zpa_segment_group" "this" {
  name            = "Example"
  description     = "Example"
  enabled         = true
}

# ZPA Server Group resource
resource "zpa_server_group" "this" {
  name              = "Example"
  description       = "Example"
  enabled           = true
  dynamic_discovery = false
  app_connector_groups {
    id = [ zpa_app_connector_group.this.id ]
  }
  depends_on = [ zpa_app_connector_group.this ]
}

# ZPA App Connector Group resource
resource "zpa_app_connector_group" "this" {
  name                          = "Example"
  description                   = "Example"
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
resource "zpa_application_segment" "this" {
    name              = "Example"
    description       = "Example"
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
    segment_group_id  = zpa_segment_group.this.id
    server_groups {
        id = [ zpa_server_group.this.id]
    }
    depends_on = [ zpa_server_group.this, zpa_segment_group.this]
}

# ZPA Segment Group resource
resource "zpa_segment_group" "this" {
  name            = "Example"
  description     = "Example"
  enabled         = true
}

# ZPA Server Group resource
resource "zpa_server_group" "this" {
  name              = "Example"
  description       = "Example"
  enabled           = true
  dynamic_discovery = false
  app_connector_groups {
    id = [ zpa_app_connector_group.this.id ]
  }
  depends_on = [ zpa_app_connector_group.this ]
}

# ZPA App Connector Group resource
resource "zpa_app_connector_group" "this" {
  name                          = "Example"
  description                   = "Example"
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
-> **NOTE:** When removing TCP and/or UDP ports, parameter must be defined but set as empty due to current API behavior.

* `tcp_port_range` - (Required) TCP port ranges used to access the app.
  * `from:`
  * `to:`

* `udp_port_range` - (Required) UDP port ranges used to access the app.
  * `from:`
  * `to:`

-> **NOTE:** Application segments must have unique ports and cannot have overlapping domain names using the same tcp/udp ports across multiple application segments.

## Attributes Reference

* `description` - (Optional) Description of the application.
* `bypass_type` - (Optional) Indicates whether users can bypass ZPA to access applications.
* `config_space` - (Optional)
* `double_encrypt` - (Optional) Whether Double Encryption is enabled or disabled for the app.
* `enabled` - (Optional) Whether this application is enabled or not.
* `health_reporting` - (Optional) Whether health reporting for the app is Continuous or On Access. Supported values: `NONE`, `ON_ACCESS`, `CONTINUOUS`.
* `health_check_type` (Optional)
* `icmp_access_type` - (Optional)
* `ip_anchored` - (Optional)
* `is_cname_enabled` - (Optional) Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors.
* `log_features` - (Optional)
* `tcp_keep_alive` (Optional) Supported values: ``1`` for Enabled and ``0`` for Disabled
* `passive_health_enabled` - (Optional) Supported values: `true`, `false`
* `select_connector_close_to_app` - (Optional) Supported values: `true`, `false`
* `use_in_dr_mode` - (Optional) Supported values: `true`, `false`
* `is_incomplete_dr_config` - (Optional) Supported values: `true`, `false`
* `select_connector_close_to_app` - (Optional) Supported values: `true`, `false`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Application Segment can be imported by using `<APPLICATION SEGMENT ID>` or `<APPLICATION SEGMENT NAME>` as the import ID.

```shell
terraform import zpa_application_segment.example <application_segment_id>
```

or

```shell
terraform import zpa_application_segment.example <application_segment_name>
```
