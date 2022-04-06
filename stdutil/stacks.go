package stdutil

import (
	"runtime"
	"strings"

	"github.com/gookit/goutil/strutil"
)

var (
	DefStackLen = 10000
	MaxStackLen = 100000
)

// GetCallStacks stacks is a wrapper for runtime.
// If all is true, Stack that attempts to recover the data for all goroutines.
//
// from glog package
func GetCallStacks(all bool) []byte {
	// We don't know how big the traces are, so grow a few times if they don't fit.
	// Start large, though.
	n := DefStackLen
	if all {
		n = MaxStackLen
	}

	// 4<<10 // 4 KB should be enough
	var trace []byte
	for i := 0; i < 10; i++ {
		trace = make([]byte, n)
		bts := runtime.Stack(trace, all)
		if bts < len(trace) {
			return trace[:bts]
		}
		n *= 2
	}
	return trace
}

// GetCallersInfo returns an array of strings containing
// the file and line number of each stack frame leading
func GetCallersInfo(skip, max int) (callers []string) {
	var (
		pc   uintptr
		ok   bool
		line int
		file string
		name string
	)

	// callers := []string{}
	for i := skip; i < max; i++ {
		pc, file, line, ok = runtime.Caller(i)
		if !ok {
			// The breaks below failed to terminate the loop, and we ran off the
			// end of the call stack.
			break
		}

		// This is a huge edge case, but it will panic if this is the case, see #180
		if file == "<autogenerated>" {
			break
		}

		f := runtime.FuncForPC(pc)
		if f == nil {
			break
		}

		name = f.Name()
		parts := strings.Split(file, "/")
		file = parts[len(parts)-1]
		if len(parts) > 1 {
			// dir := parts[len(parts)-2]
			callers = append(callers, name+" "+file+":"+strutil.MustString(line))
		}

		// Drop the package
		segments := strings.Split(name, ".")
		name = segments[len(segments)-1]
	}

	return callers
}
