package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/elisalimli/go_graphql_template/domain"
	"github.com/elisalimli/go_graphql_template/graphql"
	"github.com/elisalimli/go_graphql_template/initializers"
	customMiddleware "github.com/elisalimli/go_graphql_template/middleware"
	"github.com/elisalimli/go_graphql_template/postgres"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}

const defaultPort = "4000"

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	userRepo := postgres.UsersRepo{DB: initializers.DB}

	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(customMiddleware.AuthMiddleware(userRepo))

	d := domain.NewDomain(userRepo)

	c := graphql.Config{Resolvers: &graphql.Resolver{Domain: d}}
	// fix this error :

	// cannot use (func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) literal) (value of type func(ctx context.Context, obj interface{}, next "github.com/elisalimli/go_graphql_template/graphql".Resolver) (interface{}, error)) as func(ctx context.Context, obj interface{}, next "github.com/99designs/gqlgen/graphql".Resolver) (res interface{}, err error) value in assignmentcompilerIncompatibleAssign
	// c.Directives.Auth = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	// 	if 1 > 2 {
	// 		// block calling the next resolver
	// 		return nil, fmt.Errorf("Access denied")
	// 	}

	// 	// or let it pass through
	// 	return next(ctx)
	// }

	queryHandler := handler.NewDefaultServer(graphql.NewExecutableSchema(c))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", queryHandler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}

// srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
// http.Handle("/query", srv)

// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
// log.Fatal(http.ListenAndServe(":"+port, nil))
