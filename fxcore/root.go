package fxcore

import (
	"path/filepath"
	"runtime"
)

// RootDir returns the root dir, for a provided number of stack frames to ascend.
//
//nolint:dogsled
func RootDir(skip int) string {
	_, file, _, _ := runtime.Caller(skip)

	return filepath.Join(filepath.Dir(file), "..")
}
