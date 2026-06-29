package mappers

import (
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

// MapToDomain converts the row returned by GetUserByUsername (the login path).
// User has two distinct SQLC source rows, so the create path uses
// MapCreatedToDomain; both share the MapToDomain naming used by every mapper.
func (um *UserMapper) MapToDomain(row psql.GetUserByUsernameRow) models.User {
	return models.User{
		Id:        row.ID,
		Username:  row.Username,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Email:     row.Email,
		AvatarUrl: row.AvatarUrl,
		CreatedAt: row.CreatedAt.Time.String(),
		UpdatedAt: row.UpdatedAt.Time.String(),
	}
}

// MapCreatedToDomain converts the row returned by CreateUser (the register path).
func (um *UserMapper) MapCreatedToDomain(row psql.CreateUserRow) models.User {
	return models.User{
		Id:        row.ID,
		Username:  row.Username,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Email:     row.Email,
		AvatarUrl: row.AvatarUrl,
		CreatedAt: row.CreatedAt.Time.String(),
		UpdatedAt: row.UpdatedAt.Time.String(),
	}
}
