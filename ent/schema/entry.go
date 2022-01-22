package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Entry holds the schema definition for the Entry entity.
type Entry struct {
	ent.Schema
}

// Fields of the Entry.
func (Entry) Fields() []ent.Field {
	return []ent.Field{
		field.Int("account_id"),

		// Positive で正数のみになる制約を付与
		field.Int("amount").Positive(),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Entry.
func (Entry) Edges() []ent.Edge {
	return []ent.Edge{
		// accountスキーマのentriesを参照
		// 外部キーとして、account_id を公開
		edge.From("accounts", Account.Type).
			Ref("entries").
			Field("account_id").
			Unique().
			Required(),
	}
}

func (Entry) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("account_id"),
	}
}
