package chia

import "sync"

// Signal is a signal chan idiom helper with idempotent closing.
type Signal struct {
	// C is a channel that will be closed on Signal.Close.
	C <-chan struct{}

	c chan<- struct{}

	once sync.Once
}

// NewSignal creates new Signal helper.
func NewSignal() *Signal {
	c := make(chan struct{})

	return &Signal{
		C: c, c: c,
	}
}

// Close is an idempotent closing of Signal channel.
func (s *Signal) Close() {
	s.once.Do(func() {
		close(s.c)
	})
}

// OnClose asynchronously runs f on closing c.
// It can be used both with Signal.C and Context.Done().
func OnClose(c <-chan struct{}, f func()) {
	go func() {
		<-c
		f()
	}()
}
