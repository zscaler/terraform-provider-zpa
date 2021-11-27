---
subcategory: "Application Segment"
layout: "zpa"
page_title: "ZPA: application_segment"
description: |-
  Gets a ZPA Application Segment details.

---
# zpa_application_segment

The **zpa_application_segment** data source provides details about a specific application segments created in the Zscaler Private Access cloud. This data source must be used in the following circumstances:

1. Access policy rule

## Example Usage

```hcl
# ZPA Application Segment Data Source
data "zpa_application_segment" "foo" {
  name = "example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the application.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

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
* `segment_group_id` - (Optional)
* `segment_group_name` - (Optional)
