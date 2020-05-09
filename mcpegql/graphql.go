package mcpegql

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/midnightfreddie/McpeTool/world"
)

// Handler wrapper to allow adding headers to all responses
// concept yoinked from http://echorand.me/dissecting-golangs-handlerfunc-handle-and-defaultservemux.html
func setHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set Origin headers for CORS
		// yoinked from http://stackoverflow.com/questions/12830095/setting-http-headers-in-golang Matt Bucci's answer
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		// Since we're dynamically setting origin, don't let it get cached
		w.Header().Set("Vary", "Origin")
		handler.ServeHTTP(w, r)
	})
}

var saveGame *world.World

// SetWorld sets the module-scope world reference
func SetWorld(w *world.World) {
	saveGame = w
}

func Server(w *world.World, bindAddress, bindPort string) error {
	saveGame = w

	Schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    QueryType,
		Mutation: MutationType,
	})
	if err != nil {
		return err
	}

	// create a graphl-go HTTP handler
	graphQlHandler := handler.New(&handler.Config{
		Schema: &Schema,
		Pretty: false,
		// GraphiQL provides simple web browser query interface pulled from Internet
		GraphiQL: false,
		// Playground provides fancier web browser query interface pulled from Internet
		Playground: true,
	})

	http.Handle("/", setHeaders(graphQlHandler))
	log.Fatal(http.ListenAndServe(bindAddress+":"+bindPort, nil))
	return nil
}

// Query takes a GraphQL query and returns JSON
func Query(w *world.World, query string) (string, error) {
	saveGame = w

	Schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    QueryType,
		Mutation: MutationType,
	})
	if err != nil {
		return "", err
	}
	result := graphql.Do(graphql.Params{
		Schema:        Schema,
		RequestString: query,
	})
	out, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(out[:]), nil
}
