package comm

import (
	"sync"
)

// RingQueue ...
type RingQueue struct {
	mu       sync.Mutex
	cond     *sync.Cond
	buffer   []interface{}
	readidx  int
	writeidx int
	size     int // buffer size
	hasWait  int
	close    bool
}

// NewRingQueue, enqueue when full will eliminate the earliest element.
func NewRingQueue(size int) *RingQueue {
	q := &RingQueue{
		buffer:   make([]interface{}, size),
		readidx:  0,
		writeidx: 0,
		size:     size,
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func NewRingQueueFromSlice(arr []interface{}, length int) *RingQueue {
	q := &RingQueue{
		buffer:   arr,
		readidx:  0,
		writeidx: length,
		size:     cap(arr),
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Change queue size. If queue size is decreased, return removed elements.
func (r *RingQueue) Change(size int) (elements []interface{}) {
	buffer := make([]interface{}, size)
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.size == size {
		return
	}
	// empty
	if r.readidx == r.writeidx {
		r.size = size
		r.buffer = buffer
		r.readidx = 0
		r.writeidx = 0
		return
	}

	rdx := r.readidx % r.size
	wdx := r.writeidx % r.size
	length := r.writeidx - r.readidx
	if length > size { // shrink size, return shrink element
		eliminate := length - size
		length = size
		rdx = (r.readidx + eliminate) % r.size
		elimateIdx := r.readidx % r.size
		elements = make([]interface{}, eliminate)
		if elimateIdx < rdx {
			copy(elements, r.buffer[elimateIdx:rdx])
		} else {
			copy(elements, r.buffer[elimateIdx:])
			copy(elements[r.size-elimateIdx:], r.buffer[0:rdx])
		}
	}

	if size > 0 {
		if rdx < wdx {
			copy(buffer, r.buffer[rdx:wdx])
		} else {
			copy(buffer, r.buffer[rdx:])
			copy(buffer[r.size-rdx:], r.buffer[0:wdx])
		}
	}

	r.size = size
	r.buffer = buffer
	r.readidx = 0
	r.writeidx = length
	return
}

// Put will replace and return the replaced element when queue is full
func (r *RingQueue) Put(element interface{}) (e interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.size <= 0 {
		return element
	}

	dst := r.writeidx % r.size
	if r.writeidx == r.readidx+r.size { // queue full
		e = r.buffer[dst]
		r.readidx++ // overload readidx
	}
	r.buffer[dst] = element
	r.writeidx++
	if r.hasWait > 0 {
		r.cond.Signal()
	}
	return
}

// Pop num elements, if empty return nil
func (r *RingQueue) unsafePop(num int) (elements []interface{}) {
	if r.writeidx-r.readidx < num {
		num = r.writeidx - r.readidx
	}
	if num == 0 {
		return
	}
	elements = make([]interface{}, num)
	for idx := 0; idx < num; idx++ {
		i := (r.readidx + idx) % r.size
		elements[idx] = r.buffer[i]
		r.buffer[i] = nil
	}
	r.readidx += num
	return
}

// Pop num elements, if empty return nil
func (r *RingQueue) Pop(num int) (elements []interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.size <= 0 {
		return
	}
	return r.unsafePop(num)
}

// PopWait num elements, if empty So wait
func (r *RingQueue) PopWait(num int) (elements []interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.size <= 0 {
		return
	}

	for r.writeidx == r.readidx && !r.close {
		r.hasWait++
		r.cond.Wait()
		r.hasWait--
	}
	return r.unsafePop(num)
}

func (r *RingQueue) Iterate(f func(req interface{}) (shouldRemove bool)) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := r.readidx; i < r.writeidx; i++ {
		idx := i % r.size

		ele := r.buffer[idx]
		if shouldRemove := f(ele); shouldRemove {
			lastIdx := (r.writeidx - 1) % r.size

			r.buffer[idx], r.buffer[lastIdx] = r.buffer[lastIdx], r.buffer[idx]
			r.buffer[lastIdx] = nil

			r.writeidx--
			i--
		}
	}
}

// Len ...
func (r *RingQueue) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.writeidx - r.readidx
}

// Cap ...
func (r *RingQueue) Cap() int {
	return r.size
}

func (r *RingQueue) GetBuffer() ([]interface{}, int) {
	return r.buffer, r.Len()
}

// Close ...
func (r *RingQueue) Close() {
	r.mu.Lock()
	r.close = true
	r.mu.Unlock()
	r.cond.Broadcast()
}
