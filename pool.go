package gopartypool

import (
	"context"
	"sync"
)

type todo struct {
	input  interface{}
	output *interface{}
	wg     *sync.WaitGroup
	Do     func(interface{}) interface{}
}

type Worker struct {
	tasks chan todo
}

func (w *Worker) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-w.tasks:
			*task.output = task.Do(task.input)
			task.wg.Done()
		}
	}
}

type Pool struct {
	mu       sync.Mutex
	pool     []*Worker
	tasks    chan todo
	workerDo func(interface{}) interface{}
}

func (p *Pool) Do(input interface{}) interface{} {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	output := new(interface{})
	p.tasks <- todo{
		input:  input,
		output: output,
		Do:     p.workerDo,
		wg:     wg,
	}
	wg.Wait()
	return *output
}

func (p *Pool) AddWorker() {
	p.pool = append(p.pool, &Worker{
		tasks: p.tasks,
	})
}

func (p *Pool) Run(ctx context.Context) {
	for _, w := range p.pool {
		go w.run(ctx)
	}
}

func NewPool(workers int, workerDo func(interface{}) interface{}) *Pool {
	if workers < 1 {
		panic("workers must be more that 1, because of lock")
	}
	p := &Pool{
		mu:       sync.Mutex{},
		pool:     make([]*Worker, 0),
		tasks:    make(chan todo),
		workerDo: workerDo,
	}
	for i := 0; i < workers; i++ {
		p.AddWorker()
	}
	return p
}
