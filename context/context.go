package context

import (
	"context"

	"github.com/Hustle299/Project-0/models"
)

type privateKey string

const (
	userKey privateKey = "user"
)

// function that nhan 1 context da co va 1 user sau do tra ve context voi user la value moi
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// function de tim 1 user voi context da co
func User(ctx context.Context) *models.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}
