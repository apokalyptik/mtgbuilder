package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var v4 = &mtgv4{
	db:           mtgv4db{},
	saveTo:       ".json.db",
	downloadFrom: "https://mtgjson.com/v4/json/AllSets.json",
}

type mtgv4 struct {
	db           mtgv4db
	byName       mtgv4NameLookup
	byID         mtgv4IdLookup
	nameToId     mtgv4NameToId
	saveTo       string
	downloadFrom string
}

type mtgv4NameLookup map[string][]*mtgv4Card
type mtgv4IdLookup map[int64]*mtgv4Card
type mtgv4NameToId map[string]int64

type mtgv4db map[string]mtgv4Edition

type mtgv4Edition struct {
	BaseSetSize int64        `json:"baseSetSize"`
	BoosterV3   interface{}  `json:"boosterV3"`
	Cards       []*mtgv4Card `json:"cards"`
	Code        string       `json:"code"`
	Meta        struct {
		Date    string `json:"date"`
		Version string `json:"version"`
	} `json:"meta"`
	MtgoCode     interface{}   `json:"mtgoCode"`
	Name         string        `json:"name"`
	ReleaseDate  string        `json:"releaseDate"`
	Tokens       []interface{} `json:"tokens"`
	TotalSetSize int64         `json:"totalSetSize"`
	Type         string        `json:"type"`
}

type mtgv4Card struct {
	CustomClassification []string    `json:"custom_classification"` // This is not part of mtgjson output
	Artist               string      `json:"artist"`
	BorderColor          string      `json:"borderColor"`
	ColorIdentity        []string    `json:"colorIdentity"`
	Colors               []string    `json:"colors"`
	ConvertedManaCost    json.Number `json:"convertedManaCost"`
	FlavorText           string      `json:"flavorText"`
	ForeignData          []struct {
		Language string `json:"language"`
		Name     string `json:"name"`
		Text     string `json:"text"`
		Type     string `json:"type"`
	} `json:"foreignData"`
	FrameVersion string `json:"frameVersion"`
	HasFoil      bool   `json:"hasFoil"`
	HasNonFoil   bool   `json:"hasNonFoil"`
	IsReserved   bool   `json:"isReserved"`
	Layout       string `json:"layout"`
	Legalities   struct {
		OneV1     string `json:"1v1"`
		Brawl     string `json:"brawl"`
		Commander string `json:"commander"`
		Duel      string `json:"duel"`
		Frontier  string `json:"frontier"`
		Future    string `json:"future"`
		Legacy    string `json:"legacy"`
		Modern    string `json:"modern"`
		Pauper    string `json:"pauper"`
		Penny     string `json:"penny"`
		Standard  string `json:"standard"`
		Vintage   string `json:"vintage"`
	} `json:"legalities"`
	ManaCost     string   `json:"manaCost"`
	MultiverseID int64    `json:"multiverseId"`
	Name         string   `json:"name"`
	Number       string   `json:"number"`
	OriginalText string   `json:"originalText"`
	OriginalType string   `json:"originalType"`
	Power        string   `json:"power"`
	Printings    []string `json:"printings"`
	Rarity       string   `json:"rarity"`
	Rulings      []struct {
		Date string `json:"date"`
		Text string `json:"text"`
	} `json:"rulings"`
	Subtypes   []string `json:"subtypes"`
	Supertypes []string `json:"supertypes"`
	Text       string   `json:"text"`
	Toughness  string   `json:"toughness"`
	Type       string   `json:"type"`
	Types      []string `json:"types"`
	UUID       string   `json:"uuid"`
}

func (m *mtgv4) get() error {
	log.Println("downloading", m.downloadFrom)
	rsp, err := http.Get(m.downloadFrom)
	if rsp.Body != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		return err
	}
	if rsp.StatusCode != 200 {
		return fmt.Errorf("Status code %d for %s", rsp.StatusCode, m.downloadFrom)
	}
	fp, err := os.OpenFile(fmt.Sprintf("%s.tmp", m.saveTo), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer fp.Close()
	if _, err := io.Copy(fp, rsp.Body); err != nil {
		return err
	}
	return os.Rename(fmt.Sprintf("%s.tmp", m.saveTo), m.saveTo)
}

func (m *mtgv4) load() error {
	log.Println("parsing downloaded json")
	var newDB = mtgv4db{}
	var byName = mtgv4NameLookup{}
	var byID = mtgv4IdLookup{}
	var nameToID = mtgv4NameToId{}
	fp, err := os.Open(m.saveTo)
	if err != nil {
		return err
	}
	if err = json.NewDecoder(fp).Decode(&newDB); err == nil {
		log.Println("indexing data")
		for _, ed := range newDB {
			for _, card := range ed.Cards {
				classify(card)
				if _, ok := byName[card.Name]; ok {
					byName[card.Name] = append(byName[card.Name], card)
				} else {
					byName[card.Name] = []*mtgv4Card{card}
				}
				if _, ok := nameToID[card.Name]; ok {
					if nameToID[card.Name] < card.MultiverseID {
						nameToID[card.Name] = card.MultiverseID
					}
				} else {
					nameToID[card.Name] = card.MultiverseID
				}
				byID[card.MultiverseID] = card
			}
		}
		m.db = newDB
		m.byID = byID
		m.byName = byName
		m.nameToId = nameToID
	}
	return err
}

func (m *mtgv4) update() {
	if stat, err := os.Stat(m.saveTo); os.IsNotExist(err) {
		if err := m.get(); err != nil {
			log.Fatal(err)
		}
	} else if stat.ModTime().Before(time.Now().Add(0 - (24 * time.Hour))) {
		m.get()
	}
	if err := m.load(); err != nil {
		log.Fatal(err)
	}
}

func (m *mtgv4) stayFresh() {
	go func(m *mtgv4) {
		ticker := time.Tick(time.Hour)
		for {
			select {
			case <-ticker:
				m.update()
			}
		}
	}(m)
}

func init() {
	log.Println("v4 data loading")
	v4.update()
	v4.stayFresh()
	log.Println("v4 data loaded")
}
