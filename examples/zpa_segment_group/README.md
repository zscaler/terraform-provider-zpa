# Segment Group Example

This example will show you how to create a Segment Group in the ZPA portal.
This example codifies [this API](https://help.zscaler.com/zpa/api-reference#/segment-group-controller).

To run, configure your ZPA provider as described [Here](https://github.com/zscaler/terraform-provider-zpa/blob/master/docs/index.html.markdown)

## Run the example

From inside of this directory:

```bash
terraform init
terraform plan -out theplan
terraform apply theplan
```

## Destroy 💥

```bash
terraform destroy
```

## Deprecated Attributes

- ``policy_migrated``: "The `policy_migrated` field is now deprecated for the resource `zpa_segment_group`, please remove this attribute to prevent configuration drifts"

- ``tcp_keep_alive_enabled``: "The `tcp_keep_alive_enabled` field is now deprecated for the resource `zpa_segment_group`, please replace all uses of this within the `zpa_application_segment`resources with the attribute `tcp_keep_alive`"
