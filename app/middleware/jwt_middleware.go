package middleware

import (
	"github.com/Zainal21/my-ewallet/app/appctx"
	"github.com/Zainal21/my-ewallet/app/consts"
	"github.com/Zainal21/my-ewallet/app/entity"
	"github.com/Zainal21/my-ewallet/app/helpers"
	"github.com/Zainal21/my-ewallet/app/repositories"
	"github.com/Zainal21/my-ewallet/pkg/config"
	"github.com/gofiber/fiber/v2"
)

type JwtMiddleware struct {
	personalTokenRepository repositories.PersonalTokenRepository
}

func NewJwtMiddleware(repo repositories.PersonalTokenRepository) *JwtMiddleware {
	return &JwtMiddleware{
		personalTokenRepository: repo,
	}
}

func (t *JwtMiddleware) verifyJwt(xCtx *fiber.Ctx, tokenString string) (entity.User, appctx.Response) {
	// check token is empty
	if tokenString == "" {
		return entity.User{}, helpers.CreateErrorResponse(fiber.StatusUnauthorized, consts.UnauthorizedErrorMessage, nil)
	}

	// check token structure
	token, err := helpers.GetBearerTokenFromHeader(tokenString)
	if err != nil {
		return entity.User{}, helpers.CreateErrorResponse(fiber.StatusUnauthorized, consts.UnauthorizedErrorMessage, nil)
	}

	user, err := t.personalTokenRepository.Verify(xCtx.Context(), token)
	if err != nil {
		return entity.User{}, helpers.CreateErrorResponse(fiber.StatusUnauthorized, consts.UnauthorizedErrorMessage, nil)
	}

	xCtx.Locals("user", &entity.User{
		Id:          user.Id,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Password:    user.Password,
	})

	return *user, *appctx.NewResponse().WithCode(fiber.StatusOK)
}

func (t *JwtMiddleware) JwtVerify(xCtx *fiber.Ctx, conf *config.Config) appctx.Response {
	tokenString := xCtx.Get("Authorization")
	_, response := t.verifyJwt(xCtx, tokenString)
	if response.Code != fiber.StatusOK {
		return response
	}
	return response
}
