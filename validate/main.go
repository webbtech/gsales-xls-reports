package validate

import (
	"github.com/webbtech/gsales-xls-reports/model"
	"github.com/webbtech/gsales-xls-reports/util"
)

// SetRequest function
func SetRequest(input *model.RequestInput) (req *model.ReportRequest, err error) {

	var rt model.ReportType

	req = &model.ReportRequest{}
	req.Dates, err = util.CreateDates(input)
	if err != nil {
		return req, err
	}

	rt, err = model.ReportStringToType(input.ReportType)
	if err != nil {
		return req, err
	}

	req.ReportType = &rt
	return req, err
}
