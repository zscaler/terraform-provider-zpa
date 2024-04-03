---
page_title: "zpa_lss_config_client_types Data Source - terraform-provider-zpa"
subcategory: "Log Streaming (LSS)"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-log-streaming-service/API documentation https://help.zscaler.com/zpa/configuring-log-streaming-service-configurations-using-api
  Get information about all LSS client type details.
---

# zpa_lss_config_client_types (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-log-streaming-service)
* [API documentation](https://help.zscaler.com/zpa/configuring-log-streaming-service-configurations-using-api)

Use the **zpa_lss_config_client_types** data source to get information about all LSS client types in the Zscaler Private Access cloud. This data source is required when the defining a policy rule resource for an object type as `CLIENT_TYPE` parameter in the LSS Config Controller resource is set. To learn more see the To learn more see the [Getting Details of All LSS Status Codes](https://help.zscaler.com/zpa/log-streaming-service-configuration-use-cases#GettingLSSClientTypes)

-> **NOTE** By Default the ZPA provider will return all client types

## Example Usage

```terraform
data "zpa_lss_config_client_types" "example" {
}
```

### Read-Only

The following arguments are supported:

* `"zpn_client_type_edge_connector" = "Cloud Connector"`
* `"zpn_client_type_exporter" = "Web Browser`
* `"zpn_client_type_ip_anchoring" = "ZIA Service Edge"`
* `"zpn_client_type_machine_tunnel" = "Machine Tunnel"`
* `"zpn_client_type_slogger" = "ZPA LSS"`
* `"zpn_client_type_zapp" = "Client Connector"`

To learn more see the [Getting Details of All LSS Status Codes](https://help.zscaler.com/zpa/log-streaming-service-configuration-use-cases#GettingLSSClientTypes)
