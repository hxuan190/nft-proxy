package sqlrepo

import (
	"context"
	"errors"
	"github.com/alphabatem/nft-proxy/share"
	"github.com/alphabatem/nft-proxy/share/repository"
)

type BaseQueryRepo[Entity repository.HasTable, CondDTO any] struct {
	dbCtx DbContext
}

func NewBaseQueryRepo[Entity repository.HasTable, CondDTO any](dbCtx DbContext) *BaseQueryRepo[Entity, CondDTO] {
	return &BaseQueryRepo[Entity, CondDTO]{dbCtx: dbCtx}
}

func (repo *BaseQueryRepo[Entity, CondDTO]) FindById(ctx context.Context, id uuid.UUID, associates ...string) (*Entity, error) {
	var data Entity

	db := repo.dbCtx.GetDBConnection()

	if err := db.Table(data.TableName()).
		Where("id = ?", id).
		First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, share.ErrRecordNotFound
		}
		return nil, errors.WithStack(err)
	}

	return &data, nil
}
