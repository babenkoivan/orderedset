// Package orderedset provides a generic set implementation that maintains insertion order.
// It supports typical set operations like add, remove, union, intersection, and difference,
// while preserving the order in which elements were added.
package orderedset

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

// OrderedSet is a generic set that preserves insertion order.
type OrderedSet[T comparable] struct {
	mu     sync.RWMutex
	index  map[T]struct{}
	values []T
}

// New creates a new empty OrderedSet.
func New[T comparable]() *OrderedSet[T] {
	return &OrderedSet[T]{
		index:  make(map[T]struct{}),
		values: make([]T, 0),
	}
}

// Add inserts a value into the set if it is not already present.
func (s *OrderedSet[T]) Add(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.index[value]; !exists {
		s.index[value] = struct{}{}
		s.values = append(s.values, value)
	}
}

// Remove deletes a value from the set.
func (s *OrderedSet[T]) Remove(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.index[value]; exists {
		delete(s.index, value)
		for i, v := range s.values {
			if v == value {
				s.values = append(s.values[:i], s.values[i+1:]...)
				break
			}
		}
	}
}

// RemoveAt deletes a value by index and returns it. If index is invalid, ok is false.
func (s *OrderedSet[T]) RemoveAt(index int) (val T, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.values) {
		var zero T
		return zero, false
	}
	val = s.values[index]
	s.values = append(s.values[:index], s.values[index+1:]...)
	delete(s.index, val)
	return val, true
}

// Has reports whether the set contains the given value.
func (s *OrderedSet[T]) Has(value T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.index[value]
	return exists
}

// Len returns the number of elements in the set.
func (s *OrderedSet[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.values)
}

// Values returns a copy of the values in insertion order.
func (s *OrderedSet[T]) Values() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	valuesCopy := make([]T, len(s.values))
	copy(valuesCopy, s.values)
	return valuesCopy
}

// At returns the element at the given index.
func (s *OrderedSet[T]) At(index int) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if index < 0 || index >= len(s.values) {
		var zero T
		return zero, false
	}
	return s.values[index], true
}

// IndexOf returns the index of the given value, or -1 if not found.
func (s *OrderedSet[T]) IndexOf(value T) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i, v := range s.values {
		if v == value {
			return i
		}
	}
	return -1
}

// SortBy sorts the elements of the set in-place using the provided less function.
func (s *OrderedSet[T]) SortBy(less func(a, b T) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sort.Slice(s.values, func(i, j int) bool {
		return less(s.values[i], s.values[j])
	})
}

// Clone returns a new copy of the set.
func (s *OrderedSet[T]) Clone() *OrderedSet[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	clone := New[T]()
	for _, v := range s.values {
		clone.index[v] = struct{}{}
		clone.values = append(clone.values, v)
	}
	return clone
}

// Union returns a new set containing all elements from both sets.
func (s *OrderedSet[T]) Union(other *OrderedSet[T]) *OrderedSet[T] {
	result := s.Clone()
	for _, v := range other.Values() {
		result.Add(v)
	}
	return result
}

// Intersect returns a new set with elements common to both sets.
func (s *OrderedSet[T]) Intersect(other *OrderedSet[T]) *OrderedSet[T] {
	result := New[T]()
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.values {
		if other.Has(v) {
			result.Add(v)
		}
	}
	return result
}

// Difference returns a new set with elements in s that are not in other.
func (s *OrderedSet[T]) Difference(other *OrderedSet[T]) *OrderedSet[T] {
	result := New[T]()
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.values {
		if !other.Has(v) {
			result.Add(v)
		}
	}
	return result
}

// Slice returns a new set containing elements from index "from" (inclusive) to "to" (exclusive).
// Returns an error if indices are out of range or invalid.
func (s *OrderedSet[T]) Slice(from, to int) (*OrderedSet[T], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if from < 0 {
		return nil, fmt.Errorf("from index %d is negative", from)
	}
	if to > len(s.values) {
		return nil, fmt.Errorf("to index %d is out of range", to)
	}
	if from > to {
		return nil, fmt.Errorf("from index %d is greater than to index %d", from, to)
	}

	result := New[T]()
	for _, v := range s.values[from:to] {
		result.Add(v)
	}
	return result, nil
}

// MarshalJSON implements json.Marshaler.
func (s *OrderedSet[T]) MarshalJSON() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s.values)
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *OrderedSet[T]) UnmarshalJSON(data []byte) error {
	var raw []T
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.index = make(map[T]struct{}, len(raw))
	s.values = make([]T, 0, len(raw))
	for _, v := range raw {
		if _, exists := s.index[v]; !exists {
			s.index[v] = struct{}{}
			s.values = append(s.values, v)
		}
	}
	return nil
}
