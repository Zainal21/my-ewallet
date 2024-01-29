package helpers

import (
	"database/sql"

	"github.com/Zainal21/my-ewallet/app/appctx"
	"github.com/Zainal21/my-ewallet/pkg/logger"
)

func SuccessResponse(message string, data interface{}, statusCode int) appctx.Response {
	return *appctx.NewResponse().
		WithStatus("OK").
		WithMessage(message).
		WithCode(statusCode).
		WithData(data)
}

func HandleError(err error) error {
	if err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			logger.Error(err.Error())
			return err
		}
		logger.Info(err.Error())
	}

	return nil
}
