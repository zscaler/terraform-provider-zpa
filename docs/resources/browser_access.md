---
subcategory: "Browser Access"
layout: "zpa"
page_title: "ZPA: browser_access"
description: |-
  Creates a ZPA Browser Access.
  
---
# zpa_browser_access

The **zpa_browser_access** creates an browser access resource in the Zscaler Private Access cloud. This resource can then be referenced in an access policy rule, access policy timeout rule or access policy client forwarding rule.

## Example Usage

```hcl
Create Browser Access Application
resource "zpa_browser_access" "browser_access_apps" {
    name = "Browser Access Apps"
    description = "Browser Access Apps"
    enabled = true
    health_reporting = "ON_ACCESS"
    bypass_type = "NEVER"
    tcp_port_ranges = ["80", "80"]
    domain_names = ["sales.acme.com"]
    segment_group_id = zpa_segment_group.example.id

    clientless_apps {
        name = "sales.acme.com"
        application_protocol = "HTTP"
        application_port = "80"
        certificate_id = data.zpa_ba_certificate.sales_ba.id
        trust_untrusted_cert = true
        enabled = true
        domain = "sales.acme.com"
    }
    server_groups {
        id = [zpa_server_group.example.id]
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

data "zpa_ba_certificate" "sales_ba" {
    name = "sales.acme.com"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the application.
* `description` - (Optional) Description of the application.
* `enabled` - (Optional) Whether this app is enabled or not. Default: `false`. Boolean values: `true`, `false`.
* `bypass_type` - (Optional) Indicates whether users can bypass ZPA to access applications.
* `config_space` - (Optional) Default: `DEFAULT`. Supported values: `DEFAULT`, `SIEM`
* `domain_names` - (Required) List of domains and IPs.
* `double_encrypt` - (Optional) Whether Double Encryption is enabled or disabled for the app.
* `health_check_type` - (Optional)
* `health_reporting` - (Optional) Whether health reporting for the app is Continuous or On Access. Supported values: `NONE`, `ON_ACCESS`, `CONTINUOUS`.
* `ip_anchored` - (Optional) If Source IP Anchoring for use with ZIA, is enabled or disabled for the app.
* `is_cname_enabled` - (Optional) Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors. Default: `true`. Boolean values: `true`, `false`.
* `segment_group_id` - (Required) ID(s) of the segment group(s).
* `segment_group_name` - (Optional)
* `tcp_port_ranges` - (Required) TCP port ranges used to access the app.
* `udp_port_ranges` - (Required) UDP port ranges used to access the app.
* `clientless_apps` - (Block List) (see [below for nested schema](#nestedblock--clientless_apps))
* `server_groups` (Block List) ID of the server group. (see [below for nested schema](#nestedblock--server_groups))

<a id="nestedblock--clientless_apps"></a>

### clientless_apps

* `name` - (Required) Name of the application.
* `description` - (Optional) Description of the application.
* `enabled` - (Optional) Whether this app is enabled or not. Default: `false`. Boolean values: `true`, `false`.
* `allow_options` - (Optional)
* `application_port` - (Required)
* `application_protocol` - (Required)
* `certificate_id` - (Required)
* `certificate_name` - (Required)
* `cname` - (Optional)
* `domain` - (Required)
* `hidden` - (Optional)
* `local_domain` - (Optional)
* `path` - (Optional)
* `trust_untrusted_cert` - (Optional)
