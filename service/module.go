package service

import (
	"github.com/alphabatem/nft-proxy/service/infras/repository"
	"github.com/alphabatem/nft-proxy/service/infras/transport/ginhttp"
	"github.com/alphabatem/nft-proxy/service/usecase"
	"github.com/alphabatem/nft-proxy/share"
	"net/http"
)

func SetUpService(router *gin.Engine, sctx share.ServiceContext) {
	gorm := sctx.GetDb()
	rpc := sctx.GetCfg().GetRPCUrl()
	sql := repository.NewSqlRepository(gorm)
	v1 := sctx.GetV1()

	resizeHandler := usecase.NewResizeService()
	solanaHandler := usecase.NewSolanaService(rpc)
	solanaImgHandler := usecase.NewSolanaImageService(sql, solanaHandler)
	imgHandler := usecase.NewImageService(solanaImgHandler, resizeHandler)
	statHandler := usecase.NewStatService(sql)

	httpSvc := ginhttp.NewHttpService(imgHandler, resizeHandler, solanaHandler, solanaImgHandler, statHandler)

	router.GET("/stats", httpSvc.GetStat)

	v1.GET("tokens/:id", httpSvc.ShowNFT)
	v1.GET("tokens/:id/image", httpSvc.ShowNFTImage)
	v1.GET("tokens/:id/image.gif", httpSvc.ShowNFTImage)
	v1.GET("tokens/:id/image.png", httpSvc.ShowNFTImage)
	v1.GET("tokens/:id/image.jpg", httpSvc.ShowNFTImage)
	v1.GET("tokens/:id/image.jpeg", httpSvc.ShowNFTImage)
	v1.GET("tokens/:id/media", httpSvc.ShowNFTMedia)

	v1.GET("nfts/:id", httpSvc.ShowNFT)
	v1.GET("nfts/:id/image", httpSvc.ShowNFTImage)
	v1.GET("nfts/:id/image.gif", httpSvc.ShowNFTImage)
	v1.GET("nfts/:id/image.png", httpSvc.ShowNFTImage)
	v1.GET("nfts/:id/image.jpg", httpSvc.ShowNFTImage)
	v1.GET("nfts/:id/image.jpeg", httpSvc.ShowNFTImage)
	v1.GET("nfts/:id/media", httpSvc.ShowNFTMedia)

	router.NoRoute(func(c *gin.Context) {
		data := map[string]string{
			"code":    "PAGE_NOT_FOUND",
			"message": "Page not found",
		}
		c.JSON(http.StatusNotFound, share.ResponseData(data))
	})

}
