package models

import (
	"context"
	"time"
)

type Actor interface {
	Act(ctx context.Context, duration time.Duration) error
}
