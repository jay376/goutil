// futu code
package algo

type entry struct {
	key, value int
	hash       int
}

type node struct {
	prev, next *node
	entry      *entry
}

func (head *node) append(n *node) {
	n.next = head.next
	n.prev = head
	head.next.prev = n
	head.next = n
}

type Htable struct {
	buckets  []*node
	capacity int
}

func hash(key int) int {
	return key
}

func NewHtable(cap int) *Htable {
	ht := &Htable{
		buckets:  make([]*node, cap),
		capacity: cap,
	}
	return ht
}

func (h *Htable) Insert(key, value int) {
	hvalue := hash(key)
	target := hvalue % h.capacity
	new := &node{
		entry: &entry{
			key:   key,
			value: value,
			hash:  hvalue,
		},
	}
	new.next = new
	new.prev = new

	if h.buckets[target] == nil {
		h.buckets[target] = new
		return
	}
	begin := h.buckets[target]
	if begin.entry.key == key {
		begin.entry.value = value
		return
	}
	for p := begin.next; p != begin; p = p.next {
		if p.entry.key == key {
			p.entry.value = value
			return
		}
	}
	h.buckets[target].append(new)
	return
}

func (h *Htable) Get(key int) int {
	hvalue := hash(key)
	begin := h.buckets[hvalue%h.capacity]
	if begin.entry.key == key {
		return begin.entry.value
	}
	for p := begin.next; p != begin; p = p.next {
		if p.entry.key == key {
			return p.entry.value
		}
	}

	return -1
}
