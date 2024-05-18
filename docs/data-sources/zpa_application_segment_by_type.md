---
page_title: "zpa_application_segment_by_type Data Source - terraform-provider-zpa"
subcategory: "Application Segment By Type"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-applications
  API documentation https://help.zscaler.com/zpa/configuring-application-segments-using-api
  Get information about all configured enrollment certificate details.
---

# zpa_application_segment_by_type (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-applications)
* [API documentation](https://help.zscaler.com/zpa/configuring-application-segments-using-api)

Use the **zpa_application_segment_by_type** data source to get all configured Application Segments by Access Type (e.g., ``BROWSER_ACCESS``, ``INSPECT``, or ``SECURE_REMOTE_ACCESS``) for the specified customer.

## Example Usage

```terraform
data "zpa_application_segment_by_type" "this" {
    application_type = "BROWSER_ACCESS"
}

data "zpa_application_segment_by_type" "this" {
    application_type = "INSPECT"
}

data "zpa_application_segment_by_type" "this" {
    application_type = "SECURE_REMOTE_ACCESS"
}
```

## Schema

### Required

The following arguments are supported:

* `application_type` - (String) The name of the enrollment certificate to be exported.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The unique identifier of the Browser Access, inspection or secure remote access application.
* `app_id` - (String) The unique identifier of the application.
* `name` - (String) The name of the Browser Access, inspection or secure remote access application.
* `enabled` - (bool) Whether the Browser Access, inspection or secure remote access application is enabled or not
* `domain` - (string) The domain of the Browser Access, inspection or secure remote access application
* `application_port` - (string) The port for the Browser Access, inspection or secure remote access application
* `application_protocol` - (string) The protocol for the Browser Access, inspection or secure remote access application

* `certificate_id` - (string) The unique identifier of the Browser Access certificate
* `certificate_name` - (string) The name of the Browser Access certificate
* `microtenant_id` - (string) The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant
* `microtenant_name` - (string) The name of the Microtenant

