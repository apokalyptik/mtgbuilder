package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var jsonSize int64
var jsonChecksum []byte
var jsonFile = "json.db"
var jsonURL = "https://mtgjson.com/v4/json/AllSets.json"
var filteredFile = "processed.json"
var runAPI = false
var saveFiltered = true

type mainDB struct {
	api  *apiDB
	data *mainCardDB
	*mux.Router
}

func (m *mainDB) populateAPIDB() {
	db.data.freshen()
	data := m.data.generate()
	m.api.Lock()
	m.api.cardTypes = data
	m.api.Unlock()
}

var db = &mainDB{
	api:    new(apiDB),
	Router: mux.NewRouter(),
}

func init() {
	flag.BoolVar(&runAPI, "api", runAPI, "Run the API server")
	flag.BoolVar(&saveFiltered, "savef", saveFiltered, "Save the filtered by type card lists")
}

func main() {
	flag.Parse()
	db.populateAPIDB()
	go func() {
		t := time.Tick(time.Hour)
		for {
			select {
			case <-t:
				db.populateAPIDB()
			}
		}
	}()
	if saveFiltered {
		if err := db.api.save(); err != nil {
			log.Fatal(err.Error())
		}
	}
	if runAPI {
		db.configureRoutes()
		log.Fatal(http.ListenAndServe(":4410", db))
	}
}
