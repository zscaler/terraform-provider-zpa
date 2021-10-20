---
subcategory: "Application Segment"
layout: "zpa"
page_title: "ZPA: application_segment"
description: |-
  Creates a ZPA Application Segment details.
  
---
# zpa_application_segment

The **zpa_application_segment** resource creates an application segment in the Zscaler Private Access cloud. This resource can then be referenced in an access policy rule

## Example Usage

```hcl
# ZPA Application Segment resource
resource "zpa_application_segment" "example" {
    name = "example"
    description = "example"
    enabled = true
    health_reporting = "ON_ACCESS"
    bypass_type = "NEVER"
    is_cname_enabled = true
    tcp_port_ranges = ["8080", "8080"]
    domain_names = ["server.acme.com"]
    segment_group_id = zpa_segment_group.example.id
    server_groups {
        id = [ zpa_server_group.servergroup.id]
    }
}
```

```hcl
# ZPA Segment Group resource
resource "zpa_segment_group" "example" {
  name = "Example"
  description = "Example"
  enabled = true
  policy_migrated = true
}
```

```hcl
# ZPA Server Group resource
resource "zpa_server_group" "example" {
  name = "Example"
  description = "Example"
  enabled = false
  dynamic_discovery = false
  app_connector_groups {
    id = [data.zpa_app_connector_group.example.id]
  }
  servers {
    id = [zpa_application_server.example.id]
  }
}
```

```hcl
data "zpa_app_connector_group" "example" {
  name = "AWS-Connector"
}
```

```hcl
# ZPA Application Server resource
resource "zpa_application_server" "example" {
  name                          = "Example"
  description                   = "Example"
  address                       = "server.acme.com"
  enabled                       = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name. The name of the App Connector Group to be exported.
* `description` (Optional) Description of the application.
* `bypass_type` - (Optional) Indicates whether users can bypass ZPA to access applications.
* `config_space` - (Optional)
* `domain_names` - (Required) List of domains and IPs.
* `double_encrypt` - (Optional) Whether Double Encryption is enabled or disabled for the app.
* `enabled` - (Optional) Whether this application is enabled or not.
* `health_reporting` - (Optional) Whether health reporting for the app is Continuous or On Access. Supported values: `NONE`, `ON_ACCESS`, `CONTINUOUS`.
* `icmp_access_type` - (Optional)
* `ip_anchored` - (Optional)
* `is_cname_enabled` - (Optional) Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors.
* `log_features` - (Optional)
* `passive_health_enabled` - (Optional)
* `segment_group_id` - (Required) ID(s) of the segment group(s).
* `segment_group_name` - (Optional)
* `server_groups` - (Required) ID(s) of the server group(s).
* `tcp_port_ranges` - (Required) TCP port ranges used to access the app.
* `udp_port_ranges` - (Required) UDP port ranges used to access the app.
