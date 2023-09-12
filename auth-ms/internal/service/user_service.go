package service

import (
	"fmt"
	"time"

	"github.com/neglarken/clickhead/auth-ms/internal/auth"
	"github.com/neglarken/clickhead/auth-ms/internal/hasher"
	"github.com/neglarken/clickhead/auth-ms/internal/model"
	"github.com/neglarken/clickhead/auth-ms/internal/repo"
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

func (s *UserService) SignUp(u *model.User) error {
	passwordHash, err := s.hasher.Hash(u.UnencryptedPassword)
	if err != nil {
		return err
	}

	u.Password = passwordHash

	if err := s.userRepo.Create(u); err != nil {
		return err
	}
	return nil
}

func (s *UserService) SignIn(u *model.User) (Tokens, error) {
	passwordHash, err := s.hasher.Hash(u.UnencryptedPassword)
	if err != nil {
		return Tokens{}, err
	}

	u.Password = passwordHash

	user, err := s.userRepo.FindByLogin(u.Login)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(user.Id)
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
