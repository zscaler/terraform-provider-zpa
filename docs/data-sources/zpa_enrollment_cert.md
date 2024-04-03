---
page_title: "zpa_enrollment_cert Data Source - terraform-provider-zpa"
subcategory: "Enrollment Certificate"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-enrollment-ca-certificates
  API documentation https://help.zscaler.com/zpa/obtaining-enrollment-certificate-details-using-api
  Get information about all configured enrollment certificate details.
---

# zpa_enrollment_cert (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-enrollment-ca-certificates)
* [API documentation](https://help.zscaler.com/zpa/obtaining-enrollment-certificate-details-using-api)

Use the **zpa_enrollment_cert** data source to get information about all configured enrollment certificate details created in the Zscaler Private Access cloud. This data source is required when creating provisioning key resources.

## Example Usage

```terraform
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

## Schema

### Required

The following arguments are supported:

* `name` - (String) The name of the enrollment certificate to be exported.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The id of the enrollment certificate to be exported.
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

~> **Warning**: Notice that certificate, public and private key information are omitted from the output.
