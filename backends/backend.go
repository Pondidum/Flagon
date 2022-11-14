package backends

import "context"

type User struct {
	Key        string
	Attributes map[string]string
}

type Flag struct {
	Key          string
	DefaultValue bool
	Value        bool
}

type Backend interface {
	State(ctx context.Context, flag Flag, user User) (Flag, error)
	Close(ctx context.Context) error
}
