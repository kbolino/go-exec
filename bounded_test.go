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

func TestBounded_Do(t *testing.T) {
	Convey("Bounded.Do limits the number of running goroutines", t, func() {
		counter := int32(0)
		bounded := exec.NewBounded(2)
		go func() {
			for i := 0; i < 3; i++ {
				bounded.Do(func() {
					atomic.AddInt32(&counter, 1)
					time.Sleep(50 * time.Millisecond)
				})
			}
		}()
		time.Sleep(10 * time.Millisecond)
		So(atomic.LoadInt32(&counter), ShouldEqual, 2)
		time.Sleep(100 * time.Millisecond)
		So(atomic.LoadInt32(&counter), ShouldEqual, 3)
	})
}

func TestBounded_Try(t *testing.T) {
	Convey("Bounded.Try doesn't block", t, func() {
		counter := int32(0)
		bounded := exec.NewBounded(2)
		tries := []bool{}
		go func() {
			for i := 0; i < 3; i++ {
				try := bounded.Try(func() {
					atomic.AddInt32(&counter, 1)
					time.Sleep(50 * time.Millisecond)
				})
				tries = append(tries, try)
			}
		}()
		time.Sleep(10 * time.Millisecond)
		So(tries, ShouldResemble, []bool{true, true, false})
		So(atomic.LoadInt32(&counter), ShouldEqual, 2)
		time.Sleep(100 * time.Millisecond)
		So(atomic.LoadInt32(&counter), ShouldEqual, 2)
	})
}

func TestBounded_Wait(t *testing.T) {
	Convey("Bounded.Wait blocks until all goroutines are done", t, func() {
		counter := int32(0)
		bounded := exec.NewBounded(2)
		go func() {
			for i := 0; i < 3; i++ {
				bounded.Do(func() {
					time.Sleep(50 * time.Millisecond)
					atomic.AddInt32(&counter, 1)
				})
			}
		}()
		time.Sleep(10 * time.Millisecond)
		bounded.Wait()
		So(atomic.LoadInt32(&counter), ShouldEqual, 3)
	})
}

func TestNewBounded(t *testing.T) {
	cases := []struct {
		desc   string
		n      int
		panics bool
	}{
		{"panics when n < 0", -1, true},
		{"panics when n == 0", 0, true},
		{"does not panic when n > 0", 1, false},
	}
	Convey("NewBounded", t, func() {
		for _, c := range cases {
			Convey(c.desc, func() {
				var bounded *exec.Bounded
				f := func() {
					bounded = exec.NewBounded(c.n)
				}
				if c.panics {
					So(f, ShouldPanic)
				} else {
					So(f, ShouldNotPanic)
					So(bounded, ShouldNotBeNil)
				}

			})
		}
	})
}
