# LSS Config Example

This example will show you how to use Terraform to implement the ZPA Log Receiver (LSS Config) resource.
This example codifies [this API](https://help.zscaler.com/zpa/api-reference#/lss-config-controller-v-2).

To run, configure your ZPA provider as described [Here](https://github.com/SecurityGeekIO/terraform-provider-zpa/blob/master/docs/index.md)

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
