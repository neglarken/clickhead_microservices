package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/neglarken/clickhead/auth-ms/internal/auth"
	"github.com/neglarken/clickhead/auth-ms/internal/hasher"
	"github.com/neglarken/clickhead/auth-ms/internal/model"
	"github.com/neglarken/clickhead/auth-ms/internal/repo"
	pb "github.com/neglarken/clickhead/auth-ms/protobuf"
)

type UserService struct {
	userRepo        repo.UserRepository
	sessionRepo     repo.SessionRepository
	hasher          hasher.PasswordHasher
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUserService(userRepo repo.UserRepository, sessionRepo repo.SessionRepository, hasher *hasher.SHA1Hasher, manager *auth.Manager, accessTTL, refreshTTL time.Duration) *UserService {
	return &UserService{
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
		hasher:          hasher,
		tokenManager:    manager,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

func (s *UserService) SignUp(u *pb.AuthRequest) error {
	passwordHash, err := s.hasher.Hash(u.Password)
	if err != nil {
		return err
	}

	u.Password = passwordHash

	if err := s.userRepo.Create(
		&model.User{
			Login:    u.Login,
			Password: u.Password,
		},
	); err != nil {
		return err
	}
	return nil
}

func (s *UserService) SignIn(u *pb.AuthRequest) (*pb.TokenResponse, error) {
	res := &pb.TokenResponse{}
	passwordHash, err := s.hasher.Hash(u.Password)
	if err != nil {
		return res, err
	}

	u.Password = passwordHash

	user, err := s.userRepo.FindByLogin(u.Login)
	if err != nil {
		return res, err
	}

	if user.Password != u.Password {
		log.Println(user.Password + "\n" + u.Password)
		return res, errors.New("Wrong password or login")
	}

	tokens, err := s.createSession(user.Id)
	if err != nil {
		return res, err
	}
	res.AccessToken, res.RefreshToken = tokens.AccessToken, tokens.RefreshToken

	return res, nil
}

func (s *UserService) SetSession(user_id int, ses *model.Session) error {
	if err := s.sessionRepo.Upsert(user_id, ses); err != nil {
		return err
	}
	return nil
}

func (s *UserService) createSession(userId int) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(fmt.Sprintf("%x", userId), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := model.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.SetSession(userId, &session)

	return res, err
}
