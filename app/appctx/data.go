package appctx

import (
	"github.com/Zainal21/my-ewallet/pkg/config"
	"github.com/gofiber/fiber/v2"
)

type Data struct {
	FiberCtx *fiber.Ctx
	Cfg      *config.Config
}
