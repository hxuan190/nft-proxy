package sqlrepo

import (
	"context"
	"github.com/alphabatem/nft-proxy/share/repository"
	"github.com/pkg/errors"
)

type BaseCommandRepo[Entity, UpdateDTO repository.HasTable] struct {
	dbCtx DbContext
}

func NewBaseCommandRepo[Entity, UpdateDTO repository.HasTable](dbCtx DbContext) *BaseCommandRepo[Entity, UpdateDTO] {
	return &BaseCommandRepo[Entity, UpdateDTO]{dbCtx: dbCtx}
}

func (repo *BaseCommandRepo[Entity, UpdateDTO]) Insert(ctx context.Context, entity Entity) error {
	db := repo.dbCtx.GetDBConnection()

	if err := db.Table(entity.TableName()).Create(&entity).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *BaseCommandRepo[Entity, UpdateDTO]) Update(ctx context.Context, id uuid.UUID, dto UpdateDTO) error {
	db := repo.dbCtx.GetDBConnection()

	if err := db.Table(dto.TableName()).Where("id = ?", id).Updates(dto).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *BaseCommandRepo[Entity, UpdateDTO]) Delete(ctx context.Context, id uuid.UUID) error {
	db := repo.dbCtx.GetDBConnection()

	if err := db.Table(dto.TableName()).Where("id = ?", id).Delete(&dto).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
