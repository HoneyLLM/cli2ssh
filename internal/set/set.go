package set

type Set[T comparable] struct {
	items map[T]struct{}
}

func New[T comparable]() *Set[T] {
	return &Set[T]{items: make(map[T]struct{})}
}

// Adds an element to the set
func (s *Set[T]) Add(item T) {
	s.items[item] = struct{}{}
}

// Create a new set from a slice
func NewFromSlice[T comparable](items []T) *Set[T] {
	s := New[T]()
	for _, item := range items {
		s.Add(item)
	}
	return s
}

// Checks if an element exists in the set
func (s *Set[T]) Has(item T) bool {
	_, exists := s.items[item]
	return exists
}

// Removes an element from the set
func (s *Set[T]) Remove(item T) {
	delete(s.items, item)
}

// Returns the size of the set
func (s *Set[T]) Size() int {
	return len(s.items)
}

// Returns true if the set is empty
func (s *Set[T]) IsEmpty() bool {
	return s.Size() == 0
}
