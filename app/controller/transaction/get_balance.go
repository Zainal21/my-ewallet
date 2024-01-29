package transaction

import (
	"strconv"

	"github.com/Zainal21/my-ewallet/app/appctx"
	"github.com/Zainal21/my-ewallet/app/consts"
	"github.com/Zainal21/my-ewallet/app/controller/contract"
	"github.com/Zainal21/my-ewallet/app/entity"
	"github.com/Zainal21/my-ewallet/app/helpers"
	"github.com/Zainal21/my-ewallet/app/service"
	"github.com/Zainal21/my-ewallet/pkg/config"
	"github.com/gofiber/fiber/v2"
)

type getBalanceImpl struct {
	service  service.UserService
	transSrv service.TransactionService
	cfg      *config.Config
}

// Serve implements contract.Controller.
func (g *getBalanceImpl) Serve(xCtx appctx.Data) appctx.Response {
	ctx := xCtx.FiberCtx
	user, ok := xCtx.FiberCtx.Locals("user").(*entity.User)

	if !ok {
		return helpers.CreateErrorResponse(fiber.StatusUnauthorized, consts.UnauthorizedErrorMessage, nil)
	}

	currentBalance := 0

	balance, err := g.transSrv.GetBalance(ctx.Context(), "user_id", user.Id)

	err = helpers.HandleError(err)
	if err != nil {
		return helpers.CreateErrorResponse(fiber.StatusInternalServerError, consts.ServerErrorMessage, nil)
	}

	if balance != nil {
		currentBalance, _ = strconv.Atoi(balance.FinalDeposit)
	}

	return helpers.SuccessResponse("Get Balance Success", map[string]interface{}{
		"totalBalance": currentBalance,
	}, fiber.StatusOK)
}

func NewGetBalanceImpl(
	service service.UserService,
	transSrv service.TransactionService,
	cfg *config.Config,
) contract.Controller {
	return &getBalanceImpl{
		service:  service,
		transSrv: transSrv,
		cfg:      cfg,
	}
}
