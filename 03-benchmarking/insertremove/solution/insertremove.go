// Package insertremove contains the "Insert and Remove" exercise implementations
package insertremove

import (
	"container/list"
	"errors"
	"math/rand"
	"sort"
	"time"
)

// seed the pseudo random generator
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// declare an interface to ease testing ang benchmarking
type InserterRemover interface {
	Insert(n int) error
	Remove()

	// only for testing purposes
	Values() []int
}

// Slice implementation
type SliceImpl struct {
	values []int
}

func (s *SliceImpl) Insert(n int) error {
	if n < 0 {
		return errors.New("n must be a non negative value")
	}

	s.values = make([]int, n)

	// init slice with max int value to avoid checking for
	// sort.SearchInts returning out-of-bound index
	const maxInt = int(^uint(0) >> 1)
	for i := range s.values {
		s.values[i] = maxInt
	}

	for i := 0; i < len(s.values); i++ {
		cur := rand.Int()
		pos := sort.SearchInts(s.values, cur)
		copy(s.values[pos+1:], s.values[pos:])
		s.values[pos] = cur
	}

	return nil
}

func (s *SliceImpl) Remove() {
	for len(s.values) > 0 {
		pos := rand.Intn(len(s.values))
		s.values = append(s.values[:pos], s.values[pos+1:]...)
	}
}

func (s *SliceImpl) Values() []int {
	return s.values
}

func NewSliceImpl() *SliceImpl {
	return &SliceImpl{}
}

// List implementation
type ListImpl struct {
	list.List
}

func (l *ListImpl) Insert(n int) error {
	if n < 0 {
		return errors.New("n must be a non negative value")
	}

	// init list with guard values to avoid checking for first or last element
	first := l.List.PushFront(-1)
	last := l.List.PushBack(int(^uint(0) >> 1))

	for i := 0; i < n; i++ {
		cur := rand.Int()

		e := l.List.Front()
		for {
			e = e.Next()
			if cur <= e.Value.(int) {
				break
			}
		}

		l.List.InsertBefore(cur, e)
	}

	// remove first and last element
	l.List.Remove(first)
	l.List.Remove(last)

	return nil
}

func (l *ListImpl) Remove() {
	for l.List.Len() > 0 {
		pos := rand.Intn(l.List.Len())

		e := l.List.Front()
		for pos > 0 {
			e = e.Next()
			pos--
		}

		l.List.Remove(e)
	}
}

func (l *ListImpl) Values() []int {
	var values []int

	for e := l.List.Front(); e != nil; e = e.Next() {
		values = append(values, e.Value.(int))
	}

	return values
}

func NewListImpl() *ListImpl {
	return &ListImpl{}
}
