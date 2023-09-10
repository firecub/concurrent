package concurrent

//The interface that captures the basic function of an executor.
type Executor [T any] interface {
    //Submit submits a result-less function as a task.
    Submit(func())
    //Invoke submits a result-ful function as a task. It will return a channel from which the result can be obtained once the task has completed.
    Invoke(func() T) <-chan T
}

type work [T any] struct {
    task func() T
    resultChannel chan<- T
}

type worker [T any] struct {
    wChannel <-chan work [T]
    quitChannel chan byte
}

func newWorker [T any] (workChannel <-chan work [T]) *worker[T] {
    return &worker[T]{workChannel, make(chan byte, 1)}
}

func (w *worker[T]) start() {
    go func() {
        for newWork, ok := <- w.wChannel; ok; newWork, ok = <- w.wChannel {
            if newWork.resultChannel == nil {
                newWork.task()
            } else {
                newWork.resultChannel <- newWork.task()
            }
        }
        w.quitChannel <- byte(0)
    }()
}

//FixedSizeThreadPoolExecutor is an Executor which has a fixed number of threads.
type FixedSizeThreadPoolExecutor [T any] struct {
    workers []*worker[T]
    workChannel chan work[T]
}

//NewFixedSizeThreadPoolExecutor builds a new FixedSizeThreadPoolExecutor with the given thread pool size and returns a pointer to this new Executor.
func NewFixedSizeThreadPoolExecutor [T any] (poolSize uint64) *FixedSizeThreadPoolExecutor[T] {
    workChannel := make(chan work[T], poolSize)
    workers := make([]*worker[T], poolSize)
    for i, _ := range workers {
        w := newWorker(workChannel)
        workers[i] = w
        w.start()
    }
    return &FixedSizeThreadPoolExecutor[T]{workers, workChannel}
}

//Submit submits a result-less function a task. This method will block until a thread becomes available. Calling this method on an Executor that has been shut down results in a run-time panic.
func (e *FixedSizeThreadPoolExecutor[T]) Submit(f func()) {
    e.workChannel <- work[T]{func() T {f(); return *new(T)}, nil}
}

//Invoke submits a result-ful function a task. This method will block until a thread becomes available. Once a thread is available, this method will return a channel which will hold the result after the task has completed.
//Calling this method on an Executor that has been shut down results in a run-time panic.
func (e *FixedSizeThreadPoolExecutor[T]) Invoke(f func() T) <-chan T {
    resultChannel := make(chan T, 1)
    e.workChannel <- work[T]{f, resultChannel}
    return resultChannel
}

//ShutDownGracefully shuts the Executor down gracefully by allowing all running tasks to complete and then all threads to stop.
func (e *FixedSizeThreadPoolExecutor[T]) ShutDownGracefully() {
    close(e.workChannel)
    for _, w := range e.workers {
        <- w.quitChannel
    }
}

