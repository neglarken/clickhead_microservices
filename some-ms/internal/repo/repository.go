package repo

import (
	"database/sql"

	"github.com/neglarken/clickhead/some-ms/internal/model"
)

type Item interface {
	Create(item *model.Item) (int, error)
	Edit(item *model.Item) error
	Get() ([]*model.Item, error)
	Delete(id int) error
}

type Repository struct {
	Item *ItemRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Item: NewItemRepository(db),
	}
}
