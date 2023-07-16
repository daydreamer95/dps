package main

import "fmt"

func main() {

	heap := NewMinHeap()

	heap.Insert(Item{Priority: 2})
	heap.Insert(Item{Priority: 3})
	heap.Insert(Item{Priority: 30})
	heap.Insert(Item{Priority: 8})
	heap.Insert(Item{Priority: 4})
	heap.Insert(Item{Priority: 6})
	heap.Insert(Item{Priority: 15})
	heap.Insert(Item{Priority: 5})
	heap.Insert(Item{Priority: 7})
	heap.Insert(Item{Priority: 0})

	fmt.Println("HEAP:", heap.Data)

	poll, err := heap.Poll()
	if err != nil {
		fmt.Println("Error poll:", err)
		return
	}
	fmt.Println("Poll message:", poll)
	fmt.Println("HEAP:", heap.Data)

}
