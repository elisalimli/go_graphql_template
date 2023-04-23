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
	userRepo := postgres.UsersRepo{DB: initializers.DB, RedisClient: initializers.RedisClient}

	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4000"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(customMiddleware.AuthMiddleware(userRepo))
	// for passing http writer, reader to context
	router.Use(customMiddleware.ContextMiddleware)

	d := domain.NewDomain(userRepo)

	c := graphql.Config{Resolvers: &graphql.Resolver{Domain: d}}

	queryHandler := handler.NewDefaultServer(graphql.NewExecutableSchema(c))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", queryHandler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}

// func uploadProgress(w http.ResponseWriter, r *http.Request) {
// 	mr, err := r.MultipartReader()
// 	if err != nil {
// 		fmt.Fprintln(w, err)
// 		return
// 	}

// 	length := r.ContentLength
// 	progress := float64(0)
// 	lastProgressSent := float64(0)

// 	// Set the response headers for Server-Sent Events
// 	w.Header().Set("Content-Type", "text/event-stream")
// 	w.Header().Set("Cache-Control", "no-cache")
// 	w.Header().Set("Connection", "keep-alive")

// 	// Send an initial progress update to the client
// 	fmt.Fprintf(w, "data: %v\n\n", progress)

// 	for {
// 		var read int64
// 		part, err := mr.NextPart()

// 		if err == io.EOF {
// 			fmt.Printf("\nDone!")
// 			break
// 		}

// 		dst, err := os.OpenFile("a.pdf", os.O_WRONLY|os.O_CREATE, 0666)

// 		if err != nil {
// 			return
// 		}

// 		for {
// 			buffer := make([]byte, 100000)
// 			cBytes, err := part.Read(buffer)
// 			if err == io.EOF {
// 				fmt.Printf("\nLast buffer read!")
// 				break
// 			}

// 			read += int64(cBytes)
// 			if read > 0 {
// 				newProgress := float64(read) / float64(length) * 100

// 				// Send a progress update to the client if it's divisible by 20
// 				if int(newProgress/5) > int(lastProgressSent/5) {
// 					lastProgressSent = newProgress
// 					fmt.Fprintf(w, "data: %v\n\n", newProgress)
// 				}

// 				dst.Write(buffer[0:cBytes])
// 			} else {
// 				break
// 			}
// 		}
// 	}
// }
