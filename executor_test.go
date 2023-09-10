package concurrent

import (
    "sync/atomic"
    "testing"
)

func TestFixedSizeThreadPoolExecutorSubmitPoolSize1(t *testing.T) {
    executor := NewFixedSizeThreadPoolExecutor(1)
    var atomicBool atomic.Bool
    startChannel := make(chan byte)
    executor.Submit(func() {
        <-startChannel
        atomicBool.Store(true)
    })
    startChannel <- byte(0)
    executor.Submit(func() {
        if !atomicBool.Load() {
            t.Error("Task 2 should start after task 1 has completed.")
        }
    })
    executor.ShutDownGracefully()
}

func TestFixedSizeThreadPoolExecutorSubmitPoolSize2(t *testing.T) {
    executor := NewFixedSizeThreadPoolExecutor(2)
    startChannel := make(chan byte)
    executor.Submit(func() {
        startChannel <- byte(0)
    })
    executor.Submit(func() {
        <- startChannel
    })
    executor.ShutDownGracefully()
}

func TestFixedSizeThreadPoolExecutorInvokePoolSize1(t *testing.T) {
    executor := NewFixedSizeThreadPoolExecutor(1)
    var atomicBool atomic.Bool
    startChannel := make(chan byte)
    result1Chan := executor.Invoke(func() interface{} {
        <-startChannel
        atomicBool.Store(true)
        return 1
    })
    startChannel <- byte(0)
    result2Chan := executor.Invoke(func() interface{} {
        if !atomicBool.Load() {
            t.Error("Task 2 should start after task 1 has completed.")
        }
        return 2
    })
    result1I, result2I := <- result1Chan, <-result2Chan
    result1, result1Ok := result1I.(int)
    if !result1Ok {
        t.Error("Incorrect type returned from Invoke.")
    }
    if result1 != 1 {
        t.Error("Incorrect value returned from Invoke.")
    }
    result2, result2Ok := result2I.(int)
    if !result2Ok {
        t.Error("Incorrect type returned from Invoke.")
    }
    if result2 != 2 {
        t.Error("Incorrect value returned from Invoke.")
    }
    executor.ShutDownGracefully()
}

func TestFixedSizeThreadPoolExecutorInvokePoolSize2(t *testing.T) {
    executor := NewFixedSizeThreadPoolExecutor(2)
    startChannel := make(chan byte)
    result1Chan := executor.Invoke(func() interface{} {
        startChannel <- byte(0)
        return 1
    })
    result2Chan := executor.Invoke(func() interface{} {
        <- startChannel
        return 2
    })
    result1I, result2I := <- result1Chan, <-result2Chan
    result1, result1Ok := result1I.(int)
    if !result1Ok {
        t.Error("Incorrect type returned from Invoke.")
    }
    if result1 != 1 {
        t.Error("Incorrect value returned from Invoke.")
    }
    result2, result2Ok := result2I.(int)
    if !result2Ok {
        t.Error("Incorrect type returned from Invoke.")
    }
    if result2 != 2 {
        t.Error("Incorrect value returned from Invoke.")
    }
    executor.ShutDownGracefully()
}

