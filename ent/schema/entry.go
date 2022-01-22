package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Entry holds the schema definition for the Entry entity.
type Entry struct {
	ent.Schema
}

// Fields of the Entry.
func (Entry) Fields() []ent.Field {
	return []ent.Field{
		field.Int("account_id"),
		// Positive で正数のみなる制約を付与
		field.Int("amount").Positive(),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Entry.
func (Entry) Edges() []ent.Edge {
	return nil
}
