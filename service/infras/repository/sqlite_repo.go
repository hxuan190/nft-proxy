package repository

import (
	"context"
	nft_proxy "github.com/alphabatem/nft-proxy/services/model"
	"github.com/alphabatem/nft-proxy/share/repository/sqlrepo"
)

type sqlRepository struct {
	dbCtx sqlrepo.DbContext
}

func NewSqlRepository(dbCtx sqlrepo.DbContext) *sqlRepository {
	return &sqlRepository{
		dbCtx: dbCtx,
	}
}

func (sql *sqlRepository) CountImagesStored(ctx context.Context, out interface{}) (int64, error) {
	var total int64
	err := sql.dbCtx.GetDBConnection().Model(out).Count(&total).Error
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (repo *sqlRepository) FindMediaByMint(dest interface{}, mintKey string) error {
	return repo.dbCtx.First(dest, "mint = ?", mintKey).Error
}

func (repo *sqlRepository) DeleteMediaByMint(mintKey string) error {
	return repo.dbCtx.Delete(&nft_proxy.SolanaMedia{}, "mint = ?", mintKey).Error
}
func (repo *sqlRepository) Save(media *nft_proxy.SolanaMedia) error {
	return repo.dbCtx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "mint"}}, // key column
		UpdateAll: true,
	}).Create(media).Error
}
