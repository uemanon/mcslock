package mcslock

import (
    "fmt"
    "os"
    "reflect"
    "runtime"
    "sync"
    "testing"
    "time"
)

func TestSpinLock(t *testing.T) {
    testSpinLock(runtime.NumCPU(), 1 << 16, &SpinLock{})
}

func BenchmarkSpinLock(b *testing.B) {
    for i := 0; i < b.N; i++ {
        testSpinLock(runtime.NumCPU(), 1 << 16, &SpinLock{})
    }
}

func testSpinLock(threads, times int, l sync.Locker) {
    locker := typeOfLocker(l)
    for i := 1; i <= threads; i <<= 1 {
        var (
            wg sync.WaitGroup
            c  int64
        )
        t := time.Now()
        for k := 0; k < i; k++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                for j := 0; j < times; j++ {
                    l.Lock()
                    c++
                    c += 42
                    l.Unlock()
                }
            }()
        }
        wg.Wait()
        fmt.Fprintf(os.Stdout, "[%s] thread(s): %d times: %d cos: %dms\n", locker, i, times, time.Since(t).Milliseconds())
    }
}

func typeOfLocker(l sync.Locker) string {
    t := reflect.TypeOf(l)
    if t.Kind() == reflect.Ptr {
        return t.Elem().Name()
    }
    return t.Name()
}
