package main

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/apokalyptik/mtgbuilder/static"
)

type cConfigs []*cConfig

func (c cConfigs) init() {
	if c != nil && len(c) > 0 {
		for _, entry := range c {
			entry.init()
			classifiers = append(classifiers, entry)
		}
	}
}

type cConfig struct {
	Name string   `json:"name"`
	RxGo bool     `json:"rxgo"`
	Tags []string `json:"tags"`
	Rx   struct {
		Name     []string `json:"name"`
		NoName   []string `json:"noname"`
		Cid      []string `json:"cid"`
		NoCid    []string `json:"nocid"`
		Cmc      []string `json:"cmc"`
		NoCmc    []string `json:"nocmc"`
		Mana     []string `json:"mana"`
		NoMana   []string `json:"nomana"`
		Text     []string `json:"text"`
		NoText   []string `json:"notext"`
		Type     []string `json:"type"`
		NoType   []string `json:"notype"`
		Rarity   []string `json:"rarity"`
		NoRarity []string `json:"norarity"`
	} `json:"rx"`
	yesNameRx   []*regexp.Regexp
	noNameRx    []*regexp.Regexp
	yesCidRx    []*regexp.Regexp
	noCidRx     []*regexp.Regexp
	yesCmcRx    []*regexp.Regexp
	noCmcRx     []*regexp.Regexp
	yesManaRx   []*regexp.Regexp
	noManaRx    []*regexp.Regexp
	yesTextRx   []*regexp.Regexp
	noTextRx    []*regexp.Regexp
	yesTypeRx   []*regexp.Regexp
	noTypeRx    []*regexp.Regexp
	yesRarityRx []*regexp.Regexp
	noRarityRx  []*regexp.Regexp
}

func (c *cConfig) init() {
	c.yesNameRx = c.processRxList(c.Rx.Name)
	c.noNameRx = c.processRxList(c.Rx.NoName)
	c.yesCidRx = c.processRxList(c.Rx.Cid)
	c.noCidRx = c.processRxList(c.Rx.NoCid)
	c.yesCmcRx = c.processRxList(c.Rx.Cmc)
	c.noCmcRx = c.processRxList(c.Rx.NoCmc)
	c.yesManaRx = c.processRxList(c.Rx.Mana)
	c.noManaRx = c.processRxList(c.Rx.NoMana)
	c.yesTextRx = c.processRxList(c.Rx.Text)
	c.noTextRx = c.processRxList(c.Rx.NoText)
	c.yesTypeRx = c.processRxList(c.Rx.Type)
	c.noTypeRx = c.processRxList(c.Rx.NoType)
	c.yesRarityRx = c.processRxList(c.Rx.Type)
	c.noRarityRx = c.processRxList(c.Rx.NoType)
}

func (c *cConfig) processRxList(l []string) []*regexp.Regexp {
	var rval = []*regexp.Regexp{}
	if l != nil && len(l) > 0 {
		for _, r := range l {
			if c.RxGo {
				rval = append(rval, regexp.MustCompile(r))
			} else {
				rval = append(rval, regexp.MustCompilePOSIX(r))
			}
		}
	}
	return rval
}

func (c *cConfig) anyMatch(s string, l []*regexp.Regexp) bool {
	for _, r := range l {
		if r.MatchString(s) {
			return true
		}
	}
	return false
}

func (c *cConfig) allMatch(s string, l []*regexp.Regexp) bool {
	for _, r := range l {
		if !r.MatchString(s) {
			return false
		}
	}
	return true
}

func (c *cConfig) allNotMatch(s string, l []*regexp.Regexp) bool {
	for _, r := range l {
		if r.MatchString(s) {
			return false
		}
	}
	return true
}

func (c *cConfig) matches(s string, l []*regexp.Regexp) int {
	var matches = 0
	for _, r := range l {
		if r.MatchString(s) {
			matches++
		}
	}
	return matches
}

func (c *cConfig) classify(card *mtgv4Card) []string {
	// Bail early
	if m := c.matches(card.Name, c.noNameRx); m > 0 {
		return nil
	}
	if m := c.matches(card.Type, c.noTypeRx); m > 0 {
		return nil
	}
	if m := c.matches(card.ConvertedManaCost.String(), c.noCmcRx); m > 0 {
		return nil
	}
	if m := c.matches(strings.Join(card.ColorIdentity, ""), c.noCidRx); m > 0 {
		return nil
	}
	if m := c.matches(card.Rarity, c.noRarityRx); m > 0 {
		return nil
	}
	// Succeed late
	if m := c.matches(card.Name, c.yesNameRx); m > 0 {
		return c.Tags
	}
	if m := c.matches(card.Text, c.yesTextRx); m > 0 {
		return c.Tags
	}
	if m := c.matches(card.Type, c.yesTypeRx); m > 0 {
		return c.Tags
	}
	if m := c.matches(card.ConvertedManaCost.String(), c.yesCmcRx); m > 0 {
		return c.Tags
	}
	if m := c.matches(strings.Join(card.ColorIdentity, ""), c.yesCidRx); m > 0 {
		return c.Tags
	}
	if m := c.matches(card.Rarity, c.yesRarityRx); m > 0 {
		return c.Tags
	}
	return nil
}

type cardClassifier interface {
	classify(*mtgv4Card) []string
}

var classifiers = []cardClassifier{}

func classify(card *mtgv4Card) {
	tags := map[string]interface{}{}
	card.CustomClassification = []string{}
	for _, c := range classifiers {
		cs := c.classify(card)
		if cs != nil {
			for _, t := range cs {
				tags[t] = nil
			}
		}
	}
	for t := range tags {
		card.CustomClassification = append(card.CustomClassification, t)
	}
}

func init() {
	var cc = cConfigs{}
	if err := json.Unmarshal(static.FileClassifyJSON, &cc); err != nil {
		panic(err)
	}
	cc.init()
}
