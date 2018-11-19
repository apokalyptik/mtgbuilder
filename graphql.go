package main

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

type queryArgs struct {
	ID      *queryNullableInts
	Name    *queryNullableStrings
	CID     *queryNullableStrings
	Purpose *queryNullableStrings
	Type    *queryNullableStrings
}

func (q *query) Search(args queryArgs) *[]*v4cardResolver {
	var pool = &queryCardPool{}
	pool.init()
	pool.IncludeByIDs(args.ID)
	pool.IncludeByCID(args.CID)
	pool.IncludeByName(args.Name)
	pool.IncludeByPurpose(args.Purpose)
	pool.IncludeByType(args.Type)
	return pool.resolvers()
}
