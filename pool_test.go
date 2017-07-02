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

package exec_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/kbolino/go-exec"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewPool(t *testing.T) {
	cases := []struct {
		desc      string
		queueSize int
		workers   int
		panics    bool
	}{
		{"panics if queueSize < 0", -1, 1, true},
		{"panics if workers < 0", 1, -1, true},
		{"panics if workers == 0", 1, 0, true},
		{"does not panic if queueSize >= 0", 0, 1, false},
		{"does not panic if workers > 0", 1, 1, false},
	}
	Convey("NewPool", t, func() {
		for _, c := range cases {
			Convey(c.desc, func() {
				var pool *exec.Pool
				f := func() {
					pool = exec.NewPool(c.queueSize, c.workers)
				}
				if c.panics {
					So(f, ShouldPanic)
				} else {
					So(f, ShouldNotPanic)
					So(pool, ShouldNotBeNil)
				}
			})
		}
	})
}

func TestPool_Do(t *testing.T) {
	Convey("Pool.Do enqueues excess tasks", t, func() {
		pool := exec.NewPool(1, 2)
		counter := int32(0)
		go func() {
			for i := 0; i < 3; i++ {
				pool.Do(func() {
					atomic.AddInt32(&counter, 1)
					time.Sleep(50 * time.Millisecond)
				})
			}
			atomic.AddInt32(&counter, 1)
		}()
		time.Sleep(10 * time.Millisecond)
		So(atomic.LoadInt32(&counter), ShouldEqual, 3)
		time.Sleep(100 * time.Millisecond)
		So(atomic.LoadInt32(&counter), ShouldEqual, 4)
	})
}

func TestPool_Stop(t *testing.T) {
	Convey("Pool.Stop does not panic when called multiple times", t, func() {
		pool := exec.NewPool(1, 2)
		pool.Stop()
		f := func() {
			pool.Stop()
		}
		So(f, ShouldNotPanic)
	})
}

func TestPool_Try(t *testing.T) {
	Convey("Pool.Try doesn't block or enqueue excess tasks", t, func() {
		pool := exec.NewPool(1, 1)
		counter := int32(0)
		tries := []bool{}
		go func() {
			for i := 0; i < 3; i++ {
				try := pool.Try(func() {
					atomic.AddInt32(&counter, 1)
					time.Sleep(50 * time.Millisecond)
				})
				tries = append(tries, try)
				time.Sleep(10 * time.Millisecond)
			}
		}()
		time.Sleep(10 * time.Millisecond)
		So(atomic.LoadInt32(&counter), ShouldEqual, 1)
		time.Sleep(100 * time.Millisecond)
		So(atomic.LoadInt32(&counter), ShouldEqual, 2)
		So(tries, ShouldResemble, []bool{true, true, false})
	})
}

func TestPool_TryUntil(t *testing.T) {
	Convey("Pool.TryUntil times out when a task cannot be queued", t, func() {
		pool := exec.NewPool(1, 1)
		counter := int32(0)
		tries := []bool{}
		go func() {
			for i := 0; i < 3; i++ {
				try := pool.TryUntil(50*time.Millisecond, func() {
					atomic.AddInt32(&counter, 1)
					time.Sleep(50 * time.Millisecond)
				})
				tries = append(tries, try)
			}
			atomic.AddInt32(&counter, 1)
		}()
		time.Sleep(10 * time.Millisecond)
		So(atomic.LoadInt32(&counter), ShouldEqual, 1)
		time.Sleep(100 * time.Millisecond)
		So(atomic.LoadInt32(&counter), ShouldEqual, 3)
		So(tries, ShouldResemble, []bool{true, true, false})
	})
}

func TestPool_Wait(t *testing.T) {
	Convey("Pool.Wait blocks until all tasks complete", t, func() {
		pool := exec.NewPool(1, 2)
		counter := int32(0)
		go func() {
			for i := 0; i < 3; i++ {
				pool.Do(func() {
					time.Sleep(50 * time.Millisecond)
					atomic.AddInt32(&counter, 1)
				})
			}
		}()
		time.Sleep(10 * time.Millisecond)
		So(atomic.LoadInt32(&counter), ShouldEqual, 0)
		pool.Stop()
		pool.Wait()
		So(atomic.LoadInt32(&counter), ShouldEqual, 3)
	})
}
