---
page_title: "zpa_application_segment_browser_access Data Source - terraform-provider-zpa"
subcategory: "Application Segment"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-browser-access
  API documentation https://help.zscaler.com/zpa/configuring-browser-access-application-segments-using-api
  Get information about ZPA Browser Access Application Segment in Zscaler Private Access cloud.
---

# zpa_application_segment_browser_access (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-browser-access)
* [API documentation](https://help.zscaler.com/zpa/configuring-browser-access-application-segments-using-api)

Use the **zpa_application_segment_browser_access** data source to get information about a browser access application segment created in the Zscaler Private Access cloud. This data source can then be referenced in an Access Policy, Timeout policy, Forwarding Policy, Inspection Policy or Isolation Policy.

**NOTE:** To ensure consistent search results across data sources, please avoid using multiple spaces or special characters in your search queries.

## Zenith Community - ZPA Browser Access Application Segment

[![ZPA Terraform provider Video Series Ep8 - Browser Access Application Segment](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_browser_access_application_segments.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEGCAY/zpa-terraform-provider-video-series-ep8-zpa-browser-access-application-segment)

## Example Usage

```terraform
# ZPA Application Segment Browser Access Data Source
data "zpa_application_segment_browser_access" "example" {
  name = "example"
}
```

```terraform
# ZPA Application Segment Browser Access Data Source
data "zpa_application_segment_browser_access" "example" {
  id = "123456789"
}
```

## Schema

### Required

- `name` - (String) This field defines the name of the server.
- `id` - (String) This field defines the id of the application server.

### Read-Only

In addition to all arguments above, the following attributes are exported:

- `description` - (string) Description of the application.
- `bypass_type` - (string) Indicates whether users can bypass ZPA to access applications. Default: `NEVER`. Supported values: `ALWAYS`, `NEVER`, `ON_NET`. The value `NEVER` indicates the use of the client forwarding policy.
- `is_cname_enabled` - (bool) Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors. Default: true. Boolean: `true`, `false`.
- `health_checktype` - (string) Whether health reporting for the app is Continuous or On Access. Supported values: `NONE`, `ON_ACCESS`, `CONTINUOUS`
- `double_encrypt` - (string) Whether Double Encryption is enabled or disabled for the app. Default: false. Boolean: `true`, `false`.
- `enabled` - (Boolean) Whether this application is enabled or not. Default: false. Supported values: `true`, `false`.
- `tcp_port_ranges` - (string) TCP port ranges used to access the app.
- `udp_port_ranges` - (string) UDP port ranges used to access the app.

-> **NOTE:**  TCP and UDP ports can also be defined using the following model:

- `tcp_port_range` - (string) TCP port ranges used to access the app.
  - `from:`
  - `to:`
- `udp_port_range` - (string) UDP port ranges used to access the app.
  - `from:`
  - `to:`

- `config_space` - (string) Indicates if the configuration is created as part of the SIEM or application resource.
- `default_idle_timeout` - (string) The duration of the default idle timeout.
- `default_max_age` - (string) The default maximum age of the resource.
- `domain_names` - List of domains and IPs.
- `health_reporting` - (string) Indicates the health reporting of the application.
- `ip_anchored` - (bool) Whether Source IP Anchoring for use with ZIA is enabled (true) or disabled (false) for the application.
- `passive_health_enabled` - (bool) Indicates if passive health checks are enabled on the application.
- `segment_group_id` - (string) The unique identifier of the segment group.
- `segment_group_name` - (string) The name of the segment group
- `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
- `microtenant_name` (string) The name of the microtenant the resource is to be associated with.

- `clientless_apps`
  - `name` - (string) The name of the Browser Access application segment.
  - `application_port` - (string) The port for the Browser Access application.
  - `application_protocol` - (string) The protocol for the Browser Access application.
  - `certificate_id` - (string) The unique identifier of the certificate.
  - `certificate_name` - (string) The name of the certificate.
  - `domain` - (string)
  - `allow_options` - (bool)
  - `cname` (string) The canonical name (CNAME DNS records) of the Browser Access application.
  - `description` (string) The description of the Browser Access application segment.
  - `enabled` (bool) Indicates whether the Browser Access application segment is enabled (true) or not (false).
  - `hidden` (bool)
  - `local_domain` (string) The local domain of the Browser Access application.
  - `path` (string) The path of the Browser Access application.
  - `trust_untrusted_cert` (String) Whether the use of untrusted certificates is enabled or disabled for the Browser Access application. Supported values are `true` and `false`
  - `use_in_dr_mode` - (Boolean) Whether or not the application resource is designated for disaster recovery. Supported values: `true`, `false`
  - `is_incomplete_dr_config` - (Boolean) Indicates whether or not the disaster recovery configuration is incomplete. Supported values: `true`, `false`
  - `select_connector_close_to_app` - (Boolean) Whether the App Connector is closest to the application (true) or closest to the user (false). Supported values: `true`, `false`
  - `ext_label` (String) The domain prefix for the privileged portal URL. The supported string can include numbers, lower case characters, and only supports a hyphen (-).
  - `ext_domain` (String) The external domain name prefix of the Browser Access application that is used for Zscaler-managed certificates when creating a privileged portal. This field is returned when making GET requests to get privileged portal details or when retreiving application segment details. The supported value must be a string, but doesn't support special characters (e.g., periods) as the FQDN wouldn't match the CNAME entry.
  - `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
  - `microtenant_name` (string) The name of the microtenant the resource is to be associated with.
