package transaction

import (
	"fmt"
	"strconv"

	"github.com/Zainal21/my-ewallet/app/appctx"
	"github.com/Zainal21/my-ewallet/app/consts"
	"github.com/Zainal21/my-ewallet/app/controller/contract"
	"github.com/Zainal21/my-ewallet/app/dtos"
	"github.com/Zainal21/my-ewallet/app/helpers"
	"github.com/Zainal21/my-ewallet/app/service"
	"github.com/Zainal21/my-ewallet/pkg/config"
	"github.com/Zainal21/my-ewallet/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type callbackMidtransImpl struct {
	service  service.UserService
	transSrv service.TransactionService
	cfg      *config.Config
}

type PaymentNotificationCallback struct {
	TransactionType   string `json:"transaction_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	ReferenceID       string `json:"reference_id"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	ExpiryTime        string `json:"expiry_time"`
	Currency          string `json:"currency"`
	Acquirer          string `json:"acquirer"`
}

// Serve implements contract.Controller.
// Serve implements contract.Controller.
func (g *callbackMidtransImpl) Serve(xCtx appctx.Data) appctx.Response {
	ctx := xCtx.FiberCtx
	body := new(PaymentNotificationCallback)

	if err := xCtx.FiberCtx.BodyParser(&body); err != nil {
		logger.Error(fmt.Sprintf("Body Parser Request Callback Data Request : %v", err))
		return helpers.CreateErrorResponse(fiber.StatusBadRequest, "INVALID PAYLOAD", nil)
	}

	logger.Info(fmt.Sprintf("body callback request %v", body))

	var (
		status         string
		transactionLog bool
	)

	switch body.TransactionStatus {
	case "pending":
		status = "process"
		logger.Info(fmt.Sprintf("status transaction pending %v", body.TransactionStatus))
	case "deny", "cancel", "expire", "failure":
		status = "failed"
		logger.Info(fmt.Sprintf("status transaction cancel/expire/deny/failure %v", body.TransactionStatus))
	case "settlement", "capture":
		status = "success"
		transactionLog = true
		logger.Info(fmt.Sprintf("status transaction success %v", body.TransactionStatus))
	}

	trns, err := g.transSrv.GetTransactionByFieldName(ctx.Context(), "order_id", body.OrderID)
	if err != nil {
		logger.Error(fmt.Sprintf("error get transaction by field name %v", err.Error()))
		return helpers.CreateErrorResponse(fiber.StatusInternalServerError, consts.ServerErrorMessage, nil)
	}

	err = g.transSrv.UpdateStatusTransactionLog(ctx.Context(), status, body.OrderID)
	if err = helpers.HandleError(err); err != nil {
		return helpers.CreateErrorResponse(fiber.StatusInternalServerError, consts.ServerErrorMessage, nil)
	}

	if transactionLog {
		var currentBalance int

		balance, err := g.transSrv.GetBalance(ctx.Context(), "user_id", trns.UserID)
		if err = helpers.HandleError(err); err != nil {
			return helpers.CreateErrorResponse(fiber.StatusInternalServerError, consts.ServerErrorMessage, nil)
		}

		if balance != nil {
			currentBalance, _ = strconv.Atoi(balance.FinalDeposit)
		}

		grossAmount, err := strconv.ParseFloat(body.GrossAmount, 64)
		if err = helpers.HandleError(err); err != nil {
			return helpers.CreateErrorResponse(fiber.StatusInternalServerError, consts.ServerErrorMessage, nil)
		}

		// create deposit log
		err = g.transSrv.CreateDepositLog(ctx.Context(), dtos.LedgerDto{
			UserID:         trns.UserID,
			RefID:          body.ReferenceID,
			Type:           "TOP UP",
			CurrentDeposit: currentBalance,
			ChangeDeposit:  fmt.Sprintf("+ %.2f", grossAmount),
			FinalDeposit:   currentBalance + int(grossAmount),
			Note:           "TOP UP with REF ID : " + body.ReferenceID,
		})

		if err = helpers.HandleError(err); err != nil {
			return helpers.CreateErrorResponse(fiber.StatusInternalServerError, consts.ServerErrorMessage, nil)
		}
	}

	return *appctx.NewResponse().
		WithCode(fiber.StatusOK).
		WithMessage("success")
}

func NewCallbackMidtransImpl(
	service service.UserService,
	transSrv service.TransactionService,
	cfg *config.Config,
) contract.Controller {
	return &callbackMidtransImpl{
		service:  service,
		transSrv: transSrv,
		cfg:      cfg,
	}
}
