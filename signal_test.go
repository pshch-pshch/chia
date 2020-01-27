package chia_test

import (
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/pshch-pshch/chia"
)

func TestSignal_Close(t *testing.T) {
	signal := chia.NewSignal()

	runtime.Gosched()
	select {
	case <-signal.C:
		t.Fatal("Can receive from new Signal")
	default:
	}

	signal.Close()

	if _, ok := <-signal.C; ok {
		t.Fatal("Signal channel is not closed")
	}
}

func TestSignal_RepeatedClose(t *testing.T) {
	signal := chia.NewSignal()

	// Must not panic
	signal.Close()
	signal.Close()
}

func TestOnClose(t *testing.T) {
	signal := chia.NewSignal()

	var v int32 = 0

	chia.OnClose(signal.C, func() {
		atomic.StoreInt32(&v, 1)
	})

	runtime.Gosched()
	if val := atomic.LoadInt32(&v); val != 0 {
		t.Fatal("OnClose was called before Signal was closed")
	}

	signal.Close()

	runtime.Gosched()
	if val := atomic.LoadInt32(&v); val != 1 {
		t.Fatal("OnClose wasn't called after Signal was closed")
	}
}
