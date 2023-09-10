# concurrent
This is a package to expose types useful for managing concurrent execution of tasks.

## `Executor`
This is a basic interface type to be implementd by executors. Executors allow tasks to be executed in threads other than the main thread.

## `FixedSizeThreadPoolExecutor`
This is an exector which allocates a fixed number of threads. The number of threads is passed in to the constructor. New tasks are assigned to a free thread. If there are no free threads available then the task sumission will block until a thread becomes free.
This executor is useful and efficient if you are confident that all or most the threads will be engaged for the life of the executor, and if you wish to limit the number of threads (for instance if you want to limit the CPU workload on the host machine).

