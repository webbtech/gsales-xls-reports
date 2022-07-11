package config

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"
)

var cfg *Config

func TestSetDefaults(t *testing.T) {
	t.Run("sets defaults with local file", func(t *testing.T) {

		cfg = &Config{IsDefaultsLocal: true}
		cfg.setDefaults()
		dir, _ := os.Getwd()
		expectedFilePath := path.Join(dir, defaultsFileName)
		if expectedFilePath != defaultsFilePath {
			t.Fatalf("defaultsFilePath should be %s, have: %s", expectedFilePath, defaultsFilePath)
		}
	})

	t.Run("sets defaults with remote file", func(t *testing.T) {

		cfg = &Config{}
		err := cfg.setDefaults()
		if err != nil {
			t.Fatalf("expected nil error, have: %s", err)
		}
	})
}

// validateStage is called at various times including in setEnvVars
func TestValidateStage(t *testing.T) {
	cfg = &Config{}
	cfg.setDefaults()

	t.Run("stage set from defaults file", func(t *testing.T) {
		if cfg.Stage != ProdEnv {
			t.Fatalf("Stage value should be: %s, have: %s", ProdEnv, cfg.Stage)
		}
	})

	t.Run("stage set from environment", func(t *testing.T) {
		os.Setenv("Stage", "test")
		defer os.Unsetenv("Stage")
		cfg.setEnvVars() // calls validateStage
		if cfg.Stage != TestEnv {
			t.Fatalf("Stage value should be: %s, have: %s", TestEnv, cfg.Stage)
		}
	})

	t.Run("stage set from invalid environment variable", func(t *testing.T) {
		os.Setenv("Stage", "testit")
		defer os.Unsetenv("Stage")
		err := cfg.setEnvVars()
		if err == nil {
			t.Fatalf("Expected validateStage to return error")
		}
	})

	t.Run("stage set with SetStageEnv method", func(t *testing.T) {
		err := cfg.SetStageEnv("stage")
		if err != nil {
			t.Fatalf("Expected null error received: %s", err)
		}
	})

	t.Run("invalid stage set with SetStageEnv method", func(t *testing.T) {
		err := cfg.SetStageEnv("stageit")
		if err == nil {
			t.Fatalf("Expected validateStage error")
		}
	})
}

func TestSetSSMParams(t *testing.T) {
	cfg = &Config{IsDefaultsLocal: true}
	cfg.setDefaults()
	cfg.SetStageEnv("test")

	t.Run("DbName is accurate", func(t *testing.T) {
		err := cfg.setSSMParams()
		if err != nil {
			t.Fatalf("Expected null error, received: %s", err)
		}

		if defs.DbName == "" {
			t.Fatalf("Expected defs.DbName to have value")
		}
		if defs.CognitoClientID == "" {
			t.Fatalf("Expected defs.CognitoClientID to have value")
		}

	})
}

func TestSetEnvVars(t *testing.T) {
	cfg = &Config{IsDefaultsLocal: true}
	cfg.setDefaults()
	if defs.Stage != "prod" {
		t.Fatalf("Expected defs.Stage to be: prod, got: %s", defs.Stage)
	}

	cfg.SetStageEnv("test")
	err := cfg.setEnvVars()
	if err != nil {
		t.Fatalf("Expected null error, received: %s", err)
	}

	expectedStage := "test"
	if defs.Stage != expectedStage {
		t.Fatalf("Expected defs.Stage to be: %s, got: %s", expectedStage, defs.Stage)
	}
}

func TestInitConfig(t *testing.T) {
	cfg = &Config{IsDefaultsLocal: true}
	err := cfg.Init()
	if err != nil {
		t.Fatalf("Expected null error received: %s", err)
	}
}

func TestInitConfigProd(t *testing.T) {
	os.Setenv("Stage", "prod")
	cfg = &Config{}
	err := cfg.Init()
	if err != nil {
		t.Fatalf("Expected null error received: %s", err)
	}
}

func TestPublicGetters(t *testing.T) {
	t.Run("GetStageEnv", func(t *testing.T) {
		stg := TestEnv

		cfg = &Config{IsDefaultsLocal: true}
		cfg.setDefaults()
		cfg.SetStageEnv("test")

		receivedEnv := cfg.GetStageEnv()
		if receivedEnv != stg {
			t.Fatalf("Stage should be: %s, got: %s", stg, receivedEnv)
		}
	})

	// GetDbName()
	t.Run("GetDbName", func(t *testing.T) {
		cfg = &Config{IsDefaultsLocal: true}
		cfg.setDefaults()
		cfg.setFinal()

		expectedDbNm := "gales-sales"
		receivedDbNm := cfg.GetDbName()
		if receivedDbNm != expectedDbNm {
			t.Fatalf("DbName should be: %s, got: %s", expectedDbNm, receivedDbNm)
		}
	})

	t.Run("GetMongoConnectURL local", func(t *testing.T) {
		os.Setenv("Stage", "test")
		cfg = &Config{IsDefaultsLocal: true}
		cfg.Init()

		receivedUrl := cfg.GetMongoConnectURL()
		expectedUrl := fmt.Sprintf("mongodb://%s/?readPreference=primary&ssl=false&directConnection=true", defs.DbHost)
		if receivedUrl != expectedUrl {
			t.Fatalf("Expected url: %s, got: %s", expectedUrl, receivedUrl)
		}
	})

	t.Run("GetMongoConnectURL production", func(t *testing.T) {
		os.Setenv("Stage", "prod")
		cfg = &Config{IsDefaultsLocal: true}
		cfg.Init()

		receivedUrl := cfg.GetMongoConnectURL()
		expectedUrl := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", defs.DbUser, defs.DbPassword, defs.DbHost)
		if receivedUrl != expectedUrl {
			t.Fatalf("Expected url: %s, got: %s", expectedUrl, receivedUrl)
		}
	})
}

// This test does NOT run successfully when running the `run file tests` command, otherwise fine...
func TestUrlExpireTime(t *testing.T) {
	t.Run("sets expireTime", func(t *testing.T) {
		cfg = &Config{IsDefaultsLocal: true}
		cfg.Init()

		expectedHrs := time.Duration(time.Duration(defs.ExpireHrs) * time.Hour)
		if expectedHrs != cfg.UrlExpireTime {
			t.Fatalf("UrlExpireTime should be: %v, have: %v", expectedHrs, cfg.UrlExpireTime)
		}
	})
}
