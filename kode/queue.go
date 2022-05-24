package kode

// Queue variable type
// A queue is a first-in, first-out data structure.
type Queue []interface{}

/**
 * Check if the queue is empty. If empty, return true.
 * @return bool - True if empty, false if not.
 */
func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}

/**
 * Push an item at the beginning of the queue.
 * @param item : interface{} - The item to push.
 */
func (q *Queue) Push(item interface{}) {
	// Add the item to the end of the array.
	*q = append(*q, item)
}

/**
 * Remove and return top element of queue. Return false if stack is empty.
 * @return interface{} - The top element of the stack.
 * @return bool - True if removed, false if not.
 */
func (q *Queue) Pop() (interface{}, bool) {
	if q.IsEmpty() {
		return nil, false
	} else {
		// Remove the first item from the queue
		item := (*q)[0]
		*q = (*q)[1:]
		return item, true
	}
}

/**
 * Get the top element of the queue. Return false if queue is empty.
 * @return interface{}, bool - The top element of the queue.
 */
func (q *Queue) Peek() (interface{}, bool) {
	if q.IsEmpty() {
		return nil, false
	} else {
		// Get the first item from the stack
		item := (*q)[0]
		return item, true
	}
}

/**
 * Remove all elements from the queue.
 */
func (q *Queue) Clear() {
	*q = (*q)[:0]
}
