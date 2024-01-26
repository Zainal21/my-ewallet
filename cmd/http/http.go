package http

import (
	"fmt"

	"github.com/Zainal21/my-ewallet/pkg/app"
	"github.com/Zainal21/my-ewallet/pkg/config"
	"github.com/Zainal21/my-ewallet/pkg/logger"
)

func Start() {
	logger.SetJSONFormatter()
	cnf, err := config.LoadAllConfigs()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to load configuration file: %v", err))
	}

	app.InitializeApp(cnf)
	application := app.GetServer()

	if err := application.StartServer(); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}
