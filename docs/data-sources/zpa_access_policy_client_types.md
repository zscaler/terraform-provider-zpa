---
page_title: "zpa_access_policy_client_types Data Source - terraform-provider-zpa"
subcategory: "Policy Set Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa
  documentation https://help.zscaler.com/zpa/configuring-access-policies-using-api#getClientTypes
  Get information about all client types for the specified customer.
---

# zpa_access_policy_client_types (Data Source)

* [Official documentation](https://help.zscaler.com/zpa)
* [API documentation](https://help.zscaler.com/zpa/configuring-access-policies-using-api#getClientTypes)

Use the **zpa_access_policy_client_types** data source to get information about all client types for the specified customer in the Zscaler Private Access cloud. This data source can be optionally used when defining the following policy types:
    - ``zpa_policy_access_rule``
    - ``zpa_policy_timeout_rule``
    - ``zpa_policy_forwarding_rule``
    - ``zpa_policy_isolation_rule``
    - ``zpa_policy_inspection_rule``

The ``object_type`` attribute must be defined as "CLIENT_TYPE" in the policy operand condition. To learn more see the To learn more see the [Getting Details of All Client Types](https://help.zscaler.com/zpa/configuring-access-policies-using-api#getClientTypes)

-> **NOTE** By Default the ZPA provider will return all client types

-> **NOTE** When defining a ``zpa_policy_isolation_rule`` policy the ``object_type`` "CLIENT_TYPE" is mandatory and ``zpn_client_type_exporter`` is the only supported value.

## Example Usage

```terraform
data "zpa_access_policy_client_types" "this" {
}
```

## Schema

### Read-Only

The following values are returned:

* `"zpn_client_type_exporter": "Web Browser"`
* `"zpn_client_type_exporter_noauth": "Web Browser Unauthenticated"`
* `"zpn_client_type_machine_tunnel": "Machine Tunnel"`
* `"zpn_client_type_edge_connector": "Cloud Connector"`
* `"zpn_client_type_zia_inspection": "ZIA Inspection"`
* `"zpn_client_type_vdi": "Client Connector for VDI"`
* `"zpn_client_type_zapp": "Client Connector"`
* `"zpn_client_type_slogger": "ZPA LSS"`
* `"zpn_client_type_browser_isolation": "Cloud Browser"`
* `"zpn_client_type_ip_anchoring": "ZIA Service Edge"`
* `"zpn_client_type_zapp_partner": "Client Connector Partner"`
* `"zpn_client_type_branch_connector": "Branch Connector"`

To learn more see the [Getting Details of All Client Types](https://help.zscaler.com/zpa/configuring-access-policies-using-api#getClientTypes)
