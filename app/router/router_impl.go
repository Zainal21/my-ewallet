package router

import (
	"github.com/Zainal21/my-ewallet/app/appctx"
	"github.com/Zainal21/my-ewallet/app/bootstrap"
	"github.com/Zainal21/my-ewallet/app/controller"
	"github.com/Zainal21/my-ewallet/app/controller/auth"
	"github.com/Zainal21/my-ewallet/app/controller/contract"
	"github.com/Zainal21/my-ewallet/app/controller/transaction"
	"github.com/Zainal21/my-ewallet/app/controller/user"
	"github.com/Zainal21/my-ewallet/app/handler"
	"github.com/Zainal21/my-ewallet/app/middleware"
	"github.com/Zainal21/my-ewallet/app/provider"
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
	tokenRepo := repositories.NewPersonalTokenImpl(db, &sanctum.Token{
		Crypto: &cryptoservice.Crypto{},
	}, userRepo)
	transRepo := repositories.NewTransactionRepositoryImpl(db)

	//define services
	userSvc := service.NewUserServiceImpl(userRepo)
	transSvc := service.NewTransactionServiceImpl(userRepo, transRepo)

	//define middleware
	basicMiddleware := middleware.NewAuthMiddleware()
	jwtMiddleware := middleware.NewJwtMiddleware(tokenRepo)
	signatureMiddleware := middleware.NewSignatureMiddleware(rtr.cfg)

	//define provider
	midtransProvider := provider.NewMidtransProvider(rtr.cfg)

	//define controller
	getAllUser := user.NewGetAllUser(userSvc)
	signIn := auth.NewSignInImpl(userSvc, tokenRepo, rtr.cfg)
	signOut := auth.NewSignOutImpl(userSvc, tokenRepo, rtr.cfg)
	registration := auth.NewRegisterImpl(userSvc, tokenRepo, rtr.cfg)
	getBalance := transaction.NewGetBalanceImpl(userSvc, transSvc, rtr.cfg)
	getHistoryDeposit := transaction.NewHistoryDepositImpl(userSvc, transSvc, rtr.cfg)
	transferTransaction := transaction.NewTransferTransactionImpl(userSvc, transSvc, rtr.cfg)
	topUpTransaction := transaction.NewTopUpTransactionImpl(userSvc, transSvc, midtransProvider, rtr.cfg)
	midtransCallback := transaction.NewCallbackMidtransImpl(userSvc, transSvc, rtr.cfg)

	health := controller.NewGetHealth()
	publicApi := rtr.fiber.Group("/api/v1")
	WalletAPi := rtr.fiber.Group("/api/v1/wallet")

	rtr.fiber.Get("/ping", rtr.handle(
		handler.HttpRequest,
		health,
	))
	// authentication routes
	publicApi.Post("/auth/login", rtr.handle(
		handler.HttpRequest,
		signIn,
		// middleware
		signatureMiddleware.SignatureVerify,
	))

	publicApi.Post("/auth/registration", rtr.handle(
		handler.HttpRequest,
		registration,
		signatureMiddleware.SignatureVerify,
	))

	publicApi.Post("/auth/logout", rtr.handle(
		handler.HttpRequest,
		signOut,
		// middleware
		jwtMiddleware.JwtVerify,
		signatureMiddleware.SignatureVerify,
	))
	// get balance
	WalletAPi.Get("/balance", rtr.handle(
		handler.HttpRequest,
		getBalance,
		// middleware
		jwtMiddleware.JwtVerify,
		signatureMiddleware.SignatureVerify,
	))

	// top up deposit
	WalletAPi.Post("/top-up", rtr.handle(
		handler.HttpRequest,
		topUpTransaction,
		// middleware
		jwtMiddleware.JwtVerify,
		signatureMiddleware.SignatureVerify,
	))

	// transaction history
	WalletAPi.Get("/deposit-history", rtr.handle(
		handler.HttpRequest,
		getHistoryDeposit,
		// middleware
		jwtMiddleware.JwtVerify,
		signatureMiddleware.SignatureVerify,
	))
	// get transaction history
	WalletAPi.Get("/transactions-history", rtr.handle(
		handler.HttpRequest,
		getHistoryDeposit,
		// middleware
		jwtMiddleware.JwtVerify,
		signatureMiddleware.SignatureVerify,
	))

	//	transaction/transfer or payment
	WalletAPi.Post("/transfer", rtr.handle(
		handler.HttpRequest,
		transferTransaction,
		// middleware
		jwtMiddleware.JwtVerify,
		signatureMiddleware.SignatureVerify,
	))
	// midtrans callback
	WalletAPi.Post("/callback", rtr.handle(
		handler.HttpRequest,
		midtransCallback,
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
