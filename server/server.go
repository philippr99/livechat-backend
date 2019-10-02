package main

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"symflower/livechat_gqlgen"
	"symflower/livechat_gqlgen/auth"
	"symflower/livechat_gqlgen/routes"
)

// rebuild resolvers: go run github.com/99designs/gqlgen -v
// regenerate generated.go: go generate ./...
//
//go:generate go run github.com/99designs/gqlgen
// Put this in resolver.go over import!
// Needed for `go generate ./..` to work properly

const defaultPort = "8081"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	router.Use(auth.Middleware())
	router.Handle("/", handler.Playground("GraphQL playground", "/graphql"))
	router.Handle(
		"/graphql",
		handler.GraphQL(livechat_gqlgen.NewExecutableSchema(livechat_gqlgen.Config{Resolvers: &livechat_gqlgen.Resolver{}}),
		handler.WebsocketUpgrader(websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		})),
		)

	router.Handle("/signin", routes.CreateAccountEndPoint())

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
