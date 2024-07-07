package handlers

type ProcessorInterface interface {
	user
	game
}

type Handlers struct {
	proc ProcessorInterface
}

func NewHandler(proc ProcessorInterface) *Handlers {
	return &Handlers{
		proc: proc,
	}
}
