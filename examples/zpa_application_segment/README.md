# Application Segment Example

This example will show you how to use Terraform to implement a ZPA application segment resource.
This example codifies [this API](https://help.zscaler.com/zpa/api-reference#/application-controller/addApplication).

To run, configure your ZPA provider as described [Here](https://github.com/SecurityGeekIO/terraform-provider-zpa/blob/master/website/docs/index.html.markdown)

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
