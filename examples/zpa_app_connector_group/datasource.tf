data "zpa_app_connector_group" "aws-connector-group" {
  name = "AWS Connector Group"
}

output "get_app_connector_group" {
  value = data.zpa_app_connector_group.aws-connector-group.id
}