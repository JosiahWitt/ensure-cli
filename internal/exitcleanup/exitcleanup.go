// Package exitcleanup provides helpers for running cleanups when the program is interrupted.
package exitcleanup

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

type ExitCleaner interface {
	ToContext(ctx context.Context) context.Context
	Register(fn func() error)
}

// ExitCleanup supports registering functions to be called when the program is interrupted.
type ExitCleanup struct {
	log    *log.Logger
	signal chan os.Signal
	osExit func(int)

	terminatedCount int64 // Using int instead of bool, so we can leverage atomic
	cleanupFuncsMu  sync.Mutex
	cleanupFuncs    []func() error

	cancelFuncsMu sync.Mutex
	cancelFuncs   []context.CancelFunc
}

// New creates an ExitCleanup.
// It returns a cleanup function that should be called once right before exiting the app.
func New(log *log.Logger) (*ExitCleanup, func()) {
	cleanup := ExitCleanup{
		log:    log,
		signal: make(chan os.Signal, 1),
		osExit: os.Exit,

		cleanupFuncs: []func() error{},
	}

	signal.Notify(cleanup.signal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cleanup.signal
		log.Println()

		atomic.AddInt64(&cleanup.terminatedCount, 1)

		cleanup.cancelFuncsMu.Lock()
		defer cleanup.cancelFuncsMu.Unlock()
		for _, fn := range cleanup.cancelFuncs {
			fn()
		}
	}()

	return &cleanup, cleanup.cleanup
}

// ToContext adds the ability to cancel the context.
// It should be used to wrap any contexts that are passed to functions that leverage Register().
func (cleanup *ExitCleanup) ToContext(ctx context.Context) context.Context {
	cleanup.cleanupFuncsMu.Lock()
	defer cleanup.cleanupFuncsMu.Unlock()

	cancelCtx, cancelFunc := context.WithCancel(ctx)
	cleanup.cancelFuncs = append(cleanup.cancelFuncs, cancelFunc)
	return cancelCtx
}

// Register a callback cleanup function.
func (cleanup *ExitCleanup) Register(fn func() error) {
	cleanup.cleanupFuncsMu.Lock()
	defer cleanup.cleanupFuncsMu.Unlock()
	cleanup.cleanupFuncs = append(cleanup.cleanupFuncs, fn)
}

// Runs the cleanup functions if the app was terminated.
// It should be called once when the app is exited.
func (cleanup *ExitCleanup) cleanup() {
	if atomic.AddInt64(&cleanup.terminatedCount, 0) == 0 {
		return // Was not terminated, so don't run cleanup
	}

	cleanup.cleanupFuncsMu.Lock()
	defer cleanup.cleanupFuncsMu.Unlock()

	for _, fn := range cleanup.cleanupFuncs {
		if err := fn(); err != nil {
			cleanup.log.Printf(" -> Failed: %v", err)
		}
	}

	cleanup.osExit(1)
}
