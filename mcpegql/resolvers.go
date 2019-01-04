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
			Type:        graphql.NewList(dbObjectType),
			Description: "Get list of keys in LevelDB. Specifying multiple boolean arguments is invalid",
			Args: graphql.FieldConfigArgument{
				"isChunkKey": &graphql.ArgumentConfig{
					Type:        graphql.Boolean,
					Description: "If true/false, returns only/no chunk keys. Overridden by stringKeysOnly",
				},
				"stringKeysOnly": &graphql.ArgumentConfig{
					Type:        graphql.Boolean,
					Description: "If true, only returns readable keys",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				stringKeysOnly, okString := p.Args["stringKeysOnly"].(bool)
				isChunkKey, okChunk := p.Args["isChunkKey"].(bool)

				keyList, err := saveGame.GetKeys()
				if err != nil {
					return nil, err
				}
				var outData []DbObject
				for i := range keyList {
					thisKey := new(DbObject)
					thisKey.Key = keyList[i]
					thisKey.Fill()
					if okString && stringKeysOnly {
						if thisKey.StringKey != "" {
							outData = append(outData, *thisKey)
						}
					} else if okChunk {
						if isChunkKey == IsChunkKey(thisKey.Key) {
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
