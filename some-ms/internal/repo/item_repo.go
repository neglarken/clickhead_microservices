package repo

import (
	"database/sql"
	"errors"

	pb "github.com/neglarken/clickhead/some-ms/protobuf"
)

type ItemRepository struct {
	DB *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{
		DB: db,
	}
}

func (repo *ItemRepository) Create(item *pb.CreateItemRequest) (int32, error) {
	var id int32
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

func (repo *ItemRepository) Edit(item *pb.UpdateItemRequest) (*pb.Item, error) {
	res := &pb.Item{}
	err := repo.DB.QueryRow(
		"update items set info = $1, price = $2 where id = $3 returning id, info, price",
		item.Info,
		item.Price,
		item.Id,
	).Scan(
		&res.Id,
		&res.Info,
		&res.Price,
	)
	return res, err
}

func (repo *ItemRepository) Get() ([]*pb.Item, error) {
	items := make([]*pb.Item, 0)
	rows, err := repo.DB.Query(
		"select * from items")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := &pb.Item{}
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

func (repo *ItemRepository) Delete(id *pb.DeleteItemRequest) (int32, error) {
	var resId int32
	err := repo.DB.QueryRow(
		"delete from items where id = $1 returning id", id.Id,
	).Scan(&resId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("not found any items")
		}
		return 0, err
	}
	return resId, nil
}
