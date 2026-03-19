package ds

import (
	"testing"

	"github.com/roidaradal/tst"
)

func TestQueue(t *testing.T) {
	// NewQueue, Len, IsEmpty, NotEmpty
	q := NewQueue[int]()
	tst.AssertEqual(t, "Queue.String", q.String(), "[]")
	tst.AssertEqual(t, "Queue.Len", q.Len(), 0)
	tst.AssertEqual(t, "Queue.IsEmpty", q.IsEmpty(), true)
	tst.AssertEqual(t, "Queue.NotEmpty", q.NotEmpty(), false)
	// NewQueueFrom, Queue.Items
	items := []int{1, 2, 3}
	q = NewQueueFrom(items)
	tst.AssertEqual(t, "Queue.String", q.String(), "[1 2 3]")
	tst.AssertEqual(t, "Queue.Len", q.Len(), len(items))
	tst.AssertEqual(t, "Queue.IsEmpty", q.IsEmpty(), false)
	tst.AssertEqual(t, "Queue.NotEmpty", q.NotEmpty(), true)
	tst.AssertListEqual(t, "Queue.Items", q.Items(), items)
	// Copy
	q2 := q.Copy()
	tst.AssertListEqual(t, "Queue.Copy.Items", q2.Items(), q.Items())
	// Clear
	q2.Clear()
	tst.AssertEqual(t, "Queue.Clear.IsEmpty", q2.IsEmpty(), true)
	// Check original queue is unchanged
	tst.AssertListEqual(t, "Queue.Items", q.Items(), items)
}

func TestQueueEnqueue(t *testing.T) {
	q := NewQueue[int]()
	q.Enqueue(1)
	tst.AssertEqual(t, "Enqueue.MustFront", q.MustFront(), 1)
	q.Enqueue(2)
	q.Enqueue(3)
	tst.AssertEqual(t, "Enqueue.MustFront", q.MustFront(), 1)
	tst.AssertListEqual(t, "Queue.Items", q.Items(), []int{1, 2, 3})
}

func TestQueueFront(t *testing.T) {
	q := NewQueue[int]()
	front := q.Front()
	tst.AssertEqual(t, "Queue.Front", front.IsNil(), true)
	q.Enqueue(1)
	front = q.Front()
	tst.AssertEqual2(t, "Queue.Front", front.IsNil(), false, front.Value(), 1)
	tst.AssertEqual(t, "Queue.MustFront", q.MustFront(), 1)

	q.Dequeue()
	defer tst.AssertPanic(t, "Queue.MustFront")
	q.MustFront() // should panic
}

func TestQueueDequeue(t *testing.T) {
	q := NewQueueFrom([]int{1, 2})
	front := q.Dequeue()
	tst.AssertEqual2(t, "Queue.Dequeue", front.IsNil(), false, front.Value(), 1)
	frontItem := q.MustDequeue()
	tst.AssertEqual(t, "Queue.MustDequeue", frontItem, 2)
	front = q.Dequeue()
	tst.AssertEqual2(t, "Queue.Dequeue", front.IsNil(), true, front.Value(), 0)

	defer tst.AssertPanic(t, "Queue.MustDequeue")
	q.MustDequeue() // should panic
}
