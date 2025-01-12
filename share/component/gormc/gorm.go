package gormc

import (
	nft_proxy "github.com/alphabatem/nft-proxy/services/model"
	"github.com/alphabatem/nft-proxy/share"
	"log"
)

type sqliteService struct {
	Db       *gorm.DB
	DataBase string
}

func NewSqliteService(db string) *sqliteService {
	return &sqliteService{
		DataBase: db,
	}
}

func (ds *sqliteService) GetId() string {
	return share.SQLITE_SVC
}

func (ds *sqliteService) GetDBConnection() *gorm.DB {
	return ds.Db
}
func (ds *sqliteService) Start() (err error) {
	ds.Db, err = gorm.Open(sqlite.Open(ds.DataBase), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Printf("Error in Start: %v", err)
	}

	err = ds.db.AutoMigrate(&nft_proxy.SolanaMedia{})
	if err != nil {
		return err
	}

	return nil
}
