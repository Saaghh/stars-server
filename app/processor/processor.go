package processor

import "context"

type store interface {
	WithTx(ctx context.Context, f func(context.Context) error) error
	user
	game
}

type Processor struct {
	db store
}

func New(db store) *Processor {
	return &Processor{
		db: db,
	}
}
