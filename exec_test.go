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

func TestDirect_Do(t *testing.T) {
	Convey("Direct.Do runs task on the calling goroutine", t, func() {
		direct := exec.Direct
		counter := int32(0)
		go func() {
			direct.Do(func() {
				time.Sleep(50 * time.Millisecond)
			})
			atomic.AddInt32(&counter, 1)
		}()
		time.Sleep(10 * time.Millisecond)
		So(atomic.LoadInt32(&counter), ShouldEqual, 0)
		time.Sleep(100 * time.Millisecond)
		So(atomic.LoadInt32(&counter), ShouldEqual, 1)
	})
}
