//go:build js
// +build js

package logger

import (
	"sync/atomic"

	"github.com/ethereum/go-ethereum/log"
)

func NewLogger() log.Logger {
	println("from no logger")
	ctx := make([]interface{}, 2)
	return &NoopLog{ctx: ctx, h: new(swapHandler)}
}

// swapHandler wraps another handler that may be swapped out
// dynamically at runtime in a thread-safe fashion.
type swapHandler struct {
	handler atomic.Value
}

func (h *swapHandler) Log(r *log.Record) error {
	return (*h.handler.Load().(*log.Handler)).Log(r)
}

func (h *swapHandler) Swap(newHandler log.Handler) {
	h.handler.Store(&newHandler)
}

func (h *swapHandler) Get() log.Handler {
	return *h.handler.Load().(*log.Handler)
}

type NoopLog struct {
	ctx []interface{}
	h   *swapHandler
}

func (l *NoopLog) New(ctx ...interface{}) log.Logger {
	return &NoopLog{ctx: ctx, h: new(swapHandler)}
}

func (l *NoopLog) GetHandler() log.Handler {
	// TODO: return nil handler
	return l.h.Get()
}

func (l *NoopLog) SetHandler(h log.Handler) {}

func (l *NoopLog) Trace(msg string, ctx ...interface{}) {}

func (l *NoopLog) Debug(msg string, ctx ...interface{}) {}

func (l *NoopLog) Info(msg string, ctx ...interface{}) {}

func (l *NoopLog) Warn(msg string, ctx ...interface{}) {}

func (l *NoopLog) Error(msg string, ctx ...interface{}) {}

func (l *NoopLog) Crit(msg string, ctx ...interface{}) {}
