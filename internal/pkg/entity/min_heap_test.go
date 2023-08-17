package entity

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinHeapDeletion(t *testing.T) {
	m := NewMinHeap()

	i1 := Item{Id: uuid.New().String(), Priority: 9}
	i2 := Item{Id: uuid.New().String(), Priority: 5}
	i3 := Item{Id: uuid.New().String(), Priority: 7}
	i4 := Item{Id: uuid.New().String(), Priority: 1}
	i5 := Item{Id: uuid.New().String(), Priority: 4}
	i6 := Item{Id: uuid.New().String(), Priority: 3}
	i7 := Item{Id: uuid.New().String(), Priority: 11}
	i8 := Item{Id: uuid.New().String(), Priority: 15}
	i9 := Item{Id: uuid.New().String(), Priority: 17}

	m.Insert(i1)
	m.Insert(i2)
	m.Insert(i3)
	m.Insert(i4)
	m.Insert(i5)
	m.Insert(i6)
	m.Insert(i7)
	m.Insert(i8)
	m.Insert(i9)

	successDelete := m.Delete(i3)
	//for _, v := range m.Data {
	//	fmt.Printf(" %v ", v.Priority)
	//}
	assert.Equal(t, successDelete, true)
}

func TestMinHeapDeletionAtLastElement(t *testing.T) {
	m := NewMinHeap()

	i1 := Item{Id: uuid.New().String(), Priority: 9}
	i2 := Item{Id: uuid.New().String(), Priority: 5}
	i3 := Item{Id: uuid.New().String(), Priority: 7}
	i4 := Item{Id: uuid.New().String(), Priority: 1}
	i5 := Item{Id: uuid.New().String(), Priority: 4}

	m.Insert(i1)
	m.Insert(i2)
	m.Insert(i3)
	m.Insert(i4)
	m.Insert(i5)

	m.Delete(i2)
	assert.Equal(t, m.Data[0].Priority, i4.Priority)
	assert.Equal(t, m.Data[1].Priority, i5.Priority)
	assert.Equal(t, m.Data[2].Priority, i3.Priority)
	assert.Equal(t, m.Data[3].Priority, i1.Priority)

}

func TestMinHeapDeletionAtNodeElement(t *testing.T) {
	m := NewMinHeap()

	i1 := Item{Id: uuid.New().String(), Priority: 9}
	i2 := Item{Id: uuid.New().String(), Priority: 5}
	i3 := Item{Id: uuid.New().String(), Priority: 7}
	i4 := Item{Id: uuid.New().String(), Priority: 1}
	i5 := Item{Id: uuid.New().String(), Priority: 4}

	m.Insert(i1)
	m.Insert(i2)
	m.Insert(i3)
	m.Insert(i4)
	m.Insert(i5)

	m.Delete(i5)
	assert.Equal(t, m.Data[0].Priority, i4.Priority)
	assert.Equal(t, m.Data[1].Priority, i2.Priority)
	assert.Equal(t, m.Data[2].Priority, i3.Priority)
	assert.Equal(t, m.Data[3].Priority, i1.Priority)

}

func TestMinHeapDeletionTopElement(t *testing.T) {
	m := NewMinHeap()

	i1 := Item{Id: uuid.New().String(), Priority: 9}
	i2 := Item{Id: uuid.New().String(), Priority: 5}
	i3 := Item{Id: uuid.New().String(), Priority: 7}
	i4 := Item{Id: uuid.New().String(), Priority: 1}
	i5 := Item{Id: uuid.New().String(), Priority: 4}
	i6 := Item{Id: uuid.New().String(), Priority: 3}
	i7 := Item{Id: uuid.New().String(), Priority: 11}
	i8 := Item{Id: uuid.New().String(), Priority: 15}
	i9 := Item{Id: uuid.New().String(), Priority: 17}

	m.Insert(i1)
	m.Insert(i2)
	m.Insert(i3)
	m.Insert(i4)
	m.Insert(i5)
	m.Insert(i6)
	m.Insert(i7)
	m.Insert(i8)
	m.Insert(i9)

	successDelete := m.Delete(i4)
	//for _, v := range m.Data {
	//	fmt.Printf(" %v ", v.Priority)
	//}
	assert.Equal(t, successDelete, true)
}

func TestMinHeapInitAndInsert(t *testing.T) {
	m := NewMinHeap()

	i1 := Item{Priority: 9}
	i2 := Item{Priority: 5}
	i3 := Item{Priority: 7}
	i4 := Item{Priority: 1}
	i5 := Item{Priority: 2}

	m.Insert(i1)
	m.Insert(i2)
	m.Insert(i3)
	m.Insert(i4)
	m.Insert(i5)

	assert.Equal(t, m.Data[0].Priority, i4.Priority)
	assert.Equal(t, m.Data[1].Priority, i5.Priority)
	assert.Equal(t, m.Data[2].Priority, i3.Priority)
	assert.Equal(t, m.Data[3].Priority, i1.Priority)
	assert.Equal(t, m.Data[4].Priority, i2.Priority)
}
