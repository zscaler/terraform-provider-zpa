---
subcategory: "Application Segment"
layout: "zscaler"
page_title: "ZPA: browser_access"
description: |-
  Creates and manages ZPA Browser Access Application Segment.
---

# Resource: zpa_browser_access

The **zpa_browser_access** creates and manages a browser access application segment resource in the Zscaler Private Access cloud. This resource can then be referenced in an access policy rule, access policy timeout rule or access policy client forwarding rule.

## Example Usage

```hcl
# Retrieve Browser Access Certificate
data "zpa_ba_certificate" "test_cert" {
  name = "sales.acme.com"
}

# Create Browser Access Application
resource "zpa_browser_access" "browser_access_apps" {
    name                      = "Browser Access Apps"
    description               = "Browser Access Apps"
    enabled                   = true
    health_reporting          = "ON_ACCESS"
    bypass_type               = "NEVER"
    tcp_port_ranges           = ["80", "80"]
    domain_names              = ["sales.acme.com"]
    segment_group_id          = zpa_segment_group.example.id

    clientless_apps {
        name                  = "sales.acme.com"
        application_protocol  = "HTTP"
        application_port      = "80"
        certificate_id        = data.zpa_ba_certificate.test_cert.id
        trust_untrusted_cert  = true
        enabled               = true
        domain                = "sales.acme.com"
    }
    server_groups {
        id = [zpa_server_group.example.id]
    }
}

# ZPA Segment Group resource
resource "zpa_segment_group" "example" {
  name          = "Example"
  description   = "Example"
  enabled       = true
}

# ZPA Server Group resource
resource "zpa_server_group" "example" {
  name                  = "Example"
  description           = "Example"
  enabled               = true
  dynamic_discovery     = true
  app_connector_groups {
    id = [ data.zpa_app_connector_group.example.id ]
  }
}

data "zpa_app_connector_group" "example" {
  name = "AWS-Connector"
}

```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Name of the application.
* `domain_names` - (Required) List of domains and IPs.
* `tcp_port_ranges` - (Required) TCP port ranges used to access the app.
* `udp_port_ranges` - (Required) UDP port ranges used to access the app.

* `server_groups` - (Required) List of Server Group IDs
  * `id` - (Required)

* `segment_group_id` - (Required) List of Segment Group IDs
  * `id` - (Required)

* `clientless_apps`
  * `name` - (Required)
  * `application_port` - (Required)
  * `application_protocol` - (Required)
  * `certificate_id` - (Required)
  * `certificate_name` - (Required)
  * `domain` - (Required)
  * `allow_options` - (Optional)
  * `cname` (Optional)
  * `description` (Optional)
  * `enabled` (Optional)
  * `hidden` (Optional)
  * `local_domain` (Optional)
  * `path` (Optional)
  * `trust_untrusted_cert` (Optional)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` (Optional) Description of the application.
* `bypass_type` (Optional) Indicates whether users can bypass ZPA to access applications.
* `config_space` (Optional)
* `double_encrypt` (Optional) Whether Double Encryption is enabled or disabled for the app.
* `enabled` (Optional)
* `health_check_type` (Optional)
* `health_reporting` (Optional) Whether health reporting for the app is Continuous or On Access. Supported values: NONE, ON_ACCESS, CONTINUOUS.
* `ip_anchored` (Optional)
* `is_cname_enabled` (Optional) Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors.

## Import

**zpa_browser_access** Application Segment can be imported by using <`BROWSER ACCESS ID`> or `<<BROWSER ACCESS NAME>` as the import ID.

For example:

```shell
terraform import zpa_browser_access.example <browser_access_id>.
```

or

```shell
terraform import zpa_segment_group.example <browser_access_name>
```
