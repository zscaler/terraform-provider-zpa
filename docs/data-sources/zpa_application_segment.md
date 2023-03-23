---
subcategory: "Application Segment"
layout: "zscaler"
page_title: "ZPA: application_segment"
description: |-
  Get information about ZPA Application Segment in Zscaler Private Access cloud.
---

# Data Source: zpa_application_segment

Use the **zpa_application_segment** data source to get information about a application segment created in the Zscaler Private Access cloud. This data source can then be referenced in an Access Policy, Timeout policy, Forwarding Policy, Inspection Policy or Isolation Policy.

## Zenith Community - ZPA Application Segment

[![ZPA Terraform provider Video Series Ep7 - Application Segment](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_application_segments.svg)](https://community.zscaler.com/t/video-zpa-terraform-provider-video-series-ep-7-zpa-application-segment/18946)

## Example Usage

```hcl
# ZPA Application Segment Data Source
data "zpa_application_segment" "this" {
  name = "example"
}
```

```hcl
# ZPA Application Segment Data Source
data "zpa_application_segment" "this" {
  id = "123456789"
}
```

## Argument Reference

In addition to all arguments above, the following attributes are exported:

* `name` - (Required) Name of the application.
* `description` - (Optional) Description of the application.
* `bypass_type` - (Optional) Indicates whether users can bypass ZPA to access applications. Default: `NEVER`. Supported values: `ALWAYS`, `NEVER`, `ON_NET`. The value `NEVER` indicates the use of the client forwarding policy.
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
* `select_connector_close_to_app` - (Boolean)
* `segment_group_id` - (string)
* `segment_group_name` - (string)
* `select_connector_close_to_app` (bool)
* `use_in_dr_mode` (bool)
* `is_incomplete_dr_config` (bool)
* `select_connector_close_to_app` (bool)