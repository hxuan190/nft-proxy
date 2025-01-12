package share

import "github.com/gagliardetto/solana-go"

// i moved this to share package for better code organization
// and i refactor filename to const.go for better naming
const (
	BASE64_PREFIX  = ";base64,"
	SQLITE_SVC     = "sqlite_svc"
	SOLANA_SVC     = "solana_svc"
	SOLANA_IMG_SVC = "solana_img_svc"
	IMG_SVC        = "img_svc"
	STAT_SVC       = "stat_svc"
	RESIZE_SVC     = "resize_svc"
)

var (
	METAPLEX_CORE = solana.MustPublicKeyFromBase58("CoREENxT6tW1HoK8ypY1SxRMZTcVPm7R94rH4PZNhX7d")
	TOKEN_2022    = solana.MustPublicKeyFromBase58("TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb")

	DeleteResponseOK = `{"status": 200, "error": ""}`
)
