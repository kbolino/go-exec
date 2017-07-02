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

type empty struct{}

// Bounded is an asynchronous execution strategy that runs each task on a new
// goroutine, but limits the number of simultaneously running goroutines.
// Bounded is a useful alternative to Pool when:
//   - excess tasks should not be enqueued but instead Do should block
//     (although this behavior can be achieved with queueSize=0 on a Pool)
//   - the overhead of starting a new goroutine for each task is acceptable
//   - tasks are executed in batches, with a call to Wait after each batch
type Bounded struct {
	semaphore chan empty
}

var _ Strategy = &Bounded{}

// NewBounded creates a new bounded executor which will not allow more than n
// goroutines to be running simultaneously.
// NewBounded panics if n is less than 1.
func NewBounded(n int) *Bounded {
	if n < 1 {
		panic("n is less than 1")
	}
	semaphore := make(chan empty, n)
	return &Bounded{semaphore: semaphore}
}

// Do runs task on a new goroutine, unless the maximum number of running
// goroutines has been reached, in which case Do blocks until one or more
// running goroutines exit.
// To avoid blocking, use Try instead.
func (b *Bounded) Do(task func()) {
	b.acquire(1)
	go b.runTask(task)
}

// Try is the nonblocking version of Do.
// If task cannot be executed immediately, then Try does nothing with it
// and returns false.
func (b *Bounded) Try(task func()) bool {
	if !b.tryAcquire1() {
		return false
	}
	go b.runTask(task)
	return true
}

// Wait blocks until there are no running goroutines.
// Wait may prevent Do from starting new goroutines until Wait returns.
func (b *Bounded) Wait() {
	n := cap(b.semaphore)
	b.acquire(n)
	b.release(n)
	return
}

func (b *Bounded) acquire(n int) {
	e := empty{}
	for i := 0; i < n; i++ {
		b.semaphore <- e
	}
}

func (b *Bounded) release(n int) {
	for i := 0; i < n; i++ {
		<-b.semaphore
	}
}

func (b *Bounded) runTask(task func()) {
	defer b.release(1)
	task()
}

func (b *Bounded) tryAcquire1() bool {
	select {
	case b.semaphore <- empty{}:
		return true
	default:
		return false
	}
}
