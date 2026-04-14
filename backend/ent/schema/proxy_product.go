package schema

import (
	"github.com/meitianwang/fast-frame/ent/schema/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/shopspring/decimal"
)

// ProxyProduct holds the schema definition for a proxy rental product/plan.
type ProxyProduct struct {
	ent.Schema
}

func (ProxyProduct) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "proxy_products"},
	}
}

func (ProxyProduct) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (ProxyProduct) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(100).
			NotEmpty(),
		field.String("description").
			SchemaType(map[string]string{dialect.Postgres: "text"}).
			Default(""),
		// 有效天数，如 1、7、30
		field.Int("duration_days").
			Positive(),
		// 流量上限 GB，0 表示不限
		field.Int("traffic_limit_gb").
			Default(0),
		field.Other("price", decimal.Decimal{}).
			SchemaType(map[string]string{dialect.Postgres: "decimal(10,2)"}),
		field.Int("sort_order").
			Default(0),
		field.Bool("is_active").
			Default(true),
	}
}

func (ProxyProduct) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("rentals", ProxyRental.Type),
	}
}

func (ProxyProduct) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("is_active"),
		index.Fields("sort_order"),
	}
}
