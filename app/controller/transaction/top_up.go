package transaction

import (
	"strconv"

	"github.com/Zainal21/my-ewallet/app/appctx"
	"github.com/Zainal21/my-ewallet/app/consts"
	"github.com/Zainal21/my-ewallet/app/controller/contract"
	"github.com/Zainal21/my-ewallet/app/dtos"
	"github.com/Zainal21/my-ewallet/app/entity"
	"github.com/Zainal21/my-ewallet/app/helpers"
	"github.com/Zainal21/my-ewallet/app/service"
	"github.com/Zainal21/my-ewallet/app/utils/golvalidator"
	"github.com/Zainal21/my-ewallet/pkg/config"
	"github.com/gofiber/fiber/v2"
)

type topUpTransactionImpl struct {
	service  service.UserService
	transSrv service.TransactionService
	cfg      *config.Config
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
		RefId:  helpers.GenerateRefId(8),
		Type:   "TOP UP",
		Amount: ctx.FormValue("amount"),
	}

	errors := golvalidator.ValidateStructs(topupData, consts.Localization)

	if len(errors) > 0 {
		response := helpers.NewValidationErrorResponse(consts.ValidationMessage, errors)
		return helpers.CreateErrorResponse(fiber.StatusUnprocessableEntity, response.Message, &response.Errors)
	}

	// get current balance
	var currentBalance int = 0

	balance, err := g.transSrv.GetBalance(ctx.Context(), "user_id", user.Id)

	err = helpers.HandleError(err)
	if err != nil {
		return helpers.CreateErrorResponse(fiber.StatusInternalServerError, consts.ServerErrorMessage, nil)
	}

	if balance != nil {
		currentBalance, _ = strconv.Atoi(balance.FinalDeposit)
	}

	topupAmount, _ := strconv.Atoi(topupData.Amount)

	err = g.transSrv.CreateTransaction(ctx.Context(), dtos.LedgerDto{
		UserID:         topupData.UserId,
		RefID:          topupData.RefId,
		Type:           topupData.Type,
		CurrentDeposit: currentBalance,
		ChangeDeposit:  "+ " + topupData.Amount,
		FinalDeposit:   currentBalance + topupAmount,
		Note:           "TOP UP DEPOSIT with REF ID : " + topupData.RefId,
	})

	err = helpers.HandleError(err)
	if err != nil {
		return helpers.CreateErrorResponse(fiber.StatusInternalServerError, consts.ServerErrorMessage, nil)
	}

	return helpers.SuccessResponse("TopUp Success", topupData, fiber.StatusOK)
}

func NewTopUpTransactionImpl(
	service service.UserService,
	transSrv service.TransactionService,
	cfg *config.Config,
) contract.Controller {
	return &topUpTransactionImpl{
		service:  service,
		transSrv: transSrv,
		cfg:      cfg,
	}
}