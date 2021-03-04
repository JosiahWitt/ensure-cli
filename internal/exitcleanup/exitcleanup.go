// Package exitcleanup provides helpers for running cleanups when the program is interrupted.
package exitcleanup

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// ExitCleanup supports registering functions to be called when the program is interrupted.
type ExitCleanup struct {
	log    *log.Logger
	signal chan os.Signal
	osExit func(int)

	cleanupFuncsMu sync.Mutex
	cleanupFuncs   []func() error
}

// New creates an ExitCleanup.
func New(log *log.Logger) *ExitCleanup {
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

		cleanup.cleanupFuncsMu.Lock()
		defer cleanup.cleanupFuncsMu.Unlock()
		for _, fn := range cleanup.cleanupFuncs {
			if err := fn(); err != nil {
				log.Printf(" -> Failed: %v", err)
			}
		}
		cleanup.osExit(1)
	}()

	return &cleanup
}

// Register a callback cleanup function.
func (cleanup *ExitCleanup) Register(fn func() error) {
	cleanup.cleanupFuncsMu.Lock()
	defer cleanup.cleanupFuncsMu.Unlock()
	cleanup.cleanupFuncs = append(cleanup.cleanupFuncs, fn)
}
