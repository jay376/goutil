package algo

import "fmt"

type item struct {
	key   int
	value int
	next  *item
	prev  *item
}

func (this *item) remove() {
	this.prev.next = this.next
	this.next.prev = this.prev
}

func (this *item) appendHead(it *item) {
	it.next = this
	it.prev = this.prev
	this.prev.next = it
	this.prev = it
}

type LRUCache struct {
	head     *item
	kvs      map[int]*item
	capacity int
	size     int
}

func Constructor(capacity int) LRUCache {
	c := LRUCache{
		kvs:      make(map[int]*item),
		capacity: capacity,
		size:     0,
		head:     &item{},
	}
	c.head.next = c.head
	c.head.prev = c.head
	return c
}
func (this *LRUCache) Dump() {
	p := this.head
	for n := this.size; n > 0; n-- {
		fmt.Printf("%d->", p.prev.key)
		p = p.prev
	}
	fmt.Printf(". size:%d\n", len(this.kvs))
}

func (this *LRUCache) Get(key int) int {
	if it, ok := this.kvs[key]; ok {
		it.remove()
		this.head.appendHead(it)
		return it.value
	}
	return -1
}

func (this *LRUCache) Put(key int, value int) {
	if it, ok := this.kvs[key]; ok {
		it.value = value
		it.remove()
		this.head.appendHead(it)
	} else {
		if this.size == this.capacity {
			delete(this.kvs, this.head.next.key)
			this.head.next.remove()
			this.size--
		}
		it := &item{
			key:   key,
			value: value,
		}
		this.kvs[key] = it
		this.head.appendHead(it)
		this.size++
	}
	return
}
