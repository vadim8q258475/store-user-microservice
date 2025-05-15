package schema

import (
	"errors"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").
			Unique().
			NotEmpty().
			MaxLen(255).
			Validate(func(s string) error {
				if !strings.Contains(s, "@") {
					return errors.New("email must contain @ symbol")
				}
				return nil
			}),

		field.String("password").
			Sensitive().
			NotEmpty().
			MinLen(8).
			MaxLen(72),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
