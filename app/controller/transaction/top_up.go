package transaction

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Zainal21/my-ewallet/app/appctx"
	"github.com/Zainal21/my-ewallet/app/consts"
	"github.com/Zainal21/my-ewallet/app/controller/contract"
	"github.com/Zainal21/my-ewallet/app/dtos"
	"github.com/Zainal21/my-ewallet/app/entity"
	"github.com/Zainal21/my-ewallet/app/helpers"
	"github.com/Zainal21/my-ewallet/app/provider"
	"github.com/Zainal21/my-ewallet/app/service"
	"github.com/Zainal21/my-ewallet/app/utils/golvalidator"
	"github.com/Zainal21/my-ewallet/pkg/config"
	"github.com/Zainal21/my-ewallet/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type topUpTransactionImpl struct {
	service      service.UserService
	transSrv     service.TransactionService
	midtransProv provider.Midtrans
	cfg          *config.Config
}

// Serve implements contract.Controller.
func (g *topUpTransactionImpl) Serve(xCtx appctx.Data) appctx.Response {
	ctx := xCtx.FiberCtx
	user, ok := xCtx.FiberCtx.Locals("user").(*entity.User)

	if !ok {
		return helpers.CreateErrorResponse(fiber.StatusUnauthorized, consts.UnauthorizedErrorMessage, nil)
	}

	topupData := dtos.TopUpTransactionDto{
		UserId: user.Id,
		Type:   "TOP UP",
		Amount: ctx.FormValue("amount"),
	}
	topupAmount, _ := strconv.Atoi(topupData.Amount)

	errors := golvalidator.ValidateStructs(topupData, consts.Localization)

	if topupAmount < 150 {
		errors["amount"] = append(errors["amount"], "The amount is not sufficient")
	}

	if len(errors) > 0 {
		response := helpers.NewValidationErrorResponse(consts.ValidationMessage, errors)
		return helpers.CreateErrorResponse(fiber.StatusUnprocessableEntity, response.Message, &response.Errors)
	}

	orderId := helpers.GenerateRefId(4) + "_ORDER_ID"

	result, err := g.midtransProv.CreateCharge(ctx.Context(), map[string]interface{}{
		"amount":   topupAmount,
		"order_id": orderId,
	})

	if err != nil {
		logger.Error(fmt.Sprintf("Error Create Request Midtrans: %v", err))
		return helpers.CreateErrorResponse(fiber.StatusInternalServerError, err.Error(), nil)
	}

	var res dtos.QRISMidtransTransactionDto
	if err = json.Unmarshal([]byte(result), &res); err != nil {
		logger.Error(fmt.Sprintf("Error Unmarshal JSON or Http request Error: %v", err.Error()))
		return helpers.CreateErrorResponse(fiber.StatusInternalServerError, consts.ServerErrorMessage, nil)
	}

	topupData.RefId = res.OrderID

	if res.StatusCode != "201" {
		return helpers.CreateErrorResponse(fiber.StatusBadRequest, res.StatusMessage, nil)
	}

	err = g.transSrv.CreateTransactionLog(ctx.Context(), dtos.TransactionDto{
		OrderID:     orderId,
		UserID:      topupData.UserId,
		RefID:       topupData.RefId,
		Type:        topupData.Type,
		Status:      "process",
		Note:        "TOP UP DEPOSIT with REF ID : " + topupData.RefId,
		GrossAmount: topupData.Amount,
		Piece:       "150",
		Amount:      strconv.Itoa(topupAmount - 150),
	})

	if err = helpers.HandleError(err); err != nil {
		logger.Error(fmt.Sprintf("Error Create transaction log: %v", err.Error()))
		return helpers.CreateErrorResponse(fiber.StatusInternalServerError, consts.ServerErrorMessage, nil)
	}

	return helpers.SuccessResponse("Top Up Success", map[string]interface{}{
		"topUp": topupData,
		"payment": map[string]interface{}{
			"order_id":     res.OrderID,
			"gross_amount": res.GrossAmount,
			"status":       res.FraudStatus,
			"qris_code":    res.Actions[0].URL,
		},
	}, fiber.StatusOK)
}

func NewTopUpTransactionImpl(
	service service.UserService,
	transSrv service.TransactionService,
	midtransProv provider.Midtrans,
	cfg *config.Config,
) contract.Controller {
	return &topUpTransactionImpl{
		service:      service,
		transSrv:     transSrv,
		midtransProv: midtransProv,
		cfg:          cfg,
	}
}
