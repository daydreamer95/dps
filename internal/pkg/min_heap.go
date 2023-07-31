package pkg

import (
	"dps/internal/pkg/entity"
	"errors"
	"sync"
)

type MinHeap struct {
	mu     sync.Mutex
	Length int
	Data   []entity.Item
}

func NewMinHeap() *MinHeap {
	var m MinHeap
	m.Length = 0
	m.Data = []entity.Item{}
	return &m
}

func (m *MinHeap) Insert(msg entity.Item) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Data = append(m.Data, msg)
	m.heapifyUp(m.Length)
	m.Length++
}

func (m *MinHeap) Poll() (entity.Item, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Length == 0 {
		//fmt.Println("ha")
		return entity.Item{}, errors.New("Cant poll. No item")
	}

	out := m.Data[0]
	m.Length--

	if m.Length == 0 {
		out := m.Data[0]
		m.Data = []entity.Item{}
		return out, nil
	}

	//fmt.Println("Poll m.Data[m.Length:", m.Data[m.Length])
	m.Data[0] = m.Data[m.Length]
	m.Data = m.Data[0:m.Length]
	m.heapifyDown(0)
	return out, nil
}

func (m *MinHeap) heapifyDown(index int) {
	if index > m.Length {
		return
	}

	lIndex := leftChild(index)
	rIndex := rightChild(index)

	if index >= m.Length || lIndex >= m.Length {
		return
	}

	lV := m.Data[lIndex]
	rV := m.Data[rIndex]
	v := m.Data[index]

	if lV.Priority > rV.Priority && v.Priority > rV.Priority {
		m.Data[index] = rV
		m.Data[rIndex] = v
		m.heapifyDown(rIndex)
	} else if rV.Priority > lV.Priority && v.Priority > lV.Priority {
		m.Data[index] = lV
		m.Data[lIndex] = v
		m.heapifyDown(lIndex)
	}
}

func (m *MinHeap) heapifyUp(index int) {
	if index == 0 {
		return
	}

	p := parent(index)
	parentV := m.Data[p]
	v := m.Data[index]

	//fmt.Printf("Compare %v > %v at pIndex %v and index %v? \n", parentV.Priority, v.Priority, p, index)
	if parentV.Priority > v.Priority {
		//fmt.Println("Parent value:", parentV)
		//fmt.Println("Index value:", v)

		m.Data[index] = parentV
		m.Data[p] = v
		m.heapifyUp(p)
	}
}

func parent(index int) int {
	return (index - 1) / 2
}

func leftChild(index int) int {
	return index*2 + 1
}

func rightChild(index int) int {
	return index*2 + 2
}