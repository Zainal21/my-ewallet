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

type registerImpl struct {
	service           service.UserService
	personalTokenRepo repositories.PersonalTokenRepository
	cfg               *config.Config
}

// Serve implements contract.Controller.
func (s *registerImpl) Serve(xCtx appctx.Data) appctx.Response {
	ctx := xCtx.FiberCtx

	registrationData := dtos.UserRegistrationRequestDto{
		Name:            ctx.FormValue("name"),
		Email:           ctx.FormValue("email"),
		Password:        ctx.FormValue("password"),
		ConfirmPassword: ctx.FormValue("confirmPassword"),
		PhoneNumber:     ctx.FormValue("phoneNumber"),
	}

	errors := golvalidator.ValidateStructs(registrationData, consts.Localization)

	user, err := s.service.GetUserByFieldName(ctx.Context(), "email", registrationData.Email)

	err = helpers.HandleError(err)
	if err != nil {
		return helpers.CreateErrorResponse(fiber.StatusInternalServerError, consts.ServerErrorMessage, nil)
	}

	if user != nil {
		errors["email"] = append(errors["email"], "email already exist")
	}

	if len(errors) > 0 {
		response := helpers.NewValidationErrorResponse(consts.ValidationMessage, errors)
		return helpers.CreateErrorResponse(fiber.StatusUnprocessableEntity, response.Message, &response.Errors)
	}

	err = s.service.CreateUser(ctx.Context(), registrationData)

	err = helpers.HandleError(err)
	if err != nil {
		return helpers.CreateErrorResponse(fiber.StatusInternalServerError, consts.ServerErrorMessage, nil)
	}

	return helpers.SuccessResponse("Registration Success", nil, fiber.StatusCreated)
}

func NewRegisterImpl(
	sv service.UserService,
	repo repositories.PersonalTokenRepository,
	cfg *config.Config,
) contract.Controller {
	return &registerImpl{
		service:           sv,
		personalTokenRepo: repo,
		cfg:               cfg,
	}
}
