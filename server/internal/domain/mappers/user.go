package mappers

import (
	"github.com/mislavperi/jafa/server/internal/domain/models"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (um *UserMapper) MapFromGetByUsername(row psql.GetUserByUsernameRow) models.User {
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

func (um *UserMapper) MapFromCreate(row psql.CreateUserRow) models.User {
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
