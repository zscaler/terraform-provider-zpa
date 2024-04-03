---
subcategory: "Log Streaming (LSS)"
layout: "zscaler"
page_title: "ZPA: lss_config_status_codes"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-log-streaming-service/
  API documentation https://help.zscaler.com/zpa/configuring-log-streaming-service-configurations-using-api
  Get information about all LSS status codes details.
---

# zpa_lss_config_status_codes (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-log-streaming-service)
* [API documentation](https://help.zscaler.com/zpa/configuring-log-streaming-service-configurations-using-api)

Use the **zpa_lss_config_status_codes** data source to get information about all LSS status codes in the Zscaler Private Access cloud. This data source is required when the `filter` parameter in the LSS Config Controller resource is set. To learn more see the [Getting Details of All LSS Status Codes](https://help.zscaler.com/zpa/log-streaming-service-configuration-use-cases#GettingLSSStatusCodes)

-> **NOTE** By Default the ZPA provider will return all status codes

## Example Usage

```terraform
data "zpa_lss_config_status_codes" "this" {
}
```

## Argument Reference

The following arguments are supported:

To learn more see the [Getting Details of All LSS Status Codes](https://help.zscaler.com/zpa/log-streaming-service-configuration-use-cases#GettingLSSStatusCodes)
