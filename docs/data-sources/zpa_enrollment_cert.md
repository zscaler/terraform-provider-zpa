---
subcategory: "Enrollment Certificate"
layout: "zscaler"
page_title: "ZPA: enrollment_cert"
description: |-
  Get information about all configured enrollment certificate details.
---

# Data Source: zpa_enrollment_cert

Use the **zpa_enrollment_cert** data source to get information about all configured enrollment certificate details created in the Zscaler Private Access cloud. This data source is required when creating provisioning key resources.

## Example Usage

```hcl
data "zpa_enrollment_cert" "root" {
    name = "Root"
}

data "zpa_enrollment_cert" "client" {
    name = "Client"
}

data "zpa_enrollment_cert" "connector" {
    name = "Connector"
}

data "zpa_enrollment_cert" "service_edge" {
    name = "Service Edge"
}

data "zpa_enrollment_cert" "isolation_client" {
    name = "Isolation Client"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the enrollment certificate to be exported.
* `id` - (Optional) The id of the enrollment certificate to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `allow_signing` - (bool)
* `cname` - (string)
* `certificate` - (string) The certificate text is in PEM format.
* `client_cert_type` - (string) Returned values are:
  * `ZAPP_CLIENT`
  * `ISOLATION_CLIENT`
  * `NONE`

* `creation_time` - (string)
* `csr` - (string)
* `description` - (string)
* `issued_by` - (string)
* `issued_to` - (string)
* `modified_time` - (string)
* `modified_by` - (string)
* `parent_cert_id` - (string)
* `parent_cert_name` - (string)
* `cert_chain` - (string)
* `serial_no` - (string)
* `valid_from_in_epoch_sec` - (string)
* `valid_to_in_epochsec` - (string)

:warning: Notice that certificate, public and private key information are omitted from the output.
