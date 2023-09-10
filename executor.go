package concurrent

//The interface that captures the basic function of an executor.
type Executor interface {
    //Submit submits a result-less function as a task.
    Submit(func())
    //Invoke submits a result-ful function as a task. It will return a channel from which the result can be obtained once the task has completed.
    Invoke(func() interface{}) <-chan interface{}
}

type work struct {
    task func() interface{}
    resultChannel chan<- interface{}
}

type worker struct {
    wChannel <-chan work
    quitChannel chan byte
}

func newWorker(workChannel <-chan work) *worker {
    return &worker{workChannel, make(chan byte, 1)}
}

func (w *worker) start() {
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
type FixedSizeThreadPoolExecutor struct {
    workers []*worker
    workChannel chan work
}

//NewFixedSizeThreadPoolExecutor builds a new FixedSizeThreadPoolExecutor with the given thread pool size and returns a pointer to this new Executor.
func NewFixedSizeThreadPoolExecutor(poolSize uint64) *FixedSizeThreadPoolExecutor {
    workChannel := make(chan work, poolSize)
    workers := make([]*worker, poolSize)
    for i, _ := range workers {
        w := newWorker(workChannel)
        workers[i] = w
        w.start()
    }
    return &FixedSizeThreadPoolExecutor{workers, workChannel}
}

//Submit submits a result-less function a task. This method will block until a thread becomes available. Calling this method on an Executor that has been shut down results in a run-time panic.
func (e *FixedSizeThreadPoolExecutor) Submit(f func()) {
    e.workChannel <- work{func() interface{} {f(); return nil}, nil}
}

//Invoke submits a result-ful function a task. This method will block until a thread becomes available. Once a thread is available, this method will return a channel which will hold the result after the task has completed.
//Calling this method on an Executor that has been shut down results in a run-time panic.
func (e *FixedSizeThreadPoolExecutor) Invoke(f func() interface{}) <-chan interface{} {
    resultChannel := make(chan interface{}, 1)
    e.workChannel <- work{f, resultChannel}
    return resultChannel
}

//ShutDownGracefully shuts the Executor down gracefully by allowing all running tasks to complete and then all threads to stop.
func (e *FixedSizeThreadPoolExecutor) ShutDownGracefully() {
    close(e.workChannel)
    for _, w := range e.workers {
        <- w.quitChannel
    }
}

