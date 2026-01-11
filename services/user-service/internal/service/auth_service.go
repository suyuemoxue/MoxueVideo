package service

import (
	"context"
	"errors"
	"time"

	"example.com/MoxueVideo/user-service/internal/model"
	"example.com/MoxueVideo/user-service/internal/repo"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUsernameTaken      = errors.New("username already taken")
)

type AuthService struct {
	users  repo.UserRepo
	secret []byte
}

func NewAuthService(users repo.UserRepo, jwtSecret string) *AuthService {
	return &AuthService{users: users, secret: []byte(jwtSecret)}
}

type RegisterInput struct {
	Username string
	Password string
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (uint64, error) {
	if _, err := s.users.FindByUsername(ctx, in.Username); err == nil {
		return 0, ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	u := &model.User{
		Username:     in.Username,
		PasswordHash: string(hash),
		DisplayName:  in.Username,
	}

	if err := s.users.Create(ctx, u); err != nil {
		return 0, err
	}
	return u.ID, nil
}

type LoginInput struct {
	Username string
	Password string
}

type LoginResult struct {
	Token string
	User  *model.User
}

func (s *AuthService) Login(ctx context.Context, in LoginInput) (*LoginResult, error) {
	u, err := s.users.FindByUsername(ctx, in.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.signToken(u.ID)
	if err != nil {
		return nil, err
	}
	return &LoginResult{Token: token, User: u}, nil
}

func (s *AuthService) signToken(userID uint64) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   "user",
		Issuer:    "user-service",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		ID:        "",
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iss": claims.Issuer,
		"iat": claims.IssuedAt.Unix(),
		"exp": claims.ExpiresAt.Unix(),
	})

	return t.SignedString(s.secret)
}

func ParseUserIDFromToken(tokenStr string, secret string) (uint64, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidCredentials
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, ErrInvalidCredentials
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, ErrInvalidCredentials
	}
	sub, ok := claims["sub"]
	if !ok {
		return 0, ErrInvalidCredentials
	}
	subFloat, ok := sub.(float64)
	if !ok {
		return 0, ErrInvalidCredentials
	}
	return uint64(subFloat), nil
}
