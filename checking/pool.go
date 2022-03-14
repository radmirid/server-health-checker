package checking

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Job struct {
	Link string
}

type Result struct {
	Link         string
	StatusCode   int
	ResponseTime time.Duration
	Error        error
}

func (r Result) GetInfo() string {
	if r.Error != nil {
		return fmt.Sprintf("No access to %s:  %s", r.Link, r.Error.Error())
	}

	return fmt.Sprintf("%s is working | Response Time: %s | Status Code: %d", r.Link, r.ResponseTime.String(), r.StatusCode)
}

type Pool struct {
	worker       *worker
	workersCount int

	jobs    chan Job
	results chan Result

	wg     *sync.WaitGroup
	isStop bool
}

func (p *Pool) newWorker(id int) {
	for job := range p.jobs {
		time.Sleep(time.Second)
		p.results <- p.worker.handle(job)
		p.wg.Done()
	}

	log.Printf("Worker %d has finished", id)
}

func Run(workersCount int, timeout time.Duration, results chan Result) *Pool {
	return &Pool{
		worker:       newWorker(timeout),
		workersCount: workersCount,
		jobs:         make(chan Job),
		results:      results,
		wg:           new(sync.WaitGroup),
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.workersCount; i++ {
		go p.newWorker(i)
	}
}

func (p *Pool) Push(j Job) {
	if p.isStop {
		return
	}

	p.jobs <- j
	p.wg.Add(1)
}

func (p *Pool) Stop() {
	p.isStop = true
	close(p.jobs)
	p.wg.Wait()
}
