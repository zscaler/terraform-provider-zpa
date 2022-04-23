# Service Edge Group Example

This example will show you how to use Terraform to implement the ZPA Service Edge Group resource.
This example codifies [this API](https://help.zscaler.com/zpa/api-reference#/service-edge-group-controller).

To run, configure your ZPA provider as described [Here](https://github.com/willguibr/terraform-provider-zpa/blob/master/docs/index.html.markdown)

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
