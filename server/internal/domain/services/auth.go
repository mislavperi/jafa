package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	requestmodels "github.com/mislavperi/jafa/server/internal/domain/models/request"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	customerrors "github.com/mislavperi/jafa/server/utils/errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Queries *psql.Queries
	mapper  *mappers.UserMapper
}

func NewAuthService(queries *psql.Queries) *AuthService {
	return &AuthService{Queries: queries, mapper: mappers.NewUserMapper()}
}

func (as *AuthService) Login(username, password string) (models.User, error) {
	row, err := as.Queries.GetUserByUsername(context.Background(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, customerrors.ErrInvalidCredentials
		}
		return models.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(row.Password), []byte(password)); err != nil {
		return models.User{}, customerrors.ErrInvalidCredentials
	}
	return as.mapper.MapFromGetByUsername(row), nil
}

// DeleteAccount soft-deletes the user: the row and the user's data stay in the
// database for recovery, but the account can no longer log in and its username
// is freed for re-registration (uniqueness only covers active users).
func (as *AuthService) DeleteAccount(userID int64) error {
	rows, err := as.Queries.SoftDeleteUser(context.Background(), userID)
	if err != nil {
		return err
	}
	if rows == 0 {
		return customerrors.ErrUserNotFound
	}
	return nil
}

func (as *AuthService) Register(params requestmodels.RegisterRequest) (models.User, error) {
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
			return models.User{}, customerrors.ErrUsernameTaken
		}
		return models.User{}, err
	}
	return as.mapper.MapFromCreate(row), nil
}
