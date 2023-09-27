//go:build !smoke_test

package client

// In full run, there is no maximumStep limit. We're running until derivation complete
var maximumSteps = -1
