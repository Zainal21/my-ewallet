package auth

import (
	"github.com/Zainal21/my-ewallet/app/appctx"
	"github.com/Zainal21/my-ewallet/app/consts"
	"github.com/Zainal21/my-ewallet/app/controller/contract"
	"github.com/Zainal21/my-ewallet/app/dtos"
	"github.com/Zainal21/my-ewallet/app/helpers"
	"github.com/Zainal21/my-ewallet/app/repositories"
	"github.com/Zainal21/my-ewallet/app/service"
	"github.com/Zainal21/my-ewallet/app/utils/golvalidator"
	"github.com/Zainal21/my-ewallet/pkg/config"
	"github.com/gofiber/fiber/v2"
)

type SignInImpl struct {
	service           service.UserService
	personalTokenRepo repositories.PersonalTokenRepository
	cfg               *config.Config
}

// Serve implements contract.Controller.
func (s *SignInImpl) Serve(xCtx appctx.Data) appctx.Response {
	ctx := xCtx.FiberCtx
	signInData := dtos.UserSignInDto{
		Email:    ctx.FormValue("email"),
		Password: ctx.FormValue("password"),
	}

	errors := golvalidator.ValidateStructs(signInData, consts.Localization)

	if len(errors) > 0 {
		response := helpers.NewValidationErrorResponse(consts.ValidationMessage, errors)
		return helpers.CreateErrorResponse(fiber.StatusUnprocessableEntity, response.Message, &response.Errors)
	}

	return *appctx.NewResponse().
		WithCode(fiber.StatusOK).
		WithData(signInData)
}

func NewSignIn(
	svc service.UserService,
	pat repositories.PersonalTokenRepository,
	cfg *config.Config,
) contract.Controller {
	return &SignInImpl{
		service:           svc,
		personalTokenRepo: pat,
		cfg:               cfg,
	}
}
