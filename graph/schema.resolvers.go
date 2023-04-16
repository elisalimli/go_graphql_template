package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.29

import (
	"context"
	"fmt"

	"github.com/alisalimli/go_graphql_template/graph/models"
)

// SetMessage is the resolver for the setMessage field.
func (r *mutationResolver) SetMessage(ctx context.Context, message *string) (*string, error) {
	panic(fmt.Errorf("not implemented: SetMessage - setMessage"))
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	panic(fmt.Errorf("not implemented: Users - users"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
