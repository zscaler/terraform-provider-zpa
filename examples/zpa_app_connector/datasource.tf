data "zpa_app_connector" "example" {
  name = "AWS-VPC100-App-Connector"
}

output "zpa_app_connector" {
  value = data.zpa_app_connector.example
}