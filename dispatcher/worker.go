package dispatcher

import (
	"context"
	"fmt"
)

type Worker struct {
	id int
}

func NewWorker(id int) *Worker {
	return &Worker{
		id: id,
	}
}

func (w *Worker) Start(ctx context.Context) {
	fmt.Printf("worker %d is starting!\n", w.id)
	go func() {
		<-ctx.Done()
	}()
}

func (w *Worker) Execute(j *Job) {
	j.Do()
}
