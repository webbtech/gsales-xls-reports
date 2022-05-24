package handlers

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/opentracing/opentracing-go/log"

	"github.com/webbtech/gsales-xls-reports/config"
	lerrors "github.com/webbtech/gsales-xls-reports/errors"
	"github.com/webbtech/gsales-xls-reports/model"
	"github.com/webbtech/gsales-xls-reports/report"
	"github.com/webbtech/gsales-xls-reports/utils"
)

type Report struct {
	Cfg           *config.Config
	Db            model.DbHandler
	input         *model.RequestInput
	reportRequest *model.ReportRequest
	request       events.APIGatewayProxyRequest
	response      events.APIGatewayProxyResponse
}

const (
	CODE_SUCCESS                = "SUCCESS"
	ERR_INVALID_DATES           = "Invalid date parameters"
	ERR_INVALID_TYPE            = "Invalid report type in input"
	ERR_MISSING_REQUEST_BODY    = "Missing request body"
	ERR_FAILED_TO_CREATE_REPORT = "Failed to create report"
)

const (
	timeDayFormat   = "2006-01-02"
	timeMonthFormat = "2006-01"
)

var stage string

// ========================== Public Methods =============================== //

func (r *Report) Response(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r.request = request
	r.process()
	return r.response, nil
}

// ========================== Private Methods ============================== //

func (r *Report) process() {

	rb := responseBody{}
	var body []byte
	var err error
	var rpt *report.Report
	var statusCode int = 201
	var stdError *lerrors.StdError
	var url string

	// validate input
	if err := r.validateInput(); err != nil {
		errors.As(err, &stdError)
	}

	// create report
	if stdError == nil {
		rpt, err = report.New(r.reportRequest, r.Cfg, r.Db)
		if err != nil {
			stdError = &lerrors.StdError{
				Caller:     "handlers.process",
				Code:       lerrors.CodeApplicationError,
				Err:        err,
				Msg:        ERR_FAILED_TO_CREATE_REPORT,
				StatusCode: 500,
			}
		}
	}

	// create signed url
	if stdError == nil {
		url, err = rpt.CreateSignedURL()
		if err != nil {
			stdError = &lerrors.StdError{
				Caller:     "handlers.process",
				Code:       lerrors.CodeApplicationError,
				Err:        err,
				Msg:        ERR_FAILED_TO_CREATE_REPORT,
				StatusCode: 500,
			}
		}
	}

	// Process any errors
	if stdError != nil {
		rb.Code = stdError.Code
		rb.Message = stdError.Msg
		statusCode = stdError.StatusCode
		logError(stdError)
	} else {
		rb.Code = CODE_SUCCESS
		rb.Message = "Success"
		rb.Data = url
	}

	// Create the response object
	body, _ = json.Marshal(&rb)
	r.response = events.APIGatewayProxyResponse{
		Body:       string(body),
		Headers:    headers,
		StatusCode: statusCode,
	}
}

func (r *Report) validateInput() (err *lerrors.StdError) {

	r.reportRequest = &model.ReportRequest{}

	json.Unmarshal([]byte(r.request.Body), &r.input)

	// ensure there's a request body
	if r.input == nil {
		return &lerrors.StdError{
			Caller:     "handlers.validateInput",
			Code:       lerrors.CodeBadInput,
			Err:        errors.New(ERR_MISSING_REQUEST_BODY),
			Msg:        ERR_MISSING_REQUEST_BODY,
			StatusCode: 400,
		}
	}

	// validate and create report type
	rt, error := model.ReportStringToType(r.input.ReportType)
	if error != nil {
		return &lerrors.StdError{
			Caller:     "handlers.validateInput",
			Code:       lerrors.CodeBadInput,
			Err:        errors.New(ERR_INVALID_TYPE),
			Msg:        ERR_INVALID_TYPE,
			StatusCode: 400,
		}
	}
	r.reportRequest.ReportType = rt

	// validate and create dates
	dts, error := r.createDates()
	if error != nil {
		return &lerrors.StdError{
			Caller:     "handlers.validateInput",
			Code:       lerrors.CodeBadInput,
			Err:        errors.New(ERR_INVALID_DATES),
			Msg:        ERR_INVALID_DATES,
			StatusCode: 400,
		}
	}
	r.reportRequest.Dates = dts

	return nil
}

// createDates function
func (r *Report) createDates() (dates *model.RequestDates, err error) {

	input := r.input
	dates = &model.RequestDates{}

	if input.Date == "" && (input.DateFrom == "" || input.DateTo == "") {
		return nil, errors.New("Missing dates in handler.CreateDates")
	}

	// if it's a date, then we create a date range for the month requested
	if input.Date != "" {
		dates.DateFrom, dates.DateTo, err = utils.DatesFromMonth(input.Date)
		if err != nil {
			return dates, err
		}

		// else we should have a start and end date
	} else {
		dates.DateFrom, dates.DateTo, err = utils.DatesFromDays(input.DateFrom, input.DateTo)
		if err != nil {
			return dates, err
		}
	}

	return dates, nil
}

// NOTE: these could go into it's own package
func logError(err *lerrors.StdError) {
	if stage == "" {
		stage = os.Getenv("Stage")
	}

	if stage != "test" {
		log.Error(err)
	}
}
