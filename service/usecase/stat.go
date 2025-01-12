package usecase

import (
	"context"
	nft_proxy "github.com/alphabatem/nft-proxy/services/model"
	"github.com/alphabatem/nft-proxy/share"
	"sync/atomic"
)

type StatService interface {
	CountImagesStored(ctx context.Context, out interface{}) (int64, error)
}

type statService struct {
	imageFilesServed uint64
	mediaFilesServed uint64
	requestsServed   uint64

	repo StatService
}

func NewStatService(repo StatService) *statService {
	return &statService{
		repo: repo,
	}
}

type StatHandler interface {
	IncrementImageFileRequests()
	IncrementMediaFileRequests()
	IncrementMediaRequests()
	ServiceStats(ctx context.Context) (map[string]interface{}, error)
}

func (svc statService) Id() string {
	return share.STAT_SVC
}

func (svc *statService) IncrementImageFileRequests() {
	atomic.AddUint64(&svc.imageFilesServed, 1)
}

func (svc *statService) IncrementMediaFileRequests() {
	atomic.AddUint64(&svc.mediaFilesServed, 1)
}

func (svc *statService) IncrementMediaRequests() {
	atomic.AddUint64(&svc.requestsServed, 1)
}

func (svc *statService) ServiceStats(ctx context.Context) (map[string]interface{}, error) {
	imgCount, err := svc.repo.CountImagesStored(ctx, &nft_proxy.SolanaMedia{})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"images_stored":      imgCount,
		"requests_served":    svc.requestsServed,
		"image_files_served": svc.imageFilesServed,
		"media_files_served": svc.mediaFilesServed,
	}, nil
}
