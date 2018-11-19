package main

import (
	"sort"
	"strings"
)

// For sorting the final output
type v4cardResolverSet []*v4cardResolver

func (s v4cardResolverSet) Len() int           { return len(s) }
func (s v4cardResolverSet) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s v4cardResolverSet) Less(i, j int) bool { return s[i].mtgv4Card.Name < s[j].mtgv4Card.Name }

type queryCardPool struct {
	cards []*mtgv4Card
}

func (q *queryCardPool) init() {
	q.cards = []*mtgv4Card{}
	for _, card := range v4.byID {
		q.cards = append(q.cards, card)
	}
}

func (q *queryCardPool) IncludeByPurpose(purps *queryNullableStrings) {
	if purps == nil {
		return
	}
	if len(*purps) < 1 {
		return
	}

	set := q.cards[:0]
	purposes := purps.ToLowerStrings()
	notInPurposes := len(purposes)
	for _, card := range q.cards {
		found := 0
		for _, c := range card.CustomClassification {
			for _, purp := range purposes {
				if c == purp {
					found++
					break
				}
			}
		}
		if found == notInPurposes {
			set = append(set, card)
		}
	}
	q.cards = set
}

func (q *queryCardPool) IncludeByCID(cids *queryNullableStrings) {
	if cids == nil {
		return
	}

	set := q.cards[:0]
	cid := cids.ToUpperStrings()

	for _, card := range q.cards {
		// if the card has more colors than we want then it's obviously out
		if len(card.ColorIdentity) > len(cid) {
			continue
		}
		// if a card has no color then it's artifact or land and matches the
		// color id of the
		if len(card.ColorIdentity) == 0 {
			set = append(set, card)
			continue
		}

		l := len(card.ColorIdentity)
		for _, c := range card.ColorIdentity {
			for _, cc := range cid {
				if cc == c {
					l--
					break
				}
			}
		}
		if l > 0 {
			continue
		}
		set = append(set, card)
	}
	q.cards = set
}

func (q *queryCardPool) IncludeByType(types *queryNullableStrings) {
	if types == nil {
		return
	}
	if len(*types) == 0 {
		return
	}
	set := q.cards[:0]
	typeSet := types.ToLowerStrings()
	for _, card := range q.cards {
		cType := strings.ToLower(card.Type)
		matches := 0
		for _, t := range typeSet {
			if strings.Contains(cType, t) {
				matches++
			}
		}
		if matches < len(typeSet) {
			continue
		}
		set = append(set, card)
	}
	q.cards = set
}

func (q *queryCardPool) IncludeByName(names *queryNullableStrings) {
	if names == nil {
		return
	}
	if len(*names) == 0 {
		return
	}
	set := q.cards[:0]
	nameSet := names.ToLowerStrings()
	for _, card := range q.cards {
		cName := strings.ToLower(card.Name)
		mismatch := 0
		for _, sName := range nameSet {
			if strings.Contains(cName, sName) {
				continue
			}
			mismatch++
			break
		}
		if mismatch > 0 {
			continue
		}
		set = append(set, card)
	}
	q.cards = set
}

func (q *queryCardPool) IncludeByIDs(ids *queryNullableInts) {
	if ids == nil {
		return
	}
	if len(*ids) < 1 {
		return
	}
	idSet := ids.Int64s()
	set := q.cards[:0]
	for _, card := range q.cards {
		for _, id := range idSet {
			if id == card.MultiverseID {
				set = append(set, card)
			}
		}
	}
	q.cards = set
}

func (q *queryCardPool) resolvers() *[]*v4cardResolver {
	var rval = []*v4cardResolver{}
	for _, card := range q.cards {
		rval = append(rval, &v4cardResolver{card})
	}
	sort.Sort(v4cardResolverSet(rval))
	return &rval
}
