package entity

import (
	"errors"
	"sync"
)

type MinHeap struct {
	mu     sync.Mutex
	Length int
	Data   []Item
}

func NewMinHeap() *MinHeap {
	var m MinHeap
	m.Length = 0
	m.Data = []Item{}
	return &m
}

func (m *MinHeap) Insert(msg Item) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Data = append(m.Data, msg)
	m.heapifyUp(m.Length)
	m.Length++
}

func (m *MinHeap) Poll() (Item, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Length == 0 {
		//fmt.Println("ha")
		return Item{}, errors.New("Cant poll. No item")
	}

	out := m.Data[0]
	m.Length--

	if m.Length == 0 {
		out := m.Data[0]
		m.Data = []Item{}
		return out, nil
	}

	//fmt.Println("Poll m.Data[m.Length:", m.Data[m.Length])
	m.Data[0] = m.Data[m.Length]
	m.Data = m.Data[0:m.Length]
	m.heapifyDown(0)
	return out, nil
}

func (m *MinHeap) Delete(msg Item) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	msgIdx := m.search(msg)
	if msgIdx == -1 {
		return false
	}

	if msgIdx == m.Length-1 {
		m.Length--
		m.Data = m.Data[0:m.Length]
		return true
	}

	//Swap last and idx
	m.Data[msgIdx], m.Data[m.Length-1] = m.Data[m.Length-1], m.Data[msgIdx]
	m.Length--
	m.Data = m.Data[0:m.Length]
	if m.Data[parent(msgIdx)].Priority > m.Data[msgIdx].Priority {
		m.heapifyUp(msgIdx)
	} else {
		m.heapifyDown(msgIdx)
	}
	return true
}

// search: Search item in MinHeap based on Item.id
func (m *MinHeap) search(msg Item) int {
	for index, item := range m.Data {
		if msg.Id == item.Id {
			return index
		}
	}
	return -1
}

func (m *MinHeap) heapifyDown(index int) {
	if index > m.Length {
		return
	}

	lIndex := leftChild(index)
	rIndex := rightChild(index)

	if index >= m.Length || lIndex >= m.Length || rIndex >= m.Length {
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

// start starts goroutines for item expiry check
func (m *MinHeap) start() {

}
