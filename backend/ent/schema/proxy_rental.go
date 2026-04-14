package schema

import (
	"github.com/meitianwang/fast-frame/ent/schema/mixins"
	"github.com/meitianwang/fast-frame/internal/domain"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ProxyRental holds the schema definition for a user's proxy IP rental.
type ProxyRental struct {
	ent.Schema
}

func (ProxyRental) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "proxy_rentals"},
	}
}

func (ProxyRental) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (ProxyRental) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("user_id"),
		field.Int64("node_id"),
		field.Int64("product_id"),
		field.Int64("payment_order_id").
			Optional().
			Nillable(),
		field.String("status").
			MaxLen(20).
			Default(domain.ProxyRentalStatusPendingPayment),
		field.Time("started_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		field.Time("expires_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		// 已用流量字节数，由管理员累计更新
		field.Int64("traffic_used_bytes").
			Default(0),
		// 流量上限字节数，从套餐快照（0=不限）
		field.Int64("traffic_limit_bytes").
			Default(0),
	}
}

func (ProxyRental) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("proxy_rentals").
			Field("user_id").
			Unique().
			Required(),
		edge.From("node", ProxyNode.Type).
			Ref("rentals").
			Field("node_id").
			Unique().
			Required(),
		edge.From("product", ProxyProduct.Type).
			Ref("rentals").
			Field("product_id").
			Unique().
			Required(),
		edge.To("credential", ProxyCredential.Type).
			Unique(),
		edge.To("traffic_logs", ProxyTrafficLog.Type),
	}
}

func (ProxyRental) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "status"),
		index.Fields("node_id", "status"),
		index.Fields("expires_at"),
		index.Fields("payment_order_id"),
	}
}
