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

var dbObjectType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "DbObject",
		Fields: graphql.Fields{
			"key": &graphql.Field{
				Type: graphql.NewList(graphql.Int),
			},
			"hexKey": &graphql.Field{
				Type: graphql.String,
			},
			"stringKey": &graphql.Field{
				Type: graphql.String,
			},
			"base64Data": &graphql.Field{
				Type: graphql.String,
			},
			"base64Key": &graphql.Field{
				Type: graphql.String,
			},
			"sizeBytes": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
