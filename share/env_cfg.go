package share

import (
	"log"
	"os"
)

// EnvConfig is an interface for environment variables
// all the environment variables are defined here
type env struct{}

type EnvConfig interface {
	// Get the environment variable value by key
	GetRPCUrl() string
	GetHTTPPort() string
	GetDB() string
	InitConfig()
}

func (env *env) GetRPCUrl() string {
	return os.Getenv("RPC_URL")
}
func (env *env) GetHTTPPort() string {
	return os.Getenv("HTTP_PORT")
}
func (env *env) GetDB() string {
	return os.Getenv("DB_DATABASE")
}

func (env *env) InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func NewEnvConfig() EnvConfig {
	return &env{}
}
