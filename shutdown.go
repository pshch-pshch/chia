package chia

import "sync"

// Shutdown is similar to Signal but separates closing initialization and completeness events.
type Shutdown struct {
	// Init channel will contain "done" function once the Close or Terminate will be called for the first time.
	// Call this function to close Done channel and indicate that shutdown completes.
	// Note that there will be exactly one receive from this channel for the whole Shutdown lifecycle.
	Init <-chan func()
	// Done is a channel that will be closed when shutdown completes.
	Done <-chan struct{}

	init chan<- func()
	done *Signal

	once sync.Once
}

// NewShutdown creates new Shutdown helper.
func NewShutdown() *Shutdown {
	init := make(chan func(), 1)
	done := NewSignal()

	return &Shutdown{
		Init: init, init: init,
		Done: done.C, done: done,
	}
}

// Close sends "done" function to Init channel if it wasn't sent before.
func (t *Shutdown) Close() {
	t.once.Do(func() {
		t.init <- t.done.Close
	})
}

// CloseAndWait sends "done" function to Init channel if it wasn't sent before and waits for shutdown completeness.
// Successive calls just blocks until shutdown completes.
func (t *Shutdown) CloseAndWait() {
	t.Close()
	<-t.done.C
}

// Terminate may be used to force shutdown complete by closing Done channel.
// It will also send "done" function to Init channel if it wasn't sent before.
func (t *Shutdown) Terminate() {
	t.Close()
	t.done.Close()
}
