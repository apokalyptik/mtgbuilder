package main

//go:generate fileb0x .b0x.json

import (
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

func main() {
	schema := graphql.MustParseSchema(gqlSchema, &query{})
	server.Handle("/query", &relay.Handler{Schema: schema})
	server.catchallInit()
	log.Fatal(http.ListenAndServe(":4410", server))
}
