---
page_title: "zpa_application_segment_browser_access Resource - terraform-provider-zpa"
subcategory: "Application Segment"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-browser-access
  API documentation https://help.zscaler.com/zpa/configuring-browser-access-application-segments-using-api
  Creates and manages ZPA Browser Access Application Segment
---

# zpa_application_segment_browser_access (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-browser-access)
* [API documentation](https://help.zscaler.com/zpa/configuring-browser-access-application-segments-using-api)

The **zpa_application_segment_browser_access** creates and manages a browser access application segment resource in the Zscaler Private Access cloud. This resource can then be referenced in an access policy rule, access policy timeout rule or access policy client forwarding rule.

## Zenith Community - ZPA Browser Access Application Segment

[![ZPA Terraform provider Video Series Ep8 - Browser Access Application Segment](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_browser_access_application_segments.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEGCAY/zpa-terraform-provider-video-series-ep8-zpa-browser-access-application-segment)

## Example Usage

```terraform
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

## Schema

### Required

The following arguments are supported:

- `name` - (String) The name of the application resource.
- `domain_names` - (List of String) List of domains and IPs.
- `tcp_port_ranges` - (List of String) TCP port ranges used to access the app.
- `udp_port_ranges` - (List of String) UDP port ranges used to access the app.

-> **NOTE:**  TCP and UDP ports can also be defined using the following model:
-> **NOTE:** When removing TCP and/or UDP ports, parameter must be defined but set as empty due to current API behavior.

- `tcp_port_range` - (Block Set) TCP port ranges used to access the app.
  - `from:` (String) The starting port for a port range.
  - `to:` (String) The ending port for a port range.

- `udp_port_range` (Block Set) UDP port ranges used to access the app.
  - `from:` (String) The starting port for a port range.
  - `to:` (String) The ending port for a port range.

- `server_groups` (Block Set) List of Server Group IDs
  - `id` - (Required)

- `segment_group_id` (String) List of Segment Group IDs

- `clientless_apps` (Block Set)
  - `name` - (String) - Name of BA app.
  - `application_port` - (String) - Port for the BA app.
  - `application_protocol` - (String) - Protocol for the BA app. Supported values: `HTTP` and `HTTPS`
  - `certificate_id` - (String) - ID of the BA certificate. Refer to the data source documentation for [`zpa_ba_certificate`](https://github.com/zscaler/terraform-provider-zpa/blob/master/docs/data-sources/zpa_ba_certificate.md)
  - `domain` - (String) - Domain name or IP address of the BA app.
  - `allow_options` - (Boolean) - If you want ZPA to forward unauthenticated HTTP preflight OPTIONS requests from the browser to the app.. Supported values: `true` and `false`
  - `microtenant_id` (Boolean) The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenant_id as `0` when making requests to retrieve data from the Default Microtenant. Pass microtenant_id as null to retrieve data from all customers associated with the tenant.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

### Optional

In addition to all arguments above, the following attributes are exported:

- `description` (String) The description of the Browser Access application.
- `trust_untrusted_cert` (String) Whether the use of untrusted certificates is enabled or disabled for the Browser Access application. Supported values are `true` and `false`
- `bypass_type` (String) Indicates whether users can bypass ZPA to access applications. Default value is: `NEVER` and supported values are: `ALWAYS`, `NEVER` and `ON_NET`. The value `NEVER` indicates the use of the client forwarding policy.
- `double_encrypt` (Boolean) Whether Double Encryption is enabled or disabled for the app. Supported values are `true` and `false`
- `enabled` (Boolean) - Whether the Browser Access application is enabled or not. Supported values are `true` and `false`
- `health_check_type` (String) Whether the health check is enabled (`DEFAULT`) or disabled (`NONE`) for the application. Default: `DEFAULT`. Supported values: `DEFAULT`, `NONE`
- `health_reporting` (String) Whether health reporting for the app is Continuous or On Access. Supported values: `NONE`, `ON_ACCESS`, `CONTINUOUS`.
- `icmp_access_type` - (String) The ICMP access type. Supported values: `PING_TRACEROUTING`, `PING`, `NONE`
- `ip_anchored` (Boolean) - If Source IP Anchoring for use with ZIA, is enabled or disabled for the app. Supported values are `true` and `false`
- `fqdn_dns_check` - (Boolean) When set to Enabled, Zscaler Client Connector receives CNAME DNS records from the App Connector for FQDN applications. Supported values: `true`, `false`
-  `is_cname_enabled` (Boolean) Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors.
- `certificate_name` - (String) - Name of the BA certificate. Refer to the data source documentation for [`zpa_ba_certificate`](https://github.com/zscaler/terraform-provider-zpa/blob/master/docs/data-sources/zpa_ba_certificate.md)
- `tcp_keep_alive` (String) Supported values: ``1`` for Enabled and ``0`` for Disabled

- `select_connector_close_to_app` - (Boolean) Whether the App Connector is closest to the application (true) or closest to the user (false). Supported values: `true`, `false`

- `use_in_dr_mode` - (Boolean) Whether or not the application resource is designated for disaster recovery. Supported values: `true`, `false`
- `is_incomplete_dr_config` - (Boolean) Indicates whether or not the disaster recovery configuration is incomplete. Supported values: `true`, `false`
- `microtenant_id` (Boolean) The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenant_id as `0` when making requests to retrieve data from the Default Microtenant. Pass microtenant_id as null to retrieve data from all customers associated with the tenant.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

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
