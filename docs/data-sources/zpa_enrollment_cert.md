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

### Lookup by ID

```terraform
data "zpa_enrollment_cert" "connector" {
    id = "6573"
}
```

## Schema

### Optional

The following arguments are supported. One of `id` or `name` must be set; when both are set, `id` takes precedence and `name` is ignored.

* `id` - (String) The ID of the enrollment certificate to be exported.
* `name` - (String) The name of the enrollment certificate to be exported.
* `microtenant_id` - (String) The ID of the microtenant the enrollment certificate belongs to.

### Read-Only

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
* `serial_no` - (string)
* `valid_from_in_epoch_sec` - (string)
* `valid_to_in_epoch_sec` - (string)
* `private_key_present` - (bool) Indicates whether a private key exists for the certificate.
* `private_key` - (string) Always empty; the API does not return private key material.
* `zrsa_encrypted_private_key` - (string) Always empty; the API does not return private key material.
* `zrsa_encrypted_session_key` - (string) Always empty; the API does not return private key material.

~> **Note**: The certificate (`certificate`) and certificate signing request (`csr`) are public material and are included in the output. Private key material is never returned by the API, so `private_key`, `zrsa_encrypted_private_key`, and `zrsa_encrypted_session_key` are always empty; use `private_key_present` to check whether a private key exists.
