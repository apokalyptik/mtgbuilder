package main

import (
	"sort"
	"strings"
)

type queryNullableStrings []*string

func (q *queryNullableStrings) Strings() []string {
	if q == nil {
		return nil
	}
	var rval = []string{}
	for _, s := range *q {
		if q == nil {
			continue
		}
		rval = append(rval, *s)
	}
	sort.Strings(rval)
	return rval
}

func (q *queryNullableStrings) ToLowerStrings() []string {
	if q == nil {
		return nil
	}
	var rval = []string{}
	for _, s := range *q {
		if q == nil {
			continue
		}
		rval = append(rval, strings.ToLower(*s))
	}
	sort.Strings(rval)
	return rval
}

func (q *queryNullableStrings) ToUpperStrings() []string {
	if q == nil {
		return nil
	}
	var rval = []string{}
	for _, s := range *q {
		if q == nil {
			continue
		}
		rval = append(rval, strings.ToUpper(*s))
	}
	sort.Strings(rval)
	return rval
}
