package repo

import (
	"database/sql"
	"errors"

	"github.com/neglarken/clickhead/some-ms/internal/model"
)

type ItemRepository struct {
	DB *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{
		DB: db,
	}
}

func (repo *ItemRepository) Create(item *model.Item) (int, error) {
	var id int
	err := repo.DB.QueryRow(
		"insert into items (info, price) values ($1, $2) returning id",
		item.Info,
		item.Price,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *ItemRepository) Edit(item *model.Item) error {
	return repo.DB.QueryRow(
		"update items set info = $1, price = $2 where id = $3 returning id, info, price",
		item.Info,
		item.Price,
		item.Id,
	).Scan(
		&item.Id,
		&item.Info,
		&item.Price,
	)
}

func (repo *ItemRepository) Get() ([]*model.Item, error) {
	items := make([]*model.Item, 0)
	rows, err := repo.DB.Query(
		"select * from items")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := &model.Item{}
		if err := rows.Scan(&item.Id, &item.Info, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, errors.New("not found any items")
	}
	return items, nil
}

func (repo *ItemRepository) Delete(id int) error {
	err := repo.DB.QueryRow(
		"delete from items where id = $1 returning id", id,
	).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("not found any items")
		}
		return err
	}
	return nil
}
