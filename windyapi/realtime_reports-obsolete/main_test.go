package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckForTableExistence(t *testing.T) {
	count, _ := app.DB.CheckForTableExistence(dbName, weatherStations)
	assert.True(t, count == 1)
}

func TestCheckForInvalidTableExistence(t *testing.T) {
	count, _ := app.DB.CheckForTableExistence(dbName, "invalid_table")
	assert.True(t, count == 0)
}

func TestCheckRowExistence(t *testing.T) {
	count, _ := app.DB.CheckRowExistence()
	assert.True(t, count == 1)
}

func TestValidateDirectory(t *testing.T) {
	fileDirectories := []string{validDirectory}
	IsValidDir := validatefileDirectories(fileDirectories)
	assert.True(t, IsValidDir)
}

func TestValidateInvalidDirectory(t *testing.T) {
	fileDirectories := []string{invalidDirectory}
	IsValidDir := validatefileDirectories(fileDirectories)
	assert.False(t, IsValidDir)
}
