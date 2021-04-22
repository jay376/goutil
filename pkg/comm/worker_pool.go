package comm

// Task ...
type Task func()

// WorkerPool ...
type WorkerPool struct {
	num int
	chs []chan Task
}

// NewWorkerPool ...
func NewWorkerPool(num int) *WorkerPool {
	return &WorkerPool{
		num: num,
		chs: make([]chan Task, num),
	}
}

// Start ...
func (w *WorkerPool) Start() {
	bufNum := 500

	for i := 0; i < w.num; i++ {
		ch := make(chan Task, bufNum)
		w.chs[i] = ch

		go func() {
			for Task := range ch {
				Task()
			}
		}()
	}
}

// Stop ...
func (w *WorkerPool) Stop() {
	for i := 0; i < w.num; i++ {
		close(w.chs[i])
	}
}

// Schedule ...
func (w *WorkerPool) Schedule(t Task, i int) {
	w.chs[i%w.num] <- t
}
