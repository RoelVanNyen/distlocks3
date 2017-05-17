package distlocks3

import (
	"sync"
	"testing"
)

func testLock(p *sync.WaitGroup) {
	lockID := AquireLock("environments-test", "environments/test", "us-west-2")
	ReleaseLock("environments-test", lockID, "us-west-2")
	p.Done()
}

func TestLockingMultiple(*testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go testLock(&wg)

	}
	wg.Wait()

}
