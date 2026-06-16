package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListNotifyObjectsPagesUntilShortPage(t *testing.T) {
	var requestedPages []int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v1/notify_object/list", r.URL.Path)
		require.Equal(t, "terraform", r.URL.Query().Get("search"))
		require.Equal(t, strconv.Itoa(notifyObjectListPageSize), r.URL.Query().Get("pageSize"))

		pageIndex, err := strconv.Atoi(r.URL.Query().Get("pageIndex"))
		require.NoError(t, err)
		requestedPages = append(requestedPages, pageIndex)

		data := make([]map[string]any, notifyObjectListPageSize)
		if pageIndex == 2 {
			data = data[:2]
		}
		for i := range data {
			data[i] = map[string]any{
				"uuid": "notify_page_" + strconv.Itoa(pageIndex) + "_" + strconv.Itoa(i),
				"name": "terraform",
				"type": "simpleHTTPRequest",
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
	var content NotifyObjectListContent

	err := client.ListNotifyObjects("terraform", &content)

	require.NoError(t, err)
	require.Equal(t, []int{1, 2}, requestedPages)
	require.Len(t, content.Data, notifyObjectListPageSize+2)
	require.Equal(t, "notify_page_1_0", content.Data[0].UUID)
	require.Equal(t, "notify_page_2_1", content.Data[len(content.Data)-1].UUID)
}
