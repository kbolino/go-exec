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

// Package exec provides an execution framework for Go.
package exec

// Strategy abstracts the concept of running a task.
type Strategy interface {
	Do(task func())
}

// Direct is a synchronous execution strategy that runs each task on the
// calling goroutine.
// Direct implements Strategy.
var Direct = direct{}

type direct struct{}

var _ Strategy = direct{}

func (d direct) Do(task func()) {
	task()
}
