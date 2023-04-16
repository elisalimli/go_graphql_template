package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/gqlerror"

	"github.com/elisalimli/go_graphql_template/validator"
)

func validation(ctx context.Context, v validator.Validation) bool {
	isValid, errors := v.Validate()

	if !isValid {
		for k, e := range errors {
			graphql.AddError(ctx, &gqlerror.Error{
				Message: e,
				Extensions: map[string]interface{}{
					"field": k,
				},
			})
		}
	}

	return isValid
}
