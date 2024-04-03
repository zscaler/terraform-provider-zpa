---
page_title: "zpa_cloud_browser_isolation_certificate Resource - terraform-provider-zpa"
subcategory: "Cloud Browser Isolation"
description: |-
  Official documentation https://help.zscaler.com/isolation/about-custom-root-certificates-cloud-browser-isolation
  Creates and manages Cloud Browser Isolation Certificate.
---

# zpa_cloud_browser_isolation_certificate (Resource)

* [Official documentation](https://help.zscaler.com/isolation/about-custom-root-certificates-cloud-browser-isolation)

The **zpa_cloud_browser_isolation_certificate** resource creates a Cloud Browser Isolation certificate. This resource can then be used when creating a CBI External Profile `zpa_cloud_browser_isolation_external_profile`.`

## Example Usage

```terraform
# Retrieve CBI Banner ID
resource "zpa_cloud_browser_isolation_certificate" "this" {
    name = "CBI_Certificate"
    pem = file("cert.pem")
}

resource "zpa_cloud_browser_isolation_certificate" "this" {
    name = "CBI_Certificate"
    pem = <<CERT
    -----BEGIN CERTIFICATE-----
    MIIFYDCCBEigAwIBAgIQQAF3ITfU6UK47naqPGQKtzANBgkqhkiG9w0BAQsFADA/
CERT
}

```

## Schema

### Required

The following arguments are supported:

- `name` - (Required) The name of the CBI certificate.
- `pem` - (Required) The certificate in PEM format.
