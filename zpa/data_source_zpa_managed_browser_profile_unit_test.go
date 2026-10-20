package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/managed_browser"
)

func TestDataSourceManagedBrowserProfile(t *testing.T) {
	match := managed_browser.ManagedBrowserProfile{ID: "mb1", Name: "Chrome Enterprise", BrowserType: "CHROME"}
	testDataSource(t, dataSourceSpec{
		resource:   dataSourceManagedBrowserProfile(),
		listPath:   v1Path("managedBrowserProfile", "search"),
		idFromList: true,
		items:      []managed_browser.ManagedBrowserProfile{{ID: "other", Name: "Chrome Enterprise Legacy"}, match},
		item:       match,
		name:       "Chrome Enterprise",
		id:         "mb1",
		wantState:  map[string]any{"browser_type": "CHROME"},
	})
}
