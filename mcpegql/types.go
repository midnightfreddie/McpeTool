package mcpegql

import (
	"encoding/base64"

	"github.com/graphql-go/graphql"
)

// DbObject is used for dbKeys results
type DbObject struct {
	Key        []byte `json:"key,omitempty"`
	data       []byte
	StringKey  string `json:"stringKey,omitempty"`
	HexKey     string `json:"hexKey,omitempty"`
	Base64Data string `json:"base64Data,omitempty"`
	Base64Key  string `json:"base64Key,omitempty"`
	SizeBytes  int    `json:"sizeBytes,omitempty"`
}

// Fill is used to convert the raw byte arrays to JSON-friendly data before returning to client
func (o *DbObject) Fill() {
	o.StringKey, o.HexKey = ConvertKey(o.Key)
	o.Base64Data = base64.StdEncoding.EncodeToString(o.data)
	o.Base64Key = base64.StdEncoding.EncodeToString(o.Key)
}

var dbKeyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "DbKey",
		Fields: graphql.Fields{
			"key": &graphql.Field{
				Type: graphql.NewList(graphql.Int),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					key, ok := p.Source.([]byte)
					if ok {
						return key, nil
					}
					return nil, nil
				},
			},
			"hexKey": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					key, ok := p.Source.([]byte)
					if ok {
						_, hexKey := ConvertKey(key)
						return hexKey, nil
					}
					return nil, nil
				},
			},
			"stringKey": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					key, ok := p.Source.([]byte)
					if ok {
						stringKey, _ := ConvertKey(key)
						return stringKey, nil
					}
					return nil, nil
				},
			},
			"base64Key": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					key, ok := p.Source.([]byte)
					if ok {
						return base64.StdEncoding.EncodeToString(key), nil
					}
					return nil, nil
				},
			},
		},
	},
)
var dbObjectType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "DbObject",
		Fields: graphql.Fields{
			"data": &graphql.Field{
				Type: graphql.NewList(graphql.Int),
			},
			"base64Data": &graphql.Field{
				Type: graphql.String,
			},
			"sizeBytes": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
