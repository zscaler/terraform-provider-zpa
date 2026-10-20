package zpa

import "testing"

func TestDataSourceServiceEdgeAssistantSchedule(t *testing.T) {
	testAssistantScheduleDataSource(t, dataSourceServiceEdgeAssistantSchedule(), "serviceEdgeSchedule")
}
