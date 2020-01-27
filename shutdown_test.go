package chia_test

import (
	"runtime"
	"testing"

	"github.com/pshch-pshch/chia"
)

func TestShutdown_Close(t *testing.T) {
	shutdown := chia.NewShutdown()

	runtime.Gosched()
	select {
	case <-shutdown.Init:
		t.Fatal("Can receive from new Shutdown Init channel")
	case <-shutdown.Done:
		t.Fatal("Can receive from new Shutdown Done channel")
	default:
	}

	shutdown.Close()

	runtime.Gosched()
	select {
	case <-shutdown.Done:
		t.Fatal("Can receive from Shutdown Done channel before shutdown completes")
	default:
	}

	done := <-shutdown.Init

	runtime.Gosched()
	select {
	case <-shutdown.Done:
		t.Fatal("Can receive from Shutdown Done channel before shutdown completes")
	default:
	}

	done()

	if _, ok := <-shutdown.Done; ok {
		t.Fatal("Shutdown Done channel is not closed")
	}
}

func TestShutdown_RepeatedClose(t *testing.T) {
	shutdown := chia.NewShutdown()
	defer shutdown.Terminate()

	shutdown.Close()

	<-shutdown.Init

	shutdown.Close()

	runtime.Gosched()
	select {
	case <-shutdown.Init:
		t.Fatal("Can receive twice from Shutdown Init channel")
	default:
	}
}

func TestShutdown_Terminate(t *testing.T) {
	shutdown := chia.NewShutdown()

	shutdown.Terminate()

	<-shutdown.Init

	if _, ok := <-shutdown.Done; ok {
		t.Fatal("Shutdown Done channel is not closed")
	}
}
