// Copyright 2022 MaoLongLong. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lockfreequeue

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestLockFreeQueue(t *testing.T) {
	const n = 10000

	var (
		q   = New[int]()
		wg  sync.WaitGroup
		cnt uint32 // atomic
	)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < n; i++ {
				q.Enqueue(i)
			}
		}()
	}

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				_, ok := q.Dequeue()
				if ok {
					atomic.AddUint32(&cnt, 1)
				}
				if !ok && atomic.LoadUint32(&cnt) == 2*n {
					break
				}
			}
		}()
	}

	wg.Wait()
}
