package transaction

import (
	"fmt"
	"strconv"

	"github.com/Zainal21/my-ewallet/app/appctx"
	"github.com/Zainal21/my-ewallet/app/consts"
	"github.com/Zainal21/my-ewallet/app/controller/contract"
	"github.com/Zainal21/my-ewallet/app/dtos"
	"github.com/Zainal21/my-ewallet/app/entity"
	"github.com/Zainal21/my-ewallet/app/helpers"
	"github.com/Zainal21/my-ewallet/app/service"
	"github.com/Zainal21/my-ewallet/app/utils/paginator"
	"github.com/Zainal21/my-ewallet/pkg/config"
	"github.com/Zainal21/my-ewallet/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type historyTransactionImpl struct {
	service  service.UserService
	transSrv service.TransactionService
	cfg      *config.Config
}

// Serve implements contract.Controller.
func (g *historyTransactionImpl) Serve(xCtx appctx.Data) appctx.Response {
	ctx := xCtx.FiberCtx
	user, ok := xCtx.FiberCtx.Locals("user").(*entity.User)

	if !ok {
		return helpers.CreateErrorResponse(fiber.StatusUnauthorized, consts.UnauthorizedErrorMessage, nil)
	}

	currentPage, err := strconv.Atoi(ctx.FormValue("page"))
	if err != nil || currentPage <= 0 {
		currentPage = 1
	}

	transHistory := dtos.TransactionRequestDto{
		UserId:   user.Id,
		Search:   ctx.FormValue("search"),
		DateFrom: ctx.FormValue("dateFrom"),
		DateTo:   ctx.FormValue("dateTo"),
		Page:     currentPage,
	}

	transPaginate, totalCount, err := g.transSrv.GetTransactionHistory(ctx.Context(), transHistory)
	if err != nil {
		logger.Error(fmt.Sprintf("Error Getting User List Data: %v", err))
		return helpers.CreateErrorResponse(fiber.StatusInternalServerError, consts.ServerErrorMessage, nil)
	}

	limit := 10

	if transPaginate == nil || len(*transPaginate) == 0 {
		paginate := paginator.NewLengthAwarePaginator([]string{}, totalCount, limit, transHistory.Page, nil, false).GetStringMap()
		return *helpers.MapPaginationResponseToApiResponse(paginate)
	}

	paginate := paginator.NewLengthAwarePaginator(*transPaginate, totalCount, limit, transHistory.Page, nil, false).GetStringMap()
	return *helpers.MapPaginationResponseToApiResponse(paginate)
}

func NewHistoryTransactionImpl(
	service service.UserService,
	transSrv service.TransactionService,
	cfg *config.Config,
) contract.Controller {
	return &historyTransactionImpl{
		service:  service,
		transSrv: transSrv,
		cfg:      cfg,
	}
}
