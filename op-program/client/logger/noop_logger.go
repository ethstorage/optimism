//go:build js

package logger

import (
	"github.com/ethereum/go-ethereum/log"
)

func NewLogger() log.Logger {
	return nil
}

type NoopLog struct {
}

func (l NoopLog) New(ctx ...interface{}) log.Logger {
	return l
}

func (l NoopLog) GetHandler() log.Handler {
	// TODO: return nil handler
	return nil
}

func (l NoopLog) SetHandler(h log.Handler) {}

func (l NoopLog) Trace(msg string, ctx ...interface{}) {}

func (l NoopLog) Debug(msg string, ctx ...interface{}) {}

func (l NoopLog) Info(msg string, ctx ...interface{}) {}

func (l NoopLog) Warn(msg string, ctx ...interface{}) {}

func (l NoopLog) Error(msg string, ctx ...interface{}) {}

func (l NoopLog) Crit(msg string, ctx ...interface{}) {}
