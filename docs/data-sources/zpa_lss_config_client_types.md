---
subcategory: "Log Streaming (LSS)"
layout: "zscaler"
page_title: "ZPA: lss_config_client_types"
description: |-
  Get information about all LSS client type details.
---

# zpa_lss_config_client_types

Use the **zpa_lss_config_client_types** data source to get information about all LSS client types in the Zscaler Private Access cloud. This data source is required when the defining a policy rule resource for an object type as `CLIENT_TYPE` parameter in the LSS Config Controller resource is set. To learn more see the To learn more see the [Getting Details of All LSS Status Codes](https://help.zscaler.com/zpa/log-streaming-service-configuration-use-cases#GettingLSSClientTypes)

-> **NOTE** By Default the ZPA provider will return all client types

## Example Usage

```hcl
data "zpa_lss_config_client_types" "example" {
}
```

## Argument Reference

The following arguments are supported:

* `"zpn_client_type_edge_connector" = "Cloud Connector"`
* `"zpn_client_type_exporter" = "Web Browser`
* `"zpn_client_type_ip_anchoring" = "ZIA Service Edge"`
* `"zpn_client_type_machine_tunnel" = "Machine Tunnel"`
* `"zpn_client_type_slogger" = "ZPA LSS"`
* `"zpn_client_type_zapp" = "Client Connector"`

To learn more see the [Getting Details of All LSS Status Codes](https://help.zscaler.com/zpa/log-streaming-service-configuration-use-cases#GettingLSSClientTypes)
