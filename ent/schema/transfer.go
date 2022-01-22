package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Transfer holds the schema definition for the Transfer entity.
type Transfer struct {
	ent.Schema
}

// Fields of the Transfer.
func (Transfer) Fields() []ent.Field {
	return []ent.Field{
		field.Int("from_account_id"),
		field.Int("to_account_id"),
		field.Int("amount"),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the Transfer.
func (Transfer) Edges() []ent.Edge {
	return nil
}
