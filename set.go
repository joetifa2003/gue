package main

type set[T comparable] struct {
	values map[T]struct{}
}

func (s *set[T]) add(value T) {
	s.values[value] = struct{}{}
}

func (s *set[T]) remove(value T) {
	delete(s.values, value)
}

func (s *set[T]) forEach(fn func(T)) {
	for value := range s.values {
		fn(value)
	}
}

func newSet[T comparable]() set[T] {
	return set[T]{
		values: map[T]struct{}{},
	}
}
