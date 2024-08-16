package api

import (
	"context"

	"github.com/gin-gonic/gin"

	"log"
	"minireipaz/pkg/config"
	"minireipaz/pkg/di"
	"minireipaz/pkg/honeycomb"
	"minireipaz/pkg/interfaces/middlewares"
	"minireipaz/pkg/interfaces/routes"
	"net/http"
)

var (
	app *gin.Engine
)

func init() {
	Init()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}

func RunWebserver() {
	addr := config.GetEnv("BACKEND_ADDR", ":4020")
	err := app.Run(addr)
	if err != nil {
		log.Panicf("ERROR | Starting gin failed, %v", err)
	}
}

func Init() {
	log.Print("---- Init From Init ----")
	config.LoadEnvs(".")
	ctx := context.Background()
	tp, exp, err := honeycomb.SetupHoneyComb(ctx)

	// Handle shutdown to ensure all sub processes are closed correctly and telemetry is exported
	defer func() {
		_ = exp.Shutdown(ctx)
		_ = tp.Shutdown(ctx)
	}()

	if err != nil {
		log.Panicf("ERROR | Failed to initialize OpenTelemetry: %v", err)
	}

	gin.SetMode(gin.DebugMode)
	app = gin.New()

	worflowcontroller, authService, userController, dashboardController, authController := di.InitDependencies()
	middlewares.Register(app, authService)
	routes.Register(app, worflowcontroller, userController, dashboardController, authController)

	RunWebserver()
}

func Dummy() {
}
