package ginc

import (
	"context"
	"fmt"
	"os"
)

type httpService struct {
	BaseURL      string
	Port         int
	defaultImage []byte
	router       *gin.Engine
}

func NewHttpService(
	port int,
) *httpService {
	return &httpService{
		Port: port,
	}
}

func (svc *httpService) Id() string {
	return "http"
}
func (svc *httpService) GetGin() *gin.Engine {
	return svc.router
}

func (svc *httpService) GetDefaultImage() []byte {
	return svc.defaultImage
}
func (svc *httpService) Configure(ctx context.Context) error {
	df, err := os.ReadFile("./docs/failed_image.jpg")
	if err != nil {
		return err
	}
	svc.defaultImage = df
	svc.router = gin.Default()
	svc.router.Use(gin.Recovery())
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("Authorization")
	svc.router.Use(cors.New(config))

	return nil
}

func (svc *httpService) Start() error {
	return svc.router.Run(fmt.Sprintf(":%v", svc.Port))
}
