package services

import (
	"testing"

	"github.com/webbtech/gsales-xls-reports/config"
)

// This test isn't exactly a unit test as it depends on our Config object,
// but it seems like a minor rule to break?...
var cfg *config.Config

func getConfig(t *testing.T) {
	t.Helper()

	cfg = &config.Config{}
	cfg.Init()
}
