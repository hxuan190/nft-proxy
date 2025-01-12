package ginhttp

import "github.com/alphabatem/nft-proxy/services/usecase"

type GinHandler interface {
}
type httpService struct {
	imageHandler     usecase.ImageHandler
	resizeHandler    usecase.ResizeHandler
	solanaHandler    usecase.SolanaHandler
	solanaNFTHandler usecase.SolanaNFTHandler
	statHandler      usecase.StatHandler
}

func NewHttpService(
	imageHandler usecase.ImageHandler,
	resizeHandler usecase.ResizeHandler,
	solanaHandler usecase.SolanaHandler,
	solanaNFTHandler usecase.SolanaNFTHandler,
	statHandler usecase.StatHandler) *httpService {
	return &httpService{
		imageHandler:     imageHandler,
		resizeHandler:    resizeHandler,
		solanaHandler:    solanaHandler,
		solanaNFTHandler: solanaNFTHandler,
		statHandler:      statHandler,
	}
}
