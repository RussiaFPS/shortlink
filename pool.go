package shortlink

import "sync"

// Resetter defines the interface for objects that can be reset.
type Resetter interface {
	Reset()
}

// Pool is a generic pool for objects that implement the Resetter interface.
type Pool[T Resetter] struct {
	pool sync.Pool
}

// New creates a new Pool. The newFunc is used to create new objects when the pool is empty.
func New[T Resetter](newFunc func() any) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: newFunc,
		},
	}
}

// Get retrieves an object from the pool.
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put returns an object to the pool.
// It calls Reset() on the object before putting it back.
func (p *Pool[T]) Put(obj T) {
	obj.Reset()
	p.pool.Put(obj)
}
