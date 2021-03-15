package exitcleanup //nolint:testpackage // Internal to package, so we can test the signal channel without leaking it.

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg"
)

func TestExitCleanup(t *testing.T) {
	ensure := ensure.New(t)

	ensure.Run("with signal termination", func(ensure ensurepkg.Ensure) {
		var count int64 = 0
		var exitCode int64 = 0

		log := log.New(ioutil.Discard, "", 0)
		exitCleanup, cleanup := New(log)
		exitCleanup.osExit = func(i int) { atomic.AddInt64(&exitCode, int64(i)) }

		exitCleanup.Register(func() error {
			atomic.AddInt64(&count, 1)
			return nil
		})
		exitCleanup.Register(func() error {
			atomic.AddInt64(&count, 1)
			return errors.New("example error") // Error should not stop other cleanup from happening
		})
		exitCleanup.Register(func() error {
			time.Sleep(100 * time.Millisecond) // Simulate a larger cleanup effort to require a wait
			atomic.AddInt64(&count, 1)
			return nil
		})

		ctx1 := exitCleanup.ToContext(context.Background())
		ctx2 := exitCleanup.ToContext(context.Background())
		ensure(ctx1.Err()).IsNotError()
		ensure(ctx2.Err()).IsNotError()

		exitCleanup.signal <- os.Interrupt
		time.Sleep(10 * time.Millisecond) // Allow time for goroutine to mark as terminated

		ensure(ctx1.Err()).IsError(context.Canceled)
		ensure(ctx2.Err()).IsError(context.Canceled)
		ensure(atomic.AddInt64(&count, 0)).Equals(int64(0))
		ensure(atomic.AddInt64(&exitCode, 0)).Equals(int64(0))

		cleanup() // Trigger cleanup

		ensure(ctx1.Err()).IsError(context.Canceled)
		ensure(ctx2.Err()).IsError(context.Canceled)
		ensure(atomic.AddInt64(&count, 0)).Equals(int64(3))
		ensure(atomic.AddInt64(&exitCode, 0)).Equals(int64(1))
	})

	ensure.Run("with no termination", func(ensure ensurepkg.Ensure) {
		var count int64 = 0
		var exitCode int64 = 0

		log := log.New(ioutil.Discard, "", 0)
		exitCleanup, cleanup := New(log)
		exitCleanup.osExit = func(i int) { atomic.AddInt64(&exitCode, int64(i)) }

		exitCleanup.Register(func() error {
			atomic.AddInt64(&count, 1)
			return nil
		})
		exitCleanup.Register(func() error {
			atomic.AddInt64(&count, 1)
			return errors.New("example error") // Error should not stop other cleanup from happening
		})
		exitCleanup.Register(func() error {
			time.Sleep(100 * time.Millisecond) // Simulate a larger cleanup effort to require a wait
			atomic.AddInt64(&count, 1)
			return nil
		})

		ctx := exitCleanup.ToContext(context.Background())
		ensure(ctx.Err()).IsNotError()

		cleanup() // Trigger cleanup

		ensure(ctx.Err()).IsNotError()
		ensure(atomic.AddInt64(&count, 0)).Equals(int64(0))
		ensure(atomic.AddInt64(&exitCode, 0)).Equals(int64(0))
	})
}
