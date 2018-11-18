package main

import (
	"sort"
	"strings"
)

var gqlSchema = `
	schema {
		query: Query
		mutation: Mutation
	}
	# The query type, represents all of the entry points into our object graph
	type Query {
		search(
			id:   [Int],        # limit by card MultiverseID
			name: [String],     # limit by name fragment
			cid:  [String],     # limit by ColorIdentity
			purpose: [String],  # limit by card purpose
			type: [String],     # limit by type line fragment
		): [Card]
	}
	# The mutation type, represents all updates we can make to our data
	type Mutation {
		# createReview(episode: Episode!, review: ReviewInput!): Review
	}
	type Card {
		Colors: [String]!
		ColorIdentity: [String]!
		ConvertedManaCost: Int!
		FlavorText: String!
		ID: Int!
		Legalities: [String]!
		ManaCost: String!
		Name: String!
		Printings: [String]!
		Rulings: [Ruling]!
		Text: String!
		Type: String!
		Purpose: [String]!
	}
	type Ruling {
		Date: String!
		Text: String!
	}
`

type query struct{}

func (q *query) Search(args struct {
	ID      *[]*int32
	Name    *[]*string
	CID     *[]*string
	Purpose *[]*string
	Type    *[]*string
}) *[]*v4cardResolver {
	return q.AllInclusive(
		q.IncludeByIDs(args.ID),
		q.IncludeByName(args.Name),
		q.IncludeByCID(args.CID),
		q.IncludeByPurpose(args.Purpose),
		q.IncludeByType(args.Type),
	)
}

func (q *query) IncludeByPurpose(pref *[]*string) []*mtgv4Card {
	if pref == nil {
		return nil
	}
	if len(*pref) < 1 {
		return nil
	}
	var sets = [][]*mtgv4Card{}
	for _, ref := range *pref {
		if ref == nil {
			continue
		}
		set := []*mtgv4Card{}
		p := strings.ToLower(*ref)
		for _, card := range v4.byID {
			var match = false
			for _, c := range card.CustomClassification {
				if p == c {
					match = true
					break
				}
			}
			if match {
				set = append(set, card)
			}
		}
		sets = append(sets, set)
	}
	return q.AllInclusiveCards(sets...)
}

func (*query) IncludeByCID(cid *[]*string) []*mtgv4Card {
	if cid == nil {
		return nil
	}
	var rval = []*mtgv4Card{}
	var sid = []string{}
	for _, ps := range *cid {
		sid = append(sid, strings.ToUpper(*ps))
	}
	sort.Strings(sid)
	for _, card := range v4.byID {
		if len(card.ColorIdentity) > len(sid) {
			continue
		}
		if len(card.ColorIdentity) == 0 {
			rval = append(rval, card)
			continue
		}
		cidOK := true
		for _, cardColor := range card.ColorIdentity {
			found := false
			for _, allowColor := range sid {
				if cardColor == allowColor {
					found = true
				}
			}
			if !found {
				cidOK = false
			}
		}
		if cidOK {
			rval = append(rval, card)
		}
	}
	return rval
}

func (q *query) IncludeByType(types *[]*string) []*mtgv4Card {
	if types == nil {
		return nil
	}
	if len(*types) == 0 {
		return nil
	}
	var sets = [][]*mtgv4Card{}
	for _, typ := range *types {
		if typ == nil {
			continue
		}
		if len(*typ) < 1 {
			continue
		}
		set := []*mtgv4Card{}
		fragment := strings.ToLower(*typ)
		for _, card := range v4.byID {
			if strings.Contains(strings.ToLower(card.Type), fragment) {
				set = append(set, card)
			}
		}
		sets = append(sets, set)
	}
	return q.AllInclusiveCards(sets...)
}

func (q *query) IncludeByName(names *[]*string) []*mtgv4Card {
	if names == nil {
		return nil
	}
	if len(*names) == 0 {
		return nil
	}
	var sets = [][]*mtgv4Card{}
	for _, name := range *names {
		if name == nil {
			continue
		}
		if len(*name) < 1 {
			continue
		}
		set := []*mtgv4Card{}
		fragment := strings.ToLower(*name)
		for name, cards := range v4.byName {
			if strings.Contains(strings.ToLower(name), fragment) {
				for _, card := range cards {
					set = append(set, card)
				}
			}
		}
		sets = append(sets, set)
	}
	return q.AllInclusiveCards(sets...)
}

func (*query) IncludeByIDs(ids *[]*int32) []*mtgv4Card {
	if ids == nil {
		return nil
	}
	if len(*ids) < 1 {
		return nil
	}
	var rval = []*mtgv4Card{}
	for _, id := range *ids {
		if id == nil {
			continue
		}
		rval = append(rval, v4.byID[int64(*id)])
	}
	return rval
}

func (*query) AllInclusiveCards(sets ...[]*mtgv4Card) []*mtgv4Card {
	var rval = []*mtgv4Card{}
	var setCount = 0
	var idCounts = map[string]int{}
	var ids = map[string]*mtgv4Card{}
	for _, set := range sets {
		if set == nil {
			continue
		}
		setCount++
		for _, card := range set {
			if _, ok := idCounts[card.UUID]; !ok {
				idCounts[card.UUID] = 0
				ids[card.UUID] = card
			}
			idCounts[card.UUID]++
		}
	}
	for id, count := range idCounts {
		if count == setCount {
			rval = append(rval, ids[id])
		}
	}
	return rval
}

// For sorting the final output
type v4cardResolverSet []*v4cardResolver

func (s v4cardResolverSet) Len() int           { return len(s) }
func (s v4cardResolverSet) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s v4cardResolverSet) Less(i, j int) bool { return s[i].mtgv4Card.Name < s[j].mtgv4Card.Name }

func (q *query) AllInclusive(sets ...[]*mtgv4Card) *[]*v4cardResolver {
	var rval = []*v4cardResolver{}
	for _, card := range q.AllInclusiveCards(sets...) {
		rval = append(rval, &v4cardResolver{card})
	}
	sort.Sort(v4cardResolverSet(rval))
	return &rval
}
