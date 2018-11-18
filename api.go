package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var server = &api{
	Router: mux.NewRouter(),
}

type api struct {
	sync.RWMutex
	*mux.Router
}

func (a *api) catchallInit() {
	server.PathPrefix("/").Handler(http.FileServer(http.Dir("htdocs/build/")))
}
