//go:build tinygo
// +build tinygo

package metrics

import (
	"time"
)

// Capture new values for the Go garbage collector statistics exported in
// debug.GCStats.  This is designed to be called as a goroutine.
func CaptureDebugGCStats(r Registry, d time.Duration) {
}

// Capture new values for the Go garbage collector statistics exported in
// debug.GCStats.  This is designed to be called in a background goroutine.
// Giving a registry which has not been given to RegisterDebugGCStats will
// panic.
//
// Be careful (but much less so) with this because debug.ReadGCStats calls
// the C function runtime·lock(runtime·mheap) which, while not a stop-the-world
// operation, isn't something you want to be doing all the time.
func CaptureDebugGCStatsOnce(r Registry) {
}

// Register metrics for the Go garbage collector statistics exported in
// debug.GCStats.  The metrics are named by their fully-qualified Go symbols,
// i.e. debug.GCStats.PauseTotal.
func RegisterDebugGCStats(r Registry) {
}

// Allocate an initial slice for gcStats.Pause to avoid allocations during
// normal operation.
func init() {
}
