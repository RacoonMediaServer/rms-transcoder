package worker

import "context"

type Task interface {
	ID() string
	Do(ctx context.Context) error
}
