package domain

import (
	"errors"

	"github.com/elisalimli/go_graphql_template/graphql/models"
	"github.com/elisalimli/go_graphql_template/postgres"
)

var (
	ErrBadCredentials  = errors.New("email/password combination don't work")
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrForbidden       = errors.New("unauthorized")
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

func checkOwnership(o Ownable, user *models.User) bool {
	return o.IsOwner(user)
}
