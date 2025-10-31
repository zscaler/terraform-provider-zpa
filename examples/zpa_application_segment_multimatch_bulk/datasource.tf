data "zpa_application_segment_multimatch_bulk" "this" {
  domain_names = ["server1.bd-hashicorp.com", "server2.bd-hashicorp.com"]
}
