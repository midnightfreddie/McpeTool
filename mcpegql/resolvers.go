package mcpegql

import (
	"github.com/graphql-go/graphql"
)

var queryType = graphql.NewObject(graphql.ObjectConfig{
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
			Type:        graphql.NewList(dbKeyType),
			Description: "Get list of keys in LevelDB. Specifying multiple boolean arguments is invalid",
			Args: graphql.FieldConfigArgument{
				"isChunkKey": &graphql.ArgumentConfig{
					Type:        graphql.Boolean,
					Description: "If true/false, returns only/no chunk keys. Overridden by isStringKey",
				},
				"isStringKey": &graphql.ArgumentConfig{
					Type:        graphql.Boolean,
					Description: "If true/false, returns only/no readable keys",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				isStringKey, okString := p.Args["isStringKey"].(bool)
				isChunkKey, okChunk := p.Args["isChunkKey"].(bool)

				keyList, err := saveGame.GetKeys()
				if err != nil {
					return nil, err
				}
				if okString || okChunk {
					var outKeys [][]byte
					for i := range keyList {
						if okString {
							stringKey, _ := ConvertKey(keyList[i])
							if isStringKey == (stringKey != "") {
								outKeys = append(outKeys, keyList[i])
							}
						} else if okChunk {
							if isChunkKey == IsChunkKey(keyList[i]) {
								outKeys = append(outKeys, keyList[i])
							}
						}
					}
					return outKeys, nil
				} else {
					return keyList, nil
				}
			},
		},
	},
})
