package main

type queryNullableInts []*int32

func (q *queryNullableInts) Int64s() []int64 {
	if q == nil {
		return nil
	}
	var rval = []int64{}
	for _, i := range *q {
		if i == nil {
			continue
		}
		rval = append(rval, int64(*i))
	}
	return rval
}

func (q *queryNullableInts) Ints() []int {
	if q == nil {
		return nil
	}
	var rval = []int{}
	for _, i := range *q {
		if i == nil {
			continue
		}
		rval = append(rval, int(*i))
	}
	return rval
}
