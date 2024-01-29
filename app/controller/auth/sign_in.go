package auth

import (
	"fmt"
	"log"

	"github.com/Zainal21/my-ewallet/app/appctx"
	"github.com/Zainal21/my-ewallet/app/consts"
	"github.com/Zainal21/my-ewallet/app/controller/contract"
	"github.com/Zainal21/my-ewallet/app/dtos"
	"github.com/Zainal21/my-ewallet/app/helpers"
	"github.com/Zainal21/my-ewallet/app/repositories"
	"github.com/Zainal21/my-ewallet/app/service"
	"github.com/Zainal21/my-ewallet/app/utils/golvalidator"
	"github.com/Zainal21/my-ewallet/pkg/config"
	"github.com/Zainal21/my-ewallet/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type signInImpl struct {
	service           service.UserService
	personalTokenRepo repositories.PersonalTokenRepository
	cfg               *config.Config
}

// Serve implements contract.Controller.
func (s *signInImpl) Serve(xCtx appctx.Data) appctx.Response {
	ctx := xCtx.FiberCtx
	signInData := dtos.UserSignInRequestDto{
		Email:    ctx.FormValue("email"),
		Password: ctx.FormValue("password"),
	}

	errors := golvalidator.ValidateStructs(signInData, consts.Localization)

	// check current users
	user, err := s.service.GetUserByFieldName(ctx.Context(), "email", signInData.Email)

	if err != nil {
		errors["email"] = append(errors["email"], "You are not registered yet. please register yourself")
	}

	if user != nil {
		isValid := helpers.VerifyPassword(user.Password, signInData.Password)

		if !isValid {
			errors["password"] = append(errors["password"], "Wrong password")
		}
	}

	if len(errors) > 0 {
		response := helpers.NewValidationErrorResponse(consts.ValidationMessage, errors)
		return helpers.CreateErrorResponse(fiber.StatusUnprocessableEntity, response.Message, &response.Errors)
	}

	if err := s.personalTokenRepo.DeleteByUserId(xCtx.FiberCtx.Context(), user.Id); err != nil {
		log.Println(err.Error())
	}

	// Create token
	token, err := s.personalTokenRepo.Create(xCtx.FiberCtx.Context(), &dtos.PersonalAccessTokenDto{
		Name:        user.Name,
		TokenableId: user.Id,
	})

	if err != nil {
		logger.Error(fmt.Sprintf("Error Create Token : %v", err))
		return helpers.CreateErrorResponse(fiber.StatusInternalServerError, consts.ServerErrorMessage, nil)
	}

	return helpers.SuccessResponse("Login Success", dtos.UserSignInResponseDto{
		Id:          user.Id,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Token:       token,
	}, fiber.StatusOK)
}

func NewSignInImpl(
	svc service.UserService,
	pat repositories.PersonalTokenRepository,
	cfg *config.Config,
) contract.Controller {
	return &signInImpl{
		service:           svc,
		personalTokenRepo: pat,
		cfg:               cfg,
	}
}
