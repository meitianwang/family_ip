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

// ProxyCredential holds the generated access credentials for a proxy rental.
type ProxyCredential struct {
	ent.Schema
}

func (ProxyCredential) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "proxy_credentials"},
	}
}

func (ProxyCredential) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (ProxyCredential) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("rental_id").
			Unique(),
		field.String("http_username").
			MaxLen(64).
			NotEmpty(),
		field.String("http_password").
			MaxLen(64).
			NotEmpty(),
		field.String("vless_uuid").
			MaxLen(36).
			NotEmpty(),
		field.String("vless_link").
			SchemaType(map[string]string{dialect.Postgres: "text"}).
			NotEmpty(),
	}
}

func (ProxyCredential) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("rental", ProxyRental.Type).
			Ref("credential").
			Field("rental_id").
			Unique().
			Required(),
	}
}

func (ProxyCredential) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("rental_id"),
	}
}
