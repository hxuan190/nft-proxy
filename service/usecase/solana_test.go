package usecase

import (
	"github.com/joho/godotenv"
	"log"
	"testing"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Print("Error loading .env file")
	}
}

func TestSolanaImageService_FetchMetadata(t *testing.T) {

	pk := solana.MustPublicKeyFromBase58("CJ9AXYbSUPoR95oMvWzgCV3GbG3ZubQjFUpRHN7xqAVb")

	svc := solanaService{}
	svc.Start()

	d, _, err := svc.TokenData(pk)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", d)

}
