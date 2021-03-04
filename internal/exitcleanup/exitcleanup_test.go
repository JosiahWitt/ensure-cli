package exitcleanup //nolint:testpackage // Internal to package, so we can test the signal channel without leaking it.

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/JosiahWitt/ensure"
)

func TestNew(t *testing.T) {
	ensure := ensure.New(t)

	var count int64 = 0
	var exitCode int64 = 0

	log := log.New(ioutil.Discard, "", 0)
	exitCleanup := New(log)
	exitCleanup.osExit = func(i int) { atomic.AddInt64(&exitCode, int64(i)) }

	exitCleanup.cleanupFuncs = []func() error{
		func() error {
			atomic.AddInt64(&count, 1)
			return nil
		},
		func() error {
			atomic.AddInt64(&count, 1)
			return errors.New("example error") // Error should not stop other cleanup from happening
		},
		func() error {
			atomic.AddInt64(&count, 1)
			return nil
		},
	}

	exitCleanup.signal <- os.Interrupt
	time.Sleep(1 * time.Millisecond) // Allow time for goroutine to run functions

	ensure(atomic.AddInt64(&count, 0)).Equals(int64(3))
	ensure(atomic.AddInt64(&exitCode, 0)).Equals(int64(1))
}

func TestRegister(t *testing.T) {
	ensure := ensure.New(t)

	log := log.New(ioutil.Discard, "", 0)
	exitCleanup := New(log)

	err := errors.New("example error")
	fn1 := func() error { return nil }
	fn2 := func() error { return err }

	exitCleanup.Register(fn1)
	exitCleanup.Register(fn1)

	ensure(len(exitCleanup.cleanupFuncs)).Equals(2)
	ensure(fn1()).IsNotError()
	ensure(fn2()).IsError(err)
}
