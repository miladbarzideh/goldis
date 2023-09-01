package datastore

type HeapItem struct {
	value int64
	ref   *int32
}

type MinHeap struct {
	heap []HeapItem
}

func NewMinHeap() *MinHeap {
	return &MinHeap{}
}

func (h *MinHeap) Insert(item HeapItem) {
	*item.ref = int32(len(h.heap))
	h.heap = append(h.heap, item)
	h.heapUp(int32(len(h.heap) - 1))
}

func (h *MinHeap) Update(i int32, value int64) {
	h.heap[i].value = value
	if i > 0 && h.heap[parent(i)].value > h.heap[i].value {
		h.heapUp(i)
	} else {
		h.heapDown(i)
	}
}

func (h *MinHeap) Remove(i int32) {
	lastIndex := len(h.heap) - 1
	h.heap[i] = h.heap[lastIndex]
	h.heap = h.heap[:lastIndex]
	h.heapDown(i)
}

func (h *MinHeap) Get(i int32) *HeapItem {
	if i > int32(len(h.heap))-1 {
		return nil
	}
	return &h.heap[i]
}

func (h *MinHeap) heapUp(i int32) {
	t := h.heap[i]
	if i > 0 && t.value < h.heap[parent(i)].value {
		h.heap[i], h.heap[parent(i)] = h.heap[parent(i)], h.heap[i]
		*h.heap[i].ref, *h.heap[parent(i)].ref = i, parent(i)
		h.heapUp(parent(i))
	}
}

func (h *MinHeap) heapDown(i int32) {
	length := int32(len(h.heap))
	if length == 0 {
		return
	}
	minIndex := i
	if left(i) < length && h.heap[left(i)].value < h.heap[minIndex].value {
		minIndex = left(i)
	}
	if right(i) < length && h.heap[right(i)].value < h.heap[minIndex].value {
		minIndex = right(i)
	}
	if minIndex != i {
		h.heap[i], h.heap[minIndex] = h.heap[minIndex], h.heap[i]
		*h.heap[i].ref = i
		h.heapDown(minIndex)
	}
}

func parent(i int32) int32 {
	return ((i + 1) / 2) - 1
}

func left(i int32) int32 {
	return i*2 + 1
}

func right(i int32) int32 {
	return i*2 + 2
}
