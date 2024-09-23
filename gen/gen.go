package gen

import "go.mongodb.org/mongo-driver/bson"

func main() {

}

func Query() *query {
	return &query{
		filter: bson.M{},
	}
}

type query struct {
	filter bson.M
}

func (m *query) Where(in bson.M) *query {
	if m.filter == nil {
		m.filter = bson.M{}
	}
	for k, v := range in {
		m.filter[k] = v
	}
	return m
}

func (m *query) Build() bson.M {
	return m.filter
}

type Field[T any] struct {
	Name string
	Bson string
}

func (m Field[T]) Eq(in T) bson.M {
	return bson.M{m.Bson: in}
}
