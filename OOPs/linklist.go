package main

import "fmt"

type ListNode struct {
	val  int
	next *ListNode
}

type LinkedList struct {
	head *ListNode
}

func NewListNode(val int) *ListNode {
	return &ListNode{
		val: val,
	}
}

func (list *LinkedList) addNode(val int) {
	temp := NewListNode(val)

	if list.head == nil {
		list.head = temp
		return
	}

	// Reach the end of the list.
	curr := list.head
	for curr.next != nil {
		curr = curr.next
	}

	curr.next = temp
}

func (list *LinkedList) print() {
	curr := list.head
	for curr != nil {
		fmt.Printf("%d -> ", curr.val)
		curr = curr.next
	}
	fmt.Println("nil")
}

func main() {
	list := LinkedList{}

	list.addNode(10)
	list.addNode(30)
	list.addNode(50)
	list.addNode(80)

	list.print()
}
