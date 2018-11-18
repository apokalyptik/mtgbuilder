package main

type v4rulingsResolver struct {
	data struct {
		Date string `json:"date"`
		Text string `json:"text"`
	}
}

func (r *v4rulingsResolver) Date() string {
	return r.data.Date
}

func (r *v4rulingsResolver) Text() string {
	return r.data.Text
}

type v4cardResolver struct {
	*mtgv4Card
}

func (*v4cardResolver) helpStringSlice(s []string) []*string {
	if s == nil {
		return nil
	}
	var rval = []*string{}
	for i := range s {
		rval = append(rval, &s[i])
	}
	return rval
}

func (r *v4cardResolver) Purpose() []*string {
	return r.helpStringSlice(r.mtgv4Card.CustomClassification)
}

func (r *v4cardResolver) ColorIdentity() []*string {
	return r.helpStringSlice(r.mtgv4Card.ColorIdentity)
}

func (r *v4cardResolver) Colors() []*string {
	return r.helpStringSlice(r.mtgv4Card.Colors)
}

func (r *v4cardResolver) ConvertedManaCost() int32 {
	f, _ := r.mtgv4Card.ConvertedManaCost.Float64()
	return int32(f)
}

func (r *v4cardResolver) FlavorText() string {
	return r.mtgv4Card.FlavorText
}

func (r *v4cardResolver) ID() int32 {
	return int32(r.MultiverseID)
}

func (r *v4cardResolver) Legalities() []*string {
	var rval = []string{}
	if r.mtgv4Card.Legalities.Brawl != "" {
		rval = append(rval, "Brawl")
	}
	if r.mtgv4Card.Legalities.Commander != "" {
		rval = append(rval, "Commander")
	}
	if r.mtgv4Card.Legalities.Duel != "" {
		rval = append(rval, "Duel")
	}
	if r.mtgv4Card.Legalities.Frontier != "" {
		rval = append(rval, "Frontier")
	}
	if r.mtgv4Card.Legalities.Future != "" {
		rval = append(rval, "Future")
	}
	if r.mtgv4Card.Legalities.Legacy != "" {
		rval = append(rval, "Legacy")
	}
	if r.mtgv4Card.Legalities.Modern != "" {
		rval = append(rval, "Modern")
	}
	if r.mtgv4Card.Legalities.OneV1 != "" {
		rval = append(rval, "OneV1")
	}
	if r.mtgv4Card.Legalities.Pauper != "" {
		rval = append(rval, "Pauper")
	}
	if r.mtgv4Card.Legalities.Penny != "" {
		rval = append(rval, "Penny")
	}
	if r.mtgv4Card.Legalities.Standard != "" {
		rval = append(rval, "Standard")
	}
	if r.mtgv4Card.Legalities.Vintage != "" {
		rval = append(rval, "Vintage")
	}
	return r.helpStringSlice(rval)
}

func (r *v4cardResolver) ManaCost() string {
	return r.mtgv4Card.ManaCost
}

func (r *v4cardResolver) Name() string {
	return r.mtgv4Card.Name
}

func (r *v4cardResolver) Printings() []*string {
	return r.helpStringSlice(r.mtgv4Card.Printings)
}

func (r *v4cardResolver) Rulings() []*v4rulingsResolver {
	var rval = []*v4rulingsResolver{}
	for _, ru := range r.mtgv4Card.Rulings {
		rval = append(rval, &v4rulingsResolver{data: ru})
	}
	return rval
}

func (r *v4cardResolver) Text() string {
	return r.mtgv4Card.Text
}

func (r *v4cardResolver) Type() string {
	return r.mtgv4Card.Type
}
