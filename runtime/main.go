package main

import (
	"fmt"
	"github.com/alphabatem/nft-proxy/service"
	"github.com/alphabatem/nft-proxy/share"
	"github.com/alphabatem/nft-proxy/share/component/ginc"
	"github.com/alphabatem/nft-proxy/share/component/gormc"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	// babilu-online/common/context maybe unavailable
	"github.com/babilu-online/common/context"
	"log"
)

func main() {

	// Load environment variables from .env file when the server starts
	cfg := share.NewEnvConfig()
	cfg.InitConfig()

	//Because the "github.com/babilu-online/common/context" is not available, we can't use the same code as the original
	// i will build a new context manually to replace the original ones
	//ctx, err := context.NewCtx(
	//	&services.SqliteService{},
	//	&services.statService{},
	//	&services.ResizeService{},
	//	&services.solanaService{},
	//	&services.SolanaImageService{},
	//	&usecase.ImageService{},
	//	&services.HttpService{},
	//)
	// by using graceful shutdown,we can stop the services gracefully instead of force stop
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run services in a separate goroutine
	go func() {
		if err := ctx.Run(); err != nil {
			log.Fatalf("Failed to run services: %v", err)
		}
	}()

	// Gracefully shutdown services
	// i don't know func Shutdown() is available in context.Context???
	// i hope it available
	<-stop
	if err := ctx.Shutdown(); err != nil {
		log.Fatalf("Failed to shutdown services: %v", err)
	}
	if err != nil {
		log.Fatal(err)
		return
	}
}

func startServer(c context.Context) error {
	cfg := share.NewEnvConfig()
	cfg.InitConfig()
	dbEnv := cfg.GetDB()
	port, _ := strconv.Atoi(cfg.GetHTTPPort())

	gin := ginc.NewHttpService(port)
	if err := gin.Configure(c); err != nil {
		return fmt.Errorf("failed to configure gin: %w", err)
	}
	route := gin.GetGin()
	v1 := route.Group("/v1")

	// Initialize the database
	gorm := gormc.NewSqliteService(dbEnv)

	sctx := share.NewServiceContext(cfg, v1, gorm)

	//// @Summary Ping liquify services
	//// @Accept  json
	//// @Produce json
	//// @Router /ping [get]
	route.Get("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, share.ResponseData("pong"))
	})

	service.SetUpService(route, sctx)

	if err := gin.Start(); err != nil {
		return fmt.Errorf("failed to start gin: %w", err)
	}

}
