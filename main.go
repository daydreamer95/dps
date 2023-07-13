package main

import "fmt"

func main() {

	heap := NewMinHeap()

	heap.Insert(Message{Priority: 2})
	heap.Insert(Message{Priority: 3})
	heap.Insert(Message{Priority: 30})
	heap.Insert(Message{Priority: 8})
	heap.Insert(Message{Priority: 4})
	heap.Insert(Message{Priority: 6})
	heap.Insert(Message{Priority: 15})
	heap.Insert(Message{Priority: 5})
	heap.Insert(Message{Priority: 7})
	heap.Insert(Message{Priority: 0})

	fmt.Println("HEAP:", heap.Data)

	poll, err := heap.Poll()
	if err != nil {
		fmt.Println("Error poll:", err)
		return
	}
	fmt.Println("Poll message:", poll)
	fmt.Println("HEAP:", heap.Data)

}
