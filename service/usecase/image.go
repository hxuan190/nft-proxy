package usecase

import (
	"encoding/base64"
	"errors"
	"fmt"
	nft_proxy "github.com/alphabatem/nft-proxy/services/model"
	"github.com/alphabatem/nft-proxy/share"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type SolanaImageProvider interface {
	Media(key string, skipCache bool) (*nft_proxy.Media, error)
}
type ResizeProvider interface {
	Resize(data []byte, out io.Writer, size int) error
}

type ImageService struct {
	defaultSize  int
	httpMedia    *http.Client
	solSvc       SolanaImageProvider
	resize       ResizeProvider
	exemptImages map[string]struct{} //Some older & core tokens dont have active metadata so we shouldn't update them
}

func NewImageService(solSvc SolanaImageProvider, resize ResizeProvider) *ImageService {
	return &ImageService{
		defaultSize: 720, // Default image size, gifs will be half the size
		httpMedia:   &http.Client{Timeout: 10 * time.Second},
		solSvc:      solSvc,
		resize:      resize,
		exemptImages: map[string]struct{}{
			"2kMpEJCZL8vEDZe7YPLMCS9Y3WKSAMedXBn7xHPvsWvi": {},
			"7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU": {},
			"AFbX8oGjGpmVFywbVouvhQSRmiW2aR1mohfahi4Y2AdB": {},
			"CKfatsPMUf8SkiURsDXs7eK6GWb4Jsd6UDbs7twMCWxo": {},
			"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v": {},
			"Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB": {},
			"mSoLzYCxHdYgdzU16g5QSh3i5K3z3KZK7ytfqcJm7So":  {},
			"So11111111111111111111111111111111111111112":  {},
		},
	}
}

type ImageHandler interface {
	Id() string
	Media(key string, skipCache bool) (*nft_proxy.Media, error)
	ImageFile(c *gin.Context, key string) error
	ClearCache(key string) error
	MediaFile(c *gin.Context, key string) error
	IsSolKey(key string) bool
}

func (svc ImageService) Id() string {
	return share.IMG_SVC
}

func (svc *ImageService) Media(key string, skipCache bool) (*nft_proxy.Media, error) {
	if svc.IsSolKey(key) {
		return svc.solSvc.Media(key, skipCache)
	}

	return nil, errors.New("invalid key")
}

func (svc *ImageService) ImageFile(c *gin.Context, key string) error {
	var err error

	//Fetch the image file to see if its already in the system
	var media *nft_proxy.Media
	if svc.IsSolKey(key) {
		media, err = svc.solSvc.Media(key, false)
		if err != nil {
			return err
		}
	} else {
		return errors.New("unsupported chain")
	}

	cacheName := fmt.Sprintf("./cache/solana/%s.%s", media.Mint, media.ImageType)

	//Check for file or fetch
	ifo, err := os.Stat(cacheName)
	if err != nil || ifo.Size() == 0 { //Missing cached image
		err := svc.fetchMissingImage(media, cacheName)
		if err != nil {
			return err
		}
	}
	//log.Printf("Using cached file: %s", cacheName)

	return svc.writeFile(c, cacheName, media)
}

func (svc *ImageService) ClearCache(key string) error {
	m, err := svc.solSvc.Media(key, false)
	if err != nil {
		return err
	}

	_, exempt := svc.exemptImages[key]
	if exempt {
		//return errors.New("cache recently cleared")
		return nil
	}

	cacheName := fmt.Sprintf("./cache/solana/%s.%s", m.Mint, m.ImageType)
	err = svc.fetchMissingImage(m, cacheName)
	if err != nil {
		return err
	}

	return nil
}

func (svc *ImageService) writeFile(c *gin.Context, path string, media *nft_proxy.Media) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	ifo, err := file.Stat()
	modTime := time.Now()
	if ifo != nil {
		modTime = ifo.ModTime()
	}

	c.Header("Cache-Control", "public, max=age=172800")
	c.Header("Vary", "Accept-Encoding")
	c.Header("Last-Modified", modTime.Format("Mon, 02 Jan 2006 15:04:05 GMT")) //Mon, 03 Jun 2020 11:35:28 GMT
	c.Header("Content-Type", fmt.Sprintf("image/%s", media.ImageType))

	_, err = io.Copy(c.Writer, file)
	if err != nil {
		return err
	}
	return nil
}

func (svc *ImageService) fetchMissingImage(media *nft_proxy.Media, cacheName string) error {
	if media.ImageUri == "" {
		return errors.New("invalid image")
	}

	var err error
	var data []byte
	if strings.Contains(media.ImageUri, share.BASE64_PREFIX) {
		base64String := media.ImageUri
		// Remove the data:image/jpeg;base64, prefix if present
		if v := strings.Index(base64String, share.BASE64_PREFIX); v > -1 {
			base64String = base64String[v+len(share.BASE64_PREFIX):]
		}

		data, err = base64.StdEncoding.DecodeString(base64String)
		if err != nil {
			return err
		}
	} else {
		media.ImageUri = strings.Replace(strings.TrimSpace(media.ImageUri), ".ipfs.nftstorage.link", ".ipfs.w3s.link", 1)

		log.Println("Fetching", media.ImageUri)

		req, err := http.NewRequest("GET", media.ImageUri, nil)
		if err != nil {
			return err
		}

		req.Header.Set("User-Agent", "PostmanRuntime/7.29.2")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Encoding", "gzip,deflate,br")

		resp, err := svc.httpMedia.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return errors.New(resp.Status)
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}

	if len(data) == 0 {
		return errors.New("failed to download image")
	}

	output, err := os.Create(cacheName)
	if err != nil {
		return err
	}
	defer output.Close()

	//log.Printf("Resizing file: %s", cacheName)
	err = svc.resize.Resize(data, output, svc.defaultSize)
	if err != nil {
		return err
	}
	return nil
}

func (svc *ImageService) MediaFile(c *gin.Context, key string) error {
	var media *nft_proxy.Media
	var err error
	if svc.IsSolKey(key) {
		media, err = svc.solSvc.Media(key, false)
		if err != nil {
			return err
		}
	}

	if media.MediaUri == "" {
		return errors.New("no media for mint")
	}

	resp, err := svc.httpMedia.Get(media.MediaUri)
	if err != nil {
		return err
	}

	//Write our data
	c.Header("Cache-Control", "public, max-age=31536000")
	c.Header("Expires", time.Now().AddDate(0, 1, 0).Format(http.TimeFormat))
	c.Header("Content-Type", media.MediaType)
	err = resp.Write(c.Writer)
	if err != nil {
		return err
	}

	return nil
}

func (svc *ImageService) IsSolKey(key string) bool {
	_, err := solana.PublicKeyFromBase58(key)
	return err == nil
}
