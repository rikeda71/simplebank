package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Transfer holds the schema definition for the Transfer entity.
type Transfer struct {
	ent.Schema
}

// Fields of the Transfer.
func (Transfer) Fields() []ent.Field {
	return []ent.Field{
		// edges で定義しているので省略
		field.Int("from_account_id"),
		field.Int("to_account_id"),

		field.Int("amount"),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the Transfer.
func (Transfer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("from_accounts", Account.Type).
			Ref("from_transfers").
			Field("from_account_id").
			Unique().
			Required(),
		edge.From("to_accounts", Account.Type).
			Ref("to_transfers").
			Field("to_account_id").
			Unique().
			Required(),
	}
}

func (Transfer) Indexes() []ent.Index {
	return nil
}
