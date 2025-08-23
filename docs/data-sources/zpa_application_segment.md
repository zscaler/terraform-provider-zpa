---
page_title: "zpa_application_segment Data Source - terraform-provider-zpa"
subcategory: "Application Segment"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-applications/API documentation https://help.zscaler.com/zpa/configuring-application-segments-using-api
  Get information about ZPA Application Segment in Zscaler Private Access cloud.
---

# zpa_application_segment (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-applications)
* [API documentation](https://help.zscaler.com/zpa/configuring-application-segments-using-api)

Use the **zpa_application_segment** data source to get information about a application segment created in the Zscaler Private Access cloud. This data source can then be referenced in an Access Policy, Timeout policy, Forwarding Policy, Inspection Policy or Isolation Policy.

**NOTE:** To ensure consistent search results across data sources, please avoid using multiple spaces or special characters in your search queries.

## Zenith Community - ZPA Application Segment

[![ZPA Terraform provider Video Series Ep7 - Application Segment](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_application_segments.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEXCAY/video-zpa-terraform-provider-video-series-ep7-zpa-application-segment)

## Example Usage

```terraform
# ZPA Application Segment Data Source
data "zpa_application_segment" "this" {
  name = "example"
}
```

```terraform
# ZPA Application Segment Data Source
data "zpa_application_segment" "this" {
  id = "123456789"
}
```

## Schema

### Required

In addition to all arguments above, the following attributes are exported:

- `name` - (Required) Name of the application.

### Read-Only

- `description` - (Optional) Description of the application.
- `enabled` - (Boolean) Whether this application is enabled or not. Default: false. Supported values: `true`, `false`.
- `bypass_type` - (Optional) Indicates whether users can bypass ZPA to access applications. Default: `NEVER`. Supported values: `ALWAYS`, `NEVER`, `ON_NET`. The value `NEVER` indicates the use of the client forwarding policy.
- `bypass_on_reauth` - (Optional) Default: false. Boolean: `true`, `false`.
- `is_cname_enabled` - (Boolean) Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors. Default: true. Boolean: `true`, `false`.
- `health_check_type` - (String) Whether health reporting for the app is Continuous or On Access. Supported values: `NONE`, `ON_ACCESS`, `CONTINUOUS`
- `health_reporting` - (Optional) Supported values: `DEFAULT`, `NONE`.
- `double_encrypt` - (String) Whether Double Encryption is enabled or disabled for the app. Default: false. Boolean: `true`, `false`.
- `icmp_access_type` - (Optional) Supported values: `NONE`, `PING_TRACEROUTING`, `PING`.
- `inspect_traffic_with_zia` - (Optional) Default: false. Boolean: `true`, `false`.
- `tcp_port_ranges` - TCP port ranges used to access the app.
- `udp_port_ranges` - UDP port ranges used to access the app.
- `ip_anchored` - (Boolean) Default: false. Boolean: `true`, `false`.
- `config_space` - (String) Supported values: `DEFAULT`, `SIEM`.
- `passive_health_enabled` - (Boolean) Default: false. Boolean: `true`, `false`.
- `select_connector_close_to_app` - (Boolean) Default: false. Boolean: `true`, `false`.
- `use_in_dr_mode` (Boolean) Default: false. Boolean: `true`, `false`.
- `is_incomplete_dr_config` (Boolean) Default: false. Boolean: `true`, `false`.
- `fqdn_dns_check` - (Boolean) Supported values: `true`, `false`
- `domain_names` - List of domains and IPs.
- `segment_group_id` - (string) The unique identifier of the segment group.
- `segment_group_name` - (string) The name of the segment group
- `default_idle_timeout` - (String)
- `default_max_age` - (String)
- `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
- `microtenant_name` (string) The name of the microtenant the resource is to be associated with.
