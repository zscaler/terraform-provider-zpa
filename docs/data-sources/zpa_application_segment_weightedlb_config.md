---
page_title: "zpa_application_segment_weightedlb_config Data source - terraform-provider-zpa"
subcategory: "Application Segment"
description: |-
  Retrieve the weighted load balancer configuration for an application segment.
---

# zpa_application_segment_weightedlb_config (Data Source)

Use the **zpa_application_segment_weightedlb_config** data source to retrieve the weighted load balancer configuration for an application segment in the Zscaler Private Access cloud.

## Example Usage

```terraform
data "zpa_application_segment_weightedlb_config" "this" {
  application_name = "app02"
}
```

## Schema

### Optional

The following arguments are supported:

* `application_id` - (Optional) The unique identifier of the application segment. Either `application_id` or `application_name` must be provided.
* `application_name` - (Optional) The name of the application segment. Either `application_id` or `application_name` must be provided.
* `microtenant_id` - (Optional) The microtenant identifier to scope the request.

### Read-Only

The following attributes are exported:

* `weighted_load_balancing` - Indicates if the application load balancing configuration for application segments is enabled (`true`) or disabled (`false`).
* `application_to_server_group_mappings` - A list of server group mappings associated with the application segment.

Each `application_to_server_group_mappings` block exports the following:

* `id` - The unique identifier of the server group mapping.
* `name` - The name of the server group.
* `passive` - Indicates if the server group is in a passive state (`true`) or can load balance requests (`false`).
* `weight` - The weight assigned to the server group. Higher weights indicate that a greater number of requests are routed to the associated App Connectors; weights of `0` do not route requests.
