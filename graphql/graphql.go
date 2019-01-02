package graphql

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/midnightfreddie/McpeTool/world"
)

// DbObject is used for dbKeys results
type DbObject struct {
	key        []byte
	data       []byte
	StringKey  string `json:"stringKey,omitempty"`
	HexKey     string `json:"hexKey,omitempty"`
	Base64Data string `json:"base64Data,omitempty"`
}

// Fill is used to convert the raw byte arrays to JSON-friendly data before returning to client
func (o *DbObject) Fill() {
	o.StringKey, o.HexKey = ConvertKey(o.key)
	o.Base64Data = base64.StdEncoding.EncodeToString(o.data)
}

var dbObjectType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "DbObject",
		Fields: graphql.Fields{
			"hexKey": &graphql.Field{
				Type: graphql.String,
			},
			"stringKey": &graphql.Field{
				Type: graphql.String,
			},
			"base64Data": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// ConvertKey takes a byte array and returns a string if all characters are printable (else "")  hex-string-encoded versions of key
func ConvertKey(k []byte) (stringKey, hexKey string) {
	allAscii := true
	for i := range k {
		if k[i] < 0x20 || k[i] > 0x7e {
			allAscii = false
		}
	}
	if allAscii {
		stringKey = string(k[:])
	}
	hexKey = hex.EncodeToString(k)
	return
}

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
			// Since we're dynamically setting origin, don't let it get cached
			w.Header().Set("Vary", "Origin")
		}
		handler.ServeHTTP(w, r)
	})
}

func Server(world *world.World, bindAddress, bindPort string) error {

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"helloWorld": &graphql.Field{
				Type:        graphql.String,
				Description: "Static GraphQL sanity test",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "Hello World!", nil
				},
			},
			"dbKeys": &graphql.Field{
				Type:        graphql.NewList(dbObjectType),
				Description: "Get list of keys in LevelDB",
				Args: graphql.FieldConfigArgument{
					"stringKeysOnly": &graphql.ArgumentConfig{
						Type: graphql.Boolean,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					stringKeysOnly, ok := p.Args["stringKeysOnly"].(bool)
					keyList, err := world.GetKeys()
					if err != nil {
						return nil, err
					}
					var outData []DbObject
					for i := range keyList {
						thisKey := new(DbObject)
						thisKey.key = keyList[i]
						thisKey.Fill()
						if ok && stringKeysOnly {
							if thisKey.StringKey != "" {
								outData = append(outData, *thisKey)
							}
						} else {
							outData = append(outData, *thisKey)
						}

					}
					return outData, nil
				},
			},
		},
	})

	Schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
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
