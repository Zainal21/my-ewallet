package auth

import (
	"log"

	"github.com/Zainal21/my-ewallet/app/appctx"
	"github.com/Zainal21/my-ewallet/app/consts"
	"github.com/Zainal21/my-ewallet/app/controller/contract"
	"github.com/Zainal21/my-ewallet/app/entity"
	"github.com/Zainal21/my-ewallet/app/helpers"
	"github.com/Zainal21/my-ewallet/app/repositories"
	"github.com/Zainal21/my-ewallet/app/service"
	"github.com/Zainal21/my-ewallet/pkg/config"
	"github.com/gofiber/fiber/v2"
)

type signOutImpl struct {
	service           service.UserService
	personalTokenRepo repositories.PersonalTokenRepository
	cfg               *config.Config
}

// Serve implements contract.Controller.
func (s *signOutImpl) Serve(xCtx appctx.Data) appctx.Response {
	ctx := xCtx.FiberCtx.Context()
	tokenString := xCtx.FiberCtx.Get("Authorization")

	_, ok := xCtx.FiberCtx.Locals("user").(*entity.User)

	if !ok || tokenString == "" {
		return helpers.CreateErrorResponse(fiber.StatusUnauthorized, consts.UnauthorizedErrorMessage, nil)
	}

	if err := s.personalTokenRepo.Delete(ctx, tokenString); err != nil {
		log.Println("Error or Token not deleted: ", err)
		return helpers.CreateErrorResponse(fiber.StatusUnauthorized, consts.UnauthorizedErrorMessage, nil)
	}

	return helpers.SuccessResponse("Logout Success", nil, fiber.StatusOK)
}

func NewSignOutImpl(
	sv service.UserService,
	repo repositories.PersonalTokenRepository,
	cfg *config.Config,
) contract.Controller {
	return &signOutImpl{
		service:           sv,
		personalTokenRepo: repo,
		cfg:               cfg,
	}
}
