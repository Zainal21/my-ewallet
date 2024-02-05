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

type transferTransactionImpl struct {
	service  service.UserService
	transSrv service.TransactionService
	cfg      *config.Config
}

// Serve implements contract.Controller.
func (g *transferTransactionImpl) Serve(xCtx appctx.Data) appctx.Response {
	ctx := xCtx.FiberCtx
	user, ok := xCtx.FiberCtx.Locals("user").(*entity.User)

	if !ok {
		return helpers.CreateErrorResponse(fiber.StatusUnauthorized, consts.UnauthorizedErrorMessage, nil)
	}

	transferData := dtos.TransferRequestDto{
		UserId:             user.Id,
		RefId:              helpers.GenerateRefId(8),
		Type:               "TRANSFER",
		Amount:             ctx.FormValue("amount"),
		AccountDestination: ctx.FormValue("account_destination"),
		Bank:               ctx.FormValue("bank"),
	}

	errors := golvalidator.ValidateStructs(transferData, consts.Localization)

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

	transferAmount, _ := strconv.Atoi(transferData.Amount)
	if currentBalance < transferAmount {
		errors["amount"] = append(errors["amount"], "The balance is not sufficient")
	}

	if len(errors) > 0 {
		response := helpers.NewValidationErrorResponse(consts.ValidationMessage, errors)
		return helpers.CreateErrorResponse(fiber.StatusUnprocessableEntity, response.Message, &response.Errors)
	}

	err = g.transSrv.CreateDepositLog(ctx.Context(), dtos.LedgerDto{
		UserID:         transferData.UserId,
		RefID:          transferData.RefId,
		Type:           transferData.Type,
		CurrentDeposit: currentBalance,
		ChangeDeposit:  "- " + transferData.Amount,
		FinalDeposit:   currentBalance - transferAmount,
		Note:           "TRANSFER with REF ID : " + transferData.RefId + " TO " + transferData.AccountDestination,
	})

	if err = helpers.HandleError(err); err != nil {
		return helpers.CreateErrorResponse(fiber.StatusInternalServerError, consts.ServerErrorMessage, nil)
	}

	return helpers.SuccessResponse("Top up Success", transferData, fiber.StatusOK)
}

func NewTransferTransactionImpl(
	service service.UserService,
	transSrv service.TransactionService,
	cfg *config.Config,
) contract.Controller {
	return &transferTransactionImpl{
		service:  service,
		transSrv: transSrv,
		cfg:      cfg,
	}
}
