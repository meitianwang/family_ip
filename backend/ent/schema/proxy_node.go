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

// ProxyNode holds the schema definition for a residential proxy VPS node.
type ProxyNode struct {
	ent.Schema
}

func (ProxyNode) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "proxy_nodes"},
	}
}

func (ProxyNode) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

func (ProxyNode) Fields() []ent.Field {
	return []ent.Field{
		field.String("ip_address").
			MaxLen(45).
			NotEmpty(),
		field.String("country").
			MaxLen(100).
			Default(""),
		field.String("country_code").
			MaxLen(2).
			Default(""),
		field.String("city").
			MaxLen(100).
			Default(""),
		field.String("isp").
			MaxLen(200).
			Default(""),
		field.Int("http_port").
			Default(3128),
		field.Int("vless_port").
			Default(443),
		field.String("vless_network").
			MaxLen(10).
			Default(domain.VlessNetworkTCP),
		field.Bool("vless_tls").
			Default(false),
		field.String("vless_sni").
			MaxLen(255).
			Default(""),
		field.String("vless_ws_path").
			MaxLen(255).
			Default("/"),
		field.JSON("tags", []string{}).
			Default([]string{}),
		field.String("status").
			MaxLen(20).
			Default(domain.ProxyNodeStatusAvailable),
		field.String("description").
			SchemaType(map[string]string{dialect.Postgres: "text"}).
			Default(""),
	}
}

func (ProxyNode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("rentals", ProxyRental.Type),
	}
}

func (ProxyNode) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
		index.Fields("country_code"),
		index.Fields("deleted_at"),
	}
}
