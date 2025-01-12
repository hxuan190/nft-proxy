package main

import (
	"context"
	"fmt"
	"github.com/alphabatem/nft-proxy/service"
	"github.com/alphabatem/nft-proxy/share"
	"github.com/alphabatem/nft-proxy/share/component/ginc"
	"github.com/alphabatem/nft-proxy/share/component/gormc"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	// babilu-online/common/context maybe unavailable
	//"github.com/babilu-online/common/context"
	"log"
)

func main() {

	// Load environment variables from .env file when the server starts
	cfg := share.NewEnvConfig()
	cfg.InitConfig()

	port, err := strconv.Atoi(cfg.GetHTTPPort())
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)
	// Run services in a separate goroutine
	go func() {
		defer wg.Done()

		ginService := ginc.NewHttpService(port)
		if err := ginService.Configure(ctx); err != nil {
			log.Fatalf("Failed to configure HTTP service: %v", err)
		}

		router := ginService.GetGin()
		v1 := router.Group("/v1")

		dbEnv := cfg.GetDB()
		// Initialize the database
		gormService := gormc.NewSqliteService(dbEnv)

		sctx := share.NewServiceContext(cfg, v1, gormService)

		//// @Summary Ping liquify services
		//// @Accept  json
		//// @Produce json
		//// @Router /ping [get]
		router.Get("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, share.ResponseData("pong"))
		})

		service.SetUpService(route, sctx)

		if err := ginService.Start(); err != nil {
			log.Fatalf("Failed to start HTTP service: %v", err)
		}
	}()

	// Gracefully shutdown services
	// i don't know func Shutdown() is available in context.Context???
	// i hope it available
	<-stop
	log.Println("Shutting down services...")
	cancel()

	// Wait for all goroutines to finish
	wg.Wait()

	log.Println("Services shut down gracefully.")
}
