package domain

import (
	"github.com/elisalimli/go_graphql_template/graphql/models"
	"github.com/elisalimli/go_graphql_template/postgres"
	"github.com/elisalimli/go_graphql_template/validator"
)

var (
	ErrBadCredentials     = ("email/password combination don't work")
	ErrSomethingWentWrong = ("something went wrong")
)

type Domain struct {
	UsersRepo postgres.UsersRepo
}

func NewDomain(usersRepo postgres.UsersRepo) *Domain {
	return &Domain{UsersRepo: usersRepo}
}

type Ownable interface {
	IsOwner(user *models.User) bool
}

// common graphql error boilerplate
func NewFieldError(err validator.FieldError) *models.AuthResponse {
	return &models.AuthResponse{Ok: false, Errors: []*validator.FieldError{{Message: err.Message, Field: err.Field}}}
}
