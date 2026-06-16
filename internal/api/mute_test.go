package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetMutePagesPastTwentyPages(t *testing.T) {
	const targetUUID = "mute_target"
	var requestedPages []int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v1/monitor/mute/list", r.URL.Path)

		pageIndex, err := strconv.Atoi(r.URL.Query().Get("pageIndex"))
		require.NoError(t, err)
		require.Equal(t, strconv.Itoa(muteLookupPageSize), r.URL.Query().Get("pageSize"))
		requestedPages = append(requestedPages, pageIndex)

		data := make([]map[string]any, muteLookupPageSize)
		for i := range data {
			data[i] = map[string]any{"uuid": "mute_page_" + strconv.Itoa(pageIndex) + "_" + strconv.Itoa(i)}
		}
		if pageIndex == 21 {
			data[0] = map[string]any{
				"uuid": targetUUID,
				"name": "target mute",
				"type": "alertPolicy",
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
	var content MuteContent

	err := client.GetMute(targetUUID, &content)

	require.NoError(t, err)
	require.Equal(t, targetUUID, content.UUID)
	require.Equal(t, "target mute", content.Name)
	require.Len(t, requestedPages, 21)
	require.Equal(t, 21, requestedPages[len(requestedPages)-1])
}

func TestGetMuteStopsAfterShortPage(t *testing.T) {
	var requestedPages []int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageIndex, err := strconv.Atoi(r.URL.Query().Get("pageIndex"))
		require.NoError(t, err)
		requestedPages = append(requestedPages, pageIndex)

		data := make([]map[string]any, muteLookupPageSize)
		if pageIndex == 2 {
			data = data[:muteLookupPageSize-1]
		}
		for i := range data {
			data[i] = map[string]any{"uuid": "mute_page_" + strconv.Itoa(pageIndex) + "_" + strconv.Itoa(i)}
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

	err := client.GetMute("mute_missing", &MuteContent{})

	require.ErrorIs(t, err, Error404)
	require.Equal(t, []int{1, 2}, requestedPages)
}

func TestListMutesWithOptionsPagesUntilShortPage(t *testing.T) {
	var requestedPages []int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v1/monitor/mute/list", r.URL.Path)
		require.Equal(t, "terraform", r.URL.Query().Get("search"))
		require.Equal(t, "2", r.URL.Query().Get("workStatus"))
		require.Equal(t, "false", r.URL.Query().Get("isEnable"))
		require.Equal(t, "custom", r.URL.Query().Get("type"))
		require.Equal(t, "acnt_creator", r.URL.Query().Get("creator"))
		require.Equal(t, "acnt_updator", r.URL.Query().Get("updator"))
		require.Equal(t, strconv.Itoa(muteLookupPageSize), r.URL.Query().Get("pageSize"))

		pageIndex, err := strconv.Atoi(r.URL.Query().Get("pageIndex"))
		require.NoError(t, err)
		requestedPages = append(requestedPages, pageIndex)

		data := make([]map[string]any, muteLookupPageSize)
		if pageIndex == 2 {
			data = data[:2]
		}
		for i := range data {
			data[i] = map[string]any{
				"uuid": "mute_page_" + strconv.Itoa(pageIndex) + "_" + strconv.Itoa(i),
				"name": "terraform",
				"type": "custom",
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
	var content MuteListContent

	err := client.ListMutesWithOptions(MuteListOptions{
		Search:     "terraform",
		WorkStatus: "2",
		IsEnable:   "false",
		Type:       "custom",
		Creator:    "acnt_creator",
		Updator:    "acnt_updator",
	}, &content)

	require.NoError(t, err)
	require.Equal(t, []int{1, 2}, requestedPages)
	require.Len(t, content.Data, muteLookupPageSize+2)
	require.JSONEq(t, `{"uuid":"mute_page_1_0","name":"terraform","type":"custom"}`, string(content.Data[0]))
	require.JSONEq(t, `{"uuid":"mute_page_2_1","name":"terraform","type":"custom"}`, string(content.Data[len(content.Data)-1]))
}
