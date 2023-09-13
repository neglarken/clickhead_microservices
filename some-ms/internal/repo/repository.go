package repo

import (
	"database/sql"

	pb "github.com/neglarken/clickhead/some-ms/protobuf"
)

type Item interface {
	Create(item *pb.CreateItemRequest) (int32, error)
	Edit(item *pb.UpdateItemRequest) (*pb.Item, error)
	Get() ([]*pb.Item, error)
	Delete(id *pb.DeleteItemRequest) (int32, error)
}

type Repository struct {
	Item *ItemRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Item: NewItemRepository(db),
	}
}
