---
subcategory: "Application Segment"
layout: "zscaler"
page_title: "ZPA: application_segment_browser_access"
description: |-
  Creates and manages ZPA Browser Access Application Segment.
---

# Resource: zpa_application_segment_browser_access

The **zpa_application_segment_browser_access** creates and manages a browser access application segment resource in the Zscaler Private Access cloud. This resource can then be referenced in an access policy rule, access policy timeout rule or access policy client forwarding rule.

## Zenith Community - ZPA Browser Access Application Segment

[![ZPA Terraform provider Video Series Ep8 - Browser Access Application Segment](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_browser_access_application_segments.svg)](https://community.zscaler.com/t/zpa-terraform-provider-video-series-ep-8-zpa-browser-access-application-segment/19150)

## Example Usage

```hcl
# Retrieve Browser Access Certificate
data "zpa_ba_certificate" "test_cert" {
  name = "sales.acme.com"
}

# Create Browser Access Application
resource "zpa_application_segment_browser_access" "browser_access_apps" {
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
  * `name` - (Required) - Name of BA app.
  * `application_port` - (Required) - Port for the BA app.
  * `application_protocol` - (Required) - Protocol for the BA app. Supported values: `HTTP` and `HTTPS`
  * `certificate_id` - (Required) - ID of the BA certificate. Refer to the data source documentation for [`zpa_ba_certificate`](https://github.com/zscaler/terraform-provider-zpa/blob/master/docs/data-sources/zpa_ba_certificate.md)
  * `domain` - (Required) - Domain name or IP address of the BA app.
  * `allow_options` - (Optional) - If you want ZPA to forward unauthenticated HTTP preflight OPTIONS requests from the browser to the app.. Supported values: `true` and `false`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` (Optional) Description of the application.
* `cname` (Optional)
* `hidden` (Optional)
* `local_domain` (Optional)
* `path` (Optional)
* `trust_untrusted_cert` (Optional)
* `bypass_type` (Optional) Indicates whether users can bypass ZPA to access applications. Default value is: `NEVER` and supported values are: `ALWAYS`, `NEVER` and `ON_NET`. The value `NEVER` indicates the use of the client forwarding policy.
* `config_space` (Optional) Default: `DEFAULT`. Supported values: `DEFAULT`, `SIEM`
* `double_encrypt` (Optional) Whether Double Encryption is enabled or disabled for the app.
* `enabled` (Optional) - Whether this app is enabled or not.
* `health_check_type` (Optional) Default: `DEFAULT`. Supported values: `DEFAULT`, `NONE`
* `health_reporting` (Optional) Whether health reporting for the app is Continuous or On Access. Supported values: `NONE`, `ON_ACCESS`, `CONTINUOUS`.
* `ip_anchored` (Optional) - If Source IP Anchoring for use with ZIA, is enabled or disabled for the app. Supported values are `true` and `false`
* `is_cname_enabled` (Optional) Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors.
  * `certificate_name` - (Optional) - Name of the BA certificate. Refer to the data source documentation for [`zpa_ba_certificate`](https://github.com/zscaler/terraform-provider-zpa/blob/master/docs/data-sources/zpa_ba_certificate.md)
* `select_connector_close_to_app` - (Optional) Supported values: `true`, `false`
* `use_in_dr_mode` - (Optional) Supported values: `true`, `false`
* `is_incomplete_dr_config` - (Optional) Supported values: `true`, `false`
* `select_connector_close_to_app` - (Optional) Supported values: `true`, `false`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zpa_application_segment_browser_access** Application Segment Browser Access can be imported by using <`BROWSER ACCESS ID`> or `<<BROWSER ACCESS NAME>` as the import ID.

For example:

```shell
terraform import zpa_application_segment_browser_access.example <browser_access_id>.
```

or

```shell
terraform import zpa_application_segment_browser_access.example <browser_access_name>
```
