package handlers

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/webbtech/gsales-xls-reports/config"
	"github.com/webbtech/gsales-xls-reports/model"
	"github.com/webbtech/gsales-xls-reports/mongodb"
)

func TestCreateDates(t *testing.T) {
	t.Run("Missing all dates", func(t *testing.T) {
		r := &Report{}
		r.request.Body = `{"type": "payperiod"}`
		json.Unmarshal([]byte(r.request.Body), &r.input)

		_, err := r.createDates()
		if err == nil {
			t.Fatal("Expected error")
		}
	})

	t.Run("Missing dateFrom", func(t *testing.T) {
		r := &Report{}
		r.request.Body = `{"type": "payperiod", "dateTo": "2021-04-10"}`
		json.Unmarshal([]byte(r.request.Body), &r.input)

		_, err := r.createDates()
		if err == nil {
			t.Fatal("Expected error")
		}
	})

	t.Run("date success", func(t *testing.T) {

		r := &Report{}
		r.request.Body = `{"type": "payperiod", "date": "2022-04"}`
		json.Unmarshal([]byte(r.request.Body), &r.input)

		dts, err := r.createDates()
		if err != nil {
			t.Fatalf("Expected null error received: %s", err)
		}

		var timeType time.Time
		expectedType := reflect.TypeOf(timeType)

		if expectedType != reflect.TypeOf(dts.DateFrom) {
			t.Fatalf("Expected date type: %s, got: %s", expectedType, reflect.TypeOf(dts.DateFrom))
		}
		if expectedType != reflect.TypeOf(dts.DateTo) {
			t.Fatalf("Expected date type: %s, got: %s", expectedType, reflect.TypeOf(dts.DateTo))
		}
	})

	t.Run("dateFrom and dateTo success", func(t *testing.T) {

		r := &Report{}
		r.request.Body = `{"type": "payperiod", "dateFrom": "2022-03-21", "dateTo": "2022-04-10"}`
		json.Unmarshal([]byte(r.request.Body), &r.input)

		dts, err := r.createDates()
		if err != nil {
			t.Fatalf("Expected null error received: %s", err)
		}

		// to validate time, look at: https://stackoverflow.com/questions/6996704/how-to-check-variable-type-at-runtime-in-go-language

		var timeType time.Time
		expectedType := reflect.TypeOf(timeType)

		if expectedType != reflect.TypeOf(dts.DateFrom) {
			t.Fatalf("Expected date type: %s, got: %s", expectedType, reflect.TypeOf(dts.DateFrom))
		}
		if expectedType != reflect.TypeOf(dts.DateTo) {
			t.Fatalf("Expected date type: %s, got: %s", expectedType, reflect.TypeOf(dts.DateTo))
		}
	})
}

func TestValidateRequestType(t *testing.T) {
	t.Run("successful report type", func(t *testing.T) {

		r := &Report{}
		r.request.Body = `{"type": "payperiod"}`
		json.Unmarshal([]byte(r.request.Body), &r.input)

		rt, err := model.ReportStringToType(r.input.ReportType)
		if err != nil {
			t.Fatalf("Expected null error received: %s", err)
		}

		expectedRT := model.PayPeriodReport
		if rt != expectedRT {
			t.Fatalf("Expected report type: %d, got: %d", expectedRT, rt)
		}
	})

	t.Run("invalid report type", func(t *testing.T) {
		r := &Report{}
		r.request.Body = `{"type": "payperiodd"}`
		json.Unmarshal([]byte(r.request.Body), &r.input)

		_, err := model.ReportStringToType(r.input.ReportType)
		if err == nil {
			t.Fatal("Expected error")
		}
	})
}

func TestValidateInput(t *testing.T) {
	t.Run("Missing request body", func(t *testing.T) {
		r := &Report{}

		expectedErr := ERR_MISSING_REQUEST_BODY
		err := r.validateInput()
		if err.Msg != expectedErr {
			t.Fatalf("Error should be: %s, have: %s", expectedErr, err.Msg)
		}
	})

	t.Run("Missing report type", func(t *testing.T) {
		r := &Report{}
		r.request.Body = `{"date": "2022-03"}`

		err := r.validateInput()
		if err.Msg != ERR_INVALID_TYPE {
			t.Fatalf("Expected err: %s, got: %s", ERR_INVALID_TYPE, err.Msg)
		}
	})

	t.Run("Missing date fields", func(t *testing.T) {
		r := &Report{}
		r.request.Body = `{"type": "payperiod"}`

		err := r.validateInput()
		if err.Msg != ERR_INVALID_DATES {
			t.Fatalf("Expected err: %s, got: %s", ERR_INVALID_DATES, err.Msg)
		}
	})

	t.Run("success", func(t *testing.T) {
		r := &Report{}
		r.request.Body = `{"type": "fuelsales", "date":"2020-08"}`

		err := r.validateInput()
		if err != nil {
			t.Fatalf("Unexpected err: %s", err)
		}
	})
}

// TODO: require test here to successfully open a signed url
func TestProcess(t *testing.T) {
	os.Setenv("Stage", "test")
	cfg := &config.Config{}
	err := cfg.Init()
	if err != nil {
		t.Fatalf("Expected null error received: %s", err)
	}

	db, err := mongodb.NewDB(cfg.GetMongoConnectURL(), cfg.DbName)
	if err != nil {
		t.Fatalf("Failed to initial db with error: %s", err)
	}

	r := &Report{Cfg: cfg, Db: db}
	r.request.Body = `{"type": "fuelsales", "date":"2022-04"}`
	r.process()

	expectedStatusCode := 201
	if expectedStatusCode != r.response.StatusCode {
		t.Fatalf("Status should be: %d, have: %d", expectedStatusCode, r.response.StatusCode)
	}

	expectedSuccessCode := CODE_SUCCESS
	var responseBody = &responseBody{}
	json.Unmarshal([]byte(r.response.Body), &responseBody)

	if expectedSuccessCode != responseBody.Code {
		t.Fatalf("SuccessCode should be: %s, have: %s", expectedSuccessCode, responseBody.Code)
	}

	expectedMessageStart := "https://gsales-reports.s3.ca-central-1.amazonaws.com/"
	if !strings.HasPrefix(responseBody.Data, expectedMessageStart) {
		t.Fatalf("Expected message to start with: %s", expectedMessageStart)
	}
}
