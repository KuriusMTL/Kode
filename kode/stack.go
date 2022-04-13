package kode

// Stack variable type
// A stack is a last-in, first-out data structure.
type Stack []interface{}

/**
 * Check if the stack is empty. If empty, return true.
 * @return bool - True if empty, false if not.
 */
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

/**
 * Push an item at the end of the stack.
 * @param item : interface{} - The item to push.
 */
func (s *Stack) Push(item interface{}) {
	// Add the item to the end of the array.
	*s = append(*s, item)
}

/**
 * Remove and return top element of stack. Return false if stack is empty.
 * @return interface{} - The top element of the stack.
 * @return interface{}, bool - True if removed, false if not.
 */
func (s *Stack) Pop() (interface{}, bool) {
	if s.IsEmpty() {
		return nil, false
	} else {
		// Remove the first item from the stack
		item := (*s)[len(*s)-1]
		*s = (*s)[:len(*s)-1]
		return item, true
	}
}

/**
 * Get the top element of the stack. Return false if stack is empty.
 * @return interface{}, bool - The top element of the stack.
 */
func (s *Stack) Peek() (interface{}, bool) {
	if s.IsEmpty() {
		return nil, false
	} else {
		// Get the first item from the stack
		item := (*s)[len(*s)-1]
		return item, true
	}
}

/**
 * Remove all elements from the stack.
 */
func (s *Stack) Clear() {
	*s = (*s)[:0]
}
