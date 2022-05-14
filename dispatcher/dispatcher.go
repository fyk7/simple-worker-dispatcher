package dispatcher

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Dispatcher struct {
	jobCh        chan *Job
	workerPoolCh chan *Worker // Queue of workers that process job.
	numWorkers   int
	wgDispatcher sync.WaitGroup
	wgWorker     sync.WaitGroup
}

func NewDispatcher(numWorkers int) *Dispatcher {
	return &Dispatcher{
		jobCh:        make(chan *Job, 100),
		workerPoolCh: make(chan *Worker, numWorkers),
		numWorkers:   numWorkers,
	}
}

func (d *Dispatcher) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	d.wgDispatcher.Add(1)
	d.startWorkers(ctx)

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)
	signal.Notify(sigCh,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
		os.Interrupt,
	)

	go func() {
		<-sigCh
		d.wgWorker.Done()
		cancel()
		os.Exit(1)
	}()

	// Watch incoming jobs.
	d.FetchJobs()

	d.wgDispatcher.Wait()
}

func (d *Dispatcher) startWorkers(ctx context.Context) {
	for i := 0; i < d.numWorkers; i++ {
		w := NewWorker(i)
		w.Start(ctx)
		d.workerPoolCh <- w // Register workers to worker pool.
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case job := <-d.jobCh:
				freeWorker := <-d.workerPoolCh // Get single worker from worker pool.
				d.wgWorker.Add(1)
				go func(job *Job, w *Worker) {
					w.Execute(job)
					d.wgWorker.Done()
					d.workerPoolCh <- w // Push back single worker to worker pool.
				}(job, freeWorker)
			}
		}
	}()
}

func (d *Dispatcher) AddJob(j *Job) {
	d.jobCh <- j
}

// Overwrite Here To Generate Jobs to process.
func (d *Dispatcher) FetchJobs() {
	for i := 1; i <= 10; i++ {
		d.AddJob(&Job{ID: i, Message: "Hello Apple!"})
	}
}
