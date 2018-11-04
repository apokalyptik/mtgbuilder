package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type allSetsFileFormat map[string]struct {
	Cards []*cardInfo `json:"cards"`
}

type cardInfo struct {
	Name         string            `json:"name"`
	ColorID      []string          `json:"colorIdentity"`
	Legality     map[string]string `json:"legalities"`
	ManaCost     string            `json:"manaCost"`
	Text         string            `json:"text"`
	Type         string            `json:"type"`
	SubTypes     []string          `json:"subtypes"`
	SuperTypes   []string          `json:"supertypes"`
	Types        []string          `json:"types"`
	MultiverseID int               `json:"multiverseId"`
}

type mainCardDB struct {
	sync.RWMutex
	db map[string]*cardInfo
	rx map[string]*regexp.Regexp
}

func (m *mainCardDB) read() error {
	log.Println("mainCardDB.read()")
	m.Lock()
	defer m.Unlock()
	d := map[string]*cardInfo{}
	file := allSetsFileFormat{}
	buf, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(buf, &file); err != nil {
		return nil
	}
	for _, set := range file {
		for _, card := range set.Cards {
			if c, ok := d[card.Name]; ok {
				if c.MultiverseID < card.MultiverseID {
					d[card.Name] = card
				}
			} else {
				d[card.Name] = card
			}
		}
	}
	m.db = d
	return nil
}

func (m *mainCardDB) generate() map[string][]*cardInfo {
	log.Println("mainCardDB.generate()")
	m.RLock()
	defer m.RUnlock()
	rval := map[string][]*cardInfo{}
	for kind, rx := range m.rx {
		log.Println("generating for", kind)
		log.Println("rx:", rx.String())
		rval[kind] = m.lookup(rx)
	}
	return rval
}

func (m *mainCardDB) lookup(rx *regexp.Regexp) []*cardInfo {
	log.Println("mainCardDB.lookup()")
	rval := []*cardInfo{}
	for _, c := range m.db {
		if rx.MatchString(c.Text) {
			log.Println("\tmatched", c.MultiverseID, c.Name)
			rval = append(rval, c)
		}
	}
	log.Println("found", len(rval), "cards")
	return rval
}

func (m *mainCardDB) downloadJSON() error {
	log.Println("mainCardDB.downloadJSON()")
	rsp, err := http.Get(jsonURL)
	if rsp.Body != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		return err
	}
	if rsp.StatusCode != 200 {
		return fmt.Errorf("Status code %d for %s", rsp.StatusCode, jsonURL)
	}
	fp, err := os.OpenFile(fmt.Sprintf("%s.tmp", jsonFile), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer fp.Close()
	if _, err := io.Copy(fp, rsp.Body); err != nil {
		return err
	}
	return os.Rename(fmt.Sprintf("%s.tmp", jsonFile), jsonFile)
}

func (m *mainCardDB) freshen() {
	log.Println("mainCardDB.freshen()")
	if stat, err := os.Stat(jsonFile); os.IsNotExist(err) {
		if err := m.downloadJSON(); err != nil {
			log.Fatal(err)
		}
	} else if stat.ModTime().Before(time.Now().Add(0 - (24 * time.Hour))) {
		m.downloadJSON()
	}
	if err := m.read(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	db.data = &mainCardDB{
		db: map[string]*cardInfo{},
		rx: map[string]*regexp.Regexp{
			"ramp": regexp.MustCompile(
				strings.Join(
					[]string{
						"(?i)( (add {[^}]+}",
						"play[^.]+additional[^.]+lands?",
						"put[^.]+land[^.]*onto the battlefield",
						"search[^.]+your library for[^.]+(forest|swamp|island|plains|mountain|land)s?) )",
					},
					"|",
				),
			),
			"generic-tribe-support": regexp.MustCompile(
				strings.Join(
					[]string{
						// needs work...
						"(?i)( ((addition|is|instances?|hooses?|shares?|all)[^.]+creature types?|creature[^.]+choice) )",
					},
					"|",
				),
			),
			"counters": regexp.MustCompile(
				strings.Join(
					[]string{
						"(?i)( (proliferates?",
						"((any|another) kind of)? counters? (from|to|on|already)?",
						"Bolster \\d+",
						"Suspend \\d+",
						"{E}",
						// begin types of counters
						"([+-]\\d+/[+-]\\d+",
						"Age",
						"Aim",
						"Arrow",
						"Arrowhead",
						"Awakening",
						"Blaze",
						"Blood",
						"Bounty",
						"Bribery",
						"Brick",
						"Carrion",
						"Charge",
						"CRANK!",
						"Credit",
						"Corpse",
						"Crystal",
						"Cube",
						"Currency",
						"Death",
						"Delay",
						"Depletion",
						"Despair",
						"Devotion",
						"Divinity",
						"Doom",
						"Dream",
						"Echo",
						"Egg",
						"Elixir",
						"Energy",
						"Eon",
						"Experience",
						"Eyeball",
						"Fade",
						"Fate",
						"Feather",
						"Filibuster",
						"Flood",
						"Fungus",
						"Fuse",
						"Gem",
						"Glyph",
						"Golds",
						"Growth",
						"Hatchling",
						"Healing",
						"Hit",
						"Hoofprint",
						"Hour",
						"Hourglass",
						"Hunger",
						"Ice",
						"Incubation",
						"Infection",
						"Intervention",
						"Isolation",
						"Javelin",
						"Ki",
						"Level",
						"Lore",
						"Loyalty",
						"Luck",
						"Magnet",
						"Manifestation",
						"Mannequin",
						"Mask",
						"Matrix",
						"Mine",
						"Mining",
						"Mire",
						"Music",
						"Muster",
						"Net",
						"Omen",
						"Ore",
						"Page",
						"Pain",
						"Paralyzation",
						"Petal",
						"Petrification",
						"Phylactery",
						"Pin",
						"Plague",
						"Poison",
						"Polyp",
						"Pressure",
						"Prey",
						"Pupa",
						"Quest",
						"Rust",
						"Scream",
						"Shell",
						"Shield",
						"Silver",
						"Shred",
						"Sleep",
						"Sleight",
						"Slime",
						"Slumber",
						"Soot",
						"Spore",
						"Storage",
						"Strife",
						"Study",
						"Theft",
						"Tide",
						"Time",
						"Tower",
						"Training",
						"Trap",
						"Treasure",
						"Velocity",
						"Verse",
						"Vitality",
						"Volatile",
						"Wage",
						"Winch",
						"Wind",
						"Wish) counters?) )",
					},
					"|",
				),
			),
		},
	}
}
