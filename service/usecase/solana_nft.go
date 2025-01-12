package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/alphabatem/nft-proxy/services/model"
	"github.com/alphabatem/nft-proxy/services/model/token-metadata"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alphabatem/nft-proxy/share"
)

type SolanaMediaRepository interface {
	FindMediaByMint(dest interface{}, mintKey string) error
	DeleteMediaByMint(mintKey string) error
	Save(media *model.SolanaMedia) error
}

type SolanaProvider interface {
	TokenData(key solana.PublicKey) (*token_metadata.Metadata, uint8, error)
}

type solanaImageService struct {
	sql  SolanaMediaRepository
	http *http.Client
	sol  SolanaProvider
}

func NewSolanaImageService(sql SolanaMediaRepository, sol SolanaProvider) *solanaImageService {
	return &solanaImageService{
		sql:  sql,
		http: &http.Client{Timeout: 5 * time.Second},
		sol:  sol,
	}
}

type SolanaNFTHandler interface {
	Media(key string, skipCache bool) (*model.Media, error)
	RemoveMedia(key string) error
	FetchMetadata(key string) (*model.SolanaMedia, error)
}

func (svc *solanaImageService) Id() string {
	return share.SOLANA_IMG_SVC
}

func (svc *solanaImageService) Media(key string, skipCache bool) (*model.Media, error) {
	var media *model.SolanaMedia
	err := svc.sql.FindMediaByMint(&media, key)
	if err != nil || skipCache {
		log.Printf("FetchMetadata - %s err: %s", key, err)
		media, err = svc.FetchMetadata(key)
		if err != nil {
			return nil, err //Still cant get metadata
		}
	}

	return media.Media(), nil
}

func (svc *solanaImageService) RemoveMedia(key string) error {
	return svc.sql.DeleteMediaByMint(key)
}

func (svc *solanaImageService) FetchMetadata(key string) (*model.SolanaMedia, error) {
	metadata, err := svc._retrieveMetadata(key)
	if err != nil {
		return nil, err
	}

	media, err := svc.cache(key, metadata, "")
	if err != nil {
		return nil, err
	}

	return media, nil
}

func (svc *solanaImageService) _retrieveMetadata(key string) (*model.NFTMetadataSimple, error) {
	pk, err := solana.PublicKeyFromBase58(key)
	if err != nil {
		return nil, err
	}
	tokenData, decimals, err := svc.sol.TokenData(pk)
	if err != nil || tokenData == nil {
		log.Printf("No token data for %s - %s", pk, err)
		return nil, err
	}

	//log.Printf("TokenData retreive (%v): %+v\n", decimals, tokenData)

	switch tokenData.Protocol {
	case token_metadata.PROTOCOL_METAPLEX_CORE:
		return &model.NFTMetadataSimple{
			Image:           tokenData.Data.Uri,
			Decimals:        decimals,
			Name:            strings.Trim(tokenData.Data.Name, "\x00"),
			Symbol:          strings.Trim(tokenData.Data.Symbol, "\x00"),
			UpdateAuthority: tokenData.UpdateAuthority.String(),
		}, nil
	default:
		//Get file meta if possible
		f, err := svc.retrieveFile(tokenData.Data.Uri)
		if f != nil {
			f.Decimals = decimals
			f.UpdateAuthority = tokenData.UpdateAuthority.String()
			return f, nil
		}
		log.Printf("(%s) retrieveFile err: %s", tokenData.Data.Uri, err)
	}

	//No Metadata
	return &model.NFTMetadataSimple{
		Name:            strings.Trim(tokenData.Data.Name, "\x00"),
		Decimals:        decimals,
		Symbol:          strings.Trim(tokenData.Data.Symbol, "\x00"),
		UpdateAuthority: tokenData.UpdateAuthority.String(),
	}, nil
}

func (svc *solanaImageService) retrieveFile(uri string) (*model.NFTMetadataSimple, error) {
	file, err := svc.http.Get(strings.Trim(uri, "\x00")) //Strip crap off urls
	if err != nil {
		return nil, err
	}

	if file.StatusCode != 200 {
		return nil, err
	}

	defer file.Body.Close()
	data, err := io.ReadAll(file.Body)
	if err != nil {
		return nil, err
	}

	var metadata model.NFTMetadataSimple
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

func (svc *solanaImageService) cache(key string, metadata *model.NFTMetadataSimple, localPath string) (*model.SolanaMedia, error) {
	media := model.SolanaMedia{
		Mint:      key,
		LocalPath: localPath,
	}

	//log.Printf("Metadata: %+v\n", metadata)
	if metadata != nil {
		media.Name = metadata.Name
		media.Symbol = metadata.Symbol
		media.ImageUri = metadata.Image
		media.ImageType = svc.guessImageType(metadata)
		media.UpdateAuthority = metadata.UpdateAuthority
		media.MintDecimals = metadata.Decimals

		mediaFile := metadata.AnimationFile()
		if mediaFile != nil {
			media.MediaUri = mediaFile.URL
			mediaFile.Type = "mp4"
			if strings.Contains(mediaFile.Type, "/") {
				media.MediaType = strings.Split(mediaFile.Type, "/")[1]
			}
		}
	}

	if err := svc.sql.Save(&media); err != nil {
		return nil, fmt.Errorf("failed to cache media for mint %s: %w", key, err)
	}

	return &media, nil
}

func (svc *solanaImageService) guessImageType(metadata *model.NFTMetadataSimple) string {
	if metadata == nil {
		return "jpg"
	}

	imageType := ""
	imgFile := metadata.ImageFile()
	if imgFile != nil && strings.Contains(imgFile.Type, "/") {
		imageType = strings.Split(imgFile.Type, "/")[1]
	}
	if imageType == "" {
		parts := strings.Split(metadata.Image, ".")
		lastPart := parts[len(parts)-1]
		if strings.Contains(lastPart, "=") {
			parts := strings.Split(lastPart, "=")
			imageType = parts[len(parts)-1]
		} else {
			imageType = lastPart
		}
	}

	if strings.Contains(imageType, "?") {
		imageType = strings.Split(imageType, "?")[0]
	}

	if !svc.ValidType(imageType) {
		log.Printf("Invalid image type guessed: %s", imageType)
		return "jpg"
	}

	return imageType
}

func (svc *solanaImageService) ValidType(imageType string) bool {
	switch imageType {
	case "png":
		return true
	case "jpg":
		return true
	case "jpeg":
		return true
	case "gif":
		return true
	case "svg":
		return true
	default:
		return false
	}
}
