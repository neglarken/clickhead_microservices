package service

import (
	"time"

	"github.com/neglarken/clickhead/auth-ms/internal/auth"
	"github.com/neglarken/clickhead/auth-ms/internal/hasher"
	"github.com/neglarken/clickhead/auth-ms/internal/model"
	"github.com/neglarken/clickhead/auth-ms/internal/repo"
	pb "github.com/neglarken/clickhead/auth-ms/protobuf"
)

type User interface {
	SignUp(u *pb.AuthRequest) error
	SignIn(u *pb.AuthRequest) (*pb.TokenResponse, error)
	SetSession(user_id int, s *model.Session) error
}

type Service struct {
	User User
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func NewService(repo repo.Repository, hasher *hasher.SHA1Hasher, manager *auth.Manager, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{
		User: NewUserService(*repo.User, *repo.Session, hasher, manager, accessTTL, refreshTTL),
	}
}
