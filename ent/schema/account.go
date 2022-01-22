package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Account holds the schema definition for the Account entity.
type Account struct {
	ent.Schema
}

// Fields of the Account.
func (Account) Fields() []ent.Field {
	return []ent.Field{
		// ent はデフォルトではidカラムをautoincrementで自動生成する
		// https://entgo.io/docs/schema-fields/#id-field
		// field.Int("id"),

		// .Optional() を設定するとnullableなカラムになる
		// 今回は全てのカラムが NOT NULL 制約が付与されるため、デフォルトのままで良い
		field.String("owner"),
		field.Int("balance").Positive(),
		field.String("currency"),

		// Default でカラム作成時に付与される値を設定
		// Immutable でカラム作成時のみに値が設定される制約を付与
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the Account.
func (Account) Edges() []ent.Edge {
	return nil
}
