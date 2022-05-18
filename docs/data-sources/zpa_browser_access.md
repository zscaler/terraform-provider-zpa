---
subcategory: "Application Segment"
layout: "zscaler"
page_title: "ZPA: application_segment"
description: |-
  Get information about ZPA Browser Access Application Segment in Zscaler Private Access cloud.
---

# zpa_application_segment

Use the **zpa_application_segment** data source to get information about a browser access application segment created in the Zscaler Private Access cloud. This data source can then be referenced in an Access Policy, Timeout policy, Forwarding Policy, Inspection Policy or Isolation Policy.

## Example Usage

```hcl
# ZPA Application Segment Data Source
data "zpa_application_segment" "example" {
  name = "example"
}
```

```hcl
# ZPA Application Segment Data Source
data "zpa_application_segment" "example" {
  id = "123456789"
}
```

## Argument Reference

* `name` - (Required) This field defines the name of the server.
* `id` - (Optional) This field defines the id of the application server.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `name` - (Required) Name of the application.
* `description` - (String) Description of the application.
* `bypass_type` - (String) Indicates whether users can bypass ZPA to access applications. Default: `NEVER`. Supported values: `ALWAYS`, `NEVER`, `ON_NET`. The value `NEVER` indicates the use of the client forwarding policy.
* `is_cname_enabled` - (Boolean) Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors. Default: true. Boolean: `true`, `false`.
* `health_checktype` - (String) Whether health reporting for the app is Continuous or On Access. Supported values: `NONE`, `ON_ACCESS`, `CONTINUOUS`
* `double_encrypt` - (String) Whether Double Encryption is enabled or disabled for the app. Default: false. Boolean: `true`, `false`.
* `enabled` - (Boolean) Whether this application is enabled or not. Default: false. Supported values: `true`, `false`.
* `tcp_port_ranges` - TCP port ranges used to access the app.
* `udp_port_ranges` - UDP port ranges used to access the app.
* `config_space` - (String)
* `default_idle_timeout` - (String)
* `default_max_age` - (String)
* `domain_names` - List of domains and IPs.
* `health_reporting` - (Optional)
* `ip_anchored` - (Boolean)
* `passive_health_enabled` - (Boolean)
* `segment_group_id` - (String)
* `segment_group_name` - (String)

* `clientless_apps`
  * `name` - (String)
  * `application_port` - (String)
  * `application_protocol` - (String)
  * `certificate_id` - (String)
  * `certificate_name` - (String)
  * `domain` - (String)
  * `allow_options` - (Boolean)
  * `cname` (String)
  * `description` (String)
  * `enabled` (Boolean)
  * `hidden` (Boolean)
  * `local_domain` (String)
  * `path` (String)
  * `trust_untrusted_cert` (Boolean)
