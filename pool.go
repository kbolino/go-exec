// Copyright 2017 Kristian Bolino
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exec

import (
	"sync"
)

// Pool provides an asynchronous execution strategy using a fixed-size pool
// of worker goroutines.
// Compare Pool with Bounded as both provided a limit on the number of running
// goroutines.
type Pool struct {
	workQueue chan func()
	closeOnce sync.Once
	waitGroup sync.WaitGroup
}

// NewPool creates a new worker pool.
// If queueSize is nonzero, a buffered channel is created of that capacity.
// Otherwise, an unbuffered channel is created, meaning that Do will block
// until a worker is available.
// NewPool starts one goroutine for each worker.
// NewPool panics if queueSize is negative or workers is less than 1.
func NewPool(queueSize, workers int) *Pool {
	if queueSize < 0 {
		panic("queueSize is negative")
	} else if workers < 1 {
		panic("workers is less than 1")
	}
	workQueue := make(chan func(), queueSize)
	pool := &Pool{workQueue: workQueue}
	for i := 0; i < workers; i++ {
		pool.waitGroup.Add(1)
		go pool.workerMain()
	}
	return pool
}

// Do enqueues a task for the worker pool to execute.
// Do panics if Stop has previously been called.
func (p *Pool) Do(task func()) {
	p.workQueue <- task
}

// Stop tells the worker pool to cease accepting new tasks.
// Once all previously queued tasks have been executed, the goroutines in
// the pool will exit.
// Stop is idempotent; multiple calls to Stop have no additional effect.
func (p *Pool) Stop() {
	p.closeOnce.Do(func() {
		close(p.workQueue)
	})
}

// Wait blocks until the worker pool has finished executing all queued tasks.
// Wait can only be called after Stop.
func (p *Pool) Wait() {
	p.waitGroup.Wait()
}

func (p *Pool) workerMain() {
	defer p.waitGroup.Done()
	for task := range p.workQueue {
		task()
	}
}
