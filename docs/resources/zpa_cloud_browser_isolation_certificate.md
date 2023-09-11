---
subcategory: "Cloud Browser Isolation"
layout: "zscaler"
page_title: "ZPA: cloud_browser_isolation_external_profile"
description: |-
  Creates and manages Cloud Browser Isolation Certificate.
---

# Resource: zpa_cloud_browser_isolation_certificate

The **zpa_cloud_browser_isolation_certificate** resource creates a Cloud Browser Isolation certificate. This resource can then be used when creating a CBI External Profile `zpa_cloud_browser_isolation_external_profile`.`

## Example Usage

```hcl
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

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the CBI certificate.
* `pem` - (Required) The certificate in PEM format.
