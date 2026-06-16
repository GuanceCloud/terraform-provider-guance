package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListMonitorsWithOptionsPagesUntilShortPage(t *testing.T) {
	var requestedPages []int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v1/checker/list", r.URL.Path)
		require.Equal(t, "terraform", r.URL.Query().Get("search"))
		require.Equal(t, "trigger", r.URL.Query().Get("checkerTypes"))
		require.Equal(t, "0", r.URL.Query().Get("checkerStatus"))
		require.Equal(t, "tag_xxx", r.URL.Query().Get("tagsUUID"))
		require.Equal(t, "altpl_xxx", r.URL.Query().Get("alertPolicyUUID"))
		require.Equal(t, "dsbd_xxx", r.URL.Query().Get("dashboardUUID"))
		require.Equal(t, "rul_xxx", r.URL.Query().Get("checkerUUID"))
		require.Equal(t, strconv.Itoa(monitorListPageSize), r.URL.Query().Get("pageSize"))

		pageIndex, err := strconv.Atoi(r.URL.Query().Get("pageIndex"))
		require.NoError(t, err)
		requestedPages = append(requestedPages, pageIndex)

		data := make([]map[string]any, monitorListPageSize)
		if pageIndex == 2 {
			data = data[:2]
		}
		for i := range data {
			data[i] = map[string]any{
				"uuid": "rul_page_" + strconv.Itoa(pageIndex) + "_" + strconv.Itoa(i),
				"type": "trigger",
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
	var content MonitorListContent

	err := client.ListMonitorsWithOptions(MonitorListOptions{
		Search:          "terraform",
		Type:            "trigger",
		Status:          "0",
		TagsUUID:        "tag_xxx",
		AlertPolicyUUID: "altpl_xxx",
		DashboardUUID:   "dsbd_xxx",
		CheckerUUID:     "rul_xxx",
	}, &content)

	require.NoError(t, err)
	require.Equal(t, []int{1, 2}, requestedPages)
	require.Len(t, content.Data, monitorListPageSize+2)
	require.Equal(t, "rul_page_1_0", content.Data[0].UUID)
	require.Equal(t, "rul_page_2_1", content.Data[len(content.Data)-1].UUID)
}
