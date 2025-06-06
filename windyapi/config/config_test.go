package config

import (
	"bytes"
	"github.com/go-andiamo/cfgenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"slices"
	"strings"
	"testing"
)

func TestShowExampleConfig(t *testing.T) {
	cfg := &Config{}
	var w bytes.Buffer
	err := cfgenv.Example(&w, cfg)
	assert.NoError(t, err)
	lines := strings.Split(strings.TrimSuffix(w.String(), "\n"), "\n")
	expect := []string{
		"VERSION",
		"BUILD",
		"DATABASE_HOST",
		"DATABASE_PORT",
		"DATABASE_NAME",
		"DATABASE_USERNAME",
		"DATABASE_PASSWORD",
		"DATABASE_TIMEOUT",
		"PROCESS_REAL_TIME_DIRECTORY",
		"PROCESS_MAX_RECORDS_TO_PROCESS",
		"PROCESS_HISTORY_RECORDS",
		"PROCESS_BATCH_COUNT",
		"PROCESS_HISTORICAL_FILE_NAME",
		"PROCESS_SAVE_HISTORICAL_WEATHER_SQL_LOC",
		"PROCESS_CONCURRENT_REQUESTS",
		"PROCESS_MAX_WORKER_POOLS",
		"REQUEST_MODEL",
		"REQUEST_PARAMETERS",
		"REQUEST_LEVELS",
		"REQUEST_API_KEY",
		"WINDYAPI_ENDPOINT",
		"RESPONSE_SAVE_RESPONSE",
	}
	assert.Equal(t, len(expect), len(lines))
	for _, s := range expect {
		println(s)
		assert.True(t, slices.ContainsFunc(lines, func(l string) bool {
			return strings.HasPrefix(l, s+"=")
		}))
	}
}

func TestLoad(t *testing.T) {
	_, err := Load()
	require.Error(t, err)
}
