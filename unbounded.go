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

// Unbounded is an asynchronous execution strategy that runs each task on a
// new goroutine.
type Unbounded struct {
	waitGroup sync.WaitGroup
}

var _ Strategy = &Unbounded{}

// NewUnbounded creates a new unbounded executor.
func NewUnbounded() *Unbounded {
	return &Unbounded{}
}

// Do runs task on a new goroutine.
func (u *Unbounded) Do(task func()) {
	u.waitGroup.Add(1)
	go func(u *Unbounded) {
		defer u.waitGroup.Done()
		task()
	}(u)
}

// Wait blocks until all running tasks are done.
func (u *Unbounded) Wait() {
	u.waitGroup.Wait()
}
