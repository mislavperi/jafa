package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Queries *psql.Queries
	Mapper  *mappers.UserMapper
}

func NewAuthService(queries *psql.Queries) *AuthService {
	return &AuthService{Queries: queries, Mapper: mappers.NewUserMapper()}
}

func (as *AuthService) Login(username, password string) (models.User, error) {
	row, err := as.Queries.GetUserByUsername(context.Background(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrInvalidCredentials
		}
		return models.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(row.Password), []byte(password)); err != nil {
		return models.User{}, ErrInvalidCredentials
	}
	return as.Mapper.MapFromGetByUsername(row), nil
}

func (as *AuthService) Register(params RegisterParams) (models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	row, err := as.Queries.CreateUser(context.Background(), psql.CreateUserParams{
		Username:  params.Username,
		Password:  string(hash),
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.User{}, ErrUsernameTaken
		}
		return models.User{}, err
	}
	return as.Mapper.MapFromCreate(row), nil
}
