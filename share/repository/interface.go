package repository

import "context"

type HasTable interface {
	TableName() string
}
type QueryRepository[Entity HasTable, CondDTO any] interface {
	FindById(ctx context.Context, id uuid.UUID, associates ...string) (*Entity, error)
	Count(ctx context.Context, db *gorm.DB) (int64, error)
}

type CommandRepository[Entity HasTable, UpdateDTO HasTable] interface {
	Insert(ctx context.Context, entity Entity) error
	Update(ctx context.Context, id uuid.UUID, dto UpdateDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Repository[Entity, UpdateDTO HasTable, CondDTO any] interface {
	QueryRepository[Entity, CondDTO]
	CommandRepository[Entity, UpdateDTO]
}
