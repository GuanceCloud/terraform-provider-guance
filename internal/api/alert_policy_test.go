package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListAlertPoliciesWithOptionsPagesUntilShortPage(t *testing.T) {
	var requestedPages []int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v1/alert_policy/list", r.URL.Path)
		require.Equal(t, "terraform", r.URL.Query().Get("search"))
		require.Equal(t, "notify_xxx,notify_yyy", r.URL.Query().Get("notifyObjectUUIDs"))
		require.Equal(t, strconv.Itoa(alertPolicyListPageSize), r.URL.Query().Get("pageSize"))

		pageIndex, err := strconv.Atoi(r.URL.Query().Get("pageIndex"))
		require.NoError(t, err)
		requestedPages = append(requestedPages, pageIndex)

		data := make([]map[string]any, alertPolicyListPageSize)
		if pageIndex == 2 {
			data = data[:2]
		}
		for i := range data {
			data[i] = map[string]any{
				"uuid": "altpl_page_" + strconv.Itoa(pageIndex) + "_" + strconv.Itoa(i),
				"name": "terraform",
			}
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"code":    200,
			"success": true,
			"content": map[string]any{
				"data": data,
			},
		})
	}))
	defer server.Close()

	client := &Client{
		EndPoint:   server.URL,
		HTTPClient: server.Client(),
	}
	var content AlertPolicyListContent

	err := client.ListAlertPoliciesWithOptions(AlertPolicyListOptions{
		Search:            "terraform",
		NotifyObjectUUIDs: []string{"notify_xxx", "", "notify_yyy"},
	}, &content)

	require.NoError(t, err)
	require.Equal(t, []int{1, 2}, requestedPages)
	require.Len(t, content.Data, alertPolicyListPageSize+2)
	require.Equal(t, "altpl_page_1_0", content.Data[0].UUID)
	require.Equal(t, "altpl_page_2_1", content.Data[len(content.Data)-1].UUID)
}
