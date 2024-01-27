package router

import (
	"github.com/Zainal21/my-ewallet/app/appctx"
	"github.com/Zainal21/my-ewallet/app/bootstrap"
	"github.com/Zainal21/my-ewallet/app/controller"
	"github.com/Zainal21/my-ewallet/app/controller/auth"
	"github.com/Zainal21/my-ewallet/app/controller/contract"
	"github.com/Zainal21/my-ewallet/app/controller/user"
	"github.com/Zainal21/my-ewallet/app/handler"
	"github.com/Zainal21/my-ewallet/app/middleware"
	"github.com/Zainal21/my-ewallet/app/repositories"
	"github.com/Zainal21/my-ewallet/app/service"
	cryptoservice "github.com/Zainal21/my-ewallet/app/utils/crypto"
	"github.com/Zainal21/my-ewallet/app/utils/sanctum"
	"github.com/Zainal21/my-ewallet/pkg/config"

	"github.com/gofiber/fiber/v2"
)

type router struct {
	cfg   *config.Config
	fiber fiber.Router
}

func (rtr *router) handle(hfn httpHandlerFunc, svc contract.Controller, mdws ...middleware.MiddlewareFunc) fiber.Handler {
	return func(xCtx *fiber.Ctx) error {

		//check registered middleware functions
		if rm := middleware.FilterFunc(rtr.cfg, xCtx, mdws); rm.Code != fiber.StatusOK {
			// return response base on middleware
			res := *appctx.NewResponse().
				WithCode(rm.Code).
				WithError(rm.Errors).
				WithMessage(rm.Message)
			return rtr.response(xCtx, res)
		}

		//send to controller
		resp := hfn(xCtx, svc, rtr.cfg)
		return rtr.response(xCtx, resp)
	}
}

func (rtr *router) response(fiberCtx *fiber.Ctx, resp appctx.Response) error {
	fiberCtx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return fiberCtx.Status(resp.Code).Send(resp.Byte())
}

func (rtr *router) Route() {
	//init db
	db := bootstrap.RegistryDatabase(rtr.cfg)

	//define repositories
	userRepo := repositories.NewUserRepositoryImpl(db)
	tokenRepo := repositories.NewPersonalToken(db, &sanctum.Token{
		Crypto: &cryptoservice.Crypto{},
	}, userRepo)

	//define services
	userSvc := service.NewUserServiceImpl(userRepo)

	//define middleware
	basicMiddleware := middleware.NewAuthMiddleware()

	//define provider

	//define controller
	getAllUser := user.NewGetAllUser(userSvc)
	signIn := auth.NewSignIn(userSvc, tokenRepo, rtr.cfg)

	health := controller.NewGetHealth()
	publicApi := rtr.fiber.Group("/api/v1")

	rtr.fiber.Get("/ping", rtr.handle(
		handler.HttpRequest,
		health,
	))
	// authentication routes
	publicApi.Post("/auth/login", rtr.handle(
		handler.HttpRequest,
		signIn,
	))

	publicApi.Post("/auth/register", rtr.handle(
		handler.HttpRequest,
		signIn,
	))

	publicApi.Post("/auth/logout", rtr.handle(
		handler.HttpRequest,
		signIn,
	))
	// get balance
	publicApi.Get("/balance", rtr.handle(
		handler.HttpRequest,
		signIn,
	))

	// top up deposit
	publicApi.Post("/top-up-deposit", rtr.handle(
		handler.HttpRequest,
		signIn,
	))

	// transaction history
	publicApi.Get("/transactions", rtr.handle(
		handler.HttpRequest,
		signIn,
	))

	//	transaction/transfer or payment
	publicApi.Get("/transactions", rtr.handle(
		handler.HttpRequest,
		signIn,
	))
	// example routes
	publicApi.Get("/users", rtr.handle(
		handler.HttpRequest,
		getAllUser,
		//middleware
		basicMiddleware.Authenticate,
	))
}

func NewRouter(cfg *config.Config, fiber fiber.Router) Router {
	return &router{
		cfg:   cfg,
		fiber: fiber,
	}
}
