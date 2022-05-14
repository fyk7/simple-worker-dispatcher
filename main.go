package main

import (
	"context"

	"github.com/fyk7/simple-worker-dispatcher/dispatcher"
)

func main() {
	ctx := context.Background()
	d := dispatcher.NewDispatcher(10)
	d.Run(ctx)
}
