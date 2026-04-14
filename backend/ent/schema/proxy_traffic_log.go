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
)

// ProxyTrafficLog records admin-driven traffic usage updates for a rental.
type ProxyTrafficLog struct {
	ent.Schema
}

func (ProxyTrafficLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "proxy_traffic_logs"},
	}
}

func (ProxyTrafficLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (ProxyTrafficLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("rental_id"),
		// 本次新增流量字节数（增量，非累计）
		field.Int64("delta_bytes").
			NonNegative(),
		field.Int64("operator_id"),
		field.String("note").
			SchemaType(map[string]string{dialect.Postgres: "text"}).
			Default(""),
	}
}

func (ProxyTrafficLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("rental", ProxyRental.Type).
			Ref("traffic_logs").
			Field("rental_id").
			Unique().
			Required(),
		edge.From("operator", User.Type).
			Ref("proxy_traffic_logs").
			Field("operator_id").
			Unique().
			Required(),
	}
}

func (ProxyTrafficLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("rental_id"),
		index.Fields("created_at"),
	}
}
