package datastore

import "testing"

func TestMinHeap_Insert(t *testing.T) {
	h := NewMinHeap()
	index1, index2, index3, index4 := int32(0), int32(0), int32(0), int32(0)
	h.Insert(HeapItem{value: 5, ref: &index1})
	h.Insert(HeapItem{value: 2, ref: &index2})
	h.Insert(HeapItem{value: 9, ref: &index3})
	h.Insert(HeapItem{value: 1, ref: &index4})

	expected := []int64{1, 2, 9, 5}
	for i, value := range expected {
		if value != h.Get(int32(i)).value {
			t.Errorf("Expected value %d at index %d, got %d", value, i, h.Get(int32(i)).value)
		}
		if int32(i) != *h.Get(int32(i)).ref {
			t.Errorf("Expected index to be %d, got %d", i, *h.Get(int32(i)).ref)
		}
	}
}

func TestMinHeap_Update(t *testing.T) {
	h := NewMinHeap()
	index1, index2, index3, index4 := int32(0), int32(0), int32(0), int32(0)
	h.Insert(HeapItem{value: 5, ref: &index1})
	h.Insert(HeapItem{value: 2, ref: &index2})
	h.Insert(HeapItem{value: 9, ref: &index3})
	h.Insert(HeapItem{value: 1, ref: &index4})

	h.Update(2, 3)

	expected := []int64{1, 2, 3, 5}
	for i, value := range expected {
		if value != h.Get(int32(i)).value {
			t.Errorf("Expected value %d at index %d, got %d", value, i, h.Get(int32(i)).value)
		}
		if int32(i) != *h.Get(int32(i)).ref {
			t.Errorf("Expected index to be %d, got %d", i, *h.Get(int32(i)).ref)
		}
	}
}

func TestMinHeap_Remove(t *testing.T) {
	h := NewMinHeap()
	index1, index2, index3, index4 := int32(0), int32(0), int32(0), int32(0)
	h.Insert(HeapItem{value: 5, ref: &index1})
	h.Insert(HeapItem{value: 2, ref: &index2})
	h.Insert(HeapItem{value: 9, ref: &index3})
	h.Insert(HeapItem{value: 1, ref: &index4})

	h.Remove(0)

	expected := []int64{2, 5, 9}
	for i, value := range expected {
		if value != h.Get(int32(i)).value {
			t.Errorf("Expected value %d at index %d, got %d", value, i, h.Get(int32(i)).value)
		}
		if int32(i) != *h.Get(int32(i)).ref {
			t.Errorf("Expected index to be %d, got %d", i, *h.Get(int32(i)).ref)
		}
	}
}
