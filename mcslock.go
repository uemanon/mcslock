package mcslock

import (
    "runtime"
    "sync"
    "sync/atomic"
    "unsafe"
)

const (
    locked int32 = iota
    free
)

// SpinLock is a mcs-spinlock implementation
type SpinLock struct {
    _    [48]byte
    tail *mcs
    cur  *mcs
}

type mcs struct {
    _     [48]byte
    next  *mcs
    state int32
}

// Lock locks spinlock
func (sl *SpinLock) Lock() {
    local := new(mcs)
    prev := (*mcs)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&sl.tail)), unsafe.Pointer(local)))
    if prev == nil {
        sl.cur = local
        return
    }

    atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prev.next)), unsafe.Pointer(local))

    for !atomic.CompareAndSwapInt32(&local.state, free, locked) {
        runtime.Gosched()
    }

    sl.cur = local
}

// Unlock unlocks spinlock
func (sl *SpinLock) Unlock() {
    local := (*mcs)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&sl.cur))))
    if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&sl.tail)), unsafe.Pointer(local), unsafe.Pointer(nil)) {
        return
    }

    var next *mcs
    for next == nil {
        next = (*mcs)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&local.next))))
    }

    atomic.StoreInt32(&next.state, free)
}

var _ sync.Locker = (*SpinLock)(nil)
