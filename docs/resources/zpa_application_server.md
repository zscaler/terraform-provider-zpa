---
subcategory: "Application Server"
layout: "zscaler"
page_title: "ZPA: application_server"
description: |-
  Creates and manages ZPA Application Servers.
---

# Resource: zpa_application_server (Resource)

The **zpa_application_server** resource creates an application server in the Zscaler Private Access cloud. This resource can then be referenced in a server group.

## Zenith Community - ZPA Application Server

[![ZPA Terraform provider Video Series Ep5 - Application Server](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_application_servers.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEgCAI/video-terraform-provider-video-series-ep5-zpa-application-server)

## Example Usage

```hcl
# ZPA Application Server resource (IP Address)
resource "zpa_application_server" "test_app_server"{
  name                          = "test1-app-server"
  description                   = "test1-app-server"
  address                       = "192.168.1.1"
  enabled                       = true
}
```

```hcl
# ZPA Application Server resource (FQDN Address)
resource "zpa_application_server" "test_app_server" {
  name                          = "test1-app-server"
  description                   = "test1-app-server"
  address                       = "server1.acme.com"
  enabled                       = true
}
```

```hcl
# ZPA Application Server resource
resource "zpa_application_server" "test_app_server"{
  name                          = "test1-app-server"
  description                   = "test1-app-server"
  address                       = "192.168.1.1"
  enabled                       = true
  app_server_group_ids          = [ data.zpa_server_group.example.com ]
}

data "zpa_server_group" "example" {
    name = "Example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name. The name of the application server to be exported.
* `address` - (Required) Address. The address of the application server to be exported.

## Attributes Reference

* `app_server_group_ids` - (Optional) This field defines the list of server group IDs.
* `description` - (Optional) This field defines the description of the server.
* `enabled` - (Optional) This field defines the status of the server.
* `config_space` - (Optional)
* `microtenant_id` (Optional) The ID of the microtenant the resource is to be associated with.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Application Server can be imported by using `<APPLICATION SERVER ID>` or `<APPLICATION SERVER NAME>` as the import ID

For example:

```shell
terraform import zpa_application_server.example <application_server_id>
```

or

```shell
terraform import zpa_application_server.example <application_server_name>
```
