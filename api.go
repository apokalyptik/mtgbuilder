package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

type apiDB struct {
	sync.RWMutex
	cardTypes map[string][]*cardInfo
}

func (a *apiDB) save() error {
	fp, err := os.OpenFile(filteredFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer fp.Close()
	return json.NewEncoder(fp).Encode(a.cardTypes)
}

func (a *apiDB) filter(cards []*cardInfo, colors []string) []*cardInfo {
	rval := []*cardInfo{}
	log.Println("evaluating", len(cards), "cards")
	for _, card := range cards {
		var match = false
		log.Println(card.Name)
		if len(card.ColorID) > 0 {
			for _, has := range card.ColorID {
				for _, want := range colors {
					if has == want {
						match = true
						break
					}
				}
				if match {
					break
				}
			}
		} else {
			match = true
		}
		if match {
			rval = append(rval, card)
		}
	}
	return rval
}

func (m *mainDB) configureRoutes() {
	m.HandleFunc("/v/0/list/{type}", m.v0list)
	m.HandleFunc("/v/0/list/{type}/cid/{colors}", m.v0list)
}

func (m *mainDB) v0list(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if kind, ok := m.api.cardTypes[vars["type"]]; ok {
		w.Header().Set("Content-Type", "text/json")
		m.api.RLock()
		if cstring, ok := vars["colors"]; ok {
			json.NewEncoder(w).Encode(m.api.filter(kind, strings.Split(strings.ToUpper(cstring), "")))
		} else {
			json.NewEncoder(w).Encode(kind)
		}
		m.api.RUnlock()
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
