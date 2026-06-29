package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mislavperi/jafa/server/internal/domain/apperr"
	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	requestmodels "github.com/mislavperi/jafa/server/internal/domain/models/request"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Queries *psql.Queries
	Mapper  *mappers.UserMapper
}

func NewAuthService(pool psql.Pool) *AuthService {
	return &AuthService{Queries: psql.New(pool), Mapper: mappers.NewUserMapper()}
}

func (as *AuthService) Login(ctx context.Context, username, password string) (models.User, error) {
	row, err := as.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, apperr.ErrInvalidCredentials
		}
		return models.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(row.Password), []byte(password)); err != nil {
		return models.User{}, apperr.ErrInvalidCredentials
	}
	return as.Mapper.MapToDomain(row), nil
}

// DeleteAccount soft-deletes the user: the row and the user's data stay in the
// database for recovery, but the account can no longer log in and its username
// is freed for re-registration (uniqueness only covers active users).
func (as *AuthService) DeleteAccount(ctx context.Context, userID int64) error {
	rows, err := as.Queries.SoftDeleteUser(ctx, userID)
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperr.ErrUserNotFound
	}
	return nil
}

func (as *AuthService) Register(ctx context.Context, params requestmodels.RegisterRequest) (models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	row, err := as.Queries.CreateUser(ctx, psql.CreateUserParams{
		Username:  params.Username,
		Password:  string(hash),
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.User{}, apperr.ErrUsernameTaken
		}
		return models.User{}, err
	}
	return as.Mapper.MapCreatedToDomain(row), nil
}
