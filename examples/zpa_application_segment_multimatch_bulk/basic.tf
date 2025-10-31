resource "zpa_application_segment_multimatch_bulk" "this" {
  application_ids = ["72058304855164528", "72058304855164536"]
  match_style     = "EXCLUSIVE"
}
