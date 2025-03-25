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

- `config_space` - (string)
- `default_idle_timeout` - (string)
- `default_max_age` - (string)
- `domain_names` - List of domains and IPs.
- `health_reporting` - (string)
- `ip_anchored` - (bool)
- `passive_health_enabled` - (bool)
- `segment_group_id` - (string)
- `segment_group_name` - (string)
- `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
- `microtenant_name` (string) The name of the microtenant the resource is to be associated with.

- `clientless_apps`
  - `name` - (string)
  - `application_port` - (string)
  - `application_protocol` - (string)
  - `certificate_id` - (string)
  - `certificate_name` - (string)
  - `domain` - (string)
  - `allow_options` - (bool)
  - `cname` (string)
  - `description` (string)
  - `enabled` (bool)
  - `hidden` (bool)
  - `local_domain` (string)
  - `path` (string)
  - `trust_untrusted_cert` (bool)
  - `select_connector_close_to_app` (bool)
  - `use_in_dr_mode` (bool)
  - `is_incomplete_dr_config` (bool)
  - `select_connector_close_to_app` (bool)
  - `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
  - `microtenant_name` (string) The name of the microtenant the resource is to be associated with.
