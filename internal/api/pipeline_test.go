package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPipelineExtendMarshalOmitsUnsetKeys(t *testing.T) {
	loggingIndex := ""
	body, err := json.Marshal(&PipelineExtend{LoggingIndex: &loggingIndex})

	require.NoError(t, err)
	require.JSONEq(t, `{"loggingIndex":""}`, string(body))

	body, err = json.Marshal(&PipelineExtend{})

	require.NoError(t, err)
	require.JSONEq(t, `{}`, string(body))
}

func TestPipelineExtendUnmarshalNullKeysAreUnset(t *testing.T) {
	var extend PipelineExtend
	err := json.Unmarshal([]byte(`{"appID":null,"loggingIndex":"","measurement":null}`), &extend)

	require.NoError(t, err)
	require.Nil(t, extend.AppID)
	require.Nil(t, extend.Measurement)
	require.NotNil(t, extend.LoggingIndex)
	require.Equal(t, "", *extend.LoggingIndex)
	require.Empty(t, extend.Extra)
}

func TestPipelineExtendUnmarshalWithoutLoggingIndexKey(t *testing.T) {
	var extend PipelineExtend
	err := json.Unmarshal([]byte(`{}`), &extend)

	require.NoError(t, err)
	require.Nil(t, extend.LoggingIndex)
}

func TestPipelineExtendRoundTripKeepsUnknownKeys(t *testing.T) {
	remote := `{"loggingIndex":"default","ifAiCreate":true,"exampleType":"json","aiTargetStructure":{"fields":["a","b"]}}`

	var extend PipelineExtend
	require.NoError(t, json.Unmarshal([]byte(remote), &extend))
	require.Equal(t, "default", *extend.LoggingIndex)
	require.Len(t, extend.Extra, 3)

	body, err := json.Marshal(&extend)

	require.NoError(t, err)
	require.JSONEq(t, remote, string(body))
}

func TestPipelineOmitsNilExtend(t *testing.T) {
	body, err := json.Marshal(&Pipeline{Name: "demo"})

	require.NoError(t, err)
	require.NotContains(t, string(body), `"extend"`)
}

func TestPipelineUnmarshalExtendFromResponseContent(t *testing.T) {
	var pl Pipeline
	err := json.Unmarshal([]byte(`{"uuid":"pl_test","extend":{"loggingIndex":"","exampleSource":"x"}}`), &pl)

	require.NoError(t, err)
	require.NotNil(t, pl.Extend)
	require.Equal(t, "", *pl.Extend.LoggingIndex)
	require.JSONEq(t, `"x"`, string(pl.Extend.Extra["exampleSource"]))
}
