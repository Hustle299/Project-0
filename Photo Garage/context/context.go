package context

import (
	"context"

	"github.com/Hustle299/Project-0/models"
)

// the key has 2 component : type and value, if we create a private type even though underlying is a string
// it will be different than using just a string . string type != private key type
type privateKey string

const (
	userKey privateKey = "user"
)

// function that accepts an existing context and a user, and then returns a new context with that user set as a value.
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// function to lookup an user with given context
func User(ctx context.Context) *models.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}

